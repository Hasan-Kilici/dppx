#pragma once

#include <algorithm>
#include <queue>
#include <vector>
#include <cstdint>

struct SearchResult {
    uint64_t id;
    float score;

    bool operator>(const SearchResult& other) const {
        return score > other.score;
    }
};

struct MinScoreCompare {
    bool operator()(const SearchResult& a, const SearchResult& b) const {
        return a.score > b.score;
    }
};

class FixedMinHeap {
public:
    explicit FixedMinHeap(int capacity)
        : capacity_(capacity) {
    }

    void Add(const SearchResult& item) {
        if (static_cast<int>(heap_.size()) < capacity_) {
            heap_.push(item);
            return;
        }

        if (item.score > heap_.top().score) {
            heap_.pop();
            heap_.push(item);
        }
    }

    std::vector<SearchResult> Result() {
        std::vector<SearchResult> results;
        results.reserve(heap_.size());
        while (!heap_.empty()) {
            results.push_back(heap_.top());
            heap_.pop();
        }
        std::reverse(results.begin(), results.end());
        return results;
    }

private:
    int capacity_;
    std::priority_queue<SearchResult, std::vector<SearchResult>, MinScoreCompare> heap_;
};