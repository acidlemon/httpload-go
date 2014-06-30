package httpload

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	Parallel  int
	Seconds   int
	Urls      []string
	KeepAlive bool
}

type task struct {
	url string
}

type result struct {
	statusCode int
	transfered int64
}

func Start(config Config) {
	// prepare channels
	queue := make(chan *task)
	res := make(chan *result)

	// spawn workers
	var wg sync.WaitGroup
	for i := 0; i < config.Parallel; i++ {
		go func(id int, queue chan *task, res chan *result) {
			tr := &http.Transport{
				DisableKeepAlives: !config.KeepAlive,
			}
			client := &http.Client{Transport: tr}
			var t *task
			for {
				// get task
				t = <-queue
				if t == nil {
					break
				}
				doHttpLoad(id, client, t, res)
			}
			wg.Done()
		}(i, queue, res)
		wg.Add(1)
	}

	// result count
	count := 0

	// prepare tasks
	var tasks []task = make([]task, 0)
	for _, url := range config.Urls {

		t := new(task)
		t.url = url
		tasks = append(tasks, *t)
	}

	index := 0 // tasksのindex
	// put in initial jobs
	for i := 0; i < config.Parallel; i++ {
		queue <- &tasks[index]
		index++
		if index >= len(tasks) {
			index = 0
		} // mod計算するよりは速いはずだ…
		count++
	}

	// start timer
	timer := time.NewTimer(time.Duration(config.Seconds) * time.Second)

	// fill queue like wankosoba
	var re *result
	rescount := make(map[int]int)
	var bytecounts int64 = 0
TIMERLOOP:
	for {
		select {
		case re = <-res:
			// go go wankosoba
			queue <- &tasks[index]
			index++
			if index >= len(tasks) {
				index = 0
			}
			rescount[re.statusCode] += 1
			bytecounts += re.transfered
			count++

		case <-timer.C:
			// time up!
			fmt.Println("timeup!")
			break TIMERLOOP
		}
	}

	// get rest of res
	for i := 0; i < config.Parallel; i++ {
		re = <-res
		rescount[re.statusCode] += 1
		bytecounts += re.transfered
	}

	close(queue) // kill workers

	// wait all goroutine
	wg.Wait()

	// result!
	fmt.Printf("result: score= %.2f req/sec, statusCode=%v, total req count=%d\n",
		float64(count)/float64(config.Seconds), rescount, count)
	fmt.Printf("result: transfered: %d bytes, mean: %.2f bytes/sec, %.2f bytes/req \n", bytecounts,
		float64(bytecounts)/float64(config.Seconds),
		float64(bytecounts)/float64(count))
}

func doHttpLoad(id int, client *http.Client, t *task, res chan *result) {
	//	fmt.Println("worker: ", id, ", request=", t.url)

	resp, err := client.Get(t.url)
	if err != nil {
		re := new(result)
		fmt.Println("error: ", err)
		re.statusCode = 500
		res <- re
	} else {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)
		re := new(result)
		if err != nil {
			fmt.Println("read all error: ", err)
			re.statusCode = 500
			res <- re
		} else {
			re.statusCode = resp.StatusCode
			re.transfered = resp.ContentLength
			res <- re
		}
	}
}
