package phonedb

import "database/sql"

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id;`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		// will return an invalid ID and the error
		return -1, err
	}
	return id, nil
}

// To handle setting up the database
func Migrate(driverName, dataSource string) error {
	// Connect / Open to database
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumbTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`
	_, err := db.Exec(statement)
	return err
}

func Reset(driverName, dataSource, dbName string) error {

	// Connect / Open to database
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	// If this part is commented because it has resetDB(db, dbname)
	// so that the id won't be incremented.
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}
