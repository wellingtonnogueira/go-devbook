package main

import (
	"crypto/rand"
	"encoding/base64"
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func init() {
	config.Carregar()
	// Somente para criar a secret key a ser utilizada no package autorizacao, na geração de token
	if config.RunInit {
		chave := make([]byte, 64)
		if _, erro := rand.Read(chave); erro != nil {
			log.Fatal(erro)
		}
	
		stringBase64 := base64.StdEncoding.EncodeToString(chave)
		fmt.Println(stringBase64)
	}
}

func main() {
	host := config.Host
	portaApi := config.Porta

	fmt.Println("Rodando API na porta",host, portaApi)
	
	r := router.Gerar()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portaApi), r))
}
