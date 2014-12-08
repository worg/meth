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
	ErrNoID    = errors.New(`no id field found`)
	ErrNoSlice = errors.New(`slice expected`)
)

type (
	// Persistent is the interface that allows us to communicate with upper.io/db collection
	Persistent interface {
		Collection() db.Collection
	}
	// resFunc is a func type to operate over results
	resFunc func(db.Result)
)

// The next functions are basically alias for their upper counterparts
// the documentation is mostly stolen from upper

// Limit defines the maximum number of results, returns a resFunc usable on AllOp
func Limit(u uint) resFunc {
	return func(r db.Result) {
		r.Limit(u)
	}
}

// Skip ignores the first *n* results, returns a resFunc usable on AllOp
func Skip(n uint) resFunc {
	return func(r db.Result) {
		r.Skip(n)
	}
}

// Sort receives field names that define the order in which elements will
// be returned in a query, field names may be prefixed with a minus sign (-)
// indicating descending order; ascending order would be used by default.
// also returns a resFunc
func Sort(i ...interface{}) resFunc {
	return func(r db.Result) {
		r.Sort(i...)
	}
}

// Select defines specific fields to be fulfilled on results in this result
// set.
func Select(i ...interface{}) resFunc {
	return func(r db.Result) {
		r.Select(i...)
	}
}

// Where discards the initial filtering conditions and sets new ones.
func Where(i ...interface{}) resFunc {
	return func(r db.Result) {
		r.Where(i...)
	}
}

// Group is used to group results that have the same value in the same
// column or columns.
func Group(i ...interface{}) resFunc {
	return func(r db.Result) {
		r.Group(i...)
	}
}

// Paginate applies a limit and skip to a result set
func Paginate(limit, skip uint) resFunc {
	return func(r db.Result) {
		r.Limit(limit)
		r.Skip(skip)
	}
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

	if len(cond) < 1 {
		id, _ := getID(p)
		cond = []interface{}{db.Cond{`id`: id}}
	}

	res := col.Find(cond...)

	count, _ := res.Count()

	return count > 0
}

// AllOp works like All but applies an operation in form of resFunc [func(*db.Result)] to work
// with the result set [for things like select, group, limit]
func AllOp(p Persistent, operation resFunc, rows interface{}, cond ...interface{}) error {
	col := p.Collection()

	res := col.Find(cond...)

	operation(res)

	err := res.All(rows)

	return err
}

// OneOp works like AllOp but returns only one row
func OneOp(p Persistent, operation resFunc, row interface{}, cond ...interface{}) error {
	col := p.Collection()

	res := col.Find(cond...)

	operation(res)

	err := res.One(row)

	return err
}

// getID Returns the id field from a struct
func getID(p interface{}) (int64, error) {
	valP := getValue(p)
	i := util.GetStructFieldIndex(valP.Type(), `id`)
	if len(i) != 1 {
		return 0, ErrNoID
	}

	id := valP.FieldByIndex(i).Int()

	return id, nil
}

// getValue Returns a reflect.Value from an interface taking care of pointers when needed
func getValue(t interface{}) (rslt reflect.Value) {
	rslt = reflect.ValueOf(t)
	for rslt.Kind() == reflect.Ptr {
		rslt = rslt.Elem()
	}

	return
}

// isSlice Check if an interface is a slice
func isSlice(t interface{}) bool {
	if k := reflect.TypeOf(t).Kind(); k == reflect.Slice {
		return true
	}

	return false
}
