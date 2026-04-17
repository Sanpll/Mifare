package handler

import (
	"mifare/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "mifare/docs"
)

type Handler struct {
	apiVersion string
	services   *service.Service
}

func NewHandler(apiVersion string, services *service.Service) *Handler {
	return &Handler{
		apiVersion: apiVersion,
		services:   services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/api/"+h.apiVersion+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	api := router.Group("/api/"+h.apiVersion, h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.adminOnly, h.GetUsers)
			users.GET("/:id", h.isSameUserOrAdmin, h.GetUserById)
			users.PUT("/:id", h.isSameUserOrAdmin, h.UpdateUser)
			users.DELETE("/:id", h.isSameUserOrAdmin, h.DeleteUser)
		}

		keys := api.Group("/keys", h.adminOnly)
		{
			keys.POST("/", h.CreateKey)
			keys.GET("/", h.GetKeys)
			keys.GET("/:id", h.GetKeyById)
			keys.PUT("/:id", h.UpdateKey)
			keys.DELETE("/:id", h.DeleteKey)
		}

		cards := api.Group("/cards")
		{
			cards.POST("/", h.CreateCard)
			cards.GET("/", h.GetCards)
			cards.GET("/:id", h.GetCardById)
			cards.PUT("/:id", h.UpdateCard)
			cards.DELETE("/:id", h.DeleteCard)
		}

		terminals := api.Group("/terminals")
		{
			terminals.POST("/", h.adminOnly, h.CreateTerminal)
			terminals.GET("/", h.GetTerminals)
			terminals.GET("/:id", h.GetTerminalById)
			terminals.PUT("/:id", h.adminOnly, h.UpdateTerminal)
			terminals.DELETE("/:id", h.adminOnly, h.DeleteTerminal)
		}

		transactions := api.Group("/transactions", h.adminOnly)
		{
			transactions.POST("/", h.CreateTransaction)
			transactions.GET("/", h.GetTransactions)
			transactions.GET("/:id", h.GetTransactionById)
			transactions.PUT("/:id", h.UpdateTransaction)
			transactions.DELETE("/:id", h.DeleteTransaction)
		}
	}

	terminal := router.Group("/api/" + h.apiVersion + "/terminal")
	{
		terminal.POST("/auth", h.AuthorizeTransaction)
		terminal.GET("/keys", h.GetAllKeys)
	}

	return router
}
