package wallet

import (
	 "github.com/Ilhom5005/wallet/v1/pkg/types"
	 "github.com/google/uuid"
	 "errors"
)



var(
	ErrPhoneregistered = errors.New("phone alredy registered")
	ErrAmountMustBePossitive = errors.New("amount must be greater then zero")
	ErrAccountNotFound = errors.New("account not found")
	ErrNotEnoughBalance = errors.New("Balance not enough")
)

type Service struct {
	nextAccountID int
	accounts []*types.Account
	payments []*types.Payment
	ID []*types.Account
}


func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone{
			return nil, ErrPhoneregistered
		}
		
	}
	s.nextAccountID++
	account := &types.Account {
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePossitive
	}
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}	
	if account == nil{
		return ErrAccountNotFound
	}
	account.Balance+= amount
	return nil
}

func (s *Service) FindAccountByID(accountID int) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
		
	}
	s.nextAccountID++
	account := &types.Account {
		ID: s.nextAccountID,
	}
	s.accounts = append(s.ID, account)
	return nil, ErrAccountNotFound
}

func (s *Service) Pay(accID int, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <=0 {
		return nil, ErrAmountMustBePossitive
	}
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accID {
			account = acc
			break
		}
	}
		if account == nil{
			return nil, ErrAccountNotFound
		}
		if account.Balance < amount {
			return nil, ErrNotEnoughBalance
		}
		account.Balance -= amount
		paymentID := uuid.New().String()
		payment := &types.Payment {
			ID: paymentID,
			AccountID: accID,
			Amount: amount,
			Category: category,
			Status: types.PaymentStatusInprogress,
		}
		s.payments = append(s.payments, payment)
		return payment, nil
}