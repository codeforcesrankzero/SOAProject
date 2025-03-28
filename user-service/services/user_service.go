package services

import (
	"errors"
	"time"
	"user-service/models"
	"user-service/repository"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    userRepo    repository.UserRepositoryInterface
    jwtSecret   []byte
    tokenExpiry time.Duration
}

func NewUserService(userRepo repository.UserRepositoryInterface, jwtSecret string, tokenExpiry time.Duration) *UserService {
    return &UserService{
        userRepo:    userRepo,
        jwtSecret:   []byte(jwtSecret),
        tokenExpiry: tokenExpiry,
    }
}

func (s *UserService) Register(req models.RegisterRequest) error {
	existingUser, err := s.userRepo.GetUserByLogin(req.Login)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("пользователь с таким логином уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Login:    req.Login,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) Login(req models.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByLogin(req.Login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("неверный логин или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("неверный логин или пароль")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(s.tokenExpiry).Unix()

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

func (s *UserService) UpdateUserProfile(userID uint, req models.UpdateProfileRequest) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	if req.BirthDate != nil {
		user.BirthDate = *req.BirthDate
	}
	user.Email = req.Email
	user.Phone = req.Phone

	return s.userRepo.UpdateUser(user)
}

func (s *UserService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("недействительный токен")
}

var _ UserServiceInterface = (*UserService)(nil)