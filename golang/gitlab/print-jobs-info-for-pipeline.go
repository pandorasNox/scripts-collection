package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	// "time"
)

type Job struct {
	// ID                int        `json:"id"`
	Name string `json:"name"`
	// Status            string     `json:"status"`
	// FinishedAt        time.Time  `json:"finished_at"`
	// ArtifactsExpireAt time.Time  `json:"artifacts_expire_at"`
	Artifacts []artifact `json:"artifacts"`
}

// {"file_type": "archive", "size": 1000, "filename": "artifacts.zip", "file_format": "zip"},
type artifact struct {
	FileName  string `json:"filename"`
	SizeBytes int    `json:"size"`
}

func main() {
	fmt.Println("...start...")
	log.Println("first log")

	var dryRun bool
	var server string
	var projectId int
	var pipelineId int
	var token string

	flag.BoolVar(&dryRun, "dry-run", true, "Enable dry-run (execution without any deletions)")
	flag.StringVar(&server, "server", "", "Gitlab server")
	flag.StringVar(&token, "token", "", "Private token")
	flag.IntVar(&projectId, "projectId", 0, "Project ID")
	flag.IntVar(&pipelineId, "pipelineId", 0, "Pipeline ID")
	flag.Parse()

	log.Printf("dryRun: %v", dryRun)
	log.Printf("server: %v", server)
	// log.Printf("using token: %s", token)
	log.Printf("using projectId: %d", projectId)
	log.Printf("using pipelineId: %d", pipelineId)

	//get pipeline
	//
	//GET /projects/:id/pipelines/:pipeline_id

	jobs, err := getPipelineJobs(server, token, projectId, pipelineId)
	if err != nil {
		log.Fatalf("failed getting jobs: %s", err)
	}

	log.Println("jobs:")
	for _, job := range jobs {
		fmt.Println("  Job:", job.Name)
		// fmt.Printf("    artifacts: %v\n", job.Artifacts)
		fmt.Println("    Artifacts:")
		for _, artifact := range job.Artifacts {
			fmt.Println("      - FileName:", artifact.FileName)
			fmt.Printf("        Size: %d Mb (%d Bytes)\n", artifact.SizeBytes/1000/1000, artifact.SizeBytes )
		}

		fmt.Println("")
	}

	fmt.Println("...end...")
}

func getPipelineJobs(server string, token string, projectId int, pipelineId int) ([]Job, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%d/pipelines/%d/jobs", server, projectId, pipelineId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Job{}, fmt.Errorf("error constructing http request: %s", err)
	}

	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Job{}, fmt.Errorf("error executing http request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return []Job{}, fmt.Errorf("http error, response status: '%s'", resp.Status)
	}

	var jobs []Job
	err = json.NewDecoder(resp.Body).Decode(&jobs)
	if err != nil {
		return []Job{}, fmt.Errorf("error failed to json decode response body: %s", err)
	}

	return jobs, nil
}
