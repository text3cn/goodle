package goodlog

const (
	// 前景颜色
	black  = "\033[30m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	pink   = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"

	// 背景颜色
	blackBG  = "\033[40m"
	redBG    = "\033[41m"
	greenBG  = "\033[42m"
	yellowBG = "\033[43m"
	blueBG   = "\033[44m"
	pinkBG   = "\033[45m"
	cyanBG   = "\033[46m"
	grayBG   = "\033[47m"

	// 文本样式
	textBlob           = "\033[1m" // 粗体
	textDownplay       = "\033[2m" // 淡化
	textItalic         = "\033[3m" // 斜体
	textUnderline      = "\033[4m" // 下划线
	textBlink          = "\033[5m" // 闪烁
	textReverseDisplay = "\033[7m" // 文本背景和前景颜色对调
	textInvisible      = "\033[7m" // 隐藏文本（文本不可见，但仍占用空间）

	// 重置所有颜色和样式: fmt.Printf(reset)
	reset = "\033[0m"

	// 组合效果
	// 这些转义码可以通过组合使用来实现不同的文本效果，
	// 例如，要创建红色粗体文本，可以使用 \033[31;1m，然后使用 \033[0m 重置文本样式。
	// fmt.Sprintf(pink + greenBG + "hello world" + reset)

	// 高亮加粗白
	whiteHighlight = "\033[1;97m"
)

// 红色
func Red(output interface{}) {
	goodlogSvc.Color(red, output)
}

func Redf(output ...interface{}) {
	goodlogSvc.Colorf(red, output...)
}

// 紫色
func Pink(output interface{}) {
	goodlogSvc.Color(pink, output)
}

func Pinkf(output ...interface{}) {
	goodlogSvc.Colorf(pink, output...)
}

// 绿色
func Green(output interface{}) {
	goodlogSvc.Color(green, output)
}

func Greenf(output ...interface{}) {
	goodlogSvc.Colorf(green, output...)
}

// 青色
func Cyan(output interface{}) {
	goodlogSvc.Color(cyan, output)
}

func Cyanf(output ...interface{}) {
	goodlogSvc.Colorf(cyan, output...)
}

// 蓝色
func Blue(output interface{}) {
	goodlogSvc.Color(blue, output)
}

func Bluef(output ...interface{}) {
	goodlogSvc.Colorf(blue, output...)
}

// 黄色
func Yellow(output interface{}) {
	goodlogSvc.Color(yellow, output)
}

func Yellowf(output ...interface{}) {
	goodlogSvc.Colorf(yellow, output...)
}
