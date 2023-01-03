package helpers

import (
	"fmt"
	"time"
)

// 绑定 login登录 JSON表单
type Login struct {
	User string `form:"user" json:"user" xml:"user"  binding:"required"`
	// binding:"required" password 不能为空，binding:"-" password 可以为空
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// 不支持嵌套 结构体还带了 form
//type StructX struct {
//	X struct {} `form:"name_x"` // 有 form
//}

// 接口调用参数 如果有 field_a 将会传递到此
type StructA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

type UrlPerson struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

// 用于绑定html 复选框
type MyForm struct {
	Colors []string `form:"colors[]"`
}

func FormatAsDate(t time.Time) string {
	// 自定义模板 时间 map 格式
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func FormatAsString(ts string) string {
	// 自定义模板 时间 map 格式
	return fmt.Sprintf("%v", ts)
}
