package responses

import "time"

type HttpResponse struct {
	Time time.Time `json:"time"`
}

func (r HttpResponse) Zero() HttpResponse {
	return HttpResponse{
		Time: time.Now(),
	}
}

type ChatResponse struct {
	HttpResponse
	Message string `json:"message"`
}

func NewChatResponse(message string) *ChatResponse {
	return &ChatResponse{
		HttpResponse: HttpResponse{}.Zero(),
		Message:      message,
	}
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
