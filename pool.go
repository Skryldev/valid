package validx

import "runtime"

func defaultWorkers() int {
	return runtime.NumCPU()
}
