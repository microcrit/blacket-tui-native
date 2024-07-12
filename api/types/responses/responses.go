package responses

import (
	"crit.rip/blacket-tui/api/types/objects"
)

type GenericResponse struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

type GenericMessageResponse struct {
	Error   bool   `json:"error"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type UserResponse struct {
	Error bool         `json:"error"`
	User  objects.User `json:"user"`
}

type PackOpenResponse struct {
	Error bool   `json:"error"`
	Blook string `json:"blook"`
	New   bool   `json:"new"`
}

type Chat struct {
	Data objects.ChatData `json:"data"`
}

type ClaimRewardResponse struct {
	Error       bool `json:"error"`
	RewardIndex int  `json:"reward_index"`
}

type FriendsResponse struct {
	Error     bool                `json:"error"`
	Friends   []objects.BasicUser `json:"friends"`
	Blocks    []objects.BasicUser `json:"blocks"`
	Sending   []objects.BasicUser `json:"sending"`
	Receiving []objects.BasicUser `json:"receiving"`
}

type BazaarSearchResponse struct {
	Error  bool                 `json:"error"`
	Bazaar []objects.BazaarItem `json:"bazaar"`
}
