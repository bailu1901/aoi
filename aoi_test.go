package aoi

import (
	"math/rand"
	"testing"
)

const (
	maxCount   = 5000
	mapX       = 200
	mapY       = 200
	viewRangeX = 20
	viewRangeY = 20
)

// 事件监听
type emptyListener struct {
	leaveSet Set // 通知离开的集合
	enterSet Set // 通知进入的集合
	moveSet  Set // 通知移动的集合
}

func newEmptyListener() *emptyListener {
	return &emptyListener{
		leaveSet: make(Set, maxCount),
		enterSet: make(Set, maxCount),
		moveSet:  make(Set, maxCount),
	}
}
func (l *emptyListener) clear() {
	l.enterSet.Clear()
	l.moveSet.Clear()
	l.leaveSet.Clear()
}

func (l *emptyListener) OnEnter(id ID, s Set) {
	l.enterSet.Clear()
	for k := range s {
		l.enterSet[k] = struct{}{}
	}
	//fmt.Println("OnEnter", id, s)
}

func (l *emptyListener) OnMove(id ID, s Set) {

	l.moveSet.Clear()
	for k := range s {
		l.moveSet[k] = struct{}{}
	}
	//fmt.Println("OnMove", id, s)
}

func (l *emptyListener) OnLeave(id ID, s Set) {
	l.leaveSet.Clear()
	for k := range s {
		l.leaveSet[k] = struct{}{}
	}
	//fmt.Println("OnLeave", id, s)
}

type testPoint struct {
	id ID
	x  float32
	y  float32
}

type points []*testPoint

func (ps points) View(idx int, rx, ry float32) Set {
	ret := make(Set, 0)
	p := ps[idx]
	for i := 0; i < len(ps); i++ {
		if i == idx {
			continue
		}
		pp := ps[i]
		if Abs(pp.x-p.x) <= rx && Abs(pp.y-p.y) <= ry {
			ret[pp.id] = struct{}{}
		}
	}
	return ret
}

func TestAOIAdd(t *testing.T) {
	l := newEmptyListener()
	m := NewManager(viewRangeX, viewRangeY, maxCount, l)
	ps := make(points, 0, maxCount)

	for i := 0; i < maxCount; i++ {

		id := ID(i)
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)
		p := &testPoint{id, x, y}
		ps = append(ps, p)
		s := ps.View(i, viewRangeX, viewRangeY)

		if !m.Enter(p.id, p.x, p.y) {
			t.FailNow()
		}

		if len(l.enterSet) != len(s) {
			t.FailNow()
		}
		for k := range s {
			if _, ok := l.enterSet[k]; !ok {
				t.FailNow()
			}
		}

		l.clear()
	}
}

func TestAOIMove(t *testing.T) {
	l := newEmptyListener()
	m := NewManager(viewRangeX, viewRangeY, maxCount, l)
	ps := make(points, 0, maxCount)

	for i := 0; i < maxCount; i++ {
		id := ID(i)
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)
		p := &testPoint{id, x, y}

		m.Enter(id, x, y)

		ps = append(ps, p)
	}

	for i := 0; i < maxCount; i++ {

		id := rand.Int() % maxCount
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)

		leaveSet := ps.View(id, viewRangeX, viewRangeY)
		ps[id].x, ps[id].y = x, y
		enterSet := ps.View(id, viewRangeX, viewRangeY)

		moveSet := leaveSet.Inersect(enterSet)
		enterSet.Trim(moveSet)
		leaveSet.Trim(moveSet)

		if !m.Move(ID(id), x, y) {
			t.FailNow()
		}

		if !enterSet.Equal(l.enterSet) {
			t.FailNow()
		}

		if !moveSet.Equal(l.moveSet) {
			t.FailNow()
		}
		if !leaveSet.Equal(l.leaveSet) {
			t.FailNow()
		}

		l.clear()
	}
}

func TestAOIRange(t *testing.T) {
	//l := newEmptyListener()
	m := NewManager(viewRangeX, viewRangeY, maxCount, nil)
	ps := make(points, 0, maxCount)

	for i := 0; i < maxCount; i++ {
		id := ID(i)
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)
		p := &testPoint{id, x, y}

		m.Enter(id, x, y)

		ps = append(ps, p)
	}

	for i := 0; i < maxCount; i++ {

		id := rand.Int() % maxCount
		//x := float32(rand.Int() % mapX)
		//y := float32(rand.Int() % mapY)

		s := ps.View(id, viewRangeX, viewRangeY)

		rs := make(Set, 0)
		m.GetRange(ID(id), rs)
		if !s.Equal(rs) {
			t.FailNow()
		}
	}
}

func TestAOILeave(t *testing.T) {
	l := newEmptyListener()
	m := NewManager(viewRangeX, viewRangeY, maxCount, l)
	ps := make(points, 0, maxCount)

	for i := 0; i < maxCount; i++ {
		id := ID(i)
		x := float32(rand.Int() % mapX)
		y := float32(rand.Int() % mapY)
		p := &testPoint{id, x, y}

		m.Enter(id, x, y)

		ps = append(ps, p)
	}

	for i := 0; i < maxCount; i++ {

		id := rand.Int() % maxCount
		p := ps[id]
		//x := float32(rand.Int() % mapX)
		//y := float32(rand.Int() % mapY)

		s := ps.View(id, viewRangeX, viewRangeY)

		if !m.Leave(ID(id)) {
			t.FailNow()
		}
		if !s.Equal(l.leaveSet) {
			t.FailNow()
		}
		if !m.Enter(ID(p.id), p.x, p.y) {
			t.FailNow()
		}

		l.clear()
	}
}

func BenchmarkAdd(b *testing.B) {
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
		m.Leave(id)
		b.StartTimer()
		m.Enter(id, x, y)
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
