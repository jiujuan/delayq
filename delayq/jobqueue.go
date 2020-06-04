package delayq

import (
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiujuan/delayq/config"
	log "github.com/jiujuan/delayq/logger"
	"github.com/jiujuan/delayq/redis"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Job struct {
	Topic string `json:"topic"` //Job类型。可以理解成具体的业务名称
	ID    string `json:"id"`    //Job的唯一标识。用来检索和删除指定的Job信息。
	Delay int64  `json:"delay"`
	TTR   int64  `json:"ttr"`
	Body  string `json:"body"` //Job的内容，供消费者做具体的业务处理，以json格式存储
}

type BucketItem struct {
	jobID     string `json:"jobid"` // job ID
	timestamp int64
}

type ReadyQueueItem struct {
	jobID string `json:"jobid"` // job ID
}

func AddJob(ctx *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := ctx.Request.Body.Read(buf) // 获取body的内容
	jobJson := string(buf[0:n])
	fmt.Println(jobJson)

	var job Job
	if err := jsoniter.Unmarshal(buf[0:n], &job); err != nil {
		log.Logger.Error("job json error, ", zap.String("json", jobJson))
		ctx.JSON(200, gin.H{
			"msg":    "job error",
			"status": 601,
		})
		return
	}

	jobData := Job{
		ID:    job.ID,
		Topic: job.Topic,
		Delay: job.Delay,
		Body:  job.Body,
		TTR:   job.TTR,
	}
	// 增加到job pool
	err := jobData.AddJob(jobJson)
	if err != nil {
		log.Logger.Info("add job to redis failed",
			zap.String("job json", jobJson),
		)
	}

	// 增加到Delay Queue
	bucketItem := BucketItem{
		jobID:     job.ID,
		timestamp: time.Now().Unix() + job.Delay,
	}
	err = bucketItem.AddDelayQueue()
	if err != nil {
		log.Logger.Info("add bucketitem to redis failed",
			zap.String("bucketitem id", bucketItem.jobID),
			zap.Int64("timestamp", bucketItem.timestamp),
		)
	}

	ctx.JSON(200, gin.H{
		"msg":    "success" + job.Body,
		"status": 201,
	})
}

func PopJob(ctx *gin.Context) {

}

func DeleteJob(ctx *gin.Context) {

}

func FinishJob(ctx *gin.Context) {

}

// func GetExpireJob() ([]string, error) {

// }

// 增加到job pool里(K/V)
func (job Job) AddJob(jobjson string) error {
	_, err := redis.SET(GetJobPoolKey(job.ID), jobjson)
	return err
}

// 获取job
func (job Job) GetJob(jobid string) (string, error) {
	val, err := redis.GET(GetJobPoolKey(job.ID))
	return val, err
}

// 删除job
func (job Job) DelJob() error {
	err := redis.DEL(job.ID)
	return err
}

// 增加到延迟队列Delay Queue中 (zset)
func (item BucketItem) AddDelayQueue() error {
	_, err := redis.ZADD(GetDelayQueueKey(), item.timestamp, item.jobID)
	return err
}

// 获取延迟队列中到期的job任务
func GetExpireJob() ([]string, error) {
	return redis.ZRANGEBYSCORE(GetDelayQueueKey(), "0", time.Now().Unix())
}

// 删除延迟队列中的任务
func (item BucketItem) DelDelayQueue() error {
	err := redis.ZREM(GetDelayQueueKey(), item.jobID)
	return err
}

// 增加到Ready Queue里
func (item ReadyQueueItem) AddReadyQueue() error {
	err := redis.LPUSH(GetReadyQueueKey(), item.jobID)
	return err
}

// 从Ready Queue里取出数据
func GetReadyQueue() (string, error) {
	val, err := redis.BRPOP(GetReadyQueueKey())
	return val, err
}

// ID -> jobid
func GetJobPoolKey(ID string) string {
	return config.QConfig.DelayQ.JobPoll + ID
}

func GetReadyQueueKey() string {
	return config.QConfig.DelayQ.ReadyQueue
}

func GetDelayQueueKey() string {
	return config.QConfig.DelayQ.DelayQueue
}
