# DPPX

A Go library for diversity-aware ranking and recommendation, inspired by Determinantal Point Processes (DPP).

## What is DPPX? / DPPX Nedir?

DPPX is a developer-focused library for ranking and selecting items with a strong emphasis on relevance and diversity. It is built around a configurable search engine that combines similarity scoring with custom business scoring, using a DPP-inspired approach to encourage diverse top-k results.

DPPX, ilgili ve çeşitli sonuçları önceliklendirmek için tasarlanmış bir Go kütüphanesidir. Kütüphane, benzerlik hesaplamaları ile özel puanlama fonksiyonlarını birleştirir ve DPP benzeri bir yaklaşımı temel alır.

## Key Features / Temel Özellikler

- Parallel top-k search using CPU workers
- Custom similarity and scoring extension points
- Configurable search engine via `engine.Config`
- Compact, lightweight core based on `types.Item` and `types.Query`
- DPP-inspired diversity-friendly ranking

- CPU işçi süreçleriyle paralel top-k arama
- Özel benzerlik ve puanlama uzantı noktaları
- `engine.Config` ile yapılandırılabilir arama motoru
- `types.Item` ve `types.Query` üzerine kurulmuş hafif çekirdek
- Çeşitliliği destekleyen DPP esinli sıralama

## Installation / Kurulum

```bash
go get github.com/hasan-kilici/dppx
```

## Quick Start / Hızlı Başlangıç

```go
package main

import (
    "github.com/hasan-kilici/dppx/core/engine"
    "github.com/hasan-kilici/dppx/core/similarity"
    "github.com/hasan-kilici/dppx/types"
)

func main() {
    cfg := engine.Config{
        Similarity: similarity.Cosine,
        Scoring: func(query types.Query, item types.Item) float64 {
            return 0.0
        },
    }

    eng := engine.New(cfg)

    query := types.Query{Vector: types.Vector{0.1, 0.2, 0.3}, Norm: 0.374}
    items := []types.Item{
        {ID: "a", Vector: types.Vector{0.1, 0.2, 0.3}, Norm: 0.374},
        {ID: "b", Vector: types.Vector{0.2, 0.1, 0.4}, Norm: 0.469},
    }

    results := eng.Search(query, items, 10)
    _ = results
}
```

## Project Structure / Proje Yapısı

- `core/engine` — search engine, configuration, and selection API
- `core/topk` — fixed-size min-heap for top-k result merging
- `core/similarity` — similarity functions, including cosine similarity
- `core/scoring` — custom scoring function types
- `core/sampling` — sampling interface for future extension
- `types` — item, query, score, and vector domain models

- `core/engine` — arama motoru, yapılandırma ve seçim API'si
- `core/topk` — top-k sonuç birleştirmesi için sabit boyutlu min-heap
- `core/similarity` — kosinüs benzerliği dahil benzerlik fonksiyonları
- `core/scoring` — özel puanlama fonksiyonu türleri
- `core/sampling` — ileriye dönük genişleme için örnekleme arayüzü
- `types` — öğe, sorgu, skorlama ve vektör modelleri

## License / Lisans
This project is licensed under the MIT License.

You are free to use, modify, distribute, and publish this software, provided that the original copyright and license notice are included in substantial portions of the project.

See the `LICENSE` file for more information.