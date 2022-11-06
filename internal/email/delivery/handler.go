package delivery

import (
	"fmt"
	"gitflow/internal/email/models"
	"gitflow/internal/email/usecase"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

type Handler struct {
	usecase *usecase.EmailUseCase
	tmpl    *template.Template
	config  *template.Template
}

func newHandler(usecase *usecase.EmailUseCase) *Handler {
	tmpl, _ := template.ParseGlob("./templates/*.html")
	return &Handler{
		tmpl:    tmpl,
		usecase: usecase,
	}
}

func InitEmailRoutes(router *gin.Engine) {
	usecase := &usecase.EmailUseCase{}
	h := newHandler(usecase)
	router.POST("/sending", h.sending)
}

func (h *Handler) sending(c *gin.Context) {
	fmt.Println("here")
	var input models.UserMail
	if err := c.BindJSON(&input); err != nil {
		log.Println(err)
		return
	}
	fmt.Println("here2")
	fmt.Println(input.Emails)
	if err := h.usecase.DoTasks(input, h.tmpl); err != nil {
		log.Println(err)
		return
	}
}
