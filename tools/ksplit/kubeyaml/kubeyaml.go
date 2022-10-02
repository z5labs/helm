package kubeyaml

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

// FileOpener
type FileOpener interface {
	OpenFile(name string, flag int, perm os.FileMode) (afero.File, error)
}

// Split
func Split(r io.Reader, dir FileOpener) error {
	return nil
}
