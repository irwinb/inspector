package mem

import (
	"bytes"
	"container/heap"
	"fmt"
	"github.com/irwinb/inspector/models"
)

// A min Priority Queue that guarantees O(1) search.
// Takes O(2N) space.  A heap and map.
type projectMap struct {
	vals     []*data
	indexMap map[uint]int // Maps id to index in vals.
}

type data struct {
	key   uint
	val   models.Project
	index int
}

func newProjectMap() *projectMap {
	m := new(projectMap)
	m.indexMap = make(map[uint]int)
	m.vals = make([]*data, 0)
	heap.Init(m)
	return m
}

// Search by key.
// O(1)
func (mp *projectMap) Search(k uint) *models.Project {
	data, ok := mp.indexMap[k]
	if !ok {
		return nil
	}
	return &mp.vals[data].val
}

// Insert a k/v pair.
// O(log(n))
func (mp *projectMap) Set(k uint, v models.Project) {
	old, ok := mp.indexMap[k]
	if ok {
		mp.update(mp.vals[old], v)
	} else {
		val := data{v.Id, v, 0}
		heap.Push(mp, &val)
		mp.indexMap[k] = val.index
	}
}

// Pop the max element.
// O(log(n))
func (mp *projectMap) Max() *models.Project {
	if mp.Len() == 0 {
		return nil
	}
	item := heap.Pop(mp).(*data)
	delete(mp.indexMap, item.key)
	return &item.val
}

func (mp *projectMap) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, val := range mp.vals {
		buffer.WriteString("{")
		buffer.WriteString(fmt.Sprint(*val))
		buffer.WriteString("}")
		if i < len(mp.vals)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}

////////////////////////////////////////////////////////////////////////////
// All methods below are priority queue related.
////////////////////////////////////////////////////////////////////////////
func (mp *projectMap) Len() int {
	return len(mp.vals)
}

func (mp *projectMap) Less(i, j int) bool {
	return mp.vals[i].val.LastUpdated.Before(mp.vals[j].val.LastUpdated)
}

func (mp *projectMap) Swap(i, j int) {
	mp.vals[i], mp.vals[j] = mp.vals[j], mp.vals[i]
	mp.vals[i].index = i
	mp.vals[j].index = j
}

func (mp *projectMap) Push(v interface{}) {
	n := len(mp.vals)
	item := v.(*data)
	item.index = n
	mp.vals = append(mp.vals, item)
}

func (mp *projectMap) Pop() interface{} {
	n := len(mp.vals)
	item := mp.vals[n-1]
	item.index = -1 // for safety
	mp.vals = mp.vals[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (mp *projectMap) update(item *data, newVal models.Project) {
	heap.Remove(mp, item.index)
	mp.vals[item.index].val = newVal
	heap.Push(mp, item)
}
