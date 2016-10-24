/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
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
	type x struct {
		j string
		d *Datam
		e error
		v bool
	}

	tests := map[string]x{
		"fmNull":      x{v: true, e: nil, j: `{"table": "ntable", "data": {"null": null}}`, d: &Datam{Table: "ntable", Data: map[alphabetic]Field{alphabetic("null"): Field{Value: nil, mode: fmNull}}}},
		"fmBool":      x{v: true, e: nil, j: `{"table": "btable", "data": {"bool": false}}`, d: &Datam{Table: "btable", Data: map[alphabetic]Field{alphabetic("bool"): Field{Value: false, mode: fmBool}}}},
		"fmInt":       x{v: true, e: nil, j: `{"table": "itable", "data": {"int": 1}}`, d: &Datam{Table: "itable", Data: map[alphabetic]Field{alphabetic("int"): Field{Value: int64(1), mode: fmInt}}}},
		"fmFloat":     x{v: true, e: nil, j: `{"table": "ftable", "data": {"float": 1.0}}`, d: &Datam{Table: "ftable", Data: map[alphabetic]Field{alphabetic("float"): Field{Value: 1.0, mode: fmFloat}}}},
		"fmString":    x{v: true, e: nil, j: `{"table": "stable", "data": {"string": "str"}}`, d: &Datam{Table: "stable", Data: map[alphabetic]Field{alphabetic("string"): Field{Value: "str", mode: fmString}}}},
		"shortString": x{v: true, e: nil, j: `{"table": "stable", "data": {"string": ""}}`, d: &Datam{Table: "stable", Data: map[alphabetic]Field{alphabetic("string"): Field{Value: "", mode: fmString}}}},
		//some error varieties
		"array": x{v: false, e: errFormat, j: `{"table": "bad", "data": {"array": [1,2,3]}}`},
		"obj":   x{v: false, e: errFormat, j: `{"table": "bad", "data": {"obj": {}}}`},
		"all vars": x{
			v: true,
			j: `{"table": "table", "data": {"float": 1.0, "string": "str", "int": 1, "bool": false, "null":null}}`,
			d: &Datam{
				Table: "table",
				Data: map[alphabetic]Field{
					alphabetic("float"):  Field{Value: 1.0, mode: fmFloat},
					alphabetic("string"): Field{Value: "str", mode: fmString},
					alphabetic("int"):    Field{Value: 1, mode: fmInt},
					alphabetic("bool"):   Field{Value: false, mode: fmBool},
					alphabetic("null"):   Field{Value: nil, mode: fmNull},
				},
			},
		},
	}

	for name, x := range tests {
		t.Logf("Running checks on %q", name)
		datam := &Datam{}
		if e := json.Unmarshal([]byte(x.j), datam); e != x.e {
			t.Errorf("Returned error does not match expected.  Got %v want %v", e, x.e)
		}
		if x.e != nil { //invalid JSON should return an error - no need to check the values
			continue
		}
		if datam.Valid() != x.v {
			t.Errorf("Validity does not match")
		}
		//make sure returned object matches hardcoded variable
		if !datam.Equal(x.d) {
			spew.Dump(x.d)
			spew.Dump(datam)
			t.Errorf("Reflect do not equal")
		}
	}
}
