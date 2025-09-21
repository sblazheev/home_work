package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	from, to      string
	limit, offset int64
	totalSteps    = 50
)

type Bar struct {
	cur   int64
	total int64
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	Progress = make(chan Bar)

	go func() {
		for {
			val, ok := <-Progress
			if !ok {
				break
			}
			filled := int(val.getPercent() / 2)
			bar := "[" + strings.Repeat("#", filled) + strings.Repeat("-", totalSteps-filled) + "]"
			fmt.Printf("\r%d / %d %s %3d%%", val.cur, val.total, bar, val.getPercent())
		}
	}()

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("")
	}
}
