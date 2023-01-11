package nameservers

import "context"

type NameServer interface {
	SetRecord(ctx context.Context, domain string, overrideIp string) (string, error)
}
