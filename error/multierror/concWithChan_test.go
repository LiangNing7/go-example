package conc_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/hashicorp/go-multierror"
	"gotest.tools/assert"
)

func TestConcurrencyWithChannel(t *testing.T) {
	var wg sync.WaitGroup
	errCh := make(chan error, 10000)
	errs := &multierror.Error{}

	for range 10000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- errors.New("error") // 通过 channel 安全发送 err
		}()
	}

	// 开启子 goroutine 等待并发程序执行完成
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// main goroutine 从 channel 收到 err 并完成聚合
	for err := range errCh {
		errs = multierror.Append(errs, err) // 单 goroutine 聚合，无竞争
	}

	// 预期 10000 个错误，实际输出也是 10000
	assert.Equal(t, 10000, len(errs.Errors)) // 测试通过
}
