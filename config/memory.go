// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package config

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
)

// memoryStore implements the Store interface. It is meant primarily for testing.
type memoryStore struct {
	emitter
	commonStore

	allowEnvironmentOverrides bool
	validate                  bool
	files                     map[string][]byte
}

// MemoryStoreOptions makes configuration of the memory store explicit.
type MemoryStoreOptions struct {
	IgnoreEnvironmentOverrides bool
	SkipValidation             bool
}

// NewMemoryStore creates a new memoryStore instance.
func NewMemoryStore(options *MemoryStoreOptions) (*memoryStore, error) {
	ms := &memoryStore{
		allowEnvironmentOverrides: !options.IgnoreEnvironmentOverrides,
		validate:                  !options.SkipValidation,
		files:                     make(map[string][]byte),
	}

	if err := ms.Load(); err != nil {
		return nil, err
	}

	return ms, nil
}

// Set replaces the current configuration in its entirety.
func (ms *memoryStore) Set(newCfg *model.Config) (*model.Config, error) {
	validate := ms.commonStore.validate
	if !ms.validate {
		validate = nil
	}

	return ms.commonStore.set(newCfg, validate)
}

// Load applies environment overrides to the default config as if a re-load had occurred.
func (ms *memoryStore) Load() (err error) {
	defaultCfg := &model.Config{}
	defaultCfg.SetDefaults()

	var cfgBytes []byte
	cfgBytes, err = marshalConfig(defaultCfg)
	if err != nil {
		return errors.Wrap(err, "failed to serialize config")
	}

	f := ioutil.NopCloser(bytes.NewReader(cfgBytes))

	return ms.commonStore.load(f, false, nil, nil)
}

// Save does nothing, as there is no backing store.
func (ms *memoryStore) Save() error {
	return nil
}

// GetFile fetches the contents of a previously persisted configuration file.
func (ms *memoryStore) GetFile(name string) ([]byte, error) {
	ms.configLock.RLock()
	defer ms.configLock.RUnlock()

	data, ok := ms.files[name]
	if !ok {
		return nil, fmt.Errorf("file %s not stored", name)
	}

	return data, nil
}

// SetFile sets or replaces the contents of a configuration file.
func (ms *memoryStore) SetFile(name string, data []byte) error {
	ms.configLock.Lock()
	defer ms.configLock.Unlock()

	ms.files[name] = data

	return nil
}

// HasFile returns true if the given file was previously persisted.
func (ms *memoryStore) HasFile(name string) (bool, error) {
	ms.configLock.RLock()
	defer ms.configLock.RUnlock()

	_, ok := ms.files[name]
	return ok, nil
}

// RemoveFile remoevs a previously persisted configuration file.
func (ms *memoryStore) RemoveFile(name string) error {
	ms.configLock.Lock()
	defer ms.configLock.Unlock()

	delete(ms.files, name)

	return nil
}

// String returns a hard-coded description, as there is no backing store.
func (ms *memoryStore) String() string {
	return "mock://"
}

// Close does nothing for a mock store.
func (ms *memoryStore) Close() error {
	return nil
}
