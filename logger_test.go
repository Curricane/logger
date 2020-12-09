package logger

import (
	"testing"
)

func TestFileLogger(t *testing.T) {
	config := map[string]string{
		"log_path":       "./",
		"log_name":       "file_log",
		"log_level":      "info",
		"log_split_type": "hour",      // optional default hour
		"log_split_size": "104857600", // optional default 104857600 即100M
		"log_chan_size":  "50000",     // optional default 50000
	}
	log, err := NewFileLoger(config)
	if err != nil {
		t.Errorf("failed to NewFileLoger, err is:%#v\n", err)
	}
	log.Debug("my id is[%d]", 567)
	log.Warn("test warn log")
	log.Fatal("fatal log")
	for i := 0; i < 100; i++ {
		log.Info("wo cao")
	}
	// time.Sleep(time.Second * 2) //需要等待msg从chain中写道file中
	// log.Close()
}

func TestConsoleLogger(t *testing.T) {
	config := map[string]string{"log_level": "info"}
	log, err := NewConsoleLogger(config)
	if err != nil {
		t.Errorf("failed to NewConsoleLogger, err is:%#v\n", err)
	}
	log.Debug("my id is[%d]", 5678)
	log.Warn("test warn log")
	log.Fatal("fatal log")
	log.Close()
}
