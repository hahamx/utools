package queues

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Stack struct {
	size  int
	items []interface{}
}

func NewStack() *Stack {
	return &Stack{size: 10}
}

func (sk *Stack) Size() int {
	return sk.size
}
func (sk *Stack) IsFull() bool {
	return len(sk.items) >= sk.size
}

// 判空
func (sk *Stack) IsEmpty() bool {
	return len(sk.items) == 0
}

// 压入最近的数据
func (sk *Stack) Push(value interface{}) {
	if len(sk.items) >= sk.size {
		panic("under stack flow.")
	}
	sk.items = append(sk.items, value)
}

// 弹出最近的数据
func (sk *Stack) Pop() string {
	if sk.IsEmpty() {
		panic("under stack flow.")
	}
	lastOne := sk.items[len(sk.items)-1]
	sk.items = sk.items[:len(sk.items)-1]
	return ToString(lastOne)
}

// 查看最近的数据
func (sk *Stack) Peek() string {
	if sk.IsEmpty() {
		panic("under stack flow.")
	}
	ps := sk.items[len(sk.items)-1]
	return ToString(ps)
}

// 转换为 字符串
func ToString(inter interface{}) string {

	var key string
	if inter == nil {
		return key
	}
	switch inter.(type) {
	case float64:
		ft := inter.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := inter.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := inter.(int)
		key = strconv.Itoa(it)
	case int8:
		it := inter.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := inter.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := inter.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := inter.(int32)
		key = strconv.Itoa(int(it))
	case int64:
		it := inter.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := inter.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = inter.(string)
	case []byte:
		key = string(inter.([]byte))
	default:
		Nv, _ := json.Marshal(inter)
		key = string(Nv)
	}
	return key
}

func TestMain() {
	ns := NewStack()
	n := 10
	for i := 0; i < n; i++ {
		ns.Push(i)
	}

	fmt.Println("len after push:", len(ns.items))

	fmt.Println(ns.Pop())
	fmt.Println("len after pop:", len(ns.items))

	ns.Push("Hello")
	fmt.Println("anyone there:", ns.IsEmpty())
	fmt.Println("peek last:", ns.Peek())
	fmt.Println("len:", len(ns.items))

}
