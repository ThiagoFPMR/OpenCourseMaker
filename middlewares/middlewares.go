package middlewares

import (
	"fmt"
	"github.com/ThiagoFPMR/OpenCourseMaker/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verifica se o cookie de autenticação existe na solicitação
		tokenString := user.ExtractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Verifica se o token JWT é válido
		if err := user.TokenValid(tokenString); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Se o token JWT é válido, define o usuário atual na variável de contexto
		// Aqui estou usando o nome de usuário armazenado no token JWT, mas você pode armazenar outros dados no token
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		claims := token.Claims.(jwt.MapClaims)
		nome := claims["nome"].(string)
		email := claims["email"].(string)
		c.Set("nome", nome)
		c.Set("email", email)

		// Chama a próxima função de tratamento de solicitações na cadeia de middleware
		c.Next()
	}
}
