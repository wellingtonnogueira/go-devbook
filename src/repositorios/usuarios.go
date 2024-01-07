package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

type Usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositório de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// Criar insere usuários no banco de dados
func (repositorio Usuarios) Criar(usuario modelos.Usuario) (uint64, error) {

	statement, erro := repositorio.db.Prepare(
		`insert into usuarios (nome, nick, email, senha)
		 values (?, ?, ?, ?)`)

	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIdInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIdInserido), nil
}

// Buscar retorna usuários de acordo com filtros dados.
func (repositorio Usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick% . O escape, neste caso, para % é %% e para a string é %s

	linhas, erro := repositorio.db.Query(
		"select ID, nome, nick, email, criadoEm from usuarios where nome like ? or nick like ?",
		nomeOuNick, nomeOuNick,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	// var usuarios []modelos.Usuario
	// Ao invés de utilizar a declaração de um slice (como acima, comentado)
	// utiliza-se o make para garantir o retorno esperado de uma lista vazia.
	usuarios := make([]modelos.Usuario, 0)

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarPorId retorna dados de um usuário dado seu ID
func (repositorio Usuarios) BuscarPorId(usuarioId uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"select ID, nome, nick, email, criadoEm from usuarios where id = ?", usuarioId,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Atualizar dados de um usuário no banco de dados
func (repositorio Usuarios) Atualizar(usuarioId uint64, usuario modelos.Usuario) error {

	statement, erro := repositorio.db.Prepare(`update usuarios set nome=?, nick=?, email=? where id=?`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuarioId); erro != nil {
		return erro
	}

	return nil
}

// RemoverUsuario remove usuário da tabela
func (repositorio Usuarios) RemoverUsuario(usuarioId uint64) error {
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

// BuscarPorEmail busca um usuario dado seu email e retorna Id/senha (hash)
func (repositorio Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query(
		"select ID, email, senha, criadoEm from usuarios where email = ?", email,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(
			&usuario.ID,
			&usuario.Email,
			&usuario.Senha,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Seguir insere um novo seguidor na tabela seguidores
func (repositorio Usuarios) Seguir(usuarioId, seguidorId uint64) error {

	statement, erro := repositorio.db.Prepare(
		"insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)", // ignore: se já existir, não irá gerar erro. Somente ignorar.
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioId, seguidorId); erro != nil {
		return erro
	}

	return nil
}

// Parar de seguir remove um seguidor na tabela de seguidores
func (repositorio Usuarios) PararDeSeguir(usuarioId, seguidorId uint64) error {

	statement, erro := repositorio.db.Prepare(
		"delete from seguidores where usuario_id = ? and seguidor_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioId, seguidorId); erro != nil {
		return erro
	}

	return nil
}

// BuscarSeguidores retorna todos os seguidores de um determinado usuário
func (repositorio Usuarios) BuscarSeguidores(usuarioId uint64) ([]modelos.Usuario, error) {

	linhas, erro := repositorio.db.Query(
		`SELECT u.id, u.nome , u.nick, u.email from seguidores s 
		join usuarios u on u.id = s.seguidor_id  
		WHERE s.usuario_id = ?`,
		usuarioId,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	// var usuarios []modelos.Usuario
	// Ao invés de utilizar a declaração de um slice (como acima, comentado)
	// utiliza-se o make para garantir o retorno esperado de uma lista vazia.
	usuarios := make([]modelos.Usuario, 0)

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarSeguindo retorna todos os seguidores de um determinado usuário
func (repositorio Usuarios) BuscarSeguindo(usuarioId uint64) ([]modelos.Usuario, error) {

	linhas, erro := repositorio.db.Query(
		`SELECT u.id, u.nome , u.nick, u.email from seguidores s 
		join usuarios u on u.id = s.usuario_id   
		WHERE s.seguidor_id = ?`,
		usuarioId,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	// var usuarios []modelos.Usuario
	// Ao invés de utilizar a declaração de um slice (como acima, comentado)
	// utiliza-se o make para garantir o retorno esperado de uma lista vazia.
	usuarios := make([]modelos.Usuario, 0)

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repositorio Usuarios) BuscarSenha(usuarioId uint64) (string, error) {
	linha, erro := repositorio.db.Query("select senha from usuarios where id = ?", usuarioId)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	// var usuario modelos.Usuario
	var senha string

	if linha.Next() {
		// if erro = linha.Scan(&usuario.Senha); erro != nil {
		if erro = linha.Scan(&senha); erro != nil {
			return "", erro
		}
	}
	
	return senha, nil
}

func (repositorio Usuarios) AtualizarSenha(usuarioId uint64, senhaComHash string) error {

	statement, erro := repositorio.db.Prepare(
		"update usuarios set senha = ? where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(senhaComHash, usuarioId); erro != nil {
		return erro
	}

	return nil
}