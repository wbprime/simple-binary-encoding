package composite

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEncodeDecode(t *testing.T) {

	var s1, s2 [5]byte
	copy(s1[:], "start")
	copy(s2[:], "  end")
	p1 := Point{s1, 3.14, 1, [2]uint8{66, 77}, Truthval1.NullValue, Truthval2.T}
	p2 := Point{s2, 0.31, 2, [2]uint8{77, 88}, Truthval1.T, Truthval2.F}
	in := Composite{p1, p2}

	var cbuf = new(bytes.Buffer)
	if err := in.Encode(cbuf, binary.LittleEndian, true); err != nil {
		t.Log("Composite Encoding Error", err)
		t.Fail()
	}
	t.Log(in, " -> ", cbuf.Bytes())
	t.Log("Cap() = ", cbuf.Cap(), "Len() = \n", cbuf.Len())

	m := MessageHeader{in.SbeBlockLength(), in.SbeTemplateId(), in.SbeSchemaId(), in.SbeSchemaVersion()}
	var mbuf = new(bytes.Buffer)
	if err := m.Encode(mbuf, binary.LittleEndian); err != nil {
		t.Log("MessageHeader Encoding Error", err)
		t.Fail()
	}
	t.Log(m, " -> ", mbuf.Bytes())
	t.Log("Cap() = ", mbuf.Cap(), "Len() = \n", mbuf.Len())

	// Create a new empty MessageHeader and Composite
	m = *new(MessageHeader)
	var out Composite = *new(Composite)

	if err := m.Decode(mbuf, binary.LittleEndian, in.SbeSchemaVersion()); err != nil {
		t.Log("MessageHeader Decoding Error", err)
		t.Fail()
	}
	t.Log("MessageHeader Decodes as: ", m)
	t.Log("Cap() = ", mbuf.Cap(), "Len() = \n", mbuf.Len())

	if err := out.Decode(cbuf, binary.LittleEndian, in.SbeSchemaVersion(), in.SbeBlockLength(), true); err != nil {
		t.Log("Composite Decoding Error", err)
		t.Fail()
	}
	t.Log("Composite decodes as: ", out)
	t.Log("Cap() = ", cbuf.Cap(), "Len() = \n", cbuf.Len())

	if in != out {
		t.Logf("in != out\n%s\n%s", in, out)
		t.Fail()
	}
}
