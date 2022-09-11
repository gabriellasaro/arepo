package arepo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ID primitive.ObjectID

var NilID ID

func (id ID) Hex() string {
	return primitive.ObjectID(id).Hex()
}

func (id ID) String() string {
	return primitive.ObjectID(id).String()
}

func (id ID) IsZero() bool {
	return id == NilID
}

func IDFromHex[T ~string](s T) (ID, error) {
	id, err := primitive.ObjectIDFromHex(string(s))
	if err != nil {
		return NilID, err
	}

	return ID(id), nil
}
