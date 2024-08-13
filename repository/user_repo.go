package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"scylla/dto"
	"scylla/entity"
	"scylla/pkg/helper"
)

type UserRepo interface {
	Insert(ctx context.Context, data entity.User) error
	InsertBatch(ctx context.Context, data []entity.User, batchSize int) error
	Update(ctx context.Context, data entity.User) error
	DeleteBatch(ctx context.Context, Ids []int) error
	FindAll(ctx context.Context, dataFilter dto.UserQueryFilter) (domain []entity.User, err error)
	FindById(ctx context.Context, Id int) (data entity.User, err error)
	FindByColumns(ctx context.Context, columns []string, queries []any) (entity.User, error)
	CheckColumnExists(ctx context.Context, column string, value interface{}) bool
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) UserRepo {
	return &UserRepoImpl{db: db}
}

func (repo *UserRepoImpl) Insert(ctx context.Context, data entity.User) error {
	result := repo.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *UserRepoImpl) InsertBatch(ctx context.Context, data []entity.User, batchSize int) error {
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

func (repo *UserRepoImpl) Update(ctx context.Context, data entity.User) error {
	result := repo.db.WithContext(ctx).Updates(&data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *UserRepoImpl) DeleteBatch(ctx context.Context, Ids []int) error {
	var data entity.User
	result := repo.db.WithContext(ctx).Where("id IN (?)", Ids).Delete(&data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return nil
}

func (repo *UserRepoImpl) FindAll(ctx context.Context, dataFilter dto.UserQueryFilter) (domain []entity.User, err error) {
	query := "SELECT * FROM users"
	args := []interface{}{}

	if dataFilter.Username != "" {
		query += " WHERE username = ?"
		args = append(args, dataFilter.Username)
	}

	if dataFilter.Email != "" {
		query += " WHERE email = ?"
		args = append(args, dataFilter.Email)
	}

	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		query += " WHERE created_at BETWEEN ? AND ?"
		args = append(args, dataFilter.StartDate, dataFilter.EndDate)
	}

	rows, err := repo.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		domain = append(domain, user)
	}

	return domain, nil
}

func (repo *UserRepoImpl) FindById(ctx context.Context, Id int) (data entity.User, err error) {
	result := repo.db.WithContext(ctx).First(&data, Id)

	if result.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func (repo *UserRepoImpl) FindByColumns(ctx context.Context, columns []string, queries []any) (entity.User, error) {
	if len(columns) != len(queries) {
		return entity.User{}, errors.New("columns and queries length mismatch")
	}

	var data entity.User
	db := repo.db.WithContext(ctx).Table("users")
	for i, column := range columns {
		db = db.Where(column+" = ?", queries[i])
	}
	result := db.First(&data)

	if result.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if result.Error != nil {
		return data, errors.New("users not found")
	}

	return data, nil
}

func (repo *UserRepoImpl) CheckColumnExists(ctx context.Context, column string, value interface{}) bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM users WHERE %s = ?)", column)
	err := repo.db.WithContext(ctx).Raw(query, value).Scan(&exists).Error
	if err != nil {
		return false
	}
	return exists
}
