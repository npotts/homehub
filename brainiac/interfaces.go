/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"fmt"
	"regexp"
	"strconv"
)

/*function call to both register and store data*/
type regstore func(datam Datam) error

//something from a-z and A-Z
type alphabetic string

var realph = regexp.MustCompile("^[a-zA-Z]+$")

/*Valid returns true only if it contains valid characters*/
func (a alphabetic) Valid() bool {
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
)

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
	Table alphabetic           `json:"table"`
	Data  map[alphabetic]Field `json:"data"`
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
