package gologo

import (
	"os"
	"testing"
)

func testSetup() {
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

func TestMain(m *testing.M) {
	testSetup()
	code := m.Run()
	os.Exit(code)
}
