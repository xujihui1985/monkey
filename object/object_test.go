package object

import "testing"

func TestHashable(t *testing.T) {
	s := String{Value:"hello"}
	k := s.HashKey()
	t.Log(k)
}
