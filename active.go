package main

import (
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var covert, cancel *walk.PushButton

func ActiveWidget() []Widget {
	return []Widget{
		PushButton{
			AssignTo: &covert,
			Text:     "Covert",
			OnClicked: func() {
				ProcessUpdate(0)
				covert.SetEnabled(false)
				cancel.SetEnabled(false)

				go func() {

					defer func() {
						ProcessUpdate(1)
						time.Sleep(time.Second)
						covert.SetEnabled(true)
						cancel.SetEnabled(true)
						ProcessUpdate(0)
					}()

					if !ConfigGet().PngEnable && !ConfigGet().JpegEnable {
						ErrorBoxAction(mainWindow, "Please select one or both of the PNG or JPEG options!")
						return
					} else {
						FileConvertActive()
					}
				}()
			},
		},
		PushButton{
			AssignTo: &cancel,
			Text:     "Cancel",
			OnClicked: func() {
				CloseWindows()
			},
		},
		HSpacer{},
	}
}
