package infra

// 初始化接口
type Initializer interface {
	// 用于对象实例化后的初始化操作
	Init()
}

type InitializerRegister struct {
	Initializers []Initializer
}

func (r *InitializerRegister) Register(ai Initializer) {
	r.Initializers = append(r.Initializers, ai)
}
