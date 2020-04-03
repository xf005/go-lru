# go-lru 

go-lru is an MIT-licensed Go LRU cache bases on GroupCache, with expire time supported

## Example

Set key with expire time

```
cache := NewCache(100,60) // max entries in cache is 100, expired after 60 second
cache.Set("a", 1234) // key "a" 
```
