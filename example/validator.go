package main

import (
	"fmt"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       uint8  `validate:"gte=0,lte=150"`
	Email     string `validate:"required,email"`
}

var validate *validator.Validate
var uni *ut.UniversalTranslator

func main() {
	zh := zh.New()
	uni = ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	validate = validator.New()
	zhTrans.RegisterDefaultTranslations(validate, trans)

	user := &User{
		FirstName: "zhao",
		LastName:  "qi",
		Age:       151,
		Email:     "123.com",
	}
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			fmt.Println(errs.Translate(trans))

		}
	}
}
