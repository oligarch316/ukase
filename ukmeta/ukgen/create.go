package ukgen

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type FS interface {
	Create(name string) (io.WriteCloser, error)
	// Close() error
}

// -----

type FSFunc func(name string) (io.WriteCloser, error)

func (ff FSFunc) Create(name string) (io.WriteCloser, error) { return ff(name) }

// -----
type fsError struct{ err error }

func (fe fsError) Create(string) (io.WriteCloser, error) { return nil, fe.err }

// -----

var FSStderr = FSWriter(os.Stderr)
var FSStdout = FSWriter(os.Stdout)

func FSWriter(w io.Writer) FS { return fsWriter{Writer: w} }

type fsWriter struct{ io.Writer }

func (fsWriter) Close() error                             { return nil }
func (fw fsWriter) Create(string) (io.WriteCloser, error) { return fw, nil }

// -----

type fsRoot struct{ root *os.Root }

func FSRoot(name string) FS {
	absName, err := filepath.Abs(name)
	if err != nil {
		return fsError{err: err}
	}

	root, err := os.OpenRoot(absName)
	if err != nil {
		return fsError{err: err}
	}

	return fsRoot{root: root}
}

func (fr fsRoot) Create(name string) (io.WriteCloser, error) {
	// Normalize target file path
	target, err := fr.normalizePath(name)
	if err != nil {
		return nil, err
	}

	// Create any necessary intermediate directories
	parent := filepath.Dir(target)
	if err = fr.ensureDirectory(parent); err != nil {
		return nil, err
	}

	// Create target file
	// TODO: Be better about filemode
	flag := os.O_RDWR | os.O_CREATE | os.O_TRUNC
	return fr.root.OpenFile(target, flag, 0666)
}

func (fr fsRoot) normalizePath(target string) (string, error) {
	if !filepath.IsAbs(target) {
		return target, nil
	}

	base := fr.root.Name()
	return filepath.Rel(base, target)
}

func (fr fsRoot) ensureDirectory(target string) error {
	if target == "." || target == "/" {
		return nil
	}

	switch exists, err := fr.checkDirectory(target); {
	case err != nil:
		return err
	case exists:
		return nil
	}

	parent := filepath.Dir(target)
	if err := fr.ensureDirectory(parent); err != nil {
		return err
	}

	// TODO: Be better about filemode
	return fr.root.Mkdir(target, 0755)
}

func (fr fsRoot) checkDirectory(target string) (bool, error) {
	info, err := fr.root.Stat(target)

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, errors.New("wrong kind of thing")
	}

	return true, nil
}

// =============================================================================
// Create
// =============================================================================

type Creator interface {
	Create(name string) (io.WriteCloser, error)
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type Create func(name string) (io.WriteCloser, error)

func (c Create) Create(name string) (io.WriteCloser, error) { return c(name) }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type nopCloser struct{ io.Writer }

func (nc nopCloser) Create(string) (io.WriteCloser, error) { return nc, nil }
func (nopCloser) Close() error                             { return nil }

var CreateStdout = CreateWriter(os.Stdout)
var CreateStderr = CreateWriter(os.Stderr)

func CreateWriter(w io.Writer) Creator { return nopCloser{Writer: w} }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

func CreateRoot(name string) Create {
	// TODO: Deal with root.Close() !!!

	root, err := os.OpenRoot(name)
	if err != nil {
		return Create(func(string) (io.WriteCloser, error) { return nil, err })
	}

	f := func(x string) (io.WriteCloser, error) { return root.Create(x) }
	return f
}
