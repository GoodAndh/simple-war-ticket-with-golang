package voucherweb

import (
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Voucher struct {
	Code  int
	Max   int
	Mutex sync.RWMutex
}

type UserVoucher struct {
	Code  int
	Total int
}

type User struct {
	ID      int
	Name    string
	Voucher []UserVoucher
}

func (vcr *Voucher) decreaseTicket(max int) {
	vcr.Max -= max
}

func TakeVoucher(seller *Voucher, amount int, buyer *User, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	seller.Mutex.Lock()

	if amount > seller.Max || seller.Max-amount <= 0 {
		log.Println("kamu gagal mengambil voucher ,sisa voucher : ", seller.Max, " jumlah voucher yang kamu minta : ", amount)
		seller.Mutex.Unlock()
		return
	}
	log.Println("voucher count : ", seller.Max)
	seller.decreaseTicket(amount)
	seller.Mutex.Unlock()

	log.Println("remaining voucher inside :", seller.Max)

	seller.Mutex.Lock()
	voucherItems := &UserVoucher{
		Code:  seller.Code,
		Total: amount,
	}
	buyer.Voucher = append(buyer.Voucher, *voucherItems)
	log.Println("userNew Inside :", buyer)
	seller.Mutex.Unlock()

}
func TestMutexTf(t *testing.T) {
	wg := &sync.WaitGroup{}
	seller := &Voucher{
		Code: 908,
		Max:  10,
	}
	userNew := generatoBuyer(5)

	for _, buyer := range userNew {
		go func(user *User) {
			TakeVoucher(seller, 3, user, wg)
		}(buyer)
	}
	wg.Wait()

	time.Sleep(3 * time.Second)

	log.Println("remaining voucher  outside:", seller.Max)

}

func generatoBuyer(n int) []*User {
	buyerSlice := make([]*User, 0, n)
	for i := 1; i <= n; i++ {
		idConvert := strconv.Itoa(i)
		buyer := &User{
			ID:   i,
			Name: "Andika " + idConvert,
		}
		buyerSlice = append(buyerSlice, buyer)
	}
	return buyerSlice
}

// func checkIfAlrdExist(user *User, voucher *Voucher) bool {
// 	if len(user.Voucher) <= 0 {
// 		return false
// 	}

// 	for _, vcr := range user.Voucher {
// 		if vcr.Code == voucher.Code {
// 			return true
// 		}
// 	}

// 	return false
// }
