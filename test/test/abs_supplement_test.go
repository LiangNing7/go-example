package abs

import (
	"fmt"
	"os"
	"testing"
)

func setup() {
	fmt.Println("> setup completed")
}

func teardown() {
	fmt.Println("> teardown completed")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestAbs_setup(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	var want float64 = 1
	if got := Abs(-1); got != want {
		t.Fatalf("Abs() = %v, want %v", got, want)
	}
}

// testing.TB is the interface common to T, B, and F.
func setupTest(testing.TB) func(tb testing.TB) {
	fmt.Println(">> setup Test")

	return func(tb testing.TB) {
		fmt.Println(">> teardown test")
	}
}

func TestAbsWithTable(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "positive",
			args: args{x: 1},
			want: 1,
			// want: 2,
		},
		{
			name: "negative",
			args: args{x: -1},
			want: 1,
		},
	}
	for _, tt := range tests {
		teardownTest := setupTest(t)
		defer teardownTest(t)
		if got := Abs(tt.args.x); got != tt.want {
			t.Fatalf("Abs() = %v, want %v", got, tt.want)
		}
	}
}

func TestAbsWithTableAndSubtests(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "positive",
			args: args{x: 1},
			want: 2,
		},
		{
			name: "negative",
			args: args{x: -1},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownTest := setupTest(t)
			defer teardownTest(t)
			if got := Abs(tt.args.x); got != tt.want {
				t.Fatalf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
