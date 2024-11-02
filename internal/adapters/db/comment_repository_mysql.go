package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type CommentRepositoryMySQL struct {
	db *gorm.DB
}

func NewCommentRepositoryMySQL(db *gorm.DB) ports.CommentRepository {
	return &CommentRepositoryMySQL{db: db}
}

func (h *CommentRepositoryMySQL) Create(commentable domain.Commentable) (*domain.Commentable, error) {
	if err := h.db.Create(&commentable).Error; err != nil {
		return nil, err
	}

	if err := h.db.Preload("User").Preload("Replies.User").Preload("Hearts").First(&commentable, commentable.ID).Error; err != nil {
		return nil, err
	}

	return &commentable, nil
}

func (h *CommentRepositoryMySQL) Delete(id uint) error {
	if err := h.db.Delete(&domain.Commentable{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (h *CommentRepositoryMySQL) FindByID(id uint) (*domain.Commentable, error) {
	var comment domain.Commentable
	if err := h.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
