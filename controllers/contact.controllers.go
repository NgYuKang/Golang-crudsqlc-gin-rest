package controllers

import (
	"context"
	"database/sql"
	db "golang-crudsqlc-gin-rest/db/sqlc"
	"golang-crudsqlc-gin-rest/schemas"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	db  *db.Queries
	ctx context.Context
}

func NewContactController(db *db.Queries, ctx context.Context) *ContactController {
	return &ContactController{db, ctx}
}

func (cc *ContactController) CreateContact(ctx *gin.Context) {

	var payload *schemas.CreateContact

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed payload",
		})
		return
	}

	timeNow := time.Now()
	args := &db.CreateContactParams{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneNumber: payload.PhoneNumber,
		Street:      payload.Street,
	}
	args.CreatedAt.Scan(timeNow)
	args.UpdatedAt.Scan(timeNow)

	contact, err := cc.db.CreateContact(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed payload",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created contact", "contact": contact})

}

func (cc *ContactController) UpdateContact(ctx *gin.Context) {
	var payload *schemas.UpdateContact
	contactId := ctx.Param("contactId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "failed payload"})
		return
	}

	timeNow := time.Now()
	args := &db.UpdateContactParams{}
	args.ContactID.Scan(contactId)
	args.FirstName.Scan(payload.FirstName)
	args.FirstName.Scan(payload.LastName)
	args.FirstName.Scan(payload.PhoneNumber)
	args.FirstName.Scan(payload.Street)
	args.UpdatedAt.Scan(timeNow)
	log.Println(args)

	contact, err := cc.db.UpdateContact(ctx, *args)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "failed retrieve"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "ok", "contact": contact})

}
