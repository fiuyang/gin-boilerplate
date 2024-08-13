package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"scylla/entity"
	"scylla/pkg/helper"
)

type PassResetRepo interface {
	Insert(ctx context.Context, data entity.PasswordReset) error
	InsertBatch(ctx context.Context, data []entity.PasswordReset, batchSize int) error
	Update(ctx context.Context, data entity.PasswordReset) error
	DeleteByColumns(ctx context.Context, columns []string, queries []any) error
	DeleteBatch(ctx context.Context, Ids []int) error
	FindById(ctx context.Context, Id int) (data entity.PasswordReset, err error)
	FindByColumns(ctx context.Context, columns []string, queries []any) (entity.PasswordReset, error)
}

type PassResetRepoImpl struct {
	db *gorm.DB
}

func NewPassResetRepoImpl(db *gorm.DB) PassResetRepo {
	return &PassResetRepoImpl{db: db}
}

func (repo *PassResetRepoImpl) Insert(ctx context.Context, data entity.PasswordReset) error {
	result := repo.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *PassResetRepoImpl) InsertBatch(ctx context.Context, data []entity.PasswordReset, batchSize int) error {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.CreateInBatches(&data, batchSize).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (repo *PassResetRepoImpl) Update(ctx context.Context, data entity.PasswordReset) error {
	result := repo.db.WithContext(ctx).Updates(&data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *PassResetRepoImpl) DeleteByColumns(ctx context.Context, columns []string, queries []any) error {
	if len(columns) != len(queries) {
		return errors.New("columns and queries length mismatch")
	}

	var data entity.PasswordReset
	db := repo.db.WithContext(ctx).Table("password_resets")
	for i, column := range columns {
		db = db.Where(column+" = ?", queries[i])
	}
	result := db.Delete(&data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PassResetRepoImpl) DeleteBatch(ctx context.Context, Ids []int) error {
	var data entity.Customer
	result := repo.db.WithContext(ctx).Where("id IN (?)", Ids).Delete(&data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return nil
}

func (repo *PassResetRepoImpl) FindById(ctx context.Context, Id int) (data entity.PasswordReset, err error) {
	result := repo.db.WithContext(ctx).First(&data, Id)

	if result.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func (repo *PassResetRepoImpl) FindByColumns(ctx context.Context, columns []string, queries []any) (entity.PasswordReset, error) {
	if len(columns) != len(queries) {
		return entity.PasswordReset{}, errors.New("columns and queries length mismatch")
	}

	var data entity.PasswordReset
	db := repo.db.WithContext(ctx).Table("password_resets")
	for i, column := range columns {
		db = db.Where(column+" = ?", queries[i])
	}
	result := db.First(&data)

	if result.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if result.Error != nil {
		return data, errors.New("password reset not found")
	}

	return data, nil
}
