package server

import (
	"goauth/internal/entity"
	"goauth/internal/usecase"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authInteractor *usecase.AuthInteractor
}

func NewHandler(authInteractor *usecase.AuthInteractor) *Handler {
	return &Handler{
		authInteractor: authInteractor,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{	
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/refresh", h.Refresh)
	}

	return router
}

func (h *Handler) SignUp(c *gin.Context) {
	user := new(entity.UserInput)
	if err := c.BindJSON(user); err != nil {
		log.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	accessToken, refreshToken, err := h.authInteractor.SignUp(c.Request.Context(), user)
	if err != nil {
		log.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func (h *Handler) SignIn(c *gin.Context) {
	user := new(entity.UserInput)
	if err := c.BindJSON(user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	accessToken, refreshToken, err := h.authInteractor.SignIn(c.Request.Context(), user)
	if err!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Refresh(c *gin.Context) {
	var requestBody refreshRequest
	if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	accessToken, refreshToken, err := h.authInteractor.RefreshToken(c.Request.Context(), requestBody.RefreshToken)
	if err!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}