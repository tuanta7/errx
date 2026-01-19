package main

import (
	"fmt"

	"github.com/tuanta7/errx"
	"github.com/tuanta7/errx/errors"
	"github.com/tuanta7/errx/parsers/json"
	lang "golang.org/x/text/language"
)

func main() {
	errx.SetGlobal(errx.New())
	err := errx.Global.LoadMessages(lang.English.String(), "./static/en.json", json.Parser())
	if err != nil {
		panic(err)
	}

	err = errx.Global.LoadMessages(lang.Spanish.String(), "./static/es.json", json.Parser())
	if err != nil {
		panic(err)
	}

	err = errx.Global.LoadMessages(lang.Vietnamese.String(), "./static/vi.json", json.Parser())
	if err != nil {
		panic(err)
	}

	ErrNotFound := errors.New("default not found message").WithCode("ERR_RESOURCE_NOT_FOUND")

	fmt.Println(errx.Global.GetMessage(ErrNotFound, lang.English.String()))
	fmt.Println(errx.Global.GetMessage(ErrNotFound, lang.Spanish.String()))
	fmt.Println(errx.Global.GetMessage(ErrNotFound, lang.Vietnamese.String()))
}
