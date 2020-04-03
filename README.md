# go-lru [![Go Report Card](https://goreportcard.com/badge/github.com/xf005/go-lru)](https://goreportcard.com/report/github.com/xf005/go-lru) [![Build Status](https://travis-ci.org/xf005/go-lru.svg?branch=master)](https://travis-ci.org/xf005/go-lru) 

go-lru is an MIT-licensed Go LRU cache bases on GroupCache, with expire time supported

## Example

Set key with expire time

```
cache := NewCache(100,60) // max entries in cache is 100, expired after 1 second
cache.Set("a", 1234) // key "a" 
```

## API doc

API documentation is available via  [https://godoc.org/github.com/xf005/go-lru](https://godoc.org/github.com/xf005/go-lru)
