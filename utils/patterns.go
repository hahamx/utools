package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var (
	Ctxs    = context.Background()
	Loggers = log.New(os.Stderr, "INFO -", 13)
)

const ignoreFields = "-"

type newfieldConv struct {
	conv newconverter
	idx  int
}

type newconverter struct {
	ValueToString func(value reflect.Value) (string, bool)
	StringToValue func(value string) (reflect.Value, error)
}

type newhashConv struct {
	factory newConvFactorys
	entity  reflect.Value
}

// 从结构体到map
func (r newhashConv) ToHash() (fields map[string]string) {
	fields = make(map[string]string, len(r.factory.fields))
	for k, f := range r.factory.fields {
		ref := r.entity.Field(f.idx)
		if v, ok := f.conv.ValueToString(ref); ok {
			fields[k] = v
		}
	}
	return fields
}

// 从map到结构体
func (r newhashConv) SetFromHash(fields map[string]interface{}) error {
	for k, f := range r.factory.fields {
		v, ok := fields[k]
		if !ok {
			continue
		}
		var valuestr string
		if v == nil {
			continue
		} else {
			valuestr = v.(string)
		}
		Loggers.Printf("set f.idx:%v value:%v from hash:%v \n", f.idx, v, valuestr)
		val, err := f.conv.StringToValue(valuestr)
		if err != nil {
			return err
		}
		r.entity.Field(f.idx).Set(val)
	}
	Loggers.Printf("after entity:%#v\n", r.entity)
	return nil
}

type newConvFactorys struct {
	fields map[string]newfieldConv
}

func (f newConvFactorys) NewConverter(entity reflect.Value) newhashConv {
	return newhashConv{factory: f, entity: entity}
}

type fields struct {
	typ  reflect.Type
	name string
	idx  int
}

type schemas struct {
	key    *fields
	ver    *fields
	fields map[string]*fields
}

type HashRep[T any] struct {
	schema  schemas
	typ     reflect.Type
	factory *newConvFactorys
	prefix  string
	idx     string
}

func (r *HashRep[T]) NewEntitys() (entity *T) {
	//无需设置默认uuid 则直接返回
	var v T
	return &v
}

func (r *HashRep[T]) SetValue(ctx context.Context, entity *T, value map[string]interface{}) (err error) {
	val := reflect.ValueOf(entity).Elem()
	fmt.Printf("val:%v\n", val)
	fmt.Printf("enrity:%T, type:%v\n", entity, reflect.TypeOf(entity).Kind().String())
	err = r.factory.NewConverter(val).SetFromHash(value)
	if err != nil {
		return err
	}
	return
}

type HashRepInter[T any] interface {
	NewEntitys() (entity *T)
	SetValue(ctx context.Context, entity *T, value map[string]interface{}) (err error)
}

func NewHashRep[T any](prefix string, sm T) HashRepInter[T] {
	repo := &HashRep[T]{
		prefix: prefix,
		idx:    "hashidx:" + prefix,
		typ:    reflect.TypeOf(sm),
	}
	repo.schema = newSchemas(repo.typ)
	repo.factory = newConvFactory(repo.typ, repo.schema)
	return repo
}

func newSchemas(t reflect.Type) schemas {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("schema %q should be a struct", t))
	}

	s := schemas{fields: make(map[string]*fields, t.NumField())}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		f := parses(sf)
		if f.name == ignoreFields {
			continue
		}
		f.idx = i
		s.fields[f.name] = &f

		s.key = &f
		s.ver = &f
	}

	return s
}

func parses(f reflect.StructField) (field fields) {
	v, _ := f.Tag.Lookup("json")
	vs := strings.SplitN(v, ",", 1)
	if vs[0] == "" {
		field.name = f.Name
	} else {
		field.name = vs[0]
	}

	field.typ = f.Type
	return field
}

