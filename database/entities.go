package database

import (
	"database/sql"
	"log"
)

//Using examples from: https://earthly.dev/blog/golang-sqlite/

type Entities struct {
	ID         uint   `json:"id"`
	EntityName string `json:"entityName"`
	EntityType string `json:"entityType"`
}

func (i *InternalDB) InsertEntity(entities Entities) (int, error) {
	res, err := i.db.Exec("INSERT INTO entities VALUES(NULL, ?, ?);", entities.EntityName, entities.EntityType)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (i *InternalDB) RetrieveEntity(id int) (Entities, error) {
	log.Printf("Getting %d", id)

	//Query DB row based on ID
	row := i.db.QueryRow("SELECT id, entity_name, entity_type FROM entities WHERE id=?", id)

	//parse row into entites struct
	entities := Entities{}
	var err error
	if err = row.Scan(&entities.ID, &entities.EntityName, &entities.EntityType); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return Entities{}, err
	}
	return entities, err
}

func (i *InternalDB) DeleteEntity() (int, error) {
	sqlStatement := `
	DELETE FROM entities
	WHERE id = $1;`
	_, err := i.db.Exec(sqlStatement, 1)
	if err != nil {
		panic(err)
	}
	return 1, nil
}

func (i *InternalDB) ListEntities(offset int) ([]Entities, error) {
	rows, err := i.db.Query("SELECT * FROM entities WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []Entities{}
	for rows.Next() {
		i := Entities{}
		err = rows.Scan(&i.ID, &i.EntityName, &i.EntityType)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}
