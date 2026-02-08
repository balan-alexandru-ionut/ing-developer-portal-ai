package responses

import "time"

type HttpResponse struct {
	Time time.Time `json:"time"`
}

type GenerationStatus string

const (
	NotStarted GenerationStatus = "not_started"
	Generating GenerationStatus = "generating"
	Generated  GenerationStatus = "generated"
	Formatting GenerationStatus = "formatting"
	Done       GenerationStatus = "done"
)

type GenerationStatusResponse struct {
	HttpResponse
	Status GenerationStatus `json:"status"`
}

func NewGenerationStatusResponse(status GenerationStatus) GenerationStatusResponse {
	return GenerationStatusResponse{
		HttpResponse: HttpResponse{
			Time: time.Now(),
		},
		Status: status,
	}
}

type GeneratedFile struct {
	FilePath string `json:"filePath"`
	Code     string `json:"code"`
}

type GenerationResponse struct {
	HttpResponse
	Files []GeneratedFile `json:"files"`
}
