package database_test

import (
	"github.com/jackgris/goscrapy/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database access functions", func() {

	Context("Getting all products", func() {

		It("should return the number of product", func() {
			products := db.ReadByWholesalers("acabajo")
			Expect(len(products)).To(Equal(349))
		})
	})

	Context("Getting one product by ID", func() {

		It("should return any product", func() {
			p := database.Product{Id: "1234"}
			t := database.Product{}

			product := db.ReadById(p)

			Expect(product).To(Equal(t))
		})

		It("should return a product", func() {
			p := database.Product{Id: "https://mayorista.acabajo.com.ar/productos/taza-signos-blanca-con-dorado/"}
			t := database.Product{
				//_id: ObjectId("63d06463b27694842b0af52e"),
				Id:          "https://mayorista.acabajo.com.ar/productos/taza-signos-blanca-con-dorado/",
				Name:        "Taza Signos Blanca con Dorado",
				Image:       "https://d2r9epyceweg5n.cloudfront.net/stores/487/927/products/img_39821-f063cf06f84dc3257d16742163260990-480-0.webp",
				Description: "TAZA LINEA 12 SIGNOS  BLANCA CON DORADO Somos Acabajo, objetos con alma. Diseño y producción argentina de productos artesanales de uso cotidiano, con más de ...",
				Price:       "0",
				Stock:       "0",
				Wholesaler:  "acabajo",
			}

			product := db.ReadById(p)
			Expect(product).To(Equal(t))
		})
	})
})
