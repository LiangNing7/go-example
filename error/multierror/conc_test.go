package conc_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/hashicorp/go-multierror"
	"gotest.tools/assert"
)

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	multierrs := &multierror.Error{}

	for range 10000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			multierrs = multierror.Append(multierrs, errors.New("error"))
		}()
	}

	wg.Wait()
	// 预期 10000 个错误，实际输出可能 < 10000
	assert.Equal(t, 10000, len(multierrs.Errors)) // 测试失败
}
