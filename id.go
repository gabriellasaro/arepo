package arepo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ID = primitive.ObjectID

var NilID ID

func IDFromHex[T ~string](s T) (ID, error) {
	id, err := primitive.ObjectIDFromHex(string(s))
	if err != nil {
		return NilID, err
	}

	return id, nil
}
