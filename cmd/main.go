package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urr3-drehem-KG/Data_Pipeline_go/IE_Extractor"
	"github.com/urr3-drehem-KG/gRPC_Server_go/database"
	// "github.com/urr3-drehem-KG/gRPC_Server_go/database/"
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

	// db := NewInternalDB()
	//copy example
	f, err := os.OpenFile("./relation_input.tsv", os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	fmt.Printf("f: %v\n", f)

	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("name: %v\n", name)

	io.Copy(&buf, file)
	contents := buf.String()
	// fmt.Printf("contents: %v\n", contents)
	buf.Reset()

	println("get relations")

	createdFile, err := os.Open("relation_input.tsv")

	reader := csv.NewReader(createdFile)
	reader.Comma = '\t'

	data, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("data: %v\n", data)
	}
	fmt.Printf("data: %v\n", data)
	path := "urr3_annotations.tsv"
	destPath := "urr3_ie_annotations.tsv"

	db, err := database.NewInternalDB()
	fmt.Printf("db: %v\n", db)
	// db.InsertRelation()

	cdliParser := IE_Extractor.NewCDLIParser(path)
	RelationExtractorRB := IE_Extractor.NewRelationExtractorRB(cdliParser.Out)
	for i, row := range data {
		if i != 0 {

			relationData := IE_Extractor.NewRelationData(row[0], row[1], row[2], row[3], row[4])
			fmt.Printf("relationData: %v\n", relationData)

			RelationExtractorRB.RelationDataList = append(RelationExtractorRB.RelationDataList, *relationData)

			fmt.Printf("row: %v\n", row)
		}
	}

	dataWriter := IE_Extractor.NewDataWriter(destPath, RelationExtractorRB.Out)
	dataWriter.WaitUntilDone()
	json.NewEncoder(w).Encode(contents)

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
