package cache

import "sync"

var localQueue *LocalQueue
var once sync.Once

type LocalQueue struct {
	localBuffer chan []byte
}

func newLocalQueue() *LocalQueue {
	sync.Once
	return localQueue
}
