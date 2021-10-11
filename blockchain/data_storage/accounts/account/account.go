package account

import (
	"errors"
	"pandora-pay/cryptography/crypto"
	"pandora-pay/helpers"
)

type Account struct {
	helpers.SerializableInterface `json:"-"`
	PublicKey                     []byte              `json:"-"` //hashmap key
	Asset                         []byte              `json:"-"` //collection asset
	Index                         uint64              `json:"index"`
	Version                       uint64              `json:"version"`
	Balance                       *BalanceHomomorphic `json:"balance"`
}

func (account *Account) Validate() error {
	if account.Version != 0 {
		return errors.New("Version is invalid")
	}
	return nil
}

func (account *Account) GetBalance() (result *crypto.ElGamal) {
	return account.Balance.Amount
}

func (account *Account) Serialize(w *helpers.BufferWriter) {
	w.WriteUvarint(account.Version)
	w.WriteUvarint(account.Index)
	account.Balance.Serialize(w)
}

func (account *Account) SerializeToBytes() []byte {
	w := helpers.NewBufferWriter()
	account.Serialize(w)
	return w.Bytes()
}

func (account *Account) Deserialize(r *helpers.BufferReader) (err error) {

	var n uint64
	if n, err = r.ReadUvarint(); err != nil {
		return
	}
	if n != 0 {
		return errors.New("Invalid Account Version")
	}

	account.Version = n
	if account.Index, err = r.ReadUvarint(); err != nil {
		return
	}
	if err = account.Balance.Deserialize(r); err != nil {
		return
	}

	return
}

func NewAccount(publicKey []byte, asset []byte) *Account {

	var acckey crypto.Point
	if err := acckey.DecodeCompressed(publicKey); err != nil {
		panic(err)
	}
	acc := &Account{
		PublicKey: publicKey,
		Asset:     asset,
		Balance:   &BalanceHomomorphic{crypto.ConstructElGamal(acckey.G1(), crypto.ElGamal_BASE_G)},
	}

	return acc
}
