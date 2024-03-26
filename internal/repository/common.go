package repository

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/logger"
)

const (
	// dataFileMode is the default permissions set to the data files.
	dataFileMode = os.FileMode(0600)
	// dataDirMode is the default permissions set to the data directories.
	dataDirMode = os.FileMode(0700)
)

// HttpClient is the abstraction of a dependency of type http.Client, which allows mocking.
type HttpClient interface {
	// Do send an HTTP request and returns an HTTP response.
	Do(req *http.Request) (*http.Response, error)
}

// checkDataDir validates the given data directory.
func checkDataDir(dir string) error {
	if dir == "" {
		return ErrDirNameEmpty
	}
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return ErrIsNotDir
	}
	return nil
}

// createDataDir creates a new data directory with 0700 permissions.
// If the data directory exists, it gets validated.
func createDataDir(name string) error {
	if err := checkDataDir(name); err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return os.Mkdir(name, dataDirMode)
		default:
			return err
		}
	}
	return nil
}

// checkDataFile validates the given data file.
func checkDataFile(name string) error {
	if name == "" {
		return ErrFileNameEmpty
	}
	info, err := os.Stat(name)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return ErrFilePathIsDir
	}
	if info.Mode() != dataFileMode {
		return ErrFileModeInvalid
	}
	return nil
}

// createDataFile creates a new data file with 0600 permissions.
// If the data directory does not exist, creates a new one with 0700 permissions.
// If the data file exists with wrong permissions gets fixed to 0600 permissions.
func createDataFile(name, dir string) error {
	if err := createDataDir(dir); err != nil {
		return err
	}
	filePath := filepath.Join(dir, name)
	if err := checkDataFile(filePath); err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return os.WriteFile(filePath, nil, dataFileMode)
		case errors.Is(err, ErrFileModeInvalid):
			return os.Chmod(filePath, dataFileMode)
		default:
			return err
		}
	}
	return nil
}

// checkEndpoint consumes and validates the given endpoint.
// It consumes the given endpoint using the GET method.
// The scheme, domain and path properties are mandatory.
// The response code must be StatusOK(200), otherwise returns error
func checkEndpoint(endpoint string) error {
	uri, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return err
	}
	if uri.Path == "" {
		return &url.Error{Op: "parse", URL: endpoint, Err: ErrURLPathEmpty}
	}

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return err
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log().Error().Err(err).Msg("close response body failed")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ErrInvalidRespCode
	}
	return nil
}

// validNumType validate the number type of the given number. e.g. "Even","Odd"
func validNumType(num int, nType ct.NumberType) bool {
	switch nType {
	case ct.EvenNum:
		if num%2 == 0 {
			return true
		}
	case ct.OddNum:
		if num%2 != 0 {
			return true
		}
	}
	return false
}
