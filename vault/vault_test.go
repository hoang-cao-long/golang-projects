package vault

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

func TestVault(t *testing.T) {
	config := vault.DefaultConfig()

	config.Address = "http://127.0.0.1:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken("education")

	// key := "vault:v1:1lxZVBCbnWpjZ2W8HBTShanfkFZpYmEvAa9CBqmje7E+GpPSMvW3H227GhaA4E+M6gUMhM35z+oXM+Ev"
	data := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString([]byte("123456")),
	}

	date1, err := client.Logical().Write(fmt.Sprintf("transit/long/encrypt/longdeptrai"), data)
	if err != nil {
		log.Fatalf("cannot encrypt: %v", err)
	}

	log.Print(date1.Data["ciphertext"])
}
