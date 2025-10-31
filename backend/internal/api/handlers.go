package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"assignment3/backend/internal/auth"
	"assignment3/backend/internal/models"
	"assignment3/backend/internal/store"

	"github.com/gin-gonic/gin"
)

// Handler bundles dependencies required by HTTP handlers.
type Handler struct {
	store *store.Store
	jwt   *auth.JWTService
}

// NewHandler creates a handler instance.
func NewHandler(store *store.Store, jwt *auth.JWTService) *Handler {
	return &Handler{
		store: store,
		jwt:   jwt,
	}
}

type userResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func newUserResponse(user models.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

// Health returns a basic status payload.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register creates a new user account with role `user`.
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be at least 3 characters"})
		return
	}
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
		return
	}

	user, err := h.store.CreateUser(req.Username, req.Password, "user")
	if err != nil {
		status := http.StatusInternalServerError
		message := "failed to create user"

		switch {
		case errors.Is(err, store.ErrUserExists):
			status = http.StatusConflict
			message = "username already taken"
		case errors.Is(err, store.ErrInvalidRole):
			status = http.StatusBadRequest
			message = "invalid role"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, newUserResponse(user))
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and returns a JWT.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := h.store.Authenticate(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, store.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	token, err := h.jwt.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token: token,
		User:  newUserResponse(user),
	})
}

// ListItems returns all items.
func (h *Handler) ListItems(c *gin.Context) {
	items := h.store.ListItems()
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetItem returns a single item by ID.
func (h *Handler) GetItem(c *gin.Context) {
	item, err := h.store.GetItem(c.Param("id"))
	if err != nil {
		if errors.Is(err, store.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

type itemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// CreateItem inserts a new item belonging to the authenticated user.
func (h *Handler) CreateItem(c *gin.Context) {
	var req itemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, ok := auth.GetContextUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	item, err := h.store.CreateItem(user.Username, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateItem updates an item owned by the authenticated user or an admin.
func (h *Handler) UpdateItem(c *gin.Context) {
	var req itemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, ok := auth.GetContextUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	item, err := h.store.UpdateItem(c.Param("id"), user.Username, strings.EqualFold(user.Role, "admin"), req.Title, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrItemNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		case errors.Is(err, store.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to update this item"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem removes an item; route-level middleware ensures the caller is admin.
func (h *Handler) DeleteItem(c *gin.Context) {
	if err := h.store.DeleteItem(c.Param("id")); err != nil {
		if errors.Is(err, store.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers returns all users; route-level middleware ensures the caller is admin.
func (h *Handler) ListUsers(c *gin.Context) {
	users := h.store.ListUsers()
	response := make([]userResponse, len(users))
	for i, user := range users {
		response[i] = newUserResponse(user)
	}
	c.JSON(http.StatusOK, gin.H{"users": response})
}

// DeleteUser removes a user; route-level middleware ensures the caller is admin.
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	// Prevent admin from deleting themselves
	currentUser, ok := auth.GetContextUser(c)
	if ok && currentUser.ID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete your own account"})
		return
	}

	if err := h.store.DeleteUser(userID); err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
