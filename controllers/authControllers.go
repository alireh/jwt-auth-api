package controllers

import (
	"encoding/json"
	"jwt-auth-api/models"
	u "jwt-auth-api/utils"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request", 400))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request", 400))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := models.GetUsers()
	u.Respond(w, users)
}