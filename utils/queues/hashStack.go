package queues

import (
	"fmt"
)

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
hashfunction方法采用了简单求余方法来实现散列函数，
而冲突解决则采用 线性探测“加1”再散列函数
@param key 求余获取散列位置
*/
func (hs *HashStack) HashFuncation(key int) int {
	return key % hs.Size
}

func (hs *HashStack) ReHash(old int) int {
	return (old + 1) % hs.Size
}

// 插入
func (hs *HashStack) Put(key int, data interface{}) {
	hashPosiation := hs.HashFuncation(key)
	if hs.Slots[hashPosiation] == nil { //key 不存在。不冲突
		hs.Slots[hashPosiation] = key
		hs.Data[hashPosiation] = data
	} else {
		if hs.Slots[hashPosiation] == key { //key  存在 冲突
			hs.Data[hashPosiation] = data // 覆盖
		} else {
			nextSlot := hs.ReHash(hashPosiation)
			//一直处理冲突，通过再散列的方式，找到空槽或key
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

// 获取
func (hs *HashStack) Get(key int) string {

	//标记散列值，为查找起点
	startSlot := hs.HashFuncation(key)

	var data interface{}
	stop := false
	found := false
	posi := startSlot
	for hs.Slots[posi] != nil && !found && !stop { //查找key直到，空槽回到起点
		if hs.Slots[posi] == key {
			found = true
			data = hs.Data[posi]
		} else {
			posi = hs.ReHash(posi) //未能找到key，继续
			if posi == startSlot {
				stop = true //回到起点
			}
		}

	}
	return ToString(data)

}

func (hs *HashStack) ShowBaseInfo() {
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
	ht.ShowBaseInfo()
	ht.Put(1, 121)

	fmt.Printf("散列信息:\n")
	ht.ShowBaseInfo()
	fmt.Printf("散列PUT:\n")
	ht.Put(35, 121)
	fmt.Printf("散列1位置:%v\n", ht.Get(1))
	ht.Put(1, 1331)
	fmt.Printf("散列121位置:%v\n", ht.Get(121))
	fmt.Printf("散列信息:\n")
	ht.ShowBaseInfo()
	ht.Put(1, 1332)
	fmt.Printf("散列信息2:  \n")
	ht.ShowBaseInfo()
}
