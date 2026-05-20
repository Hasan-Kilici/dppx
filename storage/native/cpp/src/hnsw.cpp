#include "hnsw.hpp"
#include "distance.hpp"

#include <algorithm>

HNSWIndex::HNSWIndex(int dimension, int maxConnections, int efConstruction, int efSearch)
    : dimension_(dimension),
      M_(maxConnections),
      efConstruction_(efConstruction),
      efSearch_(efSearch),
      hasEntryPoint_(false),
      entryPointId_(0),
      generator_(std::random_device{}()),
      distribution_(0.0f, 1.0f) {
}

void HNSWIndex::insert(uint64_t id, const float* vector) {
    HNSWNode node;
    node.id = id;
    node.level = randomLevel();
    node.vector.assign(vector, vector + dimension_);
    node.neighbors.resize(node.level + 1);

    const size_t insertionIndex = nodes_.size();
    nodes_.push_back(std::move(node));
    idMap_[id] = insertionIndex;

    if (!hasEntryPoint_) {
        entryPointId_ = id;
        hasEntryPoint_ = true;
        return;
    }

    if (idMap_.find(entryPointId_) == idMap_.end()) {
        entryPointId_ = id;
        return;
    }

    auto currentLayerResults = searchLayer(vector, efConstruction_, node.level);
    std::vector<uint64_t> neighbors = selectNeighbors(currentLayerResults, vector, M_);

    HNSWNode& newNode = nodes_[insertionIndex];
    for (uint64_t neighborId : neighbors) {
        auto it = idMap_.find(neighborId);
        if (it == idMap_.end()) {
            continue;
        }
        size_t idx = it->second;
        newNode.neighbors[0].push_back(neighborId);
        nodes_[idx].neighbors[0].push_back(id);
    }
}

std::vector<SearchResult> HNSWIndex::search(const float* query, int k, int efSearch) const {
    if (!hasEntryPoint_ || k <= 0) {
        return {};
    }

    std::vector<SearchResult> candidates;
    auto entryIt = idMap_.find(entryPointId_);
    if (entryIt == idMap_.end()) {
        return {};
    }

    const HNSWNode& entry = nodes_[entryIt->second];
    float entryScore = similarity(query, entry.vector.data());
    candidates.push_back({entry.id, entryScore});

    std::vector<uint64_t> searchCandidates = searchLayer(query, efSearch, 0);
    for (uint64_t id : searchCandidates) {
        auto it = idMap_.find(id);
        if (it != idMap_.end()) {
            const HNSWNode& node = nodes_[it->second];
            candidates.push_back({node.id, similarity(query, node.vector.data())});
        }
    }

    std::sort(candidates.begin(), candidates.end(), [](const SearchResult& a, const SearchResult& b) {
        return a.score > b.score;
    });

    if (static_cast<int>(candidates.size()) > k) {
        candidates.resize(k);
    }

    return candidates;
}

int HNSWIndex::randomLevel() {
    float r = distribution_(generator_);
    int level = 0;
    while (r < 0.5f && level < 10) {
        level++;
        r = distribution_(generator_);
    }
    return level;
}

std::vector<uint64_t> HNSWIndex::searchLayer(const float* query, int ef, int level) const {
    std::vector<SearchResult> visited;
    visited.reserve(ef);

    for (const auto& node : nodes_) {
        if (node.level >= level) {
            visited.push_back({node.id, similarity(query, node.vector.data())});
        }
    }

    std::sort(visited.begin(), visited.end(), [](const SearchResult& a, const SearchResult& b) {
        return a.score > b.score;
    });

    if (static_cast<int>(visited.size()) > ef) {
        visited.resize(ef);
    }

    std::vector<uint64_t> result;
    result.reserve(visited.size());

    for (const auto& entry : visited) {
        result.push_back(entry.id);
    }
    return result;
}

std::vector<uint64_t> HNSWIndex::selectNeighbors(const std::vector<uint64_t>& candidates, const float* query, int m) const {
    std::vector<SearchResult> scored;
    scored.reserve(candidates.size());

    for (uint64_t candidateId : candidates) {
        auto it = idMap_.find(candidateId);
        if (it == idMap_.end()) {
            continue;
        }
        const HNSWNode& node = nodes_[it->second];
        scored.push_back({candidateId, similarity(query, node.vector.data())});
    }

    std::sort(scored.begin(), scored.end(), [](const SearchResult& a, const SearchResult& b) {
        return a.score > b.score;
    });

    std::vector<uint64_t> result;
    for (int i = 0; i < static_cast<int>(scored.size()) && i < m; ++i) {
        result.push_back(scored[i].id);
    }
    return result;
}

float HNSWIndex::similarity(const float* a, const float* b) const {
    return cosine_similarity(a, b, dimension_);
}
