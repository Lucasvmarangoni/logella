

<div align="center">

# Logella

A simple loggers and errors library.

</div>


<div align="left">

### Packages:

- <a href="#errors-package">
   Errors
  </a> 

- <a href="#router-package">
    Router
</a>       
</div>

<br><br>

## Errors Package

The `Errors` package is a custom error handling library. Its primary feature is to attach contextual information to errors, allowing them to be propagated up the call stack. 

It also provides standardized error types, such as `invalid` and `required`.

### Import
```go
import "github.com/Lucasvmarangoni/logella/err"
```

### Use
```go
errors.ErrCtx(err, "repo.InitTables")
```

### ErrCtx(err error, value string)

This function creates an error with context. Here are examples of how to use it:

```go
cfg.Db, err = pgx.ParseConfig(url)
if err != nil {
    return nil, errors.ErrCtx(err, "pgx.ParseConfig(url)")
}
```

```go
dbConnection, err := db.Connect(ctx)
if err != nil {
    return nil, errors.ErrCtx(err, "db.Connect")
}
```

The output of this function includes a field named "Operation", which is the main feature of this package.

### Standard Errors

The package provides standardized errors, such as `IsInvalidError` and `IsRequiredError`. Here's an example of how to use `IsInvalidError`:

```go
errors.IsInvalidError("Customer", "Must be google uuid")
```

<br><br>

## Router Package

The `Router` package is a logging library for initializing routes using go-chi.

### Import

```go
import "github.com/Lucasvmarangoni/logella/router"
```

### Use

```go
router := router.NewRouter()
```

### Instance Creation

To create a new instance of the `Router`, use the `NewRouter` function:

```go
func NewRouter() *Router {
    ro := &Router{}    
    return ro
}
```

### InitializeRoute Function

The `InitializeRoute` function is the main function of the `Router` package. It takes a chi router, a path, and a handler function as arguments:

```go
func (ro *Router) InitializeRoute(r chi.Router, path string, handler http.HandlerFunc)
```

```go
router.InitializeRoute(r, "/jwt", u.userHandler.Authentication)
```

### Method Function

The `Method` function sets the HTTP method for the route. It accepts a string argument, which can be either uppercase or lowercase:

```go
func (ro *Router) Method(m string) *Router
```

```go
router.Method("POST").InitializeRoute(r, "/jwt", u.userHandler.Authentication)
```

### Prefix Function

The `Prefix` function sets a prefix for the route, if there is one:

```go
func (ro *Router) Prefix(p string) *Router 
```

```go
router.Method("POST").Prefix("/authn").InitializeRoute(r, "/jwt", u.userHandler.Authentication)
```