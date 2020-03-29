# errchan [![](https://godoc.org/github.com/CyCoreSystems/errchan?status.svg)](http://godoc.org/github.com/CyCoreSystems/errchan)

## ARCHIVED:  please use the official [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc) package instead.

A simple wrapper for safely handling errors asynchronously from a function's
execution. 

Package `errchan` provides a simple wrapper for safely handling errors
asynchronously from a function's execution.

The intention of this library is to enable background execution of functions
which return an error while not losing any non-nil errors which are returned
by that function.

A typical use case if the startup of several long-running functions, any of
which might return an error which needs to be handled.

For example:

```go
	func main() {
     ctx, cancel := context.WithCancel(context.Background())
     defer cancel()

     ec := errchan.New()
	  defer ec.Close()

     ec.Wrap(func() error {
       startA(ctx)
     })
     ec.Wrap(func() error {
       startB(ctx)
     })
     ec.Wrap(func() error {
       startC(ctx)
     })

     select {
       case ctx.Done():
       case err := <-ec.Next():
         if err != nil {
           fmt.Println("Caught error; exiting", err)
         }
     }
   }
```

A functional example may be found in the [example
directory](/example).

## Conditional goroutines

A common pattern is the execution of a subroutine which should run in the
background but only after a number of subroutine-specific setups and checks have
been performed.  While all of this can be performed inside that subroutine, the
higher-order logic makes such backgrounding invisible, because the `go` keyword
is embedded within that subroutine.

By wrapping the subroutine inside an `errchan` execution, the semantics become
clear:  that the subroutine should run in the background but return an error if
it is unable to start.

`errchan.Go` provides this mechanism.
