package keyframes

//本文件主要是实现针对关键帧 key scale的处理
import (
	"github.com/psdshow/utils"
)

type KeyScale struct {
	interpolation Interpolationer
}

func NewKeyScale(interplolationer Interpolationer) *KeyScale {
	keyScale := &KeyScale{interpolation: interplolationer}
	return keyScale
}

func (this *KeyScale) GetFrameScales(inKeyFrames ImportKeyframes) []Scale {
	allframeCount := inKeyFrames.GetFramesCount()
	//获取关键帧信息
	keyFrames := inKeyFrames.GetKeyFrames()
	//获取关键帧的scale信息
	keyScales := keyFrames.GetKeyScales()
	//获取关键帧插值
	scaleFrames := this.keyScales2Scales(keyScales)
	if len(scaleFrames) == allframeCount {
		return scaleFrames
	} else if len(scaleFrames) <= 0 {
		return this.getDefalutScales(1.0, 0, allframeCount)
	}
	//获取
	firstKeyframeIndex := scaleFrames[0].GetIndex()
	firstScaleV := scaleFrames[0].GetRatio()
	//补充数据关键帧数据之前的帧数据
	resultFrames := this.getDefalutScales(firstScaleV, 0, firstKeyframeIndex)

	//插入中间关键帧的数据
	resultFrames = append(resultFrames, scaleFrames...)
	//关键帧后面的数据
	lastKeyframeIndex := scaleFrames[len(scaleFrames)-1].GetIndex()
	lastScaleV := scaleFrames[len(scaleFrames)-1].GetRatio()
	lastFrameScales := this.getDefalutScales(lastScaleV, lastKeyframeIndex+1, allframeCount)
	resultFrames = append(resultFrames, lastFrameScales...)
	utils.LogTraceI("frame scale size:%d,frame cout:%d", len(resultFrames), allframeCount)
	return resultFrames
}

//获取默认的值，如果在没有关键帧的情况下
func (this *KeyScale) getDefalutScales(scaleV float64, startIndex, endIndex int) []Scale {
	var resultFrames []Scale
	for i := startIndex; i < endIndex; i++ {
		exportFrameScale := Scale{Index: i, Ratio: scaleV}
		resultFrames = append(resultFrames, exportFrameScale)
	}
	return resultFrames
}

//对关键帧的宿放信息进行转换
func (this *KeyScale) keyScales2Scales(keyScales []Scale) []Scale {
	var frameScales []Scale
	if len(keyScales) <= 0 {
		utils.LogTraceI("the keyframe scale is empty,no scale data")
		return frameScales
	}

	for index, scaleValue := range keyScales {
		if index == 0 {
			continue
		}
		prev := keyScales[index-1]
		frameIndex := prev.Index
		numberframes := scaleValue.Index - prev.Index
		//计算这两个关键帧直接的线性插值
		scale_values := this.interpolation.Interpolation(numberframes, prev.Ratio, scaleValue.Ratio)

		for i := 0; i < numberframes; i++ {
			x := scale_values[i]
			exportFrameScale := Scale{Index: frameIndex, Ratio: x}
			frameScales = append(frameScales, exportFrameScale)
			frameIndex++
		}
	}
	//将最后一帧的关键帧记录下来
	lastKeyScaleValue := keyScales[len(keyScales)-1]
	exportFrameScale := Scale{Index: lastKeyScaleValue.Index, Ratio: lastKeyScaleValue.Ratio}
	frameScales = append(frameScales, exportFrameScale)
	return frameScales
}
