package main

import (
	"regexp"
	"strings"
)

// Blocks will be valid
// Number of opening and number of closing braces are same
// Identify all the blocks and return the starting and ending line number along with the same information about their children and grandchildren and  so on.
// Count the block nested level.
// Finally implement only the below two functions GetNginxBlock & GetNginxBlocks

type NginxBlock struct {
	StartLine         string
	EndLine           string
	NastedLevel       int
	AllContents       string
	AllLines          []string
	NestedBlocks      []NginxBlock
	TotalBlocksInside int
}

type NginxBlocks struct {
	blocks      []NginxBlock
	AllContents string
	AllLines    []string
}

func (ngBlock *NginxBlock) IsBlock(line string) bool {
	match, _ := regexp.MatchString("[\\SA-Za-z ]*{(?:(?:{(?:(?:{(?:[^{}])*})|(?:[^{}]))*})|(?:[^{}]))*}", line)
	return match
}
func (ngBlock *NginxBlock) IsLine(line string) bool {
	match, _ := regexp.MatchString(".*[{}].*", line)
	return match
}
func (ngBlock *NginxBlock) HasComment(line string) bool {
	match, _ := regexp.MatchString(".*[#]+.*", line)
	return match
}

func GetAMatch(line string, i int) string {
	re := regexp.MustCompile("[\\SA-Za-z ]*{(?:(?:{(?:(?:{(?:[^{}])*})|(?:[^{}]))*})|(?:[^{}]))*}")
	return re.FindAllString(line, -1)[i]
}
func GetTotalMatch(line string) int {
	re := regexp.MustCompile("[\\SA-Za-z ]*{(?:(?:{(?:(?:{(?:[^{}])*})|(?:[^{}]))*})|(?:[^{}]))*}")
	return len(re.FindAllString(line, -1))
}
func seekLineNumberFront(line string, lines []string) int {
	for key, value := range lines {
		if value == line {
			return key
			break
		}
	}
	return -1
}
func seekLineNumberRare(line string, lines []string) int {
	for ittr := len(lines) - 1; ittr >= 0; ittr-- {
		if lines[ittr] == line {
			return ittr
		}
	}
	return -1
}

var nginBloclsList []NginxBlock
var nastedLevelInitialCounterHolder = []int{0, 0, 0, 0, 0, 0, 0}
var nastedLevelLastCounterHolder = []int{0, 0, 0, 0, 0, 0, 0}

func GetNginxBlocks(configContent string) NginxBlocks {
	allLines := strings.Split(configContent, "\n")
	startLineIndex := 0
	endLineIndex := len(allLines)
	var nginBlocks NginxBlocks
	nginBlocks.AllContents = configContent
	nginBlocks.AllLines = allLines
	nginBlocks.blocks = GetNginxBlock(allLines, startLineIndex, endLineIndex, 10, 0)
	return nginBlocks
}

func GetNginxBlock(
	lines []string,
	startIndex,
	endIndex,
	recursionMax,
	nastedLevel int,
) []NginxBlock {
	blockToConsider := strings.Join(lines[startIndex:endIndex], "\n")
	numberOfMatch := GetTotalMatch(blockToConsider)
	nastedLevelInitialCounterHolder[nastedLevel] = numberOfMatch
	var aNginxBlock NginxBlock
	for i := nastedLevelLastCounterHolder[nastedLevel]; i < nastedLevelInitialCounterHolder[nastedLevel]; i++ {
		if recursionMax <= 0 {
			break
		}
		recursionMax -= 1
		aBlock := GetAMatch(blockToConsider, i)
		linesInABlock := strings.Split(aBlock, "\n")
		aNginxBlock.StartLine = linesInABlock[0]
		aNginxBlock.EndLine = linesInABlock[len(linesInABlock)-1]
		aNginxBlock.AllLines = linesInABlock
		aNginxBlock.AllContents = strings.Join(linesInABlock, "\n")
		aNginxBlock.NastedLevel = nastedLevel
		recurStartIndex := seekLineNumberFront(aNginxBlock.StartLine, linesInABlock)
		recurEndIndex := seekLineNumberRare(aNginxBlock.EndLine, linesInABlock)
		nastedLevelLastCounterHolder[nastedLevel] = i + 1
		aNginxBlock.NestedBlocks = GetNginxBlock(linesInABlock, recurStartIndex+1, recurEndIndex, recursionMax, nastedLevel+1)
		aNginxBlock.TotalBlocksInside = len(aNginxBlock.NestedBlocks)
		nginBloclsList = append(nginBloclsList, aNginxBlock)
	}
	nastedLevelLastCounterHolder[nastedLevel] = 0
	nastedLevel -= 1
	return nginBloclsList
}
