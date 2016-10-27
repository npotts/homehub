/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package homehub

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
		if v := Alphabetic(str).Valid(); v != valid {
			t.Errorf("With %q, expected %v, got %v", str, valid, v)
		}
	}
}

func TestFieldMode_SqlType(t *testing.T) {
	ok := map[string][]fieldmode{
		"sqlite3":  []fieldmode{fmBool, fmInt, fmFloat, fmString, fmPrimaryKey, fmDateTime},
		"postgres": []fieldmode{fmBool, fmInt, fmFloat, fmString, fmPrimaryKey, fmDateTime},
	}
	errord := map[string][]fieldmode{
		"sqlite3":  []fieldmode{fmInvalid},
		"postgres": []fieldmode{fmInvalid},
		"unknown":  []fieldmode{fmInvalid},
	}

	run := func(tests map[string][]fieldmode, err error) {
		for dialect, defined := range tests {
			for _, fm := range defined {
				if _, e := fm.sqltype(dialect); e != err {
					t.Errorf("Faied testing for %q:%v:  Wanted %v, got %v", dialect, fm, err, e)
				}
			}
		}
	}
	run(ok, nil)
	run(errord, errSQLType)
}

func TestField_UnmarshalJSON(t *testing.T) {
	type x struct {
		j string
		d *Datam
		e error
		v bool
	}

	tests := map[string]x{
		"fmNull":      x{v: true, e: nil, j: `{"table": "ntable", "data": {"null": null}}`, d: &Datam{Table: "ntable", Data: map[Alphabetic]Field{Alphabetic("null"): Field{Value: nil, mode: fmNull}}}},
		"fmBool":      x{v: true, e: nil, j: `{"table": "btable", "data": {"bool": false}}`, d: &Datam{Table: "btable", Data: map[Alphabetic]Field{Alphabetic("bool"): Field{Value: false, mode: fmBool}}}},
		"fmInt":       x{v: true, e: nil, j: `{"table": "itable", "data": {"int": 1}}`, d: &Datam{Table: "itable", Data: map[Alphabetic]Field{Alphabetic("int"): Field{Value: int64(1), mode: fmInt}}}},
		"fmFloat":     x{v: true, e: nil, j: `{"table": "ftable", "data": {"float": 1.0}}`, d: &Datam{Table: "ftable", Data: map[Alphabetic]Field{Alphabetic("float"): Field{Value: 1.0, mode: fmFloat}}}},
		"fmString":    x{v: true, e: nil, j: `{"table": "stable", "data": {"string": "str"}}`, d: &Datam{Table: "stable", Data: map[Alphabetic]Field{Alphabetic("string"): Field{Value: "str", mode: fmString}}}},
		"shortString": x{v: true, e: nil, j: `{"table": "stable", "data": {"string": ""}}`, d: &Datam{Table: "stable", Data: map[Alphabetic]Field{Alphabetic("string"): Field{Value: "", mode: fmString}}}},
		//some error varieties
		"array": x{v: false, e: errFormat, j: `{"table": "bad", "data": {"array": [1,2,3]}}`},
		"obj":   x{v: false, e: errFormat, j: `{"table": "bad", "data": {"obj": {}}}`},
		"all vars": x{
			v: true,
			j: `{"table": "table", "data": {"float": 1.0, "string": "str", "int": 1, "bool": false, "null":null}}`,
			d: &Datam{
				Table: "table",
				Data: map[Alphabetic]Field{
					Alphabetic("float"):  Field{Value: 1.0, mode: fmFloat},
					Alphabetic("string"): Field{Value: "str", mode: fmString},
					Alphabetic("int"):    Field{Value: 1, mode: fmInt},
					Alphabetic("bool"):   Field{Value: false, mode: fmBool},
					Alphabetic("null"):   Field{Value: nil, mode: fmNull},
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

func TestDatam_SqlCreate(t *testing.T) {
	type x struct {
		d       Datam
		dialect string
		expect  string
		inerror bool
	}

	tests := []x{
		x{dialect: "no idea", inerror: true},
		x{d: Datam{
			Table: "test",
			Data: map[Alphabetic]Field{
				Alphabetic("float"):  Field{Value: 1.0, mode: fmFloat},
				Alphabetic("string"): Field{Value: "str", mode: fmString},
				Alphabetic("int"):    Field{Value: 1, mode: fmInt},
				Alphabetic("bool"):   Field{Value: false, mode: fmBool},
			},
		},
			dialect: "sqlite3", inerror: false,
			expect: `CREATE TABLE IF NOT EXISTS "test" (rowid INTEGER PRIMARY KEY ASC ON CONFLICT REPLACE AUTOINCREMENT, created DATETIME DEFAULT CURRENT_TIMESTAMP, bool BOOL, float FLOAT, int INT, string TEXT);`,
		},
	}

	for i, x := range tests {
		r, e := x.d.SqlCreate(x.dialect)
		t.Logf("Running test #%d: %v", i, x.d.Valid())
		if (x.inerror && e == nil) || (!x.inerror && e != nil) {
			t.Logf("Want an error: %v Got Error: %v", x.inerror, e)
			t.Fatal("Should either error and didnt, or not and did")
		}
		if r != x.expect {
			t.Logf("Got: %v", r)
			t.Logf("Wnt: %v", x.expect)
			t.Errorf("Did not get the string I expected")
		}
	}
}

func TestDatam_NamedExec(t *testing.T) {
	type x struct {
		d       Datam
		dialect string
		expect  string
		inerror bool
	}

	tests := []x{
		x{dialect: "no idea", inerror: true},
		x{d: Datam{
			Table: "test",
			Data: map[Alphabetic]Field{
				Alphabetic("float"): Field{Value: 1.0, mode: fmFloat},
			},
		},
			dialect: "sqlite3", inerror: false,
			expect: `INSERT INTO "test" (float) VALUES (:float);`,
		},
	}

	for i, x := range tests {
		r, vals, e := x.d.NamedExec()
		t.Logf("Running test #%d: %v", i, vals)
		if (x.inerror && e == nil) || (!x.inerror && e != nil) {
			t.Logf("Want an error: %v Got Error: %v", x.inerror, e)
			t.Fatal("Should either error and didnt, or not and did")
		}
		if r != x.expect {
			t.Logf("Got: %v", r)
			t.Logf("Wnt: %v", x.expect)
			t.Errorf("Did not get the string I expected")
		}
	}
}

/*









 */
