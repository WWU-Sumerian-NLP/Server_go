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
)

// Insert Relation Patterns to the database from the react app csv/tsv file
func InsertRelationPatterns(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("path_to_csv")
	if err != nil {
		println("ERROR")
		panic(err)
	}

	//copy example
	f, err := os.OpenFile("./relation_input.tsv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal()
	}

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

	createdFile, err := os.Open("relation_input.tsv")
	if err != nil {
		log.Fatal()
	}

	reader := csv.NewReader(createdFile)
	reader.Comma = '\t'

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal()
	}

	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	fmt.Printf("db: %v\n", db)
	// db.InsertRelation()

	for i, row := range data {
		if i != 0 {
			createdRelations := &database.Relations{
				RelationType: row[0],
				SubjectTag:   row[1],
				ObjectTag:    row[2],
				RegexRules:   row[3],
				Tags:         row[4],
			}
			fmt.Printf("createdRelations: %v\n", createdRelations)

			db.InsertRelation(*createdRelations)

		}
	}
	json.NewEncoder(w).Encode(contents)
}

// Get All Relation Pattern from the database to the react app
func GetRelationPatterns(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	relations, err := db.GetAllRelations()
	fmt.Printf("relations: %v\n", relations)
	if err != nil {
		log.Fatal()
	}
	test := json.NewEncoder(w).Encode(relations)
	fmt.Printf("test: %v\n", test)

}

// Get All Relation Data from the database and run the relation extraction pipeline
func RunRelationExtraction(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	relations, err := db.ListRelations(2)
	if err != nil {
		log.Fatal()
	}
	// //Run Pipeline
	path := "urr3_annotations.tsv"
	destPath := "urr3_ie_annotations.tsv"
	cdliParser := IE_Extractor.NewCDLIParser(path)
	RelationExtractorRB := IE_Extractor.NewRelationExtractorRB(cdliParser.Out)
	for _, relation := range relations {

		relationData := IE_Extractor.NewRelationData(relation.RelationType, relation.SubjectTag, relation.ObjectTag, relation.RegexRules, relation.Tags)
		RelationExtractorRB.RelationDataList = append(RelationExtractorRB.RelationDataList, *relationData)

	}
	dataWriter := IE_Extractor.NewDataWriter(destPath, RelationExtractorRB.Out)
	dataWriter.WaitUntilDone()
	json.NewEncoder(w).Encode(relations)
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

	router.HandleFunc("/insertRelations", InsertRelationPatterns).Methods("POST")
	router.HandleFunc("/getRelations", GetRelationPatterns).Methods("GET")
	router.HandleFunc("/runRelationExtraction", RunRelationExtraction).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
