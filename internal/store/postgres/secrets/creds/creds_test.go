package creds

import (
	"context"
	"database/sql"
	"errors"
	"testing"

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

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	credentials := &types.Credentials{
		Site:     "test.com",
		Login:    "123321",
		Password: "1224",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^INSERT INTO credentials(.+)").WithArgs(userID, id, credentials.Site,
		credentials.Login, credentials.Password, credentials.Metadata).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, credentials)
	require.NoError(t, err)
}

func TestCreate_NilCredentialInfo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, nil)
	require.Equal(t, err.Error(), "repository: incorrect parameters")
}

func TestCreate_DuplicateErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	credentials := &types.Credentials{
		Site:     "test.com",
		Login:    "123321",
		Password: "1224",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^INSERT INTO credentials(.+)").WithArgs(userID, id, credentials.Site,
		credentials.Login, credentials.Password, credentials.Metadata).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, credentials)
	require.Equal(t, err, types.ErrRecordAlreadyExists)
}

func TestCreate_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	credentials := &types.Credentials{
		Site:     "test.com",
		Login:    "123321",
		Password: "1224",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^INSERT INTO credentials(.+)").WithArgs(userID, id, credentials.Site,
		credentials.Login, credentials.Password, credentials.Metadata).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, credentials)
	require.Equal(t, err, sql.ErrConnDone)
}

func TestUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	credentials := &types.Credentials{
		Site:     "test.com",
		Login:    "123321",
		Password: "1224",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^UPDATE credentials SET(.+)").WithArgs(credentials.Site,
		credentials.Login, credentials.Password, credentials.Metadata, userID, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, credentials)
	require.NoError(t, err)
}

func TestUpdate_NilCredentialInfo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, nil)
	require.Equal(t, err.Error(), "repository: incorrect parameters")
}

func TestUpdate_NotFoundErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	credentials := &types.Credentials{
		Site:     "test.com",
		Login:    "123321",
		Password: "1224",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^UPDATE credentials SET(.+)").WithArgs(credentials.Site,
		credentials.Login, credentials.Password, credentials.Metadata, userID, id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, credentials)
	require.Equal(t, err, sql.ErrNoRows)
}

func TestUpdate_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	credentials := &types.Credentials{
		Site:     "test.com",
		Login:    "123321",
		Password: "1224",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^UPDATE credentials SET(.+)").WithArgs(credentials.Site,
		credentials.Login, credentials.Password, credentials.Metadata, userID, id).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, credentials)
	require.Equal(t, err, sql.ErrConnDone)
}

func TestGet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectQuery("^SELECT site, login, password, metadata FROM credentials WHERE(.+)").WithArgs(userID, id).
		WillReturnRows(sqlmock.NewRows([]string{"site", "login", "password", "metadata"}).AddRow("123321.com", "1224", "123", "some_data"))

	store := NewRepository(db)

	ctx := context.Background()

	card, err := store.Get(ctx, userID, id)
	require.NoError(t, err)

	require.Equal(t, card.Site, "123321.com")
	require.Equal(t, card.Login, "1224")
}

func TestGet_FailId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"

	store := NewRepository(db)

	ctx := context.Background()

	_, err = store.Get(ctx, userID, "")
	require.Equal(t, err.Error(), "repository: incorrect parameters")
}

func TestGet_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectQuery("^SELECT site, login, password, metadata FROM credentials WHERE(.+)").WithArgs(userID, id).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	_, err = store.Get(ctx, userID, id)
	require.Equal(t, err, sql.ErrConnDone)
}

func TestGetKeysList_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"

	mock.ExpectQuery("^SELECT id, site FROM credentials WHERE(.+)").WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	_, err = store.GetKeysList(ctx, userID)
	require.Equal(t, err, sql.ErrConnDone)
}

func TestGetKeysList_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectQuery("^SELECT id, site FROM credentials WHERE(.+)").WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "site"}).AddRow(id, "123321.com"))

	store := NewRepository(db)

	ctx := context.Background()

	keys, err := store.GetKeysList(ctx, userID)
	require.NoError(t, err)

	require.Equal(t, len(keys), 1)
}

func TestGetKeysList_RowsErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectQuery("^SELECT id, site FROM credentials WHERE(.+)").WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "card"}).AddRow(id, "123321.com").RowError(0, sql.ErrConnDone))

	store := NewRepository(db)

	ctx := context.Background()

	_, err = store.GetKeysList(ctx, userID)
	require.Equal(t, err, sql.ErrConnDone)
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectExec("^DELETE FROM credentials WHERE (.+)").WithArgs(userID, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Delete(ctx, userID, id)
	require.NoError(t, err)
}

func TestDelete_NilCredentialId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Delete(ctx, userID, "")
	require.Equal(t, err.Error(), "repository: incorrect parameters")
}

func TestDelete_NotFoundErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectExec("^DELETE FROM credentials WHERE (.+)").WithArgs(userID, id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Delete(ctx, userID, id)
	require.Equal(t, err, sql.ErrNoRows)
}

func TestDelete_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test"
	id := "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83"

	mock.ExpectExec("^DELETE FROM credentials WHERE (.+)").WithArgs(userID, id).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Delete(ctx, userID, id)
	require.Equal(t, err, sql.ErrConnDone)
}
