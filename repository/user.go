package repository

import (
	"fasthttptest/ferror"
	"fasthttptest/log"
	"fasthttptest/model"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type User struct {
	db *pgx.ConnPool
	redis *redis.Client
}

func NewUserRepository(db *pgx.ConnPool, redis *redis.Client) User {
	return User{db: db, redis: redis}
}

func (r User) Insert(logger log.Logger, user *model.User) error {
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

func (r User) CreateSession(username string, password string) (sessionId string, err error) {
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

func (r User) GetUser(username string) (user model.User, err error) {
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

func (r User) createSession(username string) (sessionId string, err error) {
	sessionId = uuid.New().String()
	_, err = r.redis.Set(rSessionId + sessionId, username, 0).Result()
	if err != nil {return}
	_, err = r.redis.SAdd(rUserSession + username, sessionId).Result()
	return
}

func (r User) GetSession(sessionId string) (username string, err error) {
	username, err = r.redis.Get(rSessionId + sessionId).Result()
	if err == redis.Nil {
		err = ferror.New(401, ferror.Unauthorized, "You are not logged in")
	}
	return
}

func (r User) RemoveSession(sessionId string) error {
	username, err := r.redis.Get(rSessionId + sessionId).Result()
	_, err = r.redis.SRem(rUserSession + username, sessionId).Result()
	if err != nil {return err}
	_, err = r.redis.Del(rSessionId + sessionId).Result()
	return err
}

func (r User) RemoveAllSession(_ log.Logger, username string) error {
	sessions, err := r.redis.SMembers(rUserSession + username).Result()
	if err != nil {return err}
	for i := 0; i < len(sessions); i++ {
		_, err := r.redis.Del(rSessionId + sessions[i]).Result()
		if err != nil {return err}
	}
	_, err = r.redis.Del(rUserSession + username).Result()
	return err
}