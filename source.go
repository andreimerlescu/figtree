package figtree

type SourceKind int

const (
	SourceKindUnknown SourceKind = iota
	SourceKindEnv
	SourceKindFile
	SourceKindFlag
	SourceKindFlagEnv
)

type SourceConfig interface {
	Fetch() (string, error)
	Kind() SourceKind
}

func (tree *figTree) WithSource(source SourceConfig) error {
	return nil
}

func (tree *figTree) Source(name string) error {
	tree.mu.RLock()
	source, exists := tree.sources[name]
	tree.mu.RUnlock()
	if !exists {
		return ErrSourceNotFound{}
	}
	result, err := source.Fetch()
	if err != nil {
		return err
	}
	tree.StoreString(name, result)

	return nil
}

type ErrSourceNotFound struct{}

func (e ErrSourceNotFound) Error() string {
	return "source not found"
}
