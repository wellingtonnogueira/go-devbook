package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuário na rede social
type Usuario struct {
	ID uint64 `json:"id,omitempty"`
	Nome string `json:"nome,omitempty"`
	Nick string `json:"nick,omitempty"`
	Email string `json:"email,omitempty"`
	Senha string `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoEm,omitempty"`
}

// Preparar irá validar e formatar dados do usuário
func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}
	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(etapa string) error {

	if usuario.Nome == "" {
		return errors.New("nome obrigatório")
	}

	if usuario.Nick == "" {
		return errors.New("nick obrigatório")
	}

	if usuario.Email == "" {
		return errors.New("email obrigatório")
	}
	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("formato de email inválido")
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("senha obrigatório")
	}

	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if etapa == "cadastro" {
		senhaHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}
		usuario.Senha = string(senhaHash)
	}
	return nil
}