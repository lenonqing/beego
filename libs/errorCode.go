package libs

// ErrorType 错误信息
type ErrorType struct {
	StatusCode int
	ErrorCode  int
	ErrorMsg   string
}

// NewError 新建错误信息
func NewError(status, errcode int, errmsg string) (err ErrorType) {
	return ErrorType{
		StatusCode: status,
		ErrorCode:  errcode,
		ErrorMsg:   errmsg,
	}
}

// MakeError 创建错误信息
func MakeError(err error) (e ErrorType) {
	return ErrorType{
		StatusCode: 500,
		ErrorCode:  10000,
		ErrorMsg:   err.Error(),
	}
}

// 系统级错误
var (
	ErrorInvalidToken    = NewError(401, 10000, "Unauthorized")
	ErrorForbiddenAccess = NewError(403, 10000, "Forbidden Access")
)

// 逻辑错误
var (
	ErrorInternalError = NewError(200, 20000, "内部错误")
	ErrorMissParameter = NewError(200, 20001, "非法参数")
	ErrorNotFound      = NewError(200, 20002, "未找到对应数据")
)
