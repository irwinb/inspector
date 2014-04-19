package models

import (
	"container/heap"
)

// A min Priority Queue that guarantees O(1) search.
// Takes O(2N) space.  A heap and map.
type pqMap struct {
	vals     []*data
	indexMap map[uint]int // Maps id to index in vals.
}

type data struct {
	key   uint
	val   Project
	index int
}

func newPQMap() *pqMap {
	m := new(pqMap)
	m.indexMap = make(map[uint]int)
	m.vals = make([]*data, 0)
	heap.Init(m)
	return m
}

// Search by key.
// O(1)
func (mp *pqMap) Get(k uint) *Project {
	data, ok := mp.indexMap[k]
	if !ok {
		return nil
	}
	return &mp.vals[data].val
}

// Insert a k/v pair.
// O(log(n))
func (mp *pqMap) Set(k uint, v Project) {
	old, ok := mp.indexMap[k]
	if ok {
		mp.update(mp.vals[old], v)
	} else {
		val := data{v.Id, v, 0}
		mp.Push(&val)
		mp.indexMap[k] = val.index
	}
}

// Pop the max element.
// O(log(n))
func (mp *pqMap) Max() *Project {
	if mp.Len() == 0 {
		return nil
	}
	item := mp.Pop().(*data)
	delete(mp.indexMap, item.key)
	return &item.val
}

// Peek the max element.
// O(1)
func (mp *pqMap) Peek() *Project {
	n := len(mp.vals)
	return &mp.vals[n-1].val
}

////////////////////////////////////////////////////////////////////////////
// All methods below are priority queue related.
////////////////////////////////////////////////////////////////////////////
func (mp *pqMap) Len() int {
	return len(mp.vals)
}

func (mp *pqMap) Less(i, j int) bool {
	return mp.vals[i].val.LastUpdated.Before(mp.vals[j].val.LastUpdated)
}

func (mp *pqMap) Swap(i, j int) {
	mp.vals[i], mp.vals[j] = mp.vals[j], mp.vals[i]
	mp.vals[i].index = i
	mp.vals[j].index = j
}

func (mp *pqMap) Push(v interface{}) {
	n := len(mp.vals)
	item := v.(*data)
	item.index = n
	mp.vals = append(mp.vals, item)
	heap.Fix(mp, n)
}

func (mp *pqMap) Pop() interface{} {
	n := len(mp.vals)
	item := mp.vals[n-1]
	item.index = -1 // for safety
	mp.vals = mp.vals[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (mp *pqMap) update(item *data, newVal Project) {
	mp.vals[item.index].val = newVal
	heap.Fix(mp, item.index)
}
