package middleware



import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") 

func Authmiddlware(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		cookieToken, err := ctx.Cookie("Authorization")
		if err == nil {
			token = cookieToken
		}
	}
	if token == "" {
		ctx.JSON(400, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		ctx.JSON(401, gin.H{
			"message": "Invalid or expired token",
		})
		ctx.Abort()
		return
	}

	// Optionally, set claims to context
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		ctx.Set("claims", claims)
	}

	ctx.Next()
}
//fix steps for jwt auth
// 1)get token from header ctx.getHeader("Auth...")
//2)remove prefix using strings.TrimPerfix(token,"Bearer ")
//3)parsed token by jwt.parse(token,func(token *jwt.Token)(interface{},error){
//if _,ok:=token.method.(*jwt.SingingwithHMAC);!ok{
//return nil,jwt.ErrSignatureInvalid
// }
//and parsedTOKen=parsedTOken.Clims.(Jwt.MapCliams)
//})