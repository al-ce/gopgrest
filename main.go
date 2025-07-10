package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
)

// Set represents a performed set
type Set struct {
	Name    string
	Date    string
	Output  float32
	Unit    string
	Reps    int8
	Set     int8
	Notes   string
	Day     string
	Program string
	Tags    string
}

// readFile reads the contents of a file into a byte array
func readFile(fp *string) []byte {
	content, err := os.ReadFile(*fp)
	if err != nil {
		log.Fatal("Could not open json file: ", err)
	}
	return content
}

// prettyFormat formats the json data read from a file and returns it as a string
func prettyFormat(jsonContent []byte) string {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, jsonContent, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(prettyJson.Bytes())
}

// unmarshalHistory parses the exercises json data and returns exercise history as an array of Sets
func unmarshalHistory(jsonContent []byte) []Set {

	// The history file should be an array of objects that match a Set struct
	var history []Set

	err := json.Unmarshal([]byte(jsonContent), &history)
	if err != nil {
		panic(err)
	}
	return history
}

func main() {
	// Parse args
	fp := flag.String("f", "", "path to json history")
	flag.Parse()
	if *fp == "" {
		log.Fatal("path to file (-f flag) required")
	}

	// Read file content into bytes
	jsonContent := readFile(*&fp)

	// Pretty print the json file content
	formatted := prettyFormat(jsonContent)
	log.Println(formatted)
}
