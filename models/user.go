package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	DOB       time.Time `json:"dob"`
	Age       int       `json:"age" gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Age = CalculateAge(u.DOB)
	return
}

func CalculateAge(dob time.Time) int {
	today := time.Now()
	age := today.Year() - dob.Year()
	if today.YearDay() < dob.YearDay() {
		age--
	}
	return age
}
