package phonedb

import "database/sql"

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
