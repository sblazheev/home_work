package main

import (
	"errors"
	"flag"
)

var (
	from, to      string
	limit, offset int64
)

type InputParamsError struct {
	Params string
	Value  string
	err    error
}

func (e *InputParamsError) Error() string {
	return e.Params + ": " + e.Value + "-" + e.err.Error()
}

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	err := copeFile("", "", limit, offset)
	if err != nil {
		println(err.Error())
	}
}

func copeFile(from, to string, limit, offset int64) (err error) {
	if from == "" {
		return &InputParamsError{Params: "from", Value: from, err: errors.New("Empty param")}
	}
	if to == "" {
		return &InputParamsError{Params: "to", Value: from, err: errors.New("Empty param")}
	}
	return
}
