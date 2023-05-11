package storage

import (
	"os"
	"testing"
)

var storageClient = Storage{
	BaseDir: "./testdata",
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
