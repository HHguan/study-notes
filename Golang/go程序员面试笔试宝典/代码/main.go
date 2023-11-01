package main

import (
	"fmt"
)

func main() {
	chanNum := 3 // 定义通道数量为3

	// 创建通道队列切片
	chanQueue := make([]chan struct{}, chanNum)

	var result = 0                  // 定义计数器变量
	exitChan := make(chan struct{}) // 创建退出信号通道

	// 初始化通道队列
	for i := 0; i < chanNum; i++ {
		chanQueue[i] = make(chan struct{}) // 创建每个通道
		if i == chanNum-1 {
			go func(i int) {
				chanQueue[i] <- struct{}{} // 启动最后一个协程，发送通道信号
			}(i)
		}
	}

	// 启动三个协程
	for i := 0; i < chanNum; i++ {
		var lastChan, curChan chan struct{}
		if i == 0 {
			lastChan = chanQueue[chanNum-1] // 获取上一个协程的通道
		} else {
			lastChan = chanQueue[i-1]
		}
		curChan = chanQueue[i] // 获取当前协程的通道

		go func(i byte, lastChan, curChan chan struct{}) {
			for {
				<-lastChan // 等待上一个协程的通道信号
				if result >= 10 {
					exitChan <- struct{}{} // 如果已经打印完 1-10，则发送退出信号
					return
				}
				result++
				fmt.Println(result)   // 打印数字
				curChan <- struct{}{} // 发送通道信号给下一个协程
			}
		}('A'+byte(i), lastChan, curChan)
	}

	<-exitChan          // 等待退出信号
	fmt.Println("done") // 打印 "done"
}
