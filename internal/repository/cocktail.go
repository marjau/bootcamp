package repository

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/marcos-wz/capstone-go-bootcamp/internal/config"
	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/logger"
)

// Cocktail represents the Cocktail repository.
type Cocktail struct {
	csv        config.CsvDB
	dataAPI    config.DataAPI
	httpClient HttpClient
}

// NewCocktail returns a new Cocktail repository implementation.
func NewCocktail(cfg config.Config) (Cocktail, error) {
	dataAPI := cfg.HTTP.DataAPI
	csvDB := cfg.Database.Csv
	if err := checkEndpoint(dataAPI.URL()); err != nil {
		return Cocktail{}, &DataApiErr{err}
	}
	if err := createDataFile(csvDB.FileName(), csvDB.DataDir()); err != nil {
		return Cocktail{}, &CsvErr{err}
	}

	logger.Log().Debug().
		Str("csv_file", csvDB.FilePath()).
		Str("data_api", dataAPI.URL()).
		Msg("created Cocktail repository")
	return Cocktail{
		csv:        csvDB,
		dataAPI:    dataAPI,
		httpClient: &http.Client{},
	}, nil
}

// ReadAll returns all entity.Cocktail records from the CSV data file.
func (c Cocktail) ReadAll() ([]entity.Cocktail, error) {
	fd, err := os.Open(c.csv.FilePath())
	if err != nil {
		logger.Log().Error().Err(err).Str("file", c.csv.FilePath()).Msg("ReadAll: open csv file failed")
		return nil, &CsvErr{err}
	}
	defer func() {
		if err := fd.Close(); err != nil {
			logger.Log().Error().Err(err).Str("file", c.csv.FilePath()).Msg("ReadAll: close csv file failed")
		}
	}()

	reader := csv.NewReader(fd)
	reader.TrimLeadingSpace = true
	cocktails := make([]entity.Cocktail, 0)
	for {
		var rec cocktailCsvRec
		rec, err = reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logger.Log().Warn().Err(err).Msg("ReadAll: read record failed, skipped")
			continue
		}

		cocktail, err := rec.parse()
		if err != nil {
			logger.Log().Error().Err(err).Str("record", strings.Join(rec[:], ",")).Msg("ReadAll: parsing record failed, skipped")
			continue
		}
		cocktails = append(cocktails, cocktail)
	}

	return cocktails, nil
}

// ReadCC reads n number of csv records concurrently and returns a list of entity.Cocktail.
// It is based on the worker-pool pattern.
// nType: Is the number type. e.g. odd,even,...
// maxJobs: is the amount of valid csv records to be processed.
// jWorker: is the amount of jobs that each worker performs.
func (c Cocktail) ReadCC(nType ct.NumberType, maxJobs, jWorker int) ([]entity.Cocktail, error) {
	fd, err := os.Open(c.csv.FilePath())
	if err != nil {
		logger.Log().Error().Err(err).Str("file", c.csv.FilePath()).Msg("ReadCC: open csv file failed")
		return nil, &CsvErr{err}
	}
	defer func() {
		if err := fd.Close(); err != nil {
			logger.Log().Error().Err(err).Str("file", c.csv.FilePath()).Msg("ReadCC: close csv file failed")
		}
	}()
	reader := csv.NewReader(fd)
	reader.TrimLeadingSpace = true

	wp, err := newWorkerPool(nType, maxJobs, jWorker)
	if err != nil {
		return nil, &CsvErr{err}
	}

	// start worker-pool
	wp.runWorkers()
	wp.producer(reader)
	wp.consumer()

	return wp.resp, nil
}

// Fetch returns a list of entity.Cocktail records from the data API.
func (c Cocktail) Fetch() ([]entity.Cocktail, error) {
	req, err := http.NewRequest(http.MethodGet, c.dataAPI.URL(), nil)
	if err != nil {
		return nil, &DataApiErr{err}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &DataApiErr{err}
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.Log().Error().Err(err).Msg("Fetch: close response body failed")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Log().Error().Err(err).Int("code", resp.StatusCode).Msg("Fetch: bad status code, expected 200")
		return nil, &DataApiErr{ErrInvalidRespCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &DataApiErr{err}
	}

	data := drinksData{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, &DataApiErr{err}
	}

	cocktails := make([]entity.Cocktail, 0)
	for i, rec := range data.Drinks {
		cocktail, errP := rec.parse()
		if errP != nil {
			logger.Log().Error().Err(errP).Int("line", i+1).Str("record", fmt.Sprintf("%v - %v", rec.DrinkId, rec.DrinkName)).
				Msg("Fetch: parsing cocktail failed, record skipped")
			continue
		}
		cocktails = append(cocktails, cocktail)
	}

	return cocktails, nil
}

// ReplaceDB replaces the database entirely with the given entity.Cocktail records.
func (c Cocktail) ReplaceDB(cocktails []entity.Cocktail) error {
	file := c.csv.FilePath()

	f, err := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, dataFileMode)
	if err != nil {
		logger.Log().Error().Err(err).Str("file", file).Msg("ReplaceDB: open csv file failed")
		return &CsvErr{err}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Log().Error().Err(err).Str("file", file).Msg("ReplaceDB: close csv file failed")
		}
	}(f)

	w := csv.NewWriter(f)
	for i, cocktail := range cocktails {
		rec, errP := parseCsvRec(cocktail)
		if errP != nil {
			logger.Log().Error().Err(err).Int("index", i).Str("cocktail", fmt.Sprintf("ID: %d, Name: %v", cocktail.ID, cocktail.Name)).
				Msg("ReplaceDB: parsing cocktail to csv record failed, discarded record")
			continue
		}
		if err := w.Write(rec); err != nil {
			logger.Log().Error().Err(err).Int("index", i).Str("file", file).Str("record", strings.Join(rec[:], ",")).
				Msg("ReplaceDB: writing record failed, discarded record")
		}

	}
	w.Flush()
	if err := w.Error(); err != nil {
		logger.Log().Error().Err(err).Msg("ReplaceDB: flush writer failed")
		return &CsvErr{err}
	}

	return nil
}
