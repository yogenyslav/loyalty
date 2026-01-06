// Package main provides utils to add orders to accrual system.
// It is used for testing purposes only.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
)

var accrualAddr string

// Good represents a single good in an order.
type Good struct {
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

// Order represents an order to be added to the accrual system.
type Order struct {
	Number string `json:"order"`
	Goods  []Good `json:"goods"`
}

// RewardType represents the type of reward.
type RewardType string

// Possible reward types.
const (
	RewardTypePercent RewardType = "%"
	RewardTypePoints  RewardType = "pt"
)

// Reward represents a reward structure.
type Reward struct {
	Match      string  `json:"match"`
	Reward     float64 `json:"reward"`
	RewardType string  `json:"reward_type"`
}

func addReward(ctx context.Context, reward Reward) error {
	rewardRaw, err := json.Marshal(reward)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accrualAddr+"/api/goods", bytes.NewReader(rewardRaw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add reward, status code: %d", resp.StatusCode)
	}

	return nil
}

func addOrder(ctx context.Context, order Order) error {
	orderRaw, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accrualAddr+"/api/orders", bytes.NewReader(orderRaw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to add order, status code: %d", resp.StatusCode)
	}

	return nil
}

func init() {
	accrualAddrEnv, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if !ok {
		accrualAddrEnv = "http://localhost:8080"
	}
	accrualAddr = accrualAddrEnv
}

func main() {
	flags := flag.NewFlagSet("accrual", flag.ExitOnError)
	order := flags.String("o", "2844830162", "номер заказа, который будет добавлен в систему расчета бонусов")
	orderDescription := flags.String("d", "Bork", "описание товара в заказе")
	reward := flags.Float64("r", 10, "начисление бонусов за заказ")
	rewardMatch := flags.String("m", "Bork", "шаблон для совпадения заказа при начислении бонусов")
	rewardType := flags.String("t", "%", "тип начисления бонусов: pt - в баллах, % - в процентах")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err := addOrder(ctx, Order{
		Number: *order,
		Goods: []Good{
			{
				Description: *orderDescription,
				Price:       rand.Int64N(1000),
			},
		},
	})
	if err != nil {
		log.Printf("failed to add order: %v", err)
	}

	err = addReward(ctx, Reward{
		Match:      *rewardMatch,
		Reward:     *reward,
		RewardType: *rewardType,
	})
	if err != nil {
		log.Printf("failed to add reward: %v", err)
	}
}
