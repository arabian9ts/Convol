package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arabian9ts/convol"
)

func main() {
	cnvl := convol.New()
	cnvl.SetRunLevel(convol.StrictLevel)
	cnvl.Add("calcC", calcC)
	cnvl.Add("calcB", calcB)
	cnvl.Add("calcA", calcA)
	cnvl.Build()

	ctx := &convol.ConvolCtx{}
	err := cnvl.DoWithContext(ctx)
	fmt.Printf("err: %s\n", err)
	cnvl.DumpAllStatus()
}

func calcA(ctx *convol.ConvolCtx) error {
	ctx.Store["A"] = time.Now().Unix()
	return nil
}

func calcB(ctx *convol.ConvolCtx) error {
	log.Println(ctx.Store)
	ctx.Store["B"] = ctx.Store["A"].(int64) / 2
	return nil
}

func calcC(ctx *convol.ConvolCtx) error {
	ans := ctx.Store["A"].(int64) / ctx.Store["B"].(int64)
	fmt.Printf("answer: %d\n", ans)
	ctx.Store["Ans"] = ans
	return nil
}
