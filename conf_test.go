package sqsd

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestInitConf(t *testing.T) {
	c := &Conf{}
	if c.Worker.MaxProcessCount != 0 {
		t.Error("MaxProcessCount is invalid")
	}
	if c.Worker.IntervalSeconds != 0 {
		t.Error("IntervalSeconds is invalid")
	}
	c.Init()
	if c.Worker.MaxProcessCount != 1 {
		t.Error("MaxProcessCount is yet 0")
	}
	if c.Worker.IntervalSeconds != 1 {
		t.Error("IntervalSeconds not set!")
	}
}

func TestValidateConf(t *testing.T) {
	c := &Conf{}
	c.Init()
	c.SQS.QueueURL = "https://example.com/queue/hoge"
	c.SQS.Region = "ap-northeast-1"
	c.Worker.JobURL = "http://localhost:1080/run_job"
	c.Stat.ServerPort = 10000

	if err := c.Validate(); err != nil {
		t.Error("valid conf but error found", err)
	}

	c.SQS.Region = ""
	if err := c.Validate(); err == nil {
		t.Error("sqs.region is required but valid config")
	}

	c.Stat.ServerPort = 0
	if err := c.Validate(); err == nil {
		t.Error("stat.server_port is 0, but no error")
	}

	c.Stat.ServerPort = 10000
	c.SQS.QueueURL = ""
	if err := c.Validate(); err == nil {
		t.Error("SQS.QueueURL is empty, but no error")
	}
	c.SQS.QueueURL = "foo://bar/baz"
	if err := c.Validate(); err == nil {
		t.Error("SQS.QueueURL is not HTTP url, but no error")
	}
	c.SQS.QueueURL = "https://example.com/queue/hoge"
	c.SQS.Region = ""
	if err := c.Validate(); err == nil {
		t.Error("SQS.Region not exists, but no error")
	}
	c.SQS.Region = "ap-northeast-1"

	c.Worker.JobURL = ""
	if err := c.Validate(); err == nil {
		t.Error("Worker.JobURL is empty, but no error")
	}
	c.Worker.JobURL = "foo://bar/baz"
	if err := c.Validate(); err == nil {
		t.Error("Worker.JobURL is not HTTP url, but no error")
	}
}

func TestNewConf(t *testing.T) {
	d, _ := os.Getwd()
	if _, err := NewConf(filepath.Join(d, "test", "conf", "hoge.toml")); err == nil {
		t.Error("file not found")
	}
	if _, err := NewConf(filepath.Join(d, "test", "conf", "config1.toml")); err == nil {
		t.Error("invalid config but passed")
	}
	conf, err := NewConf(filepath.Join(d, "test", "conf", "config_valid.toml"))
	if err != nil {
		t.Error("invalid config??? ", err)
	}
	if conf.SQS.QueueURL != "http://example.com/queue/hoge" {
		t.Error("QueueURL not loaded correctly. " + conf.SQS.QueueURL)
	}
	if conf.Stat.ServerPort != 4080 {
		t.Error("Stat.ServerPort not loaded correctly. " + strconv.Itoa(conf.Stat.ServerPort))
	}
}
