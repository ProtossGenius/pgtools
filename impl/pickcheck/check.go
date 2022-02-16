package pickcheck

import (
	"fmt"
	"strings"
	"time"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_exec"
)

// GitLogInfo git log info detail.
type GitLogInfo struct {
	Author  string
	Date    string
	Title   string
	Tasks   string
	RevCdoe string
	date    int64
}

func (info *GitLogInfo) Show() {
	fmt.Printf("%s %s %s %s\n", info.RevCdoe[:8], info.Title, info.Author, info.Date)
}

func (info *GitLogInfo) LessThan(rhs *GitLogInfo) bool {
	return info.date < rhs.date
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

		info.RevCdoe = line[7:]

		return false
	}

	if strings.HasPrefix(line, "Author") {
		info.Author = line

		return false
	}

	if strings.HasPrefix(line, "Date") {
		info.Date = strings.TrimSpace(line[5:])
		date, err := time.ParseInLocation("Mon Jan 2 15:04:05 2006 -0700", info.Date, time.Local)
		check(err)

		info.date = date.UnixMilli()

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
func Check(mainBranch, pickBranch *string, searchBeginTime string, tasks []string) (lostLogs GitLogInfoArray) {
	pickLogs := GetGitLogInfo(pickBranch, searchBeginTime, tasks)
	mainLogs := GetGitLogInfo(mainBranch, searchBeginTime, tasks)

	return Compare(mainLogs, pickLogs)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GitCheckout(branch string) {
	_, _, err := smn_exec.DirExecGetOut(".", "git", "checkout", branch)
	check(err)
}

func GetGitLogInfo(branch *string, searchBeginTime string, tasks []string) map[string]*GitLogInfo {
	if *branch == "" {
		return make(map[string]*GitLogInfo)
	}

	GitCheckout(*branch)
	*branch = CurrentBranch()

	logDetails := getLogDetails(searchBeginTime)

	result := make(map[string]*GitLogInfo)

	for _, detail := range logDetails {
		if detail.ContainsTask(tasks) {
			result[detail.Date] = detail
		}
	}

	return result
}

func CurrentBranch() string {
	out, _, err := smn_exec.DirExecGetOut(".", "git", "status")
	check(err)

	line := strings.Split(out, "\n")[0]

	const prefix = "On branch "

	if !strings.HasPrefix(line, prefix) {
		check(fmt.Errorf("[%s] not begin with [%s]", line, prefix))
	}

	return line[len(prefix):]
}

func getLogDetails(beginTime string) []*GitLogInfo {
	out, _, err := smn_exec.DirExecGetOut(".", "git", "log", "--after", beginTime)
	check(err)

	result := make([]*GitLogInfo, 0)
	info := new(GitLogInfo)

	for _, line := range strings.Split(out, "\n") {
		if info.Parse(line) {
			result = append(result, info)
			info = new(GitLogInfo)
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

type GitLogInfoArray []*GitLogInfo

// Len is the number of elements in the collection.
func (g GitLogInfoArray) Len() int {
	return len(g)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (g GitLogInfoArray) Less(i int, j int) bool {
	return g[i].LessThan(g[j])
}

// Swap swaps the elements with indexes i andj.
func (g GitLogInfoArray) Swap(i int, j int) {
	g[i], g[j] = g[j], g[i]
}

func ShowLogs(logs GitLogInfoArray) {
	for _, log := range logs {
		log.Show()
	}
}
