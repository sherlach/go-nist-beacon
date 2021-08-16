# go-nist-beacon
NIST Randomness Beacon (https://beacon.nist.gov/home) api wrapper

[![GoDoc](https://godoc.org/github.com/ClownKnuckle/go-nist-beacon?status.svg)](https://godoc.org/github.com/ClownKnuckle/go-nist-beacon)
[![Go Report Card](https://goreportcard.com/badge/github.com/ClownKnuckle/go-nist-beacon)](https://goreportcard.com/report/github.com/ClownKnuckle/go-nist-beacon)
[![Build Status](https://travis-ci.org/ClownKnuckle/go-nist-beacon.svg?branch=master)](https://travis-ci.org/ClownKnuckle/go-nist-beacon)

## NIST Randomness Beacon (https://beacon.nist.gov/home) API wrapper

<span style="color:red">**WARNING: DO NOT USE BEACON GENERATED VALUES AS SECRET CRYPTOGRAPHIC KEYS.**</span>

### Usage example:
```
import (
  "fmt"
  "github.com/ClownKnuckle/go-nist-beacon"
  "math/big"
  "math/rand"
) 
  
func main() {
  r, err := beacon.LastRecord()
  if err != nil {
    panic(err)
  }
  
  i := new(big.Int)
  ra := rand.New(rand.NewSource(i.Rsh(&r.SeedValue, 448).Int64()))
  fmt.Println(ra.Int(), ra.Int(), i.Rsh(&r.SeedValue, 448).Int64(), r.SeedValue)
}
```
Using the same seed value the random numbers generated are the same.

A much simpler version of the same would be:
```
import (
  "fmt"
  "github.com/sherlach/go-nist-beacon"
) 
  
func main() {
  ra, err := beacon.NewUpdatedRand()
  if err != nil {
    panic(err)
  }
  
  fmt.Println(ra.Int(), ra.Int())
}
```
