package weblib

type Bytes []byte

func (b Bytes) String() string {
	return string(b)
}
