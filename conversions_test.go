package figtree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFigFlesh_ToString(t *testing.T) {
	name := t.Name()
	flesh := NewFlesh(name).AsIs()
	assert.Equal(t, name, flesh)
}

func Test_toBool(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "True bool",
			args:    args{value: true},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name:    "False bool",
			args:    args{value: false},
			want:    false,
			wantErr: assert.NoError,
		},
		{
			name:    "True string",
			args:    args{value: "true"},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name:    "False string",
			args:    args{value: "false"},
			want:    false,
			wantErr: assert.NoError,
		},
		{
			name:    "Invalid string",
			args:    args{value: "notabool"},
			want:    false,
			wantErr: assert.Error,
		},
		{
			name:    "Int should fail",
			args:    args{value: 1},
			want:    false,
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    false,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toBool(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toBool(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toBool(%v)", tt.args.value)
		})
	}
}

func Test_toFloat64(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Float64 value",
			args:    args{value: 3.14},
			want:    3.14,
			wantErr: assert.NoError,
		},
		{
			name:    "String float",
			args:    args{value: "2.718"},
			want:    2.718,
			wantErr: assert.NoError,
		},
		{
			name:    "String int as float",
			args:    args{value: "42"},
			want:    42.0,
			wantErr: assert.NoError,
		},
		{
			name:    "Invalid string",
			args:    args{value: "notafloat"},
			want:    0.0,
			wantErr: assert.Error,
		},
		{
			name:    "Bool should fail",
			args:    args{value: true},
			want:    0.0,
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    0.0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toFloat64(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toFloat64(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toFloat64(%v)", tt.args.value)
		})
	}
}

func Test_toInt(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Float64 as int",
			args:    args{value: 42.0},
			want:    42,
			wantErr: assert.NoError,
		},
		{
			name:    "String int",
			args:    args{value: "123"},
			want:    123,
			wantErr: assert.NoError,
		},
		{
			name:    "String float truncated",
			args:    args{value: "45.6"},
			want:    45,
			wantErr: assert.NoError,
		},
		{
			name:    "Invalid string",
			args:    args{value: "notanint"},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:    "Bool should fail",
			args:    args{value: true},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toInt(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toInt(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toInt(%v)", tt.args.value)
		})
	}
}

func Test_toInt64(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Float64 as int64",
			args:    args{value: 42.0},
			want:    42,
			wantErr: assert.NoError,
		},
		{
			name:    "String int64",
			args:    args{value: "1234567890"},
			want:    1234567890,
			wantErr: assert.NoError,
		},
		{
			name:    "String float truncated",
			args:    args{value: "45.6"},
			want:    45,
			wantErr: assert.NoError,
		},
		{
			name:    "Invalid string",
			args:    args{value: "notanint"},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:    "Bool should fail",
			args:    args{value: true},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toInt64(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toInt64(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toInt64(%v)", tt.args.value)
		})
	}
}

func Test_toString(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Basic string",
			args:    args{value: "hello"},
			want:    "hello",
			wantErr: assert.NoError,
		},
		{
			name:    "Float64 to string",
			args:    args{value: 3.14},
			want:    "3.14",
			wantErr: assert.NoError,
		},
		{
			name:    "Bool true to string",
			args:    args{value: true},
			want:    "true",
			wantErr: assert.NoError,
		},
		{
			name:    "Bool false to string",
			args:    args{value: false},
			want:    "false",
			wantErr: assert.NoError,
		},
		{
			name:    "Int should fail",
			args:    args{value: 42},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "list can print as string",
			args:    args{value: []string{"yah", "i am", "yahuah"}},
			want:    "yah,i am,yahuah",
			wantErr: assert.NoError,
		},
		{
			name:    "map can print as string",
			args:    args{value: map[string]string{"name": "yahuah"}},
			want:    "name=yahuah",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toString(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toString(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toString(%v)", tt.args.value)
		})
	}
}

func Test_toStringMap(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Map of strings",
			args:    args{value: map[string]interface{}{"key1": "value1", "key2": "value2"}},
			want:    map[string]string{"key1": "value1", "key2": "value2"},
			wantErr: assert.NoError,
		},
		{
			name:    "String key-value pairs",
			args:    args{value: "key1=value1,key2=value2"},
			want:    map[string]string{"key1": "value1", "key2": "value2"},
			wantErr: assert.NoError,
		},
		{
			name:    "Empty string",
			args:    args{value: ""},
			want:    map[string]string{},
			wantErr: assert.NoError,
		},
		{
			name:    "Invalid string format",
			args:    args{value: "key1value1,key2=value2"},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "Map with non-string value",
			args:    args{value: map[string]interface{}{"key": 42}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "Int should fail",
			args:    args{value: 42},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toStringMap(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toStringMap(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toStringMap(%v)", tt.args.value)
		})
	}
}

func Test_toStringSlice(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Slice of interfaces",
			args:    args{value: []interface{}{"a", "b", "c"}},
			want:    []string{"a", "b", "c"},
			wantErr: assert.NoError,
		},
		{
			name:    "Comma-separated string",
			args:    args{value: "x,y,z"},
			want:    []string{"x", "y", "z"},
			wantErr: assert.NoError,
		},
		{
			name:    "Empty string",
			args:    args{value: ""},
			want:    []string{},
			wantErr: assert.NoError,
		},
		{
			name:    "Slice with non-string",
			args:    args{value: []interface{}{"a", 42, "c"}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "Int should fail",
			args:    args{value: 42},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "Nil should fail",
			args:    args{value: nil},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toStringSlice(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("toStringSlice(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toStringSlice(%v)", tt.args.value)
		})
	}
}
