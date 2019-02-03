package service

import (
	"fasthttptest/log"
	"fasthttptest/model"
	"fasthttptest/repository"
	"time"
)

type User struct {
	Repository *repository.User
}

func (s *User) Add(logger log.Logger, user *model.User) error {
	user.CreatedOn = time.Now()
	user.UpdatedOn = time.Now()
	return s.Repository.Insert(logger, user)
}

func (s *User) GetProfile(logger log.Logger, sessionId string) (user model.User, err error) {
	username, err := s.Repository.GetSession(sessionId)
	if err != nil {return}
	return s.Repository.GetUser(username)
}

func (s *User) RemoveAllSession(logger log.Logger, sessionId string) error {
	username, err := s.Repository.GetSession(sessionId)
	if err != nil {return err}
	return s.Repository.RemoveAllSession(logger, username)
}