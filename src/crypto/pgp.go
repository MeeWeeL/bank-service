package crypto

import (
	"bytes"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

func InitPGP() (*openpgp.Entity, error) {
	entity, err := openpgp.NewEntity(
		"Bank Service",
		"PGP Encryption",
		"bank-service@example.com",
		nil,
	)
	if err != nil {
		return nil, err
	}

	for _, id := range entity.Identities {
		err := id.SelfSignature.SignUserId(
			id.UserId.Id,
			entity.PrimaryKey,
			entity.PrivateKey,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	return entity, nil
}


func EncryptPGP(data string, entity *openpgp.Entity) (string, error) {
	buf := new(bytes.Buffer)
	writer, err := armor.Encode(buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}

	plaintext, err := openpgp.Encrypt(
		writer,
		[]*openpgp.Entity{entity},
		nil,
		nil,
		nil,
	)
	if err != nil {
		return "", err
	}

	if _, err := plaintext.Write([]byte(data)); err != nil {
		return "", err
	}

	plaintext.Close()
	writer.Close()
	return buf.String(), nil
}
