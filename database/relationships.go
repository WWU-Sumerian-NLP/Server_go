package database

import (
	"database/sql"
	"log"
)

//Using examples from: https://earthly.dev/blog/golang-sqlite/

type Relationships struct {
	ID              uint   `json:"id"`
	TabletNum       string `json:"tabletNum"`
	RelationType    string `json:"relationType"`
	Subject         string `json:"subject"`
	Object          string `json:"object"`
	Providence      string `json:"providence"`
	TimePeriod      string `json:"timePeriod"`
	DatesReferenced string `json:"datesReferenced"`
}

func (i *InternalDB) InsertRelationships(relationships Relationships) (int, error) {
	res, err := i.db.Exec("INSERT INTO relationships VALUES(NULL, ?, ?, ?, ?, ?);", relationships.TabletNum, relationships.RelationType,
		relationships.Subject, relationships.Object, relationships.Providence, relationships.TimePeriod, relationships.DatesReferenced)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (i *InternalDB) RetrieveRelationship(id int) (Relationships, error) {
	log.Printf("Getting %d", id)

	//Query DB row based on ID
	row := i.db.QueryRow("SELECT id, tablet_num, relation_type, subj, obj, providence, time_period, dates_referenced FROM relations WHERE id=?", id)

	//parse row into entites struct
	relationships := Relationships{}
	var err error
	if err = row.Scan(&relationships.ID, &relationships.TabletNum, &relationships.RelationType, &relationships.Subject,
		&relationships.Object, &relationships.Providence, &relationships.TimePeriod, &relationships.DatesReferenced); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return Relationships{}, err
	}
	return relationships, err
}

func (i *InternalDB) DeleteRelationship() (int, error) {
	sqlStatement := `
	DELETE FROM relationships
	WHERE id = $1;`
	_, err := i.db.Exec(sqlStatement, 1)
	if err != nil {
		panic(err)
	}
	return 1, nil
}

func (i *InternalDB) ListRelationships(offset int) ([]Relationships, error) {
	rows, err := i.db.Query("SELECT * FROM relationships WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []Relationships{}
	for rows.Next() {
		i := Relationships{}
		err = rows.Scan(&i.ID, &i.TabletNum, &i.RelationType, &i.Subject, &i.Object, &i.Providence, &i.TimePeriod, &i.DatesReferenced)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}
