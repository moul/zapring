package zapring

import (
	"io"
	"io/ioutil"
	"sync"

	"github.com/maruel/circular"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Core is an in-memory ring buffer log that implements zapcore.Core.
type Core struct {
	zapcore.Core
	enc       zapcore.Encoder
	buffer    circular.Buffer
	setupOnce sync.Once
}

// New returns a ring-buffer with a capacity of 'size' bytes.
func New(size uint) *Core {
	return &Core{
		buffer: circular.New(int(size)),
	}
}

// Close implements zapcore.Core.
func (c *Core) Close() {
	c.buffer.Flush() // ensures all readers have caught up.
	c.buffer.Close() // gracefully closes the readers.
}

// setup sets default value for core and encoder if they are empty.
func (c *Core) setup() {
	c.setupOnce.Do(func() {
		if c.Core != nil && c.enc != nil {
			return
		}

		encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
		if c.enc == nil {
			c.enc = encoder
		}
		if c.Core == nil {
			discardCore := zapcore.NewCore(
				encoder,
				zapcore.AddSync(ioutil.Discard),
				zap.LevelEnablerFunc(func(_ zapcore.Level) bool { return true }),
			)
			c.Core = discardCore
		}
	})
}

// Enabled implements zapcore.LevelEnabler.
func (c *Core) Enabled(level zapcore.Level) bool {
	c.setup()
	return c.Core.Enabled(level)
}

func (c *Core) clone() *Core {
	return &Core{
		buffer: c.buffer,
		enc:    c.enc.Clone(),
		Core:   c.Core,
	}
}

// Sync implements zapcore.Core.
func (c *Core) Sync() error {
	if c.Core != nil {
		return c.Core.Sync()
	}
	return nil
}

// With implements zapcore.Core.
func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for _, field := range fields {
		field.AddTo(clone.enc)
	}
	return clone
}

// Check implements zapcore.Core.
func (c *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	c.setup()
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

// Write implements zapcore.Core.
func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	c.setup()
	// We need to be sure we're in 64-bit since circular buffer call LoadInt64
	// so if you're running in 32-bit the app will crash
	if 32<<(^uintptr(0)>>63) == 64 { // nolint:gomnd
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

func (c *Core) SetNextCore(core zapcore.Core) *Core {
	c.Core = core
	return c
}

func (c *Core) SetEncoder(enc zapcore.Encoder) *Core {
	c.enc = enc
	return c
}
