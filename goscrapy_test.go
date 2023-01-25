package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Goscrapy", func() {

	Context("Getting all products", func() {

		It("should return the number of product", func() {
			products := db.ReadByWholesalers("acabajo")
			Expect(len(products)).To(Equal(349))
		})
	})
})
