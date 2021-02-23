package addresses

import (
	"errors"
	"pandora-pay/config"
	"pandora-pay/crypto"
	"pandora-pay/crypto/ecdsa"
	"pandora-pay/helpers"
)

type PrivateKey struct {
	Key []byte
}

func (pk *PrivateKey) GeneratePublicKey() (publicKey []byte, err error) {

	publicKey, err = ecdsa.ComputePublicKey(pk.Key)
	if err != nil {
		return
	}

	return
}

func (pk *PrivateKey) GenerateAddress(usePublicKeyHash bool, amount uint64, paymentID []byte) (*Address, error) {

	publicKey, err := ecdsa.ComputePublicKey(pk.Key)
	if err != nil {
		return nil, errors.New("Strange error. Your private key was invalid")
	}
	if len(paymentID) != 0 && len(paymentID) != 8 {
		return nil, errors.New("Your payment ID is invalid")
	}

	var finalPublicKey []byte

	var version AddressVersion

	if usePublicKeyHash {
		finalPublicKey = crypto.RIPEMD(publicKey)
		version = AddressVersionTransparentPublicKeyHash
	} else {
		finalPublicKey = publicKey
		version = AddressVersionTransparentPublicKey
	}

	return &Address{Network: config.NETWORK_SELECTED, Version: version, PublicKey: finalPublicKey[:], Amount: amount, PaymentID: paymentID}, nil
}

func (pk *PrivateKey) Sign(message *crypto.Hash) ([]byte, error) {
	if len(message) != 32 {
		return nil, errors.New("Message must be a hash")
	}
	privateKey, err := ecdsa.ToECDSA(pk.Key)
	if err != nil {
		return nil, err
	}
	return ecdsa.Sign(message[:], privateKey)
}

func GenerateNewPrivateKey() *PrivateKey {
	return &PrivateKey{Key: helpers.RandomBytes(32)}
}
