# dProxy - document proxy

dProxy is a proxy to access `interface{}` (document) by simple query.
It is intetedd to be used with `json.Unmarshal()` or `json.NewDecorder()`.

```go
var v interface{}
json.Unmarshal(byte[](`
  "cities": [ "tokyo", 100, "osaka", 200, "hakata", 300 ],
  "data": {
    "custom": [ "male", 23, "female", 24 ]
  }
`), &v)

// s == "tokyo", got a sting.
s, _ := dproxy.New(v).M("cities").A(0).String()

// err != nil, type not matched.
_, err := dproxy.New(v).M("cities").A(0).Float64()

// n == 0, got a float64
n, _ := dproxy.New(v).M("cities").A(1).Float64()

// be able to access nested.
dproxy.New(v).M("data").M("custom").A(0).String()

_, err = dproxy.New(v).M("data").M("castom").String
// err.Error() == "not found: data.castom", wrong query.
```

## LICENSE

MIT license.  See LICENSE.
