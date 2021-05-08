package pool

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

type st struct {
	Cid   int
	Error error
}

type TestPro struct {
	msgContent string
}

// 实现发送者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

// 实现接收者
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return nil
}

//定义携程池类型
type Pool struct {
	//对外接收Task的入口
	CommonTaskChan chan *CommonTask
	//协程池最大工作数量,限定Goroutine的个数
	ChannelNum int
	//协程池内部的任务就绪队列
	JobChannels chan *CommonTask
}

//定义任务CommonTask类型,每一个任务CommonTask都可以抽象成一个函数
type CommonTask struct {
	thread func() error //一个无参的函数类型
}

//通过CreateTask来创建一个CommonTask
func CreateTask(f func() error) *CommonTask {
	t := CommonTask{
		thread: f,
	}
	return &t
}

//执行CommonTask任务的方法
func (t *CommonTask) Execute() {
	err := t.thread() //调用任务所绑定的函数
	if err != nil {
		fmt.Println("Execute err", err.Error())
	}
}

//创建一个协程池
func NewGoPool(cap int) *Pool {
	p := Pool{
		CommonTaskChan: make(chan *CommonTask),
		ChannelNum:     cap,
		JobChannels:    make(chan *CommonTask),
	}

	return &p
}

//协程池创建一个协程并且开始工作
func (p *Pool) start(taskId int) {
	//worker不断的从JobChannels内部任务队列中拿任务
	for task := range p.JobChannels {
		//如果拿到任务,则执行任务
		task.Execute()
		fmt.Println("taskId ID ", taskId, " 执行完毕任务")
	}
}

//让协程池Pool开始工作
func (p *Pool) Run() {
	//1,首先根据协程池的协程数量限定,开启固定数量的协程,
	for i := 0; i < p.ChannelNum; i++ {
		go p.start(i)
	}

	// 2、从协程池入口取外界传递过来的任务并且将任务送进JobChannels中
	for task := range p.CommonTaskChan {
		p.JobChannels <- task
	}

	//3, 执行完毕需要关闭JobChannels
	close(p.JobChannels)

	//4, 执行完毕需要关闭
	close(p.CommonTaskChan)
}

func CallDemo() {
	//创建一个Task
	t := CreateTask(func() error {
		fmt.Println(time.Now())
		return nil
	})

	//创建一个协程池,最大开启3个协程worker
	p := NewGoPool(3)

	//开一个协程 不断的向 Pool 输送打印一条时间的task任务
	go func() {
		for {
			p.CommonTaskChan <- t
			time.Sleep(10 * time.Second)
		}
	}()
	//启动协程池p
	p.Run()
}
