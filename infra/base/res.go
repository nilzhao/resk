package base

type ResCode int

const (
	ResCodeOk                 = 1000
	ResCodeValidationError    = 2000
	ResCodeRequestParamsError = 2100
	ResCodeInnerServerError   = 5000
	ResCodeBizError           = 6000
)

type Res struct {
	Code    ResCode `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data"`
}
