package handler

import (
	"log"
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
		Where("LOWER(service_name) = LOWER(?)", serviceName).
		Where("id <> ?", excludedID).
		Count(&sameSubsCount).
		Error
	if err != nil {
		return false, err
	} else if sameSubsCount != 0 {
		return false, nil
	}

	return true, nil
}

func CreateSubscription(c *gin.Context) {
	var sub model.Subscription

	log.Println("[CreateSubscription] received request")

	if err := c.ShouldBindJSON(&sub); err != nil {
		log.Printf("[CreateSubscription] JSON bind error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	log.Printf(
		"[CreateSubscription] creating subscription for user_id=%s, service_name=%s\n",
		sub.UserID,
		sub.ServiceName,
	)

	isUnique, err := isUniqueServiceName(0, sub.UserID, sub.ServiceName)
	if err != nil {
		log.Printf("[CreateSubscription] uniqueness check failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate"})
		return
	}
	if !isUnique {
		log.Printf(
			"[CreateSubscription] subscription already exists for user_id=%s and service_name=%s\n",
			sub.UserID,
			sub.ServiceName,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription already exists"})
		return
	}

	if err := repository.DB.Create(&sub).Error; err != nil {
		log.Printf("[CreateSubscription] DB create error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record in db"})
		return
	}

	log.Printf("[CreateSubscription] successfully created subscription ID=%d\n", sub.ID)
	c.JSON(http.StatusOK, gin.H{"message": "created", "id": sub.ID})
}

func ReadSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[ReadSubscription] invalid id param: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var sub model.Subscription
	log.Printf("[ReadSubscription] reading subscription id=%d\n", id)

	err = repository.DB.First(&sub, id).Error
	if err == gorm.ErrRecordNotFound {
		log.Printf("[ReadSubscription] record not found id=%d\n", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found in db"})
		return
	}
	if err != nil {
		log.Printf("[ReadSubscription] DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find record in db"})
		return
	}

	log.Printf("[ReadSubscription] found subscription id=%d\n", sub.ID)
	c.JSON(http.StatusOK, sub)
}

func UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[UpdateSubscription] invalid id param: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[UpdateSubscription] JSON bind error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	if userIDRaw, ok := input["user_id"]; ok {
		if serviceNameRaw, ok2 := input["service_name"]; ok2 {
			userIDStr, ok := userIDRaw.(string)
			if !ok {
				log.Println("[UpdateSubscription] user_id is not string")
				c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be string"})
				return
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				log.Printf("[UpdateSubscription] invalid user_id: %v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
				return
			}

			serviceName, ok := serviceNameRaw.(string)
			if !ok {
				log.Println("[UpdateSubscription] service_name is not string")
				c.JSON(http.StatusBadRequest, gin.H{"error": "service_name must be string"})
				return
			}

			isUnique, err := isUniqueServiceName(id, userID, serviceName)
			if err != nil {
				log.Printf("[UpdateSubscription] uniqueness check error: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate"})
				return
			}
			if !isUnique {
				log.Printf(
					"[UpdateSubscription] subscription already exists for user_id=%s and service_name=%s\n",
					userID,
					serviceName,
				)
				c.JSON(http.StatusBadRequest, gin.H{"error": "subscription already exists"})
				return
			}
		}
	}

	log.Printf("[UpdateSubscription] updating subscription id=%d with data: %+v\n", id, input)

	result := repository.DB.Model(&model.Subscription{}).Where("id = ?", id).Updates(input)
	if result.RowsAffected == 0 {
		log.Printf("[UpdateSubscription] no record found to update for id=%d\n", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found in db"})
		return
	}
	if result.Error != nil {
		log.Printf("[UpdateSubscription] DB update error: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update record in db"})
		return
	}

	log.Printf("[UpdateSubscription] successfully updated subscription id=%d\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[DeleteSubscription] invalid id param: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	log.Printf("[DeleteSubscription] deleting subscription id=%d\n", id)

	result := repository.DB.Delete(&model.Subscription{}, id)
	if result.RowsAffected == 0 {
		log.Printf("[DeleteSubscription] record not found id=%d\n", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found in db"})
		return
	}
	if result.Error != nil {
		log.Printf("[DeleteSubscription] DB delete error: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete record from db"})
		return
	}

	log.Printf("[DeleteSubscription] successfully deleted id=%d\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	log.Printf(
		"[ListSubscriptions] fetching subscriptions for user_id=%s, service_name=%s\n",
		userID,
		serviceName,
	)

	query := repository.DB
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	var subs []model.Subscription
	if err := query.Find(&subs).Error; err != nil {
		log.Printf("[ListSubscriptions] DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get records from db"})
		return
	}
	if len(subs) == 0 {
		log.Printf("[ListSubscriptions] no records found\n")
		c.Status(http.StatusNotFound)
		return
	}

	log.Printf("[ListSubscriptions] found %d subscriptions\n", len(subs))
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

	if err := c.ShouldBindQuery(&sumReq); err != nil {
		log.Printf("[SumSubscriptionsPrice] bind query error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	log.Printf(
		"[SumSubscriptionsPrice] calculating sum for user_id=%s, service_name=%s, start=%s, end=%s\n",
		sumReq.UserID,
		sumReq.ServiceName,
		sumReq.PeriodStart,
		sumReq.PeriodEnd,
	)

	userID, err := uuid.Parse(sumReq.UserID)
	if err != nil {
		log.Printf("[SumSubscriptionsPrice] invalid user_id: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	periodStart, err := time.Parse("01-2006", sumReq.PeriodStart)
	if err != nil {
		log.Printf("[SumSubscriptionsPrice] invalid period_start: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start"})
		return
	}

	periodEnd, err := time.Parse("01-2006", sumReq.PeriodEnd)
	if err != nil {
		log.Printf("[SumSubscriptionsPrice] invalid period_end: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end"})
		return
	}

	err = repository.DB.Model(&model.Subscription{}).
		Select("COALESCE(SUM(price), 0)").
		Where("user_id = ?", userID).
		Where("LOWER(service_name) = LOWER(?)", sumReq.ServiceName).
		Where("start_date BETWEEN ? AND ?", periodStart, periodEnd).
		Scan(&sum).
		Error
	if err != nil {
		log.Printf("[SumSubscriptionsPrice] DB error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sum"})
		return
	}

	log.Printf("[SumSubscriptionsPrice] total sum: %d\n", sum)
	c.JSON(http.StatusOK, gin.H{"sum_price": sum})
}
