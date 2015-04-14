package configreader

import (
	"fmt"
	"github.com/revel/revel"
	"strconv"
	"time"
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

func FloatFromConfig(varName string) float64 {
	value, found := revel.Config.String(varName)
	if !found {
		panic(fmt.Sprintf("%s not found in app.conf", varName))
	}
	rt, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(fmt.Sprintf("%s in invalid format in app.conf", varName))
	}
	return rt
}

func BoolFromConfig(varName string) bool {
	value, found := revel.Config.Bool(varName)
	if !found {
		panic(fmt.Sprintf("%s not found in app.conf", varName))
	}
	return value
}

func DurationFromConfig(varName string) time.Duration {
	durationString := StringFromConfig(varName)
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		panic(fmt.Sprintf("%s supposed to be valid duration. Error: %s", varName, err.Error()))
	}
	return duration
}
