package entity

type TwccJob struct {
	Job       GenericJob `json:"job"`
	TwccJobId string     `json:"jobId" example:"12345"`
}

type ContainerJob struct {
	Job             GenericJob `json:"job"`
	DockerImageName string     `json:"dockerImageName" example:"ubuntu:latest"`
}
