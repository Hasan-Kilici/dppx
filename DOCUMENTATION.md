# DPPX Documentation / DPPX Dokümantasyon

## Overview / Genel Bakış

DPPX is a Go library designed for relevance-aware ranking while preserving diversity in results. It is inspired by Determinantal Point Processes (DPP), a mathematical framework for subset selection that balances quality with diversity.

DPPX, sonuçlarda çeşitliliği korurken alaka düzeyine duyarlı sıralama için tasarlanmış bir Go kütüphanesidir. Determinantal Point Process (DPP) kavramından ilham alır; bu kavram kalite ile çeşitliliği dengeleyen alt küme seçimini tanımlar.

## What is DPP? / DPP Nedir?

Determinantal Point Processes are probabilistic models over subsets of items. A DPP assigns higher probability to diverse subsets by using a kernel matrix whose determinant reflects both individual quality and pairwise dissimilarity.

DPP, öğe alt kümeleri üzerinde olasılıksal modellerdir. DPP, bireysel kalite ile ikili benzerlik arasındaki dengeyi belirleyen determinant kullanan bir çekirdek matrisiyle çeşitli alt kümelere daha yüksek olasılık atar.

## Why DPP is Useful / DPP Neden Faydalıdır?

DPP is useful in recommendation, search, and summarization because it reduces redundancy while keeping relevant items. It is especially valuable when users expect a varied set of suggestions rather than many similar items.

DPP, öneri, arama ve özetleme alanlarında yararlıdır; çünkü ilgili öğeleri korurken tekrarları azaltır. Benzer öğelerin çoğaldığı durumlarda kullanıcıya çeşitli bir sonuç kümesi sunmak için önemlidir.

## Diversity-Based Recommendation Concepts / Çeşitlilik Tabanlı Öneri Kavramları

- Relevance: how well an item matches the query or user intent.
- Diversity: how different the selected items are from each other.
- Trade-off: a good system balances relevance and diversity to avoid redundant top results.

- Alaka: öğenin sorgu veya kullanıcı niyetiyle uyumu.
- Çeşitlilik: seçilen öğelerin birbirinden ne kadar farklı olduğu.
- Denge: iyi bir sistem, tekrar eden üst sonuçları önlemek için alaka ve çeşitlilik arasında bir denge bulur.

## How This Library Works Internally / Bu Kütüphane İçsel Olarak Nasıl Çalışır?

DPPX engine performs a parallel top-k search over item vectors.

1. `engine.New(cfg)` creates an engine instance with a configuration.
2. `Engine.Search(query, items, k)` feeds items into worker goroutines.
3. Each worker computes similarity and a custom score for each item.
4. Results are collected in per-worker min-heaps.
5. A final merge phase combines local heaps into the global top-k.

DPPX motoru, öğe vektörleri üzerinde paralel bir top-k arama gerçekleştirir.

1. `engine.New(cfg)` yapılandırmayla motor örneği oluşturur.
2. `Engine.Search(query, items, k)` öğeleri işçi gorutine'larına iletir.
3. Her işçi, her öğe için benzerlik ve özel bir puan hesaplar.
4. Sonuçlar işçi başına min-heaplerde toplanır.
5. Son olarak yerel heapler global top-k içinde birleştirilir.

### Core Components / Temel Bileşenler

- `core/engine` — search and selection API. / arama ve seçim API'si.
- `core/topk` — manages a bounded top-k result set. / sınırlı boyutlu sonuç kümesini yönetir.
- `core/similarity` — computes vector similarity. / vektör benzerliğini hesaplar.
- `core/scoring` — provides custom business logic scoring. / özel iş mantığı puanlaması sağlar.
- `types` — defines the domain model. / veri modelini tanımlar.

### Execution Flow / Çalışma Akışı

- `runtime.NumCPU()` retrieves the CPU core count. / `runtime.NumCPU()` ile CPU çekirdek sayısı alınır.
- Each worker maintains its own top-k min-heap. / Her çalışan kendi top-k min-heap'ini tutar.
- `Similarity` and `Scoring` values are combined into a final score. / `Similarity` ve `Scoring` değerleri birleştirilerek nihai skor oluşturulur.
- `topk.MinHeap` stores results with bounded capacity and helps limit memory usage on large datasets. / `topk.MinHeap` sonuçları kısıtlı bir şekilde saklar ve büyük veri setleri için bellek kullanımını sınırlamaya yardımcı olur.

## Installation / Kurulum

```bash
go get github.com/hasan-kilici/dppx
```