var converterstruct = struct {
	val   map[reflect.Kind]newconverter
	ptr   map[reflect.Kind]newconverter
	slice map[reflect.Kind]newconverter
}{
	ptr: map[reflect.Kind]newconverter{
		reflect.Int64: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.IsNil() {
					return "", false
				}
				return strconv.FormatInt(value.Elem().Int(), 10), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(&v), nil
			},
		},
		reflect.String: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.IsNil() {
					return "", false
				}
				return value.Elem().String(), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				return reflect.ValueOf(&value), nil
			},
		},
		reflect.Bool: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.IsNil() {
					return "", false
				}
				if value.Elem().Bool() {
					return "t", true
				}
				return "f", true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				b := value == "t"
				return reflect.ValueOf(&b), nil
			},
		},
	},
	val: map[reflect.Kind]newconverter{
		reflect.Int64: {
			ValueToString: func(value reflect.Value) (string, bool) {
				return strconv.FormatInt(value.Int(), 10), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(v), nil
			},
		},
		reflect.Int: {
			ValueToString: func(value reflect.Value) (string, bool) {
				return strconv.FormatInt(value.Int(), 10), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(v), nil
			},
		},
		reflect.String: {
			ValueToString: func(value reflect.Value) (string, bool) {
				return value.String(), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				return reflect.ValueOf(value), nil
			},
		},
		reflect.Bool: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.Bool() {
					return "t", true
				}
				return "f", true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				b := value == "t"
				return reflect.ValueOf(b), nil
			},
		},
	},
	slice: map[reflect.Kind]newconverter{
		reflect.Uint8: {
			ValueToString: func(value reflect.Value) (string, bool) {
				buf, ok := value.Interface().([]byte)
				if !ok {
					return "", false
				}
				return *(*string)(unsafe.Pointer(&buf)), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				buf := []byte(value)
				return reflect.ValueOf(buf), nil
			},
		},
	},
}

func newConvFactory(t reflect.Type, schema schemas) *newConvFactorys {
	factory := &newConvFactorys{fields: make(map[string]newfieldConv, len(schema.fields))}
	for name, f := range schema.fields {
		conv, ok := converterstruct.val[f.typ.Kind()]
		switch f.typ.Kind() {
		case reflect.Ptr:
			conv, ok = converterstruct.ptr[f.typ.Elem().Kind()]
		case reflect.Slice:
			conv, ok = converterstruct.slice[f.typ.Elem().Kind()]
		}
		if !ok {
			k := f.typ.Kind()
			panic(fmt.Sprintf("schema %q should not contain unsupported field type %s.", t, k))
		}
		factory.fields[name] = newfieldConv{conv: conv, idx: f.idx}
	}
	return factory
}

// 任意结构体多字段同时赋值
func ObjSets[T any](inst T, val map[string]interface{}) (T, error) {
	/*
		通用对象关系映射
		NewHashRepository并NewJSONRepository创建一个由 redis hash 或 RedisJSON 支持的 OM 存储库。
	*/
	// c := NewRueidesClient()
	// create the repo with NewHashRepository or NewJSONRepository
	repo := NewHashRep("new", inst)
	newValues := val
	exp := repo.NewEntitys()
	Loggers.Printf("before exp:%#v, name:%#v :val: %v\n", exp, reflect.ValueOf(*exp).FieldByName("Name"), val) // output 01FNH4FCXV9JTB9WTVFAAKGSYB
	err := repo.SetValue(Ctxs, exp, newValues)
	if err != nil {
		// panic(err)
		return *exp, err
	} // success
	fmt.Printf("saveinfo exp:%v\n", exp)

	Loggers.Printf("after exp:%T, name:%v \n", exp, reflect.ValueOf(*exp).FieldByName("Name").String()) // output 01FNH4FCXV9JTB9WTVFAAKGSYB

	return *exp, nil
}

// 多个同类型结构体 多字段同时赋值
func ObjsMutilSets[T any](inst T, vals []map[string]interface{}) ([]T, error) {
	/*
		通用对象映射
		NewHashRepository并NewJSONRepository创建一个由 redis hash 或 RedisJSON 支持的 OM 存储库。
	*/
	var (
		ResT []T
	)
	repo := NewHashRep("new", inst)

	for i, val := range vals {
		exp := repo.NewEntitys()
		val["id"] = fmt.Sprintf("%v", i)
		Loggers.Printf("before exp:%#v, name:%#v\n", exp, reflect.ValueOf(*exp).FieldByName("Name")) // output 01FNH4FCXV9JTB9WTVFAAKGSYB
		err := repo.SetValue(Ctxs, exp, val)
		if err != nil {
			// panic(err)
			return ResT, err
		} // success
		fmt.Printf("saveinfo exp:%v\n", exp)

		Loggers.Printf("after exp:%T, name:%v \n", exp, reflect.ValueOf(*exp).FieldByName("Name").String()) // output 01FNH4FCXV9JTB9WTVFAAKGSYB
		ResT = append(ResT, *exp)
	}

	return ResT, nil
}
