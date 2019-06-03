package aoi

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	maxCount   = 10000
	mapX       = 100
	mapY       = 100
	viewRangeX = 10
	viewRangeY = 10
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
	m := NewManager(2, 2, maxCount, &emptyListener{})

	var id ID = 1
	m.Enter(id, 1, 0)
	m.Enter(2, 0, 1)
	m.Enter(3, 1, 1)
	m.Enter(4, 3, 3)
	m.Enter(5, 4, 4)

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
	m := NewManager(viewRangeX, viewRangeY, maxCount, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Enter(ID(i), float32(i), float32(i))
	}
}

func BenchmarkMove(b *testing.B) {

	m := NewManager(viewRangeX, viewRangeY, maxCount, nil)

	for i := 0; i < maxCount; i++ {
		m.Enter(ID(i), float32(i/mapX), float32(i%mapY))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := ID(i % maxCount)
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)
		b.StartTimer()
		m.Move(id, x, y)
	}
}

func BenchmarkLeave(b *testing.B) {
	m := NewManager(viewRangeX, viewRangeY, maxCount, nil)

	for i := 0; i < maxCount; i++ {
		m.Enter(ID(i), float32(i/mapX), float32(i%mapY))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := ID(i % maxCount)
		b.StartTimer()
		m.Leave(id)
		b.StopTimer()
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)
		m.Enter(id, x, y)
		b.StartTimer()
	}
}

func BenchmarkRange(b *testing.B) {
	m := NewManager(viewRangeX, viewRangeY, maxCount, nil)

	for i := 0; i < maxCount; i++ {
		m.Enter(ID(i), float32(i/mapX), float32(i%mapY))
	}

	s := make(Set, viewRangeX*viewRangeY)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := ID(i % maxCount)
		b.StartTimer()
		m.GetRange(id, s)
	}
}
