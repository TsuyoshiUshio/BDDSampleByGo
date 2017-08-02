package cmd_test

import (
	. "."
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyVaultClient", func() {
	Context("When I request a secret", func() {
		client := KeyVaultClient{}
		It("returns the secret value", func() {
			Expect(client.GetSecretValue()).To(Equal("password"))
		})
	})
})
