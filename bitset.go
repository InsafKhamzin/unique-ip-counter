package main

import (
	"sync"

	"github.com/bits-and-blooms/bitset"
)

// ShardedBitset is a bitsets collection sharded by the first octet of ip address
// this will help to avoid frequent locking on single bitset
type ShardedBitset struct {
	shards []*bitset.BitSet
	locks  []sync.Mutex
}

// NewShardedBitset creates bitsets array of 256 bitsets and 256 mutexex accordingly
func NewShardedBitset() *ShardedBitset {
	size := 256
	shards := make([]*bitset.BitSet, size)
	locks := make([]sync.Mutex, size)
	for i := 0; i < size; i++ {
		size := uint(1) << 24 //last three octets
		shards[i] = bitset.New(size)
	}
	return &ShardedBitset{shards: shards, locks: locks}
}

// Set locks the mutex under certain index, writes to bitset
func (sb *ShardedBitset) Set(firstOctet uint8, lastThree uint) {
	sb.locks[firstOctet].Lock()
	defer sb.locks[firstOctet].Unlock()
	sb.shards[firstOctet].Set(lastThree)
}

// Count aggregates all 256 bitsets count
func (sb *ShardedBitset) Count() uint {
	var count uint = 0
	for _, v := range sb.shards {
		count += v.Count()
	}
	return count
}
