package repository

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/logger"
)

// The indexes of the cocktail fields regarding the specific column in the CSV file.
const (
	idIdx             csvColIdx = 0
	nameIdx           csvColIdx = 1
	alcoholicIdx      csvColIdx = 2
	categoryIdx       csvColIdx = 3
	ingredientsIdx    csvColIdx = 4
	instructionsIdx   csvColIdx = 5
	glassIdx          csvColIdx = 6
	ibaIdx            csvColIdx = 7
	imgAttributionIdx csvColIdx = 8
	imgSrcIdx         csvColIdx = 9
	tagsIdx           csvColIdx = 10
	thumbIdx          csvColIdx = 11
	videoIdx          csvColIdx = 12
	srcDateIdx        csvColIdx = 13
	createdAtIdx      csvColIdx = 14
	updatedAtIdx      csvColIdx = 15
)

// csvHeadersMap are the names of the fields used in the headers/columns of the CSV file
var csvHeadersMap = map[csvColIdx]string{
	idIdx:             "id",
	nameIdx:           "name",
	alcoholicIdx:      "alcoholic",
	categoryIdx:       "category",
	ingredientsIdx:    "ingredients",
	instructionsIdx:   "instructions",
	glassIdx:          "glass",
	ibaIdx:            "iba",
	imgAttributionIdx: "image_attribution",
	imgSrcIdx:         "image_source",
	tagsIdx:           "tags",
	thumbIdx:          "thumb",
	videoIdx:          "video",
	srcDateIdx:        "source_date",
	createdAtIdx:      "created_at",
	updatedAtIdx:      "updated_at",
}

// csvColIdx represents the column's position in the csv file.
type csvColIdx int

// cocktailCsvRec represents a cocktail record of type CSV
type cocktailCsvRec []string

