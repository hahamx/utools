package main

import (
	"fmt"
	"time"
)

// 散列管理结构体
type HashStack struct {
	Size  int
	Slots []interface{}
	Data  []interface{}
}

func NewHashStack(size int) *HashStack {
	return &HashStack{
		Size:  size,
		Slots: make([]interface{}, size),
		Data:  make([]interface{}, size),
	}
}

/*
hashfunction方法采用了简单随机方法来实现散列函数，
而冲突解决则采用 线性探测“加1”再散列函数
@param key 求余获取散列位置
*/
func (hs *HashStack) HashFuncation(key int) int {
	return key % hs.Size
}

func (hs *HashStack) ReHash(old int) int {
	return (old + 1) % hs.Size
}

/*
放入
当Slots[hashPosiation]中的 key 不存在时。不冲突，则直接放入
否则
当存在key时，覆盖原有值，Slots[hashPosiation] key  存在 冲突 覆盖
另一个场景是存在一个不一样的key在当前位置，因此重新散列，以寻找一个新的位置，如果位置总是不为空且不是目标key则继续寻找。
找到的新的位置 不存在冲突则放入，否则覆盖。
*/
func (hs *HashStack) Put(key int, data interface{}) {
	hashPosiation := hs.HashFuncation(key)
	if hs.Slots[hashPosiation] == nil {
		hs.Slots[hashPosiation] = key
		hs.Data[hashPosiation] = data
	} else {
		if hs.Slots[hashPosiation] == key {
			hs.Data[hashPosiation] = data
		} else {
			nextSlot := hs.ReHash(hashPosiation)

			for hs.Slots[nextSlot] != nil && hs.Slots[nextSlot] != key {
				nextSlot = hs.ReHash(nextSlot)
			}

			if hs.Slots[nextSlot] == nil {
				hs.Slots[nextSlot] = key
				hs.Data[nextSlot] = data
			} else {
				hs.Data[nextSlot] = data //覆盖
			}
		}
	}
}

/*
从散列中获取散列值
首先标记散列值，标记查找起点。
由于 查找key直到遇到空槽，否则stop = true回到起点
*/
func (hs *HashStack) Get(key int) string {

	startSlot := hs.HashFuncation(key)

	var data interface{}
	stop := false
	found := false
	posi := startSlot
	for hs.Slots[posi] != nil && !found && !stop {
		if hs.Slots[posi] == key {
			found = true
			data = hs.Data[posi]
		} else {
			posi = hs.ReHash(posi)
			if posi == startSlot {
				stop = true
			}
		}

	}
	return ToString(data)

}

// 查看基础信息
func (hs *HashStack) BaseInfo() {
	fmt.Printf("散列插槽数:%v, 散列插槽值:%#v, 散列数据:%#v\n", hs.Size, hs.Slots, hs.Data)
	fmt.Printf("散列ks:%#v\n", hs.Slots)
	fmt.Printf("散列vs:%#v ", hs.Data)
}

func TestMain() {
	ht := NewHashStack(100)
	ht.Put(1, "c") //[1] = "c"
	ht.Put(2, "b")
	ht.Put(3, "a")
	fmt.Printf("散列信息:\n")
	ht.BaseInfo()
	ht.Put(1, 121)

	fmt.Printf("散列信息:\n")
	ht.BaseInfo()
	fmt.Printf("散列PUT:\n")
	ht.Put(35, 121)
	fmt.Printf("散列1位置:%v\n", ht.Get(1))
	ht.Put(1, 1331)
	fmt.Printf("散列121位置:%v\n", ht.Get(121))
	fmt.Printf("散列信息:\n")
	ht.BaseInfo()
	ht.Put(1, 1332)
	fmt.Printf("散列信息2:  \n")
	ht.BaseInfo()

	print(time.Now().UnixMilli())
}
