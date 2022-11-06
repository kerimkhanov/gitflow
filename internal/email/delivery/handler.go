package delivery

import (
	"encoding/json"
	"fmt"
	"gitflow/internal/email/models"
	"gitflow/internal/email/usecase"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	var input models.UserMail
	if err := c.BindJSON(&input); err != nil {
		log.Println(err)
		return
	}
	go h.usecase.StartWorker()
	time.Sleep(100 * time.Millisecond)
	res, err := h.usecase.StartClient(input, h.tmpl)
	if err != nil {
		log.Println(err)
		return
	}
	daa, err := json.Marshal(res)
	if err != nil {
		return
	}
	var v interface{}
	if err = json.Unmarshal(daa, &v); err != nil {
		return
	}
	fmt.Println(string(daa), v)

	c.JSON(http.StatusOK, res)

}
