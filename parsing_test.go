package figtree

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree_Parse(t *testing.T) {
	figs := With(Options{Germinate: true})
	assert.Nil(t, figs.Parse())
}

func TestTree_ParseFile(t *testing.T) {
	exts := []string{"yaml", "json", "ini"}
	for _, ext := range exts {
		t.Run(ext, func(t *testing.T) {
			p := filepath.Join(".", "test.config."+ext)
			figs := With(Options{Germinate: true})
			figs = figs.NewString("name", "", "name")
			figs = figs.WithValidator("name", AssureStringContains("yahuah"))
			figs = figs.NewInt("age", 0, "age")
			figs = figs.WithValidator("age", AssureIntInRange(17, 47))
			figs = figs.NewString("sex", "", "sex")
			figs = figs.WithValidator("sex", AssureStringHasSuffix("male"))
			err := figs.ParseFile(p)
			assert.NoError(t, err)
		})
	}
}
