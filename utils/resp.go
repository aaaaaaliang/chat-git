package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 结构用于构造响应数据的通用格式
type H struct {
	Code  int         // 响应状态码
	Msg   string      // 响应消息
	Data  interface{} // 响应的数据
	Rows  interface{} // 响应的数据列表（通常用于列表数据）
	Total interface{} // 列表数据的总数（通常用于分页）
}

// Resp 用于生成通用响应数据
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

// RespList 用于生成列表数据的响应
func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

// RespFail 用于生成失败的响应
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

// RespOK 用于生成成功的响应
func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

// RespOKList 用于生成成功的响应
func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
