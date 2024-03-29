<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# utils

```go
import "github.com/ecshreve/slomad/pkg/utils"
```

package utils provides general utility functions.

## Index

- [func IntPtr\(i int\) \*int](<#IntPtr>)
- [func IntValOr\(ip \*int, val int\) int](<#IntValOr>)
- [func StringPtr\(s string\) \*string](<#StringPtr>)
- [func StringValOr\(sp \*string, val string\) string](<#StringValOr>)


<a name="IntPtr"></a>
## func [IntPtr](<https://github.com/ecshreve/slomad/blob/main/pkg/utils/utils.go#L19>)

```go
func IntPtr(i int) *int
```

IntPtr returns a pointer to the given int.

<details><summary>Example</summary>
<p>



```go
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
```

#### Output

```
value: 123
ptr inline: 123
ptr from func: 123
ptrs are different: true
```

</p>
</details>

<a name="IntValOr"></a>
## func [IntValOr](<https://github.com/ecshreve/slomad/blob/main/pkg/utils/utils.go#L24>)

```go
func IntValOr(ip *int, val int) int
```

IntValOr returns the value of an int pointer if it's not nil, or a default.

<details><summary>Example</summary>
<p>



```go
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
```

#### Output

```
default 123
not default 999
```

</p>
</details>

<a name="StringPtr"></a>
## func [StringPtr](<https://github.com/ecshreve/slomad/blob/main/pkg/utils/utils.go#L5>)

```go
func StringPtr(s string) *string
```

StringPtr returns a pointer to the given string.

<details><summary>Example</summary>
<p>



```go
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
```

#### Output

```
value: hello
ptr inline: hello
ptr from func: hello
ptrs are different: true
```

</p>
</details>

<a name="StringValOr"></a>
## func [StringValOr](<https://github.com/ecshreve/slomad/blob/main/pkg/utils/utils.go#L11>)

```go
func StringValOr(sp *string, val string) string
```

StringValOr returns the value of a string pointer if it's not nil, or a default value otherwise.

<details><summary>Example</summary>
<p>



```go
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
```

#### Output

```
default hello
some non default value
```

</p>
</details>

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
