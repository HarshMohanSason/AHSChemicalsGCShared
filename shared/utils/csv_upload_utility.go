package utils

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
)

const firestoreMaxBatchSize = 500

// UploadTaxRatesInBulkToFirestore reads California county tax rate data from a CSV file
// and writes the data to the "tax_rates" collection in Firestore using batch writes.
//
// CSV Requirements:
// - File path: "./extras/ca_county_details.csv"
// - Expected Columns:
//   - County: index 3
//   - City: index 4
//   - Rate: index 6
//
// Behavior:
// - Reads all CSV rows (skipping header).
// - Parses the tax rate (column index 6) as float64.
// - Batches up to 500 documents per Firestore write batch (limit per API).
// - Logs progress and any errors during the process.
func UploadTaxRatesInBulkToFirestore() {
	file, err := os.Open("./extras/ca_county_details.csv")
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	log.Println("Starting Firestore batch upload for county tax rates...")

	ctx := context.Background()
	batch := firebase_shared.FirestoreClient.Batch()
	count := 0

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		rate, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			log.Printf("Skipping row due to rate conversion error (%s - %s): %v", row[3], row[4], err)
			continue
		}

		taxDoc := map[string]any{
			"state":  "CALIFORNIA",
			"county": strings.ToUpper(row[3]),
			"city":   strings.ToUpper(row[4]),
			"rate":   rate,
		}

		docRef := firebase_shared.FirestoreClient.Collection("tax_rates").NewDoc()
		batch.Set(docRef, taxDoc)
		count++

		if count == firestoreMaxBatchSize {
			if _, err := batch.Commit(ctx); err != nil {
				log.Printf("Error committing tax rate batch: %v", err)
			}
			batch = firebase_shared.FirestoreClient.Batch()
			count = 0
		}
	}

	if count > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			log.Printf("Error committing final tax rate batch: %v", err)
		}
	}

	log.Println("Successfully uploaded all county tax rates to Firestore.")
}

// UploadPropertiesDetailsInBulkToFirestore reads property details from a CSV file
// and writes the data to the "properties" collection in Firestore using batch writes.
//
// CSV Requirements:
// - File path: "./extras/property_details.csv"
// - Expected Columns:
//   - Name: index 0
//   - Street: index 1
//   - City: index 2
//   - State: index 3
//   - County: index 4
//   - Country: index 5
//   - Postal Code: index 6
//   - SVG X Position: index 9
//   - SVG Y Position: index 10
//
// Behavior:
// - Reads all CSV rows (skipping header).
// - Converts SVG positions to float64.
// - Batches up to 500 documents per Firestore write batch.
// - Logs progress and errors throughout the process.
func UploadPropertiesDetailsInBulkToFirestore() {
	file, err := os.Open("./extras/property_details.csv")
	if err != nil {
		log.Fatal("Error opening property CSV file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	log.Println("Uploading property details to Firestore using batch writes...")

	ctx := context.Background()
	batch := firebase_shared.FirestoreClient.Batch()
	count := 0

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		xPos, err := strconv.ParseFloat(row[9], 64)
		if err != nil {
			log.Printf("Error converting X position for row %d: %v", i, err)
			continue
		}

		yPos, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			log.Printf("Error converting Y position for row %d: %v", i, err)
			continue
		}

		property := map[string]any{
			"name":             strings.ToUpper(row[0]),
			"street":           strings.ToUpper(row[1]),
			"city":             strings.ToUpper(row[2]),
			"state":            strings.ToUpper(row[3]),
			"county":           strings.ToUpper(row[4]),
			"country":          strings.ToUpper(row[5]),
			"postal":           row[6],
			"svg_circle_x_pos": xPos,
			"svg_circle_y_pos": yPos,
		}

		docRef := firebase_shared.FirestoreClient.Collection("properties").NewDoc()
		batch.Set(docRef, property)
		count++

		if count == firestoreMaxBatchSize {
			if _, err := batch.Commit(ctx); err != nil {
				log.Printf("Error committing property batch: %v", err)
			}
			batch = firebase_shared.FirestoreClient.Batch()
			count = 0
		}
	}

	if count > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			log.Printf("Error committing final property batch: %v", err)
		}
	}

	log.Println("Successfully uploaded all property details to Firestore.")
}