package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	DOB       string `json:"dob"`
	Age       int    `json:"age"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Age = CalculateAge(u.DOB)
	return
}

func CalculateAge(dobStr string) int {
	dob, err := time.Parse("2006-01-02", dobStr) // Change format according to the input format
	if err != nil {
		fmt.Println("Error parsing DOB:", err)
		return 0 // Return 0 or handle the error as per your logic
	}

	today := time.Now()
	age := today.Year() - dob.Year()
	if today.YearDay() < dob.YearDay() {
		age--
	}
	return age
}
