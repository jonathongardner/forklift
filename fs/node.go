package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jonathongardner/forklift/filetype"
)

// -------------------Reference---------------------
// everything unique to a file (i.e not mode or name)
type reference struct {
	id       string
	entry    *entry
	children map[string]*node
}

func (r *reference) storagePath(storageDir string) string {
	return filepath.Join(storageDir, r.id)
}

func (r *reference) remove(storageDir string) error {
	return os.Remove(r.storagePath(storageDir))
}

func (r *reference) create(storageDir string) (*os.File, error) {
	return os.Create(r.storagePath(storageDir))
}
func (r *reference) open(storageDir string) (*os.File, error) {
	return os.Open(r.storagePath(storageDir))
}

// -------------------Node---------------------
// unique to the file like mode path/name, etc
type node struct {
	name    string
	mode    os.FileMode
	modTime time.Time // TODO: handle
	ref     *reference
	root    bool
}

func newNodeWithEntry(name string, mode os.FileMode, entry *entry) *node {
	extracted := &atomic.Bool{}
	extracted.Store(false)
	reference := &reference{id: uuid.New().String(), entry: entry, children: make(map[string]*node)}
	return &node{name: name, mode: mode, ref: reference}
}

func newNode(name string, mode os.FileMode) *node {
	return newNodeWithEntry(name, mode, newEntry())
}

func newDirNode(name string, mode os.FileMode) *node {
	return newNodeWithEntry(name, mode, newDirEntry())
}

func (n *node) errorID() error {
	return fmt.Errorf("id: %v, name: %v, type: %v,", n.ref.id, n.name, n.ref.entry.Type)
}

func (n *node) create(storageDir, name string, perm os.FileMode) (*node, *os.File, error) {
	// NOTE:orphan this could orphin some references, might want to clean up if reference is not needed
	newNode := newNode(name, perm)
	n.ref.children[name] = newNode
	file, err := newNode.ref.create(storageDir)

	return newNode, file, err
}
func (n *node) open(storageDir string) (*os.File, error) {
	return n.ref.open(storageDir)
}

func (n *node) mkdirP(paths []string, perm os.FileMode) (*node, error) {
	if len(paths) == 0 {
		return n, nil
	}

	firstPath, paths := paths[0], paths[1:]

	child, ok := n.ref.children[firstPath]
	if !ok {
		n.ref.children[firstPath] = newDirNode(firstPath, perm)
		child = n.ref.children[firstPath]
	} else if child.ref.entry.Type != filetype.Dir {
		// NOTE:orphan some references, only issue is it could be a possible large file
		// that isnt accessible so not needed and we could delete to free space
		n.ref.children[firstPath] = newDirNode(firstPath, perm)
		child = n.ref.children[firstPath]
	}
	return child.mkdirP(paths, perm)
}

func (n *node) walkRecursive(path string, skipParent bool, callback func(string, FileInfo) error) error {
	if !skipParent {
		err := callback(path, NewFileInfo(n))
		if err == ErrDontWalk {
			return nil
		}
		if err != nil {
			return err
		}
	}
	if n.ref.children == nil {
		return nil
	}

	for name, child := range n.ref.children {
		if err := child.walkRecursive(filepath.Join(path, name), false, callback); err != nil {
			return err
		}
	}

	return nil
}

// ------------------node------------------
