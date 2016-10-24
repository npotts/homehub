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
)

/*function call to both register and store data*/
type regstore *func(table string, datam Datam) error

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

type Field struct {
	raw   []byte
	mode  fieldmode
	Value interface{}
}

/*Valid returns true if the field captured is valid*/
func (f Field) Valid() bool {
	return f.mode != fmInvalid
}

var (
	reNull   = regexp.MustCompile("^null$")
	reInt    = regexp.MustCompile(`^-?0|[1-9]\d*$`)
	reBool   = regexp.MustCompile(`^true|false$`)
	reNumber = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`)
	reString = regexp.MustCompile(`^".*"$`)
)

/*UnmarshalJSON conforms to the json.Unmarshaller interface*/
func (f *Field) UnmarshalJSON(incoming []byte) error {
	f.raw, f.mode = incoming, fmInvalid
	//null check
	if reNull.Match(incoming) {
		f.mode = fmNull
		f.Value = nil
		return nil
	}
	if reBool.Match(incoming) {
		f.mode = fmBool
		return nil
		// f.Value =
	}
	if reInt.Match(incoming) {
		f.mode = fmInt
		return nil
	}

	if reNumber.Match(incoming) {
		f.mode = fmFloat
		return nil
	}

	if reString.Match(incoming) {
		f.mode = fmString
		return nil
	}

	return fmt.Errorf("Unable to convert to a Field Value")
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

/*listener listen on something and call registerFxn when*/
type stoppable interface {
	stop()
}
