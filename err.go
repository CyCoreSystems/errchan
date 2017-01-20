// Package errchan provides a simple wrapper for safely handling errors
// asynchronously from a function's execution.
//
// The intention of this library is to enable background execution of functions
// which return an error while not losing any non-nil errors which are returned
// by that function.
//
// A typical use case if the startup of several long-running functions, any of
// which might return an error which needs to be handled.
//
// For example:
// 	func main() {
//      ctx, cancel := context.WithCancel(context.Background())
//      defer cancel()
//
//      ec := errchan.New()
//		  defer ec.Close()
//
//      ec.Wrap(func() error {
//        startA(ctx)
//      })
//      ec.Wrap(func() error {
//        startB(ctx)
//      })
//      ec.Wrap(func() error {
//        startC(ctx)
//      })
//
//      select {
//        case ctx.Done():
//        case err := <-ec.Next():
//          if err != nil {
//            fmt.Println("Caught error; exiting", err)
//          }
//      }
//    }
//
package errchan

import "sync"

// ErrChan provides a safe communication channel
// for returning asynchronous errors
type ErrChan struct {
	c      chan error
	closed bool

	mu sync.Mutex
}

// New returns a new ErrChan async error handler
func New() *ErrChan {
	return &ErrChan{
		c: make(chan error),
	}
}

// Close safely closes the ErrChan
func (e *ErrChan) Close() {
	if !e.closed {
		e.closed = true
		close(e.c)
	}
}

// Next returns the next error from the ErrChan
func (e *ErrChan) Next() <-chan error {
	return e.c
}

// Wrap takes a function which returns an error and executes
// it inside a goroutine.  This goroutine captures any (non-nil)
// error from the execution and returns it on the ErrChan.
func (e *ErrChan) Wrap(f func() error) {
	go func() {
		if err := f(); err != nil {
			if !e.closed {
				e.c <- err
			}
		}
	}()
}
