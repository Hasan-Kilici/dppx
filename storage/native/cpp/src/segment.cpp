#include "segment.hpp"

#include <cstdio>
#include <cstring>

#include "distance.hpp"

Segment::Segment(const std::string& basePath, uint32_t segmentId, int dimension)
    : basePath_(basePath),
      segmentId_(segmentId),
      dimension_(dimension) {
}

void Segment::append(uint64_t id, const float* vector, const PayloadDocument* payload) {
    SegmentNode node;
    node.id = id;
    node.payloadOffset = 0;
    node.vector.assign(vector, vector + dimension_);
    nodes_.push_back(std::move(node));

    if (payload != nullptr) {
        auto encoded = serializePayload(*payload);
        payloadBuffer_.insert(payloadBuffer_.end(), encoded.begin(), encoded.end());
    }
}

bool Segment::load() {
    const std::string file = filename();
    FILE* in = std::fopen(file.c_str(), "rb");
    if (!in) {
        return false;
    }

    SegmentHeader header;
    if (std::fread(&header, sizeof(header), 1, in) != 1) {
        std::fclose(in);
        return false;
    }

    if (header.magic != DPPX_SEGMENT_MAGIC || header.version != DPPX_SEGMENT_VERSION || header.dimension != static_cast<uint32_t>(dimension_)) {
        std::fclose(in);
        return false;
    }

    nodes_.reserve(static_cast<size_t>(header.count));
    for (uint64_t i = 0; i < header.count; i++) {
        SegmentNode node;
        if (std::fread(&node.id, sizeof(node.id), 1, in) != 1) {
            std::fclose(in);
            return false;
        }

        node.vector.resize(dimension_);
        if (std::fread(node.vector.data(), sizeof(float), dimension_, in) != static_cast<size_t>(dimension_)) {
            std::fclose(in);
            return false;
        }

        node.payloadOffset = 0;
        nodes_.push_back(std::move(node));
    }

    std::fclose(in);
    return true;
}

void Segment::persist() const {
    const std::string file = filename();
    FILE* out = std::fopen(file.c_str(), "wb");
    if (!out) {
        return;
    }

    writeHeader(out);
    for (const auto& node : nodes_) {
        std::fwrite(&node.id, sizeof(node.id), 1, out);
        std::fwrite(node.vector.data(), sizeof(float), dimension_, out);
    }

    std::fclose(out);
}

std::vector<SearchResult> Segment::search(const float* query, int k) const {
    std::vector<SearchResult> results;
    if (!query || k <= 0) {
        return results;
    }

    results.reserve(std::min<size_t>(nodes_.size(), static_cast<size_t>(k)));
    for (const auto& node : nodes_) {
        const float score = cosine_similarity(query, node.vector.data(), dimension_);
        if (static_cast<int>(results.size()) < k) {
            results.push_back({node.id, score});
            if (static_cast<int>(results.size()) == k) {
                std::sort(results.begin(), results.end(), [](const SearchResult& a, const SearchResult& b) {
                    return a.score > b.score;
                });
            }
            continue;
        }

        if (score <= results.back().score) {
            continue;
        }

        results.back() = {node.id, score};
        std::sort(results.begin(), results.end(), [](const SearchResult& a, const SearchResult& b) {
            return a.score > b.score;
        });
    }
    return results;
}

uint64_t Segment::size() const {
    return nodes_.size();
}

const std::vector<SegmentNode>& Segment::nodes() const {
    return nodes_;
}

std::string Segment::filename() const {
    return basePath_ + "/segments/segment-" + std::to_string(segmentId_) + ".seg";
}

std::string Segment::payloadFilename() const {
    return basePath_ + "/payload/payload-" + std::to_string(segmentId_) + ".bin";
}

void Segment::writeHeader(std::FILE* out) const {
    SegmentHeader header;
    header.magic = DPPX_SEGMENT_MAGIC;
    header.version = DPPX_SEGMENT_VERSION;
    header.dimension = static_cast<uint32_t>(dimension_);
    header.count = nodes_.size();
    std::fwrite(&header, sizeof(header), 1, out);
}
