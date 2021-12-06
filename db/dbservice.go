package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

type DBManager struct{
	client Database
}

var DBManagerInstance *DBManager

func (obj *DBManager)CreateConnection(){

}

func (obj *DBManager)CloseConnection(){
	obj.client.close()
}

func (obj *DBManager)Insert(dataBase, col string, doc interface{}) (result interface{}, err error){
	result, err = obj.client.insert(dataBase, col, doc)
	return result, err
}

func (obj *DBManager)Update(dataBase, col string, filter, update interface{})(result interface{}, err error){
	result, err = obj.client.update(dataBase, col, filter, update)
	return result, err
}

func (obj *DBManager)Delete(dataBase, col string, query interface{})(result interface{}, err error){
	result, err = obj.client.delete(dataBase, col, query)
	return result, err
}

func (obj *DBManager)Find( dataBase, col string, query interface{})(result []bson.D, err error){
	result, err = obj.client.find(dataBase, col, query)
	return result, err
}

func init(){
	DBManagerInstance = &DBManager{
		client: Connect("mongodb://localhost:27017"),
		}
}

