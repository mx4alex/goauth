package server

import (
	"goauth/internal/entity"
	"goauth/internal/usecase"
	"net/http"

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

	tasks := router.Group("/auth")
	{
		tasks.POST("/sign-up", h.SignUp)
		tasks.POST("/sign-in", h.SignIn)
	}

	return router
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(entity.User)
	if err := c.BindJSON(inp); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	err := h.authInteractor.SignUp(c.Request.Context(), inp)
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(entity.User)
	if err := c.BindJSON(inp); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	token, err := h.authInteractor.SignIn(c.Request.Context(), inp)
	if err!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
}