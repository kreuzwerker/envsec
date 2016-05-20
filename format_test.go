package envsec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCloudformation(t *testing.T) {

	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)

	cloudformation(buf, []string{
		"foo=bar",
	})

	assert.Equal(`"foo": {
  "Default": "bar",
  "Type": "String"
}
`, buf.String())

	buf.Reset()

	cloudformation(buf, []string{
		"foo=bar",
		"wee=gee",
	})

	var r map[string]map[string]string

	err := json.Unmarshal([]byte(fmt.Sprintf("{%s}", buf.String())), &r)

	assert.NoError(err)

	assert.Equal("bar", r["foo"]["Default"])
	assert.Equal("String", r["foo"]["Type"])
	assert.Equal("gee", r["wee"]["Default"])
	assert.Equal("String", r["wee"]["Type"])

}

func TestFormatTerraform(t *testing.T) {

	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)

	terraform(buf, []string{
		"foo=bar",
	})

	assert.Equal(`variable "foo" {
  type    = "string"
  default = "bar"
}

`, buf.String())

	buf.Reset()

	terraform(buf, []string{
		"foo=bar",
		"wee=gee",
	})

	assert.Contains(buf.String(), `variable "foo" {
  type    = "string"
  default = "bar"
}`)

	assert.Contains(buf.String(), `variable "wee" {
  type    = "string"
  default = "gee"
}`)

}

func TestFormatShell(t *testing.T) {

	buf := bytes.NewBuffer(nil)

	shell(buf, []string{
		"foo=bar",
		"wee=gee",
	})

	assert.Equal(t, `foo=bar
wee=gee
`, buf.String())

}
