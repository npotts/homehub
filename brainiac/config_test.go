/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"testing"
)

func Test_ConfigForArgs(t *testing.T) {
	tests := map[string][]string{
		"help":    []string{"--help"},
		"regular": []string{"--help"},
	}
	for tname, args := range tests {
		t.Logf("Checking %v", tname)
		c := ConfigForArgs(args)
		t.Log(c)
	}
}
