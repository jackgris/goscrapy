package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/jackgris/goscrapy/cmd/api/handlers/v1"
)

func TestHome(t *testing.T) {

	app := fiber.New()
	app.Get("/", Home)

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
			resp, _ := app.Test(req, 1)

			if resp.StatusCode != tt.status {
				t.Errorf("Home() status should be = %v, but get status %v", tt.status, resp.StatusCode)
			}
		})
	}
}
