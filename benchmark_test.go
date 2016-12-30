package ali_log

import (
	"fmt"
	"testing"
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

func BenchmarkAliLogRecordGorountine(b *testing.B) {
	conf := appConf{
		Url:             "vslog.cn-hangzhou.log.aliyuncs.com",
		AccessKey:       "xxxxxxxxxxxxxxxxxxxxxx",
		AccessKeySecret: "xxxxxxxxxxxxxxxxxxxxxx",
		LogStore:        "go-log-service",
		Topic:           "async-service",
		C:               100,
	}

	conf.AliLogMapChan = make(chan map[string]string, conf.C)
	AliLogRecordGoroutine(
		conf.Topic,
		NewLogstore(
			conf.LogStore,
			NewAliLogClient(
				conf.Url,
				conf.AccessKey,
				conf.AccessKeySecret,
			)),
		conf.AliLogMapChan)

	for i := 0; i < b.N; i++ {
		conf.AliLogMapChan <- map[string]string{
			"IP":             "192.168.1.70",
			"OccurTimeStamp": fmt.Sprintf("%v", time.Now().Unix()),
			"AppName":        "logs",
			"OptName":        "test",
			"InputParams":    "InputParams-test",
			"Err":            "success",
		}
	}

}
