package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	PageResult struct {
		List     interface{} `json:"list"`
		Total    int64       `json:"total"`
		Page     int         `json:"page"`
		PageSize int         `json:"page_size"`
	}
	Response struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
		Msg  string      `json:"msg"`
	}
	RspSystemStatistics struct {
		NormalNodeCount int64 `json:"normal_node_count"`  //正常节点数量
		FailNodeCount   int64 `json:"fail_node_count"`    //不正常节点数量
		JobExcCount     int64 `json:"job_exc_count"`      //任务执行总数
		JobRunningCount int64 `json:"job_running_count"`  //任务正在执行总数
		JobExcFailCount int64 `json:"job_exc_fail_count"` //任务执行失败总数
	}
)

const (
	SUCCESS = 200
	ERROR   = 1000

	ErrorRequestParameter = 1001
	ErrorJobFormat        = 1002
	ErrorAutoAllocateNode = 1003
	ErrorJwtInvalid       = 1011
	ErrorTokenGenerate    = 1014
	ErrorLoginStatusSet   = 1015
	ErrorRegisterFormat   = 1016
	ErrorRegister         = 1017
	ErrorUserNameExist    = 1018
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "operation success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "operation success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func FailWithMessage(code int, message string, c *gin.Context) {
	Result(code, map[string]interface{}{}, message, c)
}
func FailWithCode(code int, c *gin.Context) {
	Result(code, map[string]interface{}{}, "operation failed", c)
}
func FailWithDetailed(code int, data interface{}, message string, c *gin.Context) {
	Result(code, data, message, c)
}
