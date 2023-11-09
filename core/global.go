package core

var globalCore *ServicesContainer

// 全局服务中心
func GlobalCore() *ServicesContainer {
	if globalCore == nil {
		globalCore = New()
	}
	return globalCore
}
