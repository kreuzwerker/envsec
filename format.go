package envsec

import (
	"fmt"
	"io"
	"strings"

	"github.com/yawn/envmap"
)

// Formats exposes the Formatters known to envsec
var Formats = map[string]Formatter{
	"cloudformation": cloudformation,
	"shell":          shell,
	"terraform":      terraform,
}

// Formatter defines the output function for a given slice of environment key/value tuples
type Formatter func(w io.Writer, result []string)

// cloudformation emits a formatting suitable for AWS CloudFormation stacks
func cloudformation(w io.Writer, result []string) {

	var list []string

	for k, v := range envmap.ToMap(result) {

		list = append(list, fmt.Sprintf(`"%s": {
  "Default": "%s",
  "Type": "String"
}`, k, v))

	}

	fmt.Fprintln(w, strings.Join(list, ",\n"))

}

// shell emits a formatting suitable for shell exports
func shell(w io.Writer, result []string) {

	for _, e := range result {
		fmt.Fprintln(w, e)
	}

}

// terraform emits a formatting suitable for Terraform
func terraform(w io.Writer, result []string) {

	for k, v := range envmap.ToMap(result) {

		fmt.Fprintf(w, `variable "%s" {
  type    = "string"
  default = "%s"
}

`, k, v)

	}

}
