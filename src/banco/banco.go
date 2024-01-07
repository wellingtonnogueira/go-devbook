package banco

import (
	"api/src/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //Driver, import implícito
)

// Conectar abre conexão com base de dados
func Conectar() (*sql.DB, error) {

	db, erro := sql.Open("mysql", config.StringConexaoBanco)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}