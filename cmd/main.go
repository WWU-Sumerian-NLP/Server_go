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
	"github.com/urr3-drehem-KG/Data_Pipeline_go/CDLI_Extractor"
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
	fmt.Printf("err: %v\n", err)
	if err != nil {
		log.Fatal()
	}
	fmt.Printf("db: %v\n", db)

	for i, row := range data {
		fmt.Printf("row: %v\n", row)
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
	relations, err := db.GetAllRelations(0)
	fmt.Printf("relations: %v\n", relations)
	fmt.Printf("err: %v\n", err)
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
	relations, err := db.GetAllRelations(0)
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

// Entity Data
// Name, Tag
// Link: https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
func InsertEntities(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("path_to_entity_csv")
	if err != nil {
		println("ERROR")
		panic(err)
	}

	f, err := os.OpenFile("./entity_input.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal()
	}
	contents := buf.String()

	defer file.Close()
	io.Copy(f, file)
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("name: %v\n", name)

	io.Copy(&buf, file)
	buf.Reset()

	createdFile, err := os.Open("entity_input.csv")
	if err != nil {
		log.Fatal()
	}

	reader := csv.NewReader(createdFile)
	// reader.Comma = '\t'
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("data: %v\n", data)
	}
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}

	for i, row := range data {
		if i != 0 {
			createdEntity := &database.Entities{
				EntityName: row[0],
				EntityTag:  row[1],
			}
			fmt.Printf("createdEntity: %v\n", createdEntity)
			db.InsertEntity(*createdEntity)
		}
	}
	json.NewEncoder(w).Encode(contents)
}

func GetEntityData(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	entities, err := db.GetAllEntities(0)
	fmt.Printf("entities: %v\n", entities)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		log.Fatal()
	}
	json.NewEncoder(w).Encode(entities)
}

func RunEntityExtraction(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	entities, err := db.GetAllEntities(0)
	if err != nil {
		log.Fatal()
	}
	// //Run Pipeline
	path := "sumerian_tablets/cdli_result_20220525.txt"
	destPath := "urr3_annotations.tsv"
	atfParser := CDLI_Extractor.NewATFParser(path)
	translitCleaner := CDLI_Extractor.NewTransliterationCleaner(false, atfParser.Out)
	entityExtractor := CDLI_Extractor.NewCDLIEntityExtractor(translitCleaner.Out)

	for _, entity := range entities {
		entityExtractor.TempNERMap[entity.EntityName] = entity.EntityTag

	}
	dataWriter := CDLI_Extractor.NewDataWriter(destPath, entityExtractor.Out)
	dataWriter.WaitUntilDone()
	json.NewEncoder(w).Encode(entities)
}

///////////////////////////////////////////////////////////////////////////
func InsertRelationships(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("path_to_relationship_csv")
	if err != nil {
		println("ERROR")
		panic(err)
	}

	f, err := os.OpenFile("./relationship_input.tsv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal()
	}
	contents := buf.String()

	defer file.Close()
	io.Copy(f, file)
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("name: %v\n", name)

	io.Copy(&buf, file)
	buf.Reset()

	createdFile, err := os.Open("relationship_input.tsv")
	if err != nil {
		log.Fatal()
	}

	reader := csv.NewReader(createdFile)
	reader.Comma = '\t'
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("data: %v\n", data)
	}
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	for i, row := range data {
		if i != 0 {
			createdRelationship := &database.Relationships{
				TabletNum:       row[0],
				RelationType:    row[1],
				Subject:         row[2],
				Object:          row[3],
				Providence:      row[4],
				TimePeriod:      row[5],
				DatesReferenced: row[6],
			}
			fmt.Printf("createdRelationship: %v\n", createdRelationship)
			db.InsertRelationships(*createdRelationship)
		}
	}
	json.NewEncoder(w).Encode(contents)
}

func GetRelationships(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewInternalDB()
	if err != nil {
		log.Fatal()
	}
	relationships, err := db.GetAllRelationships(0)
	fmt.Printf("relationships: %v\n", relationships)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		log.Fatal()
	}
	json.NewEncoder(w).Encode(relationships)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/insertEntities", InsertEntities).Methods("POST")
	router.HandleFunc("/getEntities", GetEntityData).Methods("POST")
	router.HandleFunc("/runEntityExtraction", RunEntityExtraction).Methods("POST")

	router.HandleFunc("/insertRelations", InsertRelationPatterns).Methods("POST")
	router.HandleFunc("/getRelations", GetRelationPatterns).Methods("POST")
	router.HandleFunc("/runRelationExtraction", RunRelationExtraction).Methods("POST")

	router.HandleFunc("/insertRelationship", InsertRelationships).Methods("POST")
	router.HandleFunc("/getRelationships", GetRelationships).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
