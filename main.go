package main

import (
	"flag"
	"fmt"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/httpapi"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "", "监听地址")
	self := flag.Bool("selfcheck", false, "运行有界自检")
	ledgerPath := flag.String("ledger", "ledger.json", "本地 JSON 账本路径")
	flag.Parse()
	listen := httpapi.Address(*addr)
	path := *ledgerPath
	if *self {
		path = ""
	}
	l, err := store.New(path)
	if err != nil {
		log.Fatal(err)
	}
	svc := workflow.New(l)
	if *self {
		if err := svc.SelfCheck(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("selfcheck ok")
		return
	}
	if err := http.ListenAndServe(listen, httpapi.New(svc)); err != nil && !os.IsTimeout(err) {
		log.Fatal(err)
	}
}
