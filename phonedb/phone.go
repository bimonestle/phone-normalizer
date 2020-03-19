package phonedb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Phone represents the phone_numbers table in the DB
type Phone struct {
	ID     int
	Number string
}

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

func (db *DB) AllPhones() ([]Phone, error) {
	statement := `SELECT id, value FROM phone_numbers`
	rows, err := db.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	fmt.Println("findPhone()")
	var p Phone
	statement := `SELECT * FROM phone_numbers WHERE value=$1`
	row := db.db.QueryRow(statement, number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (db *DB) UpdatePhone(p *Phone) error {
	fmt.Println("updatePhone()")
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}

func (db *DB) DeletePhone(id int) error {
	fmt.Println("deletePhone()")
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.db.Exec(statement, id)
	return err
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

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	statement := `SELECT * FROM phone_numbers WHERE id=$1`
	row := db.QueryRow(statement, id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", nil
	}

	return number, nil
}
