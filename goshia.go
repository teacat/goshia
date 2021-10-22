package goshia

import (
	"reflect"
	"time"

	"github.com/teacat/rushia/v2"
	"gorm.io/gorm"
)

// Goshia
type Goshia struct {
	Gorm *gorm.DB
}

// NewGoshia
func NewGoshia(g *gorm.DB) *Goshia {
	return &Goshia{
		Gorm: g,
	}
}

// Query
func (g *Goshia) Query(q *rushia.Query, dest interface{}) error {
	query, params := rushia.Build(q)
	return g.Gorm.Raw(query, params...).Scan(dest).Error
}

// QueryCount
func (g *Goshia) QueryCount(q *rushia.Query, dest interface{}) (count int, err error) {
	countQuery, countParams := rushia.Build(q.Copy().ClearLimit())
	if err := g.Gorm.Raw(countQuery, countParams...).Scan(&count).Error; err != nil {
		return 0, err
	}
	query, params := rushia.Build(q)
	return count, g.Gorm.Raw(query, params...).Scan(dest).Error
}

// Exec
func (g *Goshia) Exec(q *rushia.Query) error {
	query, params := rushia.Build(q)
	return g.Gorm.Exec(query, params...).Error
}

// ExecID
func (g *Goshia) ExecID(q *rushia.Query) (id int, err error) {
	query, params := rushia.Build(q)

	err = g.Gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(query, params...).Error; err != nil {
			return err
		}
		if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
			return err
		}
		return nil
	}, nil)
	return
}

// Interfaces
func Interfaces(v interface{}) []interface{} {
	valuesVal := reflect.ValueOf(v)
	interfaceVals := make([]interface{}, 0)
	for i := range interfaceVals {
		interfaceVals[i] = valuesVal.Index(i).Interface()
	}
	return interfaceVals
}

// String
func String(v string) *string {
	return &v
}

// Bool
func Bool(v bool) *bool {
	return &v
}

// Int
func Int(v int) *int {
	return &v
}

// Float64
func Float64(v float64) *float64 {
	return &v
}

// Time
func Time(v time.Time) *time.Time {
	return &v
}
