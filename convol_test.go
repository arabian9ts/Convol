package convol

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func Test_identityFunc(t *testing.T) {
	t.Run("lastErr=nil,resp=nil", func(t *testing.T) {
		ctx := &ConvolCtx{}
		err := identityFunc(ctx)
		if err != nil {
			t.Errorf("result must be nil but actual is %v", err)
		}
	})

	t.Run("lastErr=err,resp=err", func(t *testing.T) {
		testErr := fmt.Errorf("test error")
		ctx := &ConvolCtx{lastErr: testErr}
		err := identityFunc(ctx)
		if err == nil {
			t.Errorf("result must be error but actual is nil")
		}
		if err != testErr {
			t.Errorf("result error is not correct")
		}
	})
}

func Test_Convol_Add(t *testing.T) {
	convol := &convol{}

	t.Run("name=empty", func(t *testing.T) {
		err := convol.Add("", identityFunc)
		if err == nil {
			t.Errorf("empty name function should not be allowd to registered")
		}
	})

	t.Run("name=test", func(t *testing.T) {
		funcName := "test"
		err := convol.Add(funcName, identityFunc)
		if err != nil {
			t.Errorf("function must be registered with no error")
		}

		if len(convol.convolutions) != 1 {
			t.Errorf("registered function count is not correct")
		}

		p := convol.convolutions[0].(*processor)

		if p.name != funcName {
			t.Errorf("registered function name is not currect")
		}

		if p.cvl != convol {
			t.Errorf("convol is not set to processor")
		}

		if p.fn == nil {
			t.Errorf("function is not set currectly")
		}

		if p.result.err != nil {
			t.Errorf("initialized err must be nil")
		}

		if p.result.status != convolProcessStatusPending {
			t.Errorf("initialized status must be pending")
		}
	})

	t.Run("name=duplicated_name", func(t *testing.T) {
		err := convol.Add("test", identityFunc)
		if err == nil {
			t.Errorf("duplicated name function should not be allowed")
		}
	})
}

func Test_Convol_Build(t *testing.T) {
	convol := &convol{}
	convol.Add("first", identityFunc)
	convol.Add("second", identityFunc)
	convol.Build()

	if convol.builtFunc == nil {
		t.Errorf("built func must be set")
	}
}

func Test_Convol_Do(t *testing.T) {
	t.Run("run_order", func(t *testing.T) {
		// prepare stdout capture
		r, w, e := os.Pipe()
		if e != nil {
			t.Errorf("os pipe error")
		}

		stdout := os.Stdout
		os.Stdout = w

		// build and do
		// each func print its name to stdout
		convol := &convol{}
		convol.Add("first", func(*ConvolCtx) error { fmt.Println("first"); return nil })
		convol.Add("second", func(*ConvolCtx) error { fmt.Println("second"); return nil })
		convol.Build()
		err := convol.Do()

		// finalize stdout
		os.Stdout = stdout
		w.Close()

		// collect function name from stdout
		outs := make([]string, 0)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			outs = append(outs, scanner.Text())
		}

		if err != nil {
			t.Errorf("run result should be nil")
		}

		// convol runs functions in reverse order of adding
		// finally, convol runs identity
		expected := []string{"second", "first"}
		for i := range outs {
			if outs[i] != expected[i] {
				t.Errorf("function calling order is not correct")
			}
		}
	})

	t.Run("permissive", func(t *testing.T) {
		convol := &convol{RunLevel: PermissiveLevel}
		firstErr := fmt.Errorf("first errir")
		secondErr := fmt.Errorf("second errir")

		convol.Add("first", func(*ConvolCtx) error { return firstErr })
		convol.Add("second", func(*ConvolCtx) error { return secondErr })
		convol.Build()

		// run order
		// second -> first -> identity
		err := convol.Do()

		if err != firstErr {
			t.Errorf("in permissive level, later functions are also run")
		}

		for i := range convol.convolutions {
			p := convol.convolutions[i].(*processor)

			if p.result.err == nil {
				t.Errorf("all functions that returns error are run")
			}

			if p.result.status != convolProcessStatusErr {
				t.Errorf("all functions that returns error are run and set status")
			}
		}
	})

	t.Run("strict", func(t *testing.T) {
		convol := &convol{RunLevel: StrictLevel}
		firstErr := fmt.Errorf("first errir")
		secondErr := fmt.Errorf("second errir")

		convol.Add("first", func(*ConvolCtx) error { return firstErr })
		convol.Add("second", func(*ConvolCtx) error { return secondErr })
		convol.Build()

		// run order
		// second -> first -> identity
		err := convol.Do()

		if err != secondErr {
			t.Errorf("in strict level, later functions should be canceled")
		}

		p_0 := convol.convolutions[0].(*processor)
		p_1 := convol.convolutions[1].(*processor)

		// in strict level, p_0 is canceled because p_1 returns err
		if p_0.result.err != nil {
			t.Errorf("later function should be canceled then err should be nil")
		}

		if p_0.result.status != convolProcessStatusPending {
			t.Errorf("later function should be canceled and keep in pending")
		}

		// p_1 returns err
		if p_1.result.err == nil {
			t.Errorf("all functions that returns error are run")
		}

		if p_1.result.status != convolProcessStatusErr {
			t.Errorf("all functions that returns error are run and set status")
		}
	})
}
