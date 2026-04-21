package mq

import "fmt"

type Topic struct {
	Name string
	Partitions int
}

func NewTopic(name string, partitions int) *Topic {
	return &Topic{Name: name, Partitions: partitions}
}

func (t *Topic) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("topic name required")
	}
	if t.Partitions < 1 {
		return fmt.Errorf("partitions must be >= 1")
	}
	return nil
}

func (t *Topic) GetPartition(key string) int {
	hash := 0
	for _, c := range key {
		hash += int(c)
	}
	return hash % t.Partitions
}
