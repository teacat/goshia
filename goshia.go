package goshia

import (
	"errors"
	"time"

	"github.com/teacat/rushia/v3"
	"gorm.io/gorm"
)

var (
	ErrAlreadyBegun   = errors.New("goshia: transaction has already begun")
	ErrNotTransaction = errors.New("goshia: processing non-transaction")
)

type Type int

const (
	TypeMySQL Type = iota + 1
	TypeSQLite
)

// Goshia
type Goshia struct {
	Gorm          *gorm.DB
	Type          Type
	isTransaction bool
}

//
type Transaction struct {
}

// New
func New(g *gorm.DB) *Goshia {
	typ := TypeMySQL
	if g.Dialector.Name() == "sqlite" {
		typ = TypeSQLite
	}
	return &Goshia{
		Gorm: g,
		Type: typ,
	}
}

// Query
func (g *Goshia) Query(q *rushia.Query, dest interface{}) error {
	query, params := rushia.Build(q)
	return g.Gorm.Raw(query, params...).Scan(dest).Error
}

// QueryCount
func (g *Goshia) QueryCount(q *rushia.Query, dest interface{}) (count int, err error) {
	countQuery, countParams := rushia.Build(q.Copy().ClearLimit().Select("COUNT(*)"))
	if err := g.Gorm.Raw(countQuery, countParams...).Scan(&count).Error; err != nil {
		return 0, err
	}
	query, params := rushia.Build(q)
	return count, g.Gorm.Raw(query, params...).Scan(dest).Error
}

// Exec
func (g *Goshia) Exec(q *rushia.Query) (err error) {
	query, params := rushia.Build(q)
	result := g.Gorm.Exec(query, params...)
	err = result.Error
	return
}

// ExecAffected
func (g *Goshia) ExecAffected(q *rushia.Query) (affectedRows int, err error) {
	query, params := rushia.Build(q)
	result := g.Gorm.Exec(query, params...)
	affectedRows = int(result.RowsAffected)
	err = result.Error
	return
}

// ExecID
func (g *Goshia) ExecID(q *rushia.Query) (id int, err error) {
	query, params := rushia.Build(q)

	err = g.Gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(query, params...).Error; err != nil {
			return err
		}
		switch g.Type {
		case TypeMySQL:
			if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
				return err
			}
		case TypeSQLite:
			if err := tx.Raw(`SELECT LAST_INSERT_ROWID()`).Scan(&id).Error; err != nil {
				return err
			}
		}
		return nil
	}, nil)
	return
}

// Transaction
func (g *Goshia) Transaction(handler func(tx *Goshia) error) error {
	return g.Gorm.Transaction(func(tx *gorm.DB) error {
		return handler(&Goshia{
			Gorm:          tx,
			isTransaction: true,
		})
	})
}

// Begin
func (g *Goshia) Begin() *Goshia {
	if g.isTransaction {
		panic(ErrAlreadyBegun)
	}
	return &Goshia{
		Gorm:          g.Gorm.Begin(),
		isTransaction: true,
	}
}

// Rollback
func (g *Goshia) Rollback() *Goshia {
	if !g.isTransaction {
		panic(ErrNotTransaction)
	}
	return &Goshia{
		Gorm:          g.Gorm.Rollback(),
		isTransaction: true,
	}
}

// RollbackTo
func (g *Goshia) RollbackTo(name string) *Goshia {
	if !g.isTransaction {
		panic(ErrNotTransaction)
	}
	return &Goshia{
		Gorm:          g.Gorm.RollbackTo(name),
		isTransaction: true,
	}
}

// Commit
func (g *Goshia) Commit() *Goshia {
	if !g.isTransaction {
		panic(ErrNotTransaction)
	}
	return &Goshia{
		Gorm:          g.Gorm.Commit(),
		isTransaction: false,
	}
}

// SavePoint
func (g *Goshia) SavePoint(name string) *Goshia {
	if !g.isTransaction {
		panic(ErrNotTransaction)
	}
	return &Goshia{
		Gorm:          g.Gorm.SavePoint(name),
		isTransaction: false,
	}
}

// String 會回傳一個指針字串，這用來填補 SQL 中的 Nullable 欄位。
func String(v string) *string {
	return &v
}

// Bool 會回傳一個指針布林值，這用來填補 SQL 中的 Nullable 欄位。
func Bool(v bool) *bool {
	return &v
}

// Int 會回傳一個指針正整數，這用來填補 SQL 中的 Nullable 欄位。
func Int(v int) *int {
	return &v
}

// SliceInt 會回傳一個指針切片正整數，這用來填補 SQL 中的 Nullable 欄位。
func SliceInt(v []int) *[]int {
	return &v
}

// SliceString 會回傳一個指針切片字串，這用來填補 SQL 中的 Nullable 欄位。
func SliceString(v []string) *[]string {
	return &v
}

// Float64 會回傳一個指針浮點數，這用來填補 SQL 中的 Nullable 欄位。
func Float64(v float64) *float64 {
	return &v
}

// Time 會回傳一個指針時間，這用來填補 SQL 中的 Nullable 欄位。
func Time(v time.Time) *time.Time {
	return &v
}
