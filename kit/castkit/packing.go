package castkit

import "github.com/spf13/cast"

type GoodleVal struct {
	Input interface{}
}

func (this *GoodleVal) ToString() string {
	return cast.ToString(this.Input)
}

func (this *GoodleVal) ToInt() int {
	return cast.ToInt(this.Input)
}

func (this *GoodleVal) ToBool() bool {
	return cast.ToBool(this.Input)
}
