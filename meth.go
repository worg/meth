// Copyright (c) 2014 Hiram Jerónimo Pérez worg{at}linuxmail[dot]org
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package meth is a thin wrapper around upper.io/db to ease some repettitive tasks
package meth

import (
	"errors"
	"reflect"
	"upper.io/db"
	"upper.io/db/util"
)

// Shared error codes
var (
	ErrNoID    = errors.New(`No id field found`)
	ErrNoSlice = errors.New(`Slice expected`)
)

// Persistent is the interface that allows us to communicate with upper.io/db collection
type Persistent interface {
	Collection() db.Collection
}

// One fills p with one element based on the id field
func One(p Persistent) error {
	col := p.Collection()

	id, err := getID(p)
	if err != nil {
		return err
	}

	res := col.Find(db.Cond{`id`: id})

	if err := res.One(p); err != nil {
		return err
	}

	return nil
}

// OneBy fills p with one element matching the given conditions
func OneBy(p Persistent, cond ...interface{}) error {
	col := p.Collection()
	res := col.Find(cond...)

	if err := res.One(p); err != nil {
		return err
	}

	return nil
}

// AllBy fills rows with data from the result set matching the given conditions
func AllBy(p Persistent, rows interface{}, cond ...interface{}) error {
	col := p.Collection()
	res := col.Find(cond...)

	if err := res.All(rows); err != nil {
		return err
	}

	return nil
}

// All works as an alias for AllBy
func All(p Persistent, rows interface{}, cond ...interface{}) error {
	return AllBy(p, rows, cond...)
}

// Exists returns true if a record exists matching either id or the given conditions
func Exists(p Persistent, cond ...interface{}) bool {
	col := p.Collection()

	id, _ := getID(p)

	if len(cond) < 1 {
		cond = []interface{}{db.Cond{`id`: id}}
	}

	res := col.Find(cond...)

	count, _ := res.Count()

	return count > 0
}

func getID(p interface{}) (int64, error) {
	valP := getValue(p)
	i := util.GetStructFieldIndex(valP.Type(), `id`)
	if len(i) != 1 {
		return 0, ErrNoId
	}

	id := valP.FieldByIndex(i).Int()

	return id, nil
}

func getValue(t interface{}) (rslt reflect.Value) {
	rslt = reflect.ValueOf(t)
	for rslt.Kind() == reflect.Ptr {
		rslt = rslt.Elem()
	}

	return
}

func isSlice(t interface{}) bool {
	if k := reflect.TypeOf(t).Kind(); k == reflect.Slice {
		return true
	}

	return false
}
