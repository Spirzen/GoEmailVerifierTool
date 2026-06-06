package main

import (
	"fmt"
	"net"
	"sort"
)

// lookupMX находит MX-записи домена через DNS и сортирует их по приоритету.
func lookupMX(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("DNS MX lookup: %w", err)
	}
	if len(mxRecords) == 0 {
		return nil, fmt.Errorf("MX-записи для домена %q не найдены", domain)
	}

	sort.Slice(mxRecords, func(i, j int) bool {
		return mxRecords[i].Pref < mxRecords[j].Pref
	})

	return mxRecords, nil
}
