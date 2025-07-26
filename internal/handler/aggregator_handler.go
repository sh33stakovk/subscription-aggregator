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

// @Summary	Создание подписки
// @Accept		json
// @Produce	json
// @Param		subscription	body		swagger.SubscriptionExample	true	"Данные подписки"
// @Success	200				{object}	swagger.MessageResponse
// @Failure	400				{object}	swagger.ErrorResponse400
// @Failure	500				{object}	swagger.ErrorResponse500
// @Router		/create [post]
func CreateSubscription(c *gin.Context) {
	var sub model.Subscription

	log.Println("[CreateSubscription] received request")

	if err := c.ShouldBindJSON(&sub); err != nil {
		log.Printf("[CreateSubscription] JSON bind error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON"})
		return
	}

	log.Printf(
		"[CreateSubscription] creating subscription for user_id=%s, service_name=%s\n",
		sub.UserID,
		sub.ServiceName,
	)

	if err := repository.DB.Create(&sub).Error; err != nil {
		log.Printf("[CreateSubscription] DB create error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record in db"})
		return
	}

	log.Printf("[CreateSubscription] successfully created subscription ID=%d\n", sub.ID)
	c.JSON(http.StatusOK, gin.H{"message": "created", "id": sub.ID})
}

// @Summary	Получить данные подписки по ID
// @Produce	json
// @Param		id	path		int	true	"ID подписки"	default(1)
// @Success	200	{object}	swagger.SubscriptionResponse
// @Failure	400	{object}	swagger.ErrorResponse400
// @Failure	404	{object}	swagger.ErrorResponse404
// @Failure	500	{object}	swagger.ErrorResponse500
// @Router		/read/{id} [get]
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

// @Summary	Обновить подписку по ID
// @Accept		json
// @Produce	json
// @Param		id				path		int									true	"ID подписки"	default(1)
// @Param		subscription	body		swagger.UpdateSubscriptionExample	true	"Новые данные подписки"
// @Success	200				{object}	swagger.MessageResponse
// @Failure	400				{object}	swagger.ErrorResponse400
// @Failure	404				{object}	swagger.ErrorResponse404
// @Failure	500				{object}	swagger.ErrorResponse500
// @Router		/update/{id} [put]
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

// @Summary	Удалить подписку по ID
// @Produce	json
// @Param		id	path		int	true	"ID подписки"	default(1)
// @Success	200	{object}	swagger.MessageResponse
// @Failure	400	{object}	swagger.ErrorResponse400
// @Failure	404	{object}	swagger.ErrorResponse404
// @Failure	500	{object}	swagger.ErrorResponse500
// @Router		/delete/{id} [delete]
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

// @Summary	Получение списка подписок (есть фильтрация по ID пользователя и по названию сервиса)
// @Produce	json
// @Param		user_id			query		string	false	"ID пользователя"	default(11111111-1111-1111-1111-111111111111)
// @Param		service_name	query		string	false	"Название сервиса"	default(Netflix)
// @Success	200				{array}		swagger.SubscriptionResponse
// @Failure	500				{object}	swagger.ErrorResponse500
// @Failure	404				{object}	swagger.ErrorResponse404
// @Router		/list [get]
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

// @Summary	Получение суммы стоимости всех подписок за выбранный период по ID пользователя и имени сервиса
// @Produce	json
// @Param		user_id			query		string	true	"ID пользователя"					default(11111111-1111-1111-1111-111111111111)
// @Param		service_name	query		string	true	"Название сервиса"					default(Netflix)
// @Param		period_start	query		string	true	"Начало периода в формате MM-YYYY"	default(06-2025)
// @Param		period_end		query		string	true	"Конец периода в формате MM-YYYY"	default(08-2025)
// @Success	200				{object}	swagger.SumResponse
// @Failure	400				{object}	swagger.ErrorResponse400
// @Failure	500				{object}	swagger.ErrorResponse500
// @Router		/sum [get]
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
