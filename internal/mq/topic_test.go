package mq

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewTopic(t *testing.T) {
	topic := NewTopic("telemetry", 3)
	assert.Equal(t, "telemetry", topic.Name)
	assert.Equal(t, 3, topic.Partitions)
}

func TestTopicValidate(t *testing.T) {
	topic := NewTopic("test", 1)
	assert.NoError(t, topic.Validate())
	
	topic2 := NewTopic("", 1)
	assert.Error(t, topic2.Validate())
	
	topic3 := NewTopic("test", 0)
	assert.Error(t, topic3.Validate())
}

func TestGetPartition(t *testing.T) {
	topic := NewTopic("test", 3)
	p1 := topic.GetPartition("gpu-1")
	p2 := topic.GetPartition("gpu-1")
	assert.Equal(t, p1, p2)
	assert.True(t, p1 >= 0 && p1 < 3)
}
