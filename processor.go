package convol

type convolProcessStatus string

const (
	convolProcessStatusPending   = "pending"
	convolProcessStatusSucceeded = "succeeded"
	convolProcessStatusErr       = "errored"
)

type convolProcessor interface {
	Process(convolFunc ConvolFunc) ConvolFunc
}

type processor struct {
	name   string
	fn     ConvolFunc
	cvl    *convol
	result ConvolResult
}

func (p *processor) Process(next ConvolFunc) ConvolFunc {
	return func(ctx *ConvolCtx) error {
		p.result.status = convolProcessStatusSucceeded
		err := p.fn(ctx)
		if err != nil {
			ctx.lastErr = err
			ctx.Errored = true
			ctx.lastFailedName = p.name

			p.result.err = err
			p.result.status = convolProcessStatusErr

			if StrictLevel <= p.cvl.RunLevel {
				return err
			}
		}

		if next == nil {
			return err
		}

		return next(ctx)
	}
}
