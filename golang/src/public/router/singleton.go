package router

//Singleton 是单例模式类
type Singleton struct{}

var singletonInst *Mux

func init() {
	singletonInst = New()
}

func Instance() *Mux {
	return singletonInst
}
