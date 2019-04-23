package vfs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/golang-migrate/migrate/source"
)

func init() {
	source.Register("httpfs", &HttpFileSystem{})
}

// HttpFileSystem satisfies the driver interface
type HttpFileSystem struct {
	path       string
	fs         http.FileSystem
	migrations *source.Migrations
}

func WithInstance(i http.FileSystem, p string) (source.Driver, error) {
	if len(p) == 0 {
		// default to root directory if no path
		p = "/"

	} else if p[0:1] == "." || p[0:1] != "/" {
		// make path absolute if relative
		p = "/" + p
	}

	rootDir, err := i.Open(p)
	if err != nil {
		return nil, err
	}

	files, err := rootDir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	nf := &HttpFileSystem{
		fs:         i,
		path:       p,
		migrations: source.NewMigrations(),
	}

	for _, fi := range files {
		if !fi.IsDir() {
			m, err := source.DefaultParse(fi.Name())
			if err != nil {
				continue // ignore files that we can't parse
			}
			if !nf.migrations.Append(m) {
				return nil, fmt.Errorf("unable to parse file %v", fi.Name())
			}
		}
	}

	return nf, nil
}

func (fs *HttpFileSystem) Close() error {
	return fs.Close()
}

func (fs *HttpFileSystem) First() (version uint, err error) {
	if v, ok := fs.migrations.First(); !ok {
		return 0, &os.PathError{"first", fs.path, os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (fs *HttpFileSystem) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := fs.migrations.Prev(version); !ok {
		return 0, &os.PathError{fmt.Sprintf("prev for version %v", version), fs.path, os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (fs *HttpFileSystem) Next(version uint) (nextVersion uint, err error) {
	if v, ok := fs.migrations.Next(version); !ok {
		return 0, &os.PathError{fmt.Sprintf("next for version %v", version), fs.path, os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (fs *HttpFileSystem) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := fs.migrations.Up(version); ok {
		r, err := fs.fs.Open(path.Join(fs.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{fmt.Sprintf("read version %v", version), fs.path, os.ErrNotExist}
}

func (fs *HttpFileSystem) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := fs.migrations.Down(version); ok {
		r, err := fs.fs.Open(path.Join(fs.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{fmt.Sprintf("read version %v", version), fs.path, os.ErrNotExist}
}

func (fs *HttpFileSystem) Open(url string) (source.Driver, error) {
	return nil, errors.New("Open() is not implemented with this backend, please use WithInstance()")
}
