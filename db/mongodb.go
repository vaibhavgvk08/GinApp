package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDB struct{
	dbclient *mongo.Client
	dbcontext context.Context
	cancel context.CancelFunc
}

func Connect(uri string)  *MongoDB{
	mongodb := &MongoDB{}
	ctx, cancelFunc := context.WithTimeout(context.Background(),
		1000*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil{
		panic(err)
	}

	mongodb.dbclient = client
	mongodb.dbcontext = ctx
	mongodb.cancel = cancelFunc
	return mongodb
}

func (obj *MongoDB)find( dataBase, col string, query interface{}) (result []bson.D, err error) { //*mongo.Cursor
	collection := obj.dbclient.Database(dataBase).Collection(col)
	cursor, err := collection.Find(obj.dbcontext, query)

	if err := cursor.All(obj.dbcontext, &result); err != nil {
		panic(err)
	}
	return
}

func (obj *MongoDB)insert(dataBase, col string, doc interface{})(interface{}, error) { //*mongo.InsertOneResult
	collection := obj.dbclient.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(obj.dbcontext, doc)
	fmt.Println(doc, result, err)
	return result, err
}

func (obj *MongoDB)update(dataBase, col string, filter, update interface{})(result interface{}, err error) { //*mongo.UpdateResult

	collection := obj.dbclient.Database(dataBase).Collection(col)
	result, err = collection.UpdateOne(obj.dbcontext, filter, update)
	return
}

func (obj *MongoDB)delete(dataBase, col string, query interface{})(result interface{}, err error) { //*mongo.DeleteResult

	collection := obj.dbclient.Database(dataBase).Collection(col)
	result, err = collection.DeleteOne(obj.dbcontext, query)
	return

}

func (obj *MongoDB)close() {
	defer obj.cancel()
	defer func() {
		if err := obj.dbclient.Disconnect(obj.dbcontext); err != nil {
			panic(err)
		}
	}()
}


