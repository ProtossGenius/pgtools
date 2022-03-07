package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/ProtossGenius/pgtools/impl/pickcheck"
)

func main() {
	mainBranch := flag.String("mb", "-", "main branch")
	pickBranch := flag.String("pb", "-", "pick branch")
	after := flag.String("after", "2021-01-24 00:00:00", "git log's begin time, as git log --after.")
	desc := flag.Bool("desc", false, "if desc, sort logs as desc; or asc")
	flag.Parse()

	lostLogs := pickcheck.Check(pickBranch, mainBranch, *after, nil)
	fmt.Println("基于", *mainBranch, "\nbuild", *pickBranch, "\npick\n")
	sort.Sort(lostLogs)

	if *desc {
		for i, j := 0, len(lostLogs)-1; i < j; i, j = i+1, j-1 {
			lostLogs[i], lostLogs[j] = lostLogs[j], lostLogs[i]
		}
	}

	pickcheck.ShowLogs(lostLogs)
}
