package repository

import (
	"fasthttptest/ferror"
	"fasthttptest/log"
	"fasthttptest/model"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *pgx.ConnPool
	redis *redis.Client
}

func NewUserRepository(db *pgx.ConnPool, redis *redis.Client) UserRepository {
	return UserRepository{db: db, redis: redis}
}

func (r UserRepository) Insert(logger log.Logger, user *model.User) error {
	sql := `
		insert into users (
			username, password, email_id, phone_number, first_name, 
			last_name, created_on, updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(sql, user.Username, user.Password, user.EmailId, user.PhoneNumber,
		user.FirstName, user.LastName, user.CreatedOn, user.UpdatedOn)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "users_pkey" {
				return ferror.New(400, ferror.UsernameNotAvailable, "Username is already taken. Change it.")
			}
		}
	}
	return err
}

func (r UserRepository) CreateSession(username string, password string) (sessionId string, err error) {
	sql := `
		select count(username) from users where username = $1 and password = $2`
	var count uint64
	err = r.db.QueryRow(sql, username, password).Scan(&count)
	if err != nil {return}
	if count != 1 {
		err = ferror.New(401, ferror.InvalidCredential, "Username or Password does not match")
		return
	}
	return r.createSession(username)
}

func (r UserRepository) GetUser(username string) (user model.User, err error) {
	sql := `
		select 
			username, password, email_id, phone_number, 
			first_name, last_name, created_on, updated_on 
		from users where username = $1`
	err = r.db.QueryRow(sql, username).Scan(
		&user.Username, &user.Password, &user.EmailId, &user.PhoneNumber,
		&user.FirstName, &user.LastName, &user.CreatedOn, &user.UpdatedOn)
	return
}

func (r UserRepository) createSession(username string) (sessionId string, err error) {
	sessionId = uuid.New().String()
	_, err = r.redis.Set(rSessionId + sessionId, username, 0).Result()
	if err != nil {return}
	_, err = r.redis.SAdd(rUserSession + username, sessionId).Result()
	return
}

func (r UserRepository) GetSession(sessionId string) (username string, err error) {
	username, err = r.redis.Get(rSessionId + sessionId).Result()
	return
}

func (r UserRepository) RemoveSession(sessionId string) error {
	username, err := r.redis.Get(rSessionId + sessionId).Result()
	_, err = r.redis.SRem(rUserSession + username, sessionId).Result()
	if err != nil {return err}
	_, err = r.redis.Del(rSessionId + sessionId).Result()
	return err
}

func (r UserRepository) RemoveAllSession(logger log.Logger, username string) error {
	sessions, err := r.redis.SMembers(rUserSession + username).Result()
	if err != nil {return err}
	for i := 0; i < len(sessions); i++ {
		logger.Printf("Session: %s", sessions[i])
		_, err := r.redis.Del(rSessionId + sessions[i]).Result()
		if err != nil {return err}
	}
	_, err = r.redis.Del(rUserSession + username).Result()
	return err
}