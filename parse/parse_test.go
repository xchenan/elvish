package parse

import (
	"fmt"
	"github.com/xiaq/elvish/util"
	"reflect"
	"testing"
)

var parseTests = []struct {
	in     string
	wanted Node
}{
	{"", newList(0)},
	{"ls", newList( // chunk
		0, newList( // pipeline
			0, &FormNode{ // form
				0, newList( // term
					0, &FactorNode{ // factor
						0, StringFactor, newString(0, "ls", "ls")}),
				newList(2), nil}))},
}

func TestParse(t *testing.T) {
	for i, tt := range parseTests {
		out, err := Parse(fmt.Sprintf("<test %d>", i), tt.in)
		if !reflect.DeepEqual(out, tt.wanted) || err != nil {
			t.Errorf("Parse(*, %q) =>\n(%s, %v), want\n(%s, nil) (up to DeepEqual)", tt.in, util.GoPrint(out), err, util.GoPrint(tt.wanted))
		}
	}
}

var completeTests = []struct {
	in     string
	wanted *Context
}{
/*
	{"", &Context{
		newList(0), newList(0),
		&FactorNode{0, StringFactor, newString(0, "", "")}}},
	{"l", &Context{
		newList(0), newList(0),
		&FactorNode{0, StringFactor, newString(0, "l", "l")}}},
	{"ls ", &Context{
		newList(0, newList(0, &FactorNode{0, StringFactor, newString(0, "ls", "ls")})),
		newList(0),
		&FactorNode{0, StringFactor, newString(0, "", "")}}},
	{"ls a", &Context{
		newList(0, newList(0, &FactorNode{0, StringFactor, newString(0, "ls", "ls")})),
		newList(0),
		&FactorNode{0, StringFactor, newString(0, "a", "a")}}},
	{"ls $a", &Context{
		newList(0, newList(0, &FactorNode{0, StringFactor, newString(0, "ls", "ls")})),
		newList(0),
		&FactorNode{0, VariableFactor, newString(0, "a", "a")}}},
*/
}

func TestComplete(t *testing.T) {
	for i, tt := range completeTests {
		out, err := Complete(fmt.Sprintf("<test %d>", i), tt.in)
		if !reflect.DeepEqual(out, tt.wanted) || err != nil {
			t.Errorf("Complete(*, %q) => (%s, %v), want (%s, nil)", tt.in, util.GoPrint(out), err, util.GoPrint(tt.wanted))
		}
	}
}