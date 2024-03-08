package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "UTS/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetAllRooms..
func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id, room_name FROM rooms"

	rows, err := db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	var room m.Room
	var rooms []m.Room
	var result m.Rooms

	for rows.Next() {
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			sendModifiedResponse(w, 400, "Scan Error")
			return
		} else {
			rooms = append(rooms, room)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.RoomsResponse
	response.Status = 200
	response.Message = "Success"
	result.Rooms = rooms
	response.Data = result
	json.NewEncoder(w).Encode(response)
}

// GetDetailRoom..
func GetDetailRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	room_id := r.URL.Query()["room_id"]

	if room_id == nil {
		sendModifiedResponse(w, 400, "Room id not provided")
		return
	}

	query := "SELECT id, room_name FROM rooms WHERE id=" + room_id[0]

	rows, err := db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	var roomParticipants m.RoomParticipants

	for rows.Next() {
		if err := rows.Scan(&roomParticipants.ID, &roomParticipants.RoomName); err != nil {
			sendModifiedResponse(w, 400, "Scan Error")
			return
		}
	}

	query = "SELECT p.id, p.id_account, a.username FROM participants p JOIN accounts a ON p.id_account = a.id WHERE p.id_room=" + room_id[0]

	rows, err = db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	var participant m.Participant
	var participants []m.Participant

	for rows.Next() {
		if err := rows.Scan(&participant.ID, &participant.AccountID, &participant.Username); err != nil {
			sendModifiedResponse(w, 400, "Query Error")
			return
		} else {
			participants = append(participants, participant)
		}
	}
	roomParticipants.Participants = participants

	w.Header().Set("Content-Type", "application/json")

	var response m.RoomParticipantsResponse
	var detail_room m.DetailRoomParticipants
	response.Status = 200
	response.Message = "Success"
	detail_room.RoomsParticipants = roomParticipants
	response.Data = detail_room
	json.NewEncoder(w).Encode(response)
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendModifiedResponse(w, 400, "Parse Form Error")
		return
	}

	room_id := r.URL.Query()["room_id"]
	acc_id := r.URL.Query()["account_id"]

	if room_id == nil || acc_id == nil {
		sendModifiedResponse(w, 400, "Room id or account id not provided")
		return
	}

	query := "SELECT id FROM rooms WHERE id=" + room_id[0]

	rows, err := db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	if !rows.Next() {
		sendModifiedResponse(w, 400, "Room not found")
		return
	}

	query = "SELECT id FROM accounts WHERE id=" + acc_id[0]

	rows, err = db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	if !rows.Next() {
		sendModifiedResponse(w, 400, "Account not found")
		return
	}

	query = "select g.max_player FROM rooms r JOIN games g ON r.id_game = g.id where r.id=" + room_id[0]

	rows, err = db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	if rows.Next() {

	}

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
