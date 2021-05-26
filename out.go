package midicatdrv

import (
	"fmt"
	"io"
	"os/exec"
	"sync"

	"gitlab.com/gomidi/midi"
)

func newOut(driver *Driver, number int, name string) midi.Out {
	o := &out{driver: driver, number: number, name: name}
	return o
}

type out struct {
	number int
	sync.RWMutex
	driver *Driver
	name   string
	wr     *io.PipeWriter
	rd     *io.PipeReader
	cmd    *exec.Cmd
}

func (o *out) fireCmd() error {
	o.Lock()
	defer o.Unlock()
	if o.cmd != nil {
		return fmt.Errorf("already running")
	}

	o.cmd = midiCatOutCmd(o.number)
	o.rd, o.wr = io.Pipe()
	o.cmd.Stdin = o.rd

	err := o.cmd.Start()
	if err != nil {
		o.rd = nil
		o.wr = nil
		o.cmd = nil
		return err
	}

	return err
}

// IsOpen returns wether the port is open
func (o *out) IsOpen() (open bool) {
	o.RLock()
	open = o.cmd != nil
	o.RUnlock()
	return
}

// Write writes a MIDI message to the MIDI output port
// If the output port is closed, it returns midi.ErrClosed
func (o *out) Write(b []byte) (int, error) {
	o.Lock()
	defer o.Unlock()
	if o.cmd == nil {
		return 0, midi.ErrPortClosed
	}

	_, err := fmt.Fprintf(o.wr, "%X\n", b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Underlying returns the underlying driver. Here it returns nil
func (o *out) Underlying() interface{} {
	return nil
}

// Number returns the number of the MIDI out port.
// Note that with rtmidi, out and in ports are counted separately.
// That means there might exists out ports and an in ports that share the same number
func (o *out) Number() int {
	return o.number
}

// String returns the name of the MIDI out port.
func (o *out) String() string {
	return o.name
}

// Close closes the MIDI out port
func (o *out) Close() (err error) {
	if !o.IsOpen() {
		return nil
	}

	o.Lock()
	defer o.Unlock()
	o.wr.Close()
	err = o.cmd.Process.Kill()
	o.cmd = nil
	o.rd.Close()
	o.wr = nil
	o.rd = nil
	return err
}

// Open opens the MIDI out port
func (o *out) Open() (err error) {
	if o.IsOpen() {
		return nil
	}

	err = o.fireCmd()

	if err != nil {
		return fmt.Errorf("can't open MIDI out port %v (%s): %v", o.number, o, err)
	}

	o.driver.Lock()
	o.driver.opened = append(o.driver.opened, o)
	o.driver.Unlock()

	return nil
}
