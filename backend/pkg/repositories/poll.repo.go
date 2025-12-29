package repositories

import (
	"database/sql"
	"pollvoting/pkg/models"
)

type PollRepository interface {
	CreatePoll(poll *models.Poll) (*models.Poll, error)
	GetPolls() ([]models.Poll, error)
	GetPollDetails(id int64) (*models.Poll, error)
	UpdatePoll(ID int64, poll *models.Poll) (*models.Poll, error)
	DeletePoll(ID int64) error
}

type pollRepository struct {
	db *sql.DB
}

func NewPollRepository(db *sql.DB) PollRepository {
	return &pollRepository{db: db}
}

func (r *pollRepository) CreatePoll(poll *models.Poll) (*models.Poll, error) {
	query := `INSERT INTO polls (title, description, created_by, expires_at, created_at) 
             VALUES ($1, $2, $3, $4, NOW()) 
			 RETURNING poll_id;`

	err := r.db.QueryRow(query, poll.Title, poll.Description, poll.CreatedBy, poll.ExpiresAt).Scan(&poll.ID)
	if err != nil{
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
