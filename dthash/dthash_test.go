package dthash

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

func TestWalkFiles(t *testing.T) {
	root := "/"
	fmt.Println("walkFiles开始", time.Now().String())
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		match, _ := regexp.MatchString("^/proc/.*|^/sys/.*|^/run/.*", path)
		if match {
			return nil
		}

		//fmt.Println(path)

		return nil
	})
	fmt.Println("walkFiles结束", time.Now().String())
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
	input := "/"
	output := "/home/zhaolong/temp/dthash"
	fi, err1 := os.Lstat(input)
	if err1 != nil {
		panic(err1)
	}
	if !fi.IsDir() {
		panic(fmt.Sprintf("%v不是目录", input))
	}

	f, err1 := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("计算sha1开始", time.Now().String())
	l, err2 := Sha1All(input)
	for line := range l {
		f.WriteString(line)
	}
	// Check whether the Walk failed.
	if err := <-err2; err != nil { // HLerrc
		fmt.Println(err)
		return
	}
	fmt.Println("计算sha1结束", time.Now().String())
}

func BenchmarkExecute(b *testing.B) {
	fmt.Println("B.N =", b.N)
	for i := 0; i < b.N; i++ {
		input := "/"
		output := "/home/zhaolong/temp/dthash"
		fi, err1 := os.Lstat(input)
		if err1 != nil {
			panic(err1)
		}
		if !fi.IsDir() {
			panic(fmt.Sprintf("%v不是目录", input))
		}

		f, err1 := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		defer f.Close()
		if err1 != nil {
			panic(err1)
		}

		fmt.Println("计算sha1开始", time.Now().String())
		l, err2 := Sha1All(input)
		for line := range l {
			f.WriteString(line)
		}
		// Check whether the Walk failed.
		if err := <-err2; err != nil { // HLerrc
			fmt.Println(err)
			return
		}
		fmt.Println("计算sha1结束", time.Now().String())
	}
}
