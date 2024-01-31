package gologo

import (
	"io"
	"os"
	"testing"

	"github.com/leedenison/gologo/log"
)

func testSetup() {
	log.InitLogger(io.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func TestMain(m *testing.M) {
	testSetup()
	code := m.Run()
	os.Exit(code)
}
