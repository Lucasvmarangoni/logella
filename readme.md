

<div align="center">

# Logella

A simple loggers and errors library.

</div>


<div align="left">

### Packages:

- <a href="#logger-package">
   Logger
  </a> 

- <a href="#errs-package">
   Errs
  </a> 

- <a href="#router-package">
    Router
</a>       
</div>

<br><br>

## Logger Package

This package defines the log settings (zerolog). Allows you to use default or customized colors configuration.

### Import
```go
import "github.com/Lucasvmarangoni/logella/config/log"
```

**ConfigDefault**: It has pre-defined log color settings, meaning there is no need to specify a any parameter.

```go
logger.ConfigDefault()
```

**ConfigCustom**: Allows you to customize log level, message, and operation colors.

```go
ConfigCustom(info, err, warn, debug, fatal, message, trace colors)
```

The parameters must be passed using the variables already defined by the package with the name of the colors.

<u>Colors</u>: Black, Red, Green, Yellow, Blue, Magenta, Cyan or White.

Example:


```go
logger.ConfigCustom(logger.Green, logger.Red, logger.Yellow, logger.Cyan, logger.Red, logger.Magenta, logger.Blue)
```

**Use Case**
```go
log.Info().Str("context", "TableRepository").Msg("Database - Created users table successfully.")
```

Output:

![alt text](img/router-example.png)

<br>

## Errs Package

The `Errs` package is a custom error handling library. Its primary feature is to attach Traceual information to errors, allowing them to be propagated up the call stack. 

It also provides standardized error types, such as `invalid` and `required`.

Output Example:
```go
log.Error().Err(errs.Trace(err).Stack()).Msg("example error")
log.Error().Err(errs.Unwrap(err).Stack()).Msg("example error")
log.Error().Err(errs.Unwrap(err).Stack()).Msg("example error")
```

![alt text](img/log.png)

➤ It is possible click on the path value (<u>test/test.go:15</u>) to go directly to file and line.

### Import
```go
import "github.com/Lucasvmarangoni/logella/err"
```

<a href="#wrap">**Wrap**</a>: Ctx is used to add the error and the operation that triggered the exception. 
The operations stack is not returned by ErrCtx, but rather persisted. 

<a href="#trace">**Trace**</a>: Used to add the trace to stack.

<a href="#trace-method">**Trace (method)**</a>: Used to add the trace manuelly to stack.

<a href="#stack">**Stack**</a>: Stack returns the error along with the operations stack. Used in internals Logs.

<a href="#toclient">**ToClient**</a>: Used to send error message to client.

<a href="#msg">**Msg**</a>: Used to add a message to error.

<a href="#unwrap">**Unwrap**</a>: It makes the type assertion and is used to access de Error Struct whitout performing other functionality.

<a href="#gethttpstatusfrompgerror">**GetHTTPStatusFromPgError**</a>: Used to determine HTTP status automatically based on database error message.

<a href="#new">**New**</a>: Used to add the trace to stack without http status code. Recommended for incializations exceptions.

<a href="#unwrap">**FailOnErrLog**</a>: Used to throw a fatal log in a simple way and with trace stack.

<a href="#panicerr">**PanicErr**</a>: Used to throw a panic based on a error value, whitout the needs to manualy use an if conditional.

<a href="#panicbool">**PanicBool**</a>: Used to throw a panic based on a boolean value, whitout the needs to manualy use an if conditional.

<a href="#isrequirederror">**IsRequiredError**</a>: Used to simplify create a "is requered" error with fmt.Errorf.

<a href="#isinvaliderror">**IsInvalidError**</a>: Used to simplify create a "is invalid" error with fmt.Errorf.

### Error Struct
```go
type Error struct {
	Cause   error   // The actual error thrown
	Code    int     // HTTP Status Code
	Message string  // Custom message
	trace error
}
```

#### Wrap
```go
func Wrap(cause error, code int) error
```

