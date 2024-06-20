package fs

import (
	"os"
	"slices"
	"sort"
	"testing"
)

const fooFile = "testdata/foo"
const fooSha512 = "0f5623276549769a63c79ca20fc573518685819fe82b39f43a3e7cf709c8baa16524daa95e006e81f7267700a88adee8a6209201be960a10c81c35ff3547e3b7"

// const barFile = "testdata/bar"
// const bazFile = "testdata/baz"

func tmpDir(fnc func(tmp string)) error {
	dname, err := os.MkdirTemp("", "virtual-testing")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dname)
	fnc(dname)
	return nil
}

type fileinfoTest struct {
	mode   os.FileMode
	sha512 string
	ftype  string
}

func assertFiles(t *testing.T, expected map[string]fileinfoTest, v *Virtual, message string) {
	v.Walk("/", func(path string, fi FileInfo) error {
		fit, ok := expected[path]
		if !ok {
			t.Fatalf("%v %v unknown path", message, path)
		}
		delete(expected, path)

		if fit.mode != fi.mode {
			t.Fatalf("%v %v mode doesnt match %v != %v", message, path, fit.mode, fi.mode)
		}

		if fit.sha512 != fi.Sha512 {
			t.Fatalf("%v %v sha512 doesnt match %v != %v", message, path, fit.sha512, fi.Sha512)
		}

		if fit.ftype != fi.Type.Mimetype {
			t.Fatalf("%v %v filetype doesnt match %v != %v", message, path, fit.ftype, fi.Type.Mimetype)
		}

		return nil
	})
	if len(expected) != 0 {
		t.Fatalf("Expected %v to be empty", expected)
	}
}

