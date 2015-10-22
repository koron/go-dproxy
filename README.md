# dProxy - document proxy

dProxy is a proxy to access `interface{}` (document) by simple query.
It is intetedd to be used with `json.Unmarshal()` or `json.NewDecorder()`.

See codes for overview.

```go
import (
  "encoding/json"

  "github.com/koron/go-dproxy"
)

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

// n == 200, got a float64
n, _ := dproxy.New(v).M("cities").A(3).Float64()

// be able to access nested.
dproxy.New(v).M("data").M("custom").A(0).String()

_, err = dproxy.New(v).M("data").M("kustom").String()
// err.Error() == "not found: data.kustom", wrong query can be verified.
```


## Getting started

1.  wrap a value (`interface{}`) with `dproxy.New()` get `dproxy.Proxy`.

    ```go
    p := dproxy.New(v) // v should be a value of interface{}
    ```

2.  query as a map by `M()`, it returns `dproxy.Proxy`.

    ```go
    p.M("cities")
    ```

3.  query as an array with `A()`, it returns `dproxy.Proxy`.

    ```go
    p.A(3)
    ```

4.  Can be chained above queris.

    ```go
    p.M("cities").A(3)
    ```

5.  Get a value finally.

    ```go
    n, _ := p.M("cities").A(3).Int64()
    ```

6.  You'll get some error, if there are some mistake in query.

    ```go
    // OOPS! "kustom" is typo, must be "custom"
    _, err := p.M("data").M("kustom").A(3).Int64()

    fmt.Println(err)
    // "not found: data.kustom"
    ```

    You can verify queries easily.


## LICENSE

MIT license.  See LICENSE.
