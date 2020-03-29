package infra

import "github.com/tietang/props/kvs"

//加载所有程序启动加载的生命周期
type BootApplication struct {
	con            kvs.ConfigSource
	starterContext StarterContext
}

//管理所有的生命周期

func (b *BootApplication) Start() {
	//初始化starter
	b.init()
	//	安装starter
	//启动starter

}

func (b *BootApplication) init() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(b.starterContext)
	}
}
