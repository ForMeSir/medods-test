package service

import (
	"encoding/base64"
	"errors"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)
 const(
	signingKey = "dfghnwmsjn12kas72hnmk9"
	//refsigningKey = "cvdgyoazlxk26782iehdg7u"
 )

 type tokenClaims struct{
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
 }


 func NewService() *Service {
	var servi *Service
	 return &Service{
				 Autorization: servi,
	 }
 }

 func (s *Service) RefreshToken( userUUID uuid.UUID,sessionUUID uuid.UUID, TimeLive time.Duration)(string){
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims {
		jwt.StandardClaims{
		ExpiresAt: time.Now().Add(TimeLive).Unix(),
	  IssuedAt: time.Now().Unix(),
	}, userUUID,sessionUUID,
})
a,_:= token.SignedString([]byte(signingKey))
	return 	a
 }
 
func(s *Service) GenerateToken(userUUID uuid.UUID, TimeLive time.Duration, RefTimeLive time.Duration) (string, string){
	sessionId:= uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims {
		jwt.StandardClaims{
		ExpiresAt: time.Now().Add(TimeLive).Unix(),
	  IssuedAt: time.Now().Unix(),
	}, userUUID,sessionId,
})

reftoken := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims {
	jwt.StandardClaims{
	ExpiresAt: time.Now().Add(RefTimeLive).Unix(),
	IssuedAt: time.Now().Unix(),
}, userUUID,sessionId,
})
ref,_:=reftoken.SignedString([]byte(signingKey))
data := []byte(ref)
refreshtoken:= base64.StdEncoding.EncodeToString(data)
accesstoken,_:= token.SignedString([]byte(signingKey))
	return 	accesstoken,refreshtoken
}
func(s *Service) ParseToken(accessToken string)(uuid.UUID, uuid.UUID, error){
	token, err:= jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err!= nil{
		return uuid.Nil,uuid.Nil, err
	}

	claims, ok :=token.Claims.(*tokenClaims)
	if !ok{
		return uuid.Nil,uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.SessionID, nil
}

// func generatePasswordHash(password string) string{
// 	hash := sha1.New()
// 	hash.Write([]byte(password))

// 	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("HASH_SALT"))))
// }