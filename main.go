package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
)

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

func readFile(fp *string) []byte {
	content, err := os.ReadFile(*fp)
	if err != nil {
		log.Fatal("Could not open json file: ", err)
	}
	return content
}

func prettyPrint(jsonContent []byte) {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, jsonContent, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Println(string(prettyJson.Bytes()))
}

func unmarshalHistory(jsonContent []byte) []Set {
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
	prettyPrint(jsonContent)
}
