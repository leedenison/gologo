package gologo

import "fmt"
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
	if object == nil {
		ts.assert.Error("Assert.That.Equals used with <nil>, use Assert.That.IsNil instead.")
	}

	if ts.subject != nil {
		if ts.subject != object {
			ts.assert.Error(fmt.Sprintf("Expected: %v, Found: %v", object, ts.subject))
		}		
	} else {
		ts.assert.Error(fmt.Sprintf("Expected: %v, Found: <nil>", object))
	}

}

func (ts *TestSubject) IsNil() {
	if ts.subject != nil {
		ts.assert.Error(fmt.Sprintf("Expected: <nil>, Found: %v", ts.subject))
	}
}
