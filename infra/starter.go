package infra

import "github.com/tietang/props/kvs"

const KeyProps = "_conf"

type StarterContext map[string]any

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有被初始化")
	}
	return p.(kvs.ConfigSource)
}

type IStarter interface {
	// 1. 系统启动，初始化一些资源
	Init(StarterContext)
	// 2. 系统基础资源安装
	Setup(StarterContext)
	// 3. 启动基础资源
	Start(StarterContext)
	// 启动器是否阻塞
	StartBlocking() bool
	// 4. 资源停止和销毁
	Stop(StarterContext)
}

// 启动器注册器
type starterRegister struct {
	starters []IStarter
}

// 注册启动器
func (r *starterRegister) Register(s IStarter) {
	r.starters = append(r.starters, s)
}

func (r *starterRegister) AllStarters() []IStarter {
	return r.starters
}

var StarterRegister *starterRegister = new(starterRegister)

func Register(s IStarter) {
	StarterRegister.Register(s)
}

type BaseStarter struct {
}

func (b *BaseStarter) Init(ctx StarterContext) {

}

func (b *BaseStarter) Setup(ctx StarterContext) {

}

func (b *BaseStarter) Start(ctx StarterContext) {

}

func (b *BaseStarter) StartBlocking() bool {
	return false
}

func (b *BaseStarter) Stop(ctx StarterContext) {

}
