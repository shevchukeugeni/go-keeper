package notes

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

	note := &types.Note{
		Key:      "test",
		Text:     "some_text",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^INSERT INTO texts(.+)").WithArgs(userID, id, note.Key,
		note.Text, note.Metadata).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, note)
	require.NoError(t, err)
}

func TestCreate_NilNoteInfo(t *testing.T) {
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

	note := &types.Note{
		Key:      "test",
		Text:     "some_text",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^INSERT INTO texts(.+)").WithArgs(userID, id, note.Key,
		note.Text, note.Metadata).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, note)
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

	note := &types.Note{
		Key:      "test",
		Text:     "some_text",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^INSERT INTO texts(.+)").WithArgs(userID, id, note.Key,
		note.Text, note.Metadata).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Create(ctx, userID, id, note)
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

	note := &types.Note{
		Key:      "test",
		Text:     "some_text",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^UPDATE texts SET(.+)").WithArgs(note.Key,
		note.Text, note.Metadata, userID, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, note)
	require.NoError(t, err)
}

func TestUpdate_NilNoteInfo(t *testing.T) {
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

	note := &types.Note{
		Key:      "test",
		Text:     "some_text",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^UPDATE texts SET(.+)").WithArgs(note.Key,
		note.Text, note.Metadata, userID, id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, note)
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

	note := &types.Note{
		Key:      "test",
		Text:     "some_text",
		Metadata: "test_meta",
	}

	mock.ExpectExec("^UPDATE texts SET(.+)").WithArgs(note.Key,
		note.Text, note.Metadata, userID, id).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Update(ctx, userID, id, note)
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

	mock.ExpectQuery("^SELECT key, data, metadata FROM texts WHERE(.+)").WithArgs(userID, id).
		WillReturnRows(sqlmock.NewRows([]string{"key", "data", "metadata"}).AddRow("123321", "1224", "some_data"))

	store := NewRepository(db)

	ctx := context.Background()

	note, err := store.Get(ctx, userID, id)
	require.NoError(t, err)

	require.Equal(t, note.Key, "123321")
	require.Equal(t, note.Text, "1224")
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

	mock.ExpectQuery("^SELECT key, data, metadata FROM texts WHERE(.+)").WithArgs(userID, id).
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

	mock.ExpectQuery("^SELECT id, key FROM texts WHERE(.+)").WithArgs(userID).
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

	mock.ExpectQuery("^SELECT id, key FROM texts WHERE(.+)").WithArgs(userID).
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

	mock.ExpectQuery("^SELECT id, key FROM texts WHERE(.+)").WithArgs(userID).
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

	mock.ExpectExec("^DELETE FROM texts WHERE (.+)").WithArgs(userID, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Delete(ctx, userID, id)
	require.NoError(t, err)
}

func TestDelete_NilNoteId(t *testing.T) {
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

	mock.ExpectExec("^DELETE FROM texts WHERE (.+)").WithArgs(userID, id).
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

	mock.ExpectExec("^DELETE FROM texts WHERE (.+)").WithArgs(userID, id).
		WillReturnError(sql.ErrConnDone)

	store := NewRepository(db)

	ctx := context.Background()

	err = store.Delete(ctx, userID, id)
	require.Equal(t, err, sql.ErrConnDone)
}
