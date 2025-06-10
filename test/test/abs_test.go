package abs

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// func TestAbs(t *testing.T) {
// 	 got := Abs(-1)
// 	 if got != 1 {
// 	 	t.Errorf("Abs(-1) = %f; want 1", got)
// 	 }
// }

// TestAbs Parallel.
func TestAbs(t *testing.T) {
	t.Parallel()
	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}

func TestAbs_TableDriver(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{
			name: "positive",
			x:    2,
			want: 2,
		},
		{
			name: "negative",
			x:    -3,
			want: 3,
		},
		// {
		// 	name: "negative",
		// 	x:    -3,
		// 	want: 33,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x); got != tt.want {
				t.Errorf("Abs(%f) = %v; want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestAbs_Skip(t *testing.T) {
	// CI 环境下跳过当前测试.
	if os.Getenv("CI") != "" {
		t.Skip("it's too slow, skip when running in CI")
	}

	t.Log(t.Skipped())

	got := Abs(-2)
	if got != 2 {
		t.Errorf("Abs(-2) = %f; want 2", got)
	}
}

func TestAbs_Parallel(t *testing.T) {
	t.Log("Parallel before")
	// 标记当前测试支持并行.
	t.Parallel()
	t.Log("Parallel after")

	got := Abs(2)
	if got != 2 {
		t.Errorf("Abs(2) = %f; want 2", got)
	}
}

func BenchmarkAbs(b *testing.B) {
	for b.Loop() {
		Abs(-1)
	}
}

func BenchmarkAbsResetTimer(b *testing.B) {
	time.Sleep(100 * time.Millisecond) // 模拟耗时的准备工作
	b.ResetTimer()
	for b.Loop() {
		Abs(-1)
	}
}

func BenchmarkAbsParallel(b *testing.B) {
	b.SetParallelism(2) // 设置并发 Goroutine 数量为 2 * GOMAXPROCS
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Abs(-1)
		}
	})
}

func ExampleAbs() {
	fmt.Println(Abs(-1))
	fmt.Println(Abs(-2))
	// Output:
	// 1
	// 2
}

func ExampleAbs_unordered() {
	fmt.Println(Abs(2))
	fmt.Println(Abs(-1))
	// Unordered Output:
	// 1
	// 2
}
