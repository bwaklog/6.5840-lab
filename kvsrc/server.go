package kvsrv

import (
	"log"
	"strings"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex

	// Your definitions here.

	// key and value of in memory store
	// is a string
	MemoryStore map[string]string
	ClientClocks map[int64]int
}

// • fetch the current value for the key
// • non existant should return an empty string
func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	kv.mu.Lock()
	if val, ok := kv.MemoryStore[args.Key]; ok {
		reply.Value = val
	} else {
		reply.Value = ""
	}
	kv.mu.Unlock()

}

// • replace or install values for a key in the sotre
func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.ClientClocks[args.ClientId] == args.LamportClock {
		reply.Value = kv.MemoryStore[args.Key]
		return
	}

	kv.ClientClocks[args.ClientId] = args.LamportClock
	kv.MemoryStore[args.Key] = args.Value
}

// • append the arg to the keys value
// • return the old value
// • append to non existant: consider existing value is 0 lenght string
func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// check if it is a repeated request
	if kv.ClientClocks[args.ClientId] == args.LamportClock {
		reply.Value= strings.Split(kv.MemoryStore[args.Key], args.Value)[0]
		return
	}

	kv.ClientClocks[args.ClientId] = args.LamportClock

	var val string
	val, ok := kv.MemoryStore[args.Key]
	if !ok {
		val = ""
	}
	reply.Value = val

	// appending the arg to the existing value
	kv.MemoryStore[args.Key] = val + args.Value
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.
	kv.MemoryStore = make(map[string]string)
	kv.ClientClocks = make(map[int64]int)
	return kv
}
