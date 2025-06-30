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
)

func CreateProduction(c *gin.Context) {
	var prod models.Production

	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prod.ID = primitive.NewObjectID()
	prod.ProductionDate = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, mat := range prod.Materials {
		filter := bson.M{"_id": mat.ItemID}
		update := bson.M{"$inc": bson.M{"quantity": -mat.QuantityUsed}}

		result := config.DB.Collection("items").FindOneAndUpdate(ctx, filter, update)
		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengurangi stok material"})
			return
		}
	}

	var productItem models.Item
	err := config.DB.Collection("items").FindOne(ctx, bson.M{"name": prod.ProductName}).Decode(&productItem)

	if err != nil {
		newItem := models.Item{
			ID:          primitive.NewObjectID(),
			Name:        prod.ProductName,
			Category:    "Hasil Produksi",
			Quantity:    prod.QuantityProduced,
			Warehouse:   "Gudang Produksi",
			Description: "Otomatis dari proses produksi",
		}
		_, err := config.DB.Collection("items").InsertOne(ctx, newItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat produk item"})
			return
		}
	} else {
		_, err := config.DB.Collection("items").UpdateOne(ctx, bson.M{"_id": productItem.ID}, bson.M{
			"$inc": bson.M{"quantity": prod.QuantityProduced},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update stok produk jadi"})
			return
		}
	}

	_, err = config.DB.Collection("productions").InsertOne(ctx, prod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data produksi"})
		return
	}

	c.JSON(http.StatusCreated, prod)
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

