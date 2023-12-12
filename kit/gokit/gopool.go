// 生产消费模型的协程池
package gokit

type taskFun = func(int) error

// 定义池
type Pool struct {
	// 对外接收 task 的入口，使用管道保存 task
	EntryChannel chan *Task

	// 协程池最大 worker 数量，即限定 Goroutine 的个数
	workerNum int

	// 协程池内部的任务就绪队列
	TasksChannel chan *Task
}

// 定义 task，即要开协程去做的事情，每一个 task 都可以抽象成一个函数
type Task struct {
	fn taskFun
}

// 创建池子
// cap 为池子容量，即最大开启多少个协程 worker
func NewPool(cap int) *Pool {
	pool := Pool{
		EntryChannel: make(chan *Task),
		TasksChannel: make(chan *Task),
		workerNum:    cap,
	}
	return &pool
}

// 创建 task
// fun 为具体要开协程去执行的函数
func NewTask(fun taskFun) *Task {
	task := Task{
		fn: fun,
	}
	return &task
}

// 启动协程池进行工作
func (pool *Pool) Run() {
	// 首先根据协程池最打容量开启对应数量的 worker，每一个 worker 用一个 goroutine 承载
	for i := 0; i < pool.workerNum; i++ {
		go pool.worker(i)
	}

	// 从 EntryChannel 取出外界传递过来的任务，然后将任务送进 TasksChannel 中
	for task := range pool.EntryChannel {
		pool.TasksChannel <- task
	}

	// 执行完毕需要关闭 EntryChannel 和 TasksChannel
	close(pool.TasksChannel)
	close(pool.EntryChannel)
}

// 从池子中拿出一个 worker 开始工作
func (p *Pool) worker(work_ID int) {
	// worker 不断的从 JobsChannel 内部任务队列中拿任务
	for task := range p.TasksChannel {
		// 如果拿到了任务则执行，即调用任务所绑定的业务函数
		task.fn(work_ID)
	}
}
