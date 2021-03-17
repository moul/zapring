package zapring

import (
	"io"

	"github.com/maruel/circular"
	"go.uber.org/zap/zapcore"
)

const ptrSize = 32 << (^uintptr(0) >> 63) // nolint:gomnd

// Core is an in-memory ring buffer log that implements zapcore.Core.
type Core struct {
	zapcore.Core
	enc    zapcore.Encoder
	buffer circular.Buffer
}

func New(size uint) *Core {
	return &Core{
		buffer: circular.New(int(size)),
	}
}

func (c *Core) Close() {
	c.buffer.Flush() // ensures all readers have caught up.
	c.buffer.Close() // gracefully closes the readers.
}

func (c *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// We need to be sure we're in 64-bit since circular buffer call LoadInt64
	// so if you're running in 32-bit the app will crash
	if ptrSize == 64 { // nolint:gomnd
		buff, err := c.enc.EncodeEntry(entry, fields)
		if err != nil {
			return err
		}
		if _, err = c.buffer.Write(buff.Bytes()); err != nil {
			return err
		}
	}

	// FIXME: add support for 32-bit systems

	return c.Core.Write(entry, fields)
}

// WriteTo implements io.WriterTo.
func (c *Core) WriteTo(w io.Writer) (n int64, err error) {
	return c.buffer.WriteTo(w)
}

func (c *Core) Wrap(core zapcore.Core, enc zapcore.Encoder) zapcore.Core {
	c.Core = core
	c.enc = enc
	return c
}
