package base

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"imooc.com/infra"
)

var props kvs.ConfigSource

//对外暴露的props方法
func Props() kvs.ConfigSource {
	return props
}

//初始化
type PropsStarter struct {
	infra.BaseStarter
}

//加载初始化配置文件
func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ini.NewIniFileCompositeConfigSource("config.ini")
}
