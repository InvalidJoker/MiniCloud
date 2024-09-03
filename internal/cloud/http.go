package cloud

import "net/http"

type BackendService struct {
	DockerService *DockerService
}

func NewBackendService(dockerService *DockerService) *BackendService {
	return &BackendService{
		DockerService: dockerService,
	}
}

func (b *BackendService) Start() {
	http.ListenAndServe(":8080", nil)

}

func (b *BackendService) start(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
