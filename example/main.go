package main

import (
	"flag"
	"github.com/lihp1603/keyframes"
	"github.com/lihp1603/utils"
)

func init() {

	//初始化解析用户参数
	flag.Parse()
}
func main() {
	//从关键帧的数据导入
	inKeyFrame := keyframes.NewImportKeyframes("import_keyframes_example2.json")
	if inKeyFrame == nil {
		utils.LogTraceE("import keyframes data failed")
		return
	}
	//生成一个转换的对象
	transform := keyframes.NewTransform()
	//生成导出的对象
	outFrame := transform.Import2Export(*inKeyFrame)
	//将数据导出成文件
	outFrame.ExportFile("out_frames_example2.json")
	//
	utils.LogTraceI("export:bgWxbgH:%dx%d,outWxoutH:%dx%d", outFrame.BgW, outFrame.BgH, outFrame.OutW, outFrame.OutH)
	return

}
