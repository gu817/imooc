package httpclient

import (
	"git.imooc.com/wendell1000/infra/lb"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpClient_Do(t *testing.T) {
	//创建一个eureka client
	conf := ini.NewIniFileConfigSource("ec_test.ini")
	client := eureka.NewClient(conf)
	client.Start()
	client.Applications, _ = client.GetApplications()

	//创建一个apps实例
	apps := &lb.Apps{Client: client}

	c := NewHttpClient(apps, &Option{
		Timeout: defaultHttpTimeout,
	})
	Convey("http客户端", t, func() {
		for i := 0; i < 10; i++ {

			req, err := c.NewRequest(http.MethodGet,
				"http://resk/",
				nil, nil)
			So(err, ShouldBeNil)
			So(req, ShouldNotBeNil)
			res, err := c.Do(req)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)

			defer res.Body.Close()
			d, err := ioutil.ReadAll(res.Body)
			So(err, ShouldBeNil)
			So(d, ShouldNotBeNil)

		}
	})
}

//作为作业留给同学们：
// 使用简单轮询算法，运行测试用例试一下，看看效果，并且能发现什么问题呢？
//想一想问题是什么？
// 造成这个问题的原因是什么？
// 那么如何来解决这个问题？
//那么在问答区或者QQ群里我们展开讨论
