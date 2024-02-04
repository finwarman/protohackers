package json

import (
	"fmt"
	"strings"

	participle "github.com/alecthomas/participle/v2"
)

// A custom (basic) json parser, for use in `finwarman/protohackers`
// This was an exercise just for fun! It's not fast, fully-featured, etc.

// Define JSON grammar (using participle syntax)

// JSONValue represents a generic JSON value.
type JSONValue struct {
	Object *JSONObject `@@`
	Array  *JSONArray  `| @@`
	Str    *string     `| @String`
	Number *float64    `| @Float | @Int`
	Bool   *Boolean    `| @("true" | "false")`
	Null   *string     `| @"null"`
}

// JSONObject represents a JSON object (key-value pairs).
type JSONObject struct {
	Pairs []*JSONPair `"{" [ @@ { "," @@ } ] "}"`
}

// JSONPair represents a key-value pair in a JSON object.
type JSONPair struct {
	Key   string     `@String ":"`
	Value *JSONValue `@@`
}

// JSONArray represents a JSON array.
type JSONArray struct {
	Values []*JSONValue `"[" [ @@ { "," @@ } ] "]"`
}

// JSONArray represents a JSON array.
func ParseJSON(input string) (*JSONValue, error) {
	parser, err := participle.Build[JSONValue](
		participle.Unquote("String"),
	)
	if err != nil {
		return nil, err
	}

	value, err := parser.ParseString("", input)
	return value, err
}

// Define boolean type to parse booleans - default 'bool' behaviour
// in participle only indicates that a match occured
type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

// == String Conversion Methods == //

func (j *JSONValue) String() string {
	return j.toString(0)
}

func (j *JSONPair) String() string {
	return j.toString(0)
}

func (j *JSONPair) toString(indent int) string {
	indentStr := strings.Repeat(" ", indent*2)
	result := fmt.Sprintf("%s\"%s\": %s", indentStr, j.Key, j.Value.toString(indent))
	return result
}

func (j *JSONValue) toString(indent int) string {
	indentStr := strings.Repeat(" ", indent*2)
	result := ""

	var comma = "," // separator, set to "" for last value of object/array
	if j.Object != nil {
		result += "{\n"
		for i, value := range j.Object.Pairs {
			if i == len(j.Object.Pairs)-1 {
				comma = ""
			}
			result += fmt.Sprintf("%s%s\n", value.toString(indent+1), comma)
		}
		result += indentStr + "}"
	} else if j.Array != nil {
		result += "[\n"
		for i, value := range j.Array.Values {
			if i == len(j.Array.Values)-1 {
				comma = ""
			}
			result += fmt.Sprintf("%s%d: %s%s\n", indentStr+"  ", i, value.toString(indent+1), comma)
		}
		result += indentStr + "]"
	} else if j.Str != nil {
		result += fmt.Sprintf("\"%s\"", *j.Str)
	} else if j.Number != nil {
		floatValue := *j.Number
		// Check if the float is "int-y" (has no fractional part)
		if floatValue == float64(int(floatValue)) {
			result += fmt.Sprintf("%d", int(floatValue)) // Format as an integer
		} else {
			result += fmt.Sprintf("%g", floatValue) // Format as a float
		}
	} else if j.Bool != nil {
		result += fmt.Sprintf("%t", *j.Bool)
	} else if j.Null != nil {
		result += "null"
	}

	return result
}
