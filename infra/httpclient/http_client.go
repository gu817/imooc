package httpclient

import (
	"bytes"
	"errors"
	"git.imooc.com/wendell1000/infra/lb"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultHttpTimeout = 30 * time.Second
)

var parseUrl = url.Parse

type Option struct {
	Timeout time.Duration
}

type HttpClient struct {
	client *http.Client
	Option Option
	apps   *lb.Apps
}

func NewHttpClient(apps *lb.Apps, opt *Option) *HttpClient {
	c := &HttpClient{
		apps: apps,
	}
	if opt == nil {
		c.Option = Option{
			Timeout: defaultHttpTimeout,
		}
	} else {
		c.Option = *opt
	}
	c.client = &http.Client{
		Timeout: c.Option.Timeout,
	}

	return c
}
func (c *HttpClient) NewRequest(
	method, url string,
	body io.Reader,
	headers http.Header) (*http.Request, error) {
	if method == "" {
		method = http.MethodGet
	}
	//解析URL
	u, err := parseUrl(url)
	if err != nil {
		return nil, err
	}
	//从解析后的URL中提取微服务名称
	name := u.Host
	//通过微服务名称从本地服务注册表中查询应用和应用实例列表
	app := c.apps.Get(name)
	if app == nil {
		return nil, errors.New("没有可用的微服务应用，应用名称：" + name + ",请求：" + url)
	}

	//通过负载均衡算法从应用实例列表中选择一个实例
	ins := app.Get(url)
	if ins == nil {
		return nil, errors.New("没有可用的应用实例，应用名称：" + name + ",请求：" + url)
	}
	//将原来URL中的域名部分替换成选择的实例IP和端口
	u.Host = ins.Address

	//使用新构造URL创建一个Request
	//fmt.Println("修改前", url)
	url = u.String()
	//fmt.Println("修改后", url)

	r, err := http.NewRequest(method, url, body)
	if len(headers) > 0 {
		for key, value := range headers {
			for _, val := range value {
				r.Header.Add(key, val)
			}
		}
	}
	return r, err

}
func (c *HttpClient) Do(r *http.Request) (*http.Response, error) {
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return res, err
}
