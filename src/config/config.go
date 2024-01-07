package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

var (
	// StringConexaoBanco é a string de conexão com MySQL
	StringConexaoBanco = ""

	// Porta onde a API vai estar em execução
	Porta = 0

	// Host configura o loopback do server
	Host = "http://172.27.55.252"

	// SecretKey é a chave que vai ser usado para assinar o token
	SecretKey []byte
)

// Carregar vai inicializar as variáveis de ambiente
func Carregar() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		Porta = 9000
	}

	StringConexaoBanco = fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USUARIO"),
		os.Getenv("DB_SENHA"),
		os.Getenv("DB_NOME"),
	)

	Host = os.Getenv("API_HOST")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}