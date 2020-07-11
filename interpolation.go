package keyframes

//本文件主要是实现各种插值算法
type Interpolationer interface {
	Interpolation(numberFrames int, p1, p2 float64) []float64
}

//线性插值
type LinearInterpolation struct {
}

func NewLinearInterpolation() *LinearInterpolation {
	return &LinearInterpolation{}
}

//线性插值运算，返回的值中会把p1带上返回
func (this *LinearInterpolation) Interpolation(numberFrames int, p1, p2 float64) []float64 {
	var values []float64
	if numberFrames <= 0 {
		return values
	}
	diff_values := p2 - p1
	cur_value := p1
	increment_value := diff_values / float64(numberFrames)
	for i := 0; i < numberFrames; i++ {
		values = append(values, cur_value)
		cur_value += increment_value
	}
	return values
}

//比塞尔插值
type BezierInterpolation struct {
}

func (this *BezierInterpolation) Interpolation(numberFrames int, p1, p2 float64) []float64 {
	var values []float64
	return values
}
