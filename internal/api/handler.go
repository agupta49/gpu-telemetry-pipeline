package api

import "fmt"

type GPU struct {
	ID string `json:"id"`
	Model string `json:"model"`
}

func NewGPU(id, model string) *GPU {
	return &GPU{ID: id, Model: model}
}

func (g *GPU) String() string {
	return fmt.Sprintf("GPU{id=%s, model=%s}", g.ID, g.Model)
}

func HealthStatus() string {
	return "ok"
}

func FilterGPUs(gpus []GPU, model string) []GPU {
	result := []GPU{}
	for _, g := range gpus {
		if g.Model == model {
			result = append(result, g)
		}
	}
	return result
}
