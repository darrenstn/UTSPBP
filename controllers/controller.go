package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "UTS/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetAllUsers..
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id, name, age, address, email FROM users"

	user_id := r.URL.Query()["user_id"]

	if user_id != nil {
		query += " WHERE id='" + user_id[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		var response m.Response
		response.Status = 400
		response.Message = "Error"
		json.NewEncoder(w).Encode(response)
		return
	}

	var user m.User
	var users []m.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

// InsertUserV2..
func InsertUserV2(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user := &m.User{Name: name, Age: age, Address: address, Email: email, Password: password}
	result := db.Create(user)

	if result.Error != nil {
		sendModifiedResponse(w, 400, "Insert Failed")
		return
	}

	sendModifiedResponse(w, 200, "Insert Success")
}

// UpdateUser..
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	id, _ := strconv.Atoi(r.Form.Get("id"))
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	query, _ := db.Prepare("UPDATE users SET name=?, age=?, address=?, email=?, password=? WHERE id=?")
	_, errQuery := query.Exec(name, age, address, email, password, id)

	var response m.Response
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		response.Status = 400
		response.Message = "Update user Failed!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUserV2..
func UpdateUserV2(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		sendModifiedResponse(w, 400, "Parse Form Error")
		return
	}
	id, _ := strconv.Atoi(r.Form.Get("id"))
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user := &m.User{ID: id, Name: name, Age: age, Address: address, Email: email, Password: password}
	var result *gorm.DB
	if err := db.First(&m.User{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			sendModifiedResponse(w, 400, "User not found")
			return
		}
		sendModifiedResponse(w, 400, "Error")
		return
	}

	result = db.Save(user)
	if result.Error != nil {
		sendModifiedResponse(w, 400, "Update Failed")
	}
	sendModifiedResponse(w, 200, "Update success")
}

// DeleteUser..
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	userId := vars["user_id"]

	_, errQuery := db.Exec("DELETE FROM users WHERE id=?",
		userId,
	)

	if errQuery == nil {
		sendModifiedResponse(w, 200, "Success")
	} else {
		sendModifiedResponse(w, 400, "Error")
	}
}

// Login..
func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	platform := r.Header.Get("platform")

	query := "SELECT name FROM users WHERE email='" + email + "' AND password='" + password + "'"

	rows, errQ := db.Query(query)

	if errQ != nil {
		sendModifiedResponse(w, 400, "Error")
		return
	}

	resMsg := "Success login from " + platform

	if rows.Next() {
		sendModifiedResponse(w, 200, resMsg)
	} else {
		sendModifiedResponse(w, 400, "Error")
	}
}

func sendModifiedResponse(w http.ResponseWriter, stat int, msg string) {
	var response m.Response
	response.Status = stat
	response.Message = msg
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
