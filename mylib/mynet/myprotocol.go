package mynet

import (
	"shoppingzone/myutil"
)

//Header size
const (
	ConstHeader       = "Headers"
	ConstHeaderLength = 7
	ConstMLength      = 4
)

//Protocol : the data package protocol struct
type Protocol struct{}

//EnPackage : encode the message to package with the protocol
func (p *Protocol) EnPackage(m []byte) []byte {
	return append(append([]byte(ConstHeader), myutil.Int2Byte(len(m))...), m...)
}

//DePackage : decode the massage from package with the protocol
func (p *Protocol) DePackage(b []byte) []byte {
	l := len(b)
	var i int
	data := make([]byte, 32)
	for i = 0; i < l; i++ {
		if l < i+ConstHeaderLength+ConstMLength {
			break
		}
		if string(b[i:i+ConstHeaderLength]) == ConstHeader {
			ml := myutil.Byte2Int(b[i+ConstHeaderLength : i+ConstHeaderLength+ConstMLength])
			if l < i+ConstHeaderLength+ConstMLength+ml {
				break
			}
			data = b[i+ConstHeaderLength+ConstMLength : i+ConstHeaderLength+ConstMLength+ml]
		}
	}
	if i == l {
		return make([]byte, 0)
	}
	return data
}
