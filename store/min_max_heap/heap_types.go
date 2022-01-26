package min_max_heap

import (
	"errors"
	"pandora-pay/helpers"
	"pandora-pay/store/hash_map"
)

type HeapElement struct {
	hash_map.HashMapElementSerializableInterface
	Key   []byte
	Score float64
}

type HeapDictElement struct {
	hash_map.HashMapElementSerializableInterface
	Key   []byte
	Index uint64
}

func (self *HeapElement) SetKey(key []byte) {
	self.Key = key
}

func (self *HeapElement) Validate() error {
	if len(self.Key) == 0 || len(self.Key) > 255 {
		return errors.New("Heap Element Key length is invalid")
	}
	return nil
}

func (self *HeapElement) Serialize(w *helpers.BufferWriter) {
	w.WriteByte(byte(len(self.Key)))
	w.Write(self.Key)
	w.WriteFloat64(self.Score)
}

func (self *HeapElement) Deserialize(r *helpers.BufferReader) (err error) {
	var count byte
	if count, err = r.ReadByte(); err != nil {
		return
	}

	if self.Key, err = r.ReadBytes(int(count)); err != nil {
		return
	}
	self.Score, err = r.ReadFloat64()
	return
}

func (self *HeapDictElement) SetKey(key []byte) {
	self.Key = key
}

func (self *HeapDictElement) Validate() error {
	return nil
}

func (self *HeapDictElement) Serialize(w *helpers.BufferWriter) {
	w.WriteUvarint(self.Index)
}

func (self *HeapDictElement) Deserialize(r *helpers.BufferReader) (err error) {
	self.Index, err = r.ReadUvarint()
	return
}
