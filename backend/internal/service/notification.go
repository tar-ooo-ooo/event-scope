package service

import (
	"errors"
	"math/rand"
)

func SendNotification() error {
	// 設定 95% 的機率成功，5% 的機率失敗
	if rand.Intn(100) < 95 {
		return nil
	}

	return errors.New("推播發送失敗！")
}
