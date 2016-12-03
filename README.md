# routing
[![Build Status](https://travis-ci.org/DavidCai1993/routing.svg?branch=master)](https://travis-ci.org/DavidCai1993/routing)
[![Coverage Status](https://coveralls.io/repos/github/DavidCai1993/routing/badge.svg?branch=master)](https://coveralls.io/github/DavidCai1993/routing?branch=master)

URL routing based on trie.

`routing` is only used to define and match URLs with some custom meta infomation. And basing on it, you can build your own http router with any additional feature.

## Installation

```
go get -u github.com/DavidCai1993/routing
```

## Support signatures

- string: `/hello`
- separated string: `/a|b|c`
- regex: `/([0-9a-f]{24})`
- named parameter: `/:id`
- named separated string: `/:id(a|b|c)`
- named regex: `/:id([0-9a-f]{24})`

## Documentation

API documentation can be found here: https://godoc.org/github.com/DavidCai1993/routing

## Usage

```go
router := routing.New()

router.Define("/:type(a|b)/:id(0-9a-f]{24})", yourHandler)

callback, params, ok := router.Match("/a/8")

fmt.Println(ok)
// -> true

fmt.Println(callback.(http.Handler))

fmt.Println(params)
// -> map[string]string{"type": "a", "id": "8"}
```
