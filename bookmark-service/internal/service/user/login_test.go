package user

import (
	"context"
	"testing"
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	mock_user "github.com/HemlockPham7/golang-system-design/internal/repository/user/mocks"
	mock_jwtgen "github.com/HemlockPham7/golang-system-design/pkg/jwtutils/mocks"
	mock_hasher "github.com/HemlockPham7/golang-system-design/pkg/utils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepository func(ctx context.Context) *mock_user.Repository
		setupMockPasswordHash   func(t *testing.T) *mock_hasher.Hasher
		setupMockJWTGen         func(t *testing.T) *mock_jwtgen.JWTGenerator

		inputUsername string
		inputPassword string

		expectedError  error
		expectedOutput string
	}{
		{
			name: "Login successfully",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "janedoe").Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Compare", "janedoe", "janedoe").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				jwtGeneratorMock := mock_jwtgen.NewJWTGenerator(t)
				tokenContent := jwt.MapClaims{
					"sub":   "d7c13097-67a7-4eae-a60e-0b9b533b7bd5",
					"email": "janedoe@gmail.com",
					"iat":   time.Now().Unix(),
					"exp":   time.Now().Add(tokenDuration).Unix(),
				}
				jwtGeneratorMock.On("GenerateJWT", tokenContent).Return("mocked_jwt_token", nil)
				return jwtGeneratorMock
			},

			inputUsername:  "janedoe",
			inputPassword:  "janedoe",
			expectedError:  nil,
			expectedOutput: "mocked_jwt_token",
		},
		{
			name: "Fail to get user by username",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "janedoe").Return(nil, ErrInvalidCredentials)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				jwtGeneratorMock := mock_jwtgen.NewJWTGenerator(t)
				return jwtGeneratorMock
			},

			inputUsername:  "janedoe",
			inputPassword:  "janedoe",
			expectedError:  ErrInvalidCredentials,
			expectedOutput: "",
		},
		{
			name: "invalid password",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "janedoe").Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Compare", "janedoe", "janedoe").Return(false)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				jwtGeneratorMock := mock_jwtgen.NewJWTGenerator(t)
				return jwtGeneratorMock
			},

			inputUsername:  "janedoe",
			inputPassword:  "janedoe",
			expectedError:  ErrInvalidCredentials,
			expectedOutput: "",
		},
		{
			name: "Fail to generate JWT",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "janedoe").Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Compare", "janedoe", "janedoe").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				jwtGeneratorMock := mock_jwtgen.NewJWTGenerator(t)
				tokenContent := jwt.MapClaims{
					"sub":   "d7c13097-67a7-4eae-a60e-0b9b533b7bd5",
					"email": "janedoe@gmail.com",
					"iat":   time.Now().Unix(),
					"exp":   time.Now().Add(tokenDuration).Unix(),
				}
				jwtGeneratorMock.On("GenerateJWT", tokenContent).Return("", ErrCannotGenerateToken)
				return jwtGeneratorMock
			},

			inputUsername:  "janedoe",
			inputPassword:  "janedoe",
			expectedError:  ErrCannotGenerateToken,
			expectedOutput: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			userRepoMock := tc.setupMockUserRepository(ctx)
			passwordHasherMock := tc.setupMockPasswordHash(t)
			jwtGeneratorMock := tc.setupMockJWTGen(t)

			userService := NewService(userRepoMock, passwordHasherMock, jwtGeneratorMock)

			output, err := userService.Login(ctx, tc.inputUsername, tc.inputPassword)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
