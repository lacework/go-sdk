package gcp

import (
	"context"

	"google.golang.org/api/oauth2/v2"
)

type Caller struct {
	Email  string
	UserID string
}

func FetchCaller(p *Preflight) error {
	oauth2Svc, err := oauth2.NewService(context.Background(), p.gcpClientOption)
	if err != nil {
		return err
	}

	tokenInfo, err := oauth2Svc.Tokeninfo().Do()
	if err != nil {
		return err
	}

	p.caller = Caller{
		Email:  tokenInfo.Email,
		UserID: tokenInfo.UserId,
	}

	return nil
}
