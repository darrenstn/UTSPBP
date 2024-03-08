package controllers

import (
	"encoding/json"
	"net/http"

	m "UTS/models"
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

func InsertRoom(w http.ResponseWriter, r *http.Request) {
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

	var maxPlayer int
	var countPlayer int

	if rows.Next() {
		if err := rows.Scan(&maxPlayer); err != nil {
			sendModifiedResponse(w, 400, "Query Error")
			return
		}
	}

	query = "Select COUNT(id) FROM participants where id_room =" + room_id[0]

	rows, err = db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	if rows.Next() {
		if err := rows.Scan(&countPlayer); err != nil {
			sendModifiedResponse(w, 400, "Query Error")
			return
		}
	}

	if countPlayer >= maxPlayer {
		sendModifiedResponse(w, 400, "Player Full")
		return
	}

	_, errQuery := db.Exec("INSERT INTO participants(id_room, id_account) values(?, ?)",
		room_id,
		acc_id,
	)

	var response m.Response
	if errQuery == nil {
		sendModifiedResponse(w, 200, "Insert Success")
	} else {
		sendModifiedResponse(w, 200, "Insert Failed")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func sendModifiedResponse(w http.ResponseWriter, stat int, msg string) {
	var response m.Response
	response.Status = stat
	response.Message = msg
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
