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
	job, err := jenkins.GetJob("test/job/mars")
	if err != nil {
		panic(err)
	}

	managedShellId := "7e260101-25d7-46b8-8a57-9c25fe9a6e4a"
	jobs, err := job.GetInnerJobs()
	for _, j := range jobs {
		jobName := j.GetName()
		if strings.HasPrefix(jobName, "sc-") && !strings.Contains(jobName, "eureka") {
			fmt.Println("try update job " + jobName)
			tJob, err := jenkins.GetJob("test/job/mars/job/" + jobName)
			if err != nil {
				panic(err)
			}
			c, err := tJob.GetConfig()
			lines := strings.Split(c, "\n")
			re := regexp.MustCompile("(>)(.*)(<)")
			for i, line := range lines {
				if strings.Contains(line, "<buildStepId>") && !strings.Contains(line, managedShellId) {
					lines[i] = re.ReplaceAllString(line, fmt.Sprintf("${1}%s${3}", managedShellId))
				}
			}
			newConfig := strings.Join(lines, "\n")
			//fmt.Println(newConfig)
			err = tJob.UpdateConfig(newConfig)
			if err != nil {
				panic(err)
			}
			break
		}
	}

}
