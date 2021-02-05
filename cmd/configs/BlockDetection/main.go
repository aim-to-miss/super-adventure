package main

import (
	"fmt"
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
	match, _ := regexp.MatchString("^#.*[\n]?", line)
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
	return 0
}
func seekLineNumberRare(line string, lines []string) int {
	for ittr := len(lines) - 1; ittr >= 0; ittr-- {
		if lines[ittr] == line {
			return ittr
		}
	}
	return 0
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

func main() {
	inputString := "user       www www;  ## Default: nobody\n" +
		"worker_processes  5;  ## Default: 1\n" +
		"error_log  logs/error.log;\n" +
		"pid        logs/nginx.pid;\n" +
		"worker_rlimit_nofile 8192;\n" +
		"\n" +
		"events {\n" +
		"  worker_connections  4096;  ## Default: 1024\n" +
		"}\n" +
		"\n" +
		"http {\n" +
		"  include    conf/mime.types;\n" +
		"  include    /etc/nginx/proxy.conf;\n" +
		"  include    /etc/nginx/fastcgi.conf;\n" +
		"  index    index.html index.htm index.php;\n" +
		"\n" +
		"  default_type application/octet-stream;\n" +
		"  log_format   main '$remote_addr - $remote_user [$time_local]  $status '\n" +
		"    '\"$request\" $body_bytes_sent \"$http_referer\" '\n" +
		"    '\"$http_user_agent\" \"$http_x_forwarded_for\"';\n" +
		"  access_log   logs/access.log  main;\n" +
		"  sendfile     on;\n  tcp_nopush   on;\n" +
		"  server_names_hash_bucket_size 128; # this seems to be required for some vhosts\n" +
		"\n" +
		"  server { # php/fastcgi\n" +
		"    listen       80;\n" +
		"    server_name  domain1.com www.domain1.com;\n" +
		"    access_log   logs/domain1.access.log  main;\n" +
		"    root         html;\n" +
		"\n" +
		"    location ~ \\.php$ {\n" +
		"      fastcgi_pass   127.0.0.1:1025;\n" +
		"    }\n" +
		"  }\n" +
		"\n" +
		"  server { # simple reverse-proxy\n" +
		"    listen       80;\n" +
		"    server_name  domain2.com www.domain2.com;\n" +
		"    access_log   logs/domain2.access.log  main;\n" +
		"\n" +
		"    # serve static files\n" +
		"    location ~ ^/(images|javascript|js|css|flash|media|static)/  {\n" +
		"      root    /var/www/virtual/big.server.com/htdocs;\n" +
		"      expires 30d;\n" +
		"    }\n" +
		"\n" +
		"    # pass requests for dynamic content to rails/turbogears/zope, et al\n" +
		"    location / {\n" +
		"      proxy_pass      http://127.0.0.1:8080;\n" +
		"    }\n" +
		"  }\n" +
		"\n" +
		"  upstream big_server_com {\n" +
		"    server 127.0.0.3:8000 weight=5;\n" +
		"    server 127.0.0.3:8001 weight=5;\n" +
		"    server 192.168.0.1:8000;\n" +
		"    server 192.168.0.1:8001;\n" +
		"  }\n" +
		"\n" +
		"  server { # simple load balancing\n" +
		"    listen          80;\n" +
		"    server_name     big.server.com;\n" +
		"    access_log      logs/big.server.access.log main;\n" +
		"\n" +
		"    location / {\n" +
		"      proxy_pass      http://big_server_com;\n" +
		"    }\n" +
		"  }\n" +
		"}"

	var AllBlocks NginxBlocks
	AllBlocks = GetNginxBlocks(inputString)

	for _, value := range AllBlocks.blocks {
		fmt.Println("----------Start Of a Block------------------\n")
		fmt.Println("Start Line: " + value.StartLine)
		fmt.Println("End Line : " + value.EndLine)
		fmt.Printf("Nasted Level : %v\n", value.NastedLevel)
		fmt.Printf("Number of nasted blocks: %v\n", value.TotalBlocksInside)
		fmt.Println("Full block:\n" + value.AllContents)
		fmt.Println("\n----------End Of a Block------------------\n")
	}

}
