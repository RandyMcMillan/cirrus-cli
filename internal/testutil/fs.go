package testutil

import (
	"github.com/otiai10/copy"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func RestoreDirLater(t *testing.T) {
	savedDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(savedDir); err != nil {
			t.Fatal(err)
		}
	})
}

// tempDir supplements an alternative to TB.TempDir()[1], which is only available in 1.15.
// [1]: https://github.com/golang/go/issues/35998
func TempDir(t *testing.T) string {
	tempDirRoot := "" // will use os.TempDir()
	if runtime.GOOS == "darwin" {
		// override the default since Docker for Mac
		// doesn't mount /var/folder by default where os.TempDir() will be located
		// See https://docs.docker.com/docker-for-mac/#file-sharing
		tempDirRoot = "/tmp"
	}

	dir, err := ioutil.TempDir(tempDirRoot, filepath.Base(t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	})

	return dir
}

// TempChdir switches to a temporary per-test directory.
func TempChdir(t *testing.T) {
	dir := TempDir(t)

	RestoreDirLater(t)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

// TempChdirPopulatedWith creates a temporary per-test directory
// filled with sourceDir contents and switches to it.
func TempDirPopulatedWith(t *testing.T, sourceDir string) string {
	tempDir := TempDir(t)

	if err := copy.Copy(sourceDir, tempDir); err != nil {
		t.Fatal(err)
	}

	return tempDir
}

// TempChdirPopulatedWith creates a temporary per-test directory
// filled with sourceDir contents and switches to it.
func TempChdirPopulatedWith(t *testing.T, sourceDir string) {
	dir := TempDirPopulatedWith(t, sourceDir)

	RestoreDirLater(t)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}
