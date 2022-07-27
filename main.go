package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"github.com/miekg/dns"
)

var records = map[string]string {
	"godns.world.": "45.90.216.199",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
			case dns.TypeA:
				ip := records[q.Name];
				if ip != "" {
					log.Printf(q.Name+": "+ip+"\n");
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip));
					if err == nil {
						m.Answer = append(m.Answer, rr);
					}
				}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg);
	m.SetReply(r);
	m.Compress = false;
	switch r.Opcode {
		case dns.OpcodeQuery:
			parseQuery(m);
	}
	w.WriteMsg(m);
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world");
}

func main() {
	dns.HandleFunc(".", handleDnsRequest);
	server := &dns.Server{Addr: ":53", Net: "udp"};
	log.Printf("Server ready\n");
	err := server.ListenAndServe();
	defer server.Shutdown();
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error());
	}
	http.HandleFunc("/", getRoot);
	weberr := http.ListenAndServe(":3333", nil);
	if weberr != nil {
		log.Fatal("Webserver error: "+weberr.Error());
	}
}