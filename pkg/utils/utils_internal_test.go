package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPtr(t *testing.T) {
	testVal := "hello"
	testPtr := StringPtr(testVal)
	assert.NotNil(t, testPtr)
	assert.Equal(t, testVal, *testPtr)
	assert.Equal(t, &testVal, testPtr)
}

func TestStringValOr(t *testing.T) {
	testDefault := "defaultval"
	testVal := "some non-default value"

	v1 := StringValOr(nil, testDefault)
	assert.Equal(t, testDefault, v1)

	v2 := StringValOr(&testVal, testDefault)
	assert.Equal(t, testVal, v2)
}

func ExampleStringPtr() {
	StringVal := "hello"
	StringP1 := &StringVal
	StringP2 := StringPtr(StringVal)

	fmt.Printf("value: %s\n", StringVal)
	fmt.Printf("ptr inline: %s\n", *StringP1)
	fmt.Printf("ptr from func: %s\n", *StringP2)
	fmt.Printf("ptrs are different: %v", (StringP1 != StringP2))

	// Output:
	// value: hello
	// ptr inline: hello
	// ptr from func: hello
	// ptrs are different: true
}

func ExampleStringValOr() {
	defaultStringVal := "default hello"

	var StringP *string
	StringVal1 := StringValOr(StringP, defaultStringVal)
	fmt.Println(StringVal1)

	stringTmp := "some non default value"
	StringP2 := &stringTmp
	StringVal2 := StringValOr(StringP2, defaultStringVal)
	fmt.Println(StringVal2)

	// Output:
	// default hello
	// some non default value
}

func TestIntPtr(t *testing.T) {
	testVal := 124
	testPtr := IntPtr(testVal)
	assert.NotNil(t, testPtr)
	assert.Equal(t, testVal, *testPtr)
	assert.Equal(t, &testVal, testPtr)
}

func TestIntValOr(t *testing.T) {
	testDefault := 124
	testVal := 999

	v1 := IntValOr(nil, testDefault)
	assert.Equal(t, testDefault, v1)

	v2 := IntValOr(&testVal, testDefault)
	assert.Equal(t, testVal, v2)
}

func ExampleIntPtr() {
	IntVal := 123
	IntP1 := &IntVal
	IntP2 := IntPtr(IntVal)

	fmt.Printf("value: %d\n", IntVal)
	fmt.Printf("ptr inline: %d\n", *IntP1)
	fmt.Printf("ptr from func: %d\n", *IntP2)
	fmt.Printf("ptrs are different: %v", (IntP1 != IntP2))

	// Output:
	// value: 123
	// ptr inline: 123
	// ptr from func: 123
	// ptrs are different: true
}

func ExampleIntValOr() {
	defaultIntVal := 123

	var IntP *int
	IntVal1 := IntValOr(IntP, defaultIntVal)
	fmt.Printf("default %d\n", IntVal1)

	intTmp := 999
	IntP2 := &intTmp
	IntVal2 := IntValOr(IntP2, defaultIntVal)
	fmt.Printf("not default %d\n", IntVal2)

	// Output:
	// default 123
	// not default 999
}
