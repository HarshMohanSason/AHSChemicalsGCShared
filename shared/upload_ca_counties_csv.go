package shared

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
)

const maxConcurrentUploads = 100

func UploadCaCountyCsvToFirestore() {

	file, err := os.Open("./extras/ca_county_details.csv")
	if err != nil {
		log.Fatal("Error occurred reading the file, please try again")
	}

	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV:", err)
		return
	}

	log.Print("Uploading the county csv data to firestore...")

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentUploads)

	// Read through the list of records
	for i, row := range records {
		if i == 0 {
			continue
		}

		wg.Add(1)
		sem <- struct{}{} //Get a spot

		go func(row []string) {
			defer wg.Done()
			defer func() { <-sem }() // release spot

			//Convert the rate to float
			rate, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				log.Fatal("Error converting the rate to float")
			}

			countyObject := map[string]any{
				"state":  "CALIFORNIA",
				"county": row[3],
				"city":   row[4],
				"rate":   rate,
			}
			_, _, err = FirestoreClient.Collection("tax_rates").Add(context.Background(), countyObject)
			if err != nil {
				log.Printf("Error uploading record (%s - %s): %v", row[3], row[4], err)
			}
		}(row)
	}

	wg.Wait()
	log.Print("Added all the county data to firestore.")

}
