package convol

type ConvolCtx struct {
	Errored        bool
	lastErr        error
	lastFailedName string
	lastResult     interface{}
	Store          map[string]interface{}
	Extra          interface{}
}

func (ctx *ConvolCtx) LastResult() interface{} {
	return ctx.lastResult
}

func (ctx *ConvolCtx) LastError() error {
	return ctx.lastErr
}
