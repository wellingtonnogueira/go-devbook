package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON formata a resposta sendo devolvida no response do request
func JSON(w http.ResponseWriter, statusCode int, dados interface{}, headers map[string] string) {

	for chave, valor := range headers {
		w.Header().Set(chave, valor)
	}	
	
	if dados != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode) // Tem que ser definido após os demais headers.
		
		if erro := json.NewEncoder(w).Encode(dados); erro != nil {
			log.Fatal(erro)
		}
	} else { // else necessário para evitar warning log `http: superfluous response.WriteHeader call from api/src/respostas.JSON (respostas.go:xx)`
		w.WriteHeader(statusCode)
	}
	
}

// ERRO formata a resposta de erro durante a execução da aplicação, sendo devolvida no response do request
func ERRO(w http.ResponseWriter, statusCode int, erro error) {
	// Aqui está sendo criado um struct anônimo que representa o erro.
	// para fins de aprendizado, está sendo exposto erro interno.
	// No mundo corporativo esta não é uma boa prática, sendo geralmente convencionado
	// logar o erro e trocar a mensagem de erro para uma mensagem de erro negocial.
	/* O struct utilizado aqui representa o seguinte código:
	type Erro struct {
		Erro string `json:"erro"`
	}
	E seu uso equivale a:
	var erro Erro
	erro.erro = erro.Error()
	*/
	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	}, map[string] string{"content-type": "application/json"})
}