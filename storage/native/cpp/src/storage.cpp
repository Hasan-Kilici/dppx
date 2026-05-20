#include "storage.hpp"

#include <algorithm>
#include <cstdio>
#include <fstream>
#include <iostream>
#include <filesystem>

#include "distance.hpp"

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

    insertInternal(id, vector);
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
    if (config_.enableWAL) {
        if (wal_) {
            wal_->appendDelete(id);
        }
    }
}

std::vector<SearchResult> IndexEngine::search(const float* query, int dim, int k, int efSearch) const {
    std::vector<SearchResult> merged;
    if (!query || dim != config_.dimension || k <= 0) {
        return merged;
    }

    if (config_.useHNSW && hnsw_) {
        merged = hnsw_->search(query, k, efSearch > 0 ? efSearch : config_.efSearch);
    }

    for (const auto& segment : segments_) {
        const auto partial = segment.search(query, k);
        merged.insert(merged.end(), partial.begin(), partial.end());
    }

    std::sort(merged.begin(), merged.end(), [](const SearchResult& a, const SearchResult& b) {
        return a.score > b.score;
    });
    if (static_cast<int>(merged.size()) > k) {
        merged.resize(k);
    }
    return merged;
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

void IndexEngine::pruneStaleSegments() {
}
