package wallet

import (
	"io"
	"strings"
	"strconv"
	"fmt"
	"os"
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
	ErrFileNotFound = errors.New("File not found")
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

	var export string

	for _, account := range s.accounts {
		export += fmt.Sprint(account.ID) + ";" + fmt.Sprint(account.Phone) + ";" + fmt.Sprint(account.Balance) + "|"
	}

	_, err = file.WriteString(export)
		if err != nil {
			log.Print(err)
			return err
		}
		log.Printf("%#v", file)
		return err
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)

	if err != nil{
		log.Print(err)
		return ErrFileNotFound
	}
	defer func (){
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil{
			log.Print(err)
			return ErrFileNotFound
		}
		content = append(content, buf[:read]...)
	}
	data := string(content)

	accounts := strings.Split(data, "|")
	accounts = accounts[:len(accounts)-1]
	for _, account := range accounts {
		value := strings.Split(account, ";")
		id, err :=strconv.Atoi(value[0])
		if err != nil {
			return err
		}
		phone := types.Phone(value[1])
		balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		addAccount := &types.Account {
			ID: int(id),
			Phone: phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, addAccount)
		log.Print(account)
	} 
	return nil
}


func (s *Service) Export(dir string) error {
	if len(s.accounts) !=0 {
	filedir1, err := os.Create(dir + "/accounts.dump")
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func(){
		err := filedir1.Close()
		if err != nil{
			log.Print(err)
		}
	}()

	var str1 string

	for _, account := range s.accounts {
		str1 += fmt.Sprint(account.ID) + ";" + fmt.Sprint(account.Phone) + ";" + ";" + fmt.Sprint(account.Balance) + "|"
	}

	_, err = filedir1.WriteString(str1)
		if err != nil {
			return err
		}
	}
		if len (s.payments) != 0 {
		filedir2, err := os.Create("/payments.dump")
		if err != nil {
			return ErrFileNotFound
		}
		defer func(){
			err := filedir2.Close()
			if err != nil {
				log.Print(err)
			}
		}()
		var str2 string

		for _, payment := range s.payments {
			str2 += fmt.Sprint(payment.ID) + ";" + fmt.Sprint(payment.AccountID) + ";" +fmt.Sprint(payment.Amount) + ";" + fmt.Sprint(payment.Category) + ";" + fmt.Sprint(payment.Status) + "|"
		}
		_, err = filedir2.WriteString(str2)
			if err != nil {
				return err
			}
		}
		if len (s.favorites) != 0{
			filedir3, err := os.Create("/favorites.dump")
			if err != nil {
				return ErrFileNotFound
			}
		defer func(){
			err := filedir3.Close()
				if err != nil {
					log.Print(err)
			}
		}()

		var str3 string
		for _, favorite := range s.favorites {
			
		str3 += fmt.Sprint(favorite.ID) + ";" + fmt.Sprint(favorite.AccountID) + ";" + fmt.Sprint(favorite.Name) + fmt.Sprint(favorite.Amount) + ";" + fmt.Sprint(favorite.Category) + "|"
		}

		_, err = filedir3.WriteString(str3)
			if err != nil {
				return err
			}
	}
	return nil
}

func (s *Service)Import(dir string) error {
	fileaccounts, err := os.Open(dir + "/accounts.dump")
	if err != nil{
		log.Print(err)
		//return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound{

	defer func (){
		if cerr := fileaccounts.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	actcontent := make([]byte, 0)
	actbuf := make([]byte, 4)
	for {
		read, err := fileaccounts.Read(actbuf)
		if err == io.EOF {
			break
		}

		if err != nil{
			log.Print(err)
			return ErrFileNotFound
		}
		actcontent = append(actcontent, actbuf[:read]...)
	}
	actdata := string(actcontent)

	accounts := strings.Split(actdata, "|")
	accounts = accounts[:len(accounts)-1]
	for _, account := range accounts {
		value := strings.Split(account, ";")
		id, err :=strconv.Atoi(value[0])
		if err != nil {
			return err
		}
		phone := types.Phone(value[1])
		balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		addAccount := &types.Account {
			ID: 		int(id),
			Phone: 		phone,
			Balance: 	types.Money(balance),
		}

		s.accounts = append(s.accounts, addAccount)
		log.Print(account)
	} 
}
	//return nil

	filepayments, err := os.Open(dir + "/payments.dump")
	if err != nil {
		log.Print(err)
		//return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {

	defer func(){
		if cerr := filepayments.Close(); cerr != nil{
			log.Print(cerr)
		}
	}()

	sqwcontent := make ([]byte,0)
	sqwbuf := make ([]byte,4)
for {
	read, err := filepayments.Read(sqwbuf)
	if err == io.EOF{
		break
	}
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	sqwcontent = append(sqwcontent, sqwbuf[:read]...)
}
sqwdata := string(sqwcontent)

payments := strings.Split(sqwdata, "|")
payments = payments[:len(payments)-1]

for _, payment := range payments {
	value := strings.Split(payment, ";")
	accountID, err :=strconv.Atoi(value[0])
	if err != nil {
		return err
	}
	paymentID := string(value[1])
	paymentAmount, err := strconv.Atoi(value[2])
	if err != nil {
		return err
	}
	paymentCategory := types.PaymentCategory(value[3])
	paymentStatus := types.PaymentStatus(value[4])
	addPayment := &types.Payment {
		ID: 		paymentID,
		AccountID: 	int(accountID),
		Amount: 	types.Money(paymentAmount),
		Category:	types.PaymentCategory(paymentCategory),
		Status:		types.PaymentStatus(paymentStatus),
	}

	s.payments = append(s.payments, addPayment)
	log.Print(payment)
	}
}

filefavorites, err := os.Open(dir + "/favorites.dump")
if err != nil {
	log.Print(err)
	//return ErrFileNotFound
	err = ErrFileNotFound
}
if err != ErrFileNotFound {
defer func (){
	if cerr := filefavorites.Close(); cerr != nil{
		log.Print(cerr)
	} 
}()

fvtcontent := make ([]byte, 0)
fvtbuf := make ([]byte, 4)
for {
	read, err := filefavorites.Read(fvtbuf)
	if err == io.EOF {
		break
	}
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	fvtcontent = append(fvtcontent, fvtbuf[:read]...)
}
fvtdata := string(fvtcontent)

favorites := strings.Split(fvtdata, "|")
favorites = favorites[:len(favorites)-1]
for _, favorite := range favorites {
	v := strings.Split(favorite, ";")
	fvtID := string(v[0])
	fvtactID, err := strconv.Atoi(v[1])
	if err != nil {
		return err
	}
	fvtName := string(v[2])
	fvtAmount, err := strconv.Atoi(v[3])
	if err != nil {
		return err
	}
	fvtCategory := types.PaymentCategory(v[4])

	addFavorite := &types.Favorite{
		ID: 	fvtID,
		AccountID: int(fvtactID),
		Name: 	fvtName,
		Amount:	types.Money(fvtAmount),
		Category: types.PaymentCategory(fvtCategory),
	}
	s.favorites = append(s.favorites, addFavorite)
	log.Print(favorite)
}
	}
	return nil
}