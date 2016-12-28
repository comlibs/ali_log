package main

import (
	"fmt"
	"github.com/comlibs/ali_log"
	"time"
)

type appConf struct {
	Url             string
	AccessKey       string
	AccessKeySecret string
	LogStore        string
	Topic           string
	C               int
	AliLogMapChan   chan map[string]string
}

func main() {
	conf := appConf{
		Url:             "vslog.cn-hangzhou.log.aliyuncs.com",
		AccessKey:       "xxxxxxxxxxxxxxxxxxxxxx",
		AccessKeySecret: "xxxxxxxxxxxxxxxxxxxxxx",
		LogStore:        "go-log-service",
		Topic:           "async-service",
		C:               10,
	}

	logStore := ali_log.NewLogstore(
		conf.LogStore,
		ali_log.NewAliLogClient(
			conf.Url,
			conf.AccessKey,
			conf.AccessKeySecret,
		))

	conf.AliLogMapChan = make(chan map[string]string, conf.C)
	ali_log.AliLogRecordGoroutine(
		conf.Topic,
		logStore,
		conf.AliLogMapChan)

	for i := 0; i < 100; i++ {
		conf.AliLogMapChan <- map[string]string{
			"IP":             "192.168.1.70",
			"OccurTimeStamp": "19191991",
			"AppName":        "logs",
			"OptName":        "test",
			"InputParams":    "InputParams-test",
			"Err":            "success",
			"CreateTime":     fmt.Sprintf("%v", time.Now().Unix()),
		}
	}

}
