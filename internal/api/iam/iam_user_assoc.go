// Code generated by api/iam/iam_gen.go. DO NOT EDIT.

package iam

import (
	"context"
	"runtime"
	"slices"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func ListUserAssociationInfo(ctx context.Context, client IIamClient, ids, names []string, document bool, filters []string) ([]UserAssociationInfo, error) {
	pols, err := client.FetchCustomerPolicies(ctx)
	if err != nil {
		return nil, err
	}
	eg, ctx := errgroup.WithContext(ctx)
	l := rate.NewLimiter(rate.Limit(50), 1)
	ich := make(chan UserAssociationInfo, runtime.NumCPU())
	var info []UserAssociationInfo
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range ich {
			info = append(info, i)
		}
	}()
	p := iam.NewListUsersPaginator(client, &iam.ListUsersInput{})
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.Users {
			item := item
			if len(ids) > 0 && !slices.Contains(ids, aws.ToString(item.UserId)) {
				continue
			}
			if len(names) > 0 && !slices.Contains(names, aws.ToString(item.UserName)) {
				continue
			}
			eg.Go(func() error {
				if err := l.Wait(ctx); err != nil {
					return err
				}
				if err := GetUserAssociationInfo(ctx, l, client, ich, item, document, filters, pols); err != nil {
					return err
				}
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		close(ich)
		return nil, err
	}
	close(ich)
	wg.Wait()
	return info, nil
}
