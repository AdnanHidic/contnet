package binders

import (
	"fmt"
	"github.com/guregu/null"
	"github.com/revel/revel"
	"reflect"
	"strconv"
)

var NullIntBinder = revel.Binder{
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		intValue, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return reflect.Zero(typ)
		}
		pValue := reflect.New(typ)
		nInt := null.NewInt(intValue, true)
		pValue.Elem().Set(reflect.ValueOf(nInt))
		return pValue.Elem()
	}),
	Unbind: func(output map[string]string, name string, val interface{}) {
		var v, ok = val.(null.Int)
		if ok {
			output[name] = fmt.Sprintf("%d", v.Int64)
		}
	},
}
var NullStringBinder = revel.Binder{
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		pValue := reflect.New(typ)
		nString := null.NewString(val, true)
		pValue.Elem().Set(reflect.ValueOf(nString))
		return pValue.Elem()
	}),
	Unbind: func(output map[string]string, name string, val interface{}) {
		var v, ok = val.(null.String)
		if ok {
			output[name] = fmt.Sprintf("%s", v.String)
		}
	},
}
var NullFloatBinder = revel.Binder{
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Zero(typ)
		}
		pValue := reflect.New(typ)
		nFloat := null.NewFloat(floatValue, true)
		pValue.Elem().Set(reflect.ValueOf(nFloat))
		return pValue.Elem()
	}),
	Unbind: func(output map[string]string, name string, val interface{}) {
		var v, ok = val.(null.Float)
		if ok {
			output[name] = fmt.Sprintf("%d", v.Float64)
		}
	},
}
var NullBoolBinder = revel.Binder{
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		boolValue, err := strconv.ParseBool(val)
		if err != nil {
			return reflect.Zero(typ)
		}
		pValue := reflect.New(typ)
		nBool := null.NewBool(boolValue, true)
		pValue.Elem().Set(reflect.ValueOf(nBool))
		return pValue.Elem()
	}),
	Unbind: func(output map[string]string, name string, val interface{}) {
		var v, ok = val.(null.Bool)
		if ok {
			output[name] = fmt.Sprintf("%b", v.Bool)
		}
	},
}

func AddTypeBinders() {
	revel.TypeBinders[reflect.TypeOf(null.Int{})] = NullIntBinder
	revel.TypeBinders[reflect.TypeOf(null.String{})] = NullStringBinder
	revel.TypeBinders[reflect.TypeOf(null.Float{})] = NullFloatBinder
	revel.TypeBinders[reflect.TypeOf(null.Bool{})] = NullBoolBinder
}
