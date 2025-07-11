package main

import (
	"embed"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 解决全局变量依赖.
func init() {
	version = "v0.1.1"
	commit = "1"
	builtGoVersion = "1.20.1"
}

func TestGetChangeLog(t *testing.T) {
	// 创建临时文件.
	// 第一个参数传 ""，表示在操作系统的临时目录下创建该文件.
	// 文件文件名会议第二个参数为前缀，剩余部分自动生成，以确保并发调用时生成的文件名不重复.
	f, err := os.CreateTemp("", "TEST_CHANGLOG")
	assert.NoError(t, err)
	defer func() {
		_ = f.Close()
		// 尽管操作系统会在某个时间自动清理临时文件，但主动清理是创建者的责任.
		_ = os.RemoveAll(f.Name())
	}()

	changeLogPath = f.Name()

	data := `
# Changelog
All notable changes to this project will be documented in this file.
`

	_, err = f.WriteString(data)
	assert.NoError(t, err)

	expected := ChangeLogSpec{
		Version:        "v0.1.1",
		Commit:         "1",
		BuiltGoVersion: "1.20.1",
		ChangeLog: `
# Changelog
All notable changes to this project will be documented in this file.
`,
	}

	actual, err := GetChangeLog()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

//go:embed testdata/CHANGELOG.md
var changeLog []byte

func TestGetChangeLog_by_embed(t *testing.T) {
	f, err := os.CreateTemp("", "TEST_CHANGLOG")
	assert.NoError(t, err)
	defer func() {
		_ = f.Close()
		_ = os.RemoveAll(f.Name())
	}()

	changeLogPath = f.Name()

	_, err = f.Write(changeLog)
	assert.NoError(t, err)

	expected := ChangeLogSpec{
		Version:        "v0.1.1",
		Commit:         "1",
		BuiltGoVersion: "1.20.1",
		ChangeLog:      string(changeLog),
	}

	actual, err := GetChangeLog()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

//go:embed testdata/CHANGELOG.md
var fs embed.FS

func TestGetChangeLogByIOReader(t *testing.T) {
	f, err := fs.Open("testdata/CHANGELOG.md")
	assert.NoError(t, err)

	data, err := io.ReadAll(f)
	assert.NoError(t, err)

	// 上面读取一遍后，指针已经走到最后，但还需要将 f 传递给 GetChangeLog, 所以这里重置一下.
	// 将数据的读取位置重置到开头.
	_, err = f.(io.ReadSeeker).Seek(0, 0)
	assert.NoError(t, err)

	expected := ChangeLogSpec{
		Version:        "v0.1.1",
		Commit:         "1",
		BuiltGoVersion: "1.20.1",
		ChangeLog:      string(data),
	}

	actual, err := GetChangeLogByIOReader(f)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
