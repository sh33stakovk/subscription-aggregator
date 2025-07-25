package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"subscription-aggregator/internal/model"
	"subscription-aggregator/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func isUniqueServiceName(excludedID int, userID uuid.UUID, serviceName string) (bool, error) {
	var sameSubsCount int64
	err := repository.DB.Model(&model.Subscription{}).
		Where("user_id = ?", userID).
		Where("service_name ILIKE ?", "%"+serviceName+"%").
		Where("id <> ?", excludedID).
		Count(&sameSubsCount).
		Error
	if err != nil {
		return false, err
	} else if sameSubsCount != 0 {
		return false, fmt.Errorf("subscription already exists")
	}

	return true, nil
}

func CreateSubscription(c *gin.Context) {
	var sub model.Subscription

	err := c.ShouldBindJSON(&sub)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	isUnique, err := isUniqueServiceName(0, sub.UserID, sub.ServiceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate"})
		return
	} else if !isUnique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription already exists"})
		return
	}

	err = repository.DB.Create(&sub).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record in db"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "created",
		"id":      sub.ID,
	})
}

func ReadSubscription(c *gin.Context) {
	var sub model.Subscription

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = repository.DB.First(&sub, id).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found in db"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find record in db"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input map[string]interface{}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	userID, err := uuid.Parse(input["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	serviceName := input["service_name"].(string)
	isUnique, err := isUniqueServiceName(id, userID, serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate"})
		return
	} else if !isUnique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription already exists"})
		return
	}

	err = repository.DB.Model(&model.Subscription{}).Where("id = ?", id).Updates(input).Error
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found in db"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update record in db"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = repository.DB.Delete(&model.Subscription{}, id).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found in db"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete record from db"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "deleted"})
}

func ListSubscriptions(c *gin.Context) {
	var subs []model.Subscription

	err := repository.DB.Find(&subs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get records from db"})
		return
	} else if len(subs) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "0 records in db"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func SumSubscriptionsPrice(c *gin.Context) {
	var sum int

	sumReq := struct {
		UserID      string `form:"user_id" binding:"required,uuid"`
		ServiceName string `form:"service_name" binding:"required"`
		PeriodStart string `form:"period_start" binding:"required"`
		PeriodEnd   string `form:"period_end" binding:"required"`
	}{}

	err := c.ShouldBindQuery(&sumReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, err := uuid.Parse(sumReq.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	periodStart, err := time.Parse("01-2006", sumReq.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start"})
		return
	}

	periodEnd, err := time.Parse("01-2006", sumReq.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end"})
		return
	}

	err = repository.DB.Model(&model.Subscription{}).
		Select("COALESCE(SUM(price), 0)").
		Where("user_id = ?", userID).
		Where("service_name ILIKE ?", "%"+sumReq.ServiceName+"%").
		Where("start_date BETWEEN ? AND ?", periodStart, periodEnd).
		Scan(&sum).
		Error
	if sum == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "records not found in db"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sum"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sum_price": sum})
}
