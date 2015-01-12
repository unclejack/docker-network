package state

import (
	"fmt"
	"strings"

	"github.com/docker/libpack"
)

// GitState uses docker/libpack and satisfies the State interface.
type GitState struct {
	db    *libpack.DB
	scope string
}

// GitStateFromFolder returns a ready-to-use GitState.
// The same folder can be used for different stores identified by storeName
// storeName can not contain slashes, as it could be mistaken for a subdirectory.
func GitStateFromFolder(folder, storeName string) (*GitState, error) {
	if strings.Contains(storeName, "/") {
		return nil, fmt.Errorf("Slashes are not allowed in storeName: %q", storeName)
	}
	r, err := libpack.Init(folder, true)
	db, err := r.DB("refs/heads/" + storeName)
	if err != nil {
		return nil, err
	}
	return &GitState{db: db}, nil
}

// Close releases resources for the underlying db.
// Subsequent calls to any methods on the same GitState will result in a panic.
// To recover the same state, call GitStateFromFolder again with the same arguments.
func (s GitState) Close() {
	s.db.Repo().Free()
	s.db = nil
}

// Get returns the value associated with `key`.
func (s GitState) Get(key string) (value string, err error) {
	t, err := s.db.Query().Scope(s.scope).Run()
	if err != nil {
		return "", err
	}

	return t.Get(key)
}

// List returns a list of keys directly under dir.
// Does not walk the tree recursively.
//
// Example: /
//		foo
//			bar
//		baz
// List("/") -> ["foo", "bar"]
func (s GitState) List(dir string) ([]string, error) {
	t, err := s.db.Query().Scope(s.scope).Run()
	if err != nil {
		return nil, err
	}

	return t.List(dir)
}

// Set sets the key `key`, to a value `value`.
// It automatically overrides the existing value if any.
func (s GitState) Set(key, value string) error {
	_, err := s.db.Query().Scope(s.scope).Set(key, value).Commit(s.db).Run()
	return err
}

// Remove deletes the value associated with `key`.
func (s GitState) Remove(key string) error {
	_, err := s.db.Query().Scope(s.scope).Delete(key).Commit(s.db).Run()
	return err
}

// Mkdir creates the directory `dir`.
func (s GitState) Mkdir(dir string) error {
	_, err := s.db.Query().Scope(s.scope).Mkdir(dir).Commit(s.db).Run()
	return err
}

// FIXME: should we stop returning error?
func (s GitState) Scope(key string) (State, error) {
	// FIXME: should we split the key on the path separator?
	return GitState{scope: key, db: s.db}, nil
}

func (s GitState) Diff(other Tree) (added, removed Tree)            { return nil, nil }
func (s GitState) Walk(func(key string, entry Value)) error         { return nil }
func (s GitState) Add(key string, overlay Tree) (Tree, error)       { return nil, nil }
func (s GitState) Subtract(key string, whiteout Tree) (Tree, error) { return nil, nil }
func (s GitState) Pipeline() Pipeline                               { return nil }
