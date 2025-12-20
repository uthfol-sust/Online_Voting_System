package repositories

import (
	"database/sql"
	"pollvoting/pkg/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id int64) (*models.User, error)
	Update(id int64, user *models.User) error
	Delete(id int64) error
}

type userRepository struct {
	post_db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{post_db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	query := ` INSERT INTO users (user_name, email, password, created_at) 
             VALUES ($1, $2, $3, NOW()) 
			 RETURNING user_id;
			  `
	err := r.post_db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
    query := `SELECT user_id, user_name, email, COALESCE(image, '') AS image, COALESCE(role, 'user') AS role, created_at FROM users`
    rows, err := r.post_db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User

    for rows.Next() {
        var u models.User
        err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Image, &u.Role, &u.CreatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
    query := `SELECT user_id, user_name, email, COALESCE(image, '') AS image, COALESCE(role, 'user') AS role, created_at 
              FROM users WHERE user_id = $1`

    u := &models.User{}
    err := r.post_db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Image, &u.Role, &u.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return u, nil
}


func (r *userRepository) Update(id int64, user *models.User) error {
	query := `UPDATE users SET user_name = $1, email = $2, password=$3, image = $4, role = $5 WHERE user_id = $6`

	_, err := r.post_db.Exec(query, user.Name, user.Email, user.Password, user.Image, user.Role, id)

	return err
}

func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := r.post_db.Exec(query, id)
	return err
}
