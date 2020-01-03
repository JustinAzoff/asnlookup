package hostdb

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/josharian/intern"
)

type Record struct {
	IPstart uint32
	IPend   uint32
	Host    string
}

type HostBackend struct {
	data       []Record
	DBFilename string
}

func NewHostBackend(db string) (*HostBackend, error) {
	b := &HostBackend{
		DBFilename: db,
	}
	err := b.reload()
	return b, err
}

func (b *HostBackend) reload() error {

	if _, err := os.Stat("host.cache"); !os.IsNotExist(err) {
		log.Printf("Trying to use cached data")
		file, err := os.Open("host.cache")
		if err != nil {
			log.Printf("Failed to open cache file: %v", err)
			goto rebuild
		}
		defer file.Close()
		d := gob.NewDecoder(file)
		err = d.Decode(&b.data)
		if err != nil {
			log.Printf("failed gob Decode: %v", err)
			b.data = make([]Record, 0)
			goto rebuild
		}
		log.Printf("Using cached data")
		return nil
	}
rebuild:

	log.Printf("Reloading db %s", b.DBFilename)
	file, err := os.Open(b.DBFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	scanner := bufio.NewScanner(gzr)
	var host string
	var ipstarts, ipends string

	var data []Record
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		host = parts[0]
		ipstarts = parts[1]
		ipends = parts[1]

		ipstart, err := strconv.Atoi(ipstarts)
		if err != nil {
			log.Printf("Invalid IP in line: %s", line)
			continue
		}
		ipend, err := strconv.Atoi(ipends)
		if err != nil {
			log.Printf("Invalid IP in line: %s", line)
			continue
		}

		//log.Printf("%s %v %v", host, ipstart, ipend)
		if err != nil {
			return err
		}
		data = append(data, Record{
			IPstart: uint32(ipstart),
			IPend:   uint32(ipend),
			Host:    intern.String(host),
		})
		if len(data)%1000000 == 0 {
			log.Printf("Loaded %d records", len(data))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	log.Printf("Loaded %d records, sorting just in case", len(data))
	sort.Slice(data, func(i, j int) bool { return data[i].IPstart < data[j].IPstart })
	log.Printf("Done")
	b.data = data
	log.Printf("Dumping to cache file")
	cache, err := os.Create("host.cache")
	if err != nil {
		return err
	}
	defer cache.Close()
	e := gob.NewEncoder(cache)
	err = e.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

var NotFound = Record{
	IPstart: 0,
	IPend:   0,
	Host:    "Not Found",
}

func (b *HostBackend) Lookup(ip string) (Record, error) {
	IPNumber := ip2Long(ip)
	index := sort.Search(len(b.data), func(i int) bool { return b.data[i].IPstart >= IPNumber })
	if index >= len(b.data) {
		return NotFound, nil
	}
	rec := b.data[index]
	if IPNumber >= rec.IPstart && IPNumber <= rec.IPend {
		return rec, nil
	}
	return NotFound, nil
}
