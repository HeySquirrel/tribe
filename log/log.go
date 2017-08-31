package log

import (
	"container/list"
	"fmt"
	"time"
)

type LogEntry struct {
	CreatedAt time.Time
	Message   string
}

func (entry *LogEntry) String() string {
	return fmt.Sprintf("[%s] %s", entry.CreatedAt, entry.Message)
}

type Log struct {
	entries *list.List
}

func New() *Log {
	log := new(Log)
	log.entries = list.New()

	return log
}

func (log *Log) Add(message string) {
	entry := new(LogEntry)
	entry.CreatedAt = time.Now()
	entry.Message = message

	log.entries.PushFront(entry)
}

func (log *Log) Entries() []*LogEntry {
	entries := make([]*LogEntry, 0)
	logs := log.entries

	for e := logs.Front(); e != nil; e = e.Next() {
		entries = append(entries, e.Value.(*LogEntry))
	}

	return entries
}
