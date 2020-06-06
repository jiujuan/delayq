package delayq

import (
	"errors"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiujuan/delayq/config"
	log "github.com/jiujuan/delayq/logger"
	"github.com/jiujuan/delayq/redis"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// Job job的元信息
type Job struct {
	Topic string `json:"topic"` //Job类型。可以理解成具体的业务名称
	ID    string `json:"id"`    //Job的唯一标识。用来检索和删除指定的Job信息。
	Delay int64  `json:"delay"`
	TTR   int64  `json:"ttr"`
	Body  string `json:"body"` //Job的内容，供消费者做具体的业务处理，以json格式存储
}

// BucketItem 存放于zset中，以时间为单位的有序队列
type BucketItem struct {
	jobID     string `json:"jobid"` // job ID
	timestamp int64  //时间单位
}

// ReadyQueueItem 存放Ready 状态的job
type ReadyQueueItem struct {
	jobID string `json:"jobid"` // job ID
}

// AddJob 增加job到Redis中
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

// PopJob 从Redis中的List获取job
func PopJob(ctx *gin.Context) {

}

//DeleteJob 删除一个job
func DeleteJob(ctx *gin.Context) {

}

func FinishJob(ctx *gin.Context) {

}

// AddJob 增加到job pool里(K/V)
func (job Job) AddJob(jobjson string) error {
	_, err := redis.SET(GetJobPoolKey(job.ID), jobjson)
	return err
}

// GetJob 获取job
func (job Job) GetJob(jobId string) (string, error) {
	val, err := redis.GET(GetJobPoolKey(jobId))
	return val, err
}

// DelJob 删除job
func (job Job) DelJob() error {
	err := redis.DEL(job.ID)
	return err
}

// AddDelayQueue 增加到延迟队列Delay Queue中 (zset)
func (item BucketItem) AddDelayQueue() error {
	_, err := redis.ZADD(GetDelayQueueKey(), item.timestamp, item.jobID)
	return err
}

// GetExpireJob 获取延迟队列中到期的job任务，实际获取的是job任务的ID
func GetExpireJob(t time.Time) ([]string, error) {
	res, err := redis.ZRANGEBYSCORE(GetDelayQueueKey(), "0", t.Unix())
	if err != nil {
		return nil, err
	}

	var jobIds []string
	if len(res) > 0 {
		for _, jobId := range res {
			jobIds = append(jobIds, jobId)
		}
		return jobIds, nil
	}
	return jobIds, errors.New("job id empty")
}

// DelDelayQueue 删除延迟队列中的任务
func (item BucketItem) DelDelayQueue(jobId string) error {
	err := redis.ZREM(GetDelayQueueKey(), jobId)
	return err
}

// AddReadyQueue 增加到Ready Queue里
func (item ReadyQueueItem) AddReadyQueue(jobId string) error {
	err := redis.LPUSH(GetReadyQueueKey(), jobId)
	return err
}

// GetReadyQueue 从Ready Queue里取出数据
func GetReadyQueue() (string, error) {
	val, err := redis.BRPOP(GetReadyQueueKey())
	return val, err
}

// GetJobPoolKey ID -> jobid
func GetJobPoolKey(ID string) string {
	return config.QConfig.DelayQ.JobPoll + ID
}

// GetReadyQueueKey 获取 ReadyQueue 的key
func GetReadyQueueKey() string {
	return config.QConfig.DelayQ.ReadyQueue
}

// GetDelayQueueKey 获取延迟队列的 key
func GetDelayQueueKey() string {
	return config.QConfig.DelayQ.DelayQueue
}
