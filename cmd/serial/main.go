package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat"
)

// record represents the structure of our data
type record struct {
	neighborhood string
	crim         float64
	zn           float64
	indus        float64
	chas         float64
	nox          float64
	rooms        float64
	age          float64
	dis          float64
	rad          float64
	tax          float64
	ptratio      float64
	lstat        float64
	mv           float64
}

type Result struct {
	XName    string
	A        float64
	B        float64
	RSquared float64
}

type IterationResult struct {
	Iteration   int
	CrimResult  Result
	RoomsResult Result
}

// checkErr is a simple error check helper function
func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// readData reads the CSV file and returns three slices: crim, rooms and mv
func readData(csvPath string) ([]float64, []float64, []float64, error) {
	// Open the file
	data, err := os.ReadFile(csvPath)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert \r line breaks to \n
	dataStr := strings.ReplaceAll(string(data), "\r", "\n")

	// Create a new CSV reader reading from the string
	reader := csv.NewReader(strings.NewReader(dataStr))
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read the header line
	_, err = reader.Read()
	if err != nil {
		return nil, nil, nil, err
	}

	// Prepare a slice for each record
	var records []record

	// Read the rest of the data
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, err
		}

		var rec record
		rec.neighborhood = line[0]
		rec.crim, err = strconv.ParseFloat(line[1], 64)
		checkErr(err)
		rec.zn, err = strconv.ParseFloat(line[2], 64)
		checkErr(err)
		rec.indus, err = strconv.ParseFloat(line[3], 64)
		checkErr(err)
		rec.chas, err = strconv.ParseFloat(line[4], 64)
		checkErr(err)
		rec.nox, err = strconv.ParseFloat(line[5], 64)
		checkErr(err)
		rec.rooms, err = strconv.ParseFloat(line[6], 64)
		checkErr(err)
		rec.age, err = strconv.ParseFloat(line[7], 64)
		checkErr(err)
		rec.dis, err = strconv.ParseFloat(line[8], 64)
		checkErr(err)
		rec.rad, err = strconv.ParseFloat(line[9], 64)
		checkErr(err)
		rec.tax, err = strconv.ParseFloat(line[10], 64)
		checkErr(err)
		rec.ptratio, err = strconv.ParseFloat(line[11], 64)
		checkErr(err)
		rec.lstat, err = strconv.ParseFloat(line[12], 64)
		checkErr(err)
		rec.mv, err = strconv.ParseFloat(line[13], 64)
		checkErr(err)

		records = append(records, rec)
	}

	// Create crim, rooms, mv slices
	var crim, rooms, mv []float64
	for _, rec := range records {
		crim = append(crim, rec.crim)
		rooms = append(rooms, rec.rooms)
		mv = append(mv, rec.mv)
	}

	return crim, rooms, mv, nil
}

func performIteration(iteration int, crim, rooms, mv []float64, resultCh chan<- IterationResult) {
	var weights []float64
	origin := false

	// calculate Crim
	a1, b1 := stat.LinearRegression(crim, mv, weights, origin)
	rsq1 := stat.RSquared(crim, mv, weights, a1, b1)

	// calculate Rooms
	a2, b2 := stat.LinearRegression(rooms, mv, weights, origin)
	rsq2 := stat.RSquared(rooms, mv, weights, a2, b2)

	resultCh <- IterationResult{
		Iteration:   iteration,
		CrimResult:  Result{XName: "Crim", A: a1, B: b1, RSquared: rsq1},
		RoomsResult: Result{XName: "Rooms", A: a2, B: b2, RSquared: rsq2},
	}
}

func performRegression(n int, crim, rooms, mv []float64, verbose bool) {
	// Create a channel for results
	resultCh := make(chan IterationResult, n)

	// Perform n iterations in serial
	for i := 0; i < n; i++ {
		performIteration(i+1, crim, rooms, mv, resultCh)
	}
	close(resultCh) // Close the channel after all iterations are done

	// Process results
	if verbose {
		for res := range resultCh {
			fmt.Printf("Iteration %d:\n", res.Iteration)
			fmt.Println("Crim vs Median Value:", fmt.Sprintf("%.2f + %.2f * %s", res.CrimResult.A, res.CrimResult.B, res.CrimResult.XName), ", R-squared:", res.CrimResult.RSquared)
			fmt.Println("Rooms vs Median Value:", fmt.Sprintf("%.2f + %.2f * %s", res.RoomsResult.A, res.RoomsResult.B, res.RoomsResult.XName), ", R-squared:", res.RoomsResult.RSquared)
		}
	}
	fmt.Println("Finished all iterations")
}

func main() {
	verbose := flag.String("verbose", "", "Print the iterations if verbose is set")
	flag.Parse()

	// Set path to the CSV file
	csvPath := "../../data/boston.csv"

	// Read the data
	crim, rooms, mv, err := readData(csvPath)
	checkErr(err)

	// Perform linear regression calculation 100 times
	performRegression(10000, crim, rooms, mv, strings.EqualFold(*verbose, "verbose"))
}
