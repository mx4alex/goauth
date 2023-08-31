package server

import (
	"goauth/internal/entity"
	"goauth/internal/usecase"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"

	_ "goauth/docs"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{	
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/refresh", h.Refresh)
	}

	return router
}

// @Summary 	SignUp
// @Tags 		auth
// @Description create account
// @ID 			create-account
// @Accept  	json
// @Produce  	json
// @Param 		input body entity.UserSignUp true "account info"
// @Success 	200 {object} outputTokens
// @Failure 	400,404 {object} errorMessage
// @Failure 	default {object} errorMessage
// @Router /auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	user := new(entity.UserSignUp)

	if err := c.BindJSON(user); err != nil {
		log.Println(err)
        c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
        return
    }

	tokens, err := h.authInteractor.SignUp(c.Request.Context(), user)
	if err != nil {
		log.Println(err)
        c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
        return
    }

	c.JSON(http.StatusOK, successResponse(tokens))
}

// @Summary 	SignIn
// @Tags 		auth
// @Description login
// @ID 			login
// @Accept  	json
// @Produce  	json
// @Param 		input body entity.UserSignIn true "credentials"
// @Success 	200 {object} outputTokens
// @Failure 	400,404 {object} errorMessage
// @Failure 	default {object} errorMessage
// @Router /auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	user := new(entity.UserSignIn)
	if err := c.BindJSON(user); err != nil {
        c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
        return
    }

	tokens, err := h.authInteractor.SignIn(c.Request.Context(), user)
	if err!= nil {
        c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
        return
    }

	c.JSON(http.StatusOK, successResponse(tokens))
}

// @Summary 	RefreshToken
// @Tags 		auth
// @Description refresh token
// @ID 			refresh-token
// @Accept  	json
// @Produce  	json
// @Param 		input body refreshRequest true "refresh token"
// @Success 	200 {object} outputTokens
// @Failure 	400,404 {object} errorMessage
// @Failure 	default {object} errorMessage
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var requestBody refreshRequest
	if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
        return
    }

	tokens, err := h.authInteractor.RefreshToken(c.Request.Context(), requestBody.RefreshToken)
	if err!= nil {
        c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
        return
    }

	c.JSON(http.StatusOK, successResponse(tokens))
}