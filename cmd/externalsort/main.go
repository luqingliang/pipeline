package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"supreulu/gointro/pipeline"
)

func main() {
	p := createNetworkPipeline("small.in", 800000000, 4)
	writeToFIle(p, "small.out")
	printFile("small.out")
}

func printFile(filename string)  {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count ++
		if count >= 100 {
			break
		}
	}
}

func writeToFIle(p <-chan int, filename string)  {
	file,err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)
}

func createPipeline(filename string, fileSize, chunkCount int) <-chan int{
	//分块去读，默认传入的值刚好可以整除
	chunkSize := fileSize / chunkCount

	//开始计时，并打印日志
	pipeline.Init()

	sortResults := []<-chan int{}
	for i := 0; i < chunkCount; i ++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err) //不知道该怎么解决索性让程序终止掉
		}
		file.Seek(int64(i * chunkSize), 0) //分块读取
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMenSort(source))
	}

	return pipeline.MergeN(sortResults...)
}

/**
 * 网络版
 */
func createNetworkPipeline(filename string, fileSize, chunkCount int) <-chan int{
	//分块去读，默认传入的值刚好可以整除
	chunkSize := fileSize / chunkCount

	//开始计时，并打印日志
	pipeline.Init()

	sortAddr := []string {}
	for i := 0; i < chunkCount; i ++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err) //不知道该怎么解决索性让程序终止掉
		}
		file.Seek(int64(i * chunkSize), 0) //分块读取
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		addr := ":" + strconv.Itoa(7000 + i)
		pipeline.NetworkSink(addr, pipeline.InMenSort(source))
		sortAddr = append(sortAddr, addr)
	}
	sortResults := []<-chan int{}
	for _, v := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(v))
	}
	return pipeline.MergeN(sortResults...)
}