package record_type

import (
	"errors"
)

func DetectRecordType(s string) (string, error) {
	switch s {
	case "ns":
		return s, nil
	case "a":
		return s, nil
	case "aaaa":
		return s, nil
	case "mx":
		return s, nil
	case "cname":
		return s, nil
	case "afsdb":
		return s, nil
	case "apl":
		return s, nil
	case "caa":
		return s, nil
	case "cert":
		return s, nil
	case "dhcid":
		return s, nil
	case "dlv":
		return s, nil
	case "dname":
		return s, nil
	case "dnskey":
		return s, nil
	case "ds":
		return s, nil
	case "hip":
		return s, nil
	case "ipseckey":
		return s, nil
	case "key":
		return s, nil
	case "kx":
		return s, nil
	case "loc":
		return s, nil
	case "naptr":
		return s, nil
	case "nsec":
		return s, nil
	case "nsec3":
		return s, nil
	case "nsec3param":
		return s, nil
	case "ptr":
		return s, nil
	case "rrsig":
		return s, nil
	case "rp":
		return s, nil
	case "sig":
		return s, nil
	case "soa":
		return s, nil
	case "spf":
		return s, nil
	case "srv":
		return s, nil
	case "sshfp":
		return s, nil
	case "ta":
		return s, nil
	case "tkey":
		return s, nil
	case "tlsa":
		return s, nil
	case "tsig":
		return s, nil
	case "txt":
		return s, nil
	case "axfr":
		return s, nil
	case "ixfr":
		return s, nil
	case "opt":
		return s, nil
	default:
		break
	}
	return s, errors.New("Unable to detect record type.")
}
