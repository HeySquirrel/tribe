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

	fileCommits := commits.ContainsFile(file.Name)
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

	commits, err := git.Log(fmt.Sprintf("-L%d,%d:%s", start, end, line.File.Name))
	workItems, err := apis.GetWorkItems(a.server, commits.RelatedWorkItems()...)
	if err != nil {
		log.Panicln(err)
	}

	return &LineAnnotation{&annotation{commits, workItems}, start, end, line}
}

type cache struct {
	annotate Annotate
	cache    gcache.Cache
	seed     chan *Line
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
	seed := make(chan *Line)
	c := &cache{annotate, gc, seed}

	c.startCacheWorkers()

	return c
}

func (c *cache) startCacheWorkers() {
	for i := 0; i < 3; i++ {
		go func() {
			for line := range c.seed {
				_, err := c.cache.Get(line)
				if err != nil {
					log.Panicln(err)
				}
			}
		}()
	}
}

func (c *cache) seedCacheForFile(file *File) {
	for i := file.Start; i <= file.End; i++ {
		c.seed <- file.GetLine(i)
	}
}

func (c *cache) seedCacheForLine(line *Line) {
	if line.Number > 1 {
		c.seed <- line.File.GetLine(line.Number - 1)
	}

	if line.Number < line.File.Len() {
		c.seed <- line.File.GetLine(line.Number + 1)
	}
}

func (c *cache) File(file *File) *FileAnnotation {
	go c.seedCacheForFile(file)

	return c.annotate.File(file)
}

func (c *cache) Line(line *Line) *LineAnnotation {
	go c.seedCacheForLine(line)

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
