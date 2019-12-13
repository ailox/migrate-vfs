package vfs

import (
	"testing"

	"github.com/ailox/migrate-vfs/testdata"
	st "github.com/golang-migrate/migrate/source/testing"
)

func Test(t *testing.T) {
	// wrap assets into Resource first

	d, err := WithInstance(testdata.Assets, "/migrations")
	if err != nil {
		t.Fatal(err)
	}
	st.Test(t, d)
}

func TestWithInstance(t *testing.T) {
	// wrap assets into Resource
	_, err := WithInstance(testdata.Assets, "/migrations")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpen(t *testing.T) {
	b := &HttpFileSystem{}
	_, err := b.Open("")
	if err == nil {
		t.Fatal("expected err, because it's not implemented yet")
	}
}
