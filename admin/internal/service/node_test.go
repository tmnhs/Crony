package service

import "testing"

type R struct {
	M map[int]bool
}

func TestName(t *testing.T) {
	r := R{
		M: make(map[int]bool),
	}
	r.M[1] = true
	r.M[2] = false
	r.M[3] = true
	t.Log(r.M)
}
