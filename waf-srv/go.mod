module waf-srv

go 1.16

require (
	github.com/gotomicro/ego v0.5.6
	github.com/gotomicro/ego-component/egorm v0.2.0
	github.com/gotomicro/ego-component/eredis v0.2.2
	github.com/hehanpeng/gofund/proto v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.15.0
	google.golang.org/grpc v1.38.0
)

replace github.com/hehanpeng/gofund/proto => ../proto
