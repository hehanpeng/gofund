module ups-srv

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/gotomicro/ego v0.5.6
	github.com/gotomicro/ego-component/egorm v0.2.0
	github.com/gotomicro/ego-component/eredis v0.2.2
	github.com/hehanpeng/gofund/common v0.0.0-00010101000000-000000000000
	github.com/hehanpeng/gofund/proto v0.0.0-20210531035004-61dc2bb696b1
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.29.1
	waf-srv v0.0.0-00010101000000-000000000000
)

replace (
	github.com/hehanpeng/gofund/common => ../common
	github.com/hehanpeng/gofund/proto => ../proto
	waf-srv => ../waf-srv
)
