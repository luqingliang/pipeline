package main

import (
	"fmt"
	"strings"
	"time"
)

/**
 * 首先定义两个接口，规范方法的实现
 */
type Reader interface {
	read(rc chan string)
}
type Writer interface {
	write(wc chan string)
}

/**
 * 然后定义两个结构体模块分别负责不同的工作
 */
type ReadFromFile struct {
	path string
}
type WriteToFile struct {
	path string
}

/**
 * 定义我们的抽象类的接口体，并以参数的形式注入两个功能模块
 */
type ClassDemo struct {
	rc chan string
	wc chan string
	read Reader
	write Writer
}

/**
 * 分别实现两个模块的接口，功能逻辑的实现
 */
func (r *ReadFromFile) read(rc chan string)  { //实现read接口
	line := "message"
	rc <- line
}
func (w *WriteToFile) write(wc chan string)  { //实现写入结构体的write接口
	fmt.Println(<-wc)
}

/**
 * 我们的抽象类的私有成员函数的实现
 */
func (c *ClassDemo) jiexiModel()  {
	m := <-c.rc
	c.wc <- strings.ToUpper(m)
}

func main() {
	//首先实例化两个模块
	r := &ReadFromFile{
		path: "user/path",
	}
	w := &WriteToFile{
		path: "supreme/path",
	}

	//然后实例化我们的类
	c := &ClassDemo{
		rc: make(chan string),
		wc: make(chan string),
		read: r,
		write: w,
	}

	//开三个协程去做这件事情，go的协程不用考虑等待或者锁之类的东西
	go c.read.read(c.rc)
	go c.jiexiModel()
	go c.write.write(c.wc)

	time.Sleep(1 * time.Second)
}
