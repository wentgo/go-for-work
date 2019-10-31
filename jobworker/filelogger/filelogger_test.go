package filelogger

import (
	"math/rand"
	"testing"
	"time"
)

func Test_FileLogger1(t *testing.T) {
	logger := New("test.log")
	//logger.SetLevel(filelogger.ERROR)
	logger.Print("output of logger.Print")
	logger.Emerg("output of logger.Emerg")
	logger.Alert("output of logger.Alert")
	logger.Crit("output of logger.Crit")
	logger.Error("output of logger.Error")
	logger.Warning("output of logger.Warning")
	logger.Notice("output of logger.Notice")
	logger.Info("output of logger.Info")
	logger.Debug("output of logger.Debug")
	logger.Flush()
}

func Test_FileLogger2(t *testing.T) {
	rand.Seed(time.Now().Unix())

	logger := New("test.log")
	//logger.SetLevel(filelogger.ERROR)
	logger.Print("output of logger.Print %d", rand.Intn(100))
	logger.Emerg("output of logger.Emerg %d", rand.Intn(100))
	logger.Alert("output of logger.Alert %d", rand.Intn(100))
	logger.Crit("output of logger.Crit %d", rand.Intn(100))
	logger.Error("output of logger.Error %d", rand.Intn(100))
	logger.Warning("output of logger.Warning %d", rand.Intn(100))
	logger.Notice("output of logger.Notice %d", rand.Intn(100))
	logger.Info("output of logger.Info %d", rand.Intn(100))
	logger.Debug("output of logger.Debug %d", rand.Intn(100))
	logger.Flush()
}
