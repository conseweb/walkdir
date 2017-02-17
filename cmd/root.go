// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var cfgFile string
var input string
var output string

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

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dthash",
	Short: "生成指定目录下所有文件的sha1值",
	Long:  `将指定目录下所有文件的sha1值计算，写入一个文件.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Println("计算sha1开始", time.Now().Format("2006-01-02 15:04:05"))
		l, err2 := Sha1All(input)
		for line := range l {
			f.WriteString(line)
		}
		// Check whether the Walk failed.
		if err := <-err2; err != nil { // HLerrc
			fmt.Println(err)
			return
		}
		fmt.Println("计算sha1结束", time.Now().Format("2006-01-02 15:04:05"))
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dthash.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&input, "input", "i", "./", "your directory ")
	RootCmd.Flags().StringVarP(&output, "output", "o", "../dthashresult", "result file path")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".dthash") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
