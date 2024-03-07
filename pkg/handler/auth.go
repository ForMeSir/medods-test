package handler

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"ruby/pkg/user"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const(
	ACtokenTimeLive = 3 * time.Hour
	REFtokenTimeLive = 22 * time.Hour
 )


	type signInInput struct{
	ID uuid.UUID `json:"id" binding:"required"`
	}
	type refres struct {
		Refreshtoken string `json:"refresh_token" binding:"required"`
	}
	
func (h *Handler) refresh(c *gin.Context){
	var acsstoken string
	var tokens refres 
	if err:= c.ShouldBindJSON(&tokens); err!=nil{
	 newErrorResponse(c,http.StatusBadRequest, err.Error())
	 return 
 }
 logrus.Info( tokens.Refreshtoken)
	 
	 data, err := base64.StdEncoding.DecodeString(tokens.Refreshtoken)
	 if err != nil {
			 fmt.Println("error:", err)
			 return
	 }
	 logrus.Info(data)

	  k:=sha512.Sum512(data)
	  hashedrefreshtoken:= string(k[:])
	 
 
	 tokens.Refreshtoken=string(data)
	 
	 logrus.Info(tokens.Refreshtoken)

	  _, refsessionId, err:= h.services.Autorization.ParseToken(string(data))

	if err!=nil{
	newErrorResponse(c,http.StatusBadRequest, err.Error())
  return
 }

 logrus.Info(refsessionId)

		findSession, err :=h.mongo.FindOne(context.Background(), refsessionId)

		if err!=nil{
			newErrorResponse(c,http.StatusBadRequest, err.Error())
	    return 
		}

		if findSession.IsBlocked{
      newErrorResponse(c,http.StatusBadRequest, "token is blocked")
	    return 
		}

		if bcrypt.CompareHashAndPassword([]byte(findSession.RefreshToken),[]byte(hashedrefreshtoken)) != nil{
			newErrorResponse(c,http.StatusBadRequest, "token is changed")
			return 
		}
		if findSession.ExpiresAt.After(time.Now()){
			acsstoken = h.services.Autorization.RefreshToken(findSession.UserId, refsessionId, ACtokenTimeLive)
		}else {
		newErrorResponse(c,http.StatusBadRequest, "token is invalid")
	  return
	}

c.JSON(http.StatusOK, map[string]interface{}{
	"AccessToken":acsstoken,
  "AccesstExpiresAt": time.Now().Add(ACtokenTimeLive) ,
})
}


func (h *Handler) signIn(c *gin.Context){
	//var input signInInput
	var id signInInput
	if err:= c.BindJSON(&id); err!=nil{
	 newErrorResponse(c,http.StatusBadRequest, err.Error())
	 return 
 }

	accesstoken,refreshtoken := h.services.Autorization.GenerateToken(id.ID, ACtokenTimeLive, REFtokenTimeLive)
	

	  userid,sessionId, err:= h.services.Autorization.ParseToken(accesstoken)
	  if err != nil { 
	  	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	   }
		 logrus.Info(userid,sessionId)
		 decodedRefresh, err := base64.StdEncoding.DecodeString(refreshtoken)
		 if err != nil {
				 fmt.Println("error:", err)
				 return
		 }
      
			k:=sha512.Sum512(decodedRefresh)
			hashedrefreshtoken:= string(k[:])
			refreshtokenhash,err := bcrypt.GenerateFromPassword([]byte(hashedrefreshtoken),10)
			if err != nil { 
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
			 }
			logrus.Info(refreshtokenhash)

			 ker:= user.CreateSessionParams{
				ID: "",
				SessionID: sessionId,
				UserId: userid,
				RefreshToken: refreshtokenhash,
				IsBlocked: false,
				ExpiresAt: time.Now().Add(REFtokenTimeLive),
			 }
			 
			a, err:= h.mongo.Create(context.Background(), ker)
			if err != nil{
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				logrus.Infof("failed create session %s",a)
			}
 
	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":accesstoken,
   "accesst_expires_at": time.Now().Add(ACtokenTimeLive) ,
	 "refresh_token": refreshtoken,
	 "refresh_expires_at":  ker.ExpiresAt,
 })
}
