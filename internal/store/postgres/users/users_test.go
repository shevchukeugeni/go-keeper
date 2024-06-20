package users

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"keeper-project/types"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &types.User{
		ID:        "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
		Login:     "test",
		Password:  "some_text",
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("^INSERT INTO users(.+)").WithArgs(user.ID, user.Login, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.CreateUser(ctx, user)
	require.NoError(t, err)
}

func TestCreate_NilUserInfo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewRepository(db)

	ctx := context.Background()

	err = store.CreateUser(ctx, nil)
	require.Equal(t, err.Error(), "repository: incorrect parameters")
}

func TestCreate_DuplicateErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &types.User{
		ID:        "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
		Login:     "test",
		Password:  "some_text",
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("^INSERT INTO users(.+)").WithArgs(user.ID, user.Login, user.Password).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.CreateUser(ctx, user)
	require.Equal(t, err, types.ErrUserAlreadyExists)
}

func TestCreate_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &types.User{
		ID:        "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
		Login:     "test",
		Password:  "some_text",
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("^INSERT INTO users(.+)").WithArgs(user.ID, user.Login, user.Password).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.CreateUser(ctx, user)
	require.Equal(t, err, sql.ErrConnDone)
}

func TestGet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &types.User{
		ID:        "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
		Login:     "test",
		Password:  "some_text",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("^SELECT id, password, created_at FROM users WHERE login(.+)").WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "created_at"}).
			AddRow("40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", user.Password, user.CreatedAt))

	store := NewRepository(db)

	ctx := context.Background()

	userFromDB, err := store.GetByLogin(ctx, user.Login)
	require.NoError(t, err)

	require.Equal(t, userFromDB.Login, user.Login)
	require.Equal(t, userFromDB.Password, user.Password)
}

func TestGet_FailLogin(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewRepository(db)

	ctx := context.Background()

	_, err = store.GetByLogin(ctx, "")
	require.Equal(t, err.Error(), "repository: incorrect parameters")
}

func TestGet_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &types.User{
		ID:        "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
		Login:     "test",
		Password:  "some_text",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("^SELECT id, password, created_at FROM users WHERE login(.+)").WithArgs(user.Login).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	_, err = store.GetByLogin(ctx, user.Login)
	require.Equal(t, err, sql.ErrConnDone)
}
