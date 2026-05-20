#pragma once

#include <cstdint>
#include <memory>
#include <string>
#include <vector>

#include "hnsw.hpp"
#include "segment.hpp"
#include "wal.hpp"
#include "payload.hpp"

struct SearchResult;

struct StorageConfig {
    int dimension = 0;
    std::string path;
    bool useHNSW = false;
    int hnswM = 16;
    int efConstruction = 200;
    int efSearch = 50;
    int maxSegmentSize = 1 << 20;
    bool enableWAL = true;
    bool storePayload = false;
};

class IndexEngine {
public:
    explicit IndexEngine(const StorageConfig& config);
    ~IndexEngine();

    void insert(uint64_t id, const float* vector);
    void insertBatch(const uint64_t* ids, const float* vectors, int count, int dim);
    void remove(uint64_t id);
    std::vector<SearchResult> search(const float* query, int dim, int k, int efSearch) const;
    void flush();
    void recover();

private:
    void insertInternal(uint64_t id, const float* vector);
    void ensureStoragePath() const;
    void loadSegments();
    void writeWalInsert(uint64_t id, const float* vector);
    void pruneStaleSegments();

    StorageConfig config_;
    std::unique_ptr<HNSWIndex> hnsw_;
    std::vector<Segment> segments_;
    std::vector<PayloadDocument> payloads_;
    std::unique_ptr<WalWriter> wal_;
    std::vector<uint64_t> idIndex_;
};
