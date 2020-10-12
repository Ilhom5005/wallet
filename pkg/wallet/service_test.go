package wallet

import (
	"fmt"
	"reflect"
	"github.com/Ilhom5005/wallet/v1/pkg/types"
	"github.com/google/uuid"
	"testing"
)

func TestService_Reject_success(t *testing.T) {
	
	s := newTestService()
	
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]

	 err = s.Reject(payment.ID)
	 if err != nil {
		t.Errorf("Reject(): can't register account, error = %v",err)
		return
	 }
}

func TestService_FindPaymentByID_success(t *testing.T) {

	s := newTestService()
	
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	 got, err := s.FindPaymentByID(payment.ID)
	 if err != nil {
		t.Errorf("FindPaymentByID(): can't found Payment by ID, error = %v",err)
		return
	 }

	 if !reflect.DeepEqual(payment, got){
		t.Errorf("FindPaymentByID(): wrong payment returned, error = %v",err)
		return
	 }

}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	 _, err = s.FindPaymentByID(uuid.New().String())
	 if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil")
		return
	 }

	 if err != ErrPaymentNotFound {
		 t.Errorf("FindPaymentByID(): must return Errpayment, error = %v",err)
		 return
	 }
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func  (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error){
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, err: %v", err)
	}
	s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit to account, err: %v", err)

	}
	return account, nil
}

type testAccount struct{
	phone types.Phone
	balance types.Money
	payments []struct{
		amount types.Money
		category types.PaymentCategory
	}
}

func (s *Service) addAccount(data testAccount)  (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, err: %v", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, err: %v", err)
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil{
			return nil, nil, fmt.Errorf("can't make payment, err: %v", err)
		}
	}
	return account, payments, nil
}

var defaultTestAccount = testAccount{
	phone: "+9923333333",
	balance: 100000000,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 10000, category: "auto"},
	},
}

func TestService_Reapet_success(t *testing.T){

s := newTestService()
acc, err := s.RegisterAccount("+987360003")
if err != nil {
	t.Errorf("RegisterAccountUser: cannot register account, error = %v", err)
	return 
}
err = s.Deposit(acc.ID, 100)
if err != nil {
	t.Errorf("can not deposit account, error = %v", err)
	return
}
payment, err := s.Pay(acc.ID, 10, "ice-cream")
if err != nil {
	t.Errorf("can not pay, error = %v", err)
	return
}
payment1, err := s.FindPaymentByID(payment.ID)
if err != nil {
	t.Errorf("method FindAccountByID returned not nil error, payment => %v", payment)
	return
}
payment1, err = s.FindPaymentByID(payment.ID)
if err != nil {
	t.Errorf("can not repeat payment, error = %v", err)
	return
}
if payment.Amount != payment1.Amount || payment.Category != payment1.Category {
	t.Error("wrong result")
}
}

func TestService_FindAccountByID_success_user(t *testing.T){
	s := newTestService()

	s.RegisterAccount("+918616330")
	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("method FindPaymentByID retuned not nil error, payment => %v", account)
		return
	}
}

func TestService_FindAccountByID_notFound_user(t *testing.T){
	s := newTestService()

	s.RegisterAccount("+918616330")
	account, err := s.FindAccountByID(2)

	if err == nil {
		t.Errorf("method FindPaymentByID returned nil error, payment => %v", account)
		return
	}
}