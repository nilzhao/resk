package base

import (
	"resk/infra"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/sirupsen/logrus"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	return validate
}

func Translate() ut.Translator {
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	cn := zh.New()
	uni := ut.New(cn, cn)
	trans, found := uni.GetTranslator("zh")
	if found {
		translator = trans
		err := zhTrans.RegisterDefaultTranslations(validate, translator)
		logrus.Error(err)
	} else {
		logrus.Error("Not found translator: zh")
	}
}
