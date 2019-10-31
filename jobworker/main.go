package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"jobworker/filelogger"

	"github.com/go-redis/redis/v7"
)

const WEBROOT = "c:/xampp/htdocs/btenew"

var logger *filelogger.FileLogger

func dpr(v interface{}) {
	fmt.Printf("%#v\n", v)
}

type Job struct {
	Time int64
	Name string
	Args []string
}

func php(args ...string) error {
	logger.Print(strings.Join(args, " "))

	//cmd := exec.Command("c:/xampp/php/php", args...)
	cmd := exec.Command("c:/xampp/php64/php", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = WEBROOT + "/job"

	err := cmd.Run()
	if err != nil {
		logger.Print(err.Error())
	}

	return err
}

func main() {
	// Setup Logger
	logfile := WEBROOT + "/app/logs/jobworker.log"
	logger = filelogger.New(logfile)

	jobQueue := NewJobQueue("job:queue")

	for {
		job, err := jobQueue.Pop(5)
		if err == nil {
			go php(job...)
		}
		jobQueue.CheckDelayJob()
		logger.Flush()
	}
}

type JobQueue struct {
	name  string
	redis *redis.Client
}

func NewJobQueue(name string) *JobQueue {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return &JobQueue{
		redis: client,
		name:  name,
	}
}

func (q JobQueue) Push(job Job) error {
	msg, err := json.Marshal(job)
	if err != nil {
		logger.Error("JobQueue.Push @1 " + err.Error())
		return err
	}

	if err = q.redis.RPush(q.name, string(msg)).Err(); err != nil {
		logger.Error("JobQueue.Push @2 " + err.Error())
		return err
	}

	return nil
}

func (q JobQueue) Pop(d time.Duration) ([]string, error) {
	result, err := q.redis.BLPop(d*time.Second, q.name).Result()
	if err != nil {
		// Queue is empty
		return nil, err
	}
	//dpr(result)

	var job Job

	err = json.Unmarshal([]byte(result[1]), &job)
	if err != nil {
		logger.Error("JobQueue.Pop " + err.Error())
		logger.Print(result[1])
		return nil, err
	}
	//dpr(job)

	var arr []string

	arr = append(arr, job.Name)
	arr = append(arr, job.Args...)

	return arr, nil
}

func (q JobQueue) CheckDelayJob() error {
	const quename = "delay:job"

	result, err := q.redis.ZRange(quename, 0, 0).Result()
	if err != nil {
		logger.Error("CheckDelayJob @1 " + err.Error())
		q.redis.ZRemRangeByRank(quename, 0, 0)
		return err
	}

	if len(result) == 0 {
		// Queue is empty
		return nil
	}

	var job Job

	err = json.Unmarshal([]byte(result[0]), &job)
	if err != nil {
		logger.Error("CheckDelayJob @2 " + err.Error())
		logger.Print(result[0])
		q.redis.ZRemRangeByRank(quename, 0, 0)
		return err
	}
	//dpr(job)

	now := time.Now().Unix()

	if job.Time <= now {
		q.Push(job)
		q.redis.ZRemRangeByRank(quename, 0, 0)
	}

	return nil
}
