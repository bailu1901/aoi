package aoi

import (
	"fmt"
	"math/rand"
	"testing"
)

// 事件监听
type emptyListener struct {
}

// OnEnter 进入
func (*emptyListener) OnEnter(id ID, s Set) {
	fmt.Println("OnEnter", id, s)
}

// OnMove 移动
func (*emptyListener) OnMove(id ID, s Set) {
	fmt.Println("OnMove", id, s)
}

// OnLeave 离开
func (*emptyListener) OnLeave(id ID, s Set) {
	fmt.Println("OnLeave", id, s)
}

func TestAOI(t *testing.T) {
	m := NewManager(2, 2, 1000, &emptyListener{})

	var id ID = 1
	m.Add(id, 1, 0)
	m.Add(2, 0, 1)
	m.Add(3, 1, 1)
	m.Add(4, 3, 3)
	m.Add(5, 4, 4)

	fmt.Println(m.head.nextX)

	s := make(Set, 100)
	m.GetRange(id, s)
	fmt.Println("m.GetRange", s)
	s.Clear()

	m.Move(id, 0, 1)
	fmt.Println("m.GetRange", s)
	s.Clear()

	m.Move(id, 3, 0)
	fmt.Println("m.GetRange", s)
	s.Clear()

	m.Move(id, 4, 4)
	fmt.Println("m.GetRange", s)
	s.Clear()

	m.Move(id, 2, 2)
	fmt.Println("m.GetRange", s)
	s.Clear()

	m.Move(id, 0, 0)
	fmt.Println("m.GetRange", s)
	s.Clear()

	fmt.Println(m.head.nextX)
	m.Leave(id)
}

func BenchmarkAdd(b *testing.B) {
	m := NewManager(100, 100, 1000, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Add(ID(i), float32(i), float32(i))
	}
}

func BenchmarkMove(b *testing.B) {
	count := 10000
	m := NewManager(20, 20, count, nil)

	for i := 0; i < count; i++ {
		m.Add(ID(i), float32(i/100), float32(i%100))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := ID(i % count)
		x := float32(rand.Int() % 100)
		y := float32(rand.Int() % 100)
		b.StartTimer()
		m.Move(id, x, y)
	}
}

func BenchmarkLeave(b *testing.B) {
	count := 10000
	m := NewManager(20, 20, count, nil)

	for i := 0; i < count; i++ {
		m.Add(ID(i), float32(i/100), float32(i%100))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := ID(i % count)
		b.StartTimer()
		m.Leave(id)
		b.StopTimer()
		x := float32(rand.Int() % 100)
		y := float32(rand.Int() % 100)
		m.Add(id, x, y)
		b.StartTimer()
	}
}

func BenchmarkRange(b *testing.B) {
	count := 10000
	m := NewManager(20, 20, count, nil)

	for i := 0; i < count; i++ {
		m.Add(ID(i), float32(i/100), float32(i%100))
	}

	s := make(Set, 400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := ID(i % count)
		b.StartTimer()
		m.GetRange(id, s)
	}
}
