package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/itchyny/gojq"
	"github.com/mitchellh/go-homedir"
	"github.com/therecipe/qt/widgets"
)

var filterValue = "."
var inputValue = `[{
  "fruit": "mango"
}, {
  "fruit": "banana"
}]`

var loadfileDialog *widgets.QFileDialog
var input *widgets.QPlainTextEdit
var filter *widgets.QLineEdit
var output *widgets.QPlainTextEdit

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	if len(os.Args) > 1 {
		filterValue = os.Args[1]
	}
	if len(os.Args) > 2 {
		b, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			log.Fatal("failed to read " + os.Args[2] + " : " + err.Error())
		}
		inputValue = string(b)
	} else {
		if m, _ := os.Stdin.Stat(); m.Mode()&os.ModeCharDevice != os.ModeCharDevice {
			b, err := ioutil.ReadAll(os.Stdin)
			if err == nil {
				inputValue = string(b)
			}
		}
	}

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(400, 500)
	window.SetWindowTitle("jqview")

	dir, _ := homedir.Dir()
	loadfileDialog = widgets.NewQFileDialog2(nil, "Select a JSON file", dir, "")
	loadfileDialog.ConnectFileSelected(func(filepath string) {
		b, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Print("failed to read " + filepath + " : " + err.Error())
		} else {
			input.SetPlainText(string(b))
			go refresh()
		}
	})

	loadfileButton := widgets.NewQPushButton2("Load", nil)
	loadfileButton.SetMaximumWidth(40)
	loadfileButton.ConnectClicked(func(_ bool) {
		loadfileDialog.Open(nil, "")
	})

	input = widgets.NewQPlainTextEdit(nil)
	input.SetPlaceholderText("JSON input")
	input.SetPlainText(inputValue)
	input.ConnectTextChanged(refresh)

	inputSection := widgets.NewQWidget(nil, 0)
	inputSection.SetLayout(widgets.NewQHBoxLayout())
	inputSection.SetMaximumHeight(150)
	inputSection.Layout().AddWidget(loadfileButton)
	inputSection.Layout().AddWidget(input)

	filter = widgets.NewQLineEdit(nil)
	filter.SetPlaceholderText("jq filter")
	filter.SetText(filterValue)
	filter.ConnectTextChanged(func(value string) {
		filterValue = value
		go refresh()
	})

	output = widgets.NewQPlainTextEdit(nil)
	output.SetSizeAdjustPolicy(widgets.QAbstractScrollArea__AdjustToContents)
	output.SetMinimumHeight(300)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	widget.Layout().AddWidget(inputSection)
	widget.Layout().AddWidget(filter)
	widget.Layout().AddWidget(output)

	window.Show()
	window.SetCentralWidget(widget)

	go refresh()

	app.Exec()
}

func refresh() {
	inputValue = input.ToPlainText()
	output.SetPlainText(runJQ(context.Background(), inputValue, filterValue))
}

func runJQ(
	ctx context.Context,
	input string,
	filter string,
) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	var object interface{}
	err := json.Unmarshal([]byte(input), &object)
	if err != nil {
		return err.Error()
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return err.Error()
	}

	iter := query.RunWithContext(ctx, object)

	var results []string
	for {
		v, exists := iter.Next()
		if !exists {
			break
		}

		if err, ok := v.(error); ok {
			return err.Error()
		}

		s, _ := json.MarshalIndent(v, "", "  ")
		results = append(results, string(s))
	}

	return strings.Join(results, "\n")
}
