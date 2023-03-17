package entity

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserHandler(t *testing.T) {

	t.Run("TestPostBadRequest", func(t *testing.T) {

		request := httptest.NewRequest("POST", "/v1/users", nil)
		response := httptest.NewRecorder()

		userService := TestingUserRepository{
			Users:  make(map[string]User),
			Errors: make(map[string]error),
		}

		handler := NewUserHandler(userService)
		handler.ServeHTTP(response, request)

		expected := http.StatusBadRequest
		actual := response.Result().StatusCode
		if expected != actual {
			t.Fatalf("Expected status code to be %d, not %d", expected, actual)
		}

	})

	t.Run("TestPostIsOk", func(t *testing.T) {

		body := `{"username":"johndoe","email":"johndoe@anchorly.co","password":"password12345!."}`

		request := httptest.NewRequest("POST", "/v1/users", strings.NewReader(body))
		response := httptest.NewRecorder()

		userService := TestingUserRepository{
			Users:  make(map[string]User),
			Errors: make(map[string]error),
		}

		handler := NewUserHandler(userService)
		handler.ServeHTTP(response, request)

		expected := http.StatusOK
		actual := response.Result().StatusCode
		if expected != actual {
			t.Fatalf("Expected status code to be %d, not %d", expected, actual)
		}

	})

	t.Run("TestPostInternalServerError", func(t *testing.T) {

		body := `{"username":"johndoe","email":"johndoe@anchorly.co","password":"password12345!."}`

		request := httptest.NewRequest("POST", "/v1/users", strings.NewReader(body))
		response := httptest.NewRecorder()

		userService := TestingUserRepository{
			Users: make(map[string]User),
			Errors: map[string]error{
				"create": errors.New("forcing an error"),
			},
		}

		handler := NewUserHandler(userService)
		handler.ServeHTTP(response, request)

		expected := http.StatusInternalServerError
		actual := response.Result().StatusCode
		if expected != actual {
			t.Fatalf("Expected status code to be %d, not %d", expected, actual)
		}

	})

	t.Run("TestGetByIdIsOk", func(t *testing.T) {

		userService := TestingUserRepository{
			Users: make(map[string]User),
		}

		user, err := userService.Create(User{
			Username: "jeandoe",
			Email:    "jeandoe@anchorly.co",
			Password: "password1234!.",
		})
		if err != nil {
			t.Fatal(err)
		}

		arguments := []map[string]any{
			{
				"target": "/v1/users/{id}",
				"vars":   map[string]string{"id": string(user.ID)},
			},
			{
				"target": "/v1/users/{email}",
				"vars":   map[string]string{"email": user.Email},
			},
		}

		for _, parameters := range arguments {

			request := httptest.NewRequest("GET", parameters["target"].(string), nil)
			request = mux.SetURLVars(request, parameters["vars"].(map[string]string))

			response := httptest.NewRecorder()

			handler := NewUserHandler(userService)
			handler.ServeHTTP(response, request)

			if http.StatusOK != response.Result().StatusCode {
				t.Fatalf("Expected status code to be %d, not %d", http.StatusOK, response.Result().StatusCode)
			}

			bytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			body := string(bytes)
			if strings.Contains(body, "\"password\"") {
				t.Fatal("Unexpected field in response")
			}

		}

	})

	t.Run("TestUnsupportedGet", func(t *testing.T) {

		userService := TestingUserRepository{
			Users: make(map[string]User),
		}

		request := httptest.NewRequest("GET", "/v1/users/{username}", nil)
		request = mux.SetURLVars(request, map[string]string{"username": "johndoe"})

		response := httptest.NewRecorder()

		handler := NewUserHandler(userService)
		handler.ServeHTTP(response, request)

		if http.StatusMethodNotAllowed != response.Result().StatusCode {
			t.Fatalf("Expected status code to be %d, not %d", http.StatusMethodNotAllowed, response.Result().StatusCode)
		}

	})

	t.Run("TestMethodNotSupported", func(t *testing.T) {

		request := httptest.NewRequest("DELETE", "/v1/users/{id}", nil)
		response := httptest.NewRecorder()

		userService := TestingUserRepository{
			Users: make(map[string]User),
		}

		handler := NewUserHandler(userService)
		handler.ServeHTTP(response, request)

		expected := http.StatusMethodNotAllowed
		actual := response.Result().StatusCode
		if expected != actual {
			t.Fatalf("Expected status code to be %d, not %d", expected, actual)
		}

	})

}
