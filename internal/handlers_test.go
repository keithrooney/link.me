package internal

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

type testAuthenticationService struct {
	err error
}

func (s testAuthenticationService) Authenticate(token string) error {
	return s.err
}

func testHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func TestAuthenticationHandler(t *testing.T) {

	t.Run("TestBadRequest", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/v1/links", nil)
		response := httptest.NewRecorder()

		handlerFunc := withAuthentication(testHandlerFunc)
		handlerFunc(response, request)

		statusCode := response.Result().StatusCode
		if http.StatusBadRequest != statusCode {
			t.Fatalf("400 != %d", response.Result().StatusCode)
		}

		bytes, err := ioutil.ReadAll(response.Result().Body)
		if err != nil {
			t.Fatal(err)
		}
		actual := string(bytes)
		if "bad request" != actual {
			t.Fatal("Unexpected response body")
		}

	})

	t.Run("TestForbidden", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/v1/links", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Authorization", "34234jfsdf.23eadsfa3rasd.23rasd88a3rakdfa")

		handlerFunc := withAuthentication(testHandlerFunc)
		handlerFunc(response, request)

		statusCode := response.Result().StatusCode
		if http.StatusForbidden != statusCode {
			t.Fatalf("200 != %d", response.Result().StatusCode)
		}

	})

	t.Run("TestIsOk", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/v1/links", nil)
		response := httptest.NewRecorder()

		claims := jwt.MapClaims{
			"iss": "anchorly.com",
			"sub": "1",
			"aud": "1",
			"exp": time.Now().Add(time.Hour * 3).UnixMilli(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		value, err := token.SignedString(secret)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Add("Authorization", value)

		handlerFunc := withAuthentication(testHandlerFunc)
		handlerFunc(response, request)

		statusCode := response.Result().StatusCode
		if http.StatusOK != statusCode {
			t.Fatalf("200 != %d", response.Result().StatusCode)
		}

	})

}
