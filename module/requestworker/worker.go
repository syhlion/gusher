package requestworker

import (
	"io/ioutil"
	"net/http"
)

type Worker struct {
	JobQuene   chan *Job
	Threads    int
	HttpClient *http.Client
}

func (w *Worker) dispatcher() {
	for j := range w.JobQuene {
		res, err := w.HttpClient.Do(j.Resq)
		if err != nil {
			j.Result <- Result{err}
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			j.Result <- Result{err}
			continue
		}

		if ret := string(body); ret != "ok" {
			e := &ErrorResult{ret}
			j.Result <- Result{e}
		}
		j.Result <- Result{nil}
	}

}

func (w *Worker) Start() {

	for i := 0; i < w.Threads; i++ {
		go w.dispatcher()
	}

}
