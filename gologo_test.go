package gologo

import (
	"io/ioutil"
	"os"
	"testing"
)

func testSetup() {
	InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func TestMain(m *testing.M) {
	testSetup()
	code := m.Run()
	os.Exit(code)
}
