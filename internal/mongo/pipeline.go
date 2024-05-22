package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func constructPipeline(collections, operations []string) mongo.Pipeline {
	var matchStage bson.D

	opTypes := bson.A{}

	for _, op := range operations {
		opTypes = append(opTypes, op)
	}

	matchStage = append(matchStage, bson.E{Key: "operationType", Value: bson.D{{Key: "$in", Value: opTypes}}})

	if len(collections) > 0 {
		namespaces := bson.A{}

		for _, col := range collections {
			namespaces = append(namespaces, col)
		}

		matchStage = append(matchStage, bson.E{Key: "ns.coll", Value: bson.D{{Key: "$in", Value: namespaces}}})
	}

	return mongo.Pipeline{bson.D{{Key: "$match", Value: matchStage}}}
}
