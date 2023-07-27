package main

import (
	"testing"
)

func BenchmarkPerformRegression(b *testing.B) {
	csvPath := "../../data/boston.csv"

	// Read the data
	crim, rooms, mv, err := readData(csvPath)
	checkErr(err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		performRegression(100, crim, rooms, mv, false)
	}
}
