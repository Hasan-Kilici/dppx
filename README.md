# DPPX

English / Türkçe

DPPX is a lightweight Go library for producing relevance-aware and diversity-friendly top-k results. It combines vector similarity with optional business scoring and supports diversity-aware selection strategies (DPP-inspired engineering) while remaining easy to integrate and extend.

DPPX, ilgili sonuçları korurken çeşitliliği teşvik eden hafif bir Go kütüphanesidir. Vektör benzerliği ile isteğe bağlı iş puanlamasını birleştirir ve DPP esintili çeşitlilik stratejilerine destek verir; kullanım ve genişletme açısından basit kalır.

---

## What is DPPX? / DPPX Nedir?

English

DPPX provides a configurable search engine to score and rank candidate items using vector similarity and optional custom scoring. The engine is implemented with parallel workers and a merge-based top-k strategy for throughput and predictable memory usage.

Türkçe

DPPX, vektör benzerliği ve isteğe bağlı özel puanlama kullanarak aday öğeleri puanlayan ve sıralayan yapılandırılabilir bir arama motoru sunar. Motor, yüksek verim için paralel işçiler ve öngörülebilir bellek kullanımı için birleştirme tabanlı top-k stratejisiyle uygulanmıştır.

---

## Why diversity / Çeşitlilik neden önemli?

English

Recommenders and search systems can return many near-duplicate items if they optimize only for relevance. Diversity-aware selection helps reduce redundancy and surfaces a broader, more useful set of items to users.

Türkçe

Sadece alaka odaklı sistemler, birbirine çok benzeyen sonuçlar döndürebilir. Çeşitlilik odaklı seçim, tekrarları azaltır ve kullanıcılara daha çeşitli ve faydalı öğeler sunar.

---

## Key Features / Temel Özellikler

- English / Türkçe
- Parallel top-k search using CPU workers / CPU işçi süreçleriyle paralel top-k arama
- Configurable similarity and business scoring hooks / Yapılandırılabilir benzerlik ve iş puanlama kancaları
- Extensible sampling/selection strategies (MMR, TopK, Random, etc.) / Genişletilebilir örnekleme/seçim stratejileri (MMR, TopK, Random vb.)
- Small core with predictable memory usage per worker / Her işçi için öngörülebilir bellek kullanımı sağlayan küçük çekirdek

---

## Installation / Kurulum

English

Install with `go get`:

```bash
go get github.com/hasan-kilici/dppx
```

Türkçe

`go get` ile kurulum:

```bash
go get github.com/hasan-kilici/dppx
```

---

## Quick Start / Hızlı Başlangıç

English

Simple example that demonstrates the in-repo `tests` usage: create an in-memory retriever, build an engine, and call `Search` with a `context.Context`.

```go
package main

import (
    "context"
    "fmt"

    "github.com/hasan-kilici/dppx/core/engine"
    mem "github.com/hasan-kilici/dppx/core/retriever/memory"
    "github.com/hasan-kilici/dppx/core/similarity"
    "github.com/hasan-kilici/dppx/types"
)

func main() {
    // prepare in-memory items (example data)
    items := []types.Item{
        {ID: "a", Vector: types.Vector{0.1, 0.2, 0.3}, Norm: 0.374},
        {ID: "b", Vector: types.Vector{0.2, 0.1, 0.4}, Norm: 0.469},
    }

    // create a simple retriever backed by memory
    retr := mem.New(items)

    cfg := engine.Config{
        Retriever:     retr,
        CandidatePool: 10,
        Similarity:    similarity.Cosine,
    }

    eng := engine.New(cfg)

    query := types.Query{Vector: types.Vector{0.1, 0.2, 0.3}, Norm: 0.374}

    // execute search with context and handle errors
    res, err := eng.Search(context.Background(), query, 10)
    if err != nil {
        panic(err)
    }

    fmt.Println(res)
}
```

Türkçe

Aşağıdaki örnek `tests` klasöründeki kullanıma uygundur: bellek tabanlı bir retriever oluşturun, motoru yapılandırın ve `context.Context` ile `Search` çağırın.


---

## Advanced Usage / Gelişmiş Kullanım

English

- `engine.Config.Scoring` accepts a function `func(query types.Query, item types.Item) float64` to apply business logic.
- `engine.Config.Similarity` accepts a function `func(a types.Vector, b types.Vector, aNorm, bNorm float32) float64` — pass precomputed norms to avoid repeated sqrt.
- `engine.Config.Sampler` supports pluggable selection strategies such as `sampling.MMR`, `sampling.TopK`, and `sampling.Random`.

Türkçe

- `engine.Config.Scoring` iş mantığı uygulamak için `func(query types.Query, item types.Item) float64` imzasını kabul eder.
- `engine.Config.Similarity` `func(a types.Vector, b types.Vector, aNorm, bNorm float32) float64` imzasını kabul eder — tekrar eden sqrt işlemlerini önlemek için normları önceden hesaplayın.
- `engine.Config.Sampler` `sampling.MMR`, `sampling.TopK`, `sampling.Random` gibi takılabilir seçim stratejilerini destekler.

