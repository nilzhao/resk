package jobs

import (
	"resk/core/envelope"
	"resk/infra"
	"time"

	"github.com/sirupsen/logrus"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker time.Ticker
}

func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	duration := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute)
	r.ticker = *time.NewTicker(duration)
}

func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			c := <-r.ticker.C
			logrus.Debug("过期红包退款开始...", c)
			// 红包过期退款逻辑
			domain := envelope.ExpiredEnvelopeDomain{}
			domain.Expired()
		}
	}()
}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}
