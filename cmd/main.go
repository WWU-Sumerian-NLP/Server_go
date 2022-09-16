package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

/*

This server will be basic as it needs to just get CSV files and then parse them
as input to our information extraction pipeline


*/

// func LoadRelationsCSV() {
// }

// func GetRelationPatterns(w http.ResponseWriter, r *http.Request) {
// 	relationPatterns := LoadRelationsCSV()
// 	json.NewEncoder(w).Encode(relationPatterns)
// }

// func GetRelationData(w http.ResponseWriter, r *http.Request) {
// 	relationData := LoadRelationsCSV()
// 	json.NewEncoder(w).Encode(relationData)
// }

// Entity Data
// Name, Tag
// Link: https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
func GetEntityData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("path_to_csv")
	if err != nil {
		println("ERROR")
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("name: %v\n", name[0])
	io.Copy(&buf, file)
	contents := buf.String()
	fmt.Printf("contents: %v\n", contents)
	buf.Reset()
	json.NewEncoder(w).Encode(contents)
	// fmt.Printf("r: %v\n", r)
	println("get entity")
	// fmt.Printf("r.Body: %v\n", r.Body)
	// parse POST body as csv
	reader := csv.NewReader(r.Body)
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("data: %v\n", data)
	}
	// fmt.Printf("reader: %v\n", reader)
	var results [][]string
	for {

		//read one row from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		//add record to result set
		results = append(results, record)
		// fmt.Printf("reader: %v\n", reader)
	}
	fmt.Printf("results: %v\n", results)
	// test := json.NewEncoder(w).Encode(results)
	// fmt.Printf("test: %v\n", test)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/entity", GetEntityData).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8000", handler))
}