Example:
```go
errs.Wrap(err, http.StatusInternalServerError)
```

Use case:
```go
cfg.Db, err = pgx.ParseConfig(url)
if err != nil {
    return nil, errs.Wrap(err, http.StatusInternalServerError)
}
```

#### Trace

```go
func Trace(err error) *Error
```

Example:
```go
errs.Trace(err)
```

Use Case:
```go
func main() {
	_, err := handler()
	log.Error().Err(errs.Unwrap(err).Stack()).Msg(fmt.Sprint(errs.Unwrap(err).Code)) 
	// OUTPUT: 2025-02-20T18:29:26-03:00 ERROR ⇝ 500 error"test error | ➜ repository; ➜ service; ➜ handler; ➜ main"
}

func handler() (string, error) {
	_, err := service()
	return "", errs.Trace(err)
}

func service() (string, error) {
	err := repository()
	if err != nil {
		return "", errs.Trace(err)
	}
	return "", nil
}

func repository() error {
	return errs.Wrap(errors.New("test error"), 500)
}
```

**Recommendation**: When an error is issued in a situation that does not involve a call, such as in a conditional comparison, it is recommended that you encapsulate this inside a function, as follows:

```go
func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is not OK. got %d (%s), want %d (%s)",
			resp.StatusCode, http.StatusText(resp.StatusCode),
			http.StatusOK, http.StatusText(http.StatusOK))
	}
	return nil
}

if err := checkStatusCode(resp); err != nil {
		return errs.Trace(err)
	}
```

Or use Trace Method, like above:


#### Trace-method
This method can be used in different situations when you needed to inject your especific trace manually.

- When an error is throw inside an anonymous function;
- When an error is throw from an external library method; 
- Whem an error is issued in a situation that does not involve a call, such as in a conditional comparison.

```go
(e *Error) Trace(trace string) *Error
```
Example:
```go
errs.Wrap(err, http.StatusBadRequest).Trace("AnonymousFunction")
```

Use Case:
```go
func ExternalMethod() error {
	return errors.New("test error")
}

func Repository() error {
err := ExternalMethod()
	if err != nil {
		return errs.Wrap(err, http.StatusBadRequest).Trace("ExternalMethod")
	}
	return nil
}
// ExternalMethod ➤ Repository ➤ service ➤ handler ➤ main"
```

#### Stack

```go
func (e *Error) Stack() error 
```

Example:
```go
log.Error().Err(errs.Trace(err).Stack())
```

Use Case:
```go
authdata, err := u.userService.VerifyTOTP(id, totpToken.Token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Unwrap(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     http.StatusText(errs.Unwrap(err).Code),
			"message":    fmt.Sprintf("%v", errs.Unwrap(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Trace(err).Stack()).Msgf("error validate totp. | (%s)", requestID)
		return
	}
```

#### ToClient
Check the Code of the error to if it is 500 return "Internal Server Error" instead of Cause (actual error)

```go
func (e *Error) ToClient() error  
```

Example:
```go
errs.Unwrap(err).ToClient()
```

```go
	log.Error().Err(errs.Unwrap(err).ToClient()).Msg(fmt.Sprint(errs.Unwrap(err).Code)) 
	// OUTPUT: 2025-02-20T18:29:26-03:00 ERROR ⇝ 500 error"Internal Server Error"
```

Use Case:
```go
authdata, err := u.userService.VerifyTOTP(id, totpToken.Token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Unwrap(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     http.StatusText(errs.Unwrap(err).Code),
			"message":    fmt.Sprintf("%v", errs.Unwrap(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Trace(err).Stack()).Msgf("error validate totp. | (%s)", requestID)
		return
	}
```

#### Msg

```go
func (e *Error) Msg(message string) 
```

Example:
```go
errs.Trace(err).Msg("Message")
message := errs.Trace(err).Message
```


#### Unwrap
Used when you should not add a trace to the error. Otherwise a duplicate trace will be added.

