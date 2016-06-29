// Copyright (c) 2014-2016 Hiram Jerónimo Pérez https://worg.xyz
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

package meth_test

import (
	"github.com/worg/meth"
	"testing"
	"time"
	"upper.io/db"
	"upper.io/db/sqlite"
)

var (
	settings   db.Settings
	conn       db.Database
	collection db.Collection
)

func init() {
	var err error

	settings = db.Settings{
		Database: `test/example.db`, // Path to a sqlite3 database file.
	}

	if conn, err = db.Open(sqlite.Adapter, settings); err != nil {
		panic(err)
	}

	if collection, err = conn.Collection(`employees`); err != nil {
		panic(err)
	}
}

// Taken from upper.io/db samples
type Employee struct {
	ID       int       `db:"id"`
	Name     string    `db:"name"`
	LastName string    `db:"last_name"`
	Born     time.Time `db:"born"`
	JobID    int       `db:"job_id"`
}

func (b *Employee) Collection() db.Collection {
	return collection
}

func TestOne(t *testing.T) {
	b := Employee{ID: 1}

	if err := meth.One(&b); err != nil {
		t.Error(err)
	}

	if b.Name != `Jonathan` {
		t.Error(`Failed to fetch one row by id`)
	}

	t.Log(`Fetch one by id OK`)
}

func TestOneBy(t *testing.T) {
	b := Employee{Name: `Linus`}

	if err := meth.OneBy(&b, db.Cond{`name`: b.Name}); err != nil {
		t.Error(err)
	}

	if b.ID != 2 {
		t.Error(`Failed to fetch one row with conditions`)
	}

	t.Log(`Fetch one with conditions OK`)
}

func TestAll(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.All(&b, &rows); err != nil {
		t.Error(err)
	}

	if len(rows) < 4 {
		t.Error(`Failed to fetch all rows`)
	}

	t.Log(`Fetch All rows OK`)
}

func TestAllBy(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllBy(&b, &rows, db.Cond{`id <=`: 2}); err != nil {
		t.Error(err)
	}

	if len(rows) != 2 {
		t.Error(`Failed to fetch all rows by condition`)
	}

	t.Log(`fetch all matching conditions OK`)
}

func TestExists(t *testing.T) {
	b := Employee{ID: 1}

	if ok := meth.Exists(&b); !ok {
		t.Error(`Failed to check existence by id`)
	}

	t.Log(`Check existence by id OK`)
}

func TestExsistsEqual(t *testing.T) {
	b := Employee{Name: `Jon`}

	if ok := meth.Exists(&b, db.Cond{`name`: b.Name}); !ok {
		t.Error(`Failed to check existence by equality on field`)
	}

	t.Log(`Check existence by equality on field OK`)

}

func TestExsistsRange(t *testing.T) {
	tm, _ := time.Parse(`2006-01-02`, `1950-08-07`)

	b := Employee{Born: tm}

	if ok := meth.Exists(&b, db.Cond{`born >=`: tm}); !ok {
		t.Error(`Failed to check existence by comparission on field`)
	}

	t.Log(`Check existence by comparission on field OK`)

}

func TestLimit(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllOp(&b, meth.Limit(2), &rows); err != nil {
		t.Error(`Failed to set limit on result`)
	}

	if len(rows) != 2 {
		t.Error(`Failed to set limit on result`)
	}

	t.Log(`Set limit on result OK`)
}

func TestSkip(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllOp(&b, meth.Skip(2), &rows); err != nil {
		t.Error(`Failed to skip on result`)
	}

	if len(rows) < 2 || rows[0].Name != `Jon` {
		t.Error(`Failed to skip on result`)
	}

	t.Log(`skip on result OK`, rows)
}

func TestSort(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllOp(&b, meth.Sort(`-born`), &rows); err != nil {
		t.Error(`Failed to sort result`)
	}

	if rows[0].Name != `Linus` {
		t.Error(`Failed to sort result`)
	}

	t.Log(`sort result OK`, rows)
}

func TestSelect(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllOp(&b, meth.Select(`job_id`), &rows); err != nil {
		t.Error(`Failed to skip on result`)
	}

	if rows[0].Name != `` || rows[0].JobID == 0 {
		t.Error(`Failed to select on result`)
	}

	t.Log(`select on result OK`, rows)
}

func TestWhere(t *testing.T) {
	var row Employee
	b := Employee{}

	if err := meth.OneOp(&b, meth.Where(db.Cond{`id`: 1}), &row); err != nil {
		t.Error(`Failed to skip on result`)
	}

	if row.Name != `Jonathan` {
		t.Error(`Failed to apply where on result`)
	}

	t.Log(`apply where on result OK`, row)
}

func TestGroup(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllOp(&b, meth.Group(`job_id`), &rows); err != nil {
		t.Error(`Failed to group result`)
	}

	if len(rows) != 2 {
		t.Error(`Failed to apply group result`)
	}

	t.Log(`apply group result OK`, rows)
}

func TestPaginate(t *testing.T) {
	var rows []Employee
	b := Employee{}

	if err := meth.AllOp(&b, meth.Paginate(1, 1), &rows); err != nil {
		t.Error(`Failed to skip on result`)
	}

	if len(rows) < 1 || rows[0].Name != `Linus` {
		t.Error(`Failed to skip on result`)
	}

	t.Log(`paginate on result OK`, rows)
}

func TestCustomFunc(t *testing.T) {
	var rows []Employee
	b := Employee{}

	cust := func(r db.Result) {
		r.Select(db.Raw{`SUM(id) as id`})
		r.Group(`job_id`).Skip(1)
	}

	if err := meth.AllOp(&b, cust, &rows); err != nil {
		t.Error(`Failed to skip on result`)
	}

	if len(rows) < 1 || rows[0].ID != 9 {
		t.Error(`Failed to skip on result`)
	}

	t.Log(`paginate on result OK`, rows)
}
