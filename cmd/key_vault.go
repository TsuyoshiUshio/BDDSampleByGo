package cmd

type KeyVaultClient struct {
}

func (c *KeyVaultClient) GetSecretValue() string {

	return "password"
}
