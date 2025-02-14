// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

type QueueEvent string

type IQueueEvent interface {
	GetEventType() QueueEvent
}

const (
	AGENT_END = "agent end"
)

const (
	LLMChunk                   QueueEvent = "llm_chunk"
	TextChunk                  QueueEvent = "text_chunk"
	AgentMessage               QueueEvent = "agent_message"
	MessageReplace             QueueEvent = "message_replace"
	MessageEnd                 QueueEvent = "message_end"
	AdvancedChatMessageEnd     QueueEvent = "advanced_chat_message_end"
	WorkflowStarted            QueueEvent = "workflow_started"
	WorkflowSucceeded          QueueEvent = "workflow_succeeded"
	WorkflowFailed             QueueEvent = "workflow_failed"
	IterationStart             QueueEvent = "iteration_start"
	IterationNext              QueueEvent = "iteration_next"
	IterationCompleted         QueueEvent = "iteration_completed"
	NodeStarted                QueueEvent = "node_started"
	NodeSucceeded              QueueEvent = "node_succeeded"
	NodeFailed                 QueueEvent = "node_failed"
	RetrieverResources         QueueEvent = "retriever_resources"
	AnnotationReply            QueueEvent = "annotation_reply"
	AgentThought               QueueEvent = "agent_thought"
	MessageFile                QueueEvent = "message_file"
	ParallelBranchRunStarted   QueueEvent = "parallel_branch_run_started"
	ParallelBranchRunSucceeded QueueEvent = "parallel_branch_run_succeeded"
	ParallelBranchRunFailed    QueueEvent = "parallel_branch_run_failed"
	Error                      QueueEvent = "error"
	Ping                       QueueEvent = "ping"
	Stop                       QueueEvent = "stop"
)

type AppQueueEvent struct {
	Event QueueEvent `json:"event"`
}

func NewAppQueueEvent(event QueueEvent) *AppQueueEvent {
	return &AppQueueEvent{Event: event}
}

func (e *AppQueueEvent) GetEventType() QueueEvent {
	return e.Event
}

type QueueErrorEvent struct {
	*AppQueueEvent
	Err error
}

type QueueLLMChunkEvent struct {
	*AppQueueEvent
	Chunk *LLMResultChunk `json:"chunk"`
}

type QueueAgentMessageEvent struct {
	*AppQueueEvent
	Chunk *LLMResultChunk `json:"chunk"`
}

type QueueAgentThoughtEvent struct {
	AgentThoughtID string `json:"agent_thought_id"`
	*AppQueueEvent
}
type QueueTextChunkEvent struct {
	*AppQueueEvent
	Text                 string    `json:"text"`
	FromVariableSelector *[]string `json:"from_variable_selector,omitempty"`
	InIterationID        *string   `json:"in_iteration_id,omitempty"`
}

type QueueMessageEndEvent struct {
	*AppQueueEvent
	LLMResult *LLMResult `json:"llm_result"`
}

type QueueAnnotationReplyEvent struct {
	*AppQueueEvent
	MessageAnnotationID string `json:"message_annotation_id"`
}

type QueueMessage struct {
	TaskID  string      `json:"task_id"`
	AppMode string      `json:"app_mode"`
	Event   IQueueEvent `json:"event"`
}

type MessageQueueMessage struct {
	TaskID         string      `json:"task_id"`
	AppMode        string      `json:"app_mode"`
	Event          IQueueEvent `json:"event"`
	MessageID      string      `json:"message_id"`
	ConversationID string      `json:"conversation_id"`
}
