#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct CSearchResult {
    uint64_t id;
    float score;
} CSearchResult;

typedef struct {
    int dimension;
    const char* path;
    int hnswM;
    int efConstruction;
    int efSearch;
    int maxSegmentSize;
    int useHNSW;
    int enableWAL;
    int storePayload;
} NativeConfig;

void* dppx_create_index(NativeConfig cfg);
void dppx_free_index(void* ptr);
void dppx_insert(void* ptr, uint64_t id, const float* vector, int dim);
void dppx_insert_batch(void* ptr, const uint64_t* ids, const float* vectors, int count, int dim);
void dppx_delete(void* ptr, uint64_t id);
int dppx_search(void* ptr, const float* query, int dim, int k, CSearchResult* out);
int dppx_search_options(void* ptr, const float* query, int dim, int k, int efSearch, CSearchResult* out);
typedef void* CSearchCursor;
CSearchCursor dppx_search_cursor(void* ptr, const float* query, int dim, int k);
CSearchCursor dppx_search_cursor_options(void* ptr, const float* query, int dim, int k, int efSearch);
int dppx_search_cursor_next(CSearchCursor cursor, CSearchResult* out);
void dppx_search_cursor_free(CSearchCursor cursor);
void dppx_flush(void* index);
void dppx_recover(void* ptr);

#ifdef __cplusplus
}
#endif
