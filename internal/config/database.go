package config

import "path/filepath"

// Database holds the configurations of the supported databases
type Database struct {
	driver string
	Csv    CsvDB
}

// Driver returns the configured database driver.
func (db Database) Driver() string {
	return db.driver
}

// CsvDB holds and retrieves the config properties of the CSV database.
type CsvDB struct {
	fileName string
	dataDir  string
}

// NewCsv returns a new CsvDB configuration implementation.
func NewCsv(fileName, dataDir string) CsvDB {
	return CsvDB{
		fileName: fileName,
		dataDir:  dataDir,
	}
}

// FileName returns the file name of the CSV database.
func (c CsvDB) FileName() string {
	return c.fileName
}

// DataDir returns the base data directory of the CSV database.
func (c CsvDB) DataDir() string {
	return c.dataDir
}

// FilePath returns the full path of the CSV database file.
// The file path is built by prefixing the data directory to the file name.
func (c CsvDB) FilePath() string {
	return filepath.Join(c.dataDir, c.fileName)
}
