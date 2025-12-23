package todoctrl

import (
	"errors"
	"fmt"
	"github.com/alirezamastery/graph_task/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

type TodoListResponse struct {
	Count     int64             `json:"count" example:"42"`
	Page      int               `json:"page" example:"1"`
	PageSize  int               `json:"page_size" example:"20"`
	PageCount int               `json:"page_count" example:"3"`
	Items     []models.TodoItem `json:"items"`
}

// GetTodoItemByID godoc
// @Summary Get a todo
// @Description Get a todo item by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.TodoItem
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [get]
func (controller *Controller) GetTodoItemByID() gin.HandlerFunc {
	type Response struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		IsDone      bool   `json:"is_done"`
	}

	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var item models.TodoItem

		if err := controller.db.First(&item, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := &Response{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			IsDone:      item.IsDone,
		}

		c.JSON(http.StatusOK, res)
	}
}

// CreateTodo godoc
// @Summary Create todo item
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param request body todoctrl.CreateTodo.Payload true "Todo payload"
// @Success 201 {object} models.TodoItem
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos [post]
func (controller *Controller) CreateTodo() gin.HandlerFunc {
	type Payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsDone      bool   `json:"is_done"`
	}

	validate := func(c *gin.Context) (*Payload, error) {
		p := &Payload{}
		if err := c.ShouldBindJSON(p); err != nil {
			return nil, err
		}

		p.Title = strings.TrimSpace(p.Title)
		if p.Title == "" {
			return nil, errors.New("\"title\" cannot be empty")
		}
		p.Description = strings.TrimSpace(p.Description)

		return p, nil
	}

	return func(c *gin.Context) {
		payload, err := validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("payload: %+v\n", payload)

		item := models.TodoItem{
			Title:       payload.Title,
			Description: payload.Description,
			IsDone:      payload.IsDone,
		}
		if err := controller.db.Create(&item).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

// GetTodoItemList godoc
// @Summary List todos
// @Description List todos with optional done filter and pagination
// @Tags todos
// @Produce json
// @Param page query int false "page number" default(1)
// @Param page_size query int false "page size" default(20)
// @Param done query bool false "Filter by is_done"
// @Success 200 {object} TodoListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos [get]
func (controller *Controller) GetTodoItemList() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.Query("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}

		pageSizeStr := c.Query("page_size")
		if pageSizeStr == "" {
			pageSizeStr = "20"
		}
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			pageSize = 20
		}

		if page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "\"page\" must be at least 1"})
			return
		}
		if pageSize < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "\"page_size\" must be at least 1"})
			return
		}
		if pageSize > 100 {
			pageSize = 100
		}

		query := controller.db.Model(&models.TodoItem{})

		if doneStr := c.Query("is_done"); doneStr != "" {
			done, err := strconv.ParseBool(doneStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid \"is_done\" query param"})
				return
			}
			query = query.Where("is_done = ?", done)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		offset := (page - 1) * pageSize

		var items []models.TodoItem
		if err := query.
			Find(&items).
			Order("created_at desc").
			Limit(pageSize).
			Offset(offset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

		c.JSON(http.StatusOK, gin.H{
			"page":       page,
			"page_size":  pageSize,
			"count":      total,
			"page_count": totalPages,
			"items":      items,
		})
	}
}

// UpdateTodoItem godoc
// @Summary Update a todo
// @Description Update a todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param payload body todoctrl.UpdateTodoItem.Payload true "Fields to update"
// @Success 200 {object} todoctrl.UpdateTodoItem.Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [patch]
func (controller *Controller) UpdateTodoItem() gin.HandlerFunc {
	type Payload struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		IsDone      *bool   `json:"is_done"`
	}

	type Response struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		IsDone      bool   `json:"is_done"`
	}

	validate := func(c *gin.Context) (*Payload, error) {
		p := &Payload{}
		if err := c.ShouldBindJSON(p); err != nil {
			return nil, err
		}

		if p.Title != nil {
			t := strings.TrimSpace(*p.Title)
			if t == "" {
				return nil, errors.New("\"title\" cannot be empty")
			}
		}
		if p.Description != nil {
			*p.Description = strings.TrimSpace(*p.Description)
		}

		return p, nil
	}

	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		payload, err := validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var item models.TodoItem
		if err := controller.db.First(&item, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updates := map[string]any{}
		if payload.Title != nil {
			updates["title"] = *payload.Title
		}
		if payload.Description != nil {
			updates["description"] = *payload.Description
		}
		if payload.IsDone != nil {
			updates["is_done"] = *payload.IsDone
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		if err := controller.db.Model(&item).Updates(updates).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		_ = controller.db.First(&item, id).Error

		res := &Response{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			IsDone:      item.IsDone,
		}
		c.JSON(http.StatusOK, res)
	}
}

// DeleteTodoItem godoc
// @Summary Delete a todo
// @Description Delete a todo item by ID
// @Tags todos
// @Produce json
// @Param id  path int true "Todo ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [delete]
func (controller *Controller) DeleteTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		res := controller.db.Delete(&models.TodoItem{}, id)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
		if res.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
