package lb

import (
	"fmt"
	"github.com/tietang/go-eureka-client/eureka"
	"strings"
)

type Apps struct {
	Client *eureka.Client
}

func (a *Apps) Get(appName string) *App {
	var app *eureka.Application
	for _, a := range a.Client.Applications.Applications {
		if a.Name == strings.ToUpper(appName) {
			app = &a
			break
		}
	}
	if app == nil {
		return nil
	}

	na := &App{
		Name:      app.Name,
		Instances: make([]*ServerInstance, 0),
		lb:        &RoundRobinBalancer{},
	}
	for _, ins := range app.Instances {
		var port int
		if ins.SecurePort.Enabled {
			port = ins.SecurePort.Port
		} else {
			port = ins.Port.Port
		}
		si := &ServerInstance{
			InstanceId: ins.InstanceId,
			AppName:    appName,
			Status:     Status(ins.Status),
			Address:    fmt.Sprintf("%s:%d", ins.IpAddr, port),
			Metadata:   make(map[string]string),
		}
		si.Metadata["rpcAddr"] = fmt.Sprintf("%s:%s", ins.IpAddr, ins.Metadata.Map["rpcPort"])
		na.Instances = append(na.Instances, si)
	}
	return na
}

type App struct {
	Name      string
	Instances []*ServerInstance
	lb        Balancer
}

func (a *App) Get(key string) *ServerInstance {
	ins := a.lb.Next(key, a.Instances)
	return ins
}

//服务实例的状态
type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
)

//服务实例
type ServerInstance struct {
	InstanceId string
	AppName    string
	Address    string
	Status     Status
	Metadata   map[string]string
}
