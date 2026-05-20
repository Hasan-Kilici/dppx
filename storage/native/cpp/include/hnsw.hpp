#pragma once

#include <cstdint>
#include <memory>
#include <random>
#include <unordered_map>
#include <vector>

#include "heap.hpp"

struct HNSWNode {
    uint64_t id;
    std::vector<float> vector;
    int level;
    std::vector<std::vector<uint64_t>> neighbors;
};

class HNSWIndex {
public:
    HNSWIndex(int dimension, int maxConnections, int efConstruction, int efSearch);

    void insert(uint64_t id, const float* vector);
    std::vector<SearchResult> search(const float* query, int k, int efSearch) const;

private:
    int randomLevel();
    std::vector<uint64_t> searchLayer(const float* query, int ef, int level) const;
    std::vector<uint64_t> selectNeighbors(const std::vector<uint64_t>& candidates, const float* query, int m) const;
    float similarity(const float* a, const float* b) const;

    int dimension_;
    int M_;
    int efConstruction_;
    int efSearch_;
    bool hasEntryPoint_;
    uint64_t entryPointId_;
    std::vector<HNSWNode> nodes_;
    std::unordered_map<uint64_t, size_t> idMap_;
    mutable std::mt19937 generator_;
    std::uniform_real_distribution<float> distribution_;
};