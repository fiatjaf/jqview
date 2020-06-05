package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/itchyny/gojq"
)

var a = app.New()
var input *widget.Entry = widget.NewEntry()
var filter *widget.Entry = widget.NewEntry()
var output *widget.Entry = widget.NewEntry()

func main() {
	initialfilter := "."
	initialcontent := `[{
  "fruit": "mango"
}, {
  "fruit": "banana"
}]`

	if len(os.Args) > 1 {
		initialfilter = os.Args[1]
	}
	if len(os.Args) > 2 {
		b, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			log.Fatal("failed to read " + os.Args[2] + " : " + err.Error())
		}
		initialcontent = string(b)
	} else {
		if m, _ := os.Stdin.Stat(); m.Mode()&os.ModeCharDevice != os.ModeCharDevice {
			b, err := ioutil.ReadAll(os.Stdin)
			if err == nil {
				initialcontent = string(b)
			}
		}
	}

	input.SetText(strings.TrimSpace(initialcontent))
	input.PlaceHolder = "JSON Input"
	input.OnChanged = refresh

	filter.SetText(initialfilter)
	filter.PlaceHolder = "jq filter"
	filter.OnChanged = refresh

	w := a.NewWindow("jqview")
	w.SetContent(widget.NewVBox(
		newScrollWithMinHeight(input, 100),
		filter,
		output,
	))

	go refresh("")

	w.ShowAndRun()
}

func refresh(_ string) {
	output.SetText(runJQ(context.Background(), input.Text, filter.Text))
}

func runJQ(
	ctx context.Context,
	input string,
	filter string,
) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
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
