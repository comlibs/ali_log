package ali_log

import (
	"github.com/comlibs/ali_log/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

const (
	LOG_NUM  = 3000              // 当日志条数达到 3000条时 触发 写入ali log 请求
	LOG_SIZE = 2.5 * 1024 * 1024 //当日志大小达到2.5M时 触发 写入 ali log 请求
)

func AliLogRecordGoroutine(topic string, logStore AliLogstore, respChan chan map[string]string) {
	go func(topic string, logStore AliLogstore, respChan chan map[string]string) {

		//ali_log 官方文档:日志数据时间戳在服务端当前处理时间前后[-7x24小时, +15分钟]小时范围内的日志,不在该时间范围内，则整个请求失败，且无任何日志数据成功写入。
		ticker := time.NewTicker(time.Duration(6*24*60*60) * time.Second)
		//ticker := time.NewTicker(1 * 24 * 60 * 60 * time.Second) //每天时触发一次 ali_log 写入请求
		defer ticker.Stop()

		aliLogGroup := &pb.LogGroup{
			Topic: proto.String(topic),
		}

		for {
			select {
			case resp := <-respChan:
				putLogs(aliLogGroup, logStore, resp)
			case <-ticker.C:
				if len(aliLogGroup.Logs) > 0 {
					aliLogBytes, _ := proto.Marshal(aliLogGroup)
					writeAliLog(logStore, aliLogBytes, aliLogGroup)
				}
			}
		}
	}(topic, logStore, respChan)
}

func putLogs(aliLogGroup *pb.LogGroup, logStore AliLogstore, logMsg map[string]string) {
	newLog := &pb.Log{
		Time: proto.Uint32(uint32(time.Now().Unix())),
	}
	for k, v := range logMsg {
		content := &pb.Log_Content{
			Key:   proto.String(k),
			Value: proto.String(v),
		}
		newLog.Contents = append(newLog.Contents, content)
	}

	aliLogBytes, _ := proto.Marshal(aliLogGroup)

	//ali_log 官方文档: 日志一次写入条数超过4096条 或大小超过3M, 超过则写入失败
	if len(aliLogGroup.Logs)+1 >= LOG_NUM || len(aliLogBytes) > LOG_SIZE {
		writeAliLog(logStore, aliLogBytes, aliLogGroup)
	}
	aliLogGroup.Logs = append(aliLogGroup.Logs, newLog)
}

func writeAliLog(logStore AliLogstore, aliLogBytes []byte, aliLogGroup *pb.LogGroup) {
	if _, err := logStore.WriteLog(aliLogBytes); err != nil { //写入日志请求失败
		go func(logStore AliLogstore, aliLogBytes []byte) {
			time.Sleep(time.Duration(1) * time.Second)
			if _, err := logStore.WriteLog(aliLogBytes); err != nil { //写入日志请求失败
				time.Sleep(time.Duration(2) * time.Second)
				logStore.WriteLog(aliLogBytes)
			}
		}(logStore, aliLogBytes)

	}
	//不管是否写入成功，都将清空
	aliLogGroup.Logs = nil
}
