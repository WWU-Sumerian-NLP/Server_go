package database

import (
	"database/sql"
	"log"
)

//Using examples from: https://earthly.dev/blog/golang-sqlite/

type Relations struct {
	ID           uint   `json:"id"`
	RelationType string `json:"relationType"`
	SubjectTag   string `json:"subjectTag"`
	ObjectTag    string `json:"objectTag"`
	RegexRules   string `json:"regexRules"`
	Tags         string `json:"tags"`
}

func (i *InternalDB) InsertRelation(relations Relations) (int, error) {
	res, err := i.db.Exec("INSERT INTO relations VALUES(NULL, ?, ?, ?, ?, ?);", relations.RelationType, relations.SubjectTag,
		relations.ObjectTag, relations.RegexRules, relations.Tags)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (i *InternalDB) RetrieveRelation(id int) (Relations, error) {
	log.Printf("Getting %d", id)

	//Query DB row based on ID
	row := i.db.QueryRow("SELECT id, relation_type, subject_tag, object_type, regex_rules, tags FROM relations WHERE id=?", id)

	//parse row into entites struct
	relations := Relations{}
	var err error
	if err = row.Scan(&relations.ID, &relations.RelationType, &relations.SubjectTag, &relations.ObjectTag,
		&relations.RegexRules, &relations.Tags); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return Relations{}, err
	}
	return relations, err
}

func (i *InternalDB) DeleteRelation() (int, error) {
	sqlStatement := `
	DELETE FROM relations
	WHERE id = $1;`
	_, err := i.db.Exec(sqlStatement, 1)
	if err != nil {
		panic(err)
	}
	return 1, nil
}

func (i *InternalDB) ListRelations(offset int) ([]Relations, error) {
	rows, err := i.db.Query("SELECT * FROM relations WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []Relations{}
	for rows.Next() {
		i := Relations{}
		err = rows.Scan(&i.ID, &i.RelationType, &i.SubjectTag, &i.ObjectTag, &i.RegexRules, &i.Tags)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}

func (i *InternalDB) GetAllRelations(offset int) ([]Relations, error) {
	rows, err := i.db.Query("SELECT * FROM relations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []Relations{}
	for rows.Next() {
		i := Relations{}
		err = rows.Scan(&i.ID, &i.RelationType, &i.SubjectTag, &i.ObjectTag, &i.RegexRules, &i.Tags)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}
