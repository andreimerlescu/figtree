package figtree

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestTree_checkAndSetFromEnv(t *testing.T) {
	const k, u = "workers-check-and-set-from-env", "usage"

	// create a new fig tree
	var figs *Tree
	figs = &Tree{
		harvest:     1,
		figs:        make(map[string]*Fig),
		tracking:    false,
		withered:    make(map[string]Fig),
		flagSet:     flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
		mu:          sync.RWMutex{},
		mutationsCh: make(chan Mutation, 1),
		filterTests: true,
	}

	// assign an int to k
	figs.NewInt(k, 10, u)
	assert.Nil(t, figs.Parse())

	// verify assignment
	assert.Equal(t, 10, *figs.Int(k))

	// clear ENV
	os.Clearenv()

	// set env for k to 17
	err := os.Setenv(k, "17")
	assert.NoError(t, err)

	figs.checkAndSetFromEnv(k)
	assert.Equal(t, 17, *figs.Int(k))
}

func TestTree_setValue(t *testing.T) {
	type fields struct {
		ConfigFilePath string
		figs           map[string]*Fig
		withered       map[string]Fig
		mu             sync.RWMutex
		tracking       bool
		mutationsCh    chan Mutation
	}
	type args struct {
		flagVal interface{}
		value   interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue interface{}
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "Set int value",
			fields: fields{
				figs:        make(map[string]*Fig),
				withered:    make(map[string]Fig),
				mutationsCh: make(chan Mutation, 1),
			},
			args: args{
				flagVal: new(int),
				value:   42,
			},
			wantValue: 42,
			wantErr:   assert.NoError,
		},
		{
			name: "Set string value",
			fields: fields{
				figs:        make(map[string]*Fig),
				withered:    make(map[string]Fig),
				mutationsCh: make(chan Mutation, 1),
			},
			args: args{
				flagVal: new(string),
				value:   "hello",
			},
			wantValue: "hello",
			wantErr:   assert.NoError,
		},
		{
			name: "Invalid type",
			fields: fields{
				figs:        make(map[string]*Fig),
				withered:    make(map[string]Fig),
				mutationsCh: make(chan Mutation, 1),
			},
			args: args{
				flagVal: new(float32), // Unsupported type
				value:   3.14,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fig := &Tree{
				ConfigFilePath: tt.fields.ConfigFilePath,
				figs:           tt.fields.figs,
				withered:       tt.fields.withered,
				mu:             tt.fields.mu,
				tracking:       tt.fields.tracking,
				mutationsCh:    tt.fields.mutationsCh,
				flagSet:        flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
				filterTests:    true,
			}
			err := fig.setValue(tt.args.flagVal, tt.args.value)
			if tt.wantErr(t, err, fmt.Sprintf("setValue(%v, %v)", tt.args.flagVal, tt.args.value)) {
				return
			}

			// Verify value only if no error is expected
			if tt.wantErr(t, nil) {
				switch v := tt.args.flagVal.(type) {
				case *int:
					assert.Equal(t, tt.wantValue, *v, "setValue should set int value")
				case *string:
					assert.Equal(t, tt.wantValue, *v, "setValue should set string value")
				}
			}
		})
	}
}

func TestTree_setValuesFromMap(t *testing.T) {
	tree := &Tree{
		figs:        make(map[string]*Fig),
		withered:    make(map[string]Fig),
		mu:          sync.RWMutex{},
		tracking:    false,
		mutationsCh: make(chan Mutation, 1),
		flagSet:     flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
		filterTests: true,
	}
	m := map[string]interface{}{
		"name": "yahuah",
		"age":  33,
		"sex":  "male",
	}
	err := tree.setValuesFromMap(m)
	assert.NoError(t, err)
	assert.Equal(t, "yahuah", m["name"])
	assert.Equal(t, "male", m["sex"])
}
