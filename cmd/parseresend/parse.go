package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
)

var mapper = map[string]string{
	"I": "初始化",
	"S": "接通",
	"U": "(强)用户拒接",
	"D": "(强)未接通",
	"R": "(强)振铃",
	"B": "(强)用户忙",
	"F": "拨打失败",
	"O": "关机",
	"P": "停机",
	"N": "空号",
	"E": "不在服务区",
	"C": "欠费未接听",
	"J": "运营商拒绝接听",
	"T": "其他原因",
	"A": "取消拨打",
	"G": "推送失败",
	"V": "语音信箱",
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// IvrSendInfo ivr send info.
type IvrSendInfo struct {
	CaseID      int
	StrategyID  int
	Provider    string
	Result      string
	TimeCreated int64
	TimeUpdated int64
	ResendTimes int
}

// toInt string to int.
func toInt(str string) int {
	if v, err := strconv.Atoi(str); err == nil {
		return v
	}

	return 0
}

// toInt64 ot uint64.
func toInt64(str string) int64 {
	const (
		base    = 10
		bitSize = 64
	)

	if v, err := strconv.ParseInt(str, base, bitSize); err == nil {
		return v
	}

	return 0
}

func (info *IvrSendInfo) Parse(lst []string) {
	info.CaseID = toInt(lst[0])
	info.StrategyID = toInt(lst[1])
	info.Provider = lst[2]
	info.Result = lst[3]
	info.TimeCreated = toInt64(lst[4])
	info.TimeUpdated = toInt64(lst[5])
	//	resendTimes := strings.Split(lst[6], `RESEND_TIMES"":""`)[1]
	//	resendTimes = strings.Split(resendTimes, `""`)[0]
	//	info.ResendTimes = toInt(resendTimes)
}

// UserInfo user info.
type UserInfo struct {
	CaseID      int
	StrategyID  int
	Infos       []*IvrSendInfo
	ResendTimes int
	DataTime    int64
}

// CanAdd can add.
func (user *UserInfo) CanAdd(sendInfo *IvrSendInfo) bool {
	return len(user.Infos) == 0 || user.CaseID == sendInfo.CaseID
}

func (user *UserInfo) Add(sendInfo *IvrSendInfo) {
	user.CaseID = sendInfo.CaseID
	user.StrategyID = sendInfo.StrategyID
	user.Infos = append(user.Infos, sendInfo)
	user.ResendTimes++
	user.DataTime = sendInfo.TimeCreated
}

// newUserInfo new user info.
func newUserInfo() *UserInfo {
	const maxResendTimes = 5

	return &UserInfo{
		CaseID:      0,
		StrategyID:  0,
		Infos:       make([]*IvrSendInfo, 0, maxResendTimes),
		ResendTimes: 0,
		DataTime:    0,
	}
}

func main() {
	data, err := smn_file.FileReadAll("./ivr_data.txt")
	check(err)

	const fieldNums = 8

	strList := strings.Split(string(data), "\n")

	userInfoList := make([]*UserInfo, 0, len(strList))

	userInfo := newUserInfo()

	for _, line := range strList {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lineSpl := strings.SplitN(line, "\t", fieldNums)
		for i := range lineSpl {
			lineSpl[i] = strings.TrimSpace(lineSpl[i])
		}

		info := new(IvrSendInfo)
		info.Parse(lineSpl)

		if !userInfo.CanAdd(info) {
			userInfoList = append(userInfoList, userInfo)

			userInfo = newUserInfo()
		}

		userInfo.Add(info)
	}

	if len(userInfo.Infos) != 0 {
		userInfoList = append(userInfoList, userInfo)
	}

	ResendTrans(userInfoList)
}

const Time20220221 = 1645459200000

// ResendTrans n times resend trans to antoher.
func ResendTrans(userInfoList []*UserInfo) {
	sum := 0
	beforeToSucSum := 0
	beforeToSucMap := map[string]int{}

	afterToSucSum := 0
	afterToSucMap := map[string]int{}

	for _, user := range userInfoList {
		for idx, sendInfo := range user.Infos {
			if idx > 1 {
				sum++

				if sendInfo.Result == "S" {
					if user.DataTime < Time20220221 {
						beforeToSucSum++
						beforeToSucMap[user.Infos[idx-1].Result]++
					} else {
						afterToSucSum++
						afterToSucMap[user.Infos[idx-1].Result]++
					}
				}
			}
		}
	}

	for k, v := range beforeToSucMap {
		fmt.Println("before ToSuc trans : ", k, " percent = ", float64(v)/float64(beforeToSucSum))
	}
	fmt.Println("--------------------------")
	for k, v := range afterToSucMap {
		fmt.Println("after ToSuc trans : ", k, " percent = ", float64(v)/float64(afterToSucSum))
	}
}
