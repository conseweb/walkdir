package dthash

import (
	"crypto/sha1"
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
			fmt.Printf("读取%v信息错误，err:%v", path, err)
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("open %v err: %v\n", path, err)
			f.Close()
			continue
		}

		h := sha1.New()
		_, err1 := io.Copy(h, f)
		if err1 != nil {
			fmt.Printf("io.copy err: %v", err1)
			f.Close()
			continue
		}
		f.Close()
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

func Sha1All(root string) (chan string, chan error) {
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

	return l, errc
}
