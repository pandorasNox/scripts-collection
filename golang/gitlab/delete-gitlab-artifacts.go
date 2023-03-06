// largely inspired by https://github.com/rafaelperoco/gitlab-artifacts-cleaner/blob/ad3f7c490e18ecbfe1d3998fadad06f89dc19942/main.go

// gitlab token needs:
//  - api access
//  - maintainer role

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Job struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Status            string    `json:"status"`
	FinishedAt        time.Time `json:"finished_at"`
	ArtifactsExpireAt time.Time `json:"artifacts_expire_at"`
}

func main() {
	var dryRun bool
	var project_id int
	var per_page int
	var pages int
	var startPage int
	var token string
	var server string
	var dontDeleteYoungerThan time.Duration
	var dontDeleteOlderThan time.Duration

	flag.BoolVar(&dryRun, "dry-run", true, "Enable dry-run (execution without any deletions)")
	flag.IntVar(&project_id, "project-id", 0, "Project ID")
	flag.IntVar(&per_page, "per-page", 100, "Number of jobs per page")
	flag.IntVar(&pages, "pages", 1, "Number of pages")
	flag.IntVar(&startPage, "start-page", 1, "Number of pages")
	flag.StringVar(&token, "token", "", "Private token")
	flag.StringVar(&server, "server", "", "Gitlab server")
	flag.DurationVar(&dontDeleteYoungerThan, "dont-delete-younger-than", 30*24*time.Hour, "sets that artifacts from jobs which are younger than the provided duration are not deleted, goes with --dont-delete-older-than and must be smaller")
	flag.DurationVar(&dontDeleteOlderThan, "dont-delete-older-than", 60*24*time.Hour, "sets that artifacts from jobs which are older than the provided duration are not deleted, goes with --dont-delete-younger-than and must be bigger")
	flag.Parse()

	fmt.Printf("dontDeleteYoungerThan: %d \n", dontDeleteYoungerThan)
	fmt.Printf("dontDeleteOlderThan: %d \n", dontDeleteOlderThan)

	if dontDeleteOlderThan <= dontDeleteYoungerThan {
		log.Fatal("EXIT: dontDeleteOlderThan must be larger than dontDeleteYoungerThan")
	}

	fmt.Printf("Fetching %d pages from %v\n", pages, server)

Loop:
	for i := startPage; i <= pages; i++ {
		url := fmt.Sprintf("%v/api/v4/projects/%v/jobs?per_page=%d&page=%d&artifact_expired=false", server, project_id, per_page, i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		println("page number:", i)
		req.Header.Set("PRIVATE-TOKEN", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()

		var jobs []Job
		err = json.NewDecoder(resp.Body).Decode(&jobs)
		if err != nil {
			log.Fatalf("error: %s, resp.status='%s'", err, resp.Status)
			return
		}

		for _, job := range jobs {
			// fmt.Println("job to process:", job.ID)

			if job.FinishedAt.IsZero() {
				continue
			}

			// job is younger than dontDeleteYoungerThan but also has expired job artifact
			if job.FinishedAt.After(time.Now().Add(-1*dontDeleteYoungerThan)) && job.ArtifactsExpireAt.Before(time.Now()) {
				_, _ = deleteJobArtifacts(server, token, project_id, job.ID)
				fmt.Println("  attemted deletion (artifacts already expired, ignored if failed bec. nice to have) for job: ", job.ID)
				continue
			}

			// job is younger than dontDeleteYoungerThan
			if job.FinishedAt.After(time.Now().Add(-1 * dontDeleteYoungerThan)) {
				fmt.Println("  skipped (not in deletion periode): ", job.ID, job.Status, job.Name, job.FinishedAt, job.FinishedAt.IsZero())
				continue
			}

			// job is older than dontDeleteOlderThan
			if job.FinishedAt.Before(time.Now().Add(-1 * dontDeleteOlderThan)) {
				fmt.Println("  stopped, found job older than dontDeleteOlderThan: ", job.ID, job.Status, job.Name, job.FinishedAt, job.FinishedAt.IsZero())
				break Loop
			}

			if dryRun {
				fmt.Println("  skipped (dry-run): ", job.ID, job.Status, job.Name, job.FinishedAt, job.FinishedAt.IsZero())
				continue
			}

			respStatus, err := deleteJobArtifacts(server, token, project_id, job.ID)
			if err != nil {
				log.Fatalf("failed deleting job(id: %s) artifacts: %s", job.ID, err)
			}

			fmt.Printf("job artifacts deleted: %d | response Status: %s \n", job.ID, respStatus)
		}
	}
}

func deleteJobArtifacts(server string, token string, projectId int, jobId int) (string, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%d/jobs/%d/artifacts", server, projectId, jobId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error constructing http request: %s", err)
	}

	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing http request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("http error, response status: '%s'", resp.Status)
	}

	return resp.Status, nil
}
