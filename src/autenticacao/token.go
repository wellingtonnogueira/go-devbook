package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go" //Alias para jwt
)

func CriarToken(usuarioId uint64) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
		
	return token.SignedString(config.SecretKey) //secret que virá do .env
}

func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)

	token, erro := jwt.Parse(tokenString, retornarChaveDeVerirficacao)
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token inválido")

}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	splittedToken := strings.Split(token, " ")

	if len(splittedToken) == 2 {
		return splittedToken[1]
	}

	return ""
}

func retornarChaveDeVerirficacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //ok é um booleano
		return nil, fmt.Errorf("metodo de assinatura inesperado! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}

func ExtrairUsuarioId(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r)

	token, erro := jwt.Parse(tokenString, retornarChaveDeVerirficacao)
	if erro != nil {
		return 0, erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		permissaoUsuarioId := fmt.Sprintf("%.0f", permissoes["usuarioId"])
		usuarioId, erro := strconv.ParseUint(permissaoUsuarioId, 10, 64)
		if erro != nil {
			return 0, erro
		}

		return usuarioId, nil
	}

	return 0, errors.New("token inválido")

}