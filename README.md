# dProxy - document proxy

dProxy is a proxy to access `interface{}` (document) by simple query.
It is intented to be used with `json.Unmarshal()` or `json.NewDecorder()`.

See codes for overview.

```go
import (
  "encoding/json"

  "github.com/koron/go-dproxy"
)

var v interface{}
json.Unmarshal(byte[](`{
  "cities": [ "tokyo", 100, "osaka", 200, "hakata", 300 ],
  "data": {
    "custom": [ "male", 23, "female", 24 ]
  }
}`), &v)

// s == "tokyo", got a string.
s, _ := dproxy.New(v).M("cities").A(0).String()

// err != nil, type not matched.
_, err := dproxy.New(v).M("cities").A(0).Float64()

// n == 200, got a float64
n, _ := dproxy.New(v).M("cities").A(3).Float64()

// can be chained.
dproxy.New(v).M("data").M("custom").A(0).String()

// err.Error() == "not found: data.kustom", wrong query can be verified.
_, err = dproxy.New(v).M("data").M("kustom").String()
```


## Getting started

1.  Wrap a value (`interface{}`) with `dproxy.New()` get `dproxy.Proxy`.

    ```go
    p := dproxy.New(v) // v should be a value of interface{}
    ```

2.  Query as a map (`map[string]interface{}`)by `M()`, returns `dproxy.Proxy`.

    ```go
    p.M("cities")
    ```

3.  Query as an array (`[]interface{}`) with `A()`, returns `dproxy.Proxy`.

    ```go
    p.A(3)
    ```

4.  Therefore, can be chained queries.

    ```go
    p.M("cities").A(3)
    ```

5.  Get a value finally.

    ```go
    n, _ := p.M("cities").A(3).Int64()
    ```

6.  You'll get an error when getting a value, if there were some mistakes.

    ```go
    // OOPS! "kustom" is typo, must be "custom"
    _, err := p.M("data").M("kustom").A(3).Int64()

    // "not found: data.kustom"
    fmt.Println(err)
    ```

7.  If you tried to get a value as different type, get an error.

    ```go
    // OOPS! "cities[3]" (=200) should be float64 or int64.
    _, err := p.M("cities").A(3).String()

    // "not matched types: expected=string actual=float64: cities[3]"
    fmt.Println(err)
    ```

8.  You can verify queries easily.


## LICENSE

MIT license.  See LICENSE.
