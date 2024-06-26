// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
package extensionerrors

import (
	"fmt"
	"strings"
	"testing"
)

func ErrorWithStack(stackDepth int, err error) error {
	if stackDepth == 0 {
		if err != nil {
			return AddStackToError(err)
		} else {
			return NewErrorWithStack("Reached the bottom of the error call stack")
		}
	} else {
		return ErrorWithStack(stackDepth-1, err)
	}
}

func ErrorWithoutStack(stackDepth int, err error) error {
	if stackDepth == 0 {
		if err != nil {
			return err
		} else {
			return fmt.Errorf("reached the bottom of the error call stack")
		}
	} else {
		return ErrorWithoutStack(stackDepth-1, err)
	}
}

func TestErrorsWrap(t *testing.T) {
	err1 := ErrorWithStack(5, fmt.Errorf("misc error"))
	err2 := ErrorWithStack(5, nil)

	t.Logf("the error stack trace is: %+v", err1)
	t.Logf("the error stack trace is: %+v", err2)

	if !strings.Contains(err1.Error(), "ErrorWithStack") {
		t.Fatal("error doesn't contain stack trace")
	}

	if !strings.Contains(err2.Error(), "ErrorWithStack") {
		t.Fatal("error doesn't contain stack trace")
	}

	if !strings.Contains(err1.Error(), "misc error") {
		t.Fatal("error doesn't base error")
	}

	if !strings.Contains(err2.Error(), "Reached the bottom of the error call stack") {
		t.Fatal("error doesn't base error")
	}
}

func TestPanicWithCallStack(t *testing.T) {
	err := ErrorWithStack(5, fmt.Errorf("misc error"))
	defer func() {
		if r := recover(); r != nil {
			errorString := fmt.Sprintf("%+v", r)
			t.Logf("Recovered %s", errorString)
			if !strings.Contains(errorString, err.Error()) {
				t.Fatal("error doesn't base error")
			}
			if !strings.Contains(errorString, "ErrorWithStack") {
				t.Fatal("error doesn't the call stack")
			}
		} else {
			t.Fatal("error not recovered")
		}
	}()

	panic(err)
}

func TestPanicWithoutCallStack(t *testing.T) {
	err := ErrorWithoutStack(5, fmt.Errorf("misc error"))
	defer func() {
		if r := recover(); r != nil {
			errorString := fmt.Sprintf("%+v", r)
			t.Logf("Recovered %s", errorString)
			if !strings.Contains(errorString, err.Error()) {
				t.Fatal("error doesn't base error")
			}
			if strings.Contains(errorString, "ErrorWithoutStack") {
				t.Fatal("error has call stack")
			}
		} else {
			t.Fatal("error not recovered")
		}
	}()

	panic(err)
}
