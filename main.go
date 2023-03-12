package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GDAL struct {
	Parameters parameters `json:"PARAMETERS"`
	Outputs    output     `json:"OUTPUTS"`
}

type parameters struct {
	Input        string `json:"INPUT"`
	InputRaster  string `json:"INPUT_RASTER"`
	RasterBand   string `json:"RASTER_BAND"`
	ColumnPrefix string `json:"COLUMN_PREFIX"`
	Statistics   string `json:"STATISTICS"`
}

type output struct {
	Output string `json:"OUTPUT"`
}

func main() {

	// load csv file "input_rasters.csv" into a slice of strings
	// each string is a path to a raster file
	// for example: "C:/Users/username/Desktop/raster.tif"
	rasters := loadCSV("input_rasters.csv")

	listFiles := []GDAL{}

	for _, raster := range rasters {
		class := getMetadata(raster)

		parameters := parameters{
			Input:        "'C:/Users/pedro/Documents/workspace/prep_2023/data/colima/agebs_repro.gpkg|layername=reproyectada'",
			InputRaster:  fmt.Sprintf("'C:/Users/pedro/Documents/workspace/prep_2023/%s'", raster),
			RasterBand:   "1",
			ColumnPrefix: "'_'",
			Statistics:   "[0,1,2]",
		}

		output := output{
			Output: fmt.Sprintf("'C:/Users/pedro/Documents/workspace/prep_2023/data/zonalStats/%s_%s.gpkg'" , class[0],class[1]),
		}

		gdal := GDAL{
			Parameters: parameters,
			Outputs:    output,
		}

		listFiles = append(listFiles, gdal)

	}

	// save listFiles as a json file
	saveJSON(listFiles, "output2.json")

}

func saveJSON(listFiles []GDAL, path string) {
	// save listFiles as a json file

	// create json file
	jsonFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// write json file
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(listFiles)
	if err != nil {
		log.Fatal(err)
	}

}

func getMetadata(path string) []string {
	// get metadata from a raster file
	// for example: "data/Rasters/Hydraulic/Depth____0.0.asc"
	// split into a slice of strings
	// for example: ["Depth", "0.0"]

	// get file name
	fileName := filepath.Base(path)

	// split file name into a slice of strings
	// for example: ["Depth____0.0.asc"]
	split := strings.Split(fileName, ".")
	// for example: ["Depth____0.0"]
	split = strings.Split(split[0], "_")

	strings := []string{}

	for _, s := range split {
		if s != "" {
			strings = append(strings, s)
		}
	}

	return strings
}

func loadCSV(path string) []string {
	// load csv file into a slice of strings
	// each string is a path to a raster file
	// for example: "data/Rasters/Hydraulic/Depth____0.0.asc"

	// open csv file
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// read csv file into a slice of strings
	r := csv.NewReader(csvFile)
	r.Comma = ','
	r.Comment = '#'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert slice of slices of strings to slice of strings
	var rasters []string
	for _, record := range records[0] {
		rasters = append(rasters, record)
	}

	return rasters
}
