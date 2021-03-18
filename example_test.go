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

func Example() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "" // used to make this test consistent (not depending on current timestamp)
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	level := zap.LevelEnablerFunc(func(_ zapcore.Level) bool { return true })
	ring := zapring.New(uint(10 * 1024 * 1024)) // 10Mb ring
	defer ring.Close()
	core := ring.Wrap(
		zapcore.NewCore(encoder, zapcore.AddSync(ioutil.Discard), level),
		encoder,
	)
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
	// -->  {"L":"INFO","C":"zapring/example_test.go:31","M":"hello world!"}
	// -->  {"L":"INFO","C":"zapring/example_test.go:32","M":"lorem ipsum"}
}
