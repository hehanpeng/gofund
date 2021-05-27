package invoker

import (
	"fmt"
	"github.com/hehanpeng/gofund/proto/fund/gen/errcodepb"

	"github.com/gotomicro/ego-component/egorm"
	"github.com/gotomicro/ego-component/eredis"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Logger    *elog.Component
	Db        *egorm.Component
	RedisStub *eredis.Component
)

func Init() error {
	Logger = elog.DefaultLogger
	Db = egorm.Load("mysql.user").Build()
	RedisStub = eredis.Load("redis.user").Build(eredis.WithStub())
	return nil
}

// 记录grpc error信息
func Error(code errcodepb.ErrCode.ErrCode, err error) error {
	Logger.Error("grpc error: ", zap.Int32("code", int32(code)), zap.Error(err))
	var cause string
	if err != nil {
		cause = err.Error()
	}
	return status.Error(codes.Code(code), fmt.Sprintf("error name: %s, cause: %s", errcodepb.ErrCode_name[int32(code)], cause))
}
