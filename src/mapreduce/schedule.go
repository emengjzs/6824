package mapreduce

import (
	"fmt"
	"sync"
)

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var nOther int // number of inputs (for reduce) or outputs (for map)

	var getTaskArgsFunc func(int) DoTaskArgs
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		nOther = nReduce
		getTaskArgsFunc = func(taskNum int) DoTaskArgs {
			return DoTaskArgs{jobName, mapFiles[taskNum], phase, taskNum, nOther}
		}
	case reducePhase:
		ntasks = nReduce
		nOther = len(mapFiles)
		getTaskArgsFunc = func(taskNum int) DoTaskArgs {
			return DoTaskArgs{jobName, "", phase, taskNum, nOther}
		}
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nOther)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//

	var finish sync.WaitGroup
	finish.Add(ntasks)

	taskSet := make(chan int)

	assignTask := func(instance string, task int) {
		taskArgs := getTaskArgsFunc(task)
		success := call(instance, "Worker.DoTask", taskArgs, nil)
		if !success {
			fmt.Printf("Fail Task on instance(%s) phase(%v) task(%d/%d)", instance, phase, task, ntasks)
			taskSet <- task
		} else {
			finish.Done()
			registerChan <- instance
		}
	}

	runTasks := func() {
		for task := range taskSet {
			instance := <-registerChan
			go assignTask(instance, task)
		}
	}

	go runTasks()

	for i := 0; i < ntasks; i++ {
		taskSet <- i
	}

	// wait until finish all tasks
	finish.Wait()
	close(taskSet)
	fmt.Printf("Schedule: %v phase done\n", phase)
}
