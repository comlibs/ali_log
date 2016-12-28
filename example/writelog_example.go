package main

import (
	"fmt"
	"github.com/comlibs/ali_log"
	"github.com/comlibs/ali_log/pb"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"time"
)

type appConf struct {
	Url             string
	AccessKey       string
	AccessKeySecret string
}

func main() {
	conf := appConf{
		Url:             "vslog.cn-hangzhou.log.aliyuncs.com",
		AccessKey:       "xxxxxxxxxxxxxxxxxxxxxx",
		AccessKeySecret: "xxxxxxxxxxxxxxxxxxxxxx",
	}
	client := ali_log.NewAliLogClient(
		conf.Url,
		conf.AccessKey,
		conf.AccessKeySecret,
	)

	logStore := ali_log.NewLogstore("go-log-service", client)

	lgdata := logMaker()

	logdata, _ := proto.Marshal(lgdata)
	resp, err := logStore.WriteLog(logdata)
	if err != nil {
		fmt.Println("[writelog_example][err]", err)
		return
	}
	fmt.Println("[writelog_example][success]", resp.StatusCode)
	respbytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("[writelog_example][body]", string(respbytes))
	fmt.Println("finished")
}

func logMaker() *pb.LogGroup {
	logMap := map[string]string{
		"IP":             "192.168.1.70",
		"OccurTimeStamp": "19191991",
		"AppName":        "logs",
		"OptName":        "test",
		"InputParams":    "InputParams-test",
		"Err":            "success",
	}
	logData := &pb.Log{
		Time: proto.Uint32(uint32(time.Now().Unix())),
	}
	for k, v := range logMap {
		content := &pb.Log_Content{
			Key:   proto.String(k),
			Value: proto.String(v),
		}
		logData.Contents = append(logData.Contents, content)
	}
	logDataGroup := &pb.LogGroup{}
	logDataGroup.Logs = append(logDataGroup.Logs, logData)
	return logDataGroup
}

func multiGroupLogMaker() *pb.LogGroupList {
	logMap := map[string]string{
		"IP":             "192.168.1.70",
		"OccurTimeStamp": "19191991",
		"AppName":        "logs",
		"OptName":        "test",
		"InputParams":    "InputParams-test",
		"Err":            "success",
	}
	logData := &pb.Log{
		Time: proto.Uint32(uint32(time.Now().Unix())),
	}
	for k, v := range logMap {
		content := &pb.Log_Content{
			Key:   proto.String(k),
			Value: proto.String(v),
		}
		logData.Contents = append(logData.Contents, content)
	}
	logDataGroup := &pb.LogGroup{}
	logDataGroup.Logs = append(logDataGroup.Logs, logData)
	logDataGroupList := &pb.LogGroupList{}
	logDataGroupList.LogGroupList = append(logDataGroupList.LogGroupList, logDataGroup)
	return logDataGroupList
}
