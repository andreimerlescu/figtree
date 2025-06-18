package figtree

type Value struct {
	Value      interface{}
	Mutagensis Mutagenesis
}

func (v *Value) Set(in string) error {
	v.Value = in
	return nil
}

func (v *Value) Flesh() Flesh {
	return &figFlesh{v.Value}
}

func (v *Value) String() string {
	f := figFlesh{v.Value}
	return f.ToString()
}
