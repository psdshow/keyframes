package keyframes

/*这里主要是针对我们最后用ffmpeg来做动画而需要导出的数据格式做定义
*create:lihaiping1603@aliyun.com
*date:2019-11-12
 */

import (
	"encoding/json"
	"github.com/psdshow/utils"
	"io/ioutil"
	"os"
)

type ExportFrameInfo struct {
	FrameIndex int   `json:"frame_index"`
	Box        []int `json:"box"` //w,h,x,y数组内部的位置按照这个来设置
}

func (this *ExportFrameInfo) GetFrameIndex() int {
	return this.FrameIndex
}

func (this *ExportFrameInfo) GetWidth() int {
	if len(this.Box) > 0 {
		return this.Box[0]
	}
	return 0
}

func (this *ExportFrameInfo) GetHeight() int {
	if len(this.Box) > 1 {
		return this.Box[1]
	}
	return 0
}

func (this *ExportFrameInfo) GetX() int {
	if len(this.Box) > 2 {
		return this.Box[2]
	}
	return 0
}

func (this *ExportFrameInfo) GetY() int {
	if len(this.Box) > 3 {
		return this.Box[3]
	}
	return 0
}

type ExportFrames struct {
	BgW    int //画布宽度
	BgH    int //画布高度
	OutW   int //`json:"output_width"`//输出目标宽度
	OutH   int //`json:"output_height"`//输出目标的高度
	Frames []ExportFrameInfo
}

func NewExportFrames(bgW, bgH, outW, outH int, frames []ExportFrameInfo) *ExportFrames {
	exportFrames := ExportFrames{BgW: bgW, BgH: bgH, OutW: outW, OutH: outH, Frames: frames}
	return &exportFrames
}

func NewExportFramesWithInfile(inFile string) *ExportFrames {
	var framesData ExportFrames
	var f *os.File
	var data []byte
	f, err := os.Open(inFile)
	if err != nil {
		utils.LogTraceI("open in file %s failed,err is ", inFile, err)
		return nil
	}

	defer f.Close()
	if data, err = ioutil.ReadAll(f); err != nil {
		utils.LogTraceI("in file %s read filed,err is ", inFile, err)
		return nil
	}

	if err = json.Unmarshal(data, &framesData.Frames); err != nil {
		utils.LogTraceI("in file %s json fmt decoding failed,err is ", inFile, err)
		return nil
	}

	return &framesData
}

func (this *ExportFrames) ExportFile(outFile string) error {
	data, err := json.Marshal(this.Frames)
	if err != nil {
		utils.LogTraceI("json fmt encoding failed,err is ", err)
		return err
	}

	if err = ioutil.WriteFile(outFile, data, 0777); err != nil {
		utils.LogTraceI("out file %s write filed,err is ", outFile, err)
		return err
	}
	return nil
}

func (this *ExportFrames) GetAllFrames() []ExportFrameInfo {
	return this.Frames
}

func (this *ExportFrames) GetFramesCount() int {
	return len(this.Frames)
}

func (this *ExportFrames) GetFrameInfo(index int) ExportFrameInfo {
	if index < this.GetFramesCount() {
		return this.Frames[index]
	}
	return ExportFrameInfo{}
}

func (this *ExportFrames) GetBgW() int {
	return this.BgW
}

func (this *ExportFrames) GetBgH() int {
	return this.BgH
}

func (this *ExportFrames) GetOutW() int {
	return this.OutW
}

func (this *ExportFrames) GetOutH() int {
	return this.OutH
}
