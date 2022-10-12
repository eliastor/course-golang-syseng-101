package main

import (
	"io"

	"github.com/adamliesko/fakelog"
)

// Following code structure is widely sed in real go projects.
// It checks during compile time (or even in IDE with gopls enabled) that fakeLogCannon satisfy io.Reader, io.Closer and Server interfaces
// Satisfying some interfaces means that fakeLogCannon has the same definitions of all methods listed in specific interface.
// For example definion of io.Reader interface is:
// type Reader interface {
// 		Read(p []byte) (n int, err error)
// }
// fakeLogCannon has method Read(p []byte) (n int, err error)
var (
	_ io.Reader = (&fakeLogCannon{})
	_ io.Closer = (&fakeLogCannon{})
	// these two interfaces can be replaced by io.ReadClose with the same result
)

type fakeLogCannon struct {
	w *io.PipeWriter
	r *io.PipeReader

	//flogger contains logic for fake logs generation.
	flogger *fakelog.Logger

	// pipeClose is funcion that will close pipe. we'll use it in Close method
	pipeClose func()
}

// Read is method for satisfying io.Reader interface. It will read data from pipe, where it was pushed by fakelog.Logger instance
func (c *fakeLogCannon) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

// Close is method for satisfying io.Close interface and it will stop c.flogger and pipe.\
// After closing fakelogCannon cannot be used
func (c *fakeLogCannon) Close() error {
	c.flogger.Stop() // it causes c.flogger.GenerateLogs() to stop according documentation of fakelog.Logger: https://pkg.go.dev/github.com/adamliesko/fakelog?utm_source=godoc#Logger.Stop
	c.w.Close()      // Closing pipe so it will stop working and consuming resources

	// c.r.Close() is not needed because underlying logic of both io.pipeReader and io.pipeWriter Close() function closes one channel.
	//   you can dig into by clicking right-mouse-key on Close() and select "Go to definition"
	//   then going to CloseWithError() definion and to closeWrite*() definition

	return nil
}

// NewFakeLogGenerator return new instance of fakeLogCannon fully prepared for work.
func NewFakeLogGenerator() *fakeLogCannon {
	cannon := new(fakeLogCannon) // here we created pointer to fakeLogCannon. another way is: cannon := &fakeLogCannon{}

	cannon.r, cannon.w = io.Pipe()

	cannon.flogger = fakelog.NewLogger(fakelog.ApacheCommonLine, cannon.w, 200)
	cannon.pipeClose = func() {
		cannon.r.Close()
		// r.Close() is not needed. Take a look on internal of io.Pipe: both r and w Close methods close one done channel
	}
	go cannon.flogger.GenerateLogs()

	return cannon
}
