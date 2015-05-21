package rque

import (
	"log"

	r "github.com/dancannon/gorethink"
	"github.com/frozzare/go-emitter"
)

// Config represent the
// configurations of the queue
type Config struct {
	Address  string
	Database string
	Emitter  *emitter.Emitter
	Table    string
}

// Job represent a job
// in the database
type Job struct {
	ID   string      `gorethink:"id,omitempty" json:"id"`
	Name string      `gorethink:"name" json:"name"`
	Data interface{} `gorethink:"data,omitempty" json:"data"`
}

var session *r.Session

// sessionInstance will return RethinkDB session
func sessionInstance(c Config) *r.Session {
	if session == nil {
		sen, err := r.Connect(r.ConnectOpts{
			Address:  c.Address,
			Database: c.Database,
		})

		if err != nil {
			log.Fatalf("Error connecting to DB: %s", err)
		}

		session = sen

		return session
	}

	return session
}

// deleteJob will delete a job and
func deleteJob(c Config, job Job) {
	_, err := r.Table(c.Table).Get(job.ID).Delete().Run(sessionInstance(c))

	if err != nil {
		panic(err)
	}
}

// runLeftovers will run all existing job in the database
// before the changes feed.
func runLeftovers(c Config) {
	res, err := r.Table(c.Table).Run(sessionInstance(c))

	if err != nil {
		log.Fatal(err)
	}

	var jobs []Job
	res.All(&jobs)

	for _, job := range jobs {
		c.Emitter.Emit(job.Name, job)
		deleteJob(c, job)
	}
}

// findJobs will find all new jobs from the changes feed.
func findJobs(c Config) {
	jobs, err := r.Table(c.Table).Changes().Field("new_val").Run(sessionInstance(c))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var job Job
		for jobs.Next(&job) {
			if job.Name != "" {
				c.Emitter.Emit(job.Name, job)
			}
			deleteJob(c, job)
		}
	}()
}

// Run the worker.
func Run(c Config) {
	quit := make(chan bool, 1)

	go runLeftovers(c)
	go findJobs(c)

	<-quit
}
