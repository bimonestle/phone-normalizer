package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bimonestle/go-exercise-projects/08.Phone-Number-Normalizer/phone/phonedb"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "rizalbimanto"
	password = "your-password"
	dbname   = "phone_normalizer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	reset := phonedb.Reset("postgres", psqlInfo, dbname)
	if reset != nil {
		log.Fatal(reset)
	}

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	reset = phonedb.Migrate("postgres", psqlInfo)
	if reset != nil {
		log.Fatal(reset)
	}

	db, err := phonedb.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(nil, err)
	}
	defer db.Close()

	if err := db.Seed(); err != nil {
		log.Fatal(err)
	}

	phones, err := db.AllPhones()
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range phones {
		fmt.Printf("Working on...%+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			if err != nil {
				log.Fatal(err)
			}
			if existing != nil {
				db.DeletePhone(p.ID)
			} else {
				db.UpdatePhone(&p)
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

// To format the inputted phone number
// into numbers format only using regexp. No other characters e.g: (),-,_
func normalize(phone string) string {
	// re := regexp.MustCompile("[^0-9]") // We only want characters between 0 & 9
	re := regexp.MustCompile("[\\D]") // Match any non digits
	return re.ReplaceAllString(phone, "")
}

// To format the inputted phone number
// into numbers format only. No other characters e.g: (),-,_
// func normalize(phone string) string {
// 	// The Correct format - 0123456789
// 	// It contains numbers only

// 	var buf bytes.Buffer

// 	// When string is iterated individually, it will output rune. Not string
// 	for _, ch := range phone {
// 		// If the string contains between these runes
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch) // write rune into the Buffer
// 		}
// 	}
// 	return buf.String() // convert it back into string
// }
