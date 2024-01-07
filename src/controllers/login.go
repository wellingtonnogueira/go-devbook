package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Login é responsável pela autenticação do usuário
func Login(w http.ResponseWriter, r *http.Request) {
	corpoDaRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoDaRequisicao, &usuario); erro != nil {
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

	usuarioSalvo, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.ERRO(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvo.Senha, usuario.Senha); erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, errors.New("usuario ou senha inválidos"))
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvo.ID)
	if erro != nil {
		respostas.ERRO(w, http.StatusUnauthorized, errors.New("usuario ou senha inválidos"))
		return
	}

	fmt.Println(usuarioSalvo, token)

	var respostaToken respostaToken

	respostaToken.Token = token

	respostas.JSON(w, http.StatusOK, respostaToken, nil)
}

type respostaToken struct {
	Token string `json:"token"`
}