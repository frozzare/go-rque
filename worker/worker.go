package worker

import (
	"log"

	r "github.com/dancannon/gorethink"
	"github.com/frozzare/go-queue/datastore"
	"github.com/frozzare/go-queue/shared"
)

// deleteJob will delete a job and
func deleteJob(c shared.Config, job shared.Job) {
	_, err := r.Table(c.Table).Get(job.ID).Delete().Run(datastore.Instance(c))

	if err != nil {
		panic(err)
	}
}

// Run the worker.
func Run(c shared.Config) {
	quit := make(chan bool, 1)

	jobs, err := r.Table(c.Table).Changes().Field("new_val").Run(datastore.Instance(c))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var job shared.Job
		for jobs.Next(&job) {
			c.Emitter.Emit(job.Name, job)
			deleteJob(c, job)
		}
	}()

	<-quit
}
