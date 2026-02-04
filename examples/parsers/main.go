package main

import (
	"fmt"

	"github.com/tuanta7/errx"
	"github.com/tuanta7/errx/parsers/json"
	"github.com/tuanta7/errx/registry"
	lang "golang.org/x/text/language"
)

func main() {
	registry.SetGlobal(registry.New())
	err := registry.Global.LoadMessages(lang.English.String(), "./static/en.json", json.Parser())
	if err != nil {
		panic(err)
	}

	err = registry.Global.LoadMessages(lang.Spanish.String(), "./static/es.json", json.Parser())
	if err != nil {
		panic(err)
	}

	err = registry.Global.LoadMessages(lang.Vietnamese.String(), "./static/vi.json", json.Parser())
	if err != nil {
		panic(err)
	}

	ErrNotFound := errx.New("default not found message").WithCode("ERR_RESOURCE_NOT_FOUND")

	fmt.Println(registry.Global.GetMessage(ErrNotFound, lang.English.String()))
	fmt.Println(registry.Global.GetMessage(ErrNotFound, lang.Spanish.String()))
	fmt.Println(registry.Global.GetMessage(ErrNotFound, lang.Vietnamese.String()))
}
