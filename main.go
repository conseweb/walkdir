package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func walkFiles(root string) (chan string, chan error) {
	paths := make(chan string, 100)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			}
			return nil
		})
	}()
	return paths, errc
}

type result struct {
	path string
	//sha1 [sha1.Size]byte
	digest hash.Hash
	size   int64
}

func fileIo(paths chan string, c chan result) {
	for path := range paths { // filenamepaths
		fileinfo, err := os.Lstat(path)
		if err != nil {
			fmt.Println("读取%v信息错误，err:%v", path, err)
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("open %v err: %v", path, err)
			continue
		}
		defer f.Close()
		h := sha1.New()
		_, err1 := io.Copy(h, f)
		if err1 != nil {
			fmt.Println("io.copy err: %v", err1)
			continue
		}
		select {
		case c <- result{path, h, fileinfo.Size()}:
		}
	}
}

func Sha1(c chan result, l chan string) {
	for result := range c { //io数据
		line := fmt.Sprintf("%v,%x,%v\n", result.path, result.digest.Sum(nil), result.size)
		select {
		case l <- line:
		}
	}
}

func Sha1All(root string) (chan string, error) {
	paths, errc := walkFiles(root)

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result, 100) //c中存放的是未经计算的io数据
	var wgIo sync.WaitGroup

	const numIos = 1 //固态硬盘20,机械硬盘1
	wgIo.Add(numIos)
	for i := 0; i < numIos; i++ {
		go func() {
			fileIo(paths, c)
			wgIo.Done()
		}()
	}

	go func() {
		wgIo.Wait()
		close(c)
	}()

	l := make(chan string, 100) //l中存放的是最终文件每行的内容
	var wgSha1 sync.WaitGroup

	const numSha1s = 20
	wgSha1.Add(numSha1s)
	for i := 0; i < numSha1s; i++ {
		go func() {
			Sha1(c, l)
			wgSha1.Done()
		}()
	}

	go func() {
		wgSha1.Wait()
		close(l)
	}()

	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return nil, err
	}
	return l, nil
}

func main() {
	// Calculate the Sha1 sum of all files under the specified directory,
	// then create file include all sha1
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	input := flag.String("input", "", "root directory")
	output := flag.String("output", "", "result file")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("参数           默认值                  说明")
		fmt.Println("-input        ./                            root directory")
		fmt.Println("-output     ../dthashresult      result file")
		return
	}
	fi, err1 := os.Lstat(*input)
	if err1 != nil {
		panic(err1)
	}
	if !fi.IsDir() {
		panic(fmt.Sprintf("%v不是目录", *input))
	}

	f, err1 := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err1 != nil {
		panic(err1)
	}
	defer f.Close()

	l, err2 := Sha1All(*input)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	for line := range l {
		f.WriteString(line)
	}

}
