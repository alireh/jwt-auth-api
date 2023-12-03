package models

import (
	u "jwt-auth-api/utils"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

// a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Firstname    string `json:"firstname"`
	Lastname string `json:"lastname"`
	Token    string `json:"token";sql:"-"`
}

// Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if len(strings.TrimSpace(account.Firstname)) == 0 {
		return u.Message(false, "Firstname is required", 400), false
	}
	if len(strings.TrimSpace(account.Lastname)) == 0 {
		return u.Message(false, "Lastname is required", 400), false
	}

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required", 400), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required", 400), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry", 500), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user.", 409), false
	}

	return u.Message(false, "Requirement passed", 500), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.", 500)
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created", 200)
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found", 500)
		}
		return u.Message(false, "Connection error. Please retry", 500)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again", 500)
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In", 200)
	resp["account"] = account
	return resp
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

func GetUsers() map[string]interface{} {
	var users []Account
	db.Table("accounts").Find(&users)

	response := u.Message(true, "Get All users", 200)
	response["account"] = users
	return response
}
