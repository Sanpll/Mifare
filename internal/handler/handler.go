package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth") 
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/", h.CreateUser)
			users.GET("/", h.GetUsers)
			users.GET("/:id", h.GetUserById)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}

		cards := api.Group("/cards")
		{
			cards.POST("/", h.CreateCard)
			cards.GET("/", h.GetCards)
			cards.GET("/:id", h.GetCardById)
			cards.PUT("/:id", h.UpdateCard)
			cards.DELETE("/:id", h.DeleteCard)
		}

		keys := api.Group("/keys")
		{
			keys.POST("/", h.CreateKey)
			keys.GET("/", h.GetKeys)
			keys.GET("/:id", h.GetKeyById)
			keys.PUT("/:id", h.UpdateKey)
			keys.DELETE("/:id", h.DeleteKey)
		}

		terminals := api.Group("/terminals")
		{
			terminals.POST("/", h.CreateTerminal)
			terminals.GET("/", h.GetTerminals)
			terminals.GET("/:id", h.GetTerminalById)
			terminals.PUT("/:id", h.UpdateTerminal)
			terminals.DELETE("/:id", h.DeleteTerminal)
		}

		transactions := api.Group("/transactions")
		{
			transactions.POST("/", h.CreateTransaction)
			transactions.GET("/", h.GetTransactions)
			transactions.GET("/:id", h.GetTransactionById)
			transactions.PUT("/:id", h.UpdateTransaction)
			transactions.DELETE("/:id", h.DeleteTransaction)
		}
	}

	terminal := router.Group("/api/v1/terminal")
	{
		terminal.POST("/authorize", h.AuthorizeTransaction)
		terminal.GET("/keys", h.GetAllKeys)
	}

	return router
}