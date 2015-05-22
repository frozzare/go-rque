package rque

import r "github.com/dancannon/gorethink"

// Config represent the
// configurations of the queue
type Config struct {
	Address  string
	Database string
	Table    string
}

// Job represent a job
// in the database
type Job struct {
	ID   string      `gorethink:"id,omitempty" json:"id"`
	Name string      `gorethink:"name" json:"name"`
	Data interface{} `gorethink:"data,omitempty" json:"data"`
}

// Que watches a job table
type Que struct {
	config  Config
	session *r.Session
	jobs    chan Job
	quit    chan bool
}

// New creates a new Que
func New(config Config) (*Que, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  config.Address,
		Database: config.Database,
	})

	if err != nil {
		return nil, err
	}

	que := Que{
		config:  config,
		session: session,
		jobs:    make(chan Job),
		quit:    make(chan bool),
	}

	// Run job fetching functions
	go que.findJobs()
	go que.runLeftovers()

	// Wait for quit
	go func() {
		<-que.quit
		que.session.Close()
	}()

	return &que, nil
}

// Jobs returns the jobs channel
func (q *Que) Jobs() <-chan Job {
	return q.jobs
}

// Quit stops the que watcher
func (q *Que) Quit() {
	q.quit <- true
	close(q.jobs)
}

// PostJob add a job to the queue
func (q *Que) PostJob(job Job) error {
	_, err := r.Table(q.config.Table).Insert(job).RunWrite(q.session)
	return err
}

// deleteJob will delete a job and
func (q *Que) deleteJob(job Job) error {
	_, err := r.Table(q.config.Table).Get(job.ID).Delete().RunWrite(q.session)
	return err
}

// runLeftovers will run all existing job in the database
// before the changes feed.
func (q *Que) runLeftovers() error {
	res, err := r.Table(q.config.Table).Run(q.session)
	defer res.Close()

	if err != nil {
		return err
	}

	var jobs []Job
	res.All(&jobs)

	for _, job := range jobs {
		q.jobs <- job
		q.deleteJob(job)
	}
	return nil
}

// findJobs will find all new jobs from the changes feed.
func (q *Que) findJobs() error {
	jobs, err := r.Table(q.config.Table).Changes().Field("new_val").Run(q.session)
	defer jobs.Close()

	if err != nil {
		return err
	}

	var job Job
	for jobs.Next(&job) {
		if job.Name != "" {
			q.jobs <- job
		}
		q.deleteJob(job)
	}
	return nil
}
