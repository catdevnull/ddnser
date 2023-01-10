package nameservers

type NameServer interface {
	SetRecord(domain string, overrideIp string) (string, error)
}
