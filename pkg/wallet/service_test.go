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

func TestService_Favorite_success_user(t *testing.T) {
	svc := Service{}

	account, err := svc.RegisterAccount("+992938151007")
	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Pay() Error() can't pay for an account(%v): %v", account, err)
	}

	favorite, err := svc.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Errorf("FavoritePayment() Error() can't for an favorite(%v): %v", favorite, err)
	}

	paymentFavorite, err := svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite() Error() can't for an favorite(%v): %v", paymentFavorite, err)
	}
}

func TestService_Export_success_user(t *testing.T) {
	var svc Service
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
  
	err := svc.ExportToFile("export.txt")
	if err != nil {
	  t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
  
  }
  
func TestService_Import_success(t *testing.T) {
	s := newTestService()

	err := s.Import("data")

	if err != nil {
		t.Errorf("Import() Error can't import error = %v", err)
	}
}

func TestService_ExportImport_success_user(t *testing.T) {
	var svc Service

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.RegisterAccount("+992000000004")
	
	err := svc.Export(".")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

	err = svc.Import(".")
	
	if err != nil {
		t.Errorf("method ImportToFile returned not nil error, err => %v", err)
	}

}

func TestService_HistoryToFiles_success(t *testing.T) {
	s := newTestService()

	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, _, err = s.addAccount(testAccount{
		phone:   "+992935444994",
		balance: 10_000_00,
		payments: []struct {
			amount   types.Money
			category types.PaymentCategory
		}{{
			amount:   1000_00,
			category: "auto",
		}, {
			amount:   1020_00,
			category: "auto",
		},
			{
				amount:   1030_00,
				category: "auto",
			}},
	})

	if err != nil {
		t.Error(err)
		return
	}

	payments, err := s.ExportAccountHistory(2)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.HistoryToFiles(payments, "data", 2)
	if err != nil {
		t.Errorf("HistoryToFiles() Error can't export to file, error = %v", err)
		return
	}
}


func BenchmarkSumPayments_user(b *testing.B) {
	var svc Service  
	want:= types.Money(0)
	for i:=0 ; i < b.N ; i++ {
		result := svc.SumPayments(2)
		if result != types.Money(want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)

		}
	}
}

