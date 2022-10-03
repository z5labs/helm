package kubeyaml

import (
	"bytes"
	"fmt"
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
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	docs := bytes.Split(b, []byte("---"))
	for i, doc := range docs {
		doc = bytes.TrimSpace(doc)
		if len(doc) == 0 {
			continue
		}

		f, err := dir.OpenFile(fmt.Sprintf("%d.yaml", i), os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		_, err = f.Write(doc)
		if err != nil {
			return err
		}
	}

	return nil
}
