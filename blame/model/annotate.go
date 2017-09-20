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

type FileAnnotation struct {
	Annotation
	File *File
}

type LineAnnotation struct {
	Annotation
	Start int
	End   int
	Line  *Line
}

func (a *annotation) GetCommits() git.Commits           { return a.commits }
func (a *annotation) GetWorkItems() []apis.WorkItem     { return a.workItems }
func (a *annotation) GetContributors() git.Contributors { return a.commits.RelatedContributors() }

type Annotate interface {
	File(file *File) *FileAnnotation
	Line(line *Line) *LineAnnotation
}

type annotate struct {
	server apis.WorkItemServer
}

func NewAnnotate(server apis.WorkItemServer) Annotate {
	return &annotate{server}
}

func (a *annotate) File(file *File) *FileAnnotation {
	commits, err := git.CommitsAfter(time.Now().AddDate(-1, 0, 0))
	if err != nil {
		log.Panicln(err)
	}

	fileCommits := commits.ContainsFile(file.Filename)
	workItems, err := apis.GetWorkItems(a.server, fileCommits.RelatedWorkItems()...)
	if err != nil {
		log.Panicln(err)
	}

	return &FileAnnotation{&annotation{fileCommits, workItems}, file}
}

func (a *annotate) Line(line *Line) *LineAnnotation {
	start := 1
	end := line.Number + 1

	if line.Number > 1 {
		start = line.Number - 1
	}

	commits, err := git.Log(fmt.Sprintf("-L%d,%d:%s", start, end, line.File.Filename))
	workItems, err := apis.GetWorkItems(a.server, commits.RelatedWorkItems()...)
	if err != nil {
		log.Panicln(err)
	}

	return &LineAnnotation{&annotation{commits, workItems}, start, end, line}
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

func (c *cache) File(file *File) *FileAnnotation { return c.annotate.File(file) }
func (c *cache) Line(line *Line) *LineAnnotation {
	value, err := c.cache.Get(line)
	if err != nil {
		log.Panicln(err)
	}

	annotation, ok := value.(*LineAnnotation)
	if !ok {
		log.Panicln("Unknown Result")
	}

	return annotation
}
