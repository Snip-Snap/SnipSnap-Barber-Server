package barber

import (
	"graphqltest/api/internal/database"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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

func (barber Barber) SaveOne() {
	insertBarber := "insert into barber (shopid, userName, hashedpassword," +
		"firstName, lastName, phonenumber, dob, gender, hiredate," +
		"dismissdate, seatnum) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	stmt, err := database.Db.Prepare(insertBarber)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	hashpw, err := bcrypt.GenerateFromPassword([]byte(barber.Password),
		bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(barber.ShopID, barber.UserName, string(hashpw),
		barber.FirstName, barber.LastName, barber.PhoneNumber, barber.Dob,
		barber.Gender, barber.HireDate, barber.DismissDate, barber.SeatNum)
	if err != nil {
		log.Fatal(err)
	}

}

// No receivers (parameters) so no func () GetAll() []Barber
func GetAll() []Barber {
	getAllBarbers := "select * from barber"
	stmt, err := database.Db.Prepare(getAllBarbers)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
		barbers = append(barbers, barber)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return barbers

}
