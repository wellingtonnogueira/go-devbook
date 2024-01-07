package modelos

import (
	"errors"
	"strings"
	"time"
)

// Publicacao representa uma publicação realizada por um usuario
type Publicacao struct {
	ID uint64 `json:"id,omitempty"`
	Titulo string `json:"titulo,omitempty"`
	Conteudo string `json:"conteudo,omitempty"`
	AuthorId uint64 `json:"autorId,omitempty"`
	AuthorNick string `json:"authorNick,omitempty"`
	Curtidas uint64 `json:"curtidas"`
	CriadaEm time.Time `json:"criadaEm,omitempty"`
}

//Preparar valida e formata dados da publicação
func (publicacao *Publicacao) Preparar() error {
	if erro := publicacao.validar(); erro != nil {
		return erro
	}

	publicacao.formatar()

	return nil
}

func (publicacao *Publicacao) validar() error {
	if publicacao.Titulo == "" {
		return errors.New("título não pode estar em branco")
	}
	if publicacao.Conteudo == "" {
		return errors.New("conteúdo não pode estar em branco")
	}
	return nil
}

func (publicacao *Publicacao) formatar() {
	publicacao.Titulo = strings.TrimSpace(publicacao.Titulo)
	publicacao.Conteudo = strings.TrimSpace(publicacao.Conteudo)
}