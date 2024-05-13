package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sony/gobreaker"
	"net/http"
	"sud_service/data"
	"sud_service/domain"
	"time"
)

type TuzilastvoClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewTuzilastvoClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) TuzilastvoClient {
	return TuzilastvoClient{
		client:  client,
		address: address,
		cb:      cb,
	}
}

func (ac TuzilastvoClient) DobaviAktivneZahtjeve(ctx context.Context) (data.Zahtevi, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := ac.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ac.address+"/dobaviZahteveZaSudskiPostupak", nil)
		if err != nil {
			return nil, err
		}
		resp, err := ac.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, domain.ErrResp{
				URL:        resp.Request.URL.String(),
				Method:     resp.Request.Method,
				StatusCode: resp.StatusCode,
			}
		}

		var zahtjevi data.Zahtevi

		if err := json.NewDecoder(resp.Body).Decode(&zahtjevi); err != nil {
			return nil, err
		}

		return zahtjevi, nil
	})
	if err != nil {
		return nil, handleHttpReqErr(err, ac.address+"/dobaviZahteveZaSudskiPostupak", http.MethodGet, timeout)
	}

	zahtjevi, ok := cbResp.(data.Zahtevi)
	if !ok {
		return nil, errors.New("invalid response type")
	}

	return zahtjevi, nil
}
