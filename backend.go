package asnlookup

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/chromicant/go-iptree"
)

type JSONOwnerMapping map[string]string

type OwnerInfo struct {
	Owner string
	CC    string
}

type OwnerMapping map[string]OwnerInfo

type AsnBackend struct {
	iptree        *iptree.IPTree
	nameMapping   OwnerMapping
	DBFilename    string
	NamesFilename string
}

type Record struct {
	IP    string
	AS    string
	Owner string
	CC    string
}

func NewAsnBackend(db, names string) (*AsnBackend, error) {
	b := &AsnBackend{
		DBFilename:    db,
		NamesFilename: names,
	}
	err := b.reload()
	return b, err
}

func (b *AsnBackend) reload() error {
	err := b.reloadDB()
	if err != nil {
		return err
	}
	return b.reloadNames()
}
func (b *AsnBackend) reloadDB() error {
	t := iptree.New()
	log.Printf("Reloading ASN db %s", b.DBFilename)
	file, err := os.Open(b.DBFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") {
			continue
		}
		if strings.Contains(line, ":") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		t.AddByString(parts[0], parts[1])
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	b.iptree = t
	return nil
}

func (b *AsnBackend) reloadNames() error {
	log.Printf("Reloading ASN owner db %s", b.NamesFilename)
	file, err := os.Open(b.NamesFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	var jsinfo JSONOwnerMapping
	info := make(OwnerMapping)

	dec := json.NewDecoder(file)
	if err = dec.Decode(&jsinfo); err != nil {
		return err
	}
	for as, ownerString := range jsinfo {
		if !strings.Contains(ownerString, ",") {
			if len(ownerString) == 2 {
				info[as] = OwnerInfo{CC: ownerString}
			} else {
				info[as] = OwnerInfo{Owner: ownerString}
			}
		} else {
			parts := strings.Split(ownerString, ",")
			info[as] = OwnerInfo{
				Owner: strings.TrimSpace(parts[0]),
				CC:    strings.TrimSpace(parts[1]),
			}
		}
	}

	b.nameMapping = info
	return nil
}

func (b *AsnBackend) Lookup(ip string) (Record, error) {
	var rec Record
	val, found, err := b.iptree.GetByString(ip)
	if err != nil {
		return rec, err
	}
	rec.IP = ip
	if !found {
		return rec, nil
	}
	rec.AS = val.(string)

	info, existed := b.nameMapping[rec.AS]
	if existed {
		rec.Owner = info.Owner
		rec.CC = info.CC
	}

	return rec, nil
}
