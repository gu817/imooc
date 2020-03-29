package gorpc

import (
	"errors"
	"git.imooc.com/wendell1000/infra/lb"
	"github.com/sirupsen/logrus"
	"net/rpc"
)

type GoRpcClient struct {
	Apps *lb.Apps
}

func (g *GoRpcClient) Call(serviceId, serviceMethod string,
	in interface{}, out interface{}) error {
	//通过微服务名称从本地服务注册表中查询应用和应用实例列表
	app := g.Apps.Get(serviceId)
	if app == nil {
		return errors.New("没有可用的微服务应用，应用名称：" + serviceId + ",请求：" + serviceMethod)
	}
	//通过负载均衡算法从应用实例列表中选择一个实例
	ins := app.Get(serviceMethod)
	if ins == nil {
		return errors.New("没有可用的应用实例，应用名称：" + serviceId + ",请求：" + serviceMethod)
	}
	//选择的实例IP和端口

	addr := ins.Metadata["rpcAddr"]

	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer c.Close()
	err = c.Call(serviceMethod, in, out)
	if err != nil {
		logrus.Error(err)
	}
	return err
}
