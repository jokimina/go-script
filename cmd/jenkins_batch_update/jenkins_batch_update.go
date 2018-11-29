package main

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"os"
	"regexp"
	"strings"
)

func main() {
	jenkins, err := gojenkins.CreateJenkins(nil, os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASS")).Init()
	if err != nil {
		panic(err)
	}
	//job, err := jenkins.GetJob("test/job/mars/job/sc-aql")
	job, err := jenkins.GetJob("platform/job/mars")
	if err != nil {
		panic(err)
	}

	managedShellId := "bb48c392-2bf9-49e0-bae9-83b3ab1f5ad0"
	jobs, err := job.GetInnerJobs()
	for _, j := range jobs {
		jobName := j.GetName()
		if strings.HasPrefix(jobName, "sc-") && !strings.Contains(jobName, "eureka") && !strings.Contains(jobName, "web") {
			tJob, err := jenkins.GetJob("platform/job/mars/job/" + jobName)
			if err != nil {
				panic(err)
			}
			c, err := tJob.GetConfig()
			if strings.Contains(c, managedShellId) {
				continue
			}
			lines := strings.Split(c, "\n")
			re := regexp.MustCompile("(>)(.*)(<)")
			for i, line := range lines {
				if strings.Contains(line, "<buildStepId>") {
					fmt.Println("try update job " + jobName + ": uat/job/mars/job/" + jobName)
					lines[i] = re.ReplaceAllString(line, fmt.Sprintf("${1}%s${3}", managedShellId))
					break
				}
			}
			// 删除中文
			newConfig := regexp.MustCompile("[\u4e00-\u9fa5]+").ReplaceAllString(strings.Join(lines, "\n"), "")
			//fmt.Println(newConfig)
			err = tJob.UpdateConfig(newConfig)
			if err != nil {
				panic(err)
			}
			//break
		}
	}

}
