package middlewares

import (
	"api/src/autenticacao"
	"api/src/respostas"
	"fmt"
	"net/http"
)

// Logger escreve informações da requisição no terminal
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s%s", r.Method, r.Host, r.RequestURI)
		fmt.Println()
		next(w, r)
	}
}

// Autenticar verifica se usuário fazendo a requisição está autenticado
func Autenticar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := autenticacao.ValidarToken(r); erro != nil {
			respostas.ERRO(w, http.StatusUnauthorized, erro)
			return
		}
		next(w, r)
	}
}