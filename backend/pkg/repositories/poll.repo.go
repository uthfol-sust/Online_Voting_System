package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	dto "pollvoting/pkg/DTO"
	"pollvoting/pkg/database"
	"pollvoting/pkg/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type PollRepository interface {
	CreatePoll(poll *models.Poll) (*models.Poll, error)
	GetPolls() ([]models.Poll, error)
	GetPollDetails(id int64) (*models.Poll, error)
	UpdatePoll(ID int64, poll *models.Poll) (*models.Poll, error)
	DeletePoll(ID int64) error
	PollWithOptionsFromDB(pollID int64) (*dto.PollResponse, error)
	PollWithOptionsFromRedis(pollID int64) (*dto.PollResponse, error)
}

type pollRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewPollRepository(db *sql.DB, rdb *redis.Client) PollRepository {
	return &pollRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *pollRepository) CreatePoll(poll *models.Poll) (*models.Poll, error) {
	query := `INSERT INTO polls (title, description, created_by, expires_at, created_at) 
             VALUES ($1, $2, $3, $4, NOW()) 
			 RETURNING poll_id;`

	err := r.db.QueryRow(query, poll.Title, poll.Description, poll.CreatedBy, poll.ExpiresAt).Scan(&poll.ID)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (r *pollRepository) GetPolls() ([]models.Poll, error) {
	query := `SELECT title, COALESCE(description,'') AS description, created_by, is_active, expires_at, created_at 
              FROM polls`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polls []models.Poll

	for rows.Next() {
		var p models.Poll
		err := rows.Scan(&p.Title, &p.Description, &p.CreatedBy, &p.IsActive, &p.ExpiresAt, &p.CreatedAt)

		if err != nil {
			return nil, err
		}
		polls = append(polls, p)
	}
	return polls, nil
}

func (r *pollRepository) GetPollDetails(id int64) (*models.Poll, error) {
	query := `SELECT title, COALESCE(description,'') AS description, created_by, is_active, expires_at, created_at 
            FROM polls WHERE poll_id=$1`

	poll := &models.Poll{}

	err := r.db.QueryRow(query, id).Scan(&poll.Title, &poll.Description, &poll.CreatedBy, &poll.IsActive, &poll.ExpiresAt, &poll.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	poll.ID = id
	return poll, nil
}

func (r *pollRepository) UpdatePoll(ID int64, poll *models.Poll) (*models.Poll, error) {
	query := `UPDATE polls SET title=$1, description=$2, is_active=$3, expires_at=$4 WHERE poll_id=$5`

	_, err := r.db.Exec(query, poll.Title, poll.Description, poll.IsActive, poll.ExpiresAt, poll.ID)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (r *pollRepository) DeletePoll(ID int64) error {
	query := `DELETE FROM polls WHERE poll_id=$1`

	_, err := r.db.Exec(query, ID)

	return err
}

func (r *pollRepository) PollWithOptionsFromDB(pollID int64) (*dto.PollResponse, error) {

	query := `
			SELECT 
				p.poll_id,
				p.title,
				p.description,
				p.is_active,
				p.expires_at,
				o.option_id,
				o.option_text,
				o.score,
				o.option_image
			FROM polls p
			LEFT JOIN poll_options o ON o.poll_id = p.poll_id
			WHERE p.poll_id = $1;`

	rows, err := r.db.Query(query, pollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poll := &dto.PollResponse{}
	options := make([]models.PollOption, 0)

	for rows.Next() {

		var (
			expiresAt   sql.NullTime
			optionID    sql.NullInt64
			optionText  sql.NullString
			score       sql.NullInt64
			optionImage sql.NullString
		)

		err := rows.Scan(
			&poll.ID,
			&poll.Title,
			&poll.Description,
			&poll.IsActive,
			&expiresAt,
			&optionID,
			&optionText,
			&score,
			&optionImage,
		)
		if err != nil {
			return nil, err
		}

		if expiresAt.Valid {
			poll.ExpiresAt = expiresAt.Time
		}

		var image *string
		if optionImage.Valid {
			image = &optionImage.String
		}
		if optionID.Valid {
			option := models.PollOption{
				ID:         optionID.Int64,
				PollID:     poll.ID,
				OptionText: optionText.String,
				Score:      int(score.Int64),
				Image:      image,
			}
			options = append(options, option)
		}
	}

	if poll.ID == 0 {
		return nil, sql.ErrNoRows
	}

	poll.Options = options

	// cache
	cacheKey := fmt.Sprintf("poll:%d", pollID)
	bytes, _ := json.Marshal(poll)
	_ = r.rdb.Set(database.Ctx, cacheKey, bytes, 5*time.Minute).Err()

	return poll, nil
}

func (r *pollRepository) PollWithOptionsFromRedis(pollID int64) (*dto.PollResponse, error) {
	cacheKey := fmt.Sprintf("poll:%d", pollID)

	val, err := r.rdb.Get(database.Ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	cached := &dto.PollResponse{}
	if err := json.Unmarshal([]byte(val), cached); err != nil {
		return nil, err
	}

	return cached, nil
}
