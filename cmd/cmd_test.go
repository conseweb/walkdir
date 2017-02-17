package cmd

import (
	"fmt"
	"testing"
)

func TestWalkFiles(t *testing.T) {
	input := "/home/zhaolong/temp/input"
	walkFiles(input)
}

func TestFileIo(t *testing.T) {
	path := make(chan string, 100)
	c := make(chan result, 100) //c中存放的是未经计算的io数据
	path <- "/home/zhaolong/temp/input/a.txt"
	path <- "/home/zhaolong/temp/input/b.txt"
	path <- "/home/zhaolong/temp/input/sub/c.txt"
	path <- "/home/zhaolong/temp/input/sub/sub1/d.txt"
	path <- "/home/zhaolong/temp/input/sub/sub1/sub2/e.txt"
	path <- "/home/zhaolong/temp/input/sub/sub1/sub2/sub3/f.txt"
	close(path)
	fileIo(path, c)
}

func TestSha1(t *testing.T) {
	path := make(chan string, 100)
	c := make(chan result, 100) //c中存放的是未经计算的io数据
	l := make(chan string, 100) //l中存放的是最终文件每行的内容
	path <- "/home/zhaolong/temp/input/a.txt"
	path <- "/home/zhaolong/temp/input/b.txt"
	path <- "/home/zhaolong/temp/input/sub/c.txt"
	path <- "/home/zhaolong/temp/input/sub/sub1/d.txt"
	path <- "/home/zhaolong/temp/input/sub/sub1/sub2/e.txt"
	path <- "/home/zhaolong/temp/input/sub/sub1/sub2/sub3/f.txt"
	close(path)
	fileIo(path, c)
	close(c)
	Sha1(c, l)
}

func TestSha1All(t *testing.T) {
	input := "/home/zhaolong/temp/input"
	Sha1All(input)
}

func TestExecute(t *testing.T) {
	Execute()
}

func BenchmarkExecute(b *testing.B) {
	fmt.Println("B.N =", b.N)
	for i := 0; i < b.N; i++ {
		Execute()
	}
}
