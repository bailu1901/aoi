package aoi

import "fmt"

// Set 集合
type Set map[ID]struct{}

// Clear 清除
func (s Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// EnterCallback 进入回掉
type EnterCallback func(ID, Set)

// MoveCallback 移动回掉
type MoveCallback func(ID, Set)

// LeaveCallback 离开回掉
type LeaveCallback func(ID, Set)

// ID id
type ID int32

type node struct {
	id ID
	x  float32
	y  float32

	prevX *node
	nextX *node
	prevY *node
	nextY *node
}

func newNode(id ID, x, y float32) *node {
	return &node{
		id: id,
		x:  x,
		y:  y,
	}
}
func (n *node) String() string {
	ret := ""
	cur := n
	for nil != cur {
		ret += fmt.Sprintf("[%d(%f,%f)]", cur.id, cur.x, cur.y)
		cur = cur.nextX
	}
	return ret
}

func (n *node) BreakX() {
	n.prevX.nextX, n.nextX.prevX = n.nextX, n.prevX
	n.prevX, n.nextX = nil, nil
}
func (n *node) BreakY() {
	n.prevY.nextY, n.nextY.prevY = n.nextY, n.prevY
	n.prevY, n.nextY = nil, nil
}
func (n *node) InsertXAfter(ohter *node) {
	ohter.prevX = n
	ohter.nextX = n.nextX
	n.nextX.prevX = ohter
	n.nextX = ohter
}

func (n *node) InsertXBefore(ohter *node) {
	ohter.nextX = n
	ohter.prevX = n.prevX
	n.prevX.nextX = ohter
	n.prevX = ohter
}
func (n *node) InsertYAfter(ohter *node) {
	ohter.prevY = n
	ohter.nextY = n.nextY
	n.nextY.prevY = ohter
	n.nextY = ohter
}

func (n *node) InsertYBefore(ohter *node) {
	ohter.nextY = n
	ohter.prevY = n.prevY
	n.prevY.nextY = ohter
	n.prevY = ohter
}

// Abs 绝对值
func Abs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

// Manager Manager
type Manager struct {
	objs map[ID]*node

	head *node
	tail *node

	rangeX float32
	rangeY float32

	leaveSet Set // 通知离开的集合
	enterSet Set // 通知进入的集合
	moveSet  Set // 通知移动的集合

	enterCb EnterCallback
	moveCb  MoveCallback
	leaveCb LeaveCallback
}

// NewManager AOI管理器
func NewManager(viewX, viewY float32, capcity int,
	ecb EnterCallback,
	mcb MoveCallback,
	lcb LeaveCallback) *Manager {
	mgr := &Manager{
		objs:     make(map[ID]*node, capcity),
		rangeX:   viewX,
		rangeY:   viewY,
		enterCb:  ecb,
		moveCb:   mcb,
		leaveCb:  lcb,
		leaveSet: make(Set, capcity),
		enterSet: make(Set, capcity),
		moveSet:  make(Set, capcity),
	}
	mgr.head = newNode(-99999999, -99999999, -99999999)
	mgr.tail = newNode(99999999, 99999999, 99999999)

	mgr.head.nextX = mgr.tail
	mgr.head.nextY = mgr.tail
	mgr.tail.prevX = mgr.head
	mgr.tail.prevY = mgr.head

	return mgr
}

// GetRange 获得视野内的对象
func (mgr *Manager) GetRange(id ID, ret Set) {
	n, ok := mgr.objs[id]
	if !ok {
		return
	}

	// 向前
	cur := n.prevX
	for nil != cur && cur != mgr.head {
		if cur.x < n.x-mgr.rangeX {
			break
		}
		if Abs(cur.y-n.y) <= mgr.rangeY {
			ret[cur.id] = struct{}{}
		}
		cur = cur.prevX
	}

	// 向后
	cur = n.nextX
	for nil != cur && cur != mgr.tail {
		if cur.x > n.x+mgr.rangeX {
			break
		}
		if Abs(cur.y-n.y) <= mgr.rangeY {
			ret[cur.id] = struct{}{}
		}
		cur = cur.nextX
	}

	return
}

// Add 添加节点
func (mgr *Manager) Add(id ID, x, y float32) bool {
	if _, ok := mgr.objs[id]; ok {
		return false
	}

	// 新节点
	newNode := newNode(id, x, y)

	// 遍历x轴，插入合适位置，同时把需要通知进入的人放入集合
	cur := mgr.head.nextX
	intertX := false
	for nil != cur {
		if !intertX && (cur == mgr.tail || cur.x > newNode.x) {
			cur.InsertXBefore(newNode)
			intertX = true
		}

		if cur == mgr.tail {
			break
		}

		diffX := cur.x - newNode.x
		if diffX > mgr.rangeX {
			break
		}

		// X轴在范围内，Y轴的也在范围，一次就找到需要通知的单位了
		if Abs(diffX) <= mgr.rangeX && Abs(cur.y-newNode.y) <= mgr.rangeY {
			mgr.enterSet[cur.id] = struct{}{}
		}
		cur = cur.nextX
	}

	// 遍历Y轴，插入合适位置
	cur = mgr.head.nextY
	for nil != cur {
		if cur == mgr.tail || cur.y > newNode.y {
			cur.InsertYBefore(newNode)
			break
		}
		cur = cur.nextY
	}

	mgr.objs[id] = newNode

	// 通知回掉
	mgr.processEvent(id)

	return true
}

// Move 移动
func (mgr *Manager) Move(id ID, x, y float32) bool {
	n, ok := mgr.objs[id]
	if !ok {
		return false
	}

	// 先获得
	mgr.GetRange(id, mgr.moveSet)

	n.x, n.y = x, y

	inRangeX := false
	if n.x < n.prevX.x {
		// 向前
		cur := n.prevX
		for cur != mgr.head {
			if n.x > cur.x {
				break
			}
			cur = cur.prevX
		}
		n.BreakX()
		// 插在这个节点的后面
		cur.InsertXAfter(n)
	} else if n.x > n.nextX.x {
		// 向后
		cur := n.nextX
		for cur != mgr.tail {
			if n.x < cur.x {
				break
			}
			cur = cur.nextX
		}
		n.BreakX()
		// 插在这个节点的后面
		cur.InsertXBefore(n)
	} else {
		inRangeX = true
	}

	inRangeY := false
	if n.y < n.prevY.y {
		// 向前
		cur := n.prevY
		for cur != mgr.head {
			if n.y > cur.y {
				break
			}
			cur = cur.prevY
		}
		n.BreakY()
		// 插在这个节点的后面
		cur.InsertYAfter(n)
	} else if n.y > n.nextY.y {
		// 向后
		cur := n.nextY
		for cur != mgr.tail {
			if n.y < cur.y {
				break
			}
			cur = cur.nextY
		}
		n.BreakY()
		// 插在这个节点的后面
		cur.InsertYBefore(n)
	} else {
		inRangeY = true
	}

	if !(inRangeX && inRangeY) { // 这次移动没出X轴的范围也没出Y轴的范围
		mgr.GetRange(id, mgr.enterSet)
		// old和new的交集就是move，剩下的是离开
		for k := range mgr.moveSet {
			if _, ok := mgr.enterSet[k]; ok {
				delete(mgr.enterSet, k)
			} else {
				mgr.leaveSet[k] = struct{}{}
				delete(mgr.moveSet, k)
			}
		}
	}

	// 回掉
	mgr.processEvent(n.id)

	return true
}

// Leave 离开
func (mgr *Manager) Leave(id ID) {
	mgr.GetRange(id, mgr.leaveSet)
	mgr.processEvent(id)
}

// processEvent 处理事件
func (mgr *Manager) processEvent(id ID) {
	if len(mgr.enterSet) > 0 {
		mgr.enterCb(id, mgr.enterSet)
	}
	if len(mgr.moveSet) > 0 {
		mgr.moveCb(id, mgr.moveSet)
	}
	if len(mgr.leaveSet) > 0 {
		mgr.leaveCb(id, mgr.leaveSet)
	}

	mgr.enterSet.Clear()
	mgr.moveSet.Clear()
	mgr.leaveSet.Clear()
}
