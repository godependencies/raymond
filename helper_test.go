package raymond

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	VERBOSE = false
)

//
// Helpers
//

func barHelper(p *HelperParams) string { return "bar" }

func barSuffixHelper(p *HelperParams) string {
	str, _ := p.Param(0).(string)
	return "bar " + str
}

func echoHelper(p *HelperParams) string {
	str, _ := p.Param(0).(string)
	nb, ok := p.Param(1).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += str
	}

	return result
}

func boolHelper(p *HelperParams) string {
	b, _ := p.Param(0).(bool)
	if b {
		return "yes it is"
	}

	return "absolutely not"
}

func gnakHelper(p *HelperParams) string {
	nb, ok := p.Param(0).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += "GnAK!"
	}

	return result
}

func linkHelper(p *HelperParams) string {
	prefix, _ := p.Param(0).(string)

	return fmt.Sprintf(`<a href="%s/%s">%s</a>`, prefix, p.DataStr("url"), p.DataStr("text"))
}

func rawHelper(p *HelperParams) string {
	result := p.Block()

	for _, param := range p.Params() {
		result += Str(param)
	}

	return result
}

func formHelper(p *HelperParams) string {
	return "<form>" + p.Block() + "</form>"
}

func formCtxHelper(p *HelperParams) string {
	return "<form>" + p.BlockWith(p.Param(0)) + "</form>"
}

func listHelper(p *HelperParams) string {
	ctx := p.Param(0)

	val := reflect.ValueOf(ctx)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() > 0 {
			result := "<ul>"
			for i := 0; i < val.Len(); i++ {
				result += "<li>"
				result += p.BlockWith(val.Index(i).Interface())
				result += "</li>"
			}
			result += "</ul>"

			return result
		}
	}

	return "<p>" + p.Inverse() + "</p>"
}

func blogHelper(p *HelperParams) string {
	return "val is " + p.ParamStr(0)
}

func equalHelper(p *HelperParams) string {
	return Str(p.ParamStr(0) == p.ParamStr(1))
}

func dashHelper(p *HelperParams) string {
	return p.ParamStr(0) + "-" + p.ParamStr(1)
}

func concatHelper(p *HelperParams) string {
	return p.ParamStr(0) + p.ParamStr(1)
}

//
// Tests
//

var helperTests = []raymondTest{
	{
		"simple helper",
		`{{foo}}`,
		nil,
		map[string]Helper{"foo": barHelper},
		`bar`,
	},
	{
		"helper with literal string param",
		`{{echo "foo"}}`,
		nil,
		map[string]Helper{"echo": echoHelper},
		`foo`,
	},
	{
		"helper with identifier param",
		`{{echo foo}}`,
		map[string]interface{}{"foo": "bar"},
		map[string]Helper{"echo": echoHelper},
		`bar`,
	},
	{
		"helper with literal boolean param",
		`{{bool true}}`,
		nil,
		map[string]Helper{"bool": boolHelper},
		`yes it is`,
	},
	{
		"helper with literal boolean param",
		`{{bool false}}`,
		nil,
		map[string]Helper{"bool": boolHelper},
		`absolutely not`,
	},
	{
		"helper with literal boolean param",
		`{{gnak 5}}`,
		nil,
		map[string]Helper{"gnak": gnakHelper},
		`GnAK!GnAK!GnAK!GnAK!GnAK!`,
	},
	{
		"helper with several parameters",
		`{{echo "GnAK!" 3}}`,
		nil,
		map[string]Helper{"echo": echoHelper},
		`GnAK!GnAK!GnAK!`,
	},
	{
		"#if helper with true literal",
		`{{#if true}}YES MAN{{/if}}`,
		nil,
		nil,
		`YES MAN`,
	},
	{
		"#if helper with false literal",
		`{{#if false}}YES MAN{{/if}}`,
		nil,
		nil,
		``,
	},
	{
		"#if helper with truthy identifier",
		`{{#if ok}}YES MAN{{/if}}`,
		map[string]interface{}{"ok": true},
		nil,
		`YES MAN`,
	},
	{
		"#if helper with falsy identifier",
		`{{#if ok}}YES MAN{{/if}}`,
		map[string]interface{}{"ok": false},
		nil,
		``,
	},
	{
		"#unless helper with true literal",
		`{{#unless true}}YES MAN{{/unless}}`,
		nil,
		nil,
		``,
	},
	{
		"#unless helper with false literal",
		`{{#unless false}}YES MAN{{/unless}}`,
		nil,
		nil,
		`YES MAN`,
	},
	{
		"#unless helper with truthy identifier",
		`{{#unless ok}}YES MAN{{/unless}}`,
		map[string]interface{}{"ok": true},
		nil,
		``,
	},
	{
		"#unless helper with falsy identifier",
		`{{#unless ok}}YES MAN{{/unless}}`,
		map[string]interface{}{"ok": false},
		nil,
		`YES MAN`,
	},
}

//
// Let's go
//

func TestHelper(t *testing.T) {
	launchRaymondTests(t, helperTests)
}
