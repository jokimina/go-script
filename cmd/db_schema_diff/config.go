package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Config struct {
	Comparelist []struct {
		SourceDsn string `json:"sourceDsn"`
		DestDsn   string `json:"destDsn"`
	} `json:"comparelist"`
	Email struct {
		SendMail bool     `json:"send_mail"`
		SMTPHost string   `json:"smtp_host"`
		SMTPPort int      `json:"smtp_port"`
		Subject  string   `json:"subject"`
		From     string   `json:"from"`
		Password string   `json:"password"`
		To       []string `json:"to"`
	} `json:"email"`
}

func (cfg *Config) String() string {
	ds, _ := json.MarshalIndent(cfg, "  ", "  ")
	return string(ds)
}

func loadJSONFile(jsonPath string, val interface{}) error {
	bs, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bs), "\n")
	var bf bytes.Buffer
	for _, line := range lines {
		lineNew := strings.TrimSpace(line)
		if (len(lineNew) > 0 && lineNew[0] == '#') || (len(lineNew) > 1 && lineNew[0:2] == "//") {
			continue
		}
		bf.WriteString(lineNew)
	}
	return json.Unmarshal(bf.Bytes(), &val)
}