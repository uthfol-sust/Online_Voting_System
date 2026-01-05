package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	dto "pollvoting/pkg/DTO"
	"github.com/redis/go-redis/v9"
)

type CastVoteRepository interface {
	CastVote(ctx context.Context, pollID, optionID, userID int64) error
	PollResultsFromDB(ctx context.Context, pollID int64) (*dto.PollResultResponse, error)
	PollResultsFromRedis(ctx context.Context, pollID int64) (*dto.PollResultResponse, error)
}

type castVoteRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewCastVoteRepository(db *sql.DB, rdb *redis.Client) CastVoteRepository {
	return &castVoteRepository{db: db, rdb: rdb}
}

func (r *castVoteRepository) CastVote(ctx context.Context, pollID, optionID, userID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO votes (poll_id, option_id, user_id)
		VALUES ($1, $2, $3)
	`, pollID, optionID, userID)

	if err != nil {
		return errors.New("user already voted or invalid poll/option")
	}

	_, err = tx.Exec(`
		UPDATE poll_options
		SET score = score + 1
		WHERE option_id = $1 AND poll_id = $2
	`, optionID, pollID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	voterKey := fmt.Sprintf("poll:%d:voters", pollID)
	optionKey := fmt.Sprintf("poll:%d:options", pollID)
	cacheKey := fmt.Sprintf("poll:%d", pollID)

	pipe := r.rdb.Pipeline()
	pipe.SAdd(ctx, voterKey, userID)
	pipe.HIncrBy(ctx, optionKey, strconv.FormatInt(optionID, 10), 1)
	pipe.Expire(ctx, voterKey, 5*time.Minute)
	pipe.Expire(ctx, optionKey, 5*time.Minute)
	pipe.Del(ctx, cacheKey)

	_, _ = pipe.Exec(ctx)

	return nil
}


func (r *castVoteRepository) PollResultsFromRedis(ctx context.Context, pollID int64) (*dto.PollResultResponse, error) {
	optionKey := fmt.Sprintf("poll:%d:options", pollID)

	vals, err := r.rdb.HGetAll(ctx, optionKey).Result()
	if err != nil || len(vals) == 0 {
		return nil, nil
	}

	results := make([]dto.PollOptionResult, 0, len(vals))
	for k, v := range vals {
		optionID, err1 := strconv.ParseInt(k, 10, 64)
		votes, err2 := strconv.ParseInt(v, 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		results = append(results, dto.PollOptionResult{
			OptionID: optionID,
			Votes:    votes,
		})
	}

	return &dto.PollResultResponse{
		PollID:  pollID,
		Results: results,
	}, nil
}


func (r *castVoteRepository) PollResultsFromDB(ctx context.Context, pollID int64) (*dto.PollResultResponse, error) {

	rows, err := r.db.Query(`
		SELECT option_id, score
		FROM poll_options
		WHERE poll_id = $1
	`, pollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []dto.PollOptionResult{}
	optionKey := fmt.Sprintf("poll:%d:options", pollID)
	pipe := r.rdb.Pipeline()

	for rows.Next() {
		var res dto.PollOptionResult
		if err := rows.Scan(&res.OptionID, &res.Votes); err != nil {
			return nil, err
		}
		results = append(results, res)
		pipe.HSet(ctx, optionKey, res.OptionID, res.Votes)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		pipe.Expire(ctx, optionKey, 5*time.Minute)
		_, _ = pipe.Exec(ctx)
	}

	return &dto.PollResultResponse{
		PollID:  pollID,
		Results: results,
	}, nil
}