func createFile(t *testing.T, v *Virtual, path string, perm os.FileMode, content string) {
	file, err := v.Create(path, perm)
	if err != nil {
		t.Fatalf("Failed to create virtual file %v: %v", path, err)
	}
	_, err = file.Write([]byte(content))
	if err != nil {
		t.Fatalf("Failed to write to virtual file %v: %v", path, err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close virtual file %v: %v", path, err)
	}
}
func assertFileCount(t *testing.T, tmp string, cnt int, msg string) {
	d, err := os.ReadDir(tmp)
	if err != nil {
		t.Fatalf("%v Failed to get file count: %v", msg, err)
	}
	if len(d) != cnt {
		t.Fatalf("%v Expected %v files, got %v", msg, cnt, len(d))

	}
}
func assertPaths(t *testing.T, expected []string, v *Virtual, message string) {
	walked := []string{}
	v.Walk("/", func(path string, fi FileInfo) error {
		walked = append(walked, path)
		return nil
	})
	sort.Strings(walked)
	if !slices.Equal(walked, expected) {
		t.Fatalf("%v Unexpected walked: %v", message, walked)
	}
}

func TestVirtual(t *testing.T) {
	err := tmpDir(func(tmp string) {
		v, err := NewVirtual(tmp, fooFile)
		if err != nil {
			t.Fatalf("Failed to create virtual function: %v", err)
		}
		expected := map[string]fileinfoTest{"/": {0664, fooSha512, "application/octet-stream"}}
		assertFiles(t, expected, v, "initial")

		v.MkdirP("/foo1/foo2", 0755)
		expected = map[string]fileinfoTest{
			"/":          {0664, fooSha512, "application/octet-stream"},
			"/foo1":      {0755, "", "directory/directory"},
			"/foo1/foo2": {0755, "", "directory/directory"},
		}
		assertFiles(t, expected, v, "after creating foo1/foo2")

		v.MkdirP("/foo1/foo2/foo3/foo4", 0700)
		expected = map[string]fileinfoTest{
			"/":                    {0664, fooSha512, "application/octet-stream"},
			"/foo1":                {0755, "", "directory/directory"},
			"/foo1/foo2":           {0755, "", "directory/directory"},
			"/foo1/foo2/foo3":      {0700, "", "directory/directory"},
			"/foo1/foo2/foo3/foo4": {0700, "", "directory/directory"},
		}
		assertFiles(t, expected, v, "after creating foo1/foo2/foo3/foo4")
		assertFileCount(t, tmp, 1, "after creating foo1/foo2/foo3/foo4")

		createFile(t, v, "/foo1/foo2/foo3/bar", 0655, "Hello, World!")
		expected = map[string]fileinfoTest{
			"/":                    {0664, fooSha512, "application/octet-stream"},
			"/foo1":                {0755, "", "directory/directory"},
			"/foo1/foo2":           {0755, "", "directory/directory"},
			"/foo1/foo2/foo3":      {0700, "", "directory/directory"},
			"/foo1/foo2/foo3/foo4": {0700, "", "directory/directory"},
			"/foo1/foo2/foo3/bar":  {0655, "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387", "text/plain; charset=utf-8"},
		}
		assertFiles(t, expected, v, "after creating foo1/foo2/foo3/bar")
		assertFileCount(t, tmp, 2, "after creating foo1/foo2/foo3/bar")
	})
	if err != nil {
		t.Fatalf("Failed to create tmp dir for testing: %v", err)
	}
}

func TestVirtualUsesReferencesForSameFile(t *testing.T) {
	err := tmpDir(func(tmp string) {
		v, err := NewVirtual(tmp, fooFile)
		if err != nil {
			t.Fatalf("Failed to create virtual function: %v", err)
		}

		createFile(t, v, "/bar", 0655, "Hello, World!")
		createFile(t, v, "/baz", 0600, "Hello, World!")
		// should get added to both bar and baz since they are the same file
		newV, err := v.VirtualFrom("/baz")
		if err != nil {
			t.Fatalf("Failed to create virtual from baz: %v", err)
		}
		createFile(t, newV, "/moreFoo", 0100, "Hello, Foo!")

		expected := map[string]fileinfoTest{
			"/":            {0664, fooSha512, "application/octet-stream"},
			"/bar":         {0655, "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387", "text/plain; charset=utf-8"},
			"/bar/moreFoo": {0100, "9b617e0675ac2ede198cfacddf0b283d378a2cee8e72e551a1ae5400cdb9a46792556187e4d2fdbedece0f0021a6b1f74a6b460b62966ef68025abf75fb7df7a", "text/plain; charset=utf-8"},
			"/baz":         {0600, "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387", "text/plain; charset=utf-8"},
			"/baz/moreFoo": {0100, "9b617e0675ac2ede198cfacddf0b283d378a2cee8e72e551a1ae5400cdb9a46792556187e4d2fdbedece0f0021a6b1f74a6b460b62966ef68025abf75fb7df7a", "text/plain; charset=utf-8"},
		}
		assertFiles(t, expected, v, "comparing")
		assertFileCount(t, tmp, 3, "comparing")
	})
	if err != nil {
		t.Fatalf("Failed to create tmp dir for testing: %v", err)
	}
}

func TestVirtualOverwriteFileWithDir(t *testing.T) {
	err := tmpDir(func(tmp string) {
		v, err := NewVirtual(tmp, fooFile)
		if err != nil {
			t.Fatalf("Failed to create virtual function: %v", err)
		}

		createFile(t, v, "/bar", 0655, "Hello, World!")
		createFile(t, v, "/baz", 0600, "Hello, World!")
		createFile(t, v, "/bar/moreFoo", 0100, "Hello, Foo!")

		expected := map[string]fileinfoTest{
			"/":            {0664, fooSha512, "application/octet-stream"},
			"/bar":         {0100, "", "directory/directory"},
			"/bar/moreFoo": {0100, "9b617e0675ac2ede198cfacddf0b283d378a2cee8e72e551a1ae5400cdb9a46792556187e4d2fdbedece0f0021a6b1f74a6b460b62966ef68025abf75fb7df7a", "text/plain; charset=utf-8"},
			"/baz":         {0600, "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387", "text/plain; charset=utf-8"},
		}
		assertFiles(t, expected, v, "comparing")
		assertFileCount(t, tmp, 3, "comparing")
	})
	if err != nil {
		t.Fatalf("Failed to create tmp dir for testing: %v", err)
	}
}

func TestVirtualFrom(t *testing.T) {
	err := tmpDir(func(tmp string) {
		v, err := NewVirtual(tmp, fooFile)
		if err != nil {
			t.Fatalf("Failed to create virtual function: %v", err)
		}

		createFile(t, v, "/bar", 0655, "Hello, World!")

		newV, err := v.VirtualFrom("/bar")
		if err != nil {
			t.Fatalf("Failed to create virtual from baz: %v", err)
		}
		newV.MkdirP("/", 0700)  // shouldnt change anything cause root
		newV.MkdirP("./", 0700) // shouldnt change anything cause root
		newV.MkdirP(".", 0700)  // shouldnt change anything cause root
		createFile(t, newV, "/moreFoo", 0100, "Hello, Foo!")

		assertPaths(t, []string{"/", "/bar", "/bar/moreFoo"}, v, "comparing")

		expected := map[string]fileinfoTest{
			"/":            {0664, fooSha512, "application/octet-stream"},
			"/bar":         {0655, "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387", "text/plain; charset=utf-8"},
			"/bar/moreFoo": {0100, "9b617e0675ac2ede198cfacddf0b283d378a2cee8e72e551a1ae5400cdb9a46792556187e4d2fdbedece0f0021a6b1f74a6b460b62966ef68025abf75fb7df7a", "text/plain; charset=utf-8"},
		}
		assertFiles(t, expected, v, "comparing")
	})
	if err != nil {
		t.Fatalf("Failed to create tmp dir for testing: %v", err)
	}
}

func TestWalk(t *testing.T) {
	err := tmpDir(func(tmp string) {
		v, err := NewVirtual(tmp, fooFile)
		if err != nil {
			t.Fatalf("Failed to create virtual function: %v", err)
		}
		v.MkdirP("/foo1/foo2", 0755)
		createFile(t, v, "/foo1/foo2/foo3/bar", 0655, "Hello, World!")

		walked := []string{}
		v.Walk("/", func(path string, fi FileInfo) error {
			walked = append(walked, path)
			return nil
		})
		sort.Strings(walked)
		if !slices.Equal(walked, []string{"/", "/foo1", "/foo1/foo2", "/foo1/foo2/foo3", "/foo1/foo2/foo3/bar"}) {
			t.Fatalf("Unexpected walked: %v", walked)
		}

		walkedChild := []string{}
		v.WalkChildren("/", func(path string, fi FileInfo) error {
			walkedChild = append(walkedChild, path)
			return nil
		})
		sort.Strings(walkedChild)
		if !slices.Equal(walkedChild, []string{"/foo1", "/foo1/foo2", "/foo1/foo2/foo3", "/foo1/foo2/foo3/bar"}) {
			t.Fatalf("Unexpected walkedChild: %v", walkedChild)
		}

	})
	if err != nil {
		t.Fatalf("Failed to create tmp dir for testing: %v", err)
	}
}
