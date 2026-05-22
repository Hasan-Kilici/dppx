#include "bridge.h"
#include "storage.hpp"

#include <cstdlib>
#include <cstring>

extern "C" {

class SearchCursor {
public:
    explicit SearchCursor(std::vector<SearchResult>&& results)
        : results_(std::move(results)), position_(0) {
    }

    bool next(CSearchResult* out) {
        if (!out || position_ >= results_.size()) {
            return false;
        }
        out->id = results_[position_].id;
        out->score = results_[position_].score;
        position_++;
        return true;
    }

private:
    std::vector<SearchResult> results_;
    size_t position_;
};

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

int dppx_search(void* ptr, const float* query, int dim, int k, CSearchResult* out) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (!engine || !query || !out || k <= 0) {
        return 0;
    }

    auto results = engine->search(query, dim, k, 0);
    const size_t count = std::min(results.size(), static_cast<size_t>(k));
    for (size_t i = 0; i < count; ++i) {
        out[i].id = results[i].id;
        out[i].score = results[i].score;
    }
    return static_cast<int>(count);
}

int dppx_search_options(void* ptr, const float* query, int dim, int k, int efSearch, CSearchResult* out) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (!engine || !query || !out || k <= 0) {
        return 0;
    }

    auto results = engine->search(query, dim, k, efSearch);
    const size_t count = std::min(results.size(), static_cast<size_t>(k));
    for (size_t i = 0; i < count; ++i) {
        out[i].id = results[i].id;
        out[i].score = results[i].score;
    }
    return static_cast<int>(count);
}

CSearchCursor dppx_search_cursor(void* ptr, const float* query, int dim, int k) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (!engine || !query || k <= 0) {
        return nullptr;
    }

    auto results = engine->search(query, dim, k, 0);
    return new SearchCursor(std::move(results));
}

CSearchCursor dppx_search_cursor_options(void* ptr, const float* query, int dim, int k, int efSearch) {
    auto* engine = static_cast<IndexEngine*>(ptr);
    if (!engine || !query || k <= 0) {
        return nullptr;
    }

    auto results = engine->search(query, dim, k, efSearch);
    return new SearchCursor(std::move(results));
}

int dppx_search_cursor_next(CSearchCursor cursor, CSearchResult* out) {
    auto* c = static_cast<SearchCursor*>(cursor);
    if (!c || !out) {
        return 0;
    }
    return c->next(out) ? 1 : 0;
}

void dppx_search_cursor_free(CSearchCursor cursor) {
    delete static_cast<SearchCursor*>(cursor);
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