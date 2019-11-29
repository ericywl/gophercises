package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	stories, err := parseAdventureJSON("gopher.json")
	if err != nil {
		panic(err)
	}

	cliPtr := flag.Bool("cli", false, "serves game on command line instead of browser")
	tmplPtr := flag.String("template", "adventure.html", "template file for HTML")
	flag.Parse()

	t := template.Must(template.ParseFiles(*tmplPtr))
	adv := Adventure{stories, t}

	if *cliPtr {
		adv.commandLine()
	} else {
		fmt.Println("Your adventure awaits at localhost:8080!")
		http.ListenAndServe(":8080", adv)
	}
}

// StoryArc contains the title, story text and next options
type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

// Adventure is a collection of StoryArcs
type Adventure struct {
	Stories  map[string]StoryArc
	Template *template.Template
}

func parseAdventureJSON(filename string) (map[string]StoryArc, error) {
	var stories map[string]StoryArc
	jsn, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsn), &stories)
	if err != nil {
		return nil, err
	}

	return stories, nil
}

func (adv Adventure) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arc := strings.TrimLeft(r.URL.Path, "/")
	if arc == "" {
		arc = "intro"
	}
	if storyArc, found := adv.Stories[arc]; found {
		err := adv.Template.Execute(w, &storyArc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (adv Adventure) commandLine() {
	arc := "intro"
	reader := bufio.NewReader(os.Stdin)
	for {
		if curr, found := adv.Stories[arc]; found {
			fmt.Println(strings.ToUpper(curr.Title) + "\n")
			for _, s := range curr.Story {
				fmt.Println(s + "\n")
			}

			fmt.Println("What do you do?")
			if len(curr.Options) != 0 {
				for i, o := range curr.Options {
					fmt.Printf("%d: %s\n", i+1, o.Text)
				}
			} else {
				fmt.Println("R: Restart game.")
			}
			fmt.Println("X: Quit game.")
			fmt.Println()

			for {
				fmt.Printf("Enter your choice: ")
				ans, _ := reader.ReadString('\n')
				ans = strings.TrimSpace(ans)
				if strings.ToLower(ans) == "x" {
					fmt.Println("Quitting game. Goodbye!")
					return
				}

				if strings.ToLower(ans) == "r" {
					fmt.Println("Restarting game.")
					arc = "intro"
					break
				}

				if c, err := strconv.Atoi(ans); err == nil {
					if c <= len(curr.Options) && c > 0 {
						arc = curr.Options[c-1].Arc
						break
					}
				}
			}

			fmt.Println()

		} else {
			fmt.Println(arc + "not found!")
			return
		}
	}
}
