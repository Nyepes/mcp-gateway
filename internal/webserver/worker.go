package ConnectionHandler



type Task struct {
	TaskID		int
	Connection	net.Conn
	request 	map
}

type WorkerPool struct {
	NumWorkers	int
	TaskQueue	chan Task
}

// 
type Request interface {
	ForwardRequest(Task)
}


func (pool *WorkerPool) Create(proxy Proxy) {
	for i := 0; i < pool.NumWorkers; i++ {
		go startWorker(proxy, pool.TaskQueue)
	}
}

func startWorker(req Request, taskChannel chan Task) {
	select {
	case tasl := <- taskChannel:
		req.ForwardRequest(task)

	}
}





