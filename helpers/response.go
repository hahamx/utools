package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Response 基础序列化器 //omitempty 如果为零值 转json时就忽略
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Size    int64       `json:"size,omitempty"`
}

// Err 通用错误处理 // 生产环境隐藏底层报错
func Err(errCode int, msg string, err error) *Response {
	res := Response{
		Code:    errCode,
		Message: msg,
	}

	if err != nil {

		if gin.Mode() != gin.ReleaseMode {
			res.Error = err.Error()
		}
	}
	return &res
}

// 成功 统一200, 如果code为0，将在这里转换，omitempty 如果为零值 转json时就忽略不必要的字段
func JSONs(ctx *gin.Context, code int, datas interface{}) {
	if code == 0 {
		code = 200
	}
	rsp := BuildResponse(code, datas)
	if rsp.Data == nil {
		rsp.Data = []struct{}{}
		rsp.Size = 0
		rsp.Total = 0
	}

	ctx.JSON(code, *rsp)
}

/*
从任意结构体去除data 和相关字段
*/
func BuildResponse(code int, datas interface{}) *Response {

	var sr = &Response{}

	if code != 0 && code != 200 {

		sr.Size = 0
		sr.Total = 0
		sr.Message = "Failed"
		sr.Error = fmt.Sprintf("%v", datas)
	} else {
		vb, err := json.Marshal(datas)
		if err != nil {
			fmt.Printf("Error On respnse json marshal datas:%v, err:%v\n", datas, err)
		}
		json.Unmarshal(vb, sr)
		sr.Message = "Success"
	}

	sr.Code = code
	return sr
}
