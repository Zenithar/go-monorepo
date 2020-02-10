package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/yaml.v2"
)

// Change represents changed file information
type Change struct {
	Action   string
	Filename string
}

// Trigger represents service target target trigger
type Trigger struct {
	matcher *regexp.Regexp
	target  string
}

// -----------------------------------------------------------------------------

type yamlRoot struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Spec       yamlSpec `yaml:"spec"`
}

type yamlSpec struct {
	Triggers yamlTriggers `yaml:"triggers"`
}

type yamlTriggers map[string][]string

// -----------------------------------------------------------------------------

func splitCommit(cr string) (string, string, error) {
	// Split commit string
	parts := strings.SplitN(cr, "..", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	return cr, "head", nil
}

func getHead(repoPath string) (*plumbing.Reference, error) {
	// Open current repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("mono: unable to open given repository: %w", err)
	}

	return r.Head()
}

func getPrevious(repoPath string, commitFrom string) (string, error) {
	// Open current repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("mono: unable to open given repository: %w", err)
	}

	c, err := r.CommitObject(plumbing.NewHash(commitFrom))
	if err != nil {
		return "", fmt.Errorf("mono: unable to find commit `%s`: %w", commitFrom, err)
	}

	previousCommit, err := c.Parent(0)
	if err != nil {
		return "", fmt.Errorf("mono: unable to previous commit `%s`: %w", commitFrom, err)
	}

	return previousCommit.ID().String(), nil
}

func changeName(ch *object.Change) string {
	if ch.From.Name != "" {
		return ch.From.Name
	}
	return ch.To.Name
}

func getChanges(repoPath string, commitFrom string, commitTo string) ([]*Change, error) {
	// Open current repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("mono: unable to open given repository: %w", err)
	}

	c, err := r.CommitObject(plumbing.NewHash(commitFrom))
	if err != nil {
		return nil, fmt.Errorf("mono: unable to find commit `%s`: %w", commitFrom, err)
	}

	prevCommit, err := r.CommitObject(plumbing.NewHash(commitTo))
	if err != nil {
		return nil, fmt.Errorf("mono: unable to retrieve previous commit: %w", err)
	}

	prevTree, err := prevCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("mono: unable to retrieve previous commit tree: %w", err)
	}

	currTree, err := c.Tree()
	if err != nil {
		return nil, fmt.Errorf("mono: unable to retrieve current commit tree: %w", err)
	}

	changes, err := currTree.Diff(prevTree)
	if err != nil {
		return nil, fmt.Errorf("mono: unable to calculate changes tree: %w", err)
	}

	var files []*Change

	for _, ch := range changes {
		action, err := ch.Action()
		if err != nil {
			continue
		}

		file := &Change{
			Action:   action.String(),
			Filename: changeName(ch),
		}
		files = append(files, file)
	}

	// No error
	return files, nil
}

// StringArray describes string array type
type StringArray []string

// -----------------------------------------------------------------------------

// Contains checks if item is in collection
func (s StringArray) Contains(item string) bool {
	for _, v := range s {
		if strings.EqualFold(item, v) {
			return true
		}
	}

	return false
}

// AddIfNotContains add item if not already in collection
func (s *StringArray) AddIfNotContains(item string) {
	if s.Contains(item) {
		return
	}
	*s = append(*s, item)
}

// -----------------------------------------------------------------------------

var (
	repoPath    string
	commitRange string
	target      string
)

func init() {
	flag.StringVar(&repoPath, "repo", ".", "Repository path to scan")
	flag.StringVar(&commitRange, "range", "", "Commit range <from>..<to>")
	flag.StringVar(&target, "target", "all", "Expected target to run")

	flag.Parse()
}

func main() {
	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	// Load descriptor from yaml file
	yamlFile, err := os.Open(".monobuild.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Read from stdin
	input, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	// Read YAML definition
	var def yamlRoot
	if err := yaml.Unmarshal(input, &def); err != nil {
		log.Fatal(err)
	}

	// Check arguments
	firstCommit, lastCommit, err := splitCommit(commitRange)
	if err != nil {
		log.Fatalln(err)
	}

	// Check head usage
	if strings.EqualFold(firstCommit, "HEAD") || strings.EqualFold(lastCommit, "HEAD") {
		headRef, err := getHead(repoPath)
		if err != nil {
			log.Fatal(err)
		}

		if strings.EqualFold(firstCommit, "HEAD") {
			firstCommit = headRef.Hash().String()
		}
		if strings.EqualFold(lastCommit, "HEAD") {
			lastCommit = headRef.Hash().String()
		}
		if firstCommit == "" {
			firstCommit, err = getPrevious(repoPath, lastCommit)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Printf("Scanning repository `%s` for changes (%s..%s) to run `%s` ...\n", repoPath, firstCommit, lastCommit, target)

	// Get all changed files
	changes, err := getChanges(repoPath, firstCommit, lastCommit)
	if err != nil {
		log.Fatal(err)
	}

	// Compile definition
	triggers, err := compile(&def)
	if err != nil {
		log.Fatal(err)
	}

	// Trigger mage targets according changes
	targets := map[string]bool{}
	for _, ch := range changes {
		fmt.Printf("%s", ch.Filename)

		var fileTriggers StringArray
		for _, trigger := range triggers {
			if trigger.matcher.MatchString(ch.Filename) {
				fileTriggers.AddIfNotContains(trigger.target)
				if _, ok := targets[trigger.target]; !ok {
					targets[trigger.target] = true
				}
			}
		}

		fmt.Printf(" [%s]\n", strings.Join(fileTriggers, ","))
	}

	// Check if expected target is triggerable
	if _, ok := targets[target]; !ok {
		log.Fatalf("Aborting, target `%s` is not concerned by build\n", target)
	}

	log.Printf("Can continue, target `%s` is concerned by build\n", target)
	os.Exit(0)
}

func compile(def *yamlRoot) ([]Trigger, error) {
	triggers := []Trigger{}

	for targetName, expressions := range def.Spec.Triggers {
		for _, expr := range expressions {
			r, err := regexp.Compile(expr)
			if err != nil {
				return nil, fmt.Errorf("unable to compile expression '%s' for target '%s': %v", expr, targetName, err)
			}

			// Add to triggers
			triggers = append(triggers, Trigger{
				matcher: r,
				target:  targetName,
			})
		}
	}

	// Return triggers
	return triggers, nil
}
