package figtree

import (
	"fmt"
)

type SourceKind int

const (
	SourceAWSSecretsManager SourceKind = iota
	SourceAESString
	SourceGPGString
	SourceOnePasswordPath
	SourceVaultPath
	SourceKeeperPath
	SourceS3PlaintextPath
	SourceS3AESPath
	SourceS3GPGPath
)

type SourceConfig interface {
	Fetch() (string, error)
	Kind() SourceKind
}

func (s SourceKind) String() string {
	switch s {
	case SourceAWSSecretsManager:
		return "aws_secrets_manager"
	case SourceAESString:
		return "aes_string"
	case SourceGPGString:
		return "gpg_string"
	case SourceOnePasswordPath:
		return "one_password_path"
	case SourceVaultPath:
		return "vault_path"
	case SourceKeeperPath:
		return "keeper_path"
	case SourceS3PlaintextPath:
		return "s3_plaintext_path"
	case SourceS3AESPath:
		return "s3_aes_path"
	case SourceS3GPGPath:
		return "s3_gpg_path"
	default:
		return "unknown"
	}
}

// LoadAllFromSource will attempt to load all the sources into the Tree
func (tree *figTree) LoadAllFromSource() error {
	return tree.fetchFromSources()
}

// fetchFromSources retrieves values from configured sources
func (tree *figTree) fetchFromSources() error {
	tree.sourceLocker.Lock()
	tree.mu.RLock()
	defer tree.sourceLocker.Unlock()
	defer tree.mu.RUnlock()
	if tree.sources == nil {
		return nil
	}
	errs := make([]error, 0)
	for propName, sourceConfig := range tree.sources {
		value, err := sourceConfig.Fetch()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to fetch from source for %s: %w", propName, err))
			continue
		}

		fruit, exists := tree.figs[propName]
		if !exists || fruit == nil {
			continue
		}

		tree.Store(fruit.Mutagenesis, fruit.name, value)
	}

	return nil
}

func (tree *figTree) WithSource(name string, source SourceConfig) Plant {
	tree.sourceLocker.Lock()
	tree.mu.Lock()
	defer tree.sourceLocker.Unlock()
	defer tree.mu.Unlock()
	tree.sources[name] = source
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
