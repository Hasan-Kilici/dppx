# DPPX

English 🇬🇧 / Türkçe 🇹🇷

DPPX is a Go library for diversity-aware ranking and recommendation. It combines vector similarity, configurable business scoring, and optional retrieval connectors to produce balanced top-k results without heavy infrastructure.

DPPX, çeşitlilik odaklı sıralama ve öneri çözümleri için geliştirilmiş bir Go kütüphanesidir. Vektör benzerliği, yapılandırılabilir iş puanlaması ve opsiyonel retriever bağlayıcıları ile ağır altyapıya gerek kalmadan dengeli top-k sonuçlar üretir.

---

## What is DPPX? / DPPX Nedir?

English 🇬🇧

DPPX is a flexible search engine framework that ranks candidate items by combining vector similarity with optional business scoring. It uses parallel workers, local top-k heaps, and a merge step for efficient scoring on large candidate sets.

Türkçe 🇹🇷

DPPX, aday öğeleri vektör benzerliği ve isteğe bağlı iş puanlamasıyla birleştiren esnek bir arama motoru çerçevesidir. Büyük aday kümelerinde verimli puanlama için paralel işçiler, yerel top-k heap'leri ve birleştirme adımı kullanır.

---

## Why diversity matters / Çeşitlilik neden önemli?

English 🇬🇧

Systems that optimize only for relevance can return redundant or overly similar items. Diversity-aware ranking helps surface varied, useful results without sacrificing relevance.

Türkçe 🇹🇷

Sadece alaka odaklı sistemler, tekrarlayan veya aşırı benzer sonuçlar döndürebilir. Çeşitlilik odaklı sıralama, alakayı korurken daha geniş ve faydalı bir sonuç seti sunmaya yardımcı olur.

---

## Key Features / Temel Özellikler

- Parallel top-k scoring with CPU workers / CPU işçileriyle paralel top-k skorlaması
- Configurable similarity and scoring hooks / Yapılandırılabilir benzerlik ve puanlama kancaları
- Pluggable retriever connectors for ANN stores / ANN depoları için takılabilir retriever bağlayıcıları
- Extensible sampling and reranking strategies / Genişletilebilir örnekleme ve yeniden sıralama stratejileri
- Small core design with predictable per-worker memory use / Her işçi için öngörülebilir bellek kullanımı sağlayan küçük çekirdek tasarım

---

## Installation / Kurulum

English 🇬🇧

Install with Go modules:

```bash
go get github.com/hasan-kilici/dppx
```

Türkçe 🇹🇷

Go modülü ile kurulum:

```bash
go get github.com/hasan-kilici/dppx
```

---

## Quick Start / Hızlı Başlangıç

English 🇬🇧

Use the in-memory retriever from `core/retriever/memory` for a minimal local example. This follows `tests/test.go` and shows the required `engine.Config.Retriever`, candidate pool, and `Search` call with `context.Context`.

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
    // Build a small in-memory candidate set.
    items := []types.Item{
        {
            ID: "a",
            Vector: types.Vector{0.1, 0.2, 0.3},
            Norm:   0.374,
        },
        {
            ID: "b",
            Vector: types.Vector{0.2, 0.1, 0.4},
            Norm:   0.469,
        },
    }

    retr := mem.New(items)

    cfg := engine.Config{
        Retriever:     retr,
        CandidatePool: 10,
        Similarity:    similarity.Dot,
    }

    eng := engine.New(cfg)

    query := types.Query{
        Vector: types.Vector{0.1, 0.2, 0.3},
        Norm:   0.374,
    }

    results, err := eng.Search(context.Background(), query, 10)
    if err != nil {
        panic(err)
    }

    fmt.Println(results)
}
```

Türkçe 🇹🇷

`tests/test.go` dosyasına uygun basit bir örnek. `engine.Config.Retriever` gereklidir ve `context.Context` ile `Search` çağrılır.

---

## Advanced Usage / Gelişmiş Kullanım

English 🇬🇧

- `engine.Config.Retriever` is required and provides candidates from memory, connectors, or external ANN stores.
- `engine.Config.Similarity` accepts `func(a types.Vector, b types.Vector, aNorm, bNorm float32) float64`.
- `engine.Config.Scoring` accepts `func(query types.Query, item types.Item) float64`.
- `engine.Config.Sampler` supports strategies such as `sampling.MMR`, `sampling.TopK`, and `sampling.Random`.

Türkçe 🇹🇷

- `engine.Config.Retriever` gereklidir ve bellek, bağlayıcılar veya harici ANN depolarından aday sağlar.
- `engine.Config.Similarity` `func(a types.Vector, b types.Vector, aNorm, bNorm float32) float64` imzasını kabul eder.
- `engine.Config.Scoring` `func(query types.Query, item types.Item) float64` imzasını kabul eder.
- `engine.Config.Sampler` `sampling.MMR`, `sampling.TopK`, `sampling.Random` gibi stratejileri destekler.

---

## Connector Example: Qdrant / Bağlayıcı Örneği: Qdrant

English 🇬🇧

For external ANN retrieval, the Qdrant connector can fetch candidates and map them into `types.Item`. This example follows the `tests/qdrant.go` implementation.

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
            scoring.Weighted{
                Func:   scoring.Popularity,
                Weight: 0.6,
            },
            scoring.Weighted{
                Func:   scoring.FreshnessBoost,
                Weight: 0.3,
            },
            scoring.Weighted{
                Func: scoring.Personalization(
                    similarity.Cosine,
                ),
                Weight: 0.8,
            },
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

func randomVector(size int) types.Vector {
    v := make(types.Vector, size)
    for i := range v {
        v[i] = 0.5
    }
    return v
}
```

