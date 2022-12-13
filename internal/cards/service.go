package cards

import (
	"context"
	"errors"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

var (
	ErrCardAlreadyAssociated = errors.New("card already associated to another account")
	ErrCardNotFound          = errors.New("card not found")
)

type Service interface {
	Save(ctx context.Context, id int, card domain.CardDto) error
	GetAll(ctx context.Context, accountID int) ([]domain.Card, error)
	GetByCardID(ctx context.Context, accountID, cardID int) (domain.Card, error)
	DeleteByCardID(ctx context.Context, cardID int) error
}

type service struct {
	cardsRepository Repository
}

func NewService(cardsRepository Repository) Service {
	return &service{
		cardsRepository: cardsRepository,
	}
}

func (s *service) Save(ctx context.Context, id int, card domain.CardDto) error {
	exists, err := s.cardsRepository.Exists(ctx, card.PAN)
	if err != nil {
		return err
	}

	if exists {
		return ErrCardAlreadyAssociated
	}

	id, err = s.cardsRepository.SaveCard(ctx, id, card)
	if err != nil {
		return err
	}

	if id == 0 {
		// todo error
		return err
	}

	return nil
}

func (s *service) GetAll(ctx context.Context, accountID int) ([]domain.Card, error) {
	cards, err := s.cardsRepository.GetAll(ctx, accountID)
	if err != nil {
		return []domain.Card{}, err
	}

	return cards, nil
}

func (s *service) GetByCardID(ctx context.Context, accountID, cardID int) (domain.Card, error) {
	cards, err := s.cardsRepository.GetByID(ctx, accountID, cardID)
	if err != nil {
		return domain.Card{}, err
	}

	return cards, nil
}

func (s *service) DeleteByCardID(ctx context.Context, cardID int) error {
	err := s.cardsRepository.DeleteByCardID(ctx, cardID)
	if err != nil {
		return err
	}

	return nil
}
