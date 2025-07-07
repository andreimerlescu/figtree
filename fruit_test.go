package figtree

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValue_Set(t *testing.T) {
	type test struct {
		In      string
		Out     interface{}
		WantErr bool
	}
	t.Run("String", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), t.Name(), false},
			{"", "", false},
			{"abc", "abc", false},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tString}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
				} else {
					assert.NoError(t, v.Set(test.In))
				}
				assert.Equal(t, test.Out, v.Value)
			})
		}
	})
	t.Run("Bool", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), true, true},
			{"", false, false},
			{"true", true, false},
			{"1", true, false},
			{"false", false, false},
			{"0", false, false},
			{"9", false, true},
			{"[]", false, true},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tBool}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.NotEqual(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
	})
	t.Run("Int", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), 0, true},
			{"", 0, false},
			{"3.14", 3, true},
			{"abc", 0, true},
			{"-3.14", -3, true},
			{fmt.Sprintf("%d", math.MaxInt64), 0, true},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tInt}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.NotEqual(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
	})
	t.Run("Int64", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), int64(0), true},
			{"", int64(0), false},
			{"abc", int64(0), true},
			{"42", int64(42), false},
			{"3.14", int64(3), false},
			{"-3.14", int64(-3), false},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tInt64}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.NotEqual(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
	})
	t.Run("Float64", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), 0.0, true},
			{"", 0.0, false},
			{"3.14", 3.14, false},
			{"3", 3.0, false},
			{"3.3.3", 0.0, true},
			{"-3.3", -3.3, false},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tFloat64}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.NotEqual(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
	})
	t.Run("Duration", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), nil, true},
			{"", time.Duration(0), false},
			{"1", time.Duration(1), false},
			{"1h", time.Hour, false},
			{"1w", time.Hour * 168, false},
			{"1d7h76s", time.Duration(31)*time.Hour + time.Duration(76)*time.Second, false},
			{"abc", nil, true},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tDuration}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
	})
	t.Run("List", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), []string{t.Name()}, false},
			{"", []string{}, false},
			{"yes,no", []string{"yes", "no"}, false},
			{"test", []string{"test"}, false},
			{"abc|cde", []string{"abc|cde"}, false},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tList}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
		tests = []test{
			{"abc|def", []string{"abc", "def"}, false},
			{"abc,def", []string{"abc,def"}, false},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				originalListSeparator := strings.Clone(ListSeparator)
				ListSeparator = "|"
				v := Value{Mutagensis: tList}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
				ListSeparator = originalListSeparator
			})
		}
	})
	t.Run("Map", func(t *testing.T) {
		testNo := 0
		tests := []test{
			{t.Name(), t.Name(), true},
			{"", map[string]string{}, false},
			{"name=yeshua,age=33", map[string]string{"name": "yeshua", "age": "33"}, false},
			{"invalid=single=format=set", map[string]string{"invalid": "single=format=set"}, false},
			{"nothing-something", map[string]string{}, true},
		}
		for _, test := range tests {
			testNo++
			t.Run(fmt.Sprintf("TEST-%d", testNo), func(t *testing.T) {
				v := Value{Mutagensis: tMap}
				if test.WantErr {
					assert.Error(t, v.Set(test.In))
					assert.NotEqual(t, test.Out, v.Value)
				} else {
					assert.NoError(t, v.Set(test.In))
					assert.Equal(t, test.Out, v.Value)
				}
			})
		}
	})
}
