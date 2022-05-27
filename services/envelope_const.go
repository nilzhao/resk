package services

const (
	// 祝福语
	DefaultBlessing = "恭喜发财!"
)

// 订单类型:发布单、退款单
type OrderType int8

const (
	OrderTypeSending OrderType = 1
	OrderTypeRefund  OrderType = 2
)

// 支付状态: 未支付、支付中、已支付、支付失败
type PayStatus int8

const (
	PayNothing PayStatus = 1
	Paying     PayStatus = 2
	Payed      PayStatus = 3
	PayFailure PayStatus = 4

	RefundNothing PayStatus = 6
	Refunding     PayStatus = 7
	Refunded      PayStatus = 8
	RefundFailure PayStatus = 9
)

// 红包订单状态: 创建 发布 过期 失效
type OrderStatus int8

const (
	OrderCreate    OrderStatus = 1
	OrderActivated OrderStatus = 2
	OrderExpired   OrderStatus = 3
	OrderDisabled  OrderStatus = 4

	// 退款成功
	OrderExpiredSuccess OrderStatus = 5
	// 退款失败
	OrderExpiredFailed OrderStatus = 6
)

// 活动状态: 创建 激活 过期 失效
type ActivityStatus int8

const (
	ActivityCreate    ActivityStatus = 1
	ActivityActivated ActivityStatus = 2
	ActivityExpired   ActivityStatus = 3
	ActivityDisabled  ActivityStatus = 4
)

// 红包类型：普通红包，碰运气红包
type EnvelopeType int

const (
	GeneralEnvelopeType EnvelopeType = 1
	LuckyEnvelopeType   EnvelopeType = 2
)
