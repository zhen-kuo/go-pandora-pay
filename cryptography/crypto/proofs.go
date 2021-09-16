package crypto

import (
	"math/big"
	"pandora-pay/cryptography/bn256"
	"pandora-pay/helpers"
)

type StatementPublicKeyIndex struct {
	Registered bool
	Index      uint64
}

type Statement struct {
	RingSize      uint64
	CLn           []*bn256.G1 //bloomed
	CRn           []*bn256.G1 //bloomed
	Publickeylist []*bn256.G1 //bloomed
	C             []*bn256.G1 // commitments
	D             *bn256.G1
	Fees          uint64
	Roothash      []byte // note roothash contains the merkle root hash of chain, when it was build
}

type Witness struct {
	SecretKey      *big.Int
	R              *big.Int
	TransferAmount uint64 // total value being transferred
	Balance        uint64 // whatever is the the amount left after transfer
	Index          []int  // index of sender in the public key list
}

func (s *Statement) Serialize(w *helpers.BufferWriter) {

	pow, err := GetPowerof2(len(s.C))
	if err != nil {
		panic(err)
	}
	w.WriteByte(byte(pow)) // len(s.Publickeylist) is always power of 2
	w.WriteUvarint(s.Fees)
	w.Write(s.D.EncodeCompressed())

	for i := 0; i < len(s.C); i++ {
		w.Write(s.CLn[i].EncodeCompressed())           //can be bloomed
		w.Write(s.CRn[i].EncodeCompressed())           //can be bloomed
		w.Write(s.Publickeylist[i].EncodeCompressed()) //can be bloomed
		w.Write(s.C[i].EncodeCompressed())
	}

	w.Write(s.Roothash)
}

func (s *Statement) Deserialize(r *helpers.BufferReader) (err error) {

	length, err := r.ReadByte()
	if err != nil {
		return
	}
	s.RingSize = 1 << length

	if s.Fees, err = r.ReadUvarint(); err != nil {
		return
	}

	if s.D, err = r.ReadBN256G1(); err != nil {
		return
	}

	s.CLn = make([]*bn256.G1, int(s.RingSize))
	s.CRn = make([]*bn256.G1, int(s.RingSize))
	s.Publickeylist = make([]*bn256.G1, int(s.RingSize))
	s.C = make([]*bn256.G1, int(s.RingSize))
	for i := 0; i < int(s.RingSize); i++ {

		if s.CLn[i], err = r.ReadBN256G1(); err != nil {
			return
		}
		if s.CRn[i], err = r.ReadBN256G1(); err != nil {
			return
		}
		if s.Publickeylist[i], err = r.ReadBN256G1(); err != nil {
			return
		}

		if s.C[i], err = r.ReadBN256G1(); err != nil {
			return
		}
	}

	if s.Roothash, err = r.ReadBytes(32); err != nil {
		return
	}

	return nil

}
