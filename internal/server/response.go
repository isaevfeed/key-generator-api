package server

import "encoding/json"

type KeyContent struct {
	Key string `json:"key"`
}

type Response struct {
	Code    int32       `json:"code"`
	Content *KeyContent `json:"content"`
}

func MakeResponseStringify(StatusCode int32, Key string) string {
	res := &Response{
		Code:    StatusCode,
		Content: &KeyContent{Key},
	}

	jsonRes, _ := json.Marshal(res)

	return string(jsonRes)
}
