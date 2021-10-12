package utils

type Reply struct {
	Code int `json:"code"`
	Result interface{} `json:"result"`
}

func SuccessResponse(data interface{}) *Reply {
	return &Reply{
		Code: 0,
		Result: data,
	}
}

func ErrorResponse(code int, data interface{}) *Reply {
	return &Reply{
		Code: code,
		Result: data,
	}
}