// parse returns a valid entity.Cocktail instance.
func (cr cocktailCsvRec) parse() (entity.Cocktail, error) {
	if len(cr) == 0 {
		return entity.Cocktail{}, ErrCSVRecEmpty
	}

	numFields := len(csvHeadersMap)
	if len(cr) != numFields {
		logger.Log().Warn().Str("required", fmt.Sprintf("%d/%d", len(cr), numFields)).Str("record", strings.Join(cr[:], ",")).
			Msg("parse: wrong number of fields")
	}

	// Preallocate a csv record matching the number of fields to avoid nil errors.
	rec := make([]string, numFields)
	copy(rec, cr)

	recID, err := strconv.Atoi(rec[idIdx])
	if err != nil {
		logger.Log().Error().Err(err).Str("id", rec[idIdx]).
			Msgf("parse: ID failure")
		return entity.Cocktail{}, err
	}

	if rec[nameIdx] == "" {
		logger.Log().Error().Err(ErrCocktailNameEmpty).Str("name", rec[nameIdx]).
			Msgf("parse: Name failure")
		return entity.Cocktail{}, ErrCocktailNameEmpty
	}

	var ingredients []entity.Ingredient
	if err := json.Unmarshal([]byte(rec[ingredientsIdx]), &ingredients); err != nil {
		logger.Log().Error().Err(err).Str("ingredients", rec[ingredientsIdx]).
			Msgf("parse: unmarshalling Ingredients failure")
		return entity.Cocktail{}, err
	}
	if len(ingredients) == 0 {
		logger.Log().Error().Err(ErrCocktailIngredientsEmpty).Str("ingredients", rec[ingredientsIdx]).
			Msgf("parse: Ingredients failure")
		return entity.Cocktail{}, ErrCocktailIngredientsEmpty
	}

	if rec[instructionsIdx] == "" {
		logger.Log().Error().Err(ErrCocktailInstructionsEmpty).Str("name", rec[instructionsIdx]).
			Msgf("parse: Instructions failure")
		return entity.Cocktail{}, ErrCocktailInstructionsEmpty
	}

	srcDate, err := time.Parse(time.DateTime, rec[srcDateIdx])
	if err != nil {
		logger.Log().Error().Err(err).Str("source_date", rec[srcDateIdx]).
			Msgf("parse: Source Date failure")
		return entity.Cocktail{}, err
	}

	createdAt, err := time.Parse(time.DateTime, rec[createdAtIdx])
	if err != nil {
		logger.Log().Error().Err(err).Str("created_at", rec[createdAtIdx]).
			Msgf("parse: Created At failure")
		return entity.Cocktail{}, err
	}

	updatedAt, err := time.Parse(time.DateTime, rec[updatedAtIdx])
	if err != nil {
		logger.Log().Error().Err(err).Str("updated_at", rec[updatedAtIdx]).
			Msgf("parse: Updated At failure")
		return entity.Cocktail{}, err
	}

	return entity.Cocktail{
		ID:             recID,
		Name:           rec[nameIdx],
		Alcoholic:      rec[alcoholicIdx],
		Category:       rec[categoryIdx],
		Instructions:   rec[instructionsIdx],
		Ingredients:    ingredients,
		Glass:          rec[glassIdx],
		IBA:            rec[ibaIdx],
		ImgAttribution: rec[imgAttributionIdx],
		ImgSrc:         rec[imgSrcIdx],
		Tags:           rec[tagsIdx],
		Thumb:          rec[thumbIdx],
		Video:          rec[videoIdx],
		SrcDate:        srcDate,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

// parseCsvRec returns a valid csv record from the given entity.Cocktail.
func parseCsvRec(c entity.Cocktail) ([]string, error) {
	ingredients, err := json.Marshal(c.Ingredients)
	if err != nil {
		return nil, err
	}

	rec := make([]string, len(csvHeadersMap))
	rec[idIdx] = strconv.Itoa(c.ID)
	rec[nameIdx] = c.Name
	rec[alcoholicIdx] = c.Alcoholic
	rec[categoryIdx] = c.Category
	rec[ingredientsIdx] = string(ingredients)
	rec[instructionsIdx] = c.Instructions
	rec[glassIdx] = c.Glass
	rec[ibaIdx] = c.IBA
	rec[imgAttributionIdx] = c.ImgAttribution
	rec[imgSrcIdx] = c.ImgSrc
	rec[tagsIdx] = c.Tags
	rec[thumbIdx] = c.Thumb
	rec[videoIdx] = c.Video
	rec[srcDateIdx] = c.SrcDate.Format(time.DateTime)
	rec[createdAtIdx] = c.CreatedAt.Format(time.DateTime)
	rec[updatedAtIdx] = c.UpdatedAt.Format(time.DateTime)
	return rec, nil
}

// drink represents the fetched JSON record from the public API.
type drink struct {
	Alcoholic        string `json:"strAlcoholic"`
	Category         string `json:"strCategory"`
	DateModified     string `json:"dateModified"`
	DrinkId          string `json:"idDrink"`
	DrinkName        string `json:"strDrink"`
	DrinkAlternate   string `json:"strDrinkAlternate"`
	DrinkThumb       string `json:"strDrinkThumb"`
	Glass            string `json:"strGlass"`
	IBA              string `json:"strIBA"`
	ImageSource      string `json:"strImageSource"`
	ImageAttribution string `json:"strImageAttribution"`
	Instructions     string `json:"strInstructions"`
	Ingredient1      string `json:"strIngredient1"`
	Ingredient2      string `json:"strIngredient2"`
	Ingredient3      string `json:"strIngredient3"`
	Ingredient4      string `json:"strIngredient4"`
	Ingredient5      string `json:"strIngredient5"`
	Ingredient6      string `json:"strIngredient6"`
	Ingredient7      string `json:"strIngredient7"`
	Ingredient8      string `json:"strIngredient8"`
	Ingredient9      string `json:"strIngredient9"`
	Ingredient10     string `json:"strIngredient10"`
	Ingredient11     string `json:"strIngredient11"`
	Ingredient12     string `json:"strIngredient12"`
	Ingredient13     string `json:"strIngredient13"`
	Ingredient14     string `json:"strIngredient14"`
	Ingredient15     string `json:"strIngredient15"`
	Measure1         string `json:"strMeasure1"`
	Measure2         string `json:"strMeasure2"`
	Measure3         string `json:"strMeasure3"`
	Measure4         string `json:"strMeasure4"`
	Measure5         string `json:"strMeasure5"`
	Measure6         string `json:"strMeasure6"`
	Measure7         string `json:"strMeasure7"`
	Measure8         string `json:"strMeasure8"`
	Measure9         string `json:"strMeasure9"`
	Measure10        string `json:"strMeasure10"`
	Measure11        string `json:"strMeasure11"`
	Measure12        string `json:"strMeasure12"`
	Measure13        string `json:"strMeasure13"`
	Measure14        string `json:"strMeasure14"`
	Measure15        string `json:"strMeasure15"`
	Tags             string `json:"strTags"`
	Video            string `json:"strVideo"`
}

// drinksData holds all the fetched drink records from the data API.
type drinksData struct {
	Drinks []drink `json:"drinks"`
}

// getIngredients returns a valid entity.Ingredient list
func (d drink) getIngredients() []entity.Ingredient {
	ingredients := make([]entity.Ingredient, 0)

	if d.Ingredient1 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient1, Measure: d.Measure1})
	}
	if d.Ingredient2 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient2, Measure: d.Measure2})
	}
	if d.Ingredient3 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient3, Measure: d.Measure3})
	}
	if d.Ingredient4 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient4, Measure: d.Measure4})
	}
	if d.Ingredient5 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient5, Measure: d.Measure5})
	}
	if d.Ingredient6 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient6, Measure: d.Measure6})
	}
	if d.Ingredient7 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient7, Measure: d.Measure7})
	}
	if d.Ingredient8 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient8, Measure: d.Measure8})
	}
	if d.Ingredient9 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient9, Measure: d.Measure9})
	}
	if d.Ingredient10 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient10, Measure: d.Measure10})
	}
	if d.Ingredient11 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient11, Measure: d.Measure11})
	}
	if d.Ingredient12 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient12, Measure: d.Measure12})
	}
	if d.Ingredient13 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient13, Measure: d.Measure13})
	}
	if d.Ingredient14 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient14, Measure: d.Measure14})
	}
	if d.Ingredient15 != "" {
		ingredients = append(ingredients, entity.Ingredient{Name: d.Ingredient15, Measure: d.Measure15})
	}
	return ingredients
}

