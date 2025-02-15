package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
)

func PrintJSON(input interface{}) {
	b, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	_, filename, line, _ := runtime.Caller(1)
	fmt.Println(filename, line, string(b))
}

func PrintJSONSimple(input interface{}) {
	b, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
		return
	}

	_, filename, line, _ := runtime.Caller(1)
	fmt.Println(filename, line, string(b))
}

func Println(input ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	fmt.Println(filename, line, input)
}

func Error(err error) {
	_, filename, line, _ := runtime.Caller(1)
	fmt.Println(filename, line, err)
}

func JSON(input interface{}) string {
	b, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(b)
}

func JSONSimple(input interface{}) string {
	b, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(b)
}
