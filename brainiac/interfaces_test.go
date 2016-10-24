/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAlphabetic_Valid(t *testing.T) {
	tests := map[string]bool{
		"ok":        true,
		"no spaces": false,
		"1@##$":     false,
		"":          false,
		"aASFDASdasdasdASDASDASdafasd": true,
	}

	for str, valid := range tests {
		if v := alphabetic(str).Valid(); v != valid {
			t.Errorf("With %q, expected %v, got %v", str, valid, v)
		}
	}
}

func TestField_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"table": "table", "data": {"float": 1.0, "string": "str", "int": 1, "bool": false, "array": [1,2,3], "obj": {}, "null": null}}`)
	datam := &Datam{}
	json.Unmarshal(data, datam)
	fmt.Println("datam =", datam)
}
