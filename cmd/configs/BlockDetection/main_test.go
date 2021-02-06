package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetAMatch(t *testing.T) {
	type args struct {
		line string
		i    int
	}
	var Allargs = []args{
		{"events {\n  worker_connections  4096;  ## Default: 1024\n}", 0},
		{"user       www www;  ## Default: nobody\nworker_processes  5;  ## Default: 1\nerror_log  logs/error.log;\npid        logs/nginx.pid;\nworker_rlimit_nofile 8192;\n\nevents {\n  worker_connections  4096;  ## Default: 1024\n}", 0},
		{"a{\nb\n}{\nc\n}", 0},
		{"a{\nb\n}{\nc\n}", 1},
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Get A Match Test One", Allargs[0], "events {\n  worker_connections  4096;  ## Default: 1024\n}"},
		{"Get A Match Test Test Two", Allargs[1], "events {\n  worker_connections  4096;  ## Default: 1024\n}"},
		{"Get A Match Test Test Three", Allargs[2], "a{\nb\n}"},
		{"Get A Match Test Test Four", Allargs[3], "{\nc\n}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAMatch(tt.args.line, tt.args.i); got != tt.want {
				t.Errorf("GetAMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTotalMatch(t *testing.T) {
	type args struct {
		line string
	}
	var Allargs = []args{
		{"events {\n  worker_connections  4096;  ## Default: 1024\n}"},
		{"user       www www;  ## Default: nobody\nworker_processes  5;  ## Default: 1\nerror_log  logs/error.log;\npid        logs/nginx.pid;\nworker_rlimit_nofile 8192;\n\nevents {\n  worker_connections  4096;  ## Default: 1024\n}"},
		{"a{\nb\n}{\nc\n}"},
		{"a{\nb\n}\nc\n}"},
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Get Total Match Test One", Allargs[0], 1},
		{"Get Total Match Test Test Two", Allargs[1], 1},
		{"Get Total Match Test Test Three", Allargs[2], 2},
		{"Get Total Match Test Test Four", Allargs[3], 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTotalMatch(tt.args.line); got != tt.want {
				t.Errorf("GetTotalMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNginxBlock_HasComment(t *testing.T) {
	type fields struct {
		StartLine         string
		EndLine           string
		NastedLevel       int
		AllContents       string
		AllLines          []string
		NestedBlocks      []NginxBlock
		TotalBlocksInside int
	}
	var testfields = []fields{
		{"user www www;  ## Default: nobody", "worker_rlimit_nofile 8192;", 0, "user www www;  ## Default: nobody", nil, nil, 0},
	}

	type args struct {
		line string
	}
	var testStringArgs = []args{
		{testfields[0].StartLine},
		{testfields[0].EndLine},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Test one", testfields[0], testStringArgs[0], true},
		{"Test two", testfields[0], testStringArgs[1], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngBlock := &NginxBlock{
				StartLine:         tt.fields.StartLine,
				EndLine:           tt.fields.EndLine,
				NastedLevel:       tt.fields.NastedLevel,
				AllContents:       tt.fields.AllContents,
				AllLines:          tt.fields.AllLines,
				NestedBlocks:      tt.fields.NestedBlocks,
				TotalBlocksInside: tt.fields.TotalBlocksInside,
			}
			if got := ngBlock.HasComment(tt.args.line); got != tt.want {
				t.Errorf("HasComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNginxBlock_IsBlock(t *testing.T) {
	type fields struct {
		StartLine         string
		EndLine           string
		NastedLevel       int
		AllContents       string
		AllLines          []string
		NestedBlocks      []NginxBlock
		TotalBlocksInside int
	}
	var testfields = []fields{
		{"", "", 0, "events {\n  worker_connections  4096;  ## Default: 1024\n}", nil, nil, 0},
		{"", "", 0, "location ~ .php$ {\nfastcgi_pass127.0.0.1:1025;}{\n}", nil, nil, 0},
		{"", "", 0, "location ~ .php$ \nfastcgi_pass127.0.0.1:1025;\n}", nil, nil, 0},
		{"", "", 0, "location ~ .php$ {\nfastcgi_pass127.0.0.1:1025;\n", nil, nil, 0},
	}
	type args struct {
		line string
	}

	var testStringArgs = []args{
		{testfields[0].AllContents},
		{testfields[1].AllContents},
		{testfields[2].AllContents},
		{testfields[3].AllContents},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Is Block Test One", testfields[0], testStringArgs[0], true},
		{"Is Block Test Two", testfields[0], testStringArgs[1], true},
		{"Is Block Test Three", testfields[0], testStringArgs[2], false},
		{"Is Block Test Four", testfields[0], testStringArgs[3], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngBlock := &NginxBlock{
				StartLine:         tt.fields.StartLine,
				EndLine:           tt.fields.EndLine,
				NastedLevel:       tt.fields.NastedLevel,
				AllContents:       tt.fields.AllContents,
				AllLines:          tt.fields.AllLines,
				NestedBlocks:      tt.fields.NestedBlocks,
				TotalBlocksInside: tt.fields.TotalBlocksInside,
			}
			if got := ngBlock.IsBlock(tt.args.line); got != tt.want {
				t.Errorf("IsBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNginxBlock_IsLine(t *testing.T) {
	type fields struct {
		StartLine         string
		EndLine           string
		NastedLevel       int
		AllContents       string
		AllLines          []string
		NestedBlocks      []NginxBlock
		TotalBlocksInside int
	}
	var testfields = []fields{
		{"", "", 0, "events {\n  worker_connections  4096;  ## Default: 1024\n}", nil, nil, 0},
		{"", "", 0, "location ~ .php$ {\nfastcgi_pass127.0.0.1:1025;}{\n}", nil, nil, 0},
		{"", "", 0, "location ~ .php$ \nfastcgi_pass127.0.0.1:1025;\n}", nil, nil, 0},
		{"", "", 0, "location ~ .php$ {\nfastcgi_pass127.0.0.1:1025;\n", nil, nil, 0},
		{"", "", 0, "location ~ .php$ \nfastcgi_pass127.0.0.1:1025;\n", nil, nil, 0},
	}
	type args struct {
		line string
	}
	var testStringArgs = []args{
		{testfields[0].AllContents},
		{testfields[1].AllContents},
		{testfields[2].AllContents},
		{testfields[3].AllContents},
		{testfields[4].AllContents},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Is Line Test One", testfields[0], testStringArgs[0], true},
		{"Is Line Test Two", testfields[0], testStringArgs[1], true},
		{"Is Line Test Three", testfields[0], testStringArgs[2], true},
		{"Is Line Test Four", testfields[0], testStringArgs[3], true},
		{"Is Line Test Five", testfields[0], testStringArgs[4], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngBlock := &NginxBlock{
				StartLine:         tt.fields.StartLine,
				EndLine:           tt.fields.EndLine,
				NastedLevel:       tt.fields.NastedLevel,
				AllContents:       tt.fields.AllContents,
				AllLines:          tt.fields.AllLines,
				NestedBlocks:      tt.fields.NestedBlocks,
				TotalBlocksInside: tt.fields.TotalBlocksInside,
			}
			if got := ngBlock.IsLine(tt.args.line); got != tt.want {
				t.Errorf("IsLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seekLineNumberFront(t *testing.T) {
	type args struct {
		line  string
		lines []string
	}
	allLines := strings.Split("user       www www;  ## Default: nobody\nworker_processes  5;  ## Default: 1\nerror_log  logs/error.log;\npid        logs/nginx.pid;\nworker_rlimit_nofile 8192;", "\n")
	var allTestArgs = []args{
		{"user       www www;  ## Default: nobody", allLines},
		{"worker_processes  5;  ## Default: 1\nerror_log  logs/error.log;", allLines},
		{"logs/nginx.pid;", allLines},
		{"pid        logs/nginx.pid;", allLines},
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Seek from front Test  One", allTestArgs[0], 0},
		{"Seek from front Test  Two", allTestArgs[1], -1},
		{"Seek from front Test  Three", allTestArgs[2], -1},
		{"Seek from front Test  Four", allTestArgs[3], 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seekLineNumberFront(tt.args.line, tt.args.lines); got != tt.want {
				t.Errorf("seekLineNumberFront() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seekLineNumberRare(t *testing.T) {
	type args struct {
		line  string
		lines []string
	}
	allLines := strings.Split("user       {www www;  ## Default: nobody\nworker_processes  5;}  ## Default: 1\n }\nerror_log  logs/error.log;\npid        \n }\nlogs/nginx.pid;\nworker_rlimit_nofile 8192;\n}", "\n")
	var allTestArgs = []args{
		{"user       {www www;  ## Default: nobody", allLines},
		{"  }", allLines},
		{" }", allLines},
		{"}", allLines},
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Seek from front Test  One", allTestArgs[0], 0},
		{"Seek from front Test  Two", allTestArgs[1], -1},
		{"Seek from front Test  Three", allTestArgs[2], 5},
		{"Seek from front Test  Four", allTestArgs[3], 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seekLineNumberRare(tt.args.line, tt.args.lines); got != tt.want {
				t.Errorf("seekLineNumberRare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNginxBlocks(t *testing.T) {
	type args struct {
		configContent string
	}
	var allArgs = []args{
		{"events {\n  worker_connections  4096;  ## Default: 1024\n}"},
	}
	var wantedNginxBlock = []NginxBlock{
		{"events {", "}", 0, allArgs[0].configContent, strings.Split(allArgs[0].configContent, "\n"), nil, 0},
	}
	var testWant = NginxBlocks{
		wantedNginxBlock,
		allArgs[0].configContent,
		strings.Split(allArgs[0].configContent, "\n"),
	}
	tests := []struct {
		name string
		args args
		want NginxBlocks
	}{
		{"GetTestGetNginxBlocks Test One", allArgs[0], testWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNginxBlocks(tt.args.configContent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNginxBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNginxBlock(t *testing.T) {
	type args struct {
		lines        []string
		startIndex   int
		endIndex     int
		recursionMax int
		nastedLevel  int
	}
	testString := "  server { # php/fastcgi\n    listen       80;\n    server_name  domain1.com www.domain1.com;\n    access_log   logs/domain1.access.log  main;\n    root         html;\n\n    location ~ .php$ \n      fastcgi_pass   127.0.0.1:1025;\n    \n  }"
	var testWant = []NginxBlock{
		{"  server { # php/fastcgi", "  }", 0, testString, strings.Split(testString, "\n"), nil, 0},
	}
	var allTestArgs = []args{
		{strings.Split(testString, "\n"), 0, 10, 5, 0},
		{strings.Split(testString, "\n"), 0, 10, 0, 0},
	}
	tests := []struct {
		name string
		args args
		want []NginxBlock
	}{
		{"Get Nginx Block Test", allTestArgs[0], testWant},
		{"Get Nginx Block Test", allTestArgs[1], nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNginxBlock(tt.args.lines, tt.args.startIndex, tt.args.endIndex, tt.args.recursionMax, tt.args.nastedLevel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNginxBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
