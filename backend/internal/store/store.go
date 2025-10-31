package store

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"assignment3/backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserExists signals that a username is already taken.
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidCredentials indicates a failed authentication attempt.
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrInvalidRole is returned when a requested role is unsupported.
	ErrInvalidRole = errors.New("invalid role")
	// ErrItemNotFound indicates that an item could not be located.
	ErrItemNotFound = errors.New("item not found")
	// ErrForbidden is returned when the caller lacks permission to perform an action.
	ErrForbidden = errors.New("forbidden")
	// ErrUserNotFound indicates that a user could not be located.
	ErrUserNotFound = errors.New("user not found")
)

// Store provides a concurrency-safe in-memory data store.
type Store struct {
	mu    sync.RWMutex
	items map[string]models.Item
	users map[string]models.User // keyed by lowercase username
}

// NewStore constructs a new store instance.
func NewStore() *Store {
	return &Store{
		items: make(map[string]models.Item),
		users: make(map[string]models.User),
	}
}

// EnsureAdminUser creates an admin user if it does not exist. If the user already
// exists, it is returned and no error is raised. The boolean indicates whether
// a new user was created.
func (s *Store) EnsureAdminUser(username, password string) (models.User, bool, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return models.User{}, false, fmt.Errorf("admin username cannot be empty")
	}
	if password == "" {
		return models.User{}, false, fmt.Errorf("admin password cannot be empty")
	}

	usernameKey := strings.ToLower(username)

	s.mu.Lock()
	defer s.mu.Unlock()

	if user, ok := s.users[usernameKey]; ok {
		// Guarantee the user retains the admin role.
		if user.Role != "admin" {
			user.Role = "admin"
			s.users[usernameKey] = user
		}
		return user, false, nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, false, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := models.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hashed),
		Role:         "admin",
		CreatedAt:    now,
	}
	s.users[usernameKey] = user

	return user, true, nil
}

// CreateUser registers a new user with the provided role.
func (s *Store) CreateUser(username, password, role string) (models.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return models.User{}, fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return models.User{}, fmt.Errorf("password cannot be empty")
	}

	role = strings.TrimSpace(role)
	if role == "" {
		role = "user"
	}
	if role != "user" && role != "admin" {
		return models.User{}, ErrInvalidRole
	}

	usernameKey := strings.ToLower(username)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[usernameKey]; exists {
		return models.User{}, ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := models.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hashed),
		Role:         role,
		CreatedAt:    now,
	}
	s.users[usernameKey] = user

	return user, nil
}

// Authenticate validates the provided credentials and returns the user on success.
func (s *Store) Authenticate(username, password string) (models.User, error) {
	usernameKey := strings.ToLower(strings.TrimSpace(username))

	s.mu.RLock()
	user, ok := s.users[usernameKey]
	s.mu.RUnlock()
	if !ok {
		return models.User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.User{}, ErrInvalidCredentials
	}

	return user, nil
}

// ListItems returns all items sorted by creation time.
func (s *Store) ListItems() []models.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]models.Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	return items
}

// GetItem returns a single item by id.
func (s *Store) GetItem(id string) (models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return models.Item{}, ErrItemNotFound
	}
	return item, nil
}

// CreateItem inserts a new item owned by the specified user.
func (s *Store) CreateItem(owner, title, description string) (models.Item, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return models.Item{}, fmt.Errorf("title cannot be empty")
	}

	now := time.Now().UTC()
	item := models.Item{
		ID:          uuid.NewString(),
		Title:       title,
		Description: strings.TrimSpace(description),
		Owner:       owner,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.mu.Lock()
	s.items[item.ID] = item
	s.mu.Unlock()

	return item, nil
}

// UpdateItem updates an existing item if the caller is the owner or an admin.
func (s *Store) UpdateItem(id, requester string, isAdmin bool, title, description string) (models.Item, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return models.Item{}, fmt.Errorf("title cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return models.Item{}, ErrItemNotFound
	}

	requester = strings.TrimSpace(requester)
	if item.Owner != requester && !isAdmin {
		return models.Item{}, ErrForbidden
	}

	item.Title = title
	item.Description = strings.TrimSpace(description)
	item.UpdatedAt = time.Now().UTC()
	s.items[id] = item

	return item, nil
}

// DeleteItem removes an item from the store.
func (s *Store) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrItemNotFound
	}

	delete(s.items, id)
	return nil
}

// ListUsers returns all users sorted by creation time.
func (s *Store) ListUsers() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})

	return users
}

// GetUser returns a single user by id.
func (s *Store) GetUser(id string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, ErrUserNotFound
}

// DeleteUser removes a user from the store.
func (s *Store) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, user := range s.users {
		if user.ID == id {
			delete(s.users, key)
			return nil
		}
	}
	return ErrUserNotFound
}
