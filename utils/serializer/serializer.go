package serializer

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	Logger = log.New(os.Stderr, "ERROR -", 18)
)

// Response 基础序列化器
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"` //omitempty 如果为零值 转json时就忽略
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Size    int64       `json:"size,omitempty"`
}

// Err 通用错误处理
func Err(errCode int, msg string, err error) *Response {
	res := Response{
		Code:    errCode,
		Message: msg,
	}

	if err != nil {
		Logger.Println(err.Error())
		// 生产环境隐藏底层报错
		if gin.Mode() != gin.ReleaseMode {
			res.Error = err.Error()
		}
	}
	return &res
}
