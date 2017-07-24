package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
)

// Store implements a K/V store for mapping distribution-related IDs
// to on-disk layer IDs and image IDs. The namespace identifies the type of
// mapping (i.e. "v1ids" or "artifacts"). MetadataStore is goroutine-safe.
type Store interface {
	// Get retrieves data by namespace and key.
	Get(namespace string, key string) ([]byte, error)
	// Set writes data indexed by namespace and key.
	Set(namespace, key string, value []byte) error
	// Delete removes data indexed by namespace and key.
	Delete(namespace, key string) error

	// List lists all keys by namespace
	List(namespace string) ([]string, error)
}

// FSMetadataStore uses the filesystem to associate metadata with layer and
// image IDs.
type FSMetadataStore struct {
	sync.RWMutex
	basePath string
}

// NewFSMetadataStore creates a new filesystem-based metadata store.
func NewFSMetadataStore(basePath string) (*FSMetadataStore, error) {
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return nil, err
	}
	return &FSMetadataStore{
		basePath: basePath,
	}, nil
}

func (store *FSMetadataStore) path(namespace, key string) string {
	return filepath.Join(store.basePath, namespace, key)
}

// Get retrieves data by namespace and key. The data is read from a file named
// after the key, stored in the namespace's directory.
func (store *FSMetadataStore) Get(namespace string, key string) ([]byte, error) {
	store.RLock()
	defer store.RUnlock()

	return ioutil.ReadFile(store.path(namespace, key))
}

// Set writes data indexed by namespace and key. The data is written to a file
// named after the key, stored in the namespace's directory.
func (store *FSMetadataStore) Set(namespace, key string, value []byte) error {
	store.Lock()
	defer store.Unlock()

	path := store.path(namespace, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutils.AtomicWriteFile(path, value, 0644)
}

// Delete removes data indexed by namespace and key. The data file named after
// the key, stored in the namespace's directory is deleted.
func (store *FSMetadataStore) Delete(namespace, key string) error {
	store.Lock()
	defer store.Unlock()

	path := store.path(namespace, key)
	return os.Remove(path)
}

// List lists all files in the sub directories in the path of namespace
func (store *FSMetadataStore) List(namespace string) ([]string, error) {
	store.RLock()
	defer store.RUnlock()

	result := make([]string, 0)
	var lastErr error = nil
	paths := store.path(namespace, "")
	total, err := ioutil.ReadDir(paths)
	if err != nil {
		return nil, err
	}
	for _, keys := range total {
		if keys.IsDir() {
			algor := keys.Name()
			sub := store.path(namespace, algor)
			digests, err := ioutil.ReadDir(sub)
			if err != nil {
				lastErr = err
				continue
			}
			for _, key := range digests {
				digest := algor + ":" + key.Name()
				result = append(result, digest)
			}
		}
	}
	return result, lastErr
}