Türkçe 🇹🇷

Qdrant bağlayıcısı, harici ANN sorgularını alır, dönüşen sonuçları `types.Item` olarak işlemenizi sağlar ve DPPX skorlamasına dahil eder.

---

## Mathematical Notes / Matematiksel Notlar

English 🇬🇧

DPPX is inspired by Determinantal Point Processes (DPP): a model that prefers subsets with both strong item quality and low pairwise similarity. The current implementation does not build a full DPP kernel matrix by default, but it supports DPP-like engineering through scoring and sampling hooks.

Türkçe 🇹🇷

DPPX, güçlü öğe kalitesini ve düşük ikili benzerliği bir araya getiren alt kümeleri tercih eden Determinantal Point Processes (DPP) ilhamlıdır. Mevcut uygulama varsayılan olarak tam bir DPP çekirdek matrisi oluşturmaz, ancak puanlama ve örnekleme kancalarıyla DPP benzeri mühendisliği destekler.

---

## Performance Considerations / Performans Hususları

English 🇬🇧

- Search is parallelized across `runtime.NumCPU()` workers.
- Each worker keeps a local min-heap of size `k` to limit contention.
- Precompute vector norms and keep `Similarity`/`Scoring` fast, since they run per candidate.

Türkçe 🇹🇷

- Arama `runtime.NumCPU()` işçileri arasında paralelleştirilir.
- Her işçi, içeriği sınırlamak için `k` boyutunda yerel bir min-heap tutar.
- `Similarity` ve `Scoring` fonksiyonlarını aday başına çalıştıkları için hızlı tutun.

---

## Project Structure / Proje Yapısı

English 🇬🇧

- `core/engine` — search engine and configuration
- `core/topk` — top-k min-heap
- `core/similarity` — similarity functions and norms
- `core/scoring` — scoring utilities and pipeline support
- `core/sampling` — sampling and reranking strategies
- `core/retriever` — retriever interfaces and memory connector
- `connectors/qdrant` — external Qdrant retriever connector
- `types` — data models for items, queries, and scored results

Türkçe 🇹🇷

- `core/engine` — arama motoru ve yapılandırma
- `core/topk` — top-k min-heap
- `core/similarity` — benzerlik fonksiyonları ve normlar
- `core/scoring` — puanlama yardımcıları ve pipeline desteği
- `core/sampling` — örnekleme ve yeniden sıralama stratejileri
- `core/retriever` — retriever arayüzleri ve bellek bağlayıcısı
- `connectors/qdrant` — harici Qdrant retriever bağlayıcısı
- `types` — öğeler, sorgular ve skorlu sonuçlar için veri modelleri

---

## Development / Geliştirme Notları

English 🇬🇧

- Keep comments concise and explain why code exists.
- Don’t change core search semantics without validation.

Türkçe 🇹🇷

- Yorumları kısa tutun ve kodun neden var olduğunu açıklayın.
- Çekirdek arama davranışını doğrulamadan değiştirmeyin.

---

## License / Lisans

English 🇬🇧 / Türkçe 🇹🇷

This project is licensed under the MIT License. See the `LICENSE` file for details.

Bu proje MIT Lisansı ile lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakın.
