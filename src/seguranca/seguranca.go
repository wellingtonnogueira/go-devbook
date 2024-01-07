package seguranca

import "golang.org/x/crypto/bcrypt"

// Hash recebe uma senha e transforma num hash
func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

// VerificarSenha valida se senha fornecida gera o mesmo hash da senha fornecida anteriormente
func VerificarSenha(senhaHash, senhaString string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(senhaString))
}