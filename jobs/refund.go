package jobs

import (
	"resk/core/envelope"
	"resk/infra"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/sirupsen/logrus"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker time.Ticker
	mutex  *redsync.Mutex
}

func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	// 初始化定时器
	duration := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute)
	r.ticker = *time.NewTicker(duration)
	// 初始化 redis 分布式锁
	timeout := ctx.Props().GetDurationDefault("redis.timeout", 20*time.Second)
	addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")
	client := goredislib.NewClient(&goredislib.Options{
		Addr:        addr,
		IdleTimeout: timeout,
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	mutexname := "locak:RefundExpiredJob"
	r.mutex = rs.NewMutex(mutexname)
}

func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			c := <-r.ticker.C
			err := r.mutex.Lock()
			if err != nil {
				logrus.Debug("已经有节点在运行该任务了")
			} else {
				logrus.Info("过期红包退款开始...", c)
				// 红包过期退款逻辑
				domain := envelope.ExpiredEnvelopeDomain{}
				domain.Expired()
			}
		}
	}()
}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}
