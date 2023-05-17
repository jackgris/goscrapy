package v1_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/business/database"
	"github.com/jackgris/goscrapy/business/dbtest"
	. "github.com/jackgris/goscrapy/cmd/api/handlers/v1"
	v1 "github.com/jackgris/goscrapy/cmd/api/handlers/v1"
	"github.com/jackgris/goscrapy/config"
)

func TestHome(t *testing.T) {

	app := fiber.New()
	app.Get("/", Home)

	response := struct{ Message string }{Message: "THIS HOME"}

	tests := []struct {
		name   string
		url    string
		method string
		status int
		resp   struct{ Message string }
	}{
		{name: "good", url: "/", method: "GET", status: http.StatusOK, resp: response},
		{name: "bad", url: "/bad", method: "GET", status: http.StatusNotFound},
		{name: "other method", url: "/", method: "POST", status: http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new http request with the route from the test case
			req := httptest.NewRequest(tt.method, tt.url, nil)

			// Perform the request plain with the app,
			// the second argument is a request latency
			// (set to -1 for no latency)
			resp, _ := app.Test(req, 1)

			if resp.StatusCode != tt.status {
				t.Errorf("Home() status should be = %v, but get status %v", tt.status, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				var r struct{ Message string }
				err := json.Unmarshal(body, &r)
				if err != nil {
					t.Errorf("Home() body should be like this { Message string } instead we have %s", r)
				}
				if r != response {
					t.Errorf("Home() body should be = %v, but get %v", r, response)
				}
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.Teardown()
	}()
	setup := config.Data{
		Dburi:  "mongodb://127.0.0.1:27017",
		Dbuser: "admin",
		Dbpass: "admin",
	}

	app := fiber.New()

	cfg := v1.Config{
		Log:   log,
		Db:    db,
		Setup: &setup,
	}

	app.Get("/", GetAllProducts(cfg))

	tests := []struct {
		name   string
		url    string
		method string
		status int
	}{
		{name: "good", url: "/", method: "GET", status: http.StatusOK},
		{name: "bad", url: "/bad", method: "GET", status: http.StatusNotFound},
		{name: "other method", url: "/", method: "POST", status: http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new http request with the route from the test case
			req := httptest.NewRequest(tt.method, tt.url, nil)

			// Perform the request plain with the app,
			// the second argument is a request latency
			// (set to -1 for no latency)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("GetAllProducts test: %s", err)
			}

			if resp.StatusCode != tt.status {
				t.Errorf("GetAllProducts status should be = %v, but get status %v", tt.status, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				products := []database.Product{}

				err := json.Unmarshal(body, &products)
				if err != nil {
					t.Errorf("GetAllProducts body should be []database.Product instead we have %s", body)
				}
				number := 1154
				if len(products) != number {
					t.Errorf("GetAllProducts number of products should be = %d, but get %d", number, len(products))
				}
			}
		})

	}
}
