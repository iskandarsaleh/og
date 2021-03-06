package og

import (
	"github.com/champii/og/lib/common"
	"strings"
)

type OgPreproc struct {
}

func (this OgPreproc) indentCount(str string) int {
	for i := range str {
		if str[i] != ' ' && str[i] != '\t' {
			return i
		}
	}
	return 0
}
func (this *OgPreproc) Run(file *common.File) {
	lineCount := 0
	file.Source = []byte(strings.Replace(string(file.Source), "\t", "  ", -1))
	lines := strings.Split(string(file.Source), "\n")
	res := []string{}
	lastIndent := 0
	indentSize := 0
	file.LineMapping = []int{0}
	for i := range lines {
		lineCount++
		if len(lines[i]) == 0 {
			continue
		}
		file.LineMapping = append(file.LineMapping, lineCount)
		indent := this.indentCount(lines[i])
		token := "=>"
		eqIf := "= if "
		if len(lines[i]) > indent+4 && lines[i][indent:indent+4] == "for " {
			lines[i] += ";"
		}
		if len(lines[i]) > indent+3 && lines[i][indent:indent+3] == "if " && !strings.Contains(lines[i], token) {
			lines[i] += ";"
		}
		if len(lines[i]) > indent+3 && strings.Contains(lines[i], eqIf) && !strings.Contains(lines[i], token) {
			lines[i] += ";"
		}
		if len(lines[i]) > indent+5 && lines[i][indent:indent+5] == "else " && !strings.Contains(lines[i], token) {
			lines[i] += ";"
		}
		if indentSize == 0 && indent != lastIndent {
			indentSize = indent
		}
		if indent > lastIndent {
			res[len(res)-1] = res[len(res)-1] + " {"
		} else if indent < lastIndent {
			indentBuff := lastIndent - indent
			for indentBuff > 0 {
				res = append(res, "}")
				file.LineMapping = append(file.LineMapping, lineCount)
				indentBuff -= indentSize
			}
		}
		lastIndent = indent
		res = append(res, lines[i])
	}
	for lastIndent > 0 {
		res = append(res, "}")
		file.LineMapping = append(file.LineMapping, lineCount)
		lastIndent -= indentSize
	}
	file.Output = strings.Join(res, "\n") + "\n"
}
func NewOgPreproc() *OgPreproc {
	return &OgPreproc{}
}
