package slomad_test

import (
	"fmt"
	"testing"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/stretchr/testify/assert"
)

func TestStringPtr(t *testing.T) {
	testVal := "hello"
	testPtr := slomad.StringPtr(testVal)
	assert.NotNil(t, testPtr)
	assert.Equal(t, testVal, *testPtr)
	assert.Equal(t, &testVal, testPtr)
}

func TestStringValOr(t *testing.T) {
	testDefault := "defaultval"
	testVal := "some non-default value"

	v1 := slomad.StringValOr(nil, testDefault)
	assert.Equal(t, testDefault, v1)

	v2 := slomad.StringValOr(&testVal, testDefault)
	assert.Equal(t, testVal, v2)
}

func ExampleStringPtr() {
	stringVal := "hello"
	stringP1 := &stringVal
	stringP2 := slomad.StringPtr(stringVal)

	fmt.Printf("value: %s\n", stringVal)
	fmt.Printf("ptr inline: %s\n", *stringP1)
	fmt.Printf("ptr from func: %s\n", *stringP2)
	fmt.Printf("ptrs are different: %v", (stringP1 != stringP2))

	// Output:
	// value: hello
	// ptr inline: hello
	// ptr from func: hello
	// ptrs are different: true
}

func ExampleStringValOr() {
	defaultStringVal := "default hello"

	var stringP *string
	stringVal1 := slomad.StringValOr(stringP, defaultStringVal)
	fmt.Println(stringVal1)

	stringTmp := "some non default value"
	stringP2 := &stringTmp
	stringVal2 := slomad.StringValOr(stringP2, defaultStringVal)
	fmt.Println(stringVal2)

	// Output:
	// default hello
	// some non default value
}

func TestIntPtr(t *testing.T) {
	testVal := 124
	testPtr := slomad.IntPtr(testVal)
	assert.NotNil(t, testPtr)
	assert.Equal(t, testVal, *testPtr)
	assert.Equal(t, &testVal, testPtr)
}

func TestIntValOr(t *testing.T) {
	testDefault := 124
	testVal := 999

	v1 := slomad.IntValOr(nil, testDefault)
	assert.Equal(t, testDefault, v1)

	v2 := slomad.IntValOr(&testVal, testDefault)
	assert.Equal(t, testVal, v2)
}

func ExampleIntPtr() {
	intVal := 123
	intP1 := &intVal
	intP2 := slomad.IntPtr(intVal)

	fmt.Printf("value: %d\n", intVal)
	fmt.Printf("ptr inline: %d\n", *intP1)
	fmt.Printf("ptr from func: %d\n", *intP2)
	fmt.Printf("ptrs are different: %v", (intP1 != intP2))

	// Output:
	// value: 123
	// ptr inline: 123
	// ptr from func: 123
	// ptrs are different: true
}

func ExampleIntValOr() {
	defaultIntVal := 123

	var intP *int
	intVal1 := slomad.IntValOr(intP, defaultIntVal)
	fmt.Printf("default %d\n", intVal1)

	intTmp := 999
	intP2 := &intTmp
	intVal2 := slomad.IntValOr(intP2, defaultIntVal)
	fmt.Printf("not default %d\n", intVal2)

	// Output:
	// default 123
	// not default 999
}
