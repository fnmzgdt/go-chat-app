package users

import (
	"context"
	"database/sql"
	"mime/multipart"
	"net/url"
	"path"

	"github.com/google/uuid"
)

type Repository interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	ExecuteGetUserPasswordFromEmail(query string, values ...interface{}) (string, error)
	ExecuteGetProfileImageFromUserId(query string, userId int) (string, error)
}

type GcBucketRepository interface {
	UploadProfileImage(ctx context.Context, objName string, imageFile multipart.File) (string, error)
	DeleteProfileImage(ctx context.Context, objName string) error
}

type Service interface {
	CreateUser(user *UserSignup) (sql.Result, error)
	GetUserPasswordFromEmail(email string) (string, error)
	UpdateProfileImage(ctx context.Context, userId int, imageUrl string) (sql.Result, error)      //in mysql
	UploadProfileImage(ctx context.Context, userId int, imageFile multipart.File) (string, error) // in google cloud bucket
	DeleteProfileImage(ctx context.Context, userId int) (error)
}

type service struct {
	mysql    Repository
	gcbucket GcBucketRepository
}

func NewService(r Repository, gcb GcBucketRepository) Service {
	return &service{r, gcb}
}

func (s *service) CreateUser(user *UserSignup) (sql.Result, error) {
	query := `INSERT INTO users (username, email, password) VALUES(?, ?, ?);`
	result, err := s.mysql.ExecuteQuery(query, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetUserPasswordFromEmail(email string) (string, error) {
	query := `SELECT CAST(password AS CHAR) FROM users WHERE email=?;`
	password, err := s.mysql.ExecuteGetUserPasswordFromEmail(query, email)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (s *service) UpdateProfileImage(ctx context.Context, userId int, imageUrl string) (sql.Result, error) { //mysql
	query := "UPDATE users SET profile_image_url = NULLIF(?, '') WHERE id = ?;"
	result, err := s.mysql.ExecuteQuery(query, imageUrl, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) UploadProfileImage(ctx context.Context, userId int, imageFile multipart.File) (string, error) { //google cloud bucket
	query := "SELECT profile_image_url FROM users WHERE id = ?;"
	imageUrl, err := s.mysql.ExecuteGetProfileImageFromUserId(query, userId)
	if err != nil {
		return "", err
	}

	objName, err := objNameFromUrl(imageUrl)
	if err != nil {
		return "", err
	}

	imageUrl, err = s.gcbucket.UploadProfileImage(ctx, objName, imageFile)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}

func (s *service) DeleteProfileImage(ctx context.Context, userId int) (error) {
	query := "SELECT profile_image_url FROM users WHERE id = ?;"

	imageUrl, err := s.mysql.ExecuteGetProfileImageFromUserId(query, userId)
	if err != nil {
		return err
	}

	if imageUrl == "" {
		return nil
	}

	objName, err := objNameFromUrl(imageUrl)
	if err != nil {
		return err
	}

	err = s.gcbucket.DeleteProfileImage(ctx, objName)
	if err != nil {
		return err
	}

	return nil
}

func objNameFromUrl(imageUrl string) (string, error) {
	if imageUrl == "" {
		objId, _ := uuid.NewRandom()
		return objId.String(), nil
	}

	urlPath, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}

	return path.Base(urlPath.Path), nil
}
