#include "distance.hpp"

#include <cmath>

#if defined(__AVX2__)
#include <immintrin.h>
#elif defined(__SSE__)
#include <xmmintrin.h>
#endif

float dot_product(const float* a, const float* b, int dim) {
    float total = 0.0f;
#if defined(__AVX2__)
    __m256 acc = _mm256_setzero_ps();
    int i = 0;
    for (; i + 7 < dim; i += 8) {
        __m256 va = _mm256_loadu_ps(a + i);
        __m256 vb = _mm256_loadu_ps(b + i);
        acc = _mm256_fmadd_ps(va, vb, acc);
    }
    alignas(32) float accum[8];
    _mm256_store_ps(accum, acc);
    for (int j = 0; j < 8; ++j) {
        total += accum[j];
    }
    for (; i < dim; ++i) {
        total += a[i] * b[i];
    }
#elif defined(__SSE__)
    __m128 acc = _mm_setzero_ps();
    int i = 0;
    for (; i + 3 < dim; i += 4) {
        __m128 va = _mm_loadu_ps(a + i);
        __m128 vb = _mm_loadu_ps(b + i);
        acc = _mm_add_ps(acc, _mm_mul_ps(va, vb));
    }
    alignas(16) float accum[4];
    _mm_store_ps(accum, acc);
    for (int j = 0; j < 4; ++j) {
        total += accum[j];
    }
    for (; i < dim; ++i) {
        total += a[i] * b[i];
    }
#else
    for (int i = 0; i < dim; i++) {
        total += a[i] * b[i];
    }
#endif
    return total;
}

float cosine_similarity(const float* a, const float* b, int dim) {
    float dot = dot_product(a, b, dim);
    float normA = dot_product(a, a, dim);
    float normB = dot_product(b, b, dim);

    if (normA == 0.0f || normB == 0.0f) {
        return 0.0f;
    }

    return dot / (std::sqrt(normA) * std::sqrt(normB));
}

float euclidean_distance(const float* a, const float* b, int dim) {
    float total = 0.0f;
#if defined(__AVX2__)
    __m256 acc = _mm256_setzero_ps();
    int i = 0;
    for (; i + 7 < dim; i += 8) {
        __m256 va = _mm256_loadu_ps(a + i);
        __m256 vb = _mm256_loadu_ps(b + i);
        __m256 diff = _mm256_sub_ps(va, vb);
        acc = _mm256_add_ps(acc, _mm256_mul_ps(diff, diff));
    }
    alignas(32) float accum[8];
    _mm256_store_ps(accum, acc);
    for (int j = 0; j < 8; ++j) {
        total += accum[j];
    }
    for (; i < dim; ++i) {
        float d = a[i] - b[i];
        total += d * d;
    }
#elif defined(__SSE__)
    __m128 acc = _mm_setzero_ps();
    int i = 0;
    for (; i + 3 < dim; i += 4) {
        __m128 va = _mm_loadu_ps(a + i);
        __m128 vb = _mm_loadu_ps(b + i);
        __m128 diff = _mm_sub_ps(va, vb);
        acc = _mm_add_ps(acc, _mm_mul_ps(diff, diff));
    }
    alignas(16) float accum[4];
    _mm_store_ps(accum, acc);
    for (int j = 0; j < 4; ++j) {
        total += accum[j];
    }
    for (; i < dim; ++i) {
        float d = a[i] - b[i];
        total += d * d;
    }
#else
    for (int i = 0; i < dim; i++) {
        float d = a[i] - b[i];
        total += d * d;
    }
#endif
    return std::sqrt(total);
}

float manhattan_distance(const float* a, const float* b, int dim) {
    float total = 0.0f;
    for (int i = 0; i < dim; i++) {
        float diff = a[i] - b[i];
        total += diff < 0 ? -diff : diff;
    }
    return total;
}
