package proxy_server

import "sync"

type IRunable interface {
	Run() error
}

func start(handle IRunable) {
	handle.Run()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go start(NewHTTP())
	go start(NewSocket5())
	wg.Wait()

}
