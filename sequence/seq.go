package sequence

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 暂支持的时间编号格式
var timeType = map[string]int{
	"yyMMdd":   1,
	"yyyyMMdd": 2,
}

type SeqInfo struct {
	Module string            `json:"module"` // 模块名称
	Expr   string            `json:"expr"`   // 表达式
	Group  map[string]string `json:"group"`  // 表达式对应 需要维护的map
	Remark string            `json:"remark"` // 备注信息
}

// ISeqData 数据库操作接口
type ISeqData interface {
	Data() []*SeqInfo
	Save(model *SeqInfo)
}

type Sequence struct {
	ISeqData
	Modules map[string]*SeqInfo
}

// New  创建数据
func New(src ISeqData) *Sequence {
	seq := &Sequence{
		ISeqData: src,
	}
	seq.Modules = map[string]*SeqInfo{}
	dataSource := src.Data()
	for _, v := range dataSource {
		seq.Modules[v.Module] = v
	}
	return seq
}

// Gen 生成对应表达式的编号
func (s *Sequence) Gen(module string, args ...string) string {
	ret := ""
	model, ok := s.Modules[module]
	if !ok {
		return ""
	}
	arr := strings.Split(model.Expr, "-")
	tType := ""
	pNum := 0
	for _, v := range arr {
		if v[0:1] == "p" {
			pNum++
		}
	}
	if len(args) != pNum {
		return ""
	}
	for _, v := range arr {
		if v[0:1] == "i" {
			ret += s.increment(v, model.Group, tType, args...)
		} else {
			rets, tTyped := s.buildStr(v, args...)
			ret += rets
			if tTyped != "" {
				tType = tTyped
			}
		}

	}
	s.Save(model)
	return ret
}

// increment 生成末尾的编号部分
func (s *Sequence) increment(v string, group map[string]string, tType string, args ...string) string {
	arr := strings.Split(v, ":")
	length, _ := strconv.ParseInt(arr[2], 10, 0)
	key := ""
	if arr[1] == "d" {
		timeStr := ""
		if tType2, ok := timeType[tType]; ok {
			if tType2 == 1 {
				timeStr = strings.ReplaceAll(time.Now().String()[2:10], "-", "")
			} else {
				timeStr = strings.ReplaceAll(time.Now().String()[:10], "-", "")
			}
		}
		key = timeStr
	} else {
		index, _ := strconv.ParseInt(arr[1], 10, 0)
		key = args[index-1]
	}
	num, ok := group[key]
	if !ok {
		for i := 0; i < int(length-1); i++ {
			num += "0"
		}
		num += "1"
		group[key] = num
		return num
	}
	n, _ := strconv.ParseInt(num, 10, 0)
	str := fmt.Sprintf("%d", n+1)
	newStr := ""
	for i := 0; i < (int(length) - len(str)); i++ {
		newStr += "0"
	}
	str = newStr + str
	group[key] = str
	return str
}

// buildStr 生成类型编号部分和日期编号部分
func (s *Sequence) buildStr(str string, args ...string) (string, string) {
	if len(str) == 0 {
		return "", ""
	}
	arr := strings.Split(str, ":")
	if arr[0] == "p" {
		if index, err := strconv.ParseInt(arr[1], 10, 0); err != nil {
			return "", ""
		} else {
			return args[index-1], ""
		}
	} else if arr[0] == "d" {
		timeStr := ""
		if tType, ok := timeType[arr[1]]; ok {
			if tType == 1 {
				return strings.ReplaceAll(time.Now().String()[2:10], "-", ""), arr[1]
			} else {
				return strings.ReplaceAll(time.Now().String()[:10], "-", ""), arr[1]
			}
		}
		return timeStr, arr[1]
	} else {
		return "", ""
	}
}

// IfAccordNumberRule 判断表达式是否符合规则
func IfAccordNumberRule(str string) bool {
	reg1 := regexp.MustCompile(`(p:[1-9]-)+d:(yyMMdd|yyyMMdd)-i:(d|[1-9]):[1-9]`)
	if reg1 == nil {
		fmt.Println("regexp err")
		return false
	}
	result1 := reg1.FindAllStringSubmatch(str, -1)
	if len(result1) == 0 {
		// 正则没有匹配到数据，编号格式不正确
		return false
	}
	strAry := strings.Split(str, "-")
	if len(str) == 0 {
		return false
	}
	// 后面尾号生成不按时间类型来生成
	if strings.Count(str, "d") != 4 {
		pNum := 0
		ruleNum := 0
		for _, v := range strAry {
			if v[0:1] == "i" {
				pointStr := strings.Split(v, ":")
				ruleNu, err := strconv.Atoi(pointStr[1])
				if err != nil {
					return false
				}
				ruleNum = ruleNu
			} else if v[0:1] == "p" {
				pNum++
			}
		}
		// 如果规则选择的按类型生成跟配置类型的个数不符合(选的第几个类型生成的数字，不能大于配置的已有类型个数)
		if pNum < ruleNum {
			return false
		}
	}
	return true
}
