# zapring

:smile: zapring

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/moul.io/zapring)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/moul/zapring/blob/master/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/moul/zapring.svg)](https://github.com/moul/zapring/releases)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/zapring.svg)](https://microbadger.com/images/moul/zapring)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)

[![Go](https://github.com/moul/zapring/workflows/Go/badge.svg)](https://github.com/moul/zapring/actions?query=workflow%3AGo)
[![Release](https://github.com/moul/zapring/workflows/Release/badge.svg)](https://github.com/moul/zapring/actions?query=workflow%3ARelease)
[![PR](https://github.com/moul/zapring/workflows/PR/badge.svg)](https://github.com/moul/zapring/actions?query=workflow%3APR)
[![GolangCI](https://golangci.com/badges/github.com/moul/zapring.svg)](https://golangci.com/r/github.com/moul/zapring)
[![codecov](https://codecov.io/gh/moul/zapring/branch/master/graph/badge.svg)](https://codecov.io/gh/moul/zapring)
[![Go Report Card](https://goreportcard.com/badge/moul.io/zapring)](https://goreportcard.com/report/moul.io/zapring)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/zapring/badge)](https://www.codefactor.io/repository/github/moul/zapring)

[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/moul/zapring)

## Usage

[embedmd]:# (example_test.go /import\ / $)
```go
import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/zapring"
)

func Example_custom() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "" // used to make this test consistent (not depending on current timestamp)
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	level := zap.LevelEnablerFunc(func(_ zapcore.Level) bool { return true })
	ring := zapring.New(uint(10 * 1024 * 1024)) // 10Mb ring
	defer ring.Close()
	core := ring.
		SetNextCore(zapcore.NewCore(encoder, zapcore.AddSync(ioutil.Discard), level)).
		SetEncoder(encoder)
	logger := zap.New(
		core,
		zap.Development(),
		zap.AddCaller(),
	)
	defer logger.Sync()
	logger.Info("hello world!")
	logger.Info("lorem ipsum")

	r, w := io.Pipe()
	go func() {
		_, err := ring.WriteTo(w)
		if err != nil && err != io.EOF {
			panic(err)
		}
		w.Close()
	}()
	scanner := bufio.NewScanner(r)
	lines := 0
	for scanner.Scan() {
		fmt.Println("--> ", scanner.Text())
		lines++
		if lines == 2 {
			break
		}
	}

	// Output:
	// -->  {"L":"INFO","C":"zapring/example_test.go:30","M":"hello world!"}
	// -->  {"L":"INFO","C":"zapring/example_test.go:31","M":"lorem ipsum"}
}

func Example_composite() {
	cli := zap.NewExample()
	cli.Info("hello cli!")
	ring := zapring.New(10 * 1024 * 1024) // 10MB ring-buffer
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "" // used to make this test consistent (not depending on current timestamp)
	ring.SetEncoder(zapcore.NewJSONEncoder(encoderConfig))
	// FIXME: ring.Info("hello ring!")
	composite := zap.New(
		zapcore.NewTee(cli.Core(), ring),
		zap.Development(),
	)
	composite.Info("hello composite!")

	r, w := io.Pipe()
	go func() {
		_, err := ring.WriteTo(w)
		if err != nil && err != io.EOF {
			panic(err)
		}
		w.Close()
	}()
	composite.Info("hello composite 2!")
	cli.Info("hello cli 2!")
	scanner := bufio.NewScanner(r)
	lines := 0
	for scanner.Scan() {
		fmt.Println("-> ", scanner.Text())
		lines++
		if lines == 2 {
			break
		}
	}

	// Output:
	// {"level":"info","msg":"hello cli!"}
	// {"level":"info","msg":"hello composite!"}
	// {"level":"info","msg":"hello composite 2!"}
	// {"level":"info","msg":"hello cli 2!"}
	// ->  {"L":"INFO","M":"hello composite!"}
	// ->  {"L":"INFO","M":"hello composite 2!"}
}

func Example_simple() {
	ring := zapring.New(10 * 1024 * 1024) // 10MB ring-buffer
	logger := zap.New(ring, zap.Development())
	logger.Info("test")
	// Output:
}
```

[embedmd]:# (.tmp/usage.txt txt /TYPES/ $)
```txt
TYPES

type Core struct {
	zapcore.Core

	// Has unexported fields.
}
    Core is an in-memory ring buffer log that implements zapcore.Core.

func New(size uint) *Core
    New returns a ring-buffer with a capacity of 'size' bytes.

func (c *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry
    Check implements zapcore.Core.

func (c *Core) Close()
    Close implements zapcore.Core.

func (c *Core) Enabled(level zapcore.Level) bool
    Enabled implements zapcore.LevelEnabler.

func (c *Core) SetEncoder(enc zapcore.Encoder) *Core

func (c *Core) SetNextCore(core zapcore.Core) *Core

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error
    Write implements zapcore.Core.

func (c *Core) WriteTo(w io.Writer) (n int64, err error)
    WriteTo implements io.WriterTo.

```

## Install

### Using go

```sh
go get moul.io/zapring
```

### Releases

See https://github.com/moul/zapring/releases

## Contribute

![Contribute <3](https://raw.githubusercontent.com/moul/moul/master/contribute.gif)

I really welcome contributions.
Your input is the most precious material.
I'm well aware of that and I thank you in advance.
Everyone is encouraged to look at what they can do on their own scale;
no effort is too small.

Everything on contribution is sum up here: [CONTRIBUTING.md](./CONTRIBUTING.md)

### Contributors ✨

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://manfred.life"><img src="https://avatars1.githubusercontent.com/u/94029?v=4" width="100px;" alt=""/><br /><sub><b>Manfred Touron</b></sub></a><br /><a href="#maintenance-moul" title="Maintenance">🚧</a> <a href="https://github.com/moul/zapring/commits?author=moul" title="Documentation">📖</a> <a href="https://github.com/moul/zapring/commits?author=moul" title="Tests">⚠️</a> <a href="https://github.com/moul/zapring/commits?author=moul" title="Code">💻</a></td>
    <td align="center"><a href="https://manfred.life/moul-bot"><img src="https://avatars1.githubusercontent.com/u/41326314?v=4" width="100px;" alt=""/><br /><sub><b>moul-bot</b></sub></a><br /><a href="#maintenance-moul-bot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors)
specification. Contributions of any kind welcome!

### Stargazers over time

[![Stargazers over time](https://starchart.cc/moul/zapring.svg)](https://starchart.cc/moul/zapring)

## License

© 2021   [Manfred Touron](https://manfred.life)

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT)
([`LICENSE-MIT`](LICENSE-MIT)), at your option.
See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`
