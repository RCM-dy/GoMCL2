package main

type HashNotSame struct {
	Need string
	Got  string
}

func (h *HashNotSame) Error() string {
	return "hash not same\ngot: " + h.Got + "\nneed: " + h.Need
}
func NewHashNotSame(need, got string) *HashNotSame {
	return &HashNotSame{Need: need, Got: got}
}
