package alarm

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Alarm interface {
	Report(receiver []string, message string) error
}

type DingDing struct {
	Url string
}

func (d DingDing) Report(receiver []string, message string) (err error) {
	baseBody := fmt.Sprintf(`{"msgtype": "markdown", 
    "markdown": {
        "title": "%s",
        "text": "%s"
     },
	"at":{
       "atMobiles": %s
     },
     "isAtAll": false
	}`, strings.Split(message, "\n")[0], message+"\n@"+strings.Join(receiver, " @"), receiver)
	req := bytes.NewBuffer([]byte(baseBody))
	_, err = http.DefaultClient.Post(d.Url, "application/json", req)
	return err
}
