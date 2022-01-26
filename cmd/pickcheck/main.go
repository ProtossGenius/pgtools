package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/ProtossGenius/pgtools/impl/pickcheck"
)

func main() {
	mainBranch := flag.String("mb", "main", "main branch")
	pickBranch := flag.String("pb", "", "pick branch")
	after := flag.String("after", "2021-01-24 00:00:00", "git log's begin time, as git log --after.")
	tasks := flag.String("tasks", "", "tasks, split by ','")
	desc := flag.Bool("desc", false, "if desc, sort logs as desc; or asc")
	flag.Parse()

	taskList := strings.Split(*tasks, ",")
	if strings.TrimSpace(*tasks) == "" {
		taskList = make([]string, 0)
	}

	for index, val := range taskList {
		taskList[index] = strings.TrimSpace(val)
	}

	lostLogs := pickcheck.Check(*mainBranch, *pickBranch, *after, taskList)
	fmt.Println("================= lost commits =====================")
	sort.Sort(lostLogs)

	if *desc {
		for i, j := 0, len(lostLogs)-1; i < j; i, j = i+1, j-1 {
			lostLogs[i], lostLogs[j] = lostLogs[j], lostLogs[i]
		}
	}

	pickcheck.ShowLogs(lostLogs)
}
