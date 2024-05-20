package controllers

import (
	"context"
	"database/sql"
	db "golang-crudsqlc-gin-rest/db/sqlc"
	"golang-crudsqlc-gin-rest/schemas"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
	// Assuming it passes the binding above, i am not checking error here.
	args.ContactID.Scan(contactId)
	args.FirstName.Scan(payload.FirstName)
	args.FirstName.Scan(payload.LastName)
	args.FirstName.Scan(payload.PhoneNumber)
	args.FirstName.Scan(payload.Street)
	args.UpdatedAt.Scan(timeNow)
	log.Println(args)

	contact, err := cc.db.UpdateContact(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "failed retrieve"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "ok", "contact": contact})

}

func (cc *ContactController) GetContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	parsedId := pgtype.UUID{}
	err := parsedId.Scan(contactId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "failed payload"})
		return
	}
	contact, err := cc.db.GetContactById(ctx, parsedId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "failed retrieve"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "ok", "contact": contact})
}

func (cc *ContactController) ListContacts(ctx *gin.Context) {
	// Using offset is not good! Future projects I will use last unique value sorting, probably....
	// But it seems unlikely/hard to do with generated sqlc go code
	// (OR LIKE USE GENERATED SQLC CODE FOR EVERYTHING, EXCEPT LIST)
	// my idea is to just copy how they do the code, manually writing it (writing the service)
	// instead of relying on generated sqlc code. but use sqlc to handle migrations.
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.ListContactsParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	contacts, err := cc.db.ListContacts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed to retrieve list of contacts"})
		return
	}

	if contacts == nil {
		contacts = []db.Contact{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"message":  "Successfully retrieved all contacts",
		"contacts": contacts,
	})
}

func (cc *ContactController) DeleteContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	parsedId := pgtype.UUID{}
	err := parsedId.Scan(contactId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "failed payload"})
		return
	}
	deletedContact, err := cc.db.GetContactById(ctx, parsedId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed to delete"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully deleted",
		"contact": deletedContact,
	})
}
