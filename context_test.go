package convol

import (
	"fmt"
	"testing"
)

func Test_ConvolCtx_LastResult(t *testing.T) {
	t.Parallel()

	t.Run("lastResult=nil,resp=nil", func(t *testing.T) {
		ctx := &ConvolCtx{}
		res := ctx.LastResult()
		if res != nil {
			t.Errorf("result must be nil, but actual is %v", res)
		}
	})

	t.Run("lastResult=struct{},resp=struct{}", func(t *testing.T) {
		ctx := &ConvolCtx{lastResult: struct{}{}}
		res := ctx.LastResult()
		if res != struct{}{} {
			t.Errorf("result must be empty struct{}{}, but actual is %v", res)
		}
	})
}

func Test_ConvolCtx_LastError(t *testing.T) {
	t.Parallel()

	t.Run("lastErr=nil,resp=nil", func(t *testing.T) {
		ctx := &ConvolCtx{}
		err := ctx.LastError()
		if err != nil {
			t.Errorf("result must be nil, but actual is %v", err)
		}
	})

	t.Run("lastErr=err,resp=err", func(t *testing.T) {
		testErr := fmt.Errorf("test error")
		ctx := &ConvolCtx{lastErr: testErr}
		err := ctx.LastError()
		if err == nil {
			t.Errorf("result must be error, but actual is nil")
		}
	})
}
