package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/retry"
)

// ProcessOrder processes an order by its number.
func (ac *AccrualClient) ProcessOrder(ctx context.Context, orderNumber string) (*AccrualInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ac.cfg.BaseURL+"/api/orders/"+orderNumber, nil)
	if err != nil {
		return nil, errs.Wrap(err, "prepare request for accrual service")
	}

	buff := &bytes.Buffer{}
	var statusCode int
	err = retry.WithLinearBackoffRetry(ctx, &retry.Config{}, func(context.Context) error {
		resp, err := ac.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		statusCode = resp.StatusCode
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("got status code: %d", resp.StatusCode)
		} else if resp.StatusCode >= http.StatusBadRequest {
			return errs.Wrap(retry.ErrUnretriable, fmt.Sprintf("got status code: %d", resp.StatusCode))
		}

		buff.Reset()
		_, err = io.Copy(buff, resp.Body)
		if err != nil {
			return errs.Wrap(retry.ErrUnretriable, err.Error())
		}

		return nil
	})
	if err != nil {
		return nil, errs.Wrap(err, "process order")
	}

	log.Debug().Str("orderNumber", orderNumber).
		Int("statusCode", statusCode).
		Msg("accrual service response")

	// order is not in the accrual system
	if statusCode == http.StatusNoContent {
		return nil, nil
	}

	var accrualInfo AccrualInfo
	err = json.Unmarshal(buff.Bytes(), &accrualInfo)
	if err != nil {
		return nil, errs.Wrap(err, "unmarshal accrual info")
	}

	return &accrualInfo, nil
}
