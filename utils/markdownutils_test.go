package utils

import (
	"testing"
)

func Test_escapeMarkdown(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"all", args{value: "_*~)(`>#+-=|{}.!]["}, "\\_\\*\\~\\)\\(\\`\\>\\#\\+\\-\\=\\|\\{\\}\\.\\!\\]\\["},
		{"_ -> \\_", args{value: "_"}, "\\_"},
		{"* -> \\*", args{value: "*"}, "\\*"},
		{"~ -> \\~", args{value: "~"}, "\\~"},
		{") -> \\)", args{value: ")"}, "\\)"},
		{"( -> \\(", args{value: "("}, "\\("},
		{"` -> \\`", args{value: "`"}, "\\`"},
		{"> -> \\>", args{value: ">"}, "\\>"},
		{"# -> \\#", args{value: "#"}, "\\#"},
		{"+ -> \\+", args{value: "+"}, "\\+"},
		{"- -> \\-", args{value: "-"}, "\\-"},
		{"= -> \\=", args{value: "="}, "\\="},
		{"| -> \\|", args{value: "|"}, "\\|"},
		{"{ -> \\{", args{value: "{"}, "\\{"},
		{"} -> \\}", args{value: "}"}, "\\}"},
		{". -> \\.", args{value: "."}, "\\."},
		{"! -> \\!", args{value: "!"}, "\\!"},
		{"] -> \\]", args{value: "]"}, "\\]"},
		{"[ -> \\[", args{value: "["}, "\\["},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeMarkdown(tt.args.value); got != tt.want {
				t.Errorf("escapeMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
