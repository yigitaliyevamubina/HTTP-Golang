package main

import (
	"backend/models"
	"backend/storage"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/user/all", GetUsers).Methods("GET")

	r.HandleFunc("/user/create", CreateUser).Methods("POST")

	r.HandleFunc("/user/{id}", GetIdByPath).Methods("GET")

	r.HandleFunc("/user/delete", DeleteUser).Methods("DELETE")

	r.HandleFunc("/user/update", UpdateUser).Methods("PUT")

	http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading body ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	if err = json.Unmarshal(content, &user); err != nil {
		log.Println("error while unmarshaling body ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	user.ID = id

	respUser, err := storage.CreateUser(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshaling response ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("error while converting page, not integer", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("error while converting limit, not integer", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := storage.GetAllUsers(intPage, intLimit)
	if err != nil {
		log.Println("error while getting all users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respUsers, err := json.Marshal(users)
	if err != nil {
		log.Println("error while marshaling all users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(respUsers)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")

	var user models.User
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(content, &user); err != nil {
		log.Println("error while unmarshaling user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatedUser, err := storage.UpdateUserById(userId, user)
	if err != nil {
		log.Println("error while updating the user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respUser, err := json.Marshal(updatedUser)
	if err != nil {
		log.Println("error while marshaling the user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")

	respUser, err := storage.DeleteUserById(userId)
	if err != nil {
		log.Println("error while deleting user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshaling the deleted user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

//id in url path

func GetIdByPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	respUser, err := storage.GetUserById(userId)
	if err != nil {
		log.Println("error while getting user by id in path", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshaling user in path", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