## Basic Usage / Temel Kullanım

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
            // Business logic can add metadata-based adjustments.
            return 0.0
        },
    }

    eng := engine.New(cfg)

    query := types.Query{Vector: types.Vector{0.1, 0.2, 0.3}, Norm: 0.374}
    items := []types.Item{
        {ID: "item1", Vector: types.Vector{0.2, 0.1, 0.4}, Norm: 0.469},
        {ID: "item2", Vector: types.Vector{0.5, 0.3, 0.2}, Norm: 0.616},
    }

    results := eng.Search(query, items, 5)
    for _, scored := range results {
        println(scored.Item.ID, scored.Score)
    }
}
```

```go
resultItems := eng.Select(items, 5)
```

## Advanced Usage / Gelişmiş Kullanım

### Custom Scoring Function / Özel Puanlama Fonksiyonu

`engine.Config.Scoring` accepts a function with signature:

```go
func(query types.Query, item types.Item) float64
```

Use this callback to incorporate metadata, business rules, or custom heuristics alongside vector similarity.

### Custom Similarity / Özel Benzerlik

`engine.Config.Similarity` accepts any function matching:

```go
func(a types.Vector, b types.Vector, aNorm float32, bNorm float32) float64
```

The repository includes `core/similarity.Cosine` as a ready-to-use cosine similarity implementation.

### Extensibility Points / Genişletilebilirlik Noktaları

The `engine.Config` struct defines the following fields:

- `Similarity`
- `Scoring`
- `Sampler`
- `CandidatePool`

Current search logic uses `Similarity` and `Scoring` directly, while `Sampler` and `CandidatePool` are defined for future extensions or custom selection strategies.

## Configuration / Yapılandırma

`engine.Config` fields:

- `Similarity` — similarity function between query and item vectors.
- `Scoring` — optional business-scoring hook.
- `Sampler` — sampler interface for custom selection logic.
- `CandidatePool` — integer value reserved for candidate pruning or pool sizing.

### Example Config / Örnek Yapılandırma

```go
cfg := engine.Config{
    Similarity: similarity.Cosine,
    Scoring: func(query types.Query, item types.Item) float64 {
        if query.UserID == item.ID {
            return 1.0
        }
        return 0.0
    },
}
```

## Mathematical Intuition / Matematiksel Sezgi

### Determinantal Point Processes

A DPP uses a kernel matrix `L` where the determinant of a subset's principal minor measures the joint quality of items and their diversity.

- Large determinant => high quality + low redundancy.
- Kernel entries combine item relevance and pairwise similarity.

This approach rewards subsets that are not only high-scoring but also diverse.

Bu yaklaşım, kümeleri sadece yüksek puanlı öğeler olarak değil, aynı zamanda birbirinden farklı olan öğeler olarak da ödüllendirir.

### DPPX Approach

While DPPX does not construct a full DPP kernel matrix in the current implementation, it is designed around the same engineering goal: deliver a top-k result set that balances relevance and variant item selection.

DPPX focuses on producing results that support diversity by combining similarity and custom scoring.

DPPX, benzerlik ve özel puanlamayı birleştirerek çeşitliliği destekleyen sonuçlar üretmeye odaklanır.

## Performance Considerations / Performans Hususları

- `Engine.Search` is parallelized across `runtime.NumCPU()` workers.
- Each worker maintains its own min-heap to limit heap contention.
- The final merge phase combines local heaps into the overall top-k.
- `core/topk.MinHeap` keeps memory usage bounded by `k` per worker.

- `Similarity` and `Scoring` should be efficient, because they execute once per item.
- Vector normalization should be precomputed to avoid repeated overhead.

## Project Structure / Proje Yapısı

- `core/engine/config.go` — engine configuration definition.
- `core/engine/engine.go` — parallel search implementation.
- `core/engine/select.go` — item selection wrapper.
- `core/topk/heap.go` — top-k min-heap implementation.
- `core/similarity/cosine.go` — cosine similarity utility.
- `core/scoring/interface.go` — scoring function signature.
- `core/sampling/interface.go` — sampler interface for extension.
- `types/item.go` — item model.
- `types/query.go` — user query model.
- `types/scoredItem.go` — scored result container.
- `types/vector.go` — vector representation.
- `types/searchResult.go` — search result model.

## Development Notes / Geliştirme Notları

- Keep comments concise and reserve them for complex logic or performance-sensitive code.
- Do not alter core search semantics unless performance improvements are validated.
- `Search` currently merges worker heaps in a separate final phase to preserve parallel throughput.
- Future enhancements can integrate `Sampler` and `CandidatePool` into the search pipeline.
