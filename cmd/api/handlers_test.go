package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	api "github.com/jackgris/goscrapy/cmd/api"
	"github.com/jackgris/goscrapy/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("Handlers", func() {

	app := fiber.New()

	Context("GetProductById", func() {

		It("return empty", func() {

			app.Get("/:id", api.GetProductById)
			req := httptest.NewRequest(http.MethodGet, "/1234", nil)

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())

			buf := new(bytes.Buffer)

			_, err = buf.ReadFrom(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			err = resp.Body.Close()
			Expect(err).NotTo(HaveOccurred())

			content := buf.String()
			badID := "{\"Message\":\"BAD ID\"}"
			Expect(content).To(Equal(badID))
		})

		It("return one product", func() {
			id := "63d06596eaeba53b3bd44bbb"
			app.Get("/:id", api.GetProductById)
			req := httptest.NewRequest(http.MethodGet, "/"+id, nil)

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()

			product := database.Product{}
			err = json.NewDecoder(resp.Body).Decode(&product)
			Expect(err).NotTo(HaveOccurred())

			oId, err := primitive.ObjectIDFromHex(id)
			Expect(err).NotTo(HaveOccurred())

			p := database.Product{Id_: oId}
			result := database.Db.ReadByMongoId(p)
			Expect(product).To(Equal(result))
		})

	})

})
