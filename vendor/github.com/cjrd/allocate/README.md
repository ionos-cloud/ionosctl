## Allocate
[![Build Status](https://travis-ci.org/cjrd/allocate.svg?branch=master)](https://travis-ci.org/cjrd/allocate)
[![Coverage Status](https://coveralls.io/repos/github/cjrd/allocate/badge.svg?branch=master)](https://coveralls.io/github/cjrd/allocate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/cjrd/allocate)](https://goreportcard.com/report/github.com/cjrd/allocate)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/mkideal/cli?status.svg)](https://godoc.org/github.com/cjrd/allocate)

Allocate provides functions for allocating golang structs so that pointer fields are pointers to zero'd values instead of `nil`. See the godoc's for more information: https://godoc.org/github.com/cjrd/allocate

### Brief Example
```go
package main

import (
    "fmt"
    "github.com/cjrd/allocate"
)

type TopLevelStruct struct {
    MyEmbeddedStruct *EmbeddedStruct
}

type EmbeddedStruct struct {
    SomeString string
    SomeInt int
}

func main() {
    topStruct := new(TopLevelStruct)
    fmt.Printf("before using allocate.Zero: %v\n", topStruct)

    allocate.Zero(&topStruct)
    fmt.Printf("post allocate.Zero: %v\n", topStruct)
    fmt.Printf("topStruct.MyEmbeddedStruct.SomeInt==%d\n", topStruct.MyEmbeddedStruct.SomeInt)
    // Note that panics would occur by executing `*topStruct.MyEmbeddedStruct` or
    //`topStruct.MyEmbeddedStruct.SomeInt`
}
```

```bash
# OUTPUT
before using allocate.Zero: &{<nil>}
post allocate.Zero: &{0x8201d2400}
topStruct.MyEmbeddedStruct.SomeInt == 0
```

### Use Cases

* Initializing structures that contain any type of pointer fields, including recursive struct fields
* Preventing panics by ensuring that all fields of a struct are initialized
* Initializing [golang protobuf struct](https://github.com/golang/protobuf) (the golang protobuf makes heavy use of pointers to embedded structs that contain pointers to embedded structs, ad infinitum)
* Initializing structs for black box testing (see also https://golang.org/pkg/testing/quick/)

