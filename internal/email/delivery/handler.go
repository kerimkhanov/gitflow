package delivery

import (
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

func NewHandler(usecase *usecase.EmailUseCase) *Handler {
	tmpl, _ := template.ParseGlob("./templates/*.html")
	return &Handler{
		tmpl:    tmpl,
		usecase: usecase,
	}
}

func (h *Handler) InitRoutes() {
	router := gin.Default()
	router.POST("/sending", h.sending)

}
func (h *Handler) sending(c *gin.Context) {
	var input models.UserMail
	if err := c.BindJSON(&input); err != nil {
		log.Println(err)
		return
	}
	if err := h.usecase.DoTasks(input, h.tmpl); err != nil {
		log.Println(err)
		return
	}
}
