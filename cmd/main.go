package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/urr3-drehem-KG/Data_Pipeline_go/IE_Extractor"

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

func GetRelationPatterns(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("path_to_csv")
	if err != nil {
		println("ERROR")
		panic(err)
	}

	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("name: %v\n", name)

	io.Copy(&buf, file)
	contents := buf.String()
	fmt.Printf("contents: %v\n", contents)
	buf.Reset()

	json.NewEncoder(w).Encode(contents)
	println("get relations")

	reader := csv.NewReader(r.Body)
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("data: %v\n", data)
	}
	path := "../CDLI_Extractor/output/urr3_annotations.tsv"
	destPath := "output/urr3_ie_annotations.tsv"

	cdliParser := IE_Extractor.NewCDLIParser(path)
	RelationExtractorRB := IE_Extractor.NewRelationExtractorRB(cdliParser.Out)

	for _, row := range data {
		fmt.Printf("row: %v\n", row)
		relationData := IE_Extractor.NewRelationData(row[0], row[1], row[2], row[3])
		RelationExtractorRB.RelationDataList = append(RelationExtractorRB.RelationDataList, *relationData)
	}
	dataWriter := IE_Extractor.NewDataWriter(destPath, RelationExtractorRB.Out)

	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		RelationExtractorRB.WaitUntilDone()
		cdliParser.WaitUntilDone()
	}()

}

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

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/entity", GetEntityData).Methods("POST")
	router.HandleFunc("/relations", GetRelationPatterns).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8000", handler))
}
