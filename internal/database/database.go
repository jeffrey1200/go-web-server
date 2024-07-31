package database

import (
	"encoding/json"
	"errors"

	// "log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Users  map[int]User  `json:"users"`
	Chirps map[int]Chirp `json:"chirps"`
}

// func checkFileExists(filePath string) bool {
// 	_, err := os.Stat(filePath)

// 	return !errors.Is(err, os.ErrNotExist)
// }

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()

	return db, err
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		ID:    id,
		Email: email,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, nil
	}
	// fmt.Println("is the dbstructure properly loaded? ", dbStructure)
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp
	// fmt.Println("is the db chirp properly populated? ", dbStructure.Chirps)
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbstucture, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{}
	for _, c := range dbstucture.Chirps {
		if c.ID == id {
			chirp = c
		}
	}
	return chirp, nil
	// chirp := dbstucture.Chirp
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	// fmt.Printf("length of chirps? %d", len(chirps))
	return chirps, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Users:  map[int]User{},
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	// os.exist
	// os.Remove()
	// db.path = "../../database.json"
	// isFileExist := checkFileExists(db.path)
	// if isFileExist {
	// 	err := os.Remove(db.path)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// newFile, err := os.Create("")

	// }
	// _, err := os.Create(db.path)
	// if err != nil {
	// 	return err
	// }

	// _, err := os.ReadFile(db.path)
	// if err != nil {
	// 	return err
	// }
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}
	// var arrayOfChirp []Chirp
	// err = json.Unmarshal(data, &dbStructure)

	// dataStruct := dbStructure
	// err = os.WriteFile(db.path,unmarshalledData)
	return nil
}
