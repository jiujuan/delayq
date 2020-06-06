package delayq

import (
	"time"
)

func Start() {
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				run(t)
			}
		}
	}()
}

func run(t time.Time) {
	jobs, err := GetExpireJob(t)
	if err != nil {
		return
	}
	if len(jobs) > 0 {
		var job = Job{}
		var bucketItem = BucketItem{}
		var readyQueueItem = ReadyQueueItem{}
		for _, jobid := range jobs {
			// 获取job的原始信息，如果不能存在就删掉 DelayQueue 中的信息
			jobJson, err := job.GetJob(jobid)
			if err != nil {
				err := bucketItem.DelDelayQueue(jobid)
				if err != nil {
					return
				}
			} else {
				readyQueueItem.AddReadyQueue(jobid)
			}
		}
	}
}
