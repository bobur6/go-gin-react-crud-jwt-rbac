package store_test

import (
	"errors"
	"testing"

	"assignment3/backend/internal/store"
)

func TestCreateUserAndAuthenticate(t *testing.T) {
	st := store.NewStore()

	user, err := st.CreateUser("alice", "password123", "user")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if user.Username != "alice" || user.Role != "user" {
		t.Fatalf("unexpected user payload: %+v", user)
	}

	if _, err := st.CreateUser("alice", "anotherpass", "user"); !errors.Is(err, store.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}

	if _, err := st.Authenticate("alice", "wrong"); !errors.Is(err, store.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}

	authUser, err := st.Authenticate("alice", "password123")
	if err != nil {
		t.Fatalf("Authenticate returned error: %v", err)
	}

	if authUser.ID != user.ID {
		t.Fatalf("expected matching IDs, got %s vs %s", authUser.ID, user.ID)
	}
}

func TestEnsureAdminUser(t *testing.T) {
	st := store.NewStore()

	user, created, err := st.EnsureAdminUser("admin", "secret")
	if err != nil {
		t.Fatalf("EnsureAdminUser returned error: %v", err)
	}
	if !created {
		t.Fatalf("expected user to be created on first call")
	}
	if user.Role != "admin" {
		t.Fatalf("expected admin role, got %s", user.Role)
	}

	userAgain, createdAgain, err := st.EnsureAdminUser("admin", "secret")
	if err != nil {
		t.Fatalf("EnsureAdminUser second call returned error: %v", err)
	}
	if createdAgain {
		t.Fatalf("expected created=false on second call")
	}
	if userAgain.ID != user.ID {
		t.Fatalf("expected same user to be returned on subsequent call")
	}
}

func TestCreateUpdateDeleteItem(t *testing.T) {
	st := store.NewStore()
	admin, _, err := st.EnsureAdminUser("admin", "secret")
	if err != nil {
		t.Fatalf("failed to seed admin: %v", err)
	}

	item, err := st.CreateItem(admin.Username, "First", "description")
	if err != nil {
		t.Fatalf("CreateItem returned error: %v", err)
	}

	updated, err := st.UpdateItem(item.ID, admin.Username, true, "Updated", "new desc")
	if err != nil {
		t.Fatalf("UpdateItem as admin returned error: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected updated title, got %s", updated.Title)
	}

	if _, err := st.UpdateItem(item.ID, "bob", false, "oops", ""); !errors.Is(err, store.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	if err := st.DeleteItem(item.ID); err != nil {
		t.Fatalf("DeleteItem returned error: %v", err)
	}

	if err := st.DeleteItem(item.ID); !errors.Is(err, store.ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound after deletion, got %v", err)
	}
}
