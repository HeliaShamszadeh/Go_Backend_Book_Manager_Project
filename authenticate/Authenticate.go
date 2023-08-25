package authenticate

import (
	"bookman/db"
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	IncorrectPasswordErr = errors.New("incorrect password")
	EmptyTokenErr        = errors.New("empty token string")
	InvalidTokenErr      = errors.New("invalid access token")
	CannotValidateToken  = errors.New("cannot validate token")
)

type Authenticate struct {
	database              *db.GormDB
	logger                *logrus.Logger
	secretKey             []byte
	jwtExpirationDuration time.Duration
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	TokenString string
}

type claims struct {
	jwt.MapClaims
	Username string `json:"username"`
}

func NewAuth(db *db.GormDB, jwtExpirationDuration time.Duration, logger *logrus.Logger) (*Authenticate, error) {
	secretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}

	// Check database
	if db == nil {
		return nil, errors.New("database cannot be nil")
	}

	return &Authenticate{
		database:              db,
		secretKey:             secretKey,
		jwtExpirationDuration: jwtExpirationDuration,
		logger:                logger,
	}, nil
}

func (a *Authenticate) Login(cred Credential) (*Token, error) {
	account, err := a.database.GetUserByUsername(cred.Username)
	if err != nil {
		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(cred.Password))
	if err != nil {
		return nil, IncorrectPasswordErr
	}

	// Create JWT token
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Username: cred.Username,
		MapClaims: jwt.MapClaims{
			"expired_at": expirationTime.Unix(),
		},
	})
	tokenString, err := tokenJWT.SignedString(a.secretKey)
	if err != nil {
		return nil, err
	}

	return &Token{TokenString: tokenString}, err
}
func generateRandomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}
	return jwtKey, nil
}

func (a *Authenticate) GetUsernameByToken(token string) (username string, err error) {

	if token == "" {
		return "", EmptyTokenErr
	}
	c := &claims{}

	jwtToken, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", InvalidTokenErr
		} else {
			return "", CannotValidateToken
		}
	}

	if !jwtToken.Valid {
		return "", InvalidTokenErr
	} else {
		return c.Username, nil
	}

}
