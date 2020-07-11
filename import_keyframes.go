package keyframes

/*本文件主要是针对用户导入的关键帧数据格式做定义
*create:lihaiping1603@aliyun.com
*date:2019-11-12
 */
import (
	"encoding/json"
	"github.com/lihp1603/utils"
	"io/ioutil"
	"os"
)

type Position struct {
	Index int     `json:"index"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

func (this *Position) GetIndex() int {
	return this.Index
}

func (this *Position) GetX() float64 {
	return this.X
}

func (this *Position) GetY() float64 {
	return this.Y
}

type Scale struct {
	Index int     `json:"index"`
	Ratio float64 `json:"ratio"` //前端带过来的值中，会x100表示
}

func (this *Scale) GetIndex() int {
	return this.Index
}

func (this *Scale) GetRatio() float64 {
	return this.Ratio
}

type KeyFrames struct {
	KeyPositions  []Position `json:"position"` //代表原视频中心点相对于原画布的位置，如果前端这个位置参数先乘以了scale的话，我们在后面处理的时候需要/scale的值
	KeyScales     []Scale    `json:"scaling"`  //代表对原视频进行的缩放操作值
	CropPositions []Position //裁剪框中心点相对于原视频的左上角的位置坐标
}

func (this *KeyFrames) GetKeyPostions() []Position {
	return this.KeyPositions
}

func (this *KeyFrames) GetKeyScales() []Scale {
	return this.KeyScales
}

func (this *KeyFrames) GetCropPostions() []Position {
	return this.CropPositions
}

func (this *KeyFrames) AdjustScales() {
	for i, scaleV := range this.KeyScales {
		ratio := scaleV.Ratio / 100.0
		this.KeyScales[i].Ratio = ratio
	}
}

//以原视频的宽高(OrigWidth，OrigHeight)大小为画布，画布固定，裁剪框固定在画布的中心区域，大小位置都固定，
//而瞄点的位置信息为将原视频进行移动以后，原视频的中心点相对于画布左上角的坐标关系
//同时宿放值为将原视频宿放的比例大小
type ImportKeyframes struct {
	OrigHeight    int       `json:"orig_height"`
	OrigWidth     int       `json:"orig_width"`
	OutputHeight  int       `json:"output_height"`
	OutputWidth   int       `json:"output_width"`
	CuttingHeight int       `json:"cutting_height"`
	CuttingWidth  int       `json:"cutting_width"`
	FrameCount    int       `json:"frame_count"`
	KeyFrame      KeyFrames `json:"motions"`
}

func NewImportKeyframes(inFile string) *ImportKeyframes {
	var keyframesData ImportKeyframes
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

	if err = json.Unmarshal(data, &keyframesData); err != nil {
		utils.LogTraceI("in file %s json fmt decoding failed,err is ", inFile, err)
		return nil
	}
	//调整一下数据值
	keyframesData.KeyFrame.AdjustScales()

	//坐一次数据变换
	// keyframesData.TransformPostion()

	return &keyframesData
}

func (this *ImportKeyframes) GetOrigWidth() int {
	return this.OrigWidth
}

func (this *ImportKeyframes) GetOrigHeight() int {
	return this.OrigHeight
}

func (this *ImportKeyframes) GetOutputWidth() int {
	return this.OutputWidth
}

func (this *ImportKeyframes) GetOutputHeight() int {
	return this.OutputHeight
}

func (this *ImportKeyframes) GetCuttingWidth() int {
	return this.CuttingWidth
}

func (this *ImportKeyframes) GetCuttingHeight() int {
	return this.CuttingHeight
}

func (this *ImportKeyframes) GetFramesCount() int {
	return this.FrameCount
}

func (this *ImportKeyframes) GetKeyFrames() KeyFrames {
	return this.KeyFrame
}

//将原视频中心点的坐标变换为裁剪框中心点相对原视频左上角原点的坐标
//即将前端的操作变换为我们的操作
func (this *ImportKeyframes) TransformPostion() {
	//先记录裁剪框中心坐标相对于原视频的最初原始位置，
	crop_x1 := float64(this.GetOrigWidth() / 2)
	crop_y1 := float64(this.GetOrigHeight() / 2)
	for _, pos := range this.KeyFrame.GetKeyPostions() {
		//计算运动矢量
		motion_x := crop_x1 - pos.X
		motion_y := crop_y1 - pos.Y
		crop_x := motion_x + crop_x1
		crop_y := motion_y + crop_y1

		//记录转换后的坐标
		cropPos := Position{Index: pos.Index, X: crop_x, Y: crop_y}
		this.KeyFrame.CropPositions = append(this.KeyFrame.CropPositions, cropPos)
	}
}
