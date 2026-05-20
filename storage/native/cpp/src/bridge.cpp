#include "bridge.h"
#include "storage.hpp"

#include <cstdlib>
#include <cstring>

extern "C" {

void* dppx_create_index(NativeConfig cfg) {
    StorageConfig config;
    config.dimension = cfg.dimension;
    config.path = cfg.path ? cfg.path : ".";
    config.hnswM = cfg.hnswM > 0 ? cfg.hnswM : 16;
    config.efConstruction = cfg.efConstruction > 0 ? cfg.efConstruction : 200;
    config.efSearch = cfg.efSearch > 0 ? cfg.efSearch : 50;
    config.maxSegmentSize = cfg.maxSegmentSize > 0 ? cfg.maxSegmentSize : (1 << 20);
    config.useHNSW = cfg.useHNSW != 0;
    config.enableWAL = cfg.enableWAL != 0;
    config.storePayload = cfg.storePayload != 0;

    return new IndexEngine(config);
}

void dppx_free_index(void* ptr) {
    delete static_cast<IndexEngine*>(ptr);
}

void dppx_insert(void* ptr, uint64_t id, const float* vector, int dim) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (engine) {
        engine->insert(id, vector);
    }
}

void dppx_insert_batch(void* ptr, const uint64_t* ids, const float* vectors, int count, int dim) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (engine) {
        engine->insertBatch(ids, vectors, count, dim);
    }
}

void dppx_delete(void* ptr, uint64_t id) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (engine) {
        engine->remove(id);
    }
}

void dppx_search(void* ptr, const float* query, int dim, int k, CSearchResult* out) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (!engine || !query || !out) {
        return;
    }

    auto results = engine->search(query, dim, k, 0);
    for (size_t i = 0; i < results.size(); ++i) {
        out[i].id = results[i].id;
        out[i].score = results[i].score;
    }
}

void dppx_search_options(void* ptr, const float* query, int dim, int k, int efSearch, CSearchResult* out) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (!engine || !query || !out) {
        return;
    }

    auto results = engine->search(query, dim, k, efSearch);
    for (size_t i = 0; i < results.size(); ++i) {
        out[i].id = results[i].id;
        out[i].score = results[i].score;
    }
}

void dppx_flush(void* index) {
    if (!index) {
        return;
    }

    auto* engine = static_cast<IndexEngine*>(index);

    engine->flush();
}

void dppx_recover(void* ptr) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (engine) {
        engine->recover();
    }
}

}