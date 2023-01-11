package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	max_size = 922337

	logger = log.New(os.Stdout, "info -", 18)
)

/*
图
属性：
Vertices  顶点总数
NumVertices  顶点总数
NumEdge  边总数

函数：
AddVertex(vert)  将顶点vert 加入图中，addEdge(fromVert,toVert) 添加有向边
AddEdge(fromVert, toVert, weight)  添加带优先级的有向边
GetVertex(vKey)  查找名称为 vKey 的顶点，getVertices() 返回图中所有顶点列表
Contains  按照vert in graph 的语句形式，返回顶点是否在图中 true/false

*/

// 图，包括顶点信息 顶点数 边数
type Graph struct {
	Vertices    map[string]*Vertex
	NumVertices int
	NumEdge     int
}

func MakeGraph() *Graph {
	return &Graph{
		Vertices: make(map[string]*Vertex),
		NumEdge:  0,
	}
}

// 添加顶点,  并返回被添加的顶点
func (gp *Graph) AddVertex(key string) *Vertex {
	if gp.NumEdge >= max_size {
		panic("too much edge or vertex.")
	}
	gp.NumVertices += 1
	NewVertex := MakeVertex(key)
	gp.Vertices[key] = NewVertex
	return NewVertex
}

// 获取顶点
func (gp *Graph) GetVertex(key string) *Vertex {
	return gp.Vertices[key]
}

// 顶点是否存在该图中
func (gp *Graph) Contains(key string) bool {
	return gp.Vertices[key] != nil
}

// 顶点是否在图中
func (gp *Graph) AddEdge(start, end string, cost int) {

	var nv *Vertex
	if !gp.Contains(start) {
		nv = gp.AddVertex(start)
	}
	if !gp.Contains(end) {
		nv = gp.AddVertex(end)
	}

	logger.Printf("added new vertex:%v\n", nv)
	// 使用顶点结构体的方法，添加邻接边
	rst := gp.Vertices[start].AddNeighbor(gp.Vertices[end], cost)
	//
	if rst {
		gp.NumEdge += 1
	}

}

// 获取顶点名称的 列表集, 和值集
func (gp *Graph) GetVertices() ([]string, []*Vertex) {
	var keys []string
	var vers []*Vertex
	for k, v := range gp.Vertices {
		keys = append(keys, k)
		vers = append(vers, v)
	}
	return keys, vers
}

// 顶点
type Vertex struct {
	Id          string
	ConnectedTo map[string]int
	Color       string
	Dist        int
	Pred        string
	Disc        int
	Fin         int
}

func MakeVertex(idNum string) *Vertex {
	return &Vertex{
		Id:          idNum,
		ConnectedTo: make(map[string]int),
		Color:       "white",
		Dist:        max_size,
		Pred:        "",
		Disc:        0,
		Fin:         0,
	}
}

/*
添加顶点邻居

顶点对象一个字典的key
设置边和权重信息
*/
func (vt *Vertex) AddNeighbor(nbr *Vertex, weight int) bool {

	if nbr == nil {
		return false
	}
	if vt.ConnectedTo[nbr.Id] != 0 {
		return false
	}
	vt.ConnectedTo[nbr.Id] = weight

	logger.Printf("vertex:%v of nbr is:%v is nil, set weight fail:%v \n", vt, nbr, weight)
	return true
}

/*
移除该顶点的一个邻接顶点
*/
func (vt *Vertex) DelNeighbor(nbr *Vertex) {
	delete(vt.ConnectedTo, nbr.Id)
}

/*
设置颜色
*/
func (vt *Vertex) SetColor(color string) {
	vt.Color = color
}

/*
设置距离
*/
func (vt *Vertex) SetDistance(d int) {
	vt.Dist = d
}

/*
设置前置
*/
func (vt *Vertex) SetPred(p string) {
	vt.Pred = p
}

/*
设置时间
*/
func (vt *Vertex) SetDiscovery(dtime int) {
	vt.Disc = dtime
}

/*
设置完成点
*/
func (vt *Vertex) SetFinish(ftime int) {
	vt.Fin = ftime
}

/*
获取最小
*/
func (vt *Vertex) GetFinish(ftime int) int {
	return vt.Fin
}

func (vt *Vertex) GetDiscovery() int {
	return vt.Disc
}

func (vt *Vertex) GetPred() string {
	return vt.Pred
}

func (vt *Vertex) GetDistance() int {
	return vt.Dist
}

func (vt *Vertex) GetColor() string {
	return vt.Color
}

