package keyframes

/*这里主要是针对输入数据和输出数据做变换，把用户数据转成我们需要的数据格式
*create:lihaiping1603@aliyun.com
*date:2019-11-12
 */

import (
	"github.com/psdshow/utils"
)

type Transform struct {
}

func NewTransform() *Transform {
	transform := Transform{}
	return &transform
}

func (this *Transform) Import2Export(inKeyFrames ImportKeyframes) *ExportFrames {
	return this.keyFramesProcess(inKeyFrames, NewLinearInterpolation())
}

func (this *Transform) Import2ExportEx(inKeyFrames ImportKeyframes, interpolationer Interpolationer) *ExportFrames {
	return this.keyFramesProcess(inKeyFrames, interpolationer)
}

func (this *Transform) keyFramesProcess(inKeyFrames ImportKeyframes, interpolationer Interpolationer) *ExportFrames {
	var frames []ExportFrameInfo
	//获取源视频宽高和目的视频宽高
	orgW := inKeyFrames.GetOrigWidth()
	orgH := inKeyFrames.GetOrigHeight()
	//获取最后输出的视频宽高
	outW := inKeyFrames.GetOutputWidth()
	outH := inKeyFrames.GetOutputHeight()
	//获取源视频中心点坐标
	orgCentreX := orgW / 2
	orgCentreY := orgH / 2
	//获取最新的背景画布大小
	bgW, bgH := this.getNewBgWH(inKeyFrames)
	utils.LogTraceI("the new backgroud width x height:%dx%d", bgW, bgH)
	//获取新画布的视频中心点
	bgCentreX := bgW / 2
	bgCentreY := bgH / 2
	//设置一个读取最后一个帧数据的值，用于处理关键帧的数据不足原视频的帧数的情况
	lastScaleV := 1.0
	//先获取scale插值
	key_scale := NewKeyScale(interpolationer)
	frameScales := key_scale.GetFrameScales(inKeyFrames)
	//再获取位置position插值
	key_position := NewKeyPostion(interpolationer)
	framePositions := key_position.GetFramePositions(inKeyFrames)
	//处理有用的信息转换
	for index, framePosV := range framePositions {
		scaleV := lastScaleV
		if len(frameScales) > index { //如果有的话
			scaleV = frameScales[index].Ratio
		}

		//+0.5是做四舍五入取整用
		cropW := int(float64(outW)/scaleV + 0.5)
		cropH := int(float64(outH)/scaleV + 0.5)
		//偶数处理
		if cropW%2 != 0 {
			cropW -= 1
		}
		if cropH%2 != 0 {
			cropH -= 1
		}
		//将framePosV.X，Y代表视频宿放以后得到的视频中心点相对于原视频画布的位置，所以这里先/scaleV,将视频进行还原，然后再/minScale，换算为新画布的坐标
		// x := float64(framePosV.X) / (scaleV * minScale)
		// y := float64(framePosV.Y) / (scaleV * minScale)
		//将framePosV.X，Y代表原视频移动以后得到的视频中心点相对于原视频画布的位置，所以直接/minScale，换算为该点位置在新画布的坐标
		// x := float64(framePosV.X) / minScale
		// y := float64(framePosV.Y) / minScale

		//然后计算原视频位移的矢量为多少
		motion_x := orgCentreX - int(framePosV.X) //得到原视频实际平移的距离
		motion_y := orgCentreY - int(framePosV.Y)
		//因为位移矢量不参与宿放值，因为我们原始视频是没有进行宿放的，只是我们的画布变大了而已
		// x := bgCentreX + motion_x
		// y := bgCentreY + motion_y
		//如果运动矢量也进行宿放的话
		scale_motion_x := int(float64(motion_x)/scaleV + 0.5)
		scale_motion_y := int(float64(motion_y)/scaleV + 0.5)
		x := bgCentreX + scale_motion_x
		y := bgCentreY + scale_motion_y
		// utils.LogTraceI("index:%d,min scale:%f,frame pos x:%f,y:%f,motion_x:%d,motion_y:%d", index, minScale, framePosV.X, framePosV.Y, motion_x, motion_y)
		// utils.LogTraceI("index:%d,scaleV:%f,cropW x cropH:%dx%d,x:%d,y:%d", index, scaleV, cropW, cropH, x, y)
		//计算
		// exportFrame := ExportFrameInfo{FrameIndex: framePosV.Index, Box: []int{cropW, cropH, int(x - float64(cropW)/2 + 0.5), int(y - float64(cropH)/2 + 0.5)}}
		exportFrame := ExportFrameInfo{FrameIndex: framePosV.Index, Box: []int{cropW, cropH, int(x - cropW/2), int(y - cropH/2)}}
		frames = append(frames, exportFrame)
		//更新最后一帧的数值
		lastScaleV = scaleV
	}

	return NewExportFrames(bgW, bgH, outW, outH, frames)
}

//获取新的画布大小
func (this *Transform) getNewBgWH(inKeyFrames ImportKeyframes) (bgW, bgH int) {
	//获取源视频宽高和目的视频宽高
	orgW := inKeyFrames.GetOrigWidth()
	orgH := inKeyFrames.GetOrigHeight()
	//获取最后输出的视频宽高
	outW := inKeyFrames.GetOutputWidth()
	outH := inKeyFrames.GetOutputHeight()
	//先获取最小的宿放比例大小
	//获取关键帧信息
	keyFrames := inKeyFrames.GetKeyFrames()
	//获取关键帧的位置坐标设置信息
	keyScales := keyFrames.GetKeyScales()
	minScale := 1.0
	for _, scaleV := range keyScales {
		if scaleV.Ratio < minScale {
			minScale = scaleV.Ratio
		}
	}
	//这个时候，说明输出的裁剪框如果同样进行同比例缩放的话，裁剪框超过了视频新画布的大小，会导致裁剪失败，所以需要将背景画布再放大一次
	out_org_ratio := float64(orgW) / float64(outW)
	out_org_ratio_h := float64(orgH) / float64(outH)
	if out_org_ratio_h < out_org_ratio {
		out_org_ratio = out_org_ratio_h
	}

	if minScale > out_org_ratio {
		minScale = out_org_ratio
	}
	//求出新的画布的大小为多少
	//这里实际上是计算一个我们所需的最小画布大小，
	//如果需要跟ae类似的话，这里我们可以建立一个虚拟的很大的背景画布，例如8000x8000，这样就可以满足大部分的关键帧动画了，实现方式不变，还是同样的方法
	//这里求最大背景大小
	bgW = int(float64(orgW) / minScale)
	bgH = int(float64(orgH) / minScale)

	if bgW%2 != 0 {
		bgW += 1
	}
	if bgH%2 != 0 {
		bgH += 1
	}
	//写死一个大的背景也是可以的，试一下:4Kx4K
	bgW = 4000
	bgH = 4000

	return
}
