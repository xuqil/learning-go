package one

import "sync"

type MyBusiness interface {
	DoSomething()
}

// singleton 单例模式（一般与接口一起使用）
type singleton struct {
}

func (s singleton) DoSomething() {
	//TODO implement me
	panic("implement me")
}

var instance *singleton
var singletonOnce sync.Once

// GetSingleton 属于懒加载
func GetSingleton() MyBusiness {
	singletonOnce.Do(func() {
		instance = &singleton{}
	})
	return instance
}

// init 属于饥饿模式
func init() {
	// 用包初始化函数取代 One
	instance = &singleton{}
}
