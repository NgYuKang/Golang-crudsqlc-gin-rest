package routes

import (
	"golang-crudsqlc-gin-rest/controllers"

	"github.com/gin-gonic/gin"
)

type ContactRouter struct {
	contractController controllers.ContactController
}

func NewContactRoute(cc controllers.ContactController) ContactRouter {
	return ContactRouter{cc}
}

func (cr *ContactRouter) ContactRouter(rg *gin.RouterGroup) {

	router := rg.Group("contacts")

	router.POST("/", cr.contractController.CreateContact)
	router.PATCH("/:contactId", cr.contractController.UpdateContact)
	router.GET("/:contactId", cr.contractController.GetContactById)
	router.GET("/", cr.contractController.ListContacts)
	router.DELETE("/:contactId", cr.contractController.DeleteContactById)
}
