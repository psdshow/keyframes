package keyframes

//本文件主要是针对key postion位置关系实现
import (
	"github.com/psdshow/utils"
)

type KeyPostion struct {
	interpolation Interpolationer
}

func NewKeyPostion(interplolationer Interpolationer) *KeyPostion {
	keyPostion := &KeyPostion{interpolation: interplolationer}
	return keyPostion
}

//补充位置信息
func (this *KeyPostion) GetFramePositions(inKeyFrames ImportKeyframes) []Position {
	allframeCount := inKeyFrames.GetFramesCount()
	inKeyFrame := inKeyFrames.GetKeyFrames()
	orgCentreX := inKeyFrames.GetOrigWidth() / 2
	orgCentreY := inKeyFrames.GetOrigHeight() / 2
	keyPostions := inKeyFrame.GetKeyPostions()
	//获取到关键帧插值以后的帧位置信息
	framePositions := this.keyPostions2Postions(keyPostions)
	if len(framePositions) == allframeCount {
		return framePositions
	} else if len(framePositions) <= 0 {
		return this.getDefaultPosition(orgCentreX, orgCentreY, 0, allframeCount)
	}
	//否则需要补充数据
	//获取
	firstKeyframeIndex := framePositions[0].GetIndex()
	firstKeyframeX := framePositions[0].GetX()
	firstKeyframeY := framePositions[0].GetY()
	//补充数据
	resultFrames := this.getDefaultPosition(int(firstKeyframeX), int(firstKeyframeY), 0, firstKeyframeIndex)
	// utils.LogTraceI("first framePostions size:%d,first cout:%d,addr:%p", len(resultFrames), firstKeyframeIndex, resultFrames)

	resultFrames = append(resultFrames, framePositions...)
	// utils.LogTraceI("key framePostions size:%d,key frame cout:%d,addr:%p,%p", len(resultFrames), len(framePositions), resultFrames, framePositions)

	lastKeyframeIndex := framePositions[len(framePositions)-1].GetIndex()
	lastKeyframeX := framePositions[len(framePositions)-1].GetX()
	lasKeyframeY := framePositions[len(framePositions)-1].GetY()
	lastFramePostions := this.getDefaultPosition(int(lastKeyframeX), int(lasKeyframeY), lastKeyframeIndex+1, allframeCount)
	// utils.LogTraceI("last framePostions size:%d,last count:%d,lastFrameIndex:%d,count:%d,addr:%p,%p", len(resultFrames), len(lastFramePostions), lastKeyframeIndex+1, allframeCount-lastKeyframeIndex, resultFrames, lastFramePostions)
	resultFrames = append(resultFrames, lastFramePostions...)

	utils.LogTraceI("framePostions size:%d,frame cout:%d", len(resultFrames), allframeCount)
	return resultFrames
}

//定义一个获取默认的参数
func (this *KeyPostion) getDefaultPosition(defaultX, defaultY int, startIndex, endIndex int) []Position {
	var resultFrames []Position
	for i := startIndex; i < endIndex; i++ {
		exportFramePosition := Position{Index: i, X: float64(defaultX), Y: float64(defaultY)}
		resultFrames = append(resultFrames, exportFramePosition)
		// utils.LogTraceI("resultFrames len:%d", len(resultFrames))
	}
	// utils.LogTraceI("resultFrames len:%d,%p", len(resultFrames), resultFrames)
	return resultFrames
}

//对关键帧位置信息进行转换
func (this *KeyPostion) keyPostions2Postions(keyPostions []Position) []Position {
	var framePositions []Position
	if len(keyPostions) <= 0 {
		utils.LogTraceI("the keyframe position is empty,no key frame position data")
		return framePositions
	}
	for index, posValue := range keyPostions {
		if index == 0 {
			continue
		}
		prev := keyPostions[index-1]
		frameIndex := prev.Index
		numberframes := posValue.Index - prev.Index
		//计算这两个关键帧直接的线性插值
		x_values := this.interpolation.Interpolation(numberframes, prev.X, posValue.X)
		y_values := this.interpolation.Interpolation(numberframes, prev.Y, posValue.Y)

		for i := 0; i < numberframes; i++ {
			x := x_values[i]
			y := y_values[i]
			if int(x)%2 != 0 {
				x -= 1
			}
			if int(y)%2 != 0 {
				y -= 1
			}
			exportFramePosition := Position{Index: frameIndex, X: x, Y: y}
			framePositions = append(framePositions, exportFramePosition)
			frameIndex++
		}
	}
	//将最后一帧的关键帧记录下来
	lastKeyPosValue := keyPostions[len(keyPostions)-1]
	exportFramePosition := Position{Index: lastKeyPosValue.Index, X: lastKeyPosValue.X, Y: lastKeyPosValue.Y}
	framePositions = append(framePositions, exportFramePosition)
	return framePositions
}
