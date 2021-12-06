package cache

import "github.com/aerospike/aerospike-client-go"

type LRUPolicy struct{
	// internal datastructure to check key presence/absence
	keymap map[*aerospike.Key]int
	queue []*aerospike.Key
	qsize int
}

func (obj *LRUPolicy) keyAccess(key *aerospike.Key) {
	if idx, ok := obj.keymap[key]; ok {
		obj.keymap[key] = 0
		obj.queue = append(obj.queue[:idx], obj.queue[idx+1:]...)
		obj.queue = append([]*aerospike.Key{key}, obj.queue...)
	}else{
		if len(obj.queue) < obj.qsize {
			obj.queue = append(obj.queue, key)
			obj.keymap[key] = len(obj.queue ) - 1
		}else{
			idx = obj.evictKey()
			obj.queue = append(obj.queue[:idx], obj.queue[idx+1:]...)
			obj.queue = append([]*aerospike.Key{key}, obj.queue...)
			obj.keymap[key] = 0
		}
	}
}

func (obj *LRUPolicy) evictKey() int{
	idx :=	len(obj.queue)-1
	return idx
}

func (obj *LRUPolicy) deleteKey(key *aerospike.Key) {
	if idx, ok := obj.keymap[key]; ok {
		obj.queue = append(obj.queue[:idx], obj.queue[idx+1:]...)
		delete(obj.keymap, key)
	}
}