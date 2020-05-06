package convol

import (
	"fmt"
)

type ConvolResult struct {
	err    error
	status convolProcessStatus
}

type ConvolFunc func(*ConvolCtx) error

func identityFunc(ctx *ConvolCtx) error {
	return ctx.lastErr
}

type ConvolRunLevel int

const (
	PermissiveLevel = iota + 1
	StrictLevel
)

type convol struct {
	RunLevel     ConvolRunLevel
	convolutions []convolProcessor
	builtFunc    ConvolFunc
}

func (c *convol) Add(name string, fn ConvolFunc) error {
	if name == "" {
		return fmt.Errorf("empty name function is not allowd")
	}

	for i := range c.convolutions {
		if c.convolutions[i].(*processor).name == name {
			return fmt.Errorf("%s is already registered", name)
		}
	}

	c.convolutions = append(c.convolutions, &processor{
		name: name,
		fn:   fn,
		cvl:  c,
		result: ConvolResult{
			err:    nil,
			status: convolProcessStatusPending,
		},
	})

	return nil
}

func (c *convol) Build() {
	c.builtFunc = func(ctx *ConvolCtx) error {
		ctx.Store = make(map[string]interface{})
		h := identityFunc
		for i := range c.convolutions {
			h = c.convolutions[i].Process(h)
		}
		return h(ctx)
	}
}

func (c *convol) DoWithContext(ctx *ConvolCtx) error {
	return c.builtFunc(ctx)
}

func (c *convol) Do() error {
	ctx := &ConvolCtx{}
	return c.builtFunc(ctx)
}

func New() *convol {
	return &convol{
		RunLevel: StrictLevel,
	}
}

func (c *convol) SetRunLevel(level ConvolRunLevel) {
	c.RunLevel = level
}

func (c *convol) DumpAllStatus() {
	for i := range c.convolutions {
		fmt.Printf("%#v\n", c.convolutions[i].(*processor).result)
	}
}
