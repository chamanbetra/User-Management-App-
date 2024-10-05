package models

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName           string `json:"first_name" validate:"required"`
	LastName            string `json:"last_name" validate:"required"`
	Email               string `gorm:"type:varchar(255);uniqueIndex" json:"email" validate:"required,email"`
	DOB                 string `json:"dob" validate:"required"`
	Age                 int    `json:"age"`
	Password            string `json:"password"`
	VerificationToken   string
	Verified            bool
	Token_GeneratedTime time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Age = CalculateAge(u.DOB)

	if tx.Statement.Context != nil {
		if r, ok := tx.Statement.Context.Value("http_request").(*http.Request); ok {
			if r.Method == "POST" {
				if err := u.hashPassword(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CalculateAge(dobStr string) int {
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		fmt.Println("Error parsing DOB:", err)
		return 0
	}

	today := time.Now()
	age := today.Year() - dob.Year()
	if today.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
