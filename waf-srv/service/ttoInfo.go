package service

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"
	"waf-srv/model"
	request2 "waf-srv/model/request"
	"waf-srv/pkg/invoker"
	"waf-srv/statistics"

	"github.com/hehanpeng/gofund/common/helper"

	"go.uber.org/zap"

	"github.com/hehanpeng/gofund/common/global/api"

	"github.com/gotomicro/ego/core/etrace"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateTtoInfo
//@description: 创建TtoInfo记录
//@param: ttoInfo model.TtoInfo
//@return: err error

func CreateTtoInfo(ttoInfo model.TtoInfo) (err error) {
	err = invoker.Db.Create(&ttoInfo).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteTtoInfo
//@description: 删除TtoInfo记录
//@param: ttoInfo model.TtoInfo
//@return: err error

func DeleteTtoInfo(ttoInfo model.TtoInfo) (err error) {
	err = invoker.Db.Delete(&ttoInfo).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteTtoInfoByIds
//@description: 批量删除TtoInfo记录
//@param: ids req.IdsReq
//@return: err error

func DeleteTtoInfoByIds(ids api.IdsReq) (err error) {
	err = invoker.Db.Delete(&[]model.TtoInfo{}, "id in ?", ids.Ids).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateTtoInfo
//@description: 更新TtoInfo记录
//@param: ttoInfo *model.TtoInfo
//@return: err error

func UpdateTtoInfo(ttoInfo model.TtoInfo) (err error) {
	err = invoker.Db.Save(&ttoInfo).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetTtoInfo
//@description: 根据id获取TtoInfo记录
//@param: id uint
//@return: err error, ttoInfo model.TtoInfo

func GetTtoInfo(id uint) (err error, ttoInfo model.TtoInfo) {
	err = invoker.Db.Where("id = ?", id).First(&ttoInfo).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetTtoInfoInfoList
//@description: 分页获取TtoInfo记录
//@param: info req.TtoInfoSearch
//@return: err error, list interface{}, total int64

func GetTtoInfoInfoList(info request2.TtoInfoSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := invoker.Db.Model(&model.TtoInfo{})
	var ttoInfos []model.TtoInfo
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&ttoInfos).Error
	return err, ttoInfos, total
}

// 执行超时转发逻辑
func DealTto(chanID uint64, ttoInfo model.TtoInfo, ch chan<- *statistics.RequestResults, wg *sync.WaitGroup) error {
	ctx := context.Background()
	startTime := time.Now()
	result := new(api.R)
	key := "tto_id" + strconv.Itoa(int(ttoInfo.ID))
	// 返回chan数据
	defer func() {
		requestTime := uint64(helper.DiffNano(startTime))
		requestResults := &statistics.RequestResults{
			Time:      requestTime,
			IsSucceed: result.IsSuccess(),
			ErrCode:   result.Code,
		}
		requestResults.SetID(chanID, uint64(ttoInfo.ID))
		ch <- requestResults
		wg.Done()
	}()
	lock, err := invoker.RedisStub.LockClient().Obtain(ctx, key, 3*time.Second)
	if err != nil {
		invoker.Logger.Warn("lock tto_id error", zap.Error(err), zap.String("key", key))
		return errors.New("lock tto_id error")
	}
	defer func() {
		if result.IsSuccess() {
			// 更新成已执行
			err := invoker.Db.Model(&ttoInfo).Update("tto_status", "1").Error
			if err != nil {
				invoker.Logger.Error("update tto error", zap.Error(err))
			}
		}
		lock.Release(ctx)
	}()
	// http 调用
	callSrvHttpComp := invoker.CallSrvHttpComps[ttoInfo.CallSrvName]
	if callSrvHttpComp == nil {
		invoker.Logger.Error("callSrvHttpComp is not register", zap.String("key", ttoInfo.CallSrvName))
		return errors.New("callSrvHttpComp is not register")
	}
	span, ctx := etrace.StartSpanFromContext(ctx, "callHTTP()")
	defer span.Finish()
	req := callSrvHttpComp.R()
	// Inject traceId Into Header
	c1 := etrace.HeaderInjector(ctx, req.Header)
	// 默认2s超时
	info, err := req.SetContext(c1).
		SetBody(ttoInfo).
		SetResult(&api.R{}).
		ExpectContentType("application/json").
		Post(ttoInfo.CallMethod)
	if err != nil {
		invoker.Logger.Error("callSrvHttpComp post error", zap.Error(err))
		return errors.New("callSrvHttpComp post error")
	}
	result = info.Result().(*api.R)
	return nil
}

func RegisterTto(ttoInfo model.TtoInfo) (err error) {
	ttoInfo.RegisterTime = time.Now()
	ttoInfo.ExecuteTime = ttoInfo.RegisterTime.Add(time.Duration(ttoInfo.ExpiredTime) * time.Millisecond)
	ttoInfo.TtoStatus = "0"
	err = invoker.Db.Create(&ttoInfo).Error
	return err
}

func CancelTto(id uint) (err error) {
	update := invoker.Db.Model(&model.TtoInfo{}).Where("id = ?", id).Update("tto_status", "1")
	if update.Error != nil {
		return update.Error
	} else if update.RowsAffected == 0 {
		return errors.New("记录不存在")
	}
	return nil
}
