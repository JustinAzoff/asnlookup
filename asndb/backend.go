package asndb

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/chromicant/go-iptree"
	"github.com/josharian/intern"
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
	AS     int
	Prefix string
	IP     string
	Owner  string
	CC     string
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
	err := b.reloadNames()
	if err != nil {
		return err
	}
	return b.reloadDB()
}
func (b *AsnBackend) reloadDB() error {
	t := iptree.New()
	log.Printf("Reloading AS db %s", b.DBFilename)
	file, err := os.Open(b.DBFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var prefix, as string
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
		prefix = parts[0]
		as = parts[1]

		asint, err := strconv.Atoi(as)
		if err != nil {
			log.Printf("Invalid AS in line: %s", line)
			continue
		}

		rec := Record{
			AS:     asint,
			Prefix: prefix,
		}
		info, existed := b.nameMapping[as]
		if existed {
			rec.Owner = intern.String(info.Owner)
			rec.CC = intern.String(info.CC)
		}
		t.AddByString(prefix, rec)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	b.iptree = t
	return nil
}

func (b *AsnBackend) reloadNames() error {
	log.Printf("Reloading AS owner db %s", b.NamesFilename)
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
	if !found {
		return Record{IP: ip}, nil
	}
	rec = val.(Record)
	rec.IP = ip

	return rec, nil
}
