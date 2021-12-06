package main

import (
	"GinApp/cache"
	"GinApp/db"
	"GinApp/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

type GinApp struct {
	cache *cache.Cache
	dbmanager *db.DBManager
}

var ginapp *GinApp


func createDocument(c *gin.Context) {
	var emp model.Employee
	if err := c.BindJSON(&emp); err != nil {
		panic(err)
	}

	ginapp.cache.Put(emp.Empid, emp)
	ginapp.dbmanager.Insert("test","Employee", bson.D{
		{"Age"	,	emp.Age},
		{"Name"	,	emp.Name},
		{"Company"	,	emp.Company},
		{"Empid"	,	emp.Empid},
	})

	c.String(200, "success")
}

func readDocument(c *gin.Context) {
	queryparams := c.Request.URL.Query()
	empid,_ :=  strconv.Atoi(queryparams["Empid"][0])
	res, err := ginapp.cache.Get(empid)
	if err==nil{
		fmt.Println("Entry found in cache")
		c.JSON(200, res)
		return
	}

	filter := bson.D{}
	for k,v := range queryparams{
		if k == "Company" || k == "Name"{
			filter = append(filter, bson.E{k, v[0]})
		}else {
			newv,_ :=  strconv.Atoi(v[0])
			filter = append(filter, bson.E{k, newv})
		}
	}

	result, err := ginapp.dbmanager.Find("test", "Employee", filter)
	if err != nil {
		panic(err)
	}

	//for _, doc := range result {
	//	fmt.Println(doc)
	//}
	c.JSON(200, result)
}

func updateDocument(c *gin.Context) {
	queryparams := c.Request.URL.Query()
	filter := bson.D{}
	for k,v := range queryparams{
		if k == "Company" || k == "Name"{
			filter = append(filter, bson.E{k, v[0]})
		}else {
			newv,_ :=  strconv.Atoi(v[0])
			filter = append(filter, bson.E{k, newv})
		}
	}

	value := bson.D{}
	var emp model.Employee
	if err := c.BindJSON(&emp); err != nil {
		panic(err)
	}
	if emp.Name != ""{
		value = append(value, bson.E{"Name", emp.Name})
	}
	if emp.Age != 0{
		value = append(value, bson.E{"Age", emp.Age})
	}
	if emp.Empid != 0 {
		value = append(value, bson.E{"Empid", emp.Empid})
	}
	if emp.Company != ""{
		value = append(value, bson.E{"Company", emp.Company})
	}
	update := bson.D{
		{"$set", value},
	}

	ginapp.cache.Put(emp.Empid, emp) //todo make it better
	ginapp.dbmanager.Update("test","Employee", filter, update)
	c.String(200, "success")
}

func deleteDocument(c *gin.Context) {
	queryparams := c.Request.URL.Query()
	filter := bson.D{}
	for k,v := range queryparams{
		if k == "Company" || k == "Name"{
			filter = append(filter, bson.E{k, v[0]})
		}else {
			newv,_ :=  strconv.Atoi(v[0])
			filter = append(filter, bson.E{k, newv})
		}
	}

	k,_ :=  strconv.Atoi(queryparams["Empid"][0])
	ginapp.cache.Delete(k)
	ginapp.dbmanager.Delete("test","Employee", filter)
	c.String(200, "success")
}

func init(){
	ginapp = &GinApp{}
	ginapp.cache = cache.InitializeCache("LRU", 10)
	ginapp.dbmanager = db.DBManagerInstance
}

func CloseConnection(){
	fmt.Println("Closing all connections")
	ginapp.dbmanager.CloseConnection()
	cache.CloseCache(ginapp.cache)
}

func main() {
	route := gin.Default()

	route.POST("/createDoc", createDocument)
	route.GET("/readDoc", readDocument)
	route.POST("/updateDoc", updateDocument)
	route.DELETE("/deleteDoc", deleteDocument)

	route.Run(":8085")
	defer CloseConnection()
}