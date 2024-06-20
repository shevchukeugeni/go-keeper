package creds

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

func NewRepository(db *sql.DB) store.Secrets[types.Credentials] {
	return &repo{db: db}
}

func (repo *repo) Create(ctx context.Context, userID, id string, creds *types.Credentials) error {
	if creds == nil {
		return errors.New("repository: incorrect parameters")
	}

	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO credentials(user_id, id, site, login, password, metadata) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, id, creds.Site, creds.Login, creds.Password, creds.Metadata)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return types.ErrRecordAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *repo) Get(ctx context.Context, userID, id string) (*types.Credentials, error) {
	if id == "" {
		return nil, errors.New("repository: incorrect parameters")
	}

	ret := types.Credentials{}

	err := repo.db.QueryRowContext(ctx, "SELECT site, login, password, metadata FROM credentials WHERE user_id=$1 and id=$2",
		userID, id).Scan(&ret.Site, &ret.Login, &ret.Password, &ret.Metadata)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (repo *repo) GetKeysList(ctx context.Context, userID string) ([]types.Key, error) {
	var ret []types.Key

	rows, err := repo.db.QueryContext(ctx, "SELECT id, site FROM credentials WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, site string
		err = rows.Scan(&id, &site)
		if err != nil {
			return nil, err
		}

		ret = append(ret, types.Key{Id: id, Key: site})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (repo *repo) Update(ctx context.Context, userID, id string, creds *types.Credentials) error {
	if id == "" || creds == nil {
		return errors.New("repository: incorrect parameters")
	}

	result, err := repo.db.ExecContext(ctx, "UPDATE credentials SET site=$1, login=$2, password=$3, metadata=$4 WHERE user_id=$5 and id =$6;",
		creds.Site, creds.Login, creds.Password, creds.Metadata, userID, id)
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

func (repo *repo) Delete(ctx context.Context, userID, id string) error {
	if id == "" {
		return errors.New("repository: incorrect parameters")
	}

	result, err := repo.db.ExecContext(ctx, "DELETE FROM credentials WHERE user_id=$1 and id =$2;",
		userID, id)
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
