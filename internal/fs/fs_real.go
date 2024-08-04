package fs

import "os"

type realFS struct {
	fp goFilepath
}

func RealFS() (FS, error) {
	var fp goFilepath
	if cwd, err := os.Getwd(); err == nil {
		fp.cwd = cwd
	} else {
		fp.cwd = "/"
	}

	return &realFS{
		fp: fp,
	}, nil
}

func (fs *realFS) Cwd() string {
	return fs.fp.cwd
}
