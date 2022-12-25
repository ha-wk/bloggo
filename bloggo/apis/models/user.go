package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail" //package for email validation
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt" //USE FOR HASHING AND PROTECTING PASS AND IMP FIELDS
	//"gorm.io/gorm"               //NEWER GORM PACKAGE THEN ABOVE ONE,,MAY HV DECLARATION ERR
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;unique" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// hashing password with default cost of 10,if not we return error
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// comapring the hashed password with original one,if doesn't matches we return err
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// storing hashed password in hashedPassword var
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// validating user,if empty return error,3 cases \\update,login,default
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")

		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

// after having user details,storing it in DB
// db is name of database,gorm.DB is the package providing it.
// SHOWING ERROR NOW,SOME PROBLEM WITH Create() in DB...may be rectified after db creation
func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil

}

// all users insertion in database!!!
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{} //storing all users in user array
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserById(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		// handle record not found
		return &User{}, errors.New("User Not found")
	}
	/*	if gorm.IsRecordNotFoundError(err) {
		return &User{},errors.New("User Not found")
	 }*/
	return u, err
}

// TAKES ID OF USER IT WAS UPDATING THE EXIXTING CREDS WITH GIVEN CREDS....
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	//to hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Model(&User{}).Where("id=?", uid).Updates(User{Password: u.Password, Nickname: u.Nickname, Email: u.Email, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &User{}, err
	}

	/*db = db.Debug().Model(&User{}).Where("id=?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)*/
	if db.Error != nil {
		return &User{}, db.Error
	}
	err = db.Debug().Model(&User{}).Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// TAKES ID OF USER AND DELETE ITS CREDS...
func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id=?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
