package modelos

// Senha representa o formato da atualização de senha
type Senha struct {
	Nova string `json:"nova"`
	Atual string `json:"atual"`
}