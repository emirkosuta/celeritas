package storage

import (
	"testing"
)

func TestStorage_MoveFile(t *testing.T) {

	err := storageClient.CreateDir("dir1")
	if err != nil {
		t.Error(err)
	}

	err = storageClient.CreateDir("dir2")
	if err != nil {
		t.Error(err)
	}

	err = storageClient.CreateFile("dir1/test.txt", nil)
	if err != nil {
		t.Error(err)
	}

	err = storageClient.Move("dir1/test.txt", "dir2")
	if err != nil {
		t.Error(err)
	}

	err = storageClient.Move("dir2/test.txt", "dir3")
	if err != nil {
		t.Error(err)
	}

	err = storageClient.Copy("dir3", "dir2")
	if err != nil {
		t.Error(err)
	}

	err = storageClient.DeleteDir("dir3")
	if err != nil {
		t.Error(err)
	}

	err = storageClient.EmptyDir("")
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}
}
