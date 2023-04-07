package backgroundserver

import (
	"net/url"
	"os"
	"os/user"
)

const (
	StateNew = 1 + iota
	StateRunning
	StateStopped
	StateFailed
)

// BackgroundServer is a helper to run a real cockroach node.
type BackgroundServer interface {
	// Start starts the server.
	Start() error
	// Stop stops the server and cleans up any associated resources.
	Stop()
	// Stdout returns the entire contents of the process' stdout.
	Stdout() string
	// Stderr returns the entire contents of the process' stderr.
	Stderr() string
	// ConnURL returns the connection URL to this server.
	ConnURL() *url.URL
	// WaitForInit retries until a SQL connection is successfully established to this server.
	WaitForInit() error
}

type LogWriter interface {
	Write(p []byte) (n int, err error)
	String() string
	Len() int64
	Close() error
}

type FileLogWriter struct {
	filename string
	file     *os.File
}

func NewFileLogWriter(file string) (*FileLogWriter, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	return &FileLogWriter{
		filename: file,
		file:     f,
	}, nil
}

func (w FileLogWriter) Close() error {
	return w.file.Close()
}

func (w FileLogWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w FileLogWriter) String() string {
	b, err := os.ReadFile(w.filename)
	if err == nil {
		return string(b)
	}
	return ""
}

func (w FileLogWriter) Len() int64 {
	s, err := os.Stat(w.filename)
	if err == nil {
		return s.Size()
	}
	return 0
}

func DefaultEnv() map[string]string {
	vars := map[string]string{}
	u, err := user.Current()
	if err == nil {
		if _, ok := vars["USER"]; !ok {
			vars["USER"] = u.Username
		}
		if _, ok := vars["UID"]; !ok {
			vars["UID"] = u.Uid
		}
		if _, ok := vars["GID"]; !ok {
			vars["GID"] = u.Gid
		}
		if _, ok := vars["HOME"]; !ok {
			vars["HOME"] = u.HomeDir
		}
	}
	if _, ok := vars["PATH"]; !ok {
		vars["PATH"] = os.Getenv("PATH")
	}
	return vars
}
