package zapring_test

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
