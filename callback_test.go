package figtree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTree_WithCallback(t *testing.T) {
	figs := With(Options{Germinate: true, Pollinate: false, IgnoreEnvironment: true})
	figs = figs.NewString(t.Name(), t.Name(), "usage")
	figs = figs.WithCallback(t.Name(), CallbackAfterVerify, func(value interface{}) error {
		if value == nil {
			return nil
		}
		switch v := value.(type) {
		case *string:
			t.Logf("CallbackAfterVerify: %s on %s", *v, time.Now())
		case string:
			t.Logf("CallbackAfterVerify: %s on %s", v, time.Now())
		}
		return nil
	})
	figs = figs.WithCallback(t.Name(), CallbackAfterRead, func(value interface{}) error {
		if value == nil {
			return nil
		}
		switch v := value.(type) {
		case *string:
			t.Logf("CallbackAfterRead: %s on %s", *v, time.Now())
		case string:
			t.Logf("CallbackAfterRead: %s on %s", v, time.Now())
		}
		return nil
	})
	figs = figs.WithCallback(t.Name(), CallbackAfterChange, func(value interface{}) error {
		if value == nil {
			return nil
		}
		switch v := value.(type) {
		case *string:
			t.Logf("CallbackAfterChange: %s on %s", *v, time.Now())
		case string:
			t.Logf("CallbackAfterChange: %s on %s", v, time.Now())
		}
		return nil
	})
	time.Sleep(369 * time.Millisecond)
	assert.NoError(t, figs.Load())
	time.Sleep(369 * time.Millisecond)
	property := *figs.String(t.Name())
	assert.NotNil(t, property)
	time.Sleep(369 * time.Millisecond)
	assert.NotNil(t, figs.FigFlesh(t.Name()))
	figs.StoreString(t.Name(), "new value")
	time.Sleep(369 * time.Millisecond)

}
