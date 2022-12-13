package cards

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) GetAll(ctx context.Context, accountID int) ([]domain.Card, error) {
	args := r.Called(ctx, accountID)
	return args.Get(0).([]domain.Card), args.Error(1)
}

func (r *repositoryMock) GetByID(ctx context.Context, accountID, cardID int) (domain.Card, error) {
	args := r.Called(ctx, accountID, cardID)
	return args.Get(0).(domain.Card), args.Error(1)
}

func (r *repositoryMock) Exists(ctx context.Context, pan string) (bool, error) {
	args := r.Called(ctx, pan)
	return args.Bool(0), args.Error(1)
}

func (r *repositoryMock) SaveCard(ctx context.Context, id int, card domain.CardDto) (int, error) {
	args := r.Called(ctx, id, card)
	return args.Int(0), args.Error(1)
}

func (r *repositoryMock) DeleteByCardID(ctx context.Context, cardID int) error {
	args := r.Called(ctx, cardID)
	return args.Error(0)
}

func Test_service_Save(t *testing.T) {
	var ctx = context.Background()
	card := domain.CardDto{
		PAN: "123456",
	}
	id := 1
	testCases := []struct {
		name          string
		repoMock      func(m *mock.Mock)
		expectedError error
	}{
		{
			name: "Error exists card",
			repoMock: func(m *mock.Mock) {
				m.On("Exists", ctx, card.PAN).Return(false, errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		{
			name: "Card exists true",
			repoMock: func(m *mock.Mock) {
				m.On("Exists", ctx, card.PAN).Return(true, nil)
			},
			expectedError: errors.New("card already associated to another account"),
		},
		{
			name: "Card save error",
			repoMock: func(m *mock.Mock) {
				m.On("Exists", ctx, card.PAN).Return(false, nil)
				m.On("SaveCard", ctx, id, card).Return(0, errors.New("error save card"))
			},
			expectedError: errors.New("error save card"),
		},
		{
			name: "Card save cero",
			repoMock: func(m *mock.Mock) {
				m.On("Exists", ctx, card.PAN).Return(false, nil)
				m.On("SaveCard", ctx, id, card).Return(0, nil)
			},
			expectedError: nil,
		},
		{
			name: "Card save successfully",
			repoMock: func(m *mock.Mock) {
				m.On("Exists", ctx, card.PAN).Return(false, nil)
				m.On("SaveCard", ctx, id, card).Return(1, nil)
			},
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repoMock := new(repositoryMock)
			testCase.repoMock(&repoMock.Mock)

			cardsService := NewService(repoMock)

			err := cardsService.Save(ctx, id, card)

			assert.Equal(t, testCase.expectedError, err)

		})
	}
}

func Test_service_GetAll(t *testing.T) {
	var ctx = context.Background()
	accountID := 1
	testCases := []struct {
		name           string
		repoMock       func(m *mock.Mock)
		id             int
		expectedError  error
		expectedResult []domain.Card
	}{
		{
			name: "Error get all card repository",
			repoMock: func(m *mock.Mock) {
				m.On("GetAll", ctx, accountID).Return([]domain.Card{}, errors.New("error"))
			},
			expectedError:  errors.New("error"),
			expectedResult: []domain.Card{},
		},
		{
			name: "Get all succesfully",
			repoMock: func(m *mock.Mock) {
				m.On("GetAll", ctx, accountID).Return([]domain.Card{
					{
						ID:             1,
						AccountID:      1,
						PAN:            "12345",
						HolderName:     "Miguel",
						ExpirationDate: "2025-10-10",
						CID:            "124",
						Type:           "Master Card",
					},
					{
						ID:             2,
						AccountID:      1,
						PAN:            "323432",
						HolderName:     "Miguel",
						ExpirationDate: "2025-10-10",
						CID:            "143",
						Type:           "Visa",
					},
				}, nil)
			},
			expectedError: nil,
			expectedResult: []domain.Card{
				{
					ID:             1,
					AccountID:      1,
					PAN:            "12345",
					HolderName:     "Miguel",
					ExpirationDate: "2025-10-10",
					CID:            "124",
					Type:           "Master Card",
				},
				{
					ID:             2,
					AccountID:      1,
					PAN:            "323432",
					HolderName:     "Miguel",
					ExpirationDate: "2025-10-10",
					CID:            "143",
					Type:           "Visa",
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repoMock := new(repositoryMock)
			testCase.repoMock(&repoMock.Mock)

			cardsService := NewService(repoMock)

			cards, err := cardsService.GetAll(ctx, accountID)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, cards)
		})
	}
}

func Test_service_GetByCardID(t *testing.T) {
	var ctx = context.Background()
	accountID := 1
	cardID := 1
	testCases := []struct {
		name           string
		repoMock       func(m *mock.Mock)
		id             int
		expectedError  error
		expectedResult domain.Card
	}{
		{
			name: "Error get by id card repository",
			repoMock: func(m *mock.Mock) {
				m.On("GetByID", ctx, accountID, cardID).Return(domain.Card{}, errors.New("error"))
			},
			expectedError:  errors.New("error"),
			expectedResult: domain.Card{},
		},
		{
			name: "Get by id card succesfully",
			repoMock: func(m *mock.Mock) {
				m.On("GetByID", ctx, accountID, cardID).Return(domain.Card{
					ID:             1,
					AccountID:      1,
					PAN:            "12345",
					HolderName:     "Miguel",
					ExpirationDate: "2025-10-10",
					CID:            "124",
					Type:           "Master Card",
				}, nil)
			},
			expectedError: nil,
			expectedResult: domain.Card{
				ID:             1,
				AccountID:      1,
				PAN:            "12345",
				HolderName:     "Miguel",
				ExpirationDate: "2025-10-10",
				CID:            "124",
				Type:           "Master Card",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repoMock := new(repositoryMock)
			testCase.repoMock(&repoMock.Mock)

			cardsService := NewService(repoMock)

			card, err := cardsService.GetByCardID(ctx, accountID, cardID)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, card)
		})
	}
}

func Test_service_DeleteByCardID(t *testing.T) {
	var ctx = context.Background()
	cardID := 1
	testCases := []struct {
		name          string
		repoMock      func(m *mock.Mock)
		id            int
		expectedError error
	}{
		{
			name: "Error delete by id card repository",
			repoMock: func(m *mock.Mock) {
				m.On("DeleteByCardID", ctx, cardID).Return(errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		{
			name: "Delete by id card successfully",
			repoMock: func(m *mock.Mock) {
				m.On("DeleteByCardID", ctx, cardID).Return(nil)
			},
		},

		// TODO: Add test cases.
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repoMock := new(repositoryMock)
			testCase.repoMock(&repoMock.Mock)

			cardsService := NewService(repoMock)

			err := cardsService.DeleteByCardID(ctx, cardID)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
