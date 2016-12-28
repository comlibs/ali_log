package ali_log

import (
	"errors"
	"fmt"
	"net/http"
)

type AliLogstore interface {
	WriteLog(logdata interface{}) (resp *http.Response, err error)
}

type Logstore struct {
	name   string
	client LogClient
}

func NewLogstore(name string, client LogClient) AliLogstore {
	if 0 == len(name) {
		panic("ali-log: logstore name could not be empty")
	}

	logstore := new(Logstore)
	logstore.name = name
	logstore.client = client
	return logstore
}

func (l *Logstore) WriteLog(logdata interface{}) (resp *http.Response, err error) {
	if resp, err = l.client.Send("POST", nil, logdata, fmt.Sprintf("logstores/%s", l.name)); err != nil {
		err = errors.New("[ali_log][WriteLog][Send] " + err.Error())
		return
	}
	return
}
