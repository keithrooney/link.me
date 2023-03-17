package entity

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type userHandler struct {
	userService UserService
}

func (h userHandler) create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	} else if value, err := h.userService.Create(user); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(rw).Encode(value)
		rw.WriteHeader(http.StatusOK)
	}
}

func (h userHandler) getBy(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var user User
	var err error
	if id, present := vars["id"]; present {
		user, err = h.userService.GetById(id)
	} else if email, present := vars["email"]; present {
		user, err = h.userService.GetByEmail(email)
	} else {
		err = errors.New("unsupported query")
	}
	if err != nil {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		json.NewEncoder(rw).Encode(user)
		rw.WriteHeader(http.StatusOK)
	}
}

func (h userHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if http.MethodPost == req.Method {
		h.create(rw, req)
	} else if http.MethodGet == req.Method {
		h.getBy(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func NewUserHandler(service UserService) http.Handler {
	return userHandler{
		userService: service,
	}
}
