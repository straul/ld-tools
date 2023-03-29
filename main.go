package main

import (
	"bufio"
	"bytes"
	"fmt"
	inf "github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components/input/text"
	"github.com/fzdwx/infinite/theme"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"strings"
)

// 源文件: i.csv
// 输出文件: o.csv
func main() {
	// 退出
	waitExit := inf.NewText(
		text.WithPrompt("执行完成，按任意键退出"),
		text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
		//text.WithRequired(),
	)

	// 录入源文件路径
	inputInFile := inf.NewText(
		text.WithPrompt("把源文件名放在这里（只能是同文件夹文件，默认是当前文件夹内的 i.csv）"),
		text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
		text.WithDefaultValue("i.csv"),
		//text.WithRequired(),
	)

	_, _ = inputInFile.Display()

	fmt.Printf("源文件路径: %s\n", inputInFile.Value())
	inFilePath := inputInFile.Value()

	//myFile, err := os.Open("i.csv")
	myFile, err := os.Open(inFilePath)
	if err != nil {
		fmt.Println("读取源文件错误: ", err)
		_, _ = waitExit.Display()
		return
	}
	defer myFile.Close()

	// 录入输出文件路径
	inputOutFile := inf.NewText(
		text.WithPrompt("把输出文件名放在这里（只能是同文件夹文件，默认是当前文件夹内的 o.csv）"),
		text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
		text.WithDefaultValue("o.csv"),
		//text.WithRequired(),
	)

	_, _ = inputOutFile.Display()

	fmt.Printf("输出文件路径: %s\n", inputOutFile.Value())
	outFilePath := inputOutFile.Value()

	// 重复检验
	m := make(map[string]string)

	scanner := bufio.NewScanner(myFile)

	for scanner.Scan() {
		// 读出每一行的内容
		line := scanner.Text()
		lineGbk, _ := GbkToUtf8([]byte(line))
		line = string(lineGbk)

		// 提取出每一行第一个逗号左边的文字
		lineArray := strings.Split(line, ",")
		left := lineArray[0]

		// 把第一个逗号左边的内容去掉
		right := strings.Replace(line, left+",", "", -1)
		// 去掉两端的","
		right = strings.Trim(right, ",")
		// 把右边的英文逗号转为中文逗号，避免csv文件中被认为是两列
		right = strings.Replace(right, ",", "，", -1)

		// 判断当前的 left 是否已经处理过
		if rightHistory, ok := m[left]; ok {
			// 处理过，把本行内容追加进去
			m[left] = rightHistory + "_" + right
		} else {
			// 没处理过，把本行内容放进去
			m[left] = right
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file: ", err)
	}

	// 输出
	//outFile := "o.csv"
	outFile := outFilePath
	file, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		_, _ = waitExit.Display()
		fmt.Println("打开输出文件错误: ", err)
		return
	}
	defer file.Close()

	write := bufio.NewWriter(file)

	for k, v := range m {
		newLine := k + "," + v + "\n"
		// 转gbk
		newLineGbk, _ := Utf8ToGbk([]byte(newLine))
		//write.WriteString(newLine)
		write.WriteString(string(newLineGbk))
	}

	write.Flush()

	_, _ = waitExit.Display()
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
