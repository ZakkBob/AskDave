package daveapiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ZakkBob/AskDave/crawler/fetcher"
	"github.com/ZakkBob/AskDave/crawler/taskrunner"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
)

type DaveApiClient struct {
	TaskRunner taskrunner.TaskRunner
	Url        string
}

func (d *DaveApiClient) FetchTasks() *tasks.Tasks {
	f := fetcher.NetFetcher{}
	resp, err := f.Fetch(d.Url + "/work")
	if err != nil {
		fmt.Print(err.Error())
	}
	data := resp.Body

	fmt.Println(data)

	var t tasks.Tasks
	json.Unmarshal([]byte(data), &t)
	return &t
}

func (d *DaveApiClient) UploadTasks() error {
	jsonData, _ := json.MarshalIndent(&d.TaskRunner.Results, "", "  ")
	postData := strings.NewReader(string(jsonData))
	resp, err := http.Post(d.Url+"/work", "application/json", postData)

	if err != nil {
		return fmt.Errorf("uploading tasks: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("uploading tasks: %w", err)
	}

	fmt.Print(string(respBody))
	return nil
}

func (d *DaveApiClient) Run() {
	t := d.FetchTasks()
	d.TaskRunner.Tasks = *t

	f := fetcher.NetFetcher{
		Debug: true,
	}

	d.TaskRunner.Fetcher = &f
	d.TaskRunner.Results = tasks.Results{
		Robots:     make(map[string]*tasks.RobotsResult),
		Pages:      make(map[string]*tasks.PageResult),
		RobotsChan: make(chan *tasks.RobotsResult, 5),
		// SitemapsChan: make(chan *string, 5),
		PagesChan:      make(chan *tasks.PageResult, 5),
		RobotsFinished: make(chan bool, 1),
		PagesFinished:  make(chan bool, 1),
	}

	d.TaskRunner.Run(100)

	err := d.UploadTasks()

	if err != nil {
		fmt.Println(err.Error())
	}

	j, _ := json.MarshalIndent(&d.TaskRunner.Results, "", "  ")

	fmt.Println(string(j))
}

func Create(u string) DaveApiClient {
	return DaveApiClient{
		TaskRunner: taskrunner.TaskRunner{
			Tasks:   tasks.Tasks{},
			Results: tasks.Results{},
			Fetcher: &fetcher.NetFetcher{},
		},
		Url: u,
	}
}
