/*
	管理 应用加载启动生命周期
*/

package infra

import "github.com/tietang/props/kvs"

type BootApplication struct {
	conf           kvs.ConfigSource
	starterContext StarterContext
}

func New(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterContext: StarterContext{},
	}
	b.starterContext[KeyProps] = conf
	return b
}

func (b *BootApplication) Start() {
	// 1、初始化 starter
	b.init()
	// 2、安装 starter
	b.setup()
	// 3、启动 starter
	b.start()
}

func (b *BootApplication) init() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	for index, starter := range StarterRegister.AllStarters() {
		if starter.StartBlocking() {
			// 如果是最后一个 starter，是可以阻塞的，直接启动并阻塞
			if index+1 == len(StarterRegister.AllStarters()) {
				starter.Start(b.starterContext)
				// 如果不是，使用goroutine，防止阻塞后续 starter
			} else {
				go starter.Start(b.starterContext)
			}
		} else {
			starter.Start(b.starterContext)

		}
	}
}
