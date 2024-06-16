package cards

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

func NewRepository(db *sql.DB) store.Secrets[types.CardInfo] {
	return &repo{db: db}
}

func (repo *repo) Create(ctx context.Context, userID, id string, cardInfo *types.CardInfo) error {
	if cardInfo == nil {
		return errors.New("repository: incorrect parameters")
	}

	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO cards(user_id, id, card, expiration, cvv, metadata) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, id, cardInfo.Number, cardInfo.Expiration, cardInfo.CVV, cardInfo.Metadata)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return types.ErrRecordAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *repo) Get(ctx context.Context, userID, id string) (*types.CardInfo, error) {
	if id == "" {
		return nil, errors.New("repository: incorrect parameters")
	}

	ret := types.CardInfo{}

	err := repo.db.QueryRowContext(ctx, "SELECT card, expiration, cvv, metadata FROM cards WHERE user_id=$1 and id=$2",
		userID, id).Scan(&ret.Number, &ret.Expiration, &ret.CVV, &ret.Metadata)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (repo *repo) GetKeysList(ctx context.Context, userID string) ([]types.Key, error) {
	var ret []types.Key

	rows, err := repo.db.QueryContext(ctx, "SELECT id, card FROM cards WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, number string
		err = rows.Scan(&id, &number)
		if err != nil {
			return nil, err
		}

		ret = append(ret, types.Key{Id: id, Key: number})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (repo *repo) Update(ctx context.Context, userID, id string, cardInfo *types.CardInfo) error {
	if id == "" {
		return errors.New("repository: incorrect parameters")
	}

	result, err := repo.db.ExecContext(ctx, "UPDATE cards SET card=$1, expiration=$2, cvv = $3, metadata=$4 WHERE user_id=$5 and id =$6;",
		cardInfo.Number, cardInfo.Expiration, cardInfo.CVV, cardInfo.Metadata, userID, id)
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

	result, err := repo.db.ExecContext(ctx, "DELETE FROM cards WHERE user_id=$1 and id=$2;",
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
