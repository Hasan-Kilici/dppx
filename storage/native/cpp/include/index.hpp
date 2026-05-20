#pragma once

#include "vector.hpp"
#include "heap.hpp"

#include <vector>

class FlatIndex {
public:

    explicit FlatIndex(int dim);

    void insert(
        uint64_t id,
        const float* vector
    );

    std::vector<SearchResult> search(
        const float* query,
        int k
    );

    void remove(
        uint64_t id
    );

private:

    int dimension;

    std::vector<VectorNode> items;
};