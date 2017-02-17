package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestWalkDir(t *testing.T) {
	input := "/home/zhaolong/temp/input"
	p = NewPool(runtime.NumCPU()*2, 100, 100)
	p.Add(1)
	go WalkDir(input)
	p.Wait()
}

func TestGetFileSha1(t *testing.T) {
	file := "/home/zhaolong/temp/input/a.txt"
	sha1 := GetFileSha1(file)
	fmt.Printf("%v的sha1哈希值是%v", file, sha1)
}

func TestDealFileName(t *testing.T) {
	p = NewPool(runtime.NumCPU()*2, 100, 100)
	p.WalkDirDone()
	p.fn <- "/home/zhaolong/temp/input/a.txt"
	p.fn <- "/home/zhaolong/temp/input/b.txt"
	p.fn <- "/home/zhaolong/temp/input/sub/c.txt"
	p.fn <- "/home/zhaolong/temp/input/sub/sub1/d.txt"
	p.fn <- "/home/zhaolong/temp/input/sub/sub1/sub2/e.txt"
	p.fn <- "/home/zhaolong/temp/input/sub/sub1/sub2/sub3/f.txt"
	DealFileName()
	p.Wait()
}

func TestMakeResultFile(t *testing.T) {
	output := "/home/zhaolong/temp/dthash"
	p = NewPool(runtime.NumCPU()*2, 100, 100)
	p.dfnd = p.dfnn
	p.rs <- "c.txt,c585520b4415e99bf4e94345a2830b8e29bf4833,137\n"
	p.rs <- "a.txt,80607f6a6597c98b5409711d525c2cd2625e7cd0,94\n"
	p.rs <- "b.txt,e1432e814307c7b0d532e903ae7d22c578396a5b,150\n"
	p.rs <- "d.txt,c7d4dae5e554951106b453891220cf619ef297c1,191\n"
	p.rs <- "e.txt,65a138edd2a8cafcfc2405284d18019e884ba1ee,161\n"
	p.rs <- "f.txt,a6926a641d9cb08d81946eb1ed5725efb258e573,107\n"
	p.Add(1)
	go MakeResultFile(output)
	p.Wait()
}

func BenchmarkWalkDir(b *testing.B) {
	fmt.Println("B.N =", b.N)
	for i := 0; i < b.N; i++ {
		input := "/home/zhaolong/temp/input"
		p = NewPool(runtime.NumCPU()*2, 100, 100)
		p.Add(1)
		go WalkDir(input)
		p.Wait()
	}
}

func BenchmarkGetFileSha1(b *testing.B) {
	fmt.Println("B.N =", b.N)
	for i := 0; i < b.N; i++ {
		file := "/home/zhaolong/temp/input/a.txt"
		GetFileSha1(file)

	}
}

func BenchmarkDealFileName(b *testing.B) {
	fmt.Println("B.N =", b.N)
	for i := 0; i < b.N; i++ {
		p = NewPool(runtime.NumCPU()*2, 100, 100)
		p.WalkDirDone()
		p.fn <- "/home/zhaolong/temp/input/a.txt"
		p.fn <- "/home/zhaolong/temp/input/b.txt"
		p.fn <- "/home/zhaolong/temp/input/sub/c.txt"
		p.fn <- "/home/zhaolong/temp/input/sub/sub1/d.txt"
		p.fn <- "/home/zhaolong/temp/input/sub/sub1/sub2/e.txt"
		p.fn <- "/home/zhaolong/temp/input/sub/sub1/sub2/sub3/f.txt"
		DealFileName()
		p.Wait()

	}
}

func BenchmarkMakeResultFile(b *testing.B) {
	fmt.Println("B.N =", b.N)
	for i := 0; i < b.N; i++ {
		output := "/home/zhaolong/temp/dthash"
		p = NewPool(runtime.NumCPU()*2, 100, 100)
		p.dfnd = p.dfnn
		p.rs <- "c.txt,c585520b4415e99bf4e94345a2830b8e29bf4833,137\n"
		p.rs <- "a.txt,80607f6a6597c98b5409711d525c2cd2625e7cd0,94\n"
		p.rs <- "b.txt,e1432e814307c7b0d532e903ae7d22c578396a5b,150\n"
		p.rs <- "d.txt,c7d4dae5e554951106b453891220cf619ef297c1,191\n"
		p.rs <- "e.txt,65a138edd2a8cafcfc2405284d18019e884ba1ee,161\n"
		p.rs <- "f.txt,a6926a641d9cb08d81946eb1ed5725efb258e573,107\n"
		p.Add(1)
		go MakeResultFile(output)
		p.Wait()

	}
}
