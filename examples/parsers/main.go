package main

import (
	"fmt"

	"github.com/tuanta7/errx"
	"github.com/tuanta7/errx/errors"
	"github.com/tuanta7/errx/parsers/json"
)

func main() {
	errx.SetGlobal(errx.New())
	err := errx.Global.LoadMessages("./errors.json", json.Parser())
	if err != nil {
		panic(err)
	}

	ErrNotFound := errors.New("default not found message").WithCode("ERR_RESOURCE_NOT_FOUND")

	fmt.Println(errx.Global.GetMessage(ErrNotFound, "en"))
	fmt.Println(errx.Global.GetMessage(ErrNotFound, "es"))
	fmt.Println(errx.Global.GetMessage(ErrNotFound, "vi"))
}
