package mongodb

import (
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/hasan-kilici/dppx/types"
)

func mapDocument(
    doc map[string]any,
    vectorField string,
) types.Item {
    item := types.Item{
        Metadata: map[string]any{},
    }

    if id, ok := doc["_id"]; ok {
        item.ID = fmt.Sprintf("%v", id)
    }

    if rawVector, ok := doc[vectorField].(bson.A); ok {
        vector := make(
            types.Vector,
            len(rawVector),
        )

        for i, v := range rawVector {
            vector[i] = float32(v.(float64))
        }

        item.Vector = vector
    }

    for k, v := range doc {
        if k == vectorField {
            continue
        }

        item.Metadata[k] = v
    }

    return item
}