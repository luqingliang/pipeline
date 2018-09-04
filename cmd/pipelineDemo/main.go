package main

import (
	"bufio"
	"fmt"
	"os"
	"supreulu/gointro/pipeline"
)

func main() {
	const filename = "small.in"
	const n = 100000000

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := pipeline.RandomSource(n)
	writer := bufio.NewWriter(file) //没有buffer运行效率会慢很多
	pipeline.WriterSink(writer, p)
	writer.Flush() //不完全倾倒出去生成的文件会少一点

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count ++
		if count >= 100 {
			break
		}
	}
}

func mergeDemo()  {
	p := pipeline.Merge(pipeline.InMenSort(pipeline.ArraySource(3, 2, 5, 6, 1)),
		pipeline.InMenSort(pipeline.ArraySource(9, 4, 2, 6, 7, 10, 2, 18)))

	for v :=  range p {
		fmt.Println(v)
	}
}