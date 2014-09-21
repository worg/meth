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

	conn, err = db.Open(sqlite.Adapter, settings)
	if err != nil {
		panic(err)
	}

	collection, err = conn.Collection(`birthdays`)
	if err != nil {
		panic(err)
	}
}

// Taken from upper.io/db samples
type Birthday struct {
	ID   int       `db:"id"`
	Name string    `db:"name"`
	Born time.Time `db:"born"`
}

func (b *Birthday) Collection() db.Collection {
	return collection
}

func TestOne(t *testing.T) {
	b := Birthday{ID: 1}

	if err := meth.One(&b); err != nil {
		t.Error(err)
	}

	if b.Name != `Jonathan Ive` {
		t.Error(`Failed to fetch one row by id`)
	}

	t.Log(`Fetch one by id OK`)
}

func TestOneBy(t *testing.T) {
	b := Birthday{Name: `Linus Torvalds`}

	if err := meth.OneBy(&b, db.Cond{`name`: b.Name}); err != nil {
		t.Error(err)
	}

	if b.ID != 2 {
		t.Error(`Failed to fetch one row with conditions`)
	}

	t.Log(`Fetch one with conditions OK`)
}

func TestAll(t *testing.T) {
	var rows []Birthday
	b := Birthday{}

	if err := meth.All(&b, &rows); err != nil {
		t.Error(err)
	}

	if len(rows) != 3 {
		t.Error(`Failed to fetch all rows`)
	}

	t.Log(`Fetch All rows OK`)
}

func TestAllBy(t *testing.T) {
	var rows []Birthday
	b := Birthday{}

	if err := meth.AllBy(&b, &rows, db.Cond{`id <=`: 2}); err != nil {
		t.Error(err)
	}

	if len(rows) != 2 {
		t.Error(`Failed to fetch all rows by condition`)
	}

	t.Log(`fetch all matching conditions OK`)
}

func TestExists(t *testing.T) {
	b := Birthday{ID: 1}

	if ok := meth.Exists(&b); !ok {
		t.Error(`Failed to check existence by id`)
	}

	t.Log(`Check existence by id OK`)
}

func TestExsistsEqual(t *testing.T) {
	b := Birthday{Name: `Jon Hall`}

	if ok := meth.Exists(&b, db.Cond{`name`: b.Name}); !ok {
		t.Error(`Failed to check existence by equality on field`)
	}

	t.Log(`Check existence by equality on field OK`)

}

func TestExsistsRange(t *testing.T) {
	tm, _ := time.Parse(`2006-01-02`, `1950-08-07`)

	b := Birthday{Born: tm}

	if ok := meth.Exists(&b, db.Cond{`born >=`: tm}); !ok {
		t.Error(`Failed to check existence by comparission on field`)
	}

	t.Log(`Check existence by comparission on field OK`)

}
