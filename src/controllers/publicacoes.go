package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/config"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarPublicacao adiciona uma publicação
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	var publicacao modelos.Publicacao

	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	publicacao.AuthorId = usuarioId

	if erro = publicacao.Preparar(); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacaoId, erro := repositorio.Criar(publicacao)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	host := config.Host
	portaApi := config.Porta

	headers := map[string]string{
		"location": fmt.Sprintf("%s:%d/publicacoes/%d", host, portaApi, publicacaoId),
	}
	respostas.JSON(w, http.StatusCreated, nil, headers)
}

// BuscarPublicacoes retorna todas as publicações
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.Buscar(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, publicacoes, nil)
}

// BuscarPublicacao retorna uma publicação
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacao, erro := repositorio.BuscarPorId(publicacaoId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	if publicacao.ID == 0 {
		respostas.ERRO(w, http.StatusNotFound, errors.New("publicação não encontrada"))
		return
	}

	respostas.JSON(w, http.StatusOK, publicacao, nil)
}

// AtualizarPublicacao atualiza uma publicação
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacaoSalvaNoBanco, erro := repositorio.BuscarPorId(publicacaoId)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if publicacaoSalvaNoBanco.AuthorId != usuarioId {
		respostas.ERRO(w, http.StatusForbidden, errors.New("não foi possível atualizar a publicação"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao

	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if erro = publicacao.Preparar(); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.Atualizar(publicacaoId, publicacao); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil, nil)
}

// RemoverPublicacao remove uma publicação
func RemoverPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacaoSalvaNoBanco, erro := repositorio.BuscarPorId(publicacaoId)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if publicacaoSalvaNoBanco.AuthorId != usuarioId {
		respostas.ERRO(w, http.StatusForbidden, errors.New("não foi possível remover a publicação"))
		return
	}

	if erro = repositorio.RemoverPublicacao(publicacaoId); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil, nil)
}

// BuscarPublicacoesPorUsuario traz todas as publicações de um determinado usuário
func BuscarPublicacoesPorUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacoes, erro := repositorio.BuscarPorUsuario(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, publicacoes, nil)
}

// CurtirPublicacao incrementa curtida em 1
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	if erro = repositorio.Curtir(publicacaoId); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	host := config.Host
	portaApi := config.Porta

	headers := map[string]string{
		"location": fmt.Sprintf("%s:%d/publicacoes/%d", host, portaApi, publicacaoId),
	}
	respostas.JSON(w, http.StatusNoContent, nil, headers)
}

// DescurtirPublicacao decrementa curtida em 1
func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	if erro = repositorio.Descurtir(publicacaoId); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	host := config.Host
	portaApi := config.Porta

	headers := map[string]string{
		"location": fmt.Sprintf("%s:%d/publicacoes/%d", host, portaApi, publicacaoId),
	}
	respostas.JSON(w, http.StatusNoContent, nil, headers)
}