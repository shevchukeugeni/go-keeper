package notes

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"keeper-project/internal/store"
	"keeper-project/types"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) store.Secrets[types.Note] {
	return &repo{db: db}
}

func (repo *repo) Create(ctx context.Context, userID, id string, text *types.Note) error {
	if text == nil {
		return errors.New("repository: incorrect parameters")
	}

	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO texts(user_id, id, key, data, metadata) VALUES ($1, $2, $3, $4, $5)",
		userID, id, text.Key, text.Text, text.Metadata)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return types.ErrRecordAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *repo) Get(ctx context.Context, userID, id string) (*types.Note, error) {
	if id == "" {
		return nil, errors.New("repository: incorrect parameters")
	}

	ret := types.Note{}

	err := repo.db.QueryRowContext(ctx, "SELECT key, data, metadata FROM texts WHERE user_id=$1 and id=$2", userID, id).Scan(
		&ret.Key, &ret.Text, &ret.Metadata)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (repo *repo) GetKeysList(ctx context.Context, userID string) ([]types.Key, error) {
	var ret []types.Key

	rows, err := repo.db.QueryContext(ctx, "SELECT id, key FROM texts WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, key string
		err = rows.Scan(&id, &key)
		if err != nil {
			return nil, err
		}

		ret = append(ret, types.Key{Id: id, Key: key})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (repo *repo) Update(ctx context.Context, userID, id string, text *types.Note) error {
	if id == "" || text == nil {
		return errors.New("repository: incorrect parameters")
	}

	result, err := repo.db.ExecContext(ctx, "UPDATE texts SET key=$1, data=$2, metadata = $3 WHERE user_id=$4 and id =$5;",
		text.Key, text.Text, text.Metadata, userID, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *repo) Delete(ctx context.Context, userID, key string) error {
	if key == "" {
		return errors.New("repository: incorrect parameters")
	}

	result, err := repo.db.ExecContext(ctx, "DELETE FROM texts WHERE user_id=$1 and id =$2;",
		userID, key)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return sql.ErrNoRows
	}
	return nil
}
