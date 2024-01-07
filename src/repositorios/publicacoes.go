package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type Publicacoes struct {
	db *sql.DB
}

// NovoRepositorioDePublicacoes cria um repositório de publicações
func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

// Criar insere publicação no banco de dados
func (repositorio Publicacoes) Criar(publicacao modelos.Publicacao) (uint64, error) {

	statement, erro := repositorio.db.Prepare(
		`insert into publicacoes (titulo, conteudo, autor_id)
		 values (?, ?, ?)`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.AuthorId)
	if erro != nil {
		return 0, erro
	}

	ultimoIdInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIdInserido), nil
}

// Buscar retorna publicações de acordo com filtros dados.
func (repositorio Publicacoes) Buscar(usuarioId uint64) ([]modelos.Publicacao, error) {

	linhas, erro := repositorio.db.Query(
		`select distinct p.*, u.nick from publicacoes p 
		inner join usuarios u on u.id = p.autor_id 
		inner join seguidores s on s.usuario_id = p.autor_id  
		where u.id = ? or s.seguidor_id = ?
		order by p.id desc`, usuarioId, usuarioId,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	publicacoes := make([]modelos.Publicacao, 0)

	for linhas.Next() {
		var publicacao modelos.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AuthorId,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

// BuscarPorId retorna dados de uma publicação dado seu ID
func (repositorio Publicacoes) BuscarPorId(publicacaoId uint64) (modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(
		`select p.*, u.nick from publicacoes p
		inner join usuarios u on u.id = p.autor_id
		where p.id = ?`, publicacaoId,
	)
	if erro != nil {
		return modelos.Publicacao{}, erro
	}
	defer linhas.Close()

	var publicacao modelos.Publicacao

	if linhas.Next() {
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AuthorId,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AuthorNick,
		); erro != nil {
			return modelos.Publicacao{}, erro
		}
	}

	return publicacao, nil
}

// Atualizar dados de um usuário no banco de dados
func (repositorio Publicacoes) Atualizar(usuarioId uint64, publicacao modelos.Usuario) error {

	statement, erro := repositorio.db.Prepare(`update usuarios set nome=?, nick=?, email=? where id=?`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publicacao.Nome, publicacao.Nick, publicacao.Email, usuarioId); erro != nil {
		return erro
	}

	return nil
}

// RemoverUsuario remove usuário da tabela
func (repositorio Publicacoes) RemoverUsuario(usuarioId uint64) error {
	statement, erro := repositorio.db.Prepare(`delete from usuarios where id=?`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioId); erro != nil {
		return erro
	}

	return nil
}
