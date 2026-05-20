#include "index.hpp"
#include "distance.hpp"

#include <queue>
#include <algorithm>

FlatIndex::FlatIndex(int dim)
    : dimension(dim) {}

void FlatIndex::insert(
    uint64_t id,
    const float* vector
) {
    VectorNode node;

    node.id = id;

    node.vector.assign(
        vector,
        vector + dimension
    );

    items.push_back(node);
}

std::vector<SearchResult> FlatIndex::search(
    const float* query,
    int k
) {
    std::priority_queue<
        SearchResult,
        std::vector<SearchResult>,
        std::greater<SearchResult>
    > heap;

    for (const auto& item : items) {

        float sim = cosine_similarity(
            query,
            item.vector.data(),
            dimension
        );

        if ((int)heap.size() < k) {
            heap.push({item.id, sim});
        }
        else if (sim > heap.top().score) {
            heap.pop();
            heap.push({item.id, sim});
        }
    }

    std::vector<SearchResult> results;

    while (!heap.empty()) {
        results.push_back(heap.top());
        heap.pop();
    }

    std::reverse(
        results.begin(),
        results.end()
    );

    return results;
}

void FlatIndex::remove(
    uint64_t id
) {
    items.erase(
        std::remove_if(
            items.begin(),
            items.end(),
            [id](const VectorNode& node) {
                return node.id == id;
            }
        ),
        items.end()
    );
}