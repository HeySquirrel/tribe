package model

import (
	"errors"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/git"
	"log"
	"time"
)

type Annotation interface {
	GetCommits() git.Commits
	GetWorkItems() []apis.WorkItem
	GetContributors() git.Contributors
}
type annotation struct {
	commits   git.Commits
	workItems []apis.WorkItem
}

func (a *annotation) GetCommits() git.Commits           { return a.commits }
func (a *annotation) GetWorkItems() []apis.WorkItem     { return a.workItems }
func (a *annotation) GetContributors() git.Contributors { return a.commits.RelatedContributors() }

type Annotate interface {
	File(file *File) Annotation
	Line(line *Line) Annotation
}

type annotate struct {
	server apis.WorkItemServer
}

func NewAnnotate(server apis.WorkItemServer) Annotate {
	return &annotate{server}
}

func (a *annotate) File(file *File) Annotation {
	commits, err := git.CommitsAfter(time.Now().AddDate(-1, 0, 0))
	if err != nil {
		log.Panicln(err)
	}

	fileCommits := commits.ContainsFile(file.Filename)
	workItems, err := apis.GetWorkItems(a.server, fileCommits.RelatedWorkItems()...)
	if err != nil {
		log.Panicln(err)
	}

	return &annotation{fileCommits, workItems}
}

func (a *annotate) Line(line *Line) Annotation {
	start := 1
	end := line.Number + 1

	if line.Number > 1 {
		start = line.Number - 1
	}

	commits, err := git.Log(fmt.Sprintf("-L%d,%d:%s", start, end, line.Filename))
	workItems, err := apis.GetWorkItems(a.server, commits.RelatedWorkItems()...)
	if err != nil {
		log.Panicln(err)
	}

	return &annotation{commits, workItems}
}

type cache struct {
	annotate Annotate
	cache    gcache.Cache
}

func NewCachingAnnotate(annotate Annotate) Annotate {
	gc := gcache.New(100).
		LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			line, ok := key.(*Line)
			if ok {
				return annotate.Line(line), nil
			}
			return nil, errors.New("Unknown line")
		}).
		Build()

	return &cache{annotate, gc}
}

func (c *cache) File(file *File) Annotation { return c.annotate.File(file) }
func (c *cache) Line(line *Line) Annotation {
	value, err := c.cache.Get(line)
	if err != nil {
		log.Panicln(err)
	}

	annotation, ok := value.(Annotation)
	if !ok {
		log.Panicln("Unknown Result")
	}

	return annotation
}
