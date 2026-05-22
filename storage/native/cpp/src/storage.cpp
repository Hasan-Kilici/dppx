#include "storage.hpp"

#include <algorithm>
#include <cstdio>
#include <fstream>
#include <iostream>
#include <filesystem>
#include <unordered_map>

#include "distance.hpp"

namespace {

class CandidateStream {
public:
    virtual ~CandidateStream() = default;
    virtual bool next(SearchResult& out) = 0;
};

class HNSWSearchStream : public CandidateStream {
public:
    HNSWSearchStream(const HNSWIndex* index, const float* query, int k, int efSearch)
        : position_(0) {
        if (index) {
            results_ = index->search(query, k, efSearch);
        }
    }

    bool next(SearchResult& out) override {
        if (position_ >= results_.size()) {
            return false;
        }
        out = results_[position_++];
        return true;
    }

private:
    std::vector<SearchResult> results_;
    size_t position_;
};

class SegmentSearchStream : public CandidateStream {
public:
    SegmentSearchStream(const Segment* segment, const float* query, int k)
        : position_(0) {
        if (segment) {
            results_ = segment->search(query, k);
        }
    }

    bool next(SearchResult& out) override {
        if (position_ >= results_.size()) {
            return false;
        }
        out = results_[position_++];
        return true;
    }

private:
    std::vector<SearchResult> results_;
    size_t position_;
};

std::vector<SearchResult> mergeTopK(const std::unordered_map<uint64_t, float>& bestScores, int k) {
    FixedMinHeap heap(k);
    for (const auto& entry : bestScores) {
        heap.Add({entry.first, entry.second});
    }
    return heap.Result();
}

} // namespace

IndexEngine::IndexEngine(const StorageConfig& config)
    : config_(config) {
    ensureStoragePath();

    if (config_.enableWAL) {
        wal_ = std::make_unique<WalWriter>(config_.path + "/wal/wal.log");
    }

    if (config_.useHNSW) {
        hnsw_ = std::make_unique<HNSWIndex>(
            config_.dimension,
            config_.hnswM,
            config_.efConstruction,
            config_.efSearch);
    }

    loadSegments();
    recover();
}

IndexEngine::~IndexEngine() {
    flush();
}

void IndexEngine::ensureStoragePath() const {
    std::filesystem::create_directories(config_.path);
    std::filesystem::create_directories(config_.path + "/segments");
    std::filesystem::create_directories(config_.path + "/wal");
    if (config_.storePayload) {
        std::filesystem::create_directories(config_.path + "/payload");
    }
}

void IndexEngine::loadSegments() {
    for (const auto& entry : std::filesystem::directory_iterator(config_.path + "/segments")) {
        if (entry.is_regular_file()) {
            std::string name = entry.path().filename().string();
            if (name.rfind("segment-", 0) == 0 && name.find(".seg") != std::string::npos) {
                uint32_t segmentId = 0;
                try {
                    segmentId = std::stoul(name.substr(8, name.find(".seg") - 8));
                } catch (...) {
                    continue;
                }

                Segment segment(config_.path, segmentId, config_.dimension);
                if (segment.load()) {
                    segments_.push_back(std::move(segment));
                }
            }
        }
    }
}

void IndexEngine::recover() {
    if (!wal_) {
        return;
    }

    std::ifstream reader(config_.path + "/wal/wal.log", std::ios::binary);
    if (!reader.is_open()) {
        return;
    }

    while (reader.good()) {
        WalRecordHeader header;
        reader.read(reinterpret_cast<char*>(&header), sizeof(header));
        if (reader.gcount() != sizeof(header)) {
            break;
        }

        if (header.op == static_cast<uint8_t>(WalOp::Insert)) {
            uint64_t id;
            reader.read(reinterpret_cast<char*>(&id), sizeof(id));
            std::vector<float> vector(config_.dimension);
            reader.read(reinterpret_cast<char*>(vector.data()), sizeof(float) * config_.dimension);
            insertInternal(id, vector.data());
        } else if (header.op == static_cast<uint8_t>(WalOp::Delete)) {
            uint64_t id;
            reader.read(reinterpret_cast<char*>(&id), sizeof(id));
            remove(id);
        } else {
            reader.seekg(header.length, std::ios::cur);
        }
    }
}


void IndexEngine::insertInternal(uint64_t id, const float* vector) {
    if (segments_.empty()) {
        segments_.emplace_back(config_.path, 0, config_.dimension);
    }

    Segment& active = segments_.back();
    active.append(id, vector, nullptr);

    if (active.size() >= static_cast<uint64_t>(config_.maxSegmentSize)) {
        active.persist();
        segments_.emplace_back(config_.path, static_cast<uint32_t>(segments_.size()), config_.dimension);
    }
}

void IndexEngine::insertBatch(const uint64_t* ids, const float* vectors, int count, int dim) {
    if (!ids || !vectors || count <= 0 || dim != config_.dimension) {
        return;
    }

    for (int i = 0; i < count; i++) {
        insert(ids[i], vectors + static_cast<size_t>(i) * dim);
    }
}

void IndexEngine::remove(uint64_t id) {
    if (wal_) {
        wal_->appendDelete(id);
    }
    deletedIds_.insert(id);
}

std::vector<SearchResult> IndexEngine::search(const SearchRequest& request) const {
    std::vector<SearchResult> merged;
    if (!request.query || request.dim != config_.dimension || request.k <= 0) {
        return merged;
    }

    std::unordered_map<uint64_t, float> uniqueScores;
    std::vector<std::unique_ptr<CandidateStream>> streams;

    if (request.includeHNSW && config_.useHNSW && hnsw_) {
        streams.push_back(std::make_unique<HNSWSearchStream>(
            hnsw_.get(),
            request.query,
            request.k,
            request.efSearch > 0 ? request.efSearch : config_.efSearch
        ));
    }

    if (request.includeSegments) {
        for (const auto& segment : segments_) {
            streams.push_back(std::make_unique<SegmentSearchStream>(
                &segment,
                request.query,
                request.k
            ));
        }
    }

    for (const auto& stream : streams) {
        SearchResult result;
        while (stream->next(result)) {
            if (deletedIds_.count(result.id)) {
                continue;
            }
            auto it = uniqueScores.find(result.id);
            if (it == uniqueScores.end() || result.score > it->second) {
                uniqueScores[result.id] = result.score;
            }
        }
    }

    return mergeTopK(uniqueScores, request.k);
}

std::vector<SearchResult> IndexEngine::search(const float* query, int dim, int k, int efSearch) const {
    SearchRequest request;
    request.query = query;
    request.dim = dim;
    request.k = k;
    request.efSearch = efSearch;
    return search(request);
}

void IndexEngine::flush() {
    for (auto& segment : segments_) {
        if (segment.size() > 0) {
            segment.persist();
        }
    }
    if (wal_) {
        wal_->flush();
    }
}

void IndexEngine::writeWalInsert(uint64_t id, const float* vector) {
    if (wal_) {
        wal_->appendInsert(id, vector, config_.dimension);
    }
}

void IndexEngine::insert(uint64_t id, const float* vector) {
    if (config_.dimension <= 0 || vector == nullptr) {
        return;
    }

    if (config_.enableWAL) {
        writeWalInsert(id, vector);
    }

    if (config_.useHNSW && hnsw_) {
        hnsw_->insert(id, vector);
    }

    deletedIds_.erase(id);
    insertInternal(id, vector);
}

void IndexEngine::pruneStaleSegments() {
}
