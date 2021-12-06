package cache

import "github.com/aerospike/aerospike-client-go"

type EvictionPolicy interface {
	keyAccess(key *aerospike.Key)
	evictKey() int
	deleteKey(key *aerospike.Key)
}
