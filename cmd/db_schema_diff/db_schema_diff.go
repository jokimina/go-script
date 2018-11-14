package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	cfg *Config
	saveDir   = time.Now().Format("2006-01-02")
	fDiffText = saveDir + string(os.PathSeparator) + "diff.txt"
	fDiffSql  = saveDir + string(os.PathSeparator) + "diff.sql"
	c = make(chan int, 10)
	q = make(chan bool)
)

func init(){
	os.Mkdir(saveDir, 0600)
	cfgFile := filepath.Join(filepath.Dir(os.Args[0]),"config.json")
	err := loadJSONFile(cfgFile, &cfg)
	checkErr(err)
}

func getCompareTables(ds string) (tables []string) {
	db, err := sql.Open("mysql", ds)
	defer db.Close()
	checkErr(err)

	rows, err := db.Query("SHOW TABLES")
	checkErr(err)
	var t string
	for rows.Next() {
		rows.Scan(&t)
		tables = append(tables, t)
	}
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func urlGoToJdbc(url string) string {
	r := strings.NewReplacer("@tcp(", "@", ")", "")
	url = r.Replace(url)
	return url
}

// Need install mysql-utilities
func mysqlDiff(sourceDsn string, destDsn string) {
	sDsList := strings.Split(sourceDsn, "/")
	dDsList := strings.Split(destDsn, "/")
	sDs, sDbName := sDsList[0], sDsList[1]
	dDs, dDbName := dDsList[0], dDsList[1]
	diffStdout, err := os.OpenFile(fDiffText, os.O_CREATE|os.O_WRONLY, 0644)
	sqlStdout, err := os.OpenFile(fDiffSql, os.O_CREATE|os.O_WRONLY, 0644)
	strBuf := bytes.NewBufferString("")
	checkErr(err)
	defer diffStdout.Close()
	defer sqlStdout.Close()

	go func() {
		for {
			select {
			case <-c:
			case <-time.After(15 * time.Second):
				q <- true
				break
			}
		}
	}()

	for _, table := range getCompareTables(sourceDsn) {
		command := []string{
			fmt.Sprintf("--server1=%s", urlGoToJdbc(sDs)),
			fmt.Sprintf("--server2=%s", urlGoToJdbc(dDs)),
			fmt.Sprintf("%[1]s.%[3]s:%[2]s.%[3]s", sDbName, dDbName, table),
			"--skip-table-options",
		}
		// diff format
		go func() {
			cmd := exec.Command("mysqldiff", command...)
			multiWriter := io.MultiWriter(os.Stdout, diffStdout, strBuf)
			cmd.Stdout = multiWriter
			cmd.Stderr = multiWriter
			//fmt.Println(cmd.Args)
			cmd.Run()
			c <- 0
			fmt.Println("Done diff -- ", table)
		}()

		// sql format
		go func() {
			cmd := exec.Command("mysqldiff", command...)
			multiWriter := io.MultiWriter(os.Stdout, sqlStdout)
			cmd.Args = append(cmd.Args, "-d", "sql")
			cmd.Stdout = multiWriter
			cmd.Stderr = multiWriter
			//fmt.Println(cmd.Args)
			cmd.Run()
			c <- 0
			fmt.Println("Done sql -- ", table)
		}()
	}
	<-q
	if cfg.Email.SendMail {
		cfg.Email.Subject = fmt.Sprintf("【Mysqldiff】%s %s", strings.Split(urlGoToJdbc(sourceDsn), "@")[1], strings.Split(urlGoToJdbc(destDsn), "@")[1])
		body, err := ioutil.ReadAll(strBuf)
		checkErr(err)
		sendMail(string(body))
	}
}

func sendMail(body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Email.From)
	m.SetHeader("Subject", cfg.Email.Subject)
	m.SetHeader("To", cfg.Email.To[0])
	m.SetBody("text/plain", fmt.Sprintf("变更SQL文件见附件！\n %s", body))
	if len(cfg.Email.To) > 1 {
		m.SetHeader("Cc", cfg.Email.To[1:]...)
	}
	m.Attach(fDiffSql)
	d := gomail.NewDialer(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.From, cfg.Email.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	checkErr(err)
}

func main() {
	for _, l := range cfg.Comparelist {
		mysqlDiff(l.SourceDsn, l.DestDsn)
	}
}
