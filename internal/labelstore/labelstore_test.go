package labelstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/labelstore"
)

func tmpStore(t *testing.T) *labelstore.Store {
	t.Helper()
	dir := t.TempDir()
	return labelstore.New(filepath.Join(dir, "labels.json"))
}

func TestSetAndGet(t *testing.T) {
	s := tmpStore(t)
	k := labelstore.Key{Proto: "tcp", Port: 80}
	s.Set(k, labelstore.Label{Name: "http", Comment: "web"})
	l, ok := s.Get(k)
	if !ok {
		t.Fatal("expected label to exist")
	}
	if l.Name != "http" {
		t.Fatalf("got name %q, want http", l.Name)
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	s := tmpStore(t)
	_, ok := s.Get(labelstore.Key{Proto: "tcp", Port: 9999})
	if ok {
		t.Fatal("expected no label")
	}
}

func TestDelete(t *testing.T) {
	s := tmpStore(t)
	k := labelstore.Key{Proto: "udp", Port: 53}
	s.Set(k, labelstore.Label{Name: "dns"})
	s.Delete(k)
	_, ok := s.Get(k)
	if ok {
		t.Fatal("expected label to be deleted")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")

	s1 := labelstore.New(path)
	s1.Set(labelstore.Key{Proto: "tcp", Port: 443}, labelstore.Label{Name: "https"})
	if err := s1.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	s2 := labelstore.New(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	l, ok := s2.Get(labelstore.Key{Proto: "tcp", Port: 443})
	if !ok || l.Name != "https" {
		t.Fatalf("expected https label, got %+v ok=%v", l, ok)
	}
}

func TestLoadMissingFileReturnsError(t *testing.T) {
	s := labelstore.New("/nonexistent/path/labels.json")
	if err := s.Load(); !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}