// 返回该顶点的连接的全部key
/*
@param item = key 返回全部key
*/
func (vt *Vertex) GetConnections(item string) []int {
	var keys []string
	var vts []int
	for k, v := range vt.ConnectedTo {
		keys = append(keys, k)
		vts = append(vts, v)
	}
	return vts

}

// 返回该顶点边的 全部 权重信息
func (vt *Vertex) GetAllWeights(nbrKey string) []int {
	var weights []int
	for _, w := range vt.ConnectedTo {
		weights = append(weights, w)
	}
	return weights
}

func (vt *Vertex) GetWeight(nbrKey string) int {
	return vt.ConnectedTo[nbrKey]
}

func (vt *Vertex) GetStr(nbr string) string {
	return fmt.Sprintf("%v:color:%v:disc:%v:fin:%v:dist:%v:connect and weight info:%v:pred:%v",
		vt.Id, vt.Color, vt.Disc, vt.Fin, vt.Dist, vt.ConnectedTo, vt.Pred)
}

func (vt *Vertex) GetId() string {
	return vt.Id
}

func TestSetup(t *testing.T) {

	Gp := MakeGraph()
	Gp.AddVertex(fmt.Sprintf("%v", 1))
	Gp.AddVertex(fmt.Sprintf("%v", 2))
	//设置顶点1的 成本 2
	Gp.Vertices[fmt.Sprintf("%v", 1)].SetDistance(2)
	//获取一条边的成本
	for i := 1; i < 3; i++ {
		t.Logf("vertices 1 connections:%v \n", Gp.Vertices[fmt.Sprintf("%v", i)].GetConnections("key"))

	}

	//顶点1 添加一条边 到顶点2 成本1
	Gp.AddEdge(fmt.Sprintf("%v", 1), fmt.Sprintf("%v", 2), 1)
	//重复添加一条边 到顶点2
	Gp.AddEdge(fmt.Sprintf("%v", 1), fmt.Sprintf("%v", 2), 1)
	//添加一条边 到顶点3  成本2
	Gp.AddEdge(fmt.Sprintf("%v", 1), fmt.Sprintf("%v", 3), 2)
	//顶点1 添加一条边 到顶点4  成本1
	Gp.AddEdge(fmt.Sprintf("%v", 1), fmt.Sprintf("%v", 4), 1)
	//顶点2 添加一条边 到顶点3  成本2
	Gp.AddEdge(fmt.Sprintf("%v", 2), fmt.Sprintf("%v", 3), 2)
	//顶点2 添加一条边 到顶点4  成本3
	Gp.AddEdge(fmt.Sprintf("%v", 2), fmt.Sprintf("%v", 4), 3)
	//顶点3 添加一条边 到顶点4  成本4
	Gp.AddEdge(fmt.Sprintf("%v", 3), fmt.Sprintf("%v", 4), 4)
	//顶点2 添加一条边 到顶点1  成本2
	Gp.AddEdge(fmt.Sprintf("%v", 2), fmt.Sprintf("%v", 1), 2)

	//设置顶点3 距离为4
	Gp.Vertices[fmt.Sprintf("3")].SetDistance(4)
	//设置顶点2 距离为3
	Gp.Vertices[fmt.Sprintf("2")].SetDistance(3)
	t.Logf("顶点总数:%v, 边总数:%v, \n顶点1的信息:%v", Gp.NumVertices, Gp.NumEdge, Gp.GetVertex(fmt.Sprintf("%v", 1)))

	for k, v := range Gp.Vertices {
		t.Logf("vertex:%#v edge:%#v \n", k, v)
	}

	if Gp.NumVertices != 4 || Gp.NumEdge != 7 {
		t.Fatalf("vertices number should be 4, have:%v, edge number should be 7   have:%v \n", Gp.NumVertices, Gp.NumEdge)
	}

	vert := Gp.GetVertex(fmt.Sprintf("1"))

	new_color := "blcak"
	vert.SetColor(new_color)

	if nc := vert.GetColor(); nc != new_color {
		t.Fatalf("color set fail hope:%v, have:%v \n", new_color, nc)
	}

	//删除
	vert.DelNeighbor(Gp.GetVertex(fmt.Sprintf("2")))
	Gp.NumEdge -= 1
	if Gp.NumVertices != 4 || Gp.NumEdge != 6 {
		t.Fatalf("vertices number should be 4, have:%v, edge number should be 7   have:%v \n", Gp.NumVertices, Gp.NumEdge)
	}

}
