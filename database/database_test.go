package database_test

import (
	"github.com/jackgris/goscrapy/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("Database access functions", func() {

	oId, err := primitive.ObjectIDFromHex("63d06463b27694842b0af52e")
	Expect(err).NotTo(HaveOccurred())
	t := database.Product{
		Id_:         oId,
		Id:          "https://mayorista.acabajo.com.ar/productos/taza-signos-blanca-con-dorado/",
		Name:        "Taza Signos Blanca con Dorado",
		Image:       "https://d2r9epyceweg5n.cloudfront.net/stores/487/927/products/img_39821-f063cf06f84dc3257d16742163260990-480-0.webp",
		Description: "TAZA LINEA 12 SIGNOS  BLANCA CON DORADO Somos Acabajo, objetos con alma. Diseño y producción argentina de productos artesanales de uso cotidiano, con más de ...",
		Price:       "0",
		Stock:       "0",
		Wholesaler:  "acabajo",
	}

	Context("Getting all products", func() {

		It("should return the number of product", func() {
			products := db.ReadByWholesalers("acabajo")
			Expect(len(products)).To(Equal(349))
		})

		It("should return 0 product", func() {
			products := db.ReadByWholesalers("testing")
			Expect(len(products)).To(Equal(0))
		})

	})

	Context("Getting one product by ID", func() {

		It("should return an empty product", func() {
			p := database.Product{Id: "1234"}
			t := database.Product{}

			product := db.ReadById(p)

			Expect(product).To(Equal(t))
		})

		It("should return a product", func() {
			p := database.Product{Id: "https://mayorista.acabajo.com.ar/productos/taza-signos-blanca-con-dorado/"}

			product := db.ReadById(p)
			Expect(product).To(Equal(t))
		})
	})

	Context("Getting one product by Mongo ID", func() {

		It("should not return a product", func() {
			badId, err := primitive.ObjectIDFromHex("64d06463b27694842b0af53e")
			Expect(err).NotTo(HaveOccurred())

			p := database.Product{Id_: badId}
			product := db.ReadByMongoId(p)
			Expect(product).NotTo(Equal(t))
		})

		It("should return a product", func() {
			p := database.Product{Id_: oId}
			product := db.ReadByMongoId(p)
			Expect(product).To(Equal(t))
		})
	})

	Context("Saved and delete one product", func() {

		newId, err := primitive.ObjectIDFromHex("64d06463b27694842b0af53e")

		tProduct := t
		tProduct.Id = "testingID"
		tProduct.Id_ = newId

		It("Create newId like mongo ObjectID", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("Saved new product on database", func() {
			err = db.Create(tProduct)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Read new product from database", func() {
			newP := db.ReadByMongoId(database.Product{Id_: newId})
			Expect(newP).To(Equal(tProduct))
		})

		It("Delete the new product created for test propose", func() {
			err = db.Delete(database.Product{Id_: newId})
			Expect(err).NotTo(HaveOccurred())
		})

		It("Check new product was deleted correctly", func() {
			new := db.ReadByMongoId(database.Product{Id_: newId})
			Expect(new).To(Equal(database.Product{}))
		})

		It("Delete not found product", func() {
			err = db.Delete(database.Product{Id_: newId})
			Expect(err).NotTo(BeNil())
		})

	})
})
