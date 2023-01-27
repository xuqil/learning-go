package sync

import "sync"

type MiBiz struct {
	once sync.Once
}

// Init 接收器使用指针或者 one 成员使用指针(Once 是不要被复制使用的)
func (m *MiBiz) Init() {
	m.once.Do(func() {

	})
}

type MyBiz1 struct {
	once *sync.Once
}

// Init 接收器使用指针或者 one 成员使用指针(Once 是不要被复制使用的)
func (m MyBiz1) Init() {
	m.once.Do(func() {

	})
}

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

var s *singleton
var singletonOnce sync.Once

// GetSingleton 属于懒加载
func GetSingleton() MyBusiness {
	singletonOnce.Do(func() {
		s = &singleton{}
	})
	return s
}

//// GetSingleton 直接返回 *singleton 会有告警，可以返回接口
//func GetSingleton() *singleton {
//	singletonOnce.Do(func() {
//		s = &singleton{}
//	})
//	return s
//}

// init 属于饥饿模式
func init() {
	// 用包初始化函数取代 One
	s = &singleton{}
}
