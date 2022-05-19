package base

import (
	"fmt"
	"resk/infra"
	"sync"

	"github.com/tietang/props/kvs"
)

var props kvs.ConfigSource
var systemAccountOnce sync.Once

func Props() kvs.ConfigSource {
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("初始化配置成功")
	GetSystemAccount()
}

type SystemAccount struct {
	UserId    string
	Username  string
	AccountNo string
}

// 系统账户配置
var systemAccount *SystemAccount

func GetSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		systemAccount = new(SystemAccount)
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	link, err := Props().Get("envelope.link")
	if err != nil {
		panic(err)
	}
	return link
}

func GetEnvelopeActivityDomain() string {
	domain, err := Props().Get("envelope.domain")
	if err != nil {
		panic(err)
	}
	return domain
}
