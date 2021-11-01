package utils

func SuccessResponse(data interface{}) interface{} {
	return map[string]interface{}{
		"code": 0,
		"result": data,
	}
}

func ErrorResponse(code int, data interface{}) interface{} {
	return map[string]interface{}{
		"code": code,
		"result": data,
	}
}