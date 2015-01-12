package meth // import "github.com/worg/meth"
import (
	"upper.io/db"
)

// The next functions are basically alias for their upper counterparts
// the documentation is mostly stolen from upper

// Limit defines the maximum number of results
// returns a resFunc
func Limit(u uint) resFunc {
	return func(r db.Result) {
		r.Limit(u)
	}
}

// Skip ignores the first *n* results
// returns a resFunc
func Skip(n uint) resFunc {
	return func(r db.Result) {
		r.Skip(n)
	}
}

// Sort receives field names that define the order in which elements will
// be returned in a query, field names may be prefixed with a minus sign (-)
// indicating descending order; ascending order would be used by default.
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

// AllOp works like All but applies an operation in form of:
//     func(db.Result)
// to work with the result set [for things like select, group, limit]
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
