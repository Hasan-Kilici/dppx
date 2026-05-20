#pragma once

#include <cstdint>
#include <string>
#include <vector>
#include <algorithm>

#include "vector.hpp"
#include "payload.hpp"
#include "heap.hpp"

static constexpr uint32_t DPPX_SEGMENT_MAGIC = 0x44505058; // 'DPPX'
static constexpr uint32_t DPPX_SEGMENT_VERSION = 1;

struct SegmentHeader {
    uint32_t magic;
    uint32_t version;
    uint32_t dimension;
    uint64_t count;
};

struct SegmentNode {
    uint64_t id;
    std::vector<float> vector;
    uint64_t payloadOffset;
};

class Segment {
public:
    Segment(const std::string& basePath, uint32_t segmentId, int dimension);

    void append(uint64_t id, const float* vector, const PayloadDocument* payload = nullptr);
    bool load();
    void persist() const;
    std::vector<SearchResult> search(const float* query, int k) const;
    uint64_t size() const;
    const std::vector<SegmentNode>& nodes() const;

private:
    std::string filename() const;
    std::string payloadFilename() const;
    void writeHeader(std::FILE* out) const;

    std::string basePath_;
    uint32_t segmentId_;
    int dimension_;
    std::vector<SegmentNode> nodes_;
    std::vector<uint8_t> payloadBuffer_;
};
