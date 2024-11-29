package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type SimplifiedData struct {
	ID string `json:"id"`
}

func processJSON(filePath string) ([]SimplifiedData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	token, err := decoder.Token()
	if err != nil {
		log.Fatalln("Error reading opening token: ", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		log.Fatalln("Expeceted start of JSON array")
	}

	var out []SimplifiedData

	for decoder.More() {
		var simplified SimplifiedData
		err = decoder.Decode(&simplified)
		if err != nil {
			return nil, err
		}
		out = append(out, simplified)
	}

	// Read the closing bracket of the JSON array
	token, err = decoder.Token()
	if err != nil {
		log.Fatalln("Error reading closing token:", err)
	}

	// Check if the closing token is the end of the array
	if delim, ok := token.(json.Delim); !ok || delim != ']' {
		log.Fatalln("Expected end of JSON array")
	}

	return out, nil
}

func main() {
	out, err := processJSON("large-file.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(out))
}
