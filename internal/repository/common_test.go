package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
	workDir string
}

func TestFileTestSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}

func (s *FileTestSuite) SetupSuite() {
	workDir := "testdata"
	if _, err := os.Stat(workDir); errors.Is(err, os.ErrNotExist) {
		require.NoError(s.T(), os.Mkdir(workDir, os.ModePerm),
			"create the work directory is mandatory")
	} else {
		require.Nil(s.T(), err)
	}
	s.workDir = workDir
}

func (s *FileTestSuite) TearDownSuite() {
	assert.NoError(s.T(), os.RemoveAll(s.workDir),
		fmt.Sprintf("remove the work directory %q is mandatory", s.workDir))
}

func (s *FileTestSuite) TearDownTest() {
	files, err := os.ReadDir(s.workDir)
	assert.Nil(s.T(), err)
	for _, f := range files {
		assert.NoError(s.T(), os.RemoveAll(filepath.Join(s.workDir, f.Name())),
			fmt.Sprintf("remove file %q is mandatory", f.Name()))
	}
}

func (s *FileTestSuite) TestCheckDataDir() {
	tests := []struct {
		name     string
		dirName  string
		err      error
		wantDir  bool
		wantFile bool
	}{
		{
			name:    "Empty",
			dirName: "",
			err:     ErrDirNameEmpty,
			wantDir: false,
		},
		{
			name:    "Arbitrary",
			dirName: filepath.Join(s.workDir, "foo"),
			err:     fs.ErrNotExist,
			wantDir: false,
		},
		{
			name:     "Is not a directory",
			dirName:  filepath.Join(s.workDir, "foo"),
			err:      ErrIsNotDir,
			wantDir:  false,
			wantFile: true,
		},
		{
			name:    "Valid",
			dirName: filepath.Join(s.workDir, "valid-dir"),
			err:     nil,
			wantDir: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.wantFile {
				require.NoError(t, os.WriteFile(tt.dirName, nil, dataFileMode))
			}
			if tt.wantDir {
				require.NoError(t, os.Mkdir(tt.dirName, dataDirMode))
			}
			err := checkDataDir(tt.dirName)
			if tt.err != nil {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func (s *FileTestSuite) TestCreateDataDir() {
	type dir struct {
		name string
		mode os.FileMode
	}
	tests := []struct {
		name    string
		dir     dir
		err     error
		wantDir bool
	}{
		{
			name:    "Empty",
			dir:     dir{name: ""},
			err:     ErrDirNameEmpty,
			wantDir: false,
		},
		{
			name: "Exists",
			dir: dir{
				name: filepath.Join(s.workDir, "foo"),
				mode: dataDirMode,
			},
			err:     nil,
			wantDir: true,
		},
		{
			name:    "Created",
			dir:     dir{name: filepath.Join(s.workDir, "foo")},
			err:     nil,
			wantDir: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.wantDir {
				require.NoError(t, os.Mkdir(tt.dir.name, tt.dir.mode))
			}
			err := createDataDir(tt.dir.name)
			if tt.err != nil {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func (s *FileTestSuite) TestCheckDataFile() {
	type dataFile struct {
		name string
		mode os.FileMode
	}
	tests := []struct {
		name     string
		file     dataFile
		err      error
		wantFile bool
	}{
		{
			name:     "Empty",
			file:     dataFile{name: ""},
			err:      ErrFileNameEmpty,
			wantFile: false,
		},
		{
			name: "Arbitrary",
			file: dataFile{
				name: filepath.Join(s.workDir, "foo.txt"),
				mode: 0000,
			},
			err:      os.ErrNotExist,
			wantFile: false,
		},
		{
			name: "Is Directory",
			file: dataFile{
				name: s.workDir,
				mode: 0000,
			},
			err:      ErrFilePathIsDir,
			wantFile: false,
		},
		{
			name: "Invalid mode",
			file: dataFile{
				name: filepath.Join(s.workDir, "invalid-mode.txt"),
				mode: 0666,
			},
			err:      ErrFileModeInvalid,
			wantFile: true,
		},
		{
			name: "Valid",
			file: dataFile{
				name: filepath.Join(s.workDir, "valid.txt"),
				mode: dataFileMode,
			},
			err:      nil,
			wantFile: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.wantFile {
				require.NoError(t, os.WriteFile(tt.file.name, nil, tt.file.mode))
			}

			err := checkDataFile(tt.file.name)
			if err != nil {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func (s *FileTestSuite) TestCreateDataFile() {
	type args struct {
		name string
		dir  string
	}
	tests := []struct {
		name     string
		args     args
		err      error
		fileMode os.FileMode
		wantFile bool
	}{
		{
			name:     "All empty",
			args:     args{name: "", dir: ""},
			err:      ErrDirNameEmpty,
			wantFile: false,
		},
		{
			name:     "Directory name empty",
			args:     args{name: "foo", dir: ""},
			err:      ErrDirNameEmpty,
			wantFile: false,
		},
		{
			name:     "File name empty",
			args:     args{name: "", dir: s.workDir},
			err:      ErrFilePathIsDir,
			wantFile: false,
		},
		{
			name:     "Exists with bad permissions",
			args:     args{name: "bad-mode", dir: s.workDir},
			fileMode: 0666,
			err:      nil,
			wantFile: true,
		},
		{
			name:     "Exists",
			args:     args{name: "valid", dir: s.workDir},
			fileMode: dataFileMode,
			err:      nil,
			wantFile: true,
		},
		{
			name:     "Created",
			args:     args{name: "created", dir: s.workDir},
			err:      nil,
			wantFile: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.wantFile {
				require.NoError(t,
					os.WriteFile(filepath.Join(tt.args.dir, tt.args.name), nil, tt.fileMode),
				)
			}
			err := createDataFile(tt.args.name, tt.args.dir)
			if tt.err != nil {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestCheckEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		err      error
	}{
		{
			name:     "Empty",
			endpoint: "",
			err:      &url.Error{},
		},
		{
			name:     "Bad Scheme",
			endpoint: "htttp://thecocktaildb.com/api/json/v1/1/search.php?f=a",
			err:      &url.Error{},
		},
		{
			name:     "Bad Domain",
			endpoint: "https://foo,.com/api/someendpoint",
			err:      &url.Error{},
		},
		{
			name:     "Path empty",
			endpoint: "https://thecocktaildb.com",
			err:      &url.Error{Err: ErrURLPathEmpty},
		},
		{
			name:     "Invalid Response Code",
			endpoint: "https://thecocktaildb.com/api/json/v1/1/foo.php",
			err:      ErrInvalidRespCode,
		},
		{
			name:     "Valid",
			endpoint: "https://thecocktaildb.com/api/json/v1/1/search.php?f=a",
			err:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkEndpoint(tt.endpoint)
			if tt.err != nil {
				assert.IsType(t, tt.err, err)
				if errW := errors.Unwrap(tt.err); errW != nil {
					assert.IsType(t, errW, errors.Unwrap(err))
				}
			}
		})
	}
}
