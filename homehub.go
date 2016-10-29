/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package homehub

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*Stoppable is anything that can be stopped and somehow acts
as a proxy to receieve data.*/
type Stoppable interface {
	Stop()
}

/*The RegStore function call either registers or stores a Datam somewhere*/
type RegStore func(datam Datam) error

//Alphabetic is something from a-z and A-Z
type Alphabetic string

var realph = regexp.MustCompile("^[a-zA-Z]+[0-9]?$")

/*Valid returns true only if it contains valid characters*/
func (a Alphabetic) Valid() bool {
	return realph.MatchString(string(a))
}

type fieldmode int

const (
	fmInvalid fieldmode = iota
	fmNull
	fmBool
	fmInt
	fmFloat
	fmString
	//the rest are not defined for importing via JSON, but are used internally
	fmPrimaryKey
	fmDateTime
)

var errSQLType = fmt.Errorf("Unknown SQL Type")

/*sqltype is functionally a look up table for known dialects to sql create types*/
func (f fieldmode) sqltype(dialect string) (string, error) {
	switch dialect {
	case "sqlite3":
		switch f {
		case fmBool:
			return "BOOL", nil
		case fmInt:
			return "INT", nil
		case fmFloat:
			return "FLOAT", nil
		case fmString:
			return "TEXT", nil
		case fmPrimaryKey:
			return "INTEGER PRIMARY KEY ASC ON CONFLICT REPLACE AUTOINCREMENT", nil
		case fmDateTime:
			return "DATETIME DEFAULT CURRENT_TIMESTAMP", nil
		default:
			return "", errSQLType
		}
	case "postgres":
		switch f {
		case fmBool:
			return "BOOLEAN", nil
		case fmInt:
			return "BIGINT", nil
		case fmFloat:
			return "FLOAT8", nil
		case fmString:
			return "TEXT", nil
		case fmPrimaryKey:
			return "BIGSERIAL PRIMARY KEY", nil
		case fmDateTime:
			return "TIMESTAMP WITH TIME ZONE DEFAULT (now() at time zone 'utc')", nil
		default:
			return "", errSQLType
		}
	default:
		return "", errSQLType
	}

}

/*Field is a JSON parsable*/
type Field struct {
	mode  fieldmode
	Value interface{}
}

/*Valid returns true if the field captured is valid*/
func (f Field) Valid() bool {
	return f.mode != fmInvalid
}

var (
	reNull    = regexp.MustCompile("^null$")
	reInt     = regexp.MustCompile(`^-?0|[1-9]\d*$`)
	reBool    = regexp.MustCompile(`^true|false$`)
	reNumber  = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`)
	reString  = regexp.MustCompile(`^".*"$`)
	errFormat = fmt.Errorf("Unable to convert to a Field Value")
)

/*UnmarshalJSON conforms to the json.Unmarshaller interface*/
func (f *Field) UnmarshalJSON(incoming []byte) (err error) {
	raws := string(incoming)

	if reNull.Match(incoming) { //null check
		f.mode, f.Value = fmNull, nil
		return
	}

	if reBool.Match(incoming) { //bool check
		f.mode = fmBool
		f.Value, err = strconv.ParseBool(raws)
		return
	}
	if reInt.Match(incoming) {
		f.mode = fmInt
		f.Value, err = strconv.ParseInt(raws, 10, 64)
		return
	}

	if reNumber.Match(incoming) { //floatting point
		f.mode = fmFloat
		f.Value, err = strconv.ParseFloat(raws, 64)
		return
	}

	if reString.Match(incoming) {
		f.mode = fmString
		f.Value = raws[1 : len(raws)-1]
		return
	}

	return errFormat
}

/*Datam is what all insertable things should map to*/
type Datam struct {
	Table Alphabetic           `json:"table"`
	Data  map[Alphabetic]Field `json:"data"`
}

/*Valid is true if the fields in Datam are valid*/
func (d Datam) Valid() bool {
	ok := d.Table.Valid()
	for label, value := range d.Data {
		ok = ok && label.Valid() && value.Valid()
	}
	return ok
}

/*Equal returns true if a is the same as d*/
func (d *Datam) Equal(a *Datam) bool {
	same := d.Table == a.Table
	for key, val := range d.Data {
		_, ok := a.Data[key]
		same = same && key.Valid() && val.Valid() && ok
	}
	for key, val := range a.Data {
		_, ok := d.Data[key]
		same = same && key.Valid() && val.Valid() && ok
	}
	return same
}

/*SQLCreate forms a SQL statement to store the datam into somde database. It will
prepend a primary key 'rowid' key, timestamp as createdat and any additional data.
Dialect should be one of the following:
 - "sqlite3"
 - "postgres"
Others have yet to be defined, but should be added to fieldmode.sqltype()
*/
func (d *Datam) SqlCreate(dialect string) (r string, err error) {
	pk, err1 := fmPrimaryKey.sqltype(dialect)
	date, err2 := fmDateTime.sqltype(dialect)
	if err1 != nil || err2 != nil || !d.Valid() {
		return "", fmt.Errorf("Cannot form SqlCreate")
	}

	//fetch labels and sort them
	labels := sort.StringSlice{}
	for label, val := range d.Data {
		if txpt, err := val.mode.sqltype(dialect); err == nil {
			labels = append(labels, fmt.Sprintf("%s %s", label, txpt))
		}
	}
	labels.Sort()

	r = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %q (rowid %s, created %s, %s);`, d.Table, pk, date, strings.Join(labels, ", "))
	return r, nil
}

/*NamedExec returns a SQL statement can be be fed into a sqlx.NamedExec along with a set of matching values,
and a non-nil error if it cannot form such a statement.*/
func (d *Datam) NamedExec() (r string, vals map[string]interface{}, err error) {
	if !d.Valid() || len(d.Data) == 0 {
		return "", nil, fmt.Errorf("Cannot insert invalid data")
	}
	vals = map[string]interface{}{}
	//fetch labels and sort them
	labels := []string{}
	for label, value := range d.Data {
		vals[string(label)] = value.Value
		labels = append(labels, string(label))
	}

	r = fmt.Sprintf(`INSERT INTO %q (%s) VALUES (:%s);`, d.Table, strings.Join(labels, ","), strings.Join(labels, ",:"))
	return r, vals, nil
}
