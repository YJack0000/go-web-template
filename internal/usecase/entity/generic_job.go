package entity

type GenericJob struct {
	ID     string `json:"jobId"       example:"12345"`
	Name   string `json:"jobName"       example:"name"`
	Status string `json:"jobStatus"       example:"jobStatus"`
}
