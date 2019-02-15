// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package config

import (
	"bytes"
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
}

// NewMemoryStore creates a new memoryStore instance.
func NewMemoryStore(allowEnvironmentOverrides bool, validate bool) (*memoryStore, error) {
	ms := &memoryStore{
		allowEnvironmentOverrides: allowEnvironmentOverrides,
		validate:                  validate,
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

// String returns a hard-coded description, as there is no backing store.
func (ms *memoryStore) String() string {
	return "mock://"
}

// Close does nothing for a mock store.
func (ms *memoryStore) Close() error {
	return nil
}
