package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type FileItem struct {
	Index      int
	InputFile  string
	OutputFile string
	Status     string

	checked bool
}

type FileModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items []*FileItem
}

func (n *FileModel) RowCount() int {
	return len(n.items)
}

func (n *FileModel) Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.InputFile
	case 2:
		return item.OutputFile
	case 3:
		return item.Status
	}
	panic("unexpected col")
}

func (n *FileModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *FileModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *FileModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)
		case 1:
			return c(a.InputFile < b.InputFile)
		case 2:
			return c(a.OutputFile < b.OutputFile)
		case 3:
			return c(a.Status < b.Status)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

const (
	STATUS_UNDO = ""
	STATUS_DONE = "done"
	STATUS_FAIL = "failed"
)

var consoleFileTable *FileModel
var tableView *walk.TableView

func init() {
	consoleFileTable = new(FileModel)
	consoleFileTable.items = make([]*FileItem, 0)
}

func FileTableActive(isPNG bool, isJpeg bool) {
	lt := consoleFileTable

	lt.Lock()
	defer lt.Unlock()

	if ConfigGet().InputDir == "" {
		ErrorBoxAction(mainWindow, "Please set input directory first!")
		return
	}

	if ConfigGet().OutputDir == "" {
		ErrorBoxAction(mainWindow, "Please set output directory first!")
		return
	}

	fileList, err := ReadFileList(ConfigGet().InputDir)
	if err != nil {
		ErrorBoxAction(mainWindow, err.Error())
		return
	}

	tableView.SetCurrentIndex(-1)
	lt.items = make([]*FileItem, 0)
	lt.PublishRowsReset()
	lt.Sort(lt.sortColumn, lt.sortOrder)

	totalNumber := 0
	if isPNG {
		totalNumber = totalNumber + len(fileList)
	}

	if isJpeg {
		totalNumber = totalNumber + len(fileList)
	}

	index := 0

	for _, file := range fileList {
		timestamp := time.Now().Format("2006-01-02T15-04-05.000000")

		if isPNG {
			item := new(FileItem)
			item.Index = index
			item.InputFile = file
			item.Status = STATUS_UNDO
			item.OutputFile = filepath.Join(ConfigGet().OutputDir,
				fmt.Sprintf("%s.png", timestamp))

			err := ConvertHeic2Png(item.InputFile, item.OutputFile, ConfigGet().PngCompLevel)
			if err != nil {
				logs.Error("covert heic to png fail, %s", err.Error())
				item.Status = STATUS_FAIL
			} else {
				item.Status = STATUS_DONE
			}

			lt.items = append(lt.items, item)
			index++

			ProcessUpdate(float32(index) / float32(totalNumber))
		}

		if isJpeg {
			item := new(FileItem)
			item.Index = index
			item.InputFile = file
			item.Status = STATUS_UNDO
			item.OutputFile = filepath.Join(ConfigGet().OutputDir,
				fmt.Sprintf("%s.jpeg", timestamp))

			err := ConvertHeic2Jpeg(item.InputFile, item.OutputFile, ConfigGet().JpegQuality)
			if err != nil {
				logs.Error("covert heic to jpeg fail, %s", err.Error())
				item.Status = STATUS_FAIL
			} else {
				item.Status = STATUS_DONE
			}

			lt.items = append(lt.items, item)
			index++

			ProcessUpdate(float32(index) / float32(totalNumber))
		}

		lt.PublishRowsReset()
		lt.Sort(lt.sortColumn, lt.sortOrder)
	}
}

func TableWidget() []Widget {
	return []Widget{
		Label{
			Text: "Output File List: ",
		},
		TableView{
			AssignTo:         &tableView,
			AlternatingRowBG: true,
			ColumnsOrderable: true,
			CheckBoxes:       false,
			OnItemActivated: func() {
			},
			Columns: []TableViewColumn{
				{Title: "No", Width: 30},
				{Title: "InputFile", Width: 200},
				{Title: "OutputFile", Width: 200},
				{Title: "Status", Width: 60},
			},
			StyleCell: func(style *walk.CellStyle) {
				if style.Row()%2 == 0 {
					style.BackgroundColor = walk.RGB(248, 248, 255)
				} else {
					style.BackgroundColor = walk.RGB(220, 220, 220)
				}
			},
			Model: consoleFileTable,
		},
	}
}