---

## Connector Example: Qdrant / Bağlayıcı Örneği: Qdrant

English

If you use an external ANN store like Qdrant, DPPX supports a retriever connector that fetches candidates from the vector DB and converts them into `types.Item` objects for reranking. The following example is adapted from the repository `tests/qdrant.go` file.

```go
package main

import (
    "context"
    "fmt"

    "github.com/hasan-kilici/dppx/connectors/qdrant"
    "github.com/hasan-kilici/dppx/core/engine"
    "github.com/hasan-kilici/dppx/core/sampling"
    "github.com/hasan-kilici/dppx/core/scoring"
    "github.com/hasan-kilici/dppx/core/similarity"
    "github.com/hasan-kilici/dppx/types"
)

func main() {
    // create a Qdrant retriever (see tests/qdrant.go)
    retr, err := qdrant.New(qdrant.Config{
        Host:       "localhost",
        Port:       6334,
        Collection: "blogs",
    })
    if err != nil {
        panic(err)
    }

    eng := engine.New(engine.Config{
        Retriever:     retr,
        CandidatePool: 200,
        Similarity:    similarity.Cosine,
        Scoring: scoring.Combine(
            scoring.Weighted{Func: scoring.Popularity, Weight: 0.6},
            scoring.Weighted{Func: scoring.FreshnessBoost, Weight: 0.3},
            scoring.Weighted{Func: scoring.Personalization, Weight: 0.8},
        ),
        Sampler: sampling.MMR{Lambda: 0.7},
    })

    query := types.Query{UserID: "user-1", Vector: randomVector(768)}
    query.Norm = similarity.Norm(query.Vector)

    results, err := eng.Search(context.Background(), query, 10)
    if err != nil {
        panic(err)
    }

    for _, r := range results {
        fmt.Printf("ITEM=%s SCORE=%.4f\n", r.Item.ID, r.Score)
    }
}

// randomVector is a simple placeholder — replace with real embeddings in production.
func randomVector(size int) types.Vector {
    v := make(types.Vector, size)
    for i := range v { v[i] = 0.5 }
    return v
}
```

Türkçe

Qdrant gibi harici bir ANN deposu kullanıyorsanız, DPPX bu veritabanından adayları alıp `types.Item`'a dönüştüren bir retriever bağlayıcısını destekler. Yukarıdaki örnek `tests/qdrant.go` dosyasından uyarlanmıştır.


## Mathematical Notes / Matematiksel Notlar

English

DPPX is inspired by Determinantal Point Processes (DPP): a mathematical model that prefers subsets with both high-quality items and low pairwise similarity. This implementation does not build a full DPP kernel matrix by default, but it provides patterns and hooks (scoring, sampling) that enable diversity-aware selection (e.g., using MMR).

Türkçe

DPPX, hem yüksek kalite hem de düşük ikili benzerlik içeren alt kümeleri tercih eden Determinantal Point Processes (DPP) ilhamlıdır. Bu uygulama varsayılan olarak tam bir DPP çekirdek matrisi oluşturmaz, ancak MMR gibi çeşitlilik odaklı seçimleri destekleyecek puanlama ve örnekleme kancaları sunar.

---

## Performance Considerations / Performans Hususları

English

- Search is parallelized across `runtime.NumCPU()` workers; each worker keeps a local min-heap of size `k` to limit contention and memory.
- Precompute vector norms and keep `Similarity`/`Scoring` implementations efficient — they are invoked per candidate.

Türkçe

- Arama `runtime.NumCPU()` işçisi arasında paralelleştirilir; her işçi, içeriği azaltmak için `k` boyutunda yerel bir min-heap tutar.
- Vektör normlarını önceden hesaplayın ve `Similarity`/`Scoring` uygulamalarını verimli tutun — bunlar aday başına çağrılır.

---

## Project Structure / Proje Yapısı

English

- `core/engine` — search engine and configuration
- `core/topk` — top-k min-heap
- `core/similarity` — similarity functions
- `core/scoring` — scoring utilities and hooks
- `core/sampling` — selection/sampling strategies
- `types` — `Item`, `Query`, `Vector`, `ScoredItem`

Türkçe

- `core/engine` — arama motoru ve yapılandırma
- `core/topk` — top-k min-heap
- `core/similarity` — benzerlik fonksiyonları
- `core/scoring` — puanlama yardımcıları ve kancalar
- `core/sampling` — seçim/örnekleme stratejileri
- `types` — `Item`, `Query`, `Vector`, `ScoredItem`

---

## Development / Geliştirme Notları

English

- Keep comments concise and focused on why a piece of code exists or why an optimization is necessary.
- Avoid changing core search semantics without benchmarks.

Türkçe

- Yorumları kısa ve kodun neden var olduğunu veya hangi optimizasyonun gerekli olduğunu açıklayacak şekilde tutun.
- Benchmark olmadan çekirdek arama davranışını değiştirmeyin.

---

## License / Lisans

English & Türkçe

This project is licensed under the MIT License. See the `LICENSE` file for details.

Proje MIT Lisansı ile lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakın.
