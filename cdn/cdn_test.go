package cdn

import "testing"

func TestFetchCloudflare(t *testing.T) {
	fetcher := CdnIPFetcher{}
	fetcher.FetchCloudFlare()
}
