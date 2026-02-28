package service

import "sql-injection-subtle/internal/repository"

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Search forwards the search term to the repository. No SQL is built in this layer.
func (s *UserService) Search(username string) ([]repository.User, error) {
	return s.repo.FindByUsername(username)
}

// ListSorted forwards sort column and direction to the repository.
func (s *UserService) ListSorted(sortColumn, sortDir string) ([]repository.User, error) {
	return s.repo.FindWithSort(sortColumn, sortDir)
}
