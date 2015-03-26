package configreader

import (
	"fmt"
	"github.com/revel/revel"
)

func StringFromConfig(varName string) string {
	value, found := revel.Config.String(varName)
	if !found {
		panic(fmt.Sprintf("%s not found in app.conf", varName))
	}
	return value
}

func IntFromConfig(varName string) int {
	value, found := revel.Config.Int(varName)
	if !found {
		panic(fmt.Sprintf("%s not found in app.conf", varName))
	}
	return value
}

func BoolFromConfig(varName string) bool {
	value, found := revel.Config.Bool(varName)
	if !found {
		panic(fmt.Sprintf("%s not found in app.conf", varName))
	}
	return value
}
