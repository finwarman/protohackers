package json

import (
	"fmt"
	"reflect"
	"testing"
)

func floatPtr(f float64) *float64 {
	return &f
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *Boolean {
	_bool := Boolean(b)
	return &_bool
}

func TestBasicParsing(t *testing.T) {
	parsedValue, err := ParseJSON("{\"method\":\"isPrime\",\"number\":970747}")
	if err != nil {
		t.Fatalf("Parsing Error:\n%v\n\n", err)
	}

	fmt.Printf("Parsed JSON Value:\n%v\n\n", parsedValue)

	expectedValue := &JSONValue{
		Object: &JSONObject{
			Pairs: []*JSONPair{
				{
					Key:   "method",
					Value: &JSONValue{Str: strPtr("isPrime")},
				},
				{
					Key:   "number",
					Value: &JSONValue{Number: floatPtr(970747)},
				},
			},
		},
	}

	if !reflect.DeepEqual(parsedValue, expectedValue) {
		t.Errorf("Parsed JSON does not match the expected structure."+
			"\nGot:\n%v\n\nExpected:\n%v", parsedValue, expectedValue)
	}
}

func TestMultilineParsing(t *testing.T) {
	parsedValue, err := ParseJSON(`
		{
			"method": "isPrime",
			"number": 970747
		}
	`)
	if err != nil {
		t.Fatalf("Parsing Error:\n %v\n\n", err)
	}

	expectedValue := &JSONValue{
		Object: &JSONObject{
			Pairs: []*JSONPair{
				{
					Key:   "method",
					Value: &JSONValue{Str: strPtr("isPrime")},
				},
				{
					Key:   "number",
					Value: &JSONValue{Number: floatPtr(970747)},
				},
			},
		},
	}

	if !reflect.DeepEqual(parsedValue, expectedValue) {
		t.Errorf("Parsed JSON does not match the expected structure."+
			"\nGot:\n%v\n\nExpected:\n%v", parsedValue, expectedValue)
	}
}

func TestNestedArrayObjectParsing(t *testing.T) {
	parsedValue, err := ParseJSON(`
	[
		{
			"method": "foobar",
			"number": 100000
		},
		[
			{ "values": [1, 2] },
			"barfoo",
			true,
			false,
			null,
			1234,
			1000.25
		]
	]
	`)
	if err != nil {
		t.Fatalf("Parsing Error:\n %v\n\n", err)
	}

	expectedValue := &JSONValue{
		Array: &JSONArray{
			Values: []*JSONValue{
				{
					Object: &JSONObject{
						Pairs: []*JSONPair{
							{
								Key:   "method",
								Value: &JSONValue{Str: strPtr("foobar")},
							},
							{
								Key:   "number",
								Value: &JSONValue{Number: floatPtr(100000)},
							},
						},
					},
				},
				&JSONValue{
					Array: &JSONArray{
						Values: []*JSONValue{
							{
								Object: &JSONObject{
									Pairs: []*JSONPair{
										{
											Key: "values",
											Value: &JSONValue{
												Array: &JSONArray{
													Values: []*JSONValue{
														{Number: floatPtr(1)},
														{Number: floatPtr(2)},
													},
												},
											},
										},
									},
								},
							},
							{
								Str: strPtr("barfoo"),
							},
							{
								Bool: boolPtr(true),
							},
							{
								Bool: boolPtr(false),
							},
							{
								Null: strPtr("null"),
							},
							{
								Number: floatPtr(1234),
							},
							{
								Number: floatPtr(1000.25),
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(parsedValue, expectedValue) {
		t.Errorf("Parsed JSON does not match the expected structure."+
			"\nGot:\n%v\n\nExpected:\n%v", parsedValue, expectedValue)
	}
}