// parse returns a valid entity.Cocktail instance
func (d drink) parse() (entity.Cocktail, error) {
	if d == (drink{}) {
		return entity.Cocktail{}, ErrJsonRecEmpty
	}

	id, err := strconv.Atoi(d.DrinkId)
	if err != nil {
		logger.Log().Error().Err(err).Str("id", d.DrinkId).Msgf("parse: ID failure")
		return entity.Cocktail{}, err
	}

	if d.DrinkName == "" {
		logger.Log().Error().Err(ErrCocktailNameEmpty).Msgf("parse: Name failure")
		return entity.Cocktail{}, ErrCocktailNameEmpty
	}

	if d.Instructions == "" {
		logger.Log().Error().Err(ErrCocktailInstructionsEmpty).Msgf("parse: Instructions failure")
		return entity.Cocktail{}, ErrCocktailInstructionsEmpty
	}

	ingredients := d.getIngredients()
	if len(ingredients) == 0 {
		logger.Log().Error().Err(ErrCocktailIngredientsEmpty).Msgf("parse: Ingredients failure")
		return entity.Cocktail{}, ErrCocktailIngredientsEmpty
	}

	srcDate, err := time.Parse(time.DateTime, d.DateModified)
	if err != nil {
		logger.Log().Error().Err(err).Msgf("parse: Source Date failure")
		return entity.Cocktail{}, err
	}

	return entity.Cocktail{
		ID:             id,
		Name:           d.DrinkName,
		Alcoholic:      d.Alcoholic,
		Category:       d.Category,
		Ingredients:    ingredients,
		Instructions:   d.Instructions,
		Glass:          d.Glass,
		IBA:            d.IBA,
		ImgAttribution: d.ImageAttribution,
		ImgSrc:         d.ImageSource,
		Tags:           d.Tags,
		Thumb:          d.DrinkThumb,
		Video:          d.Video,
		SrcDate:        srcDate,
	}, nil
}

// workerPool represents the Worker Pool pattern.
type workerPool struct {
	nType      ct.NumberType
	maxJobs    int
	maxWorkers int
	jobsWorker int
	resp       []entity.Cocktail
	jobs       chan cocktailCsvRec
	results    chan entity.Cocktail
}

// newWorkerPool returns a new workerPool implementation.
func newWorkerPool(nType ct.NumberType, maxJobs, jobsWorker int) (workerPool, error) {
	if nType == ct.InvalidNum || maxJobs == 0 || jobsWorker == 0 || jobsWorker > maxJobs {
		return workerPool{}, ErrWPInvalidArgs
	}
	maxWorkers := int(
		math.Ceil(float64(maxJobs) / float64(jobsWorker)),
	)
	return workerPool{
		nType:      nType,
		maxJobs:    maxJobs,
		maxWorkers: maxWorkers,
		jobsWorker: jobsWorker,
		resp:       make([]entity.Cocktail, 0),
		jobs:       make(chan cocktailCsvRec, maxJobs),
		results:    make(chan entity.Cocktail, maxJobs),
	}, nil
}

// runWorkers starts the workers and dispatches jobs to workers.
// It waits for all workers to finish and then closes the results channel
func (wp *workerPool) runWorkers() {
	wg := new(sync.WaitGroup)
	for i := 1; i <= wp.maxWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			wp.worker(id)
		}(i)
	}
	go func() {
		wg.Wait()
		close(wp.results)
	}()
}

// producer reads a valid csv record and send it to the jobs queue.
// If reach the end of file, stop sending jobs.
// Any bad record is skipped.
func (wp *workerPool) producer(r *csv.Reader) {
	go func() {
		for i := 0; i < wp.maxJobs; {
			rec, err := r.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				logger.Log().Warn().Err(err).Msgf("producer: read csv record failed, skipped")
				continue
			}
			wp.jobs <- rec
			i++
		}
		close(wp.jobs)
	}()
}

// consumer receives a valid response from results queue and add it to the responses list.
func (wp *workerPool) consumer() {
	for i := 0; i < wp.maxJobs; i++ {
		resp, open := <-wp.results
		if !open {
			break
		}
		wp.resp = append(wp.resp, resp)
	}
}

// worker processes jobs from queue and send them to the results queue
func (wp *workerPool) worker(id int) {
	for i := 0; i < wp.jobsWorker; i++ {
		job, open := <-wp.jobs
		if !open {
			break
		}
		resp, err := job.parse()
		if err != nil {
			logger.Log().Error().Err(err).Str("record", strings.Join(job[:], ",")).
				Msgf("worker(%d): parsing record failed, skipped.", id)
			continue
		}
		if validNumType(resp.ID, wp.nType) {
			wp.results <- resp
		}
	}

}
