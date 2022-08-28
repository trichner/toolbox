package freshdesk

import (
	"encoding/json"
	"fmt"
)

type Ticket struct {
	Id              int
	Subject         string
	DescriptionText string
	Tags            []string
	Conversation    []*ConversationItem
}

type ConversationSource int

const (
	ConversationSourceUnknown = iota
	ConversationSourceReply
	ConversationSourceNote
	ConversationSourceForwardedEmail
)

type ConversationItem struct {
	Id              int
	DescriptionText string
	Source          ConversationSource
}

type conversationResponse struct {
	Id       int    `json:"id"`
	BodyText string `json:"body_text"`
	Source   int    `json:"source"`
}

type ticketResponse struct {
	Id              int                     `json:"id"`
	Subject         string                  `json:"subject"`
	DescriptionText string                  `json:"description_text"`
	Tags            []string                `json:"tags"`
	Conversations   []*conversationResponse `json:"conversations"`
	/*
		{
		  "cc_emails" : ["user@cc.com"],
		  "fwd_emails" : [ ],
		  "reply_cc_emails" : ["user@cc.com"],
		  "email_config_id" : null,
		  "fr_escalated" : false,
		  "group_id" : null,
		  "priority" : 1,
		  "requester_id" : 1,
		  "responder_id" : null,
		  "source" : 2,
		  "spam" : false,
		  "status" : 2,
		  "subject" : "",
		  "company_id" : 1,
		  "id" : 20,
		  "type" : null,
		  "to_emails" : null,
		  "product_id" : null,
		  "created_at" : "2015-08-24T11:56:51Z",
		  "updated_at" : "2015-08-24T11:59:05Z",
		  "due_by" : "2015-08-27T11:30:00Z",
		  "fr_due_by" : "2015-08-25T11:30:00Z",
		  "is_escalated" : false,
		  "association_type" : null,
		  "description_text" : "Not given.",
		  "description" : "<div>Not given.</div>",
		  "custom_fields" : {
		    "category" : "Primary"
		  },
		  "tags" : [ ],
		  "attachments" : [ ]
		}
	*/
}

func (c *Client) GetTicket(id int) (*Ticket, error) {
	url := fmt.Sprintf("/api/v2/tickets/%d?include=conversations", id)

	res, err := c.do(&request{
		method: "GET",
		url:    url,
	})
	if err != nil {
		return nil, err
	}

	var body ticketResponse
	err = json.Unmarshal(res.body, &body)
	if err != nil {
		return nil, err
	}

	return mapTicketResponse(&body), nil
}

func mapTicketResponse(res *ticketResponse) *Ticket {
	var conversation []*ConversationItem
	for _, c := range res.Conversations {
		conversation = append(conversation, mapConversationResponse(c))
	}
	return &Ticket{
		Id:              res.Id,
		Subject:         res.Subject,
		DescriptionText: res.DescriptionText,
		Tags:            res.Tags,
		Conversation:    conversation,
	}
}

func mapConversationResponse(res *conversationResponse) *ConversationItem {
	return &ConversationItem{
		Id:              res.Id,
		DescriptionText: res.BodyText,
		Source:          mapConversationSource(res.Source),
	}
}

func mapConversationSource(source int) ConversationSource {
	switch source {
	case 0:
		return ConversationSourceReply
	case 2:
		return ConversationSourceNote
	case 8:
		return ConversationSourceForwardedEmail
	}
	return ConversationSourceUnknown
}
