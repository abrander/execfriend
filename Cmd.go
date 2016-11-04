package execfriend

import (
	"os/exec"
)

type (
	// Cmd represents an external command being prepared or run.
	Cmd struct {
		*exec.Cmd
		updates []chan []byte
		Output  []byte
	}
)

// Command returns the Cmd struct to execute the named program with the given
// arguments.
func Command(name string, arg ...string) *Cmd {
	c := &Cmd{}

	c.Cmd = exec.Command(name, arg...)

	// We use ourself as io.Writers for both stdout and stderr.
	c.Stdout = c
	c.Stderr = c

	return c
}

// Write implement io.Writer. This allows ud to easily read output from command.
func (c *Cmd) Write(p []byte) (n int, err error) {
	for _, ch := range c.updates {
		// Save combined output.
		c.Output = append(c.Output, p...)

		// Notify channel.
		ch <- p
	}

	return len(p), nil
}

// UpdateChan will return a channel which will transport updates to output. The
// channel will be closed when there's no more data written.
func (c *Cmd) UpdateChan() chan []byte {
	ch := make(chan []byte)

	c.updates = append(c.updates, ch)

	return ch
}

// Wait waits for the command to exit. It must have been started by Start. This
// will also close the channel(s) returned by UpdateChan() if any.
func (c *Cmd) Wait() error {
	err := c.Cmd.Wait()
	for _, ch := range c.updates {
		close(ch)
	}

	return err
}
