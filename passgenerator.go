package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql" // Underscore: importeer alleen 'database/sql'
)

/*
Dit is een struct (verzameling van velden) die de instellingen voor de wachtwoordgenerator bijhoudt.
*/
type PasswordGenerator struct {
	length  int
	letters string
	numbers bool
	symbols bool
	db      *sql.DB
}

/*
Deze methode genereert een willekeurig wachtwoord op basis van de instellingen van de PasswordGenerator struct.
*/
func (pg *PasswordGenerator) GeneratePassword() string {
	rand.Seed(time.Now().UnixNano())

	possibleChars := pg.letters
	if pg.numbers {
		possibleChars += "0123456789"
	}
	if pg.symbols {
		possibleChars += "!@#$%^&*()_+{}:\"<>?,./;'[]\\=-`~"
	}

	password := ""
	for i := 0; i < pg.length; i++ {
		randomIndex := rand.Intn(len(possibleChars))
		password += string(possibleChars[randomIndex])
	}

	return password
}

/*
Deze methode past de lengte van het gegenereerde wachtwoord aan.
*/
func (pg *PasswordGenerator) SetPasswordLength(length int) {
	pg.length = length
}

/*
Deze methode past de tekens aan die worden gebruikt om het wachtwoord te genereren.
*/
func (pg *PasswordGenerator) SetLetters(letters string) {
	pg.letters = letters
}

/*
Deze methode stelt in of er cijfers moeten worden opgenomen in het gegenereerde wachtwoord.
*/
func (pg *PasswordGenerator) SetNumbers(includeNumbers bool) {
	pg.numbers = includeNumbers
}

/*
Deze methode stelt in of er symbolen moeten worden opgenomen in het gegenereerde wachtwoord.
*/
func (pg *PasswordGenerator) SetSymbols(includeSymbols bool) {
	pg.symbols = includeSymbols
}

/*
Deze methode maakt een verbinding met de database op basis van een databaseverbindingsreeks (DSN) en slaat de verbinding op in de PasswordGenerator struct.
*/
func (pg *PasswordGenerator) SetDatabaseConnection(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	pg.db = db

	return nil
}

/*
Deze methode controleert of het gegeven wachtwoord al in de database voorkomt door een
SQL query uit te voeren om het aantal rijen te tellen die het gegeven wachtwoord bevatten.
*/
func (pg *PasswordGenerator) PasswordExists(password string) (bool, error) {
	if pg.db == nil {
		return false, fmt.Errorf("database connection not established")
	}

	query := "SELECT COUNT(*) FROM passwords WHERE password = ?"
	row := pg.db.QueryRow(query, password)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

/*
Deze methode voegt het gegeven wachtwoord toe aan de database. Het voegt een nieuwe rij toe aan de passwords tabel met het wachtwoord en een tijdstempel.
*/
func (pg *PasswordGenerator) AddPasswordToDatabase(password string) error {
	if pg.db == nil {
		return fmt.Errorf("database connection not established")
	}

	_, err := pg.db.Exec("INSERT INTO passwords (password) VALUES (?)", password)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	pg := PasswordGenerator{ // instantie van de PasswordGenerator struct
		length:  12,
		letters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		numbers: false,
		symbols: false,
		db:      nil,
	}

	err := pg.SetDatabaseConnection("root@tcp(localhost:3306)/passgenerator")
	if err != nil {
		log.Fatal(err)
	}

	password := ""
	exists := true
	for exists {
		password = pg.GeneratePassword()          // wachtwoord generatie
		exists, err = pg.PasswordExists(password) // check
		if err != nil {
			log.Fatal(err)
		}
	}

	err = pg.AddPasswordToDatabase(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(password)
}
