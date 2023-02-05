package database_test

import (
	"strconv"

	"github.com/jackgris/goscrapy/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("Database access for wholesalers", func() {

	id, _ := primitive.ObjectIDFromHex("63dd8cce380a7efc7b39aea5")

	wholesaler := database.Wholesalers{
		Id:         id,
		Name:       "userExample",
		Login:      "http://example/login",
		User:       "user@example.com",
		Pass:       "1234",
		Searchpage: "http://example/products/?page=",
	}

	Context("Save and delete one wholesaler", func() {

		It("should save one", func() {
			err := db.InsertWholesaer(wholesaler)
			Expect(err).NotTo(HaveOccurred())

		})

		It("should return the whosaler saved", func() {
			w := db.GetWhosalerById(wholesaler)
			Expect(w).To(Equal(wholesaler))

			err := db.DeleteWhosaler(w)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Get all whosalers", func() {
		wholesalers := []database.Wholesalers{}
		count := 10

		It("Create a list of wholesalers", func() {
			for i := 0; i < count; i++ {
				id, _ := primitive.ObjectIDFromHex(strconv.Itoa(i + 1))
				w := wholesaler
				w.Id = id
				wholesalers = append(wholesalers, w)

				err := db.InsertWholesaer(w)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("The number of provider should be equal to the amount created", func() {
			ws := db.FindWholesalers()
			Expect(len(ws)).To(Equal(count))
		})

		It("This need to clean all without errors", func() {
			wholesalers := db.FindWholesalers()
			for _, w := range wholesalers {
				err := db.DeleteWhosaler(w)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
