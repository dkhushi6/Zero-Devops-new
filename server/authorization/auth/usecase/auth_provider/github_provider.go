package AuthProvider

import (
	"context"
	"golang.org/x/sync/errgroup"
	"golang.org/x/oauth2"
	"golang.org/x/oauth/github"
)

type Github interface {
	
}

func Github(ctx context.Context){
	g , ctx = errgroup.WithContext(ctx)
	g.Go(func() error {
		
		return nil
	})
}