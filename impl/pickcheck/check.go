package pickcheck

import (
	"fmt"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_exec"
)

// GitLogInfo git log info detail.
type GitLogInfo struct {
	Author  string
	Date    string
	Title   string
	Tasks   string
	RevCdoe string
}

func newInfo() *GitLogInfo {
	return &GitLogInfo{Author: "", Date: "", Title: "", Tasks: "", RevCdoe: ""}
}

func (info *GitLogInfo) ContainsTask(tasks []string) bool {
	if len(tasks) == 0 {
		return true
	}

	for _, infoTask := range strings.Split(info.Tasks, ",") {
		for _, task := range tasks {
			if strings.TrimSpace(infoTask) == task {
				return true
			}
		}
	}

	return false
}

func (info *GitLogInfo) Parse(line string) (parseFinish bool) {
	if strings.HasPrefix(line, "commit") {
		if info.RevCdoe != "" {
			return true
		}

		info.RevCdoe = line

		return false
	}

	if strings.HasPrefix(line, "Author") {
		info.Author = line

		return false
	}

	if strings.HasPrefix(line, "Date") {
		info.Date = line

		return false
	}

	const (
		mainphestTasks = "Maniphest Tasks:"
		splitNum       = 2
	)

	if strings.Contains(line, mainphestTasks) {
		info.Tasks = strings.TrimSpace(strings.SplitN(line, mainphestTasks, splitNum)[1])

		return false
	}

	if info.Title == "" && strings.TrimSpace(line) != "" {
		info.Title = line
	}

	return false
}

// Check .
func Check(mainBranch, pickBranch, searchBeginTime string, tasks []string) (lostLogs []*GitLogInfo) {
	mainLogs := GetGitLogInfo(mainBranch, searchBeginTime, tasks)
	pickLogs := GetGitLogInfo(pickBranch, searchBeginTime, tasks)

	return Compare(mainLogs, pickLogs)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GetGitLogInfo(branch, searchBeginTime string, tasks []string) map[string]*GitLogInfo {
	if branch == "" {
		return make(map[string]*GitLogInfo)
	}

	err := smn_exec.EasyDirExec(".", "git", "checkout", branch)
	check(err)

	logDetails := getLogDetails(searchBeginTime)

	result := make(map[string]*GitLogInfo)

	for _, detail := range logDetails {
		if detail.ContainsTask(tasks) {
			result[detail.Date] = detail
		}
	}

	return result
}

func getLogDetails(beginTime string) []*GitLogInfo {
	out, _, err := smn_exec.DirExecGetOut(".", "git", "log", "--after", beginTime)
	check(err)

	result := make([]*GitLogInfo, 0)
	info := newInfo()

	for _, line := range strings.Split(out, "\n") {
		if info.Parse(line) {
			result = append(result, info)
			info = newInfo()
			info.Parse(line)

			continue
		}
	}

	result = append(result, info)

	return result
}

func Compare(mainLogs, pickLogs map[string]*GitLogInfo) (lostLogs []*GitLogInfo) {
	lostLogs = make([]*GitLogInfo, 0, len(mainLogs))

	for key, info := range mainLogs {
		if _, exist := pickLogs[key]; !exist {
			lostLogs = append(lostLogs, info)
		}
	}

	return lostLogs
}

func ShowLogs(logs []*GitLogInfo) {
	for _, log := range logs {
		fmt.Println("log info detail : ", log)
	}
}
