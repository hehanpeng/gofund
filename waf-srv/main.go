package main

import (
	invoker2 "waf-srv/pkg/invoker"
	job2 "waf-srv/pkg/job"
	router2 "waf-srv/pkg/router"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egovernor"
)

//  export EGO_DEBUG=true && go run main.go --config=config/dev.toml --job=install
func main() {
	if err := ego.New().
		Invoker(invoker2.Init).
		Job(
			job2.InstallComponent(),
		).
		Serve(
			egovernor.Load("server.governor").Build(),
			router2.ServeGRPC(),
		).
		Run(); err != nil {
		elog.Panic(err.Error())
	}
}
