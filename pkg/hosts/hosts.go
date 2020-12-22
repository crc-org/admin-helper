package hosts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goodhosts/hostsfile"
)

type Hosts struct {
	File *hostsfile.Hosts
}

func New() (*Hosts, error) {
	file, err := hostsfile.NewHosts()
	if err != nil {
		return nil, err
	}
	if !file.IsWritable() {
		return nil, fmt.Errorf("host file not writable, try running with elevated privileges")
	}
	return &Hosts{
		File: &file,
	}, nil
}

func (h *Hosts) Add(ip string, hosts []string) error {
	uniqueHosts := map[string]bool{}
	for i := 0; i < len(hosts); i++ {
		uniqueHosts[hosts[i]] = true
	}

	var hostEntries []string
	for key := range uniqueHosts {
		hostEntries = append(hostEntries, key)
	}

	sort.Strings(hostEntries)

	if err := h.File.Add(ip, hostEntries...); err != nil {
		return err
	}
	return h.File.Flush()
}

func (h *Hosts) Remove(hosts []string) error {
	uniqueHosts := map[string]bool{}
	for i := 0; i < len(hosts); i++ {
		uniqueHosts[hosts[i]] = true
	}

	var hostEntries []string
	for key := range uniqueHosts {
		hostEntries = append(hostEntries, key)
	}

	for _, host := range hostEntries {
		if err := h.File.RemoveByHostname(host); err != nil {
			return err
		}
	}
	return h.File.Flush()
}

func (h *Hosts) Clean(rawSuffixes []string) error {
	var suffixes []string
	for _, suffix := range rawSuffixes {
		if !strings.HasPrefix(suffix, ".") {
			return fmt.Errorf("suffix should start with a dot")
		}
		suffixes = append(suffixes, suffix)
	}

	var toDelete []string
	for _, line := range h.File.Lines {
		for _, host := range line.Hosts {
			for _, suffix := range suffixes {
				if strings.HasSuffix(host, suffix) {
					toDelete = append(toDelete, host)
					break
				}
			}
		}
	}

	for _, host := range toDelete {
		if err := h.File.RemoveByHostname(host); err != nil {
			return err
		}
	}
	return h.File.Flush()
}

func (h *Hosts) Contains(ip, host string) bool {
	return h.File.Has(ip, host)
}
