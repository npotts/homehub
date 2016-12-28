/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package sql

import (
	_ "github.com/go-sql-driver/mysql" //mysql support
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           //postgres support
	_ "github.com/mattn/go-sqlite3" //sqlite3  support

	"github.com/npotts/homehub"
)

/*SQLBackend wraps a database and functions as a homehub.Backend*/
type SQLBackend struct {
	dialect string
	db      *sqlx.DB //database backend
}

/*Backend returns a backend and nil error if successful*/
func Backend(driver, source string) (homehub.Backend, error) {
	return New(driver, source)
}

/*New returns an intialized Brianiac or a non-nil error*/
func New(driver, source string) (*SQLBackend, error) {
	db, err := sqlx.Connect(driver, source)
	if err != nil {
		return nil, err
	}
	return &SQLBackend{dialect: driver, db: db}, nil
}

/*Register attempts to register the passed piece of data
into the database - usually this means creating a table*/
func (q *SQLBackend) Register(datam homehub.Datam) error {
	sql, err := datam.SqlCreate(q.dialect)
	if err != nil {
		return err
	}
	//convert ?'s to whatever is natively used
	sql = q.db.Rebind(sql)
	_, err = q.db.Exec(sql)
	return err
}

/*Store attempts to store the passed piece of data
into the database*/
func (q *SQLBackend) Store(datam homehub.Datam) error {
	query, args, err := datam.NamedExec()
	if err != nil {
		return err
	}
	//convert ? -> whatever is natively used
	query = q.db.Rebind(query)
	_, err = q.db.NamedExec(query, args)
	return err
}

/*Stop shuts down the database*/
func (q *SQLBackend) Stop() {
	q.db.Close()
}
