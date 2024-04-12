package arepo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrNotUpdated = errors.New("not updated")
	ErrNotDeleted = errors.New("not deleted")
)

// FindOne realiza uma busca no banco de dados MongoDB e retorna um único documento correspondente ao filtro especificado.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a busca será realizada.
//   - filter (any): O filtro para a busca no banco de dados.
//   - opts (...*options.FindOneOptions): Opções adicionais para a operação de busca. (Opcional)
//
// Retorna:
//   - (*T): Um ponteiro para o documento encontrado, se existir.
//   - error: Retorna ErrNotFound se nenhum documento correspondente for encontrado. Retorna o erro do MongoDB se ocorrer um problema durante a busca.
func FindOne[T any](ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.FindOneOptions) (*T, error) {
	result := collection.FindOne(ctx, filter, opts...)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, result.Err()
	}

	var data T

	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// FindOneByID realiza uma busca no banco de dados MongoDB e retorna um único documento com base no ID especificado.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a busca será realizada.
//   - id (primitive.ObjectID): O ID do documento a ser encontrado.
//   - opts (...*options.FindOneOptions): Opções adicionais para a operação de busca. (Opcional)
//
// Retorna:
//   - (*T): Um ponteiro para o documento encontrado, se existir.
//   - error: Retorna ErrNotFound se nenhum documento correspondente for encontrado. Retorna o erro do MongoDB se ocorrer um problema durante a busca.
func FindOneByID[T any](ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, opts ...*options.FindOneOptions) (*T, error) {
	return FindOne[T](ctx, collection, bson.M{
		"_id": id,
	})
}

// FindOneAndUpdate executa uma operação de atualização atômica em um único documento no MongoDB e retorna o documento antes da atualização.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a operação será executada.
//   - filter (any): O filtro que seleciona o documento a ser atualizado.
//   - update (any): As atualizações a serem aplicadas no documento.
//   - opts (...*options.FindOneAndUpdateOptions): Opções adicionais para a operação de atualização. (Opcional)
//
// Retorna:
//   - (*T): Um ponteiro para o documento antes da atualização, se existir.
//   - error: Retorna ErrNotFound se nenhum documento correspondente for encontrado. Retorna o erro do MongoDB se ocorrer um problema durante a operação de atualização.
func FindOneAndUpdate[T any](ctx context.Context, collection *mongo.Collection, filter, update any, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	result := collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, result.Err()
	}

	var data T

	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Find executa uma operação de consulta no MongoDB e retorna uma lista de documentos que correspondem ao filtro especificado.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a operação será executada.
//   - filter (any): O filtro que seleciona os documentos a serem recuperados.
//   - opts (...*options.FindOptions): Opções adicionais para a operação de consulta. (Opcional)
//
// Retorna:
//   - ([]*T): Uma lista de ponteiros para os documentos correspondentes ao filtro.
//   - error: Retorna o erro do MongoDB se ocorrer um problema durante a operação de consulta.
func Find[T any](ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.FindOptions) ([]*T, error) {
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var list []*T

	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// UpdateOneByID atualiza um único documento no MongoDB com base no ID especificado.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a operação será executada.
//   - id (primitive.ObjectID): O ID do documento a ser atualizado.
//   - update (any): As atualizações a serem aplicadas ao documento.
//
// Retorna:
//   - error: Retorna ErrNotFound ou ErrNotUpdated se nenhum documento correspondente for encontrado ou atualizado.
//     Retorna o erro do MongoDB se ocorrer um problema durante a operação de atualização.
func UpdateOneByID(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, update any) error {
	r, err := collection.UpdateOne(
		ctx,
		bson.M{
			"_id": id,
		},
		update,
	)
	if err != nil {
		return err
	}

	if r.MatchedCount == 0 {
		return ErrNotFound
	}

	if r.ModifiedCount == 0 {
		return ErrNotUpdated
	}

	return nil
}

// DeleteOne exclui um único documento do MongoDB com base no filtro especificado.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a operação será executada.
//   - filter (any): O filtro para identificar o documento a ser excluído.
//   - opts ([]*options.DeleteOptions): Opções adicionais para a operação de exclusão.
//
// Retorna:
//   - error: Retorna ErrNotDeleted se o valor de 'DeletedCount' for igual a zero.
//     Retorna o erro do MongoDB se ocorrer um problema durante a operação de exclusão.
func DeleteOne(ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.DeleteOptions) error {
	r, err := collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	if r.DeletedCount == 0 {
		return ErrNotDeleted
	}

	return nil
}

// DeleteOneByID exclui um único documento do MongoDB com base no ID especificado.
//
// Parâmetros:
//   - ctx (context.Context): O contexto da execução da função.
//   - collection (*mongo.Collection): A coleção MongoDB na qual a operação será executada.
//   - id (primitive.ObjectID): O ID do documento a ser excluído.
//   - opts ([]*options.DeleteOptions): Opções adicionais para a operação de exclusão.
//
// Retorna:
//   - error: Retorna ErrNotDeleted se o valor de 'DeletedCount' for igual a zero.
//     Retorna o erro do MongoDB se ocorrer um problema durante a operação de exclusão.
func DeleteOneByID(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, opts ...*options.DeleteOptions) error {
	return DeleteOne(
		ctx,
		collection,
		bson.M{
			"_id": id,
		},
		opts...,
	)
}
