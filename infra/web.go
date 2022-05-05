package infra

var apiInitializerRegister *InitializerRegister = new(InitializerRegister)

// 注册 web api 初始化对象
func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebStarter struct {
	BaseStarter
}

func (w *WebStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializers() {
		v.Init()
	}
}
