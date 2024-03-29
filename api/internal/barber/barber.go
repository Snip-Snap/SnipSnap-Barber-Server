package barber

import (
	"api/internal/database"
	"api/internal/methods"
	"database/sql"
	"fmt"
	"log"
)

// Barber represents a barber in a barbershop.
type Barber struct {
	BarberID    string
	ShopID      int
	UserName    string
	Password    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Gender      *string
	Dob         string
	HireDate    string
	DismissDate *string
	SeatNum     int
}

// SaveOne inserts a specified new barber into the DB. Returns err if
// it has encountered an error. Else returns nil.
func (barber Barber) SaveOne() error {
	insertBarber := "insert into barber (shopid, userName, hashedpassword," +
		"firstName, lastName, phonenumber, dob, gender, hiredate," +
		"dismissdate, seatnum) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	stmt, err := database.Db.Prepare(insertBarber)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// hashpw, err := bcrypt.GenerateFromPassword([]byte(barber.Password),
	// 	bcrypt.DefaultCost)
	hashedpw, err := methods.HashPassword(barber.Password)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(barber.ShopID, barber.UserName, hashedpw,
		barber.FirstName, barber.LastName, barber.PhoneNumber, barber.Dob,
		barber.Gender, barber.HireDate, barber.DismissDate, barber.SeatNum)
	if err != nil {
		return err
	}
	return nil

}

// GetAll selects all barbers from DB and returns them to resolver. Returns
// err if it has encountered an error. Else returns nil.
func GetAll() ([]Barber, error) {
	getAllBarbers := "select * from barber"
	stmt, err := database.Db.Prepare(getAllBarbers)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []Barber
	for rows.Next() {
		var barber Barber
		// Save directly into arguments of Scan
		err := rows.Scan(&barber.BarberID, &barber.ShopID, &barber.UserName,
			&barber.Password, &barber.FirstName, &barber.LastName,
			&barber.PhoneNumber, &barber.Gender, &barber.Dob, &barber.HireDate,
			&barber.DismissDate, &barber.SeatNum)
		if err != nil {
			return nil, err
		}
		barbers = append(barbers, barber)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return barbers, nil

}

// Get selects a specified barber via its username and modifies the receiver barber.
func (barber *Barber) Get() error {
	selectBarber := "select * from barber where username = $1"

	row := database.Db.QueryRow(selectBarber, barber.UserName)

	err := row.Scan(&barber.BarberID, &barber.ShopID, &barber.UserName,
		&barber.Password, &barber.FirstName, &barber.LastName,
		&barber.PhoneNumber, &barber.Gender, &barber.Dob,
		&barber.HireDate, &barber.DismissDate, &barber.SeatNum)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned.")
	case nil:
		barber.Dob = methods.RemoveSuffix(barber.Dob)
		barber.HireDate = methods.RemoveSuffix((barber.HireDate))
	default:
		return err
	}
	return nil
}

// GetBarberIDByUsername will return the id of the barber if it's found
// in the db via its username. Otherwise, it will return -1 when an error is
// encountered. Will return 0 if username not found in db.
func GetBarberIDByUsername(username string) (int, error) {
	getBarberID := "select barberid from barber where userName = $1"

	stmt, err := database.Db.Prepare(getBarberID)
	if err != nil {
		return -1, err
	}

	row := stmt.QueryRow(username)

	var ID int
	err = row.Scan(&ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return -1, err
		}
		return 0, err
	}

	return ID, nil
}

// Authenticate checks whether the inputBarber's raw password matches the
// Hashed password in the database's matching username.
func (barber *Barber) Authenticate() bool {
	getHashedPW := "select hashedpassword from barber WHERE username = $1"
	stmt, err := database.Db.Prepare(getHashedPW)
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(barber.UserName)

	var hashedPW string
	err = row.Scan(&hashedPW)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return methods.CheckPasswordHash(barber.Password, hashedPW)
}
