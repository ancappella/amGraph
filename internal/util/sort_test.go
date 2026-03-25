package util

import (
	"testing"
)

func TestSort(t *testing.T) {
	cases := []struct {
	   Name string
	   A int
	   B int
	   Want int
	}{
		{"1+1",1,1,2},
		{"2+3",2,3,5},
		{"负数",-1,1,0}
	}
}
