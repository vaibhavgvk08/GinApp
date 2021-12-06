package cache

import (
	"GinApp/model"
	"fmt"
	aero "github.com/aerospike/aerospike-client-go"
)

// This is only for this example.
// Please handle errors properly.
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type Cache struct {
	client         *aero.Client
	evictionpolicy EvictionPolicy
}

func InitializeCache(evictionpolicytype string, cachesize int) *Cache{
	c := &Cache{
		client: nil,
		evictionpolicy: nil,
	}
	if evictionpolicytype == "LRU"{
		c.evictionpolicy = &LRUPolicy{
			keymap: make(map[*aero.Key]int, cachesize),
			queue: make([]*aero.Key, cachesize),
			qsize: cachesize,
		}
	}

	var err error
	c.client, err = aero.NewClient("0.0.0.0", 3000)
	panicOnError(err)

	return c
}

func CloseCache(c *Cache){
	c.client.Close()
}

func (c *Cache)CreateBinMap(value interface{}) aero.BinMap{
	bins := aero.BinMap{}
	emp := value.(model.Employee)
	fmt.Println(emp)
	if emp.Name != ""{
		bins["Name"] = emp.Name
	}
	if emp.Age != 0{
		bins["Age"] = emp.Age
	}
	if emp.Empid != 0 {
		bins["Empid"] = emp.Empid
	}
	if emp.Company != ""{
		bins["Company"] = emp.Company
	}
	fmt.Println(bins)
	return bins
}

func (c *Cache)Put(key int, value interface{}){ // what will be key and value here ?
	aerokey, err := aero.NewKey("hermes", "aerospike", key)
	c.evictionpolicy.keyAccess(aerokey)

	bins := c.CreateBinMap(value)

	// write the bins
	err = c.client.Put(nil, aerokey, bins)
	panicOnError(err)
}

func (c *Cache)Get(key int) (aero.BinMap, error){
	aerokey, _ := aero.NewKey("hermes", "aerospike", key)
	c.evictionpolicy.keyAccess(aerokey)
	rec, err := c.client.Get(nil, aerokey)
	return rec.Bins, err
}

func (c *Cache)Delete(key int){
	aerokey, _ := aero.NewKey("hermes", "aerospike", key)
	c.evictionpolicy.deleteKey(aerokey)
	_, err := c.client.Delete(nil, aerokey)
	panicOnError(err)
}
