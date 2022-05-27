package resk

import (
	_ "resk/apis/web"
	_ "resk/core/account"
	_ "resk/core/envelope"
	"resk/infra"
	"resk/infra/base"
	"resk/jobs"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&infra.WebStarter{})
	infra.Register(&base.IrisServerStarter{})
}
