package goodle

import (
	"fmt"
	"github.com/text3cn/goodle/kit/strkit"
	"os"
)

var (
	fontColorBlack       = 30
	fontColorRed         = 31
	fontColorGreen       = 32
	fontColorYellow      = 33
	fontColorBlue        = 34
	fontColorAmaranth    = 35 // 紫红色
	fontColorUltramarine = 36 // 青蓝色
	fontColorWhite       = 37
	bgColorBlack         = 40
	bgColorRed           = 41
	bgColorGreen         = 42
	bgColorYellow        = 43
	bgColorBlue          = 44
	bgColorAmaranth      = 45 // 紫红色
	bgColorUltramarine   = 46 // 青蓝色
	bgColorWhitel        = 47
	bgColorTransparency  = 97 // 透明背景，无背景
)

func drawControl() {
	fmt.Fprintf(os.Stdout, "\033[K")
	borderBackground := 97 // 无背景色
	borderFontColor := 33
	statusSuccessFontColor := 32
	statusFailFontColor := 31
	style := 0
	text := "--------------------- Goodle Framework ------------------------"
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, style, borderBackground, borderFontColor, text, 0x1B)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, borderBackground, borderFontColor, "|", 0x1B)
	fmt.Fprintf(os.Stdout, "\033[?25h")
	// 从行首到监听一共45个字符/空格
	// 从行首到状态一共96个字符/空格
	// 不能用tab制表符号来排列布局，统一使用空格计算
	text = " http" + "           " + "listen" + "               " + "status" + "  cpu" + "   memory    "
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, borderBackground, 36, text, 0x1B)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, style, borderBackground, borderFontColor, "|", 0x1B)
	text = "---------------------------------------------------------------"
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, style, borderBackground, borderFontColor, text, 0x1B)

	// http
	listen := "8000"
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, borderBackground, borderFontColor, "|", 0x1B)
	fmt.Printf(" Http" + strkit.StrRepeat(" ", 15-4) + listen + strkit.StrRepeat(" ", 86-45-len(listen)))
	if true {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, 97, statusSuccessFontColor, "ok", 0x1B)
	} else {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, 97, statusFailFontColor, "no", 0x1B)
	}
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, style, borderBackground, borderFontColor, "  |", 0x1B)

	// discovery
	listen = "8888"
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, borderBackground, borderFontColor, "|", 0x1B)
	fmt.Printf(" Consul" + strkit.StrRepeat(" ", 15-6) + listen + strkit.StrRepeat(" ", 86-45-len(listen)))
	if true {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, 97, statusSuccessFontColor, "ok", 0x1B)
	} else {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, style, 97, statusFailFontColor, "no", 0x1B)
	}
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, style, borderBackground, borderFontColor, "  |", 0x1B)

	text = "---------------------------------------------------------------"
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n\n", 0x1B, style, borderBackground, borderFontColor, text, 0x1B)

}
