package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func ConsoleWidget() []Widget {
	var inputDir, outputDir *walk.LineEdit
	var checkPNG, checkJPEG *walk.CheckBox
	var qualityNum, taskNum *walk.NumberEdit
	var compressLevel *walk.ComboBox

	return []Widget{
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text:    "HEIC File Directory: ",
					MinSize: Size{Width: 100},
				},
				LineEdit{
					AssignTo: &inputDir,
					Text:     ConfigGet().InputDir,
					OnEditingFinished: func() {
						dir := inputDir.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "HEIC file directory is not exist")
								inputDir.SetText(ConfigGet().InputDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "HEIC file directory is not directory")
								inputDir.SetText(ConfigGet().InputDir)
								return
							}
							return
						}
						InputDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().InputDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as HEIC file directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as HEIC file directory", dlgDir.FilePath)
							inputDir.SetText(dlgDir.FilePath)
							InputDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text:    "Output Directory: ",
					MinSize: Size{Width: 100},
				},
				LineEdit{
					AssignTo: &outputDir,
					Text:     ConfigGet().OutputDir,
					OnEditingFinished: func() {
						dir := outputDir.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "Output directory is not exist")
								outputDir.SetText(ConfigGet().OutputDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "Output directory is not directory")
								inputDir.SetText(ConfigGet().OutputDir)
								return
							}
						}
						OutputDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().OutputDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as output directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as output directory", dlgDir.FilePath)
							outputDir.SetText(dlgDir.FilePath)
							OutputDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text:    "Output Options: ",
					MinSize: Size{Width: 100},
				},
				CheckBox{
					AssignTo: &checkPNG,
					Checked:  ConfigGet().PngEnable,
					Text:     "PNG",
					OnCheckedChanged: func() {
						compressLevel.SetEnabled(checkPNG.Checked())
						err := PngEnableSave(checkPNG.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				Label{
					Text: "Compression level: ",
				},
				ComboBox{
					AssignTo:     &compressLevel,
					CurrentIndex: 1,
					Model:        []string{"Low", "Middle", "High"},
					OnCurrentIndexChanged: func() {

					},
				},
				VSpacer{
					MinSize: Size{Width: 30},
				},
				CheckBox{
					AssignTo: &checkJPEG,
					Checked:  ConfigGet().JpegEnable,
					Text:     "JPEG",
					OnCheckedChanged: func() {
						qualityNum.SetEnabled(checkJPEG.Checked())
						err := JpegEnableSave(checkJPEG.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				Label{
					Text: "JPEG Quality: ",
				},
				NumberEdit{
					AssignTo:    &qualityNum,
					Value:       float64(ConfigGet().JpegQuality),
					ToolTipText: "10~100",
					MaxValue:    100,
					MinValue:    10,
					OnValueChanged: func() {
						err := JpegQualitySave(int(qualityNum.Value()))
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				VSpacer{
					MinSize: Size{Width: 30},
				},
				Label{
					Text: "Parallel Tasks: ",
				},
				NumberEdit{
					AssignTo:    &taskNum,
					Value:       float64(ConfigGet().TaskNum),
					ToolTipText: "1~20",
					MaxValue:    20,
					MinValue:    1,
					OnValueChanged: func() {
						err := TaskNumSave(int(taskNum.Value()))
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
			},
		},
		Composite{
			Layout:   HBox{MarginsZero: true},
			Children: ProcessWidget(),
		},
	}
}
