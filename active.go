package main

import (
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
				covert.SetEnabled(false)
				cancel.SetEnabled(false)

				go func() {
					FileTableActive(true, true)
					covert.SetEnabled(true)
					cancel.SetEnabled(true)
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
