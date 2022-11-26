package factory

/*
A helper to build test factories
*/

import (
	"reflect"
	"strings"
	"time"

	"src.goblgobl.com/utils/typed"
)

type SQLStorage interface {
	MustExec(sql string, args ...any)
	Placeholder(i int) string
}

var DB SQLStorage

type Table struct {
	Truncate func() Table
	Insert   func(args ...any) typed.Typed
}

func NewTable(name string, builder func(KV) KV, pks ...string) Table {
	obj := builder(KV{})
	keys := make([]string, len(obj))
	placeholders := make([]string, len(obj))

	i := 0
	for k := range obj {
		keys[i] = k
		placeholders[i] = DB.Placeholder(i)
		i++
	}

	insertSQL := "insert into " + name + " (" + strings.Join(keys, ",") + ")"
	insertSQL += "\nvalues (" + strings.Join(placeholders, ",") + ")"
	if len(pks) > 0 {
		insertSQL += "\non conflict (" + strings.Join(pks, ",") + ") do update set "
		insertSQL += keys[0] + " = excluded." + keys[0]
		for _, k := range keys[1:] {
			insertSQL += ", " + k + " = excluded." + k
		}
	}

	deleteSQL := "delete from " + name

	t := Table{}

	t.Truncate = func() Table {
		DB.MustExec(deleteSQL)
		return t
	}

	t.Insert = func(args ...any) typed.Typed {
		obj := builder(ToKV(args))
		values := make([]any, len(obj))
		for i, k := range keys {
			values[i] = obj[k]
		}
		DB.MustExec(insertSQL, values...)
		return typed.Typed(obj)
	}

	return t
}

type KV map[string]any

func ToKV(opts []any) KV {
	args := make(KV, len(opts)/2)
	for i := 0; i < len(opts); i += 2 {
		args[opts[i].(string)] = opts[i+1]
	}
	return args
}

func (kv KV) UUID(key string, deflt ...string) any {
	if value, exists := kv[key]; exists {
		return value.(string)
	}
	if len(deflt) == 1 {
		return deflt[0]
	}
	return nil
}

func (kv KV) Int(key string, deflt ...int) any {
	if value, exists := kv[key]; exists {
		return value.(int)
	}
	if len(deflt) == 1 {
		return deflt[0]
	}
	return nil
}

func (kv KV) UInt16(key string, deflt ...uint16) any {
	if value, exists := kv[key]; exists {
		switch v := value.(type) {
		case int:
			return uint16(v)
		case uint16:
			return v
		default:
			return uint16(reflect.ValueOf(value).Uint())
		}
	}
	if len(deflt) == 1 {
		return deflt[0]
	}
	return nil
}

func (kv KV) Bool(key string, deflt ...bool) any {
	if value, exists := kv[key]; exists {
		return value.(bool)
	}
	if len(deflt) == 1 {
		return deflt[0]
	}
	return nil
}

func (kv KV) String(key string, deflt ...string) any {
	if value, exists := kv[key]; exists {
		return value.(string)
	}
	if len(deflt) == 1 {
		return deflt[0]
	}
	return nil
}

func (kv KV) Time(key string, deflt ...time.Time) any {
	if value, exists := kv[key]; exists {
		return value.(time.Time)
	}
	if len(deflt) == 1 {
		return deflt[0]
	}
	return nil
}
