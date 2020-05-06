package convol

import (
	"fmt"
	"testing"
	"time"
)

func Test_Processor_Process(t *testing.T) {
	t.Run("next=nil,level=strict", func(t *testing.T) {
		convol := &convol{RunLevel: StrictLevel}
		caseName := fmt.Sprintf("%d", time.Now().Unix())
		p := &processor{name: caseName, fn: identityFunc, cvl: convol}

		result := p.Process(nil)

		t.Run("lastErr=nil", func(t *testing.T) {
			ctx := &ConvolCtx{lastErr: nil}
			err := result(ctx)

			if err != nil {
				t.Errorf("err must be nil, but actual is not nil")
			}

			if p.result.err != nil {
				t.Errorf("err must be nil, but actual is not nil")
			}

			if p.result.status != convolProcessStatusSucceeded {
				t.Errorf("processor status must be succeeded")
			}
		})

		t.Run("lastErr=error", func(t *testing.T) {
			now := time.Now().Unix()
			e := fmt.Errorf("test-%d", now)
			ctx := &ConvolCtx{lastErr: e}
			err := result(ctx)

			if err == nil {
				t.Errorf("err must be error, but actual is nil")
			}

			if p.result.err == nil {
				t.Errorf("err must be error, but actual is nil")
			}

			if p.result.status != convolProcessStatusErr {
				t.Errorf("processor status must be errored")
			}
		})
	})

	t.Run("next=nil,level=permissive", func(t *testing.T) {
		convol := &convol{RunLevel: PermissiveLevel}
		caseName := fmt.Sprintf("%d", time.Now().Unix())
		p := &processor{name: caseName, fn: identityFunc, cvl: convol}

		result := p.Process(nil)

		t.Run("lastErr=nil", func(t *testing.T) {
			ctx := &ConvolCtx{lastErr: nil}
			err := result(ctx)

			if err != nil {
				t.Errorf("err must be nil, but actual is not nil")
			}

			if p.result.err != nil {
				t.Errorf("err must be nil, but actual is not nil")
			}

			if p.result.status != convolProcessStatusSucceeded {
				t.Errorf("processor status must be succeeded")
			}
		})

		t.Run("lastErr=error", func(t *testing.T) {
			now := time.Now().Unix()
			e := fmt.Errorf("test-%d", now)
			ctx := &ConvolCtx{lastErr: e}
			err := result(ctx)

			if err == nil {
				t.Errorf("err must be error, but actual is nil")
			}

			if p.result.err == nil {
				t.Errorf("err must be error, but actual is nil")
			}

			if p.result.status != convolProcessStatusErr {
				t.Errorf("processor status must be errored")
			}
		})
	})

	t.Run("next=non_nil,level=strict", func(t *testing.T) {
		convol := &convol{RunLevel: StrictLevel}
		caseName := fmt.Sprintf("%d", time.Now().Unix())
		p := &processor{name: caseName, fn: identityFunc, cvl: convol}

		nextErr := fmt.Errorf("test-%d", time.Now().Unix())
		next := func(*ConvolCtx) error { return nextErr }
		result := p.Process(next)

		t.Run("lastErr=nil", func(t *testing.T) {
			ctx := &ConvolCtx{lastErr: nil}
			err := result(ctx)

			// inspect next process error
			// because lastErr is nil, next is called
			if err != nextErr {
				t.Errorf("err must be equal to next func error")
			}

			// inspect fn error
			if p.result.err != nil {
				t.Errorf("err must be nil, but actual is not nil")
			}

			// inspect fn error
			if p.result.status != convolProcessStatusSucceeded {
				t.Errorf("processor status must be succeeded")
			}
		})

		t.Run("lastErr=error", func(t *testing.T) {
			now := time.Now().Unix()
			e := fmt.Errorf("test-%d", now)
			ctx := &ConvolCtx{lastErr: e}
			err := result(ctx)

			// in strict level, next should not be called and processor returns fn error soon
			if err != e {
				t.Errorf("err must be error that is defined in context")
			}

			if p.result.err == nil {
				t.Errorf("err must be error, but actual is nil")
			}

			if p.result.status != convolProcessStatusErr {
				t.Errorf("processor status must be errored")
			}
		})
	})

	t.Run("next=non_nil,level=permissive", func(t *testing.T) {
		convol := &convol{RunLevel: PermissiveLevel}
		caseName := fmt.Sprintf("%d", time.Now().Unix())
		p := &processor{name: caseName, fn: identityFunc, cvl: convol}

		nextErr := fmt.Errorf("test-%d", time.Now().Unix())
		next := func(*ConvolCtx) error { return nextErr }
		result := p.Process(next)

		t.Run("lastErr=nil", func(t *testing.T) {
			ctx := &ConvolCtx{lastErr: nil}
			err := result(ctx)

			// inspect next process error
			// because lastErr is nil, next is called
			if err != nextErr {
				t.Errorf("err must be equal to next func error")
			}

			// inspect fn error
			if p.result.err != nil {
				t.Errorf("err must be nil, but actual is not nil")
			}

			// inspect fn error
			if p.result.status != convolProcessStatusSucceeded {
				t.Errorf("processor status must be succeeded")
			}
		})

		t.Run("lastErr=error", func(t *testing.T) {
			now := time.Now().Unix()
			e := fmt.Errorf("test-%d", now)
			ctx := &ConvolCtx{lastErr: e}
			err := result(ctx)

			// in permissive level, next should be called and processor returns next error
			if err != nextErr {
				t.Errorf("err must be error that next returns")
			}

			if p.result.err == nil {
				t.Errorf("err must be error, but actual is nil")
			}

			if p.result.status != convolProcessStatusErr {
				t.Errorf("processor status must be errored")
			}
		})
	})
}
