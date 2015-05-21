package shared

// Job represent a job
// in the database
type Job struct {
	ID   string      `gorethink:"id,omitempty" json:"id"`
	Name string      `gorethink:"name" json:"name"`
	Data interface{} `gorethink:"data,omitempty" json:"data"`
}
