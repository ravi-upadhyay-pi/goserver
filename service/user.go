package service

import (
	"fasthttptest/log"
	"fasthttptest/model"
	"fasthttptest/repository"
	"time"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func (s *UserService) Add(logger log.Logger, user *model.User) error {
	user.CreatedOn = time.Now()
	user.UpdatedOn = time.Now()
	return s.userRepository.Insert(logger, user)
}

func (s *UserService) CreateSession(logger log.Logger, username string, password string) (string, error) {
	return s.userRepository.CreateSession(username, password)
}

func (s *UserService) GetSession(logger log.Logger, sessionId string) (string, error) {
	return s.userRepository.GetSession(sessionId)
}

func (s *UserService) GetProfile(logger log.Logger, sessionId string) (user model.User, err error) {
	username, err := s.userRepository.GetSession(sessionId)
	if err != nil {return}
	return s.userRepository.GetUser(username)
}

func (s *UserService) RemoveSession(logger log.Logger, sessionId string) error {
	return s.userRepository.RemoveSession(sessionId)
}

func (s *UserService) RemoveAllSession(logger log.Logger, sessionId string) error {
	username, err := s.userRepository.GetSession(sessionId)
	if err != nil {return err}
	return s.userRepository.RemoveAllSession(logger, username)
}