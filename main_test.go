package main

import (
	"reflect"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"true", true},
		{"false", false},
		{"True", true},
		{"False", false},
		{"1", int64(1)},
		{"0", int64(0)},
		{"-1", int64(-1)},
		{"1234567890", int64(1234567890)},
		{"-1234567890", int64(-1234567890)},
		{"3.1415", float64(3.1415)},
		{"-3.1415", float64(-3.1415)},
		{"0.0", float64(0.0)},
		{"hello", "hello"},
		{"", ""},
		{" ", " "},
		{"truee", "truee"},   // Should be parsed as string
		{"123abc", "123abc"}, // Should be parsed as string
	}

	for _, test := range tests {
		result := parseValue(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseValue(%q) = %v (%T); expected %v (%T)",
				test.input, result, result, test.expected, test.expected)
		}
	}
}

func TestParseDynamicFlags(t *testing.T) {
	tests := []struct {
		args     []string
		expected map[string]interface{}
	}{
		{
			args: []string{"--bool=true", "--int=123", "--float=45.67", "--string=hello"},
			expected: map[string]interface{}{
				"bool":   true,
				"int":    int64(123),
				"float":  float64(45.67),
				"string": "hello",
			},
		},
		{
			args: []string{"--name=John", "--age", "30", "--active", "false"},
			expected: map[string]interface{}{
				"name":   "John",
				"age":    int64(30),
				"active": false,
			},
		},
		{
			args: []string{"--empty="},
			expected: map[string]interface{}{
				"empty": "",
			},
		},
		{
			args: []string{"--list", "item1,item2,item3"},
			expected: map[string]interface{}{
				"list": "item1,item2,item3",
			},
		},
		{
			args: []string{"--number=00123"},
			expected: map[string]interface{}{
				"number": int64(123),
			},
		},
		{
			args: []string{"--negativeFloat", "-123.456"},
			expected: map[string]interface{}{
				"negativeFloat": float64(-123.456),
			},
		},
		{
			args: []string{"--weirdBool", "truee"},
			expected: map[string]interface{}{
				"weirdBool": "truee",
			},
		},
	}

	for _, test := range tests {
		result := parseDynamicFlags(test.args)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseDynamicFlags(%v) = %v; expected %v",
				test.args, result, test.expected)
		}
	}
}
