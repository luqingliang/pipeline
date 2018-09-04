package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)
var startTime time.Time
func Init()  {
	startTime = time.Now()
}

func ArraySource(a ...int) <-chan int{
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMenSort(in <-chan int) <-chan int{
	out := make(chan int, 1024) //1024是不让他一个一个的发，提高效率
	go func() {
		//Read into memory
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read done : ", time.Now().Sub(startTime))

		//Sort
		sort.Ints(a)
		fmt.Println("InMemSort done : ", time.Now().Sub(startTime))

		//Output
		for _, v := range a {
			out <- v
		}

		close(out)
	}()
	return out
}
//归并节点
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 < v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge done : ", time.Now().Sub(startTime))
	}()
	return out
}

func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		//8byte是64字节，64位机器的int也是64位字节
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			//chunkSize = -1则表示可以无限读
			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}

		close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

//生成随机的数据源
func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

//归并多个节点(只能两两归并)
func MergeN(inputs ...<-chan int) <-chan int{
	//只有一个节点则直接返回
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}