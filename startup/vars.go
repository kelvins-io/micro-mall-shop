package startup

import (
	"gitee.com/cristiane/micro-mall-shop/pkg/util/goroutine"
	"gitee.com/cristiane/micro-mall-shop/vars"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	vars.GPool = goroutine.NewPool(20,100)
	return err
}
