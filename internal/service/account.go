package service

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// 	"github.com/realdanielursul/simbir-go/internal/entity"
// 	"github.com/realdanielursul/simbir-go/internal/repository"
// 	"github.com/realdanielursul/simbir-go/pkg/hasher"
// )

// type TokenClaims struct {
// 	jwt.StandardClaims
// 	UserID  int64
// 	IsAdmin bool
// }

// type AccountService struct {
// 	accountRepo    repository.Account
// 	tokenRepo      repository.Token
// 	passwordHasher hasher.PasswordHasher
// 	signKey        string
// 	tokenTTL       time.Duration
// }

// func NewAccountService(accountRepo repository.Account, tokenRepo repository.Token, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AccountService {
// 	return &AccountService{
// 		accountRepo:    accountRepo,
// 		tokenRepo:      tokenRepo,
// 		passwordHasher: passwordHasher,
// 		signKey:        signKey,
// 		tokenTTL:       tokenTTL,
// 	}
// }

// func (s *AccountService) SignUp(ctx context.Context, input *AccountInput) (int64, error) {
// 	// check username uniqueness
// 	account, err := s.accountRepo.GetByUsername(ctx, input.Username)
// 	if err != nil {
// 		return -1, err
// 	}

// 	if account != nil {
// 		return -1, ErrUsernameAlreadyExists
// 	}

// 	// validate username
// 	if err := validateUsername(input.Username); err != nil {
// 		return -1, err
// 	}

// 	// validate password
// 	if err := validatePassword(input.Password); err != nil {
// 		return -1, err
// 	}

// 	// create account in db
// 	id, err := s.accountRepo.Create(ctx, &entity.Account{
// 		Username:     input.Username,
// 		PasswordHash: s.passwordHasher.Hash(input.Password),
// 	})
// 	if err != nil {
// 		return -1, err
// 	}

// 	return id, nil
// }

// func (s *AccountService) SignIn(ctx context.Context, input *AccountInput) (string, error) {
// 	// get account from db
// 	account, err := s.accountRepo.GetByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
// 	if err != nil {
// 		return "", err
// 	}

// 	// create token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 		},
// 		UserID:  account.ID,
// 		IsAdmin: account.IsAdmin,
// 	})

// 	tokenString, err := token.SignedString([]byte(s.signKey))
// 	if err != nil {
// 		return "", ErrCannotSignToken
// 	}

// 	err = s.tokenRepo.CreateToken(ctx, &entity.Token{
// 		UserID:      account.ID,
// 		TokenString: tokenString,
// 		IsValid:     true,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func (s *AccountService) SignOut(ctx context.Context, id int64) error {
// 	if err := s.tokenRepo.InvalidateUserTokens(ctx, id); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *AccountService) GetAccount(ctx context.Context, id int64) (*AccountOutput, error) {
// 	account, err := s.accountRepo.GetByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if account == nil {
// 		return nil, ErrAccountNotFound
// 	}

// 	return &AccountOutput{
// 		ID:        account.ID,
// 		Username:  account.Username,
// 		Balance:   float64(account.Balance) / 100,
// 		CreatedAt: account.CreatedAt,
// 		UpdatedAt: account.UpdatedAt,
// 	}, nil
// }

// func (s *AccountService) UpdateAccount(ctx context.Context, id int64, input *AccountInput) error {
// 	// check username uniqueness
// 	account, err := s.accountRepo.GetByUsername(ctx, input.Username)
// 	if err != nil {
// 		return err
// 	}

// 	if account != nil {
// 		return ErrUsernameAlreadyExists
// 	}

// 	// validate username
// 	if err := validateUsername(input.Username); err != nil {
// 		return err
// 	}

// 	// validate password
// 	if err := validatePassword(input.Password); err != nil {
// 		return err
// 	}

// 	// update account
// 	err = s.accountRepo.Update(ctx, &entity.Account{
// 		ID:           id,
// 		Username:     input.Username,
// 		PasswordHash: s.passwordHasher.Hash(input.Password),
// 		IsAdmin:      account.IsAdmin,
// 		Balance:      account.Balance,
// 		UpdatedAt:    time.Now(),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *AccountService) ValidateToken(ctx context.Context, tokenString string) (int64, bool, error) {
// 	_, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(s.signKey), nil
// 	})
// 	if err != nil {
// 		return -1, false, ErrCannotParseToken
// 	}

// 	token, err := s.tokenRepo.GetToken(ctx, tokenString)
// 	if err != nil {
// 		return -1, false, err
// 	}

// 	return token.UserID, token.IsValid, nil
// }

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////

// func validateUsername(username string) error {
// 	return nil
// }

// func validatePassword(password string) error {
// 	return nil
// }
