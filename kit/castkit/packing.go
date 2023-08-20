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

func (this *GoodleVal) ToInt32() int32 {
	return cast.ToInt32(this.Input)
}

func (this *GoodleVal) ToInt64() int64 {
	return cast.ToInt64(this.Input)
}

func (this *GoodleVal) ToFloat32() float32 {
	return cast.ToFloat32(this.Input)
}

func (this *GoodleVal) ToFloat64() float64 {
	return cast.ToFloat64(this.Input)
}

func (this *GoodleVal) ToBool() bool {
	return cast.ToBool(this.Input)
}
