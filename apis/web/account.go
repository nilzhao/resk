package web

import (
	"resk/infra"
	"resk/infra/base"
	"resk/services"

	"github.com/kataras/iris/v12"
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct{}

func (a *AccountApi) Init() {
	groupRouter := base.Iris().Party("/api/v1/account")
	groupRouter.Post("/create", createHandler)
	groupRouter.Post("/transfer", transferHandler)
	groupRouter.Get("/envelope/get", getEnvelopeAccountHandler)
	groupRouter.Get("/get", getAccountHandler)
}

func createHandler(ctx iris.Context) {
	// 获取请求参数
	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	res := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		res.Code = base.ResCodeRequestParamsError
		res.Message = err.Error()
		ctx.JSON(&res)
		return
	}
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		res.Code = base.ResCodeInnerServerError
		res.Message = err.Error()
	}
	res.Data = dto
	ctx.JSON(&res)
}

func transferHandler(ctx iris.Context) {
	// 获取请求参数
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	res := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		res.Code = base.ResCodeRequestParamsError
		res.Message = err.Error()
		ctx.JSON(&res)
		return
	}
	service := services.GetAccountService()
	status, err := service.Transfer(account)
	if err != nil {
		res.Code = base.ResCodeInnerServerError
		res.Message = err.Error()
	}
	res.Data = status
	if status != services.TransferStatusSuccess {
		res.Code = base.ResCodeBizError
		res.Message = err.Error()
	}
	ctx.JSON(&res)
}

func getEnvelopeAccountHandler(ctx iris.Context) {
	userId := ctx.URLParam("userId")
	res := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		res.Code = base.ResCodeRequestParamsError
		res.Message = "用户 ID 不能为空"
		ctx.JSON(&res)
		return
	}
	service := services.GetAccountService()
	account := service.GetAccountByUserId(userId)
	res.Data = account
	ctx.JSON(&res)
}

func getAccountHandler(ctx iris.Context) {
	accountNo := ctx.URLParam("accountNo")
	res := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		res.Code = base.ResCodeRequestParamsError
		res.Message = "账户编号不能为空"
		ctx.JSON(&res)
		return
	}
	service := services.GetAccountService()
	account := service.GetAccount(accountNo)
	res.Data = account
	ctx.JSON(&res)
}
