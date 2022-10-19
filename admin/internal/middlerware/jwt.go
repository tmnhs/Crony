package middlerware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/common/pkg/logger"
	"golang.org/x/sync/singleflight"
	"time"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	ID       int
	UserName string
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

var control = &singleflight.Group{}

func NewJWT() *JWT {
	return &JWT{
		[]byte("S0dEdN9tqG0AAAAHdElNRQfmCgwBDCSd2zTMAAAA"),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: 86400, // buffer time 1 day buffer time will get a new token refresh token. In this case, a user will have two valid tokens, but only one will be left at the front end and the other will be lost.
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,   // effective time of signature
			ExpiresAt: time.Now().Unix() + 604800, // Expiration time 7 days profile
			Issuer:    "tmnhs",                    // the publisher of the signature
		},
	}
	return claims
}

// create a token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// Replacing old token with new token using merging and origin-pull to avoid concurrency problems
func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
	v, err, _ := control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

func GetClaims(c *gin.Context) (*CustomClaims, error) {
	token := c.Request.Header.Get("Authorization")
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		logger.GetLogger().Error("Failed to obtain parsing information from jwt from Context of Gin. Please check whether Authorization exists in the request header and whether claims is the specified structure.")
	}
	return claims, err
}

// Get the user roles parsed from jwt from the Context of Gin
func GetUserInfo(c *gin.Context) *CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// We have jwt authentication header information to return token information when Authorization logs in. Here,
		//the front end needs to store the token in cookie or local localStorage,
		//but you need to negotiate the expiration time with the back end.
		//You can agree to refresh the token or log in again.
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			resp.FailWithDetailed(resp.ERROR, gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				resp.FailWithDetailed(resp.ERROR, gin.H{"reload": true}, "授权已过期", c)
				c.Abort()
				return
			}
			resp.FailWithDetailed(resp.ERROR, gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
