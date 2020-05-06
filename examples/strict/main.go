package main

import (
	"fmt"
	"log"

	"github.com/arabian9ts/convol"
)

func main() {
	cnvl := convol.New()
	cnvl.Add("handler1", handler_1)
	cnvl.Add("handler2", handler_2)
	cnvl.Add("handler3", handler_3)
	cnvl.Build()

	ctx := &convol.ConvolCtx{}
	err := cnvl.DoWithContext(ctx)
	fmt.Printf("err: %s\n", err)
	cnvl.DumpAllStatus()
}

func handler_1(*convol.ConvolCtx) error {
	log.Println("handler 1")
	return fmt.Errorf("error occurred at1")
}

func handler_2(*convol.ConvolCtx) error {
	log.Println("handler 2")
	return fmt.Errorf("error occurred at2")
}

func handler_3(*convol.ConvolCtx) error {
	log.Println("handler 3")
	return nil
}
