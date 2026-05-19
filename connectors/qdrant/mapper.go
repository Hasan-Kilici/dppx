package qdrant

import (
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"

	"github.com/hasan-kilici/dppx/core/similarity"
	"github.com/hasan-kilici/dppx/types"
)

func mapPoint(
	point *qdrant.ScoredPoint,
) types.Item {

	item := types.Item{
		Metadata: map[string]any{},
	}

	/*ID*/
	if point.Id != nil {
		switch id := point.Id.PointIdOptions.(type) {
		case *qdrant.PointId_Num:
			item.ID = fmt.Sprintf(
				"%d",
				id.Num,
			)
		case *qdrant.PointId_Uuid:
			item.ID = id.Uuid
		}
	}

	/*Vector*/
	if point.Vectors != nil &&
		point.Vectors.GetVector() != nil {

		vector := point.Vectors.GetVector().Data

		item.Vector = vector

		item.Norm = similarity.Norm(
			vector,
		)
	}

	/*Payload*/
	for k, v := range point.Payload {
		switch val := v.Kind.(type) {
		case *qdrant.Value_StringValue:
			item.Metadata[k] = val.StringValue
		case *qdrant.Value_DoubleValue:
			item.Metadata[k] = val.DoubleValue
		case *qdrant.Value_IntegerValue:
			item.Metadata[k] = val.IntegerValue
		case *qdrant.Value_BoolValue:
			item.Metadata[k] = val.BoolValue
		}
	}

	return item
}