```go
func Unwrap(err error) *Error 
```

Example:
```go
errs.Unwrap(err).Msg("Message")
message := errs.Unwrap(err).Message
code := errs.Unwrap(err).Code
cause := errs.Unwrap(err).Cause
```

Use Case
```go
log.Error().Err(errs.Trace(err).Stack()).Msg(fmt.Sprint(errs.Unwrap(err).Code))
```

```go
authdata, err := u.userService.VerifyTOTP(id, totpToken.Token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Unwrap(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     http.StatusText(errs.Unwrap(err).Code),
			"message":    fmt.Sprintf("%v", errs.Unwrap(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Trace(err).Stack()).Msgf("%s. | (%s)", errs.Unwrap(err).Message, requestID)
		return
	}
```

#### GetHTTPStatusFromPgError 
**Compatible with PGX V5 library**

A function to determine HTTP status automatically based on database error message.

```go
func GetHTTPStatusFromPgError(err error) int
```
Example: 
```go
return errs.Wrap(err, "row.Scan", errs.GetHTTPStatusFromPgError(err))
```

Use Case: 
```go
func (r *UserRepositoryDb) UpdateOTP(user *entities.User, ctx Trace.Trace) error {
	sql := `UPDATE users SET otp_auth_url = encrypt($2::BYTES, $4::BYTES, 'aes'), otp_secret = encrypt($3::BYTES, $4::BYTES, 'aes') WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.OtpAuthUrl,
			user.OtpSecret,
			r.key,
		)
		if err != nil {
			return errs.Wrap(err, errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return errs.Trace(err)
	}
	return nil
}
```


#### New

```go
New(cause error) *Error
```

Example:
```go
err := errs.New(errors.New("error"))
```

Use Case:
```go
func (r *TableRepositoryDb) initUserTable(ctx context.Context) error {
	_, err := r.tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY NOT NULL,
			name VARCHAR(15) NOT NULL,
			last_name BYTEA NOT NULL,				
			email BYTEA UNIQUE NOT NULL,				
			password TEXT NOT NULL,		
			created_at TIMESTAMP NOT NULL			
		)`)
	if err != nil {
		return errs.New(err)
	}
	log.Info().Str("context", "TableRepository").Msg("Database - Created users table successfully.")
	return nil
}
```

#### FailOnErrLog


```go
FailOnErrLog(err error, msg string)
```

Example:
```go
errs.FailOnErrLog(err, "database failled")
```

Use Case:
```go
if err = Database(); err != nil {
	errs.FailOnErrLog(err, "database failled")
}
```

#### PanicErr

```go
PanicErr(err error, msg string)
```

#### PanicBool

```go
PanicBool(boolean bool, msg string)
```

#### Standard Errors

The package provides standardized errors, such as `IsInvalidError` and `IsRequiredError`. Here's an example of how to use `IsInvalidError`:

```go
errs.IsInvalidError("Customer", "Must be google uuid")
```

- IsInvalidError(fieldName, msg string) error
- IsRequiredError(fieldName, msg string) error
- FailOnErrLog(err error, msg string)
- PanicErr(err error, ctx string)
- PanicBool(boolean bool, msg string)

<br>


## Router Package

The `Router` is a logging package for initializing routes using go-chi.

![Alt text](img/router.png)

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
router.InitializeRoute(r, "/login", u.userHandler.Authentication)
```

### Method Function

The `Method` function sets the HTTP method for the route. It accepts a string argument, which can be either uppercase or lowercase:

```go
func (ro *Router) Method(m string) *Router
```

```go
router.Method(router.POST).InitializeRoute(r, "/login", u.userHandler.Authentication)
```

### Prefix Function

The `Prefix` function sets a prefix for the route, if there is one:

```go
func (ro *Router) Prefix(p string) *Router 
```

```go
router.Method(router.POST).Prefix("/api").InitializeRoute(r, "/login", u.userHandler.Authentication)
```