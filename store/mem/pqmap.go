package mem

import (
	"bytes"
	"container/heap"
	"fmt"
	"inspector/models"
)

// A min Priority Queue that guarantees O(1) search.
// Takes O(2N) space.  A heap and map.
type ProjectMap struct {
	vals     []*data
	indexMap map[uint]int // Maps id to index in vals.
}

type data struct {
	val   models.Project
	index int
}

func NewProjectMap() *ProjectMap {
	m := new(ProjectMap)
	m.indexMap = make(map[uint]int)
	m.vals = make([]*data, 0)
	heap.Init(m)
	return m
}

func (mp *ProjectMap) List() []models.Project {
	num := len(mp.vals)
	list := make([]models.Project, num, num)
	for i := 0; i < num; i++ {
		list[i] = mp.vals[i].val
	}

	return list
}

// Search by key.
// O(1)
func (mp *ProjectMap) Search(k uint) *models.Project {
	data, ok := mp.indexMap[k]
	if !ok {
		return nil
	}
	return &mp.vals[data].val
}

// Insert a k/v pair.  Key being the project's id and value being the project.
// O(log(n))
func (mp *ProjectMap) Set(v *models.Project) {
	old, ok := mp.indexMap[v.Id]
	if ok {
		mp.update(mp.vals[old], *v)
	} else {
		val := data{*v, -1}
		heap.Push(mp, &val)
	}
}

// Pop the max element.
// O(log(n))
func (mp *ProjectMap) Max() *models.Project {
	if mp.Len() == 0 {
		return nil
	}
	item := heap.Pop(mp).(*data)
	delete(mp.indexMap, item.val.Id)
	return &item.val
}

func (mp *ProjectMap) String() string {
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
func (mp *ProjectMap) Len() int {
	return len(mp.vals)
}

func (mp *ProjectMap) Less(i, j int) bool {
	return mp.vals[i].val.LastUpdated.Before(mp.vals[j].val.LastUpdated)
}

func (mp *ProjectMap) Swap(i, j int) {
	mp.vals[i], mp.vals[j] = mp.vals[j], mp.vals[i]
	mp.vals[i].index = i
	mp.vals[j].index = j
	mp.indexMap[mp.vals[i].val.Id] = i
	mp.indexMap[mp.vals[j].val.Id] = j
}

func (mp *ProjectMap) Push(v interface{}) {
	n := len(mp.vals)
	item := v.(*data)
	item.index = n
	mp.vals = append(mp.vals, item)
	mp.indexMap[item.val.Id] = item.index
}

func (mp *ProjectMap) Pop() interface{} {
	n := len(mp.vals)
	item := mp.vals[n-1]
	item.index = -1 // for safety
	mp.vals = mp.vals[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (mp *ProjectMap) update(item *data, newVal models.Project) {
	heap.Remove(mp, item.index)
	delete(mp.indexMap, item.val.Id)
	item.val = newVal
	heap.Push(mp, item)
}
