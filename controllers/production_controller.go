package controllers

import (
	"context"
	"net/http"
	"time"

	"production-warehouse-api/config"
	"production-warehouse-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"production-warehouse-api/job"
)

func CreateProduction(c *gin.Context) {
	var prod models.Production
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prod.ID = primitive.NewObjectID()
	prod.ProductionDate = time.Now()
	prod.Status = "pending"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.DB.Collection("productions").InsertOne(ctx, prod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data produksi"})
		return
	}

	prodID := prod.ID.Hex()
	jobCtx, jobCancel := context.WithCancel(context.Background())
	job.AddJob(prodID, jobCancel)

	go func() {
	batchSize := prod.QuantityProduced
	seconds := (batchSize / 100) * 10
	if seconds == 0 {
		seconds = 10 
	}

	select {
	case <-time.After(time.Duration(seconds) * time.Second):
		processProduction(prod)
		job.CancelJob(prodID)

	case <-jobCtx.Done():
		fmt.Println("Produksi", prodID, "dibatalkan.")

		config.DB.Collection("productions").UpdateOne(context.Background(), bson.M{"_id": prod.ID}, bson.M{
			"$set": bson.M{"status": "canceled"},
		})

		logProductionAction(prod.ID, "canceled", "Produksi dibatalkan oleh pengguna")
	}
}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Produksi dimulai", "id": prodID})
}


func GetProductions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.DB.Collection("productions").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produksi"})
		return
	}
	defer cursor.Close(ctx)

	var productions []models.Production
	for cursor.Next(ctx) {
		var prod models.Production
		if err := cursor.Decode(&prod); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
			return
		}
		productions = append(productions, prod)
	}

	c.JSON(http.StatusOK, productions)
}

func GetProductionByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var prod models.Production
	err = config.DB.Collection("productions").FindOne(ctx, bson.M{"_id": objID}).Decode(&prod)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data produksi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, prod)
}

func processProduction(prod models.Production) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, mat := range prod.Materials {
		filter := bson.M{"_id": mat.ItemID}
		update := bson.M{"$inc": bson.M{"quantity": -mat.QuantityUsed}}
		config.DB.Collection("items").UpdateOne(ctx, filter, update)
	}

	var existing models.Item
	err := config.DB.Collection("items").FindOne(ctx, bson.M{"name": prod.ProductName}).Decode(&existing)
	if err != nil {
		newItem := models.Item{
			ID:          primitive.NewObjectID(),
			Name:        prod.ProductName,
			Category:    "Produk Jadi",
			Quantity:    prod.QuantityProduced,
			Warehouse:   "Gudang Produksi",
			Description: "Dari proses produksi otomatis",
		}
		config.DB.Collection("items").InsertOne(ctx, newItem)
	} else {
		config.DB.Collection("items").UpdateOne(ctx, bson.M{"_id": existing.ID}, bson.M{
			"$inc": bson.M{"quantity": prod.QuantityProduced},
		})
	}

	config.DB.Collection("productions").UpdateOne(ctx, bson.M{"_id": prod.ID}, bson.M{
		"$set": bson.M{"status": "completed"},
	})

	logProductionAction(prod.ID, "completed", "Produksi selesai otomatis")
}



func CancelProduction(c *gin.Context) {
	id := c.Param("id")
	success := job.CancelJob(id)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Produksi dibatalkan"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produksi tidak ditemukan atau sudah selesai"})
	}
}

func logProductionAction(prodID primitive.ObjectID, action, note string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logEntry := models.ProductionLog{
		ID:           primitive.NewObjectID(),
		ProductionID: prodID,
		Action:       action,
		Timestamp:    time.Now(),
		Note:         note,
	}

	config.DB.Collection("production_logs").InsertOne(ctx, logEntry)
}

