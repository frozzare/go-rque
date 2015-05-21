package worker

import (
	"log"

	r "github.com/dancannon/gorethink"
	"github.com/frozzare/go-rque/datastore"
	"github.com/frozzare/go-rque/shared"
)

// deleteJob will delete a job and
func deleteJob(c shared.Config, job shared.Job) {
	_, err := r.Table(c.Table).Get(job.ID).Delete().Run(datastore.Instance(c))

	if err != nil {
		panic(err)
	}
}

// runLeftovers will run all existing job in the database
// before the changes feed.
func runLeftovers(c shared.Config) {
	res, err := r.Table(c.Table).Run(datastore.Instance(c))

	if err != nil {
		log.Fatal(err)
	}

	var jobs []shared.Job
	res.All(&jobs)

	for _, job := range jobs {
		c.Emitter.Emit(job.Name, job)
		deleteJob(c, job)
	}
}

// findJobs will find all new jobs from the changes feed.
func findJobs(c shared.Config) {
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
}

// Run the worker.
func Run(c shared.Config) {
	quit := make(chan bool, 1)

	go runLeftovers(c)
	go findJobs(c)

	<-quit
}
