package core

var frameContainer *ServicesContainer

// F全局服务中心，非 web 服务开发用
func FrameContainer() *ServicesContainer {
	if frameContainer == nil {
		frameContainer = New()
	}
	return frameContainer
}
