package inputs

import "fmt"

type ConfigParams struct {
	GistToken string
	Encrypt   bool
}

func (input ConfigParams) GetItems() []string {
	return []string{input.GistToken, fmt.Sprintf("%v", input.Encrypt)}
}

func (input ConfigParams) GetInputType() InputType {
	return CONFIG
}

func (input *ConfigParams) GetSummary() string {
	return ""
}

type SyncParams struct {
	// create a public gist
	Public bool
	// select a single item to sync
	SingleItem string
}

func (input *SyncParams) GetItems() []string {
	return []string{fmt.Sprintf("%v", input.Public)}
}

func (input *SyncParams) GetInputType() InputType {
	return SYNC
}

func (input *SyncParams) GetSummary() string {
	return input.SingleItem
}

type SwitchParam struct {
	Target string
}

func (input *SwitchParam) GetItems() []string {
	return []string{}
}

func (input *SwitchParam) GetInputType() InputType {
	return SWITCH
}

func (input *SwitchParam) GetSummary() string {
	return input.Target
}

type PrintParam struct {
	Target string
}

func (input *PrintParam) GetItems() []string {
	return []string{}
}

func (input *PrintParam) GetInputType() InputType {
	return PRINT
}

func (input *PrintParam) GetSummary() string {
	return input.Target
}

type ListParam struct {
}

func (input *ListParam) GetItems() []string {
	return []string{}
}

func (input *ListParam) GetInputType() InputType {
	return LIST
}

func (input *ListParam) GetSummary() string {
	return ""
}
