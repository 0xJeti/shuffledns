package massdns

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/shuffledns/internal/store"
	"github.com/projectdiscovery/shuffledns/pkg/parser"
	"github.com/rs/xid"
)

// Process runs the actual enumeration process returning a file
func (c *Client) Process() error {
	// Check for blank input file or non-existent input file
	blank, err := IsBlankFile(c.config.InputFile)
	if err != nil {
		return err
	}
	if blank {
		return errors.New("blank input file specified")
	}

	// Create a store for storing ip metadata
	shstore := store.New()
	defer shstore.Close()

	// Create a temporary file for the massdns output
	temporaryOutput := path.Join(c.config.TempDir, xid.New().String())

	gologger.Infof("Creating temporary massdns output file: %s\n", temporaryOutput)
	gologger.Infof("Executing massdns on %s\n", c.config.Domain)

	now := time.Now()
	// Run the command on a temp file and wait for the output
	cmd := exec.Command(c.config.MassdnsPath, []string{"-r", c.config.ResolversFile, "-t", "A", c.config.InputFile, "-w", temporaryOutput, "-s", strconv.Itoa(c.config.Threads)}...)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not execute massdns: %w", err)
	}
	gologger.Infof("Massdns execution took %s\n", time.Now().Sub(now))

	massdnsOutput, err := os.Open(temporaryOutput)
	if err != nil {
		return fmt.Errorf("could not open massdns output file: %w", err)
	}
	defer massdnsOutput.Close()

	gologger.Infof("Parsing output and removing wildcards\n")

	// at first we need the full structure in memory to elaborate it in parallell
	err = parser.Parse(massdnsOutput, func(domain string, ip []string) {
		for _, ip := range ip {
			// Check if ip exists in the store. If not,
			// add the ip to the map and continue with the next ip.
			if !shstore.Exists(ip) {
				shstore.New(ip, domain)
				continue
			}

			// Get the IP meta-information from the store.
			record := shstore.Get(ip)

			// Put the new hostname and increment the counter by 1.
			record.Hostnames[domain] = struct{}{}
			record.Counter++
		}
	})

	if err != nil {
		return fmt.Errorf("could not parse massdns output: %w", err)
	}

	// start to works in parallel on wildcards
	var (
		wildcardWG sync.WaitGroup
	)
	workchan := make(chan *store.IPMeta)
	for i := 0; i < c.config.WildcardsThreads; i++ {
		wildcardWG.Add(1)
		go func() {
			defer wildcardWG.Done()
			for ipm := range workchan {
				// We've stumbled upon a wildcard, just ignore it.
				if _, ok := c.wildcardIPMap[ipm.IP]; ok {
					continue
				}

				// // If the same ip has been found more than 5 times, perform wildcard detection
				// // on it now, if an IP is found in the wildcard we add it to the wildcard map
				// // so that further runs don't require such filtering again.
				if ipm.Counter >= 5 && !ipm.Validated {
					for host := range ipm.Hostnames {
						wildcard, ips := c.wildcardResolver.LookupHost(host)
						if wildcard {
							for ip := range ips {
								c.wildcardIPMap[ip] = struct{}{}
							}

							continue
						}
						ipm.Validated = true
					}
				}
			}
		}()
	}

	// process all the items
	wildcardWG.Add(1)
	go func() {
		defer close(workchan)
		defer wildcardWG.Done()
		for _, record := range shstore.IP {
			workchan <- record
		}
	}()

	wildcardWG.Wait()

	// drop all wildcard from the store
	for wildcardIP := range c.wildcardIPMap {
		shstore.Delete(wildcardIP)
	}

	gologger.Infof("Finished enumeration, started writing output\n")

	// Parse the massdns output
	return c.writeOutput(shstore)
}

func (c *Client) writeOutput(store *store.Store) error {
	// Write the unique deduplicated output to the file or stdout
	// depending on what the user has asked.
	var output *os.File
	var w *bufio.Writer
	var err error

	if c.config.OutputFile != "" {
		output, err = os.Create(c.config.OutputFile)
		if err != nil {
			return fmt.Errorf("could not create massdns output file: %v", err)
		}
		w = bufio.NewWriter(output)
	}
	buffer := &strings.Builder{}

	uniqueMap := make(map[string]struct{})

	for _, record := range store.IP {
		for hostname := range record.Hostnames {
			// Skip if we already printed this subdomain once
			if _, ok := uniqueMap[hostname]; ok {
				continue
			}
			uniqueMap[hostname] = struct{}{}

			buffer.WriteString(hostname)
			buffer.WriteString("\n")
			data := buffer.String()

			if output != nil {
				w.WriteString(data)
			}
			gologger.Silentf("%s", data)
			buffer.Reset()
		}
	}

	// Close the files and return
	if output != nil {
		w.Flush()
		output.Close()
	}
	return nil
}
