package mem

import (
	"math/rand"
	"testing"
	"time"
)

func TestSetAndMin(t *testing.T) {
	m := newPQMap()
	m.Set(1, Project{Id: 1, Name: "1", LastUpdated: time.Now()})
	if p := m.Search(1); p == nil || p.Name != "1" {
		t.Error("Set failed.")
	}
	if m.Len() != 1 {
		t.Error("length should be 1 but was", m.Len())
	}
	m.Set(2, Project{Id: 2, Name: "2", LastUpdated: time.Now()})
	if p := m.Search(2); p == nil || p.Name != "2" {
		t.Error("Set failed.")
	}
	m.Set(3, Project{Id: 3, Name: "3", LastUpdated: time.Now()})
	if p := m.Search(3); p == nil || p.Name != "3" {
		t.Error("Set failed.")
	}
	m.Set(4, Project{Id: 4, Name: "4", LastUpdated: time.Now()})
	if p := m.Search(4); p == nil || p.Name != "4" {
		t.Error("Set failed.")
	}
	m.Set(5, Project{Id: 5, Name: "5", LastUpdated: time.Now()})
	if p := m.Search(5); p == nil || p.Name != "5" {
		t.Error("Set failed.")
	}
	if m.Len() != 5 {
		t.Error("Length should have been 5, but was", m.Len())
	}

	if max := m.Max(); max.Id != 5 {
		t.Error("max ID should have been 5 but was")
	}
	if m.Len() != 4 {
		t.Error("Length should have been 4, but was", m.Len())
	}

	if max := m.Max(); max.Id != 4 {
		t.Error("max ID should have been 4 but was", max.Id)
	}

	if max := m.Max(); max.Id != 3 {
		t.Error("max ID should have been 3 but was", max.Id)
	}

	if max := m.Max(); max.Id != 2 {
		t.Error("max ID should have been 2 but was", max.Id)
	}

	if max := m.Max(); max.Id != 1 {
		t.Error("max ID should have been 1 but was", max.Id)
	}
	if m.Len() != 0 {
		t.Error("Length should have been 0, but was", m.Len())
	}
}

func BenchmarkSet25k(b *testing.B) {
	bennchmarkSetnk(25000, b)
}

func BenchmarkSet50k(b *testing.B) {
	bennchmarkSetnk(50000, b)
}

func BenchmarkSet100k(b *testing.B) {
	bennchmarkSetnk(100000, b)
}

func BenchmarkSet200k(b *testing.B) {
	bennchmarkSetnk(200000, b)
}

func BenchmarkSet300k(b *testing.B) {
	bennchmarkSetnk(300000, b)
}

func BenchmarkSet400k(b *testing.B) {
	bennchmarkSetnk(400000, b)
}

func BenchmarkSet500k(b *testing.B) {
	bennchmarkSetnk(500000, b)
}

func BenchmarkSet600k(b *testing.B) {
	bennchmarkSetnk(600000, b)
}

func BenchmarkSet700k(b *testing.B) {
	bennchmarkSetnk(700000, b)
}

func BenchmarkSet800k(b *testing.B) {
	bennchmarkSetnk(800000, b)
}

func BenchmarkSet900k(b *testing.B) {
	bennchmarkSetnk(900000, b)
}

func BenchmarkSet1000k(b *testing.B) {
	bennchmarkSetnk(1000000, b)
}

func BenchmarkSet1100k(b *testing.B) {
	bennchmarkSetnk(1100000, b)
}

func BenchmarkSet1200k(b *testing.B) {
	bennchmarkSetnk(1200000, b)
}

func bennchmarkSetnk(n int, b *testing.B) {
	m := newPQMap()
	rand.Seed(time.Now().Unix())
	for i := 1; i <= n; i++ {
		j := uint(rand.Intn(i))
		m.Set(j, Project{Id: j, Name: "name", LastUpdated: time.Now()})
	}
}
