package errchan

import (
	"fmt"
	"os"
)

func ExampleNewErrChan() {
	ec := NewErrChan()

	ec.Wrap(func() error {
		return os.Remove("/tmp/testmenow")
	})

	ec.Wrap(func() error {
		return os.Remove("/tmp/testmelater")
	})

	err := <-ec.Next()
	if err != nil {
		fmt.Println("Got an error", err)
	}
}
