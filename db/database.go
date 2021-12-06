package db

import "go.mongodb.org/mongo-driver/bson"

type Database interface {
	//insert()
	//update()
	//delete()
	//find()

	insert(string, string, interface{}) (result interface{}, err error)
	update(string, string, interface{}, interface{}) (result interface{}, err error)
	delete(string, string, interface{}) (result interface{}, err error)
	find(string, string, interface{}) ( result []bson.D, err error)

	close()
}
