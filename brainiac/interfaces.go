/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"regexp"
)

/*function call to both register and store data*/
type regstore *func(table string, datam Datam) error

//something from a-z and A-Z
type alphabetic string

var realph = regexp.MustCompile("[a-zA-Z]+") //at least one char

/*Valid returns true only if it contains valid characters*/
func (a alphabetic) Valid() bool {
	return realph.MatchString(string(a))
}

type fieldmode int

const (
	fmInvalid fieldmode = iota
	fmNull
	fmInt
	fmFloat
	fmBool
	fmString
)

type Field struct {
	raw  []byte
	mode fieldmode
}

func (f *Field) UnmarshalJSON(incoming []byte) error {
	f.raw, f.mode = incoming, fmInvalid
	//
}

/*Datam is what all insertable things should map to*/
type Datam struct {
	Table alphabetic                 `json:"table"`
	Data  map[alphabetic]interface{} `json:"data"`
}

func (d Datam) Valid() bool {
	ok := d.Table.Valid()
	for label := range d.Data {
		ok = ok && label.Valid()
	}
	return ok
}

/*listener listen on something and call registerFxn when*/
type stoppable interface {
	stop()
}
