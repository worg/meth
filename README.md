METH
====

**M**​aladroit **E**​lusive **T**​ransaction **H**​ub

A wrapper around [upper.io/db](http://github.com/upper/db) to ease some repetitive tasks.


### Reasoning Behind

Working with [upper.io/db](http://github.com/upper/db) sometimes becomes a little bit repetitive, so METH is an attempt to minimize the scaffolding needed in some basic tasks [fetching by id, by certain conditions].  

## Docs
[![GoDoc](https://godoc.org/github.com/worg/meth?status.svg)](https://godoc.org/github.com/worg/meth)

## Usage

### Install

```
go get github.com/worg/meth
```

Any type must implement `Persistent` interface and return an initalized and valid db.Collection in order to use METH functions.


### Example

``` go
package person

import (
    "upper.io/db"
    "upper.io/db/some_driver"
    
)

type Person struct {
    ID      int    `db:"id"`
    Name    string `db:"name"`
    Address string `db:"address"`
}

var (
    collection db.Collection
)

func init() {
    //init DB and set collection
}

func (p *Person) Collection() db.Collection {
    return collection
}

func (p *Person) ById(id int) {
    if err := meth.One(&p); err := nil {
        // handle err
    }
    // by this point p should've been filled with db data
}

```