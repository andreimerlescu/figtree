package figtree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFigTree_ReadFile(t *testing.T) {
	figs := With(Options{Germinate: true})
	assert.NoError(t, figs.ReadFrom(filepath.Join(".", "test.config.yaml")))
	assert.Equal(t, "yahuah", *figs.String("name"))
}

func TestFigTree_SaveTo(t *testing.T) {
	testFile := filepath.Join(".", "test.saved.config.yaml")

	figs := With(Options{Germinate: true})
	figs.NewString("name", "unknown", "name")
	figs.NewInt("age", 0, "age")
	figs.NewString("sex", "unknown", "sex")
	figs.StoreString("name", t.Name())
	figs.StoreInt("age", 33)
	figs.StoreString("sex", "male")
	assert.NoError(t, figs.SaveTo(testFile))

	fig2 := With(Options{Germinate: true})
	assert.NoError(t, fig2.ReadFrom(testFile))
	nameFig := fig2.Fig("name")
	assert.NotNil(t, nameFig)
	name := fig2.String("name")
	assert.Equal(t, t.Name(), *name)
	assert.NoError(t, os.RemoveAll(testFile))
}
