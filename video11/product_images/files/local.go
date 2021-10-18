package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

//Local is an implementation of the storage interface which works with
// local disk of the current machine
type Local struct {
	maxFileSize int //maximum number of bytes for files
	basePath    string
}

// NewLocal creates a local filesystem with the given basepath
// basePath or "base" in this case is the base directory to save files
// to maxfilesize number of bytes that a file can be
func NewLocal(base string, size int) (*Local, error) {
	//we didnt use size int here, just a basic demonstration
	p, err := filepath.Abs(base) //filepath is an inbuilt package
	if err != nil {
		return nil, err
	}
	return &Local{basePath: p}, nil
}

// Save the contents of the writer to the given path
// implementing storage interface
// doubt, what's the benefit of implementing storage interface
func (l *Local) Save(path string, content io.Reader) error {
	//get the full path for the file
	fp := l.fullPath(path)

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.Mkdir(d, os.ModePerm) //not advisable to use 0o777 code with is modeperm because it can help inject malicious code into the server
	if err != nil {
		return xerrors.Errorf("unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err1 := os.Remove(fp)
		if err1 != nil {
			return xerrors.Errorf("unable to delete file: %w", err1)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return xerrors.Errorf("unable to get file info: %w", err)
	}

	//create a new file path
	f, err := os.Create(fp)
	if err != nil {
		xerrors.Errorf("Unable to create file: %w", err)
	}

	// write the contents to the new file
	// ensure that we are not writing greater than max bytes
	_, err = io.Copy(f, content)
	if err != nil {
		xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil

}

// get the file at the given path and return a Reader
// the calling function is responsible for closing the reader
func (l *Local) Get(path string) (*os.File, error) {
	// get the full path for the file
	fp := l.fullPath(path)

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("unable to open file:: %w", err)
	}

	return f, nil
}

//returns absolute path
func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}
