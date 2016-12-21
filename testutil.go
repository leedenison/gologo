package gologo

import "fmt"
import "reflect"
import "testing"

type Assert testing.T

type TestSubject struct {
	assert *Assert
	subject interface{}
}

func (a *Assert) That(subject interface{}) *TestSubject {
	return &TestSubject { assert: a, subject: subject }
}

func (ts *TestSubject) Equals(object interface{}) {
	if reflect.ValueOf(object).IsNil() {
		ts.assert.Error("Assert.That.Equals used with <nil>, use Assert.That.IsNil instead.")
	}

	if !reflect.ValueOf(ts.subject).IsNil() {
		if ts.subject != object {
			ts.assert.Error(fmt.Sprintf("Expected: %T(%v), Found: %T(%v)", object, object, ts.subject, ts.subject))
		}		
	} else {
		ts.assert.Error(fmt.Sprintf("Expected: %T(%v), Found: <nil>", object, object))
	}

}

func (ts *TestSubject) IsNil() {
	if !reflect.ValueOf(ts.subject).IsNil() {
		ts.assert.Error(fmt.Sprintf("Expected: <nil>, Found: %T(%v)", ts.subject, ts.subject))
	}
}
