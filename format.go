package envsec

import (
	"fmt"
	"strings"

	"github.com/yawn/envmap"
)

var Formats = map[string]Formatter{
	"cloudformation": cloudformation,
	"shell":          shell,
	"terraform":      terraform,
}

type Formatter func(result []string)

func cloudformation(result []string) {

	var list []string

	for k, v := range envmap.ToMap(result) {

		list = append(list, fmt.Sprintf(`"%s": {
  "Default": "%s",
  "Type": "String"
}`, k, v))

	}

	fmt.Println(strings.Join(list, ",\n"))

}

func shell(result []string) {

	for _, e := range result {
		fmt.Println(e)
	}

}

func terraform(result []string) {

	for k, v := range envmap.ToMap(result) {

		fmt.Printf(`variable "%s" {
  type    = "string"
  default = "%s"
}

`, k, v)

	}

}
