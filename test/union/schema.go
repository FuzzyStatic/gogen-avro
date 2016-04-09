package avro

import "fmt"
import "io"
import "math"

type PrimitiveUnionTestRecord struct {
	UnionField UnionIntLongFloatDoubleStringBoolBytesNull
}

func (r PrimitiveUnionTestRecord) Serialize(w io.Writer) error {
	return writePrimitiveUnionTestRecord(r, w)
}

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

type UnionIntLongFloatDoubleStringBoolBytesNull struct {
	Int       int32
	Long      int64
	Float     float32
	Double    float64
	String    string
	Bool      bool
	Bytes     []byte
	Null      interface{}
	UnionType UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum
}

type UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum int

const (
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumInt    UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 0
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumLong   UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 1
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumFloat  UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 2
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumDouble UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 3
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumString UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 4
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBool   UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 5
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBytes  UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 6
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumNull   UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 7
)

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeBytes(r []byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}

func writePrimitiveUnionTestRecord(r PrimitiveUnionTestRecord, w io.Writer) error {
	var err error
	err = writeUnionIntLongFloatDoubleStringBoolBytesNull(r.UnionField, w)
	if err != nil {
		return err
	}

	return nil
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

func writeUnionIntLongFloatDoubleStringBoolBytesNull(r UnionIntLongFloatDoubleStringBoolBytesNull, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumInt:
		return writeInt(r.Int, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumLong:
		return writeLong(r.Long, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumFloat:
		return writeFloat(r.Float, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumDouble:
		return writeDouble(r.Double, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumString:
		return writeString(r.String, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBool:
		return writeBool(r.Bool, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBytes:
		return writeBytes(r.Bytes, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumNull:
		return writeNull(r.Null, w)

	}
	return fmt.Errorf("Invalid value for UnionIntLongFloatDoubleStringBoolBytesNull")
}