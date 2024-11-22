package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var processBar *walk.ProgressBar

func ProcessWidget() []Widget {
	return []Widget{
		ProgressBar{
			AssignTo: &processBar,
			MaxValue: 100,
			MinValue: 0,
			MinSize:  Size{Height: 5},
			MaxSize:  Size{Height: 5},
		},
	}
}

// 0.00 - 1.00
func ProcessUpdate(value float32) {
	processBar.SetValue(int(value * 100))
}
