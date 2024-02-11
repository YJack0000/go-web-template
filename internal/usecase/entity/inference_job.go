package entity

type InferenceJob struct {
	Job        GenericJob `json:"job"`
	TwccCCSId  string     `json:"jobId" example:"12345"`
	EntryPoint string     `json:"entryPoint" example:"12345"`
}
