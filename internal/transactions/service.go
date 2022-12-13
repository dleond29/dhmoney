package transactions

import (
	"context"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

type Service interface {
	GetTransactionsLastFive(ctx context.Context, id int) ([]domain.TransactionInfo, error)
}

type service struct {
	transactionsRepository Repository
}

func NewService(transactionsRepository Repository) Service {
	return &service{
		transactionsRepository: transactionsRepository,
	}
}

func (s *service) GetTransactionsLastFive(ctx context.Context, id int) ([]domain.TransactionInfo, error) {
	trx, err := s.transactionsRepository.GetAllByIDLimit(ctx, id, 5)
	if err != nil {
		return []domain.TransactionInfo{}, err
	}

	return trx, nil
}
