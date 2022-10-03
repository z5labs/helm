package kubeyaml

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	t.Run("it should completely write each document", func(t *testing.T) {
		t.Run("if each document ends with a ---", func(t *testing.T) {
			docOne := `hello: world
a: b`
			docTwo := `good: bye
1: 2`
			docs := docOne + "---\n" + docTwo + "---\n"

			dir := afero.NewMemMapFs()
			err := Split(strings.NewReader(docs), dir)
			if !assert.Nil(t, err) {
				return
			}

			one, err := afero.ReadFile(dir, "0.yaml")
			if !assert.Nil(t, err) {
				return
			}
			if !assert.Equal(t, docOne, string(one)) {
				return
			}

			two, err := afero.ReadFile(dir, "1.yaml")
			if !assert.Nil(t, err) {
				return
			}
			if !assert.Equal(t, docTwo, string(two)) {
				return
			}
		})

		t.Run("even if the last document doesn't end with a ---", func(t *testing.T) {
			docOne := `hello: world
a: b`
			docTwo := `good: bye
1: 2`
			docs := docOne + "---\n" + docTwo

			dir := afero.NewMemMapFs()
			err := Split(strings.NewReader(docs), dir)
			if !assert.Nil(t, err) {
				return
			}

			one, err := afero.ReadFile(dir, "0.yaml")
			if !assert.Nil(t, err) {
				return
			}
			if !assert.Equal(t, docOne, string(one)) {
				return
			}

			two, err := afero.ReadFile(dir, "1.yaml")
			if !assert.Nil(t, err) {
				return
			}
			if !assert.Equal(t, docTwo, string(two)) {
				return
			}
		})
	})
}
