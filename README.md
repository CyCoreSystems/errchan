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

A functional example may be found in the [test directory](/test).
