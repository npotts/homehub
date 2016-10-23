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
type alphabet string

var realph = regexp.MustCompile("[a-zA-Z]+") //at least one char

/*Valid returns true only if it contains valid characters*/
func (a alphabet) Valid() bool {
	return realph.MatchString(string(a))
}

/*Datam is what all insertable things should map to*/
type Datam struct {
	Table alphabet                 `json:"table"`
	Data  map[alphabet]interface{} `json:"data"`
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
