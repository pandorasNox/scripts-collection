// largely inspired by https://github.com/rafaelperoco/gitlab-artifacts-cleaner/blob/ad3f7c490e18ecbfe1d3998fadad06f89dc19942/main.go

// gitlab token needs:
//  - api access
//  - maintainer role

// usage e.g. in gitlab
//
// delete-artifacts:
//   image: golang:1.20.1-alpine3.17
//   stage: some-stage
//   interruptible: true
//   timeout: 15m
//   script:
//     - set -e
//     - test ${GITLAB_PROJECT_ACCESS_TOKEN} # should be set/provided in gitlab project settings
//     - go run cicd/scripts/delete-gitlab-artifacts.go
//       --dry-run=false
//       --skip-tag-pipeline-jobs=true
//       --server=${CI_SERVER_URL}
//       --token=${GITLAB_PROJECT_ACCESS_TOKEN}
//       --project-id=${CI_PROJECT_ID}
//       --per-page=100
//       --pages=100
//       --dont-delete-younger-than=$((30*24))h
//       --dont-delete-older-than=$((12*30*24))h

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
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	Status            string     `json:"status"`
	FinishedAt        time.Time  `json:"finished_at"`
	ArtifactsExpireAt time.Time  `json:"artifacts_expire_at"`
	Artifacts         []Artifact `json:"artifacts"`
	Pipeline          Pipeline   `json:"pipeline"`
}

type Artifact struct {
	FileType   string `json:"file_type"`
	Size       uint64 `json:"size"`
	Filename   string `json:"filename"`
	FileFormat string `json:"file_format"`
}

type Pipeline struct {
	ID int `json:"id"`
	// "project_id": 1,
	Ref string `json:"ref"`
	// "sha": "0ff3ae198f8601a285adcf5c0fff204ee6fba5fd",
	// "status": "pending"
}

// https://docs.gitlab.com/ee/api/tags.html
type Tag struct {
	Name string `json:"name"`
}

// todo:
// * https://docs.gitlab.com/ee/api/rest/#offset-based-pagination VS https://docs.gitlab.com/ee/api/rest/#keyset-based-pagination
// * add `&sort=asc` / `&sort=desc` see: https://docs.gitlab.com/ee/api/rest/#keyset-based-pagination
// * use `x-total-pages` to get + use it to figure out how many pages should be used if given a date (recursive)

func main() {
	var dryRun bool
	var skipTagPipelineJobs bool
	var project_id int
	var per_page int
	var pages int
	var startPage int
	var token string
	var server string
	var dontDeleteYoungerThan time.Duration
	var dontDeleteOlderThan time.Duration

	flag.BoolVar(&dryRun, "dry-run", true, "Enable dry-run (execution without any deletions)")
	flag.BoolVar(&skipTagPipelineJobs, "skip-tag-pipeline-jobs", true, "Skip artifact deletion for jobs of pipelines which were triggered by tag")
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

	tags, err := collectProjectTags(project_id, server, token)
	if err != nil {
		log.Fatalf("couldn't collect tags: %s", err)
	}

	fmt.Printf("tags found: %v \n", tags)

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

		if len(jobs) == 0 {
			fmt.Printf("  stop execution as no jobs are found on page %d\n", i)
			break Loop
		}

		for _, job := range jobs {
			// fmt.Println("job to process:", job.ID)

			if job.FinishedAt.IsZero() {
				continue
			}

			if skipTagPipelineJobs && tagListContains(tags, job.Pipeline.Ref) {
				fmt.Printf("  skipped (skipTagPipelineJobs=true): id=%d status='%s' name='%s' job.Pipeline.Ref='%s'\n", job.ID, job.Status, job.Name, job.Pipeline.Ref)
				continue
			}

			// job is younger than dontDeleteYoungerThan but also has expired job artifact
			if job.FinishedAt.After(time.Now().Add(-1*dontDeleteYoungerThan)) && job.ArtifactsExpireAt.Before(time.Now()) {
				_, _ = deleteJobArtifacts(job, server, token, project_id)
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

			respStatus, err := deleteJobArtifacts(job, server, token, project_id)
			if err != nil {
				log.Fatalf("failed deleting job(id: %d) artifacts: %s", job.ID, err)
			}

			fmt.Printf("  job artifacts deleted: %d | response Status: %s \n", job.ID, respStatus)
		}
	}
}

func deleteJobArtifacts(job Job, server string, token string, projectId int) (string, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%d/jobs/%d/artifacts", server, projectId, job.ID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error constructing http request: %s", err)
	}

	req.Header.Set("PRIVATE-TOKEN", token)

	//artifact.FileType == "archive"
	for _, artifact := range job.Artifacts {
		fmt.Printf("      artifact.name: %s \n", artifact.Filename)
	}

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

func collectProjectTags(project_id int, server, token string) ([]Tag, error) {
	startPage := 0
	perPage := 1000
	maxPages := 100

	allTags := []Tag{}

	for i := startPage; i <= maxPages; i++ {
		url := fmt.Sprintf("%v/api/v4/projects/%v/repository/tags?per_page=%d&page=%d", server, project_id, perPage, i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []Tag{}, fmt.Errorf("failed building GET http request: %s", err)
		}
		req.Header.Set("PRIVATE-TOKEN", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []Tag{}, fmt.Errorf("failed doing GET http request: %s", err)
		}
		defer resp.Body.Close()

		var tags []Tag
		err = json.NewDecoder(resp.Body).Decode(&tags)
		if err != nil {
			return []Tag{}, fmt.Errorf("failed decoding tags from response: %s", err)
		}

		if len(tags) == 0 {
			break
		}

		allTags = append(allTags, tags...)
	}

	return allTags, nil
}

func tagListContains(tags []Tag, name string) bool {
	for _, tag := range tags {
		if tag.Name == name {
			return true
		}
	}

	return false
}
