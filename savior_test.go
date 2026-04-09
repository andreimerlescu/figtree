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
	nameFig := fig2.FigFlesh("name")
	assert.NotNil(t, nameFig)
	name := fig2.String("name")
	assert.Equal(t, t.Name(), *name)
	assert.NoError(t, os.RemoveAll(testFile))
}

func TestFigTree_SaveTo_MapRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.yaml")

	figs := With(Options{Germinate: true})
	figs.NewMap("cfg", map[string]string{"key": "value", "foo": "bar"}, "usage")
	// intentionally do NOT call StoreMap — use the raw MapFlag state from NewMap
	assert.NoError(t, figs.SaveTo(path))

	figs2 := With(Options{Germinate: true})
	figs2.NewMap("cfg", map[string]string{}, "usage")
	assert.NoError(t, figs2.ReadFrom(path))

	result := *figs2.Map("cfg")
	assert.Equal(t, map[string]string{"key": "value", "foo": "bar"}, result,
		"MapFlag should be unwrapped before serialization — got %v", result)
}

func TestFigTree_SaveTo_ListRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.yaml")

	figs := With(Options{Germinate: true})
	figs.NewList("items", []string{"one", "two", "three"}, "usage")
	// intentionally do NOT call StoreList — use the raw ListFlag state from NewList
	assert.NoError(t, figs.SaveTo(path))

	figs2 := With(Options{Germinate: true})
	figs2.NewList("items", []string{}, "usage")
	assert.NoError(t, figs2.ReadFrom(path))

	result := *figs2.List("items")
	assert.Equal(t, []string{"one", "three", "two"}, result,
		"ListFlag should be unwrapped before serialization — got %v", result)
}
