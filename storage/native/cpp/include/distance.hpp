#pragma once

#include <cstdint>

float dot_product(
    const float* a,
    const float* b,
    int dim
);

float cosine_similarity(
    const float* a,
    const float* b,
    int dim
);

float euclidean_distance(
    const float* a,
    const float* b,
    int dim
);

float manhattan_distance(
    const float* a,
    const float* b,
    int dim
);
