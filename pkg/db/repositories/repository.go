package repositories

import (
	"database/sql"
)

type Repository interface {
	User() UserRepository
	Publisher() PublisherRepository
}

type repositoryImpl struct {
	userRepo      UserRepository
	publisherRepo PublisherRepository
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{
		userRepo:      NewUserRepository(db),
		publisherRepo: NewPublisherRepository(db),
	}
}

func (r *repositoryImpl) User() UserRepository {
	return r.userRepo
}

func (r *repositoryImpl) Publisher() PublisherRepository {
	return r.publisherRepo
}
