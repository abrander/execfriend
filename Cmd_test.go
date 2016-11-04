package execfriend

import (
	"os"
	"sync"
	"testing"
)

func TestCommand(t *testing.T) {
	c := Command("echo", "buh")
	ch := c.UpdateChan()

	c.Start()
	var output string

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			payload, more := <-ch
			if !more {
				break
			}

			output += string(payload)
		}

		wg.Done()
	}()

	c.Wait()

	if output != "buh\n" {
		t.Fatalf("Unexpected output, got '%s', expected 'buh'", output)
	}
	wg.Wait()
}

func ExampleCommand() {
	c := Command("echo", "hello world")
	c.Start()
	ch := c.UpdateChan()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			update, more := <-ch
			if !more {
				break
			}

			// Do something clever with update...
			os.Stdout.Write(update)
		}
		wg.Done()
	}()

	c.Wait()
	wg.Wait()

	// Output: hello world
}
