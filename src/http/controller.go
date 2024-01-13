package http

import (
	"github.com/gin-gonic/gin"
	"krisha/src/internal/service/parallel"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	router *gin.Engine
}

func NewController() *Controller {
	router := gin.Default()
	controller := &Controller{router}

	router.GET("/requests", controller.handleEndpoint1)

	return controller
}

func (c *Controller) handleEndpoint1(context *gin.Context) {
	jobs := make([]func() string, 0)
	for i := 0; i < 10; i++ {
		jobs = append(jobs, func() string {
			i := doRequest(1)
			return strconv.Itoa(i)
		})
	}
	result := parallel.DoJobs(jobs, 4)
	context.String(http.StatusOK, "Requests completed: "+strings.Join(result, ","))
}

func (c *Controller) Start() {
	c.router.Run(":8083")
}

func doRequest(d int) int {
	time.Sleep(time.Duration(d) * time.Second)
	return rand.Intn(20)
}
