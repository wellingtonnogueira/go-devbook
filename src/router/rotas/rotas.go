package rotas

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

//Rota representa as rotas da API
type Rota struct {
	URI string
	Metodo string
	Funcao func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

// Configurar coloca rotas dentro do router
func Configurar(r *mux.Router) *mux.Router {
	rotas := rotaUsuarios
	rotas = append(rotas, rotaLogin)
	rotas = append(rotas, rotasPublicacoes...)

	for _, rota := range rotas {

		funcao := rota.Funcao

		if rota.RequerAutenticacao {
			funcao = middlewares.Autenticar(rota.Funcao)
		}
		
		r.HandleFunc(rota.URI, middlewares.Logger(funcao)).Methods(rota.Metodo)
	}

	return r
}