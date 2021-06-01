package invoker

import (
	"errors"
	"fmt"
	"github.com/gotomicro/ego-component/egorm"
	"github.com/gotomicro/ego-component/eredis"
	"github.com/gotomicro/ego-component/eredis/ecronlock"
	"github.com/gotomicro/ego/client/ehttp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"github.com/hehanpeng/gofund/proto/fund/gen/errcodepb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Logger      *elog.Component
	Gin         *egin.Component
	Db          *egorm.Component
	RedisStub   *eredis.Component
	UpsHttpComp *ehttp.Component
	EcronLocker *ecronlock.Component
	//UpsSrvGrpc   upssrv.UpsClient
	//EtcdClient   *eetcd.Component
	//EtcdRegistry *registry.Component
	CallSrvHttpComp = make(map[string]*ehttp.Component)
)

func Init() error {
	Logger = elog.DefaultLogger
	Gin = egin.Load("server.http").Build()
	Db = egorm.Load("mysql.waf").Build()
	RedisStub = eredis.Load("redis.waf").Build(eredis.WithStub())
	EcronLocker = ecronlock.DefaultContainer().Build(ecronlock.WithClient(RedisStub))
	UpsHttpComp = ehttp.Load("http.ups").Build()
	RegisterCallSrvHttpComp("wafs", UpsHttpComp)
	//EtcdClient = eetcd.Load("etcd").Build()
	//EtcdRegistry = registry.Load("registry").Build(registry.WithClientEtcd(EtcdClient))

	// 必须注册在grpc前面
	//resolver.Register("etcd", EtcdRegistry)
	//userConn := egrpc.Load("grpc.upssrv").Build().ClientConn
	//UpsSrvGrpc = upssrv.NewUpsClient(userConn)
	return nil
}

// 记录grpc error信息
func Error(code errcodepb.ErrCode, err error) error {
	Logger.Error("grpc error: ", zap.Int32("code", int32(code)), zap.Error(err))
	var cause string
	if err != nil {
		cause = err.Error()
	}
	return status.Error(codes.Code(code), fmt.Sprintf("error name: %s, cause: %s", errcodepb.ErrCode_name[int32(code)], cause))
}

func RegisterCallSrvHttpComp(key string, comp *ehttp.Component) (err error) {
	if CallSrvHttpComp[key] != nil {
		Logger.Errorf("%v 已注册comp,无法重复注册", key)
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CallSrvHttpComp[key] = comp
		return nil
	}
}
