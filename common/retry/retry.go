package retry

import (
	"time"

	"git.code.oa.com/grpc-go/grpc-filter/slime/retry"
	"git.code.oa.com/grpc-go/grpc-go/errs"
)

// GetNetworkErrRetry 获取网络错误重试
func GetNetworkErrRetry(retryNum int) *retry.Retry {
	if retryNum <= 0 {
		retryNum = 1
	}

	// 网络错误重试, 间隔一毫秒重试一次 https://iwiki.woa.com/pages/viewpage.action?pageId=276029299
	// 101,111,131,141
	r, _ := retry.New(retryNum, []int{errs.RetClientNetErr, errs.RetClientTimeout, errs.RetClientConnectFail,
		errs.RetClientRouteErr},
		retry.WithLinearBackoff(time.Millisecond))
	return r
}
