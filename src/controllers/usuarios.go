package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/config"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario recurso para criar usuário
func CriarUsuario(w http.ResponseWriter, r *http.Request) {

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioId, erro := repositorio.Criar(usuario)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	host := config.Host
	portaApi := config.Porta

	headers := map[string]string{
		"location": fmt.Sprintf("%s:%d/usuarios/%d", host, portaApi, usuarioId),
	}
	respostas.JSON(w, http.StatusCreated, nil, headers)
}

// BuscarUsuarios recurso para buscar todos os usuários
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios, nil)
}

// BuscarUsuario recurso para buscar dados de um usuário
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
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

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuario, erro := repositorio.BuscarPorId(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	if usuario.ID == 0 {
		respostas.ERRO(w, http.StatusNotFound, errors.New("usuário inexistente"))
		return
	}

	respostas.JSON(w, http.StatusOK, usuario, nil)
}

// AtualizarUsuario recurso para atualizar dados de um usuário
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIdNoToken, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioId != usuarioIdNoToken {
		respostas.ERRO(w, http.StatusForbidden, errors.New("token inválido"))
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuarioNaBase, erro := repositorio.BuscarPorId(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	if usuarioNaBase.ID == 0 {
		respostas.ERRO(w, http.StatusBadRequest, errors.New("usuário inexistente"))
		return
	}

	if erro = repositorio.Atualizar(usuarioId, usuario); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil, nil)
}

// RemoverUsuario recurso para remover os dados de um usuário
func RemoverUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIdNoToken, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioId != usuarioIdNoToken {
		respostas.ERRO(w, http.StatusForbidden, errors.New("token inválido"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuarioNaBase, erro := repositorio.BuscarPorId(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	if usuarioNaBase.ID == 0 {
		respostas.ERRO(w, http.StatusBadRequest, errors.New("usuário inexistente"))
		return
	}

	if erro = repositorio.RemoverUsuario(usuarioId); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil, nil)
}

// SeguirUsuario permite que um usuário siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorId, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioId == seguidorId {
		respostas.ERRO(w, http.StatusForbidden, errors.New("não é possível seguir a si mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro := repositorio.Seguir(seguidorId, usuarioId); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil, nil)
}

// PararDeSeguirUsuario permite um usuário deixar de seguir outro
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)

	seguidorId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorId == usuarioId {
		respostas.ERRO(w, http.StatusForbidden, errors.New("não é possível parar de seguir a si mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	if erro := repositorio.PararDeSeguir(usuarioId, seguidorId); erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil, nil)
}

// BuscarSeguidores traz todos os seguidores de um usuário
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
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

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores, nil)
}

// BuscarSeguidos retorna a lista dos usuarios que seguem o usuário do request
func BuscarSeguidos(w http.ResponseWriter, r *http.Request) {
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

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguindo(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores, nil)
}

// AtualizarSenha atualiza a senha do usuário
func AtualizarSenha (w http.ResponseWriter, r *http.Request) {

	usuarioIdNoToken, erro := autenticacao.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}	

	if usuarioId != usuarioIdNoToken {
		respostas.ERRO(w, http.StatusForbidden, errors.New("não foi possível atualizar a senha"))
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequest, &senha); erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioId)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, errors.New("senha inválida"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.ERRO(w, http.StatusBadRequest, erro)
		return
	}

	erro = repositorio.AtualizarSenha(usuarioId, string(senhaComHash)); if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, nil, nil)
}