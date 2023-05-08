package webapi

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

type ConfirmAPI struct {
	t       string
	storage map[string]string
	mu      sync.RWMutex
}

func New(token string) *ConfirmAPI {
	return &ConfirmAPI{
		t:       token,
		storage: make(map[string]string),
	}
}

func (v *ConfirmAPI) GenerateCode(number string) (string, error) {
	code := strconv.Itoa(rangeIn(1000, 9999))

	v.mu.Lock()
	defer v.mu.Unlock()
	v.storage[number] = code

	return code, nil
}

func (v *ConfirmAPI) ConfirmNumber(number, code string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	c, ok := v.storage[number]
	if !ok {
		return fmt.Errorf("number not found: %s", number)
	}

	if c != code {
		return fmt.Errorf("confirmation code wrong")
	}

	delete(v.storage, number)

	return nil
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
