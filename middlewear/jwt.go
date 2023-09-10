package middlewear

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var (
	TokenExpired = errors.New("Token is expired")
)

//授权

// 指定加密密钥
var jwtSecret = []byte("Khalil")

// 用户的状态和元数据
type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

// 根据用户的用户名产生token
func GenerateToken(userId uint, iss string) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	expiredTime := nowTime.Add(7 * 24 * time.Hour)

	claims := Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			//过期时间
			ExpiresAt: expiredTime.Unix(),
			//指定token发行人
			Issuer: iss,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//生成字符串然后获取完整的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 根据传入的token值获取到claims对象信息
func ParseToken(token string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最后返回token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// 鉴权
func JWY() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		user := c.Query("userId")
		userId, err := strconv.Atoi(user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "your userId is not legitimate",
			})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "please login",
			})
			c.Abort()
			return
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "token is invalid",
				})
				c.Abort()
				return
			} else if time.Now().Unix() > claims.ExpiresAt {
				err = TokenExpired
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "the authorization has expired",
				})
				c.Abort()
				return
			}
			if claims.UserID != uint(userId) {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "your login is invalid",
				})
				c.Abort()
				return
			}

			fmt.Println("token success")
			c.Next()
		}
	}
}
