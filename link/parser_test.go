package link

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"testing"
)

type args struct {
	htmlBytes []byte
}

type testCase struct {
	name string
	args args
	want []Link
}

func setupTests() ([]testCase, error) {
	var tests []testCase
	expected := [][]Link{
		// ex1
		{
			{
				Href: "/other-page",
				Text: "A link to another page",
			},
		},
		// ex2
		{
			{
				Href: "https://www.twitter.com/joncalhoun",
				Text: "Check me out on twitter",
			},
			{
				Href: "https://github.com/gophercises",
				Text: "Gophercises is on Github!",
			},
		},
		// ex3
		{
			{
				Href: "#",
				Text: "Login",
			},
			{
				Href: "/lost",
				Text: "Lost? Need help?",
			},
			{
				Href: "https://twitter.com/marcusolsson",
				Text: "@marcusolsson",
			},
		},
		// ex4
		{
			{
				Href: "/dog-cat",
				Text: "dog cat",
			},
		},
	}

	for i := 1; i < 5; i++ {
		iStr := strconv.Itoa(i)
		bytes, err := ioutil.ReadFile("./ex" + iStr + ".html")
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{
			name: "ex" + iStr,
			args: args{bytes},
			want: expected[i-1],
		})
	}

	return tests, nil
}

func Test_parseHTML(t *testing.T) {
	tests, err := setupTests()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseHTML(tt.args.htmlBytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
