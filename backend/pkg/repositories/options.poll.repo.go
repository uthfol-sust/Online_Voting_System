package repositories

import (
	"database/sql"
	"fmt"
	"pollvoting/pkg/database"
	"pollvoting/pkg/models"

	"github.com/redis/go-redis/v9"
)

type PollOptionsRepository interface {
	Create(options *models.PollOption) (*models.PollOption, error)
	Delete(ID int64) error
}

type pollOptionsRepository struct {
	post_db  *sql.DB
	redis_db *redis.Client
}

func NewPollOptionsRepository(db *sql.DB, rdb *redis.Client) PollOptionsRepository {
	return &pollOptionsRepository{
		post_db:  db,
		redis_db: rdb,
	}
}

func (r *pollOptionsRepository) Create(options *models.PollOption) (*models.PollOption, error) {

	query := `
		INSERT INTO poll_options (poll_id, option_image, option_text)
		VALUES ($1, $2, $3)
		RETURNING option_id;
	`

	err := r.post_db.QueryRow(
		query,
		options.PollID,
		options.Image,
		options.OptionText,
	).Scan(&options.ID)

	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("poll:%d", options.PollID)
	_ = r.redis_db.Del(database.Ctx, cacheKey).Err()

	return options, nil
}

func (r *pollOptionsRepository) Delete(optionID int64) error {

	query := `
		DELETE FROM poll_options
		WHERE option_id = $1
		RETURNING poll_id;
	`

	var pollID int64

	err := r.post_db.QueryRow(query, optionID).Scan(&pollID)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("poll:%d", pollID)
	_ = r.redis_db.Del(database.Ctx, cacheKey).Err()

	return nil
}
