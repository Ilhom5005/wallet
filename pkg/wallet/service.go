package wallet

import (
	"fmt"
	"os"
	"strconv"
	"log"
	 "github.com/Ilhom5005/wallet/v1/pkg/types"
	 "github.com/google/uuid"
	 "errors"
)



var(
	ErrPhoneregistered = errors.New("phone alredy registered")
	ErrAmountMustBePossitive = errors.New("amount must be greater then zero")
	ErrAccountNotFound = errors.New("account not found")
	ErrNotEnoughBalance = errors.New("Balance not enough")
	ErrPaymentNotFound = errors.New("Payment not Found")
	ErrFavoriteNotFound = errors.New("Favorite not Found")
)

type Service struct {
	nextAccountID int
	accounts []*types.Account
	payments []*types.Payment
	ID []*types.Account
	favorites []*types.Favorite
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

func (s *Service) Reject(paymentID string) error {
	var targetPayment *types.Payment

	for _, payment := range s.payments {
		if payment.ID == paymentID{
			targetPayment = payment
			break
		}
		
	}

	if targetPayment == nil{
		return ErrPaymentNotFound
	}

	var targetAccount *types.Account
	for _, account := range s.accounts {
		if account.ID == targetPayment.AccountID{
			targetAccount = account
			break
		}
		
	}
	if targetAccount == nil {
		return ErrAccountNotFound
	}
	targetPayment.Status = types.PaymentStatusFail
	targetAccount.Balance += targetPayment.Amount
	return nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID{
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}


func (s *Service) Repeat(paymentID string)(*types.Payment, error)  {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment1, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return payment1, nil
}




func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID{
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}


func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	
	if err != nil {
		return nil, err
	}
	favoriteID := uuid.New().String()
	newFavorite := &types.Favorite{
		ID: 		favoriteID,
		AccountID: 	payment.AccountID,
		Name: 		name,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
	}
	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}


func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) ExportToFile(path string) error {

	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func(){
		if cerr:=file.Close(); cerr != nil{
			log.Print(err)
		}
	}()

	var account *types.Account
	var export string

	for _, account := range s.accounts {
		export += fmt.Sprint(account.ID) + ";" + fmt.Sprint(account.Phone) + ";" + fmt.Sprint(account.Balance) + "|"
	}

	_, err = file.WriteString(strconv.FormatInt(int64(account.ID), 10))
		if err != nil {
			log.Print(err)
			return err
		}
		log.Printf("%#v", file)
		return err
}