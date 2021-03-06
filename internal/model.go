package internal

import (
	"encoding/json"
)

const refType = "jira"

type user struct {
	AccountID    string  `json:"accountId"` // AccountID not available in hosted jira.
	Self         string  `json:"self"`
	Name         string  `json:"name"` // Name not available from jira cloud.
	Key          string  `json:"key"`
	EmailAddress string  `json:"emailAddress"` // EmailAddress not available from jira cloud.
	Avatars      Avatars `json:"avatarUrls"`
	DisplayName  string  `json:"displayName"`
	Active       bool    `json:"active"`
	Timezone     string  `json:"timeZone"`

	Groups struct {
		Groups []userGroup `json:"items,omitempty"`
	} `json:"groups"`
}

type userGroup struct {
	Name string `json:"name,omitempty"`
}

func (s user) IsZero() bool {
	return s.RefID() == ""
}

func (s user) RefID() string {
	if s.AccountID != "" {
		return s.AccountID
	}
	return s.Key
}

// Avatars is a type that describes a set of avatar image properties
type Avatars struct {
	XSmall string `json:"16x16"`
	Small  string `json:"24x24"`
	Medium string `json:"32x32"`
	Large  string `json:"48x48"`
}

type changeLogItem struct {
	Field      string `json:"field"`
	FieldType  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}

type issueSource struct {
	ID  string `json:"id"`
	Key string `json:"key"`

	// Using map here instead of the Fields struct declared below,
	// since we extract custom fields which could have keys prefixed
	// with customfield_.
	Fields    map[string]interface{} `json:"fields"`
	Changelog struct {
		Histories []struct {
			ID      string          `json:"id"`
			Author  user            `json:"author"`
			Created string          `json:"created"`
			Items   []changeLogItem `json:"items"`
		} `json:"histories"`
	} `json:"changelog"`
	Transitions []transitionSource `json:"transitions"`
}

type transitionSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   struct {
		StatusCategory statusCategory `json:"statusCategory"`
	} `json:"to"`
}

type statusCategory struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

const (
	statusCategoryNew          = "new"
	statusCategoryIntermediate = "indeterminate"
	statusCategoryDone         = "done"
)

type linkedIssue struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type issueFields struct {
	Project struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	} `json:"project"`
	Description json.RawMessage `json:"description"`
	Comment     struct {
		Comments []comment
	} `json:"comment"`
	Summary string `json:"summary"`
	DueDate string `json:"duedate"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Parent  *struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	} `json:"parent,omitempty"`
	Priority struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"priority"`
	IssueType struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"issuetype"`
	Status struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"status"`
	Resolution struct {
		Name string `json:"name"`
	} `json:"resolution"`
	Creator    user
	Reporter   user
	Assignee   user
	Labels     []string `json:"labels"`
	IssueLinks []struct {
		ID   string `json:"id"`
		Type struct {
			//ID   string `json:"id"`
			Name string `json:"name"` // Using Name instead of ID for mapping
		} `json:"type"`
		OutwardIssue linkedIssue `json:"outwardIssue"`
		InwardIssue  linkedIssue `json:"inwardIssue"`
	} `json:"issuelinks"`
	Attachment []struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
		Author   struct {
			Key       string `json:"key"`
			AccountID string `json:"accountId"`
		} `json:"author"`
		Created   string `json:"created"`
		Size      int    `json:"size"`
		MimeType  string `json:"mimeType"`
		Content   string `json:"content"`
		Thumbnail string `json:"thumbnail"`
	} `json:"attachment"`
}

type issueQueryResult struct {
	Total  int           `json:"total"`
	Issues []issueSource `json:"issues"`
}

type project struct {
	Expand      string `json:"expand"`
	Self        string `json:"self"`
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
	IssueTypes  []struct {
		Self        string `json:"self"`
		ID          string `json:"id"`
		Description string `json:"description"`
		IconURL     string `json:"iconUrl"`
		Name        string `json:"name"`
		Subtask     bool   `json:"subtask"`
		AvatarID    int    `json:"avatarId,omitempty"`
	} `json:"issueTypes"`
	Name       string `json:"name"`
	AvatarUrls struct {
		Four8X48  string `json:"48x48"`
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
	} `json:"avatarUrls"`
	ProjectKeys     []string `json:"projectKeys"`
	ProjectCategory struct {
		Self        string `json:"self"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"projectCategory,omitempty"`
	ProjectTypeKey string `json:"projectTypeKey"`
	Simplified     bool   `json:"simplified"`
	Style          string `json:"style"`
	IsPrivate      bool   `json:"isPrivate"`
	EntityID       string `json:"entityId,omitempty"`
	UUID           string `json:"uuid,omitempty"`
	Insight        *struct {
		TotalIssueCount     int    `json:"totalIssueCount"`
		LastIssueUpdateTime string `json:"lastIssueUpdateTime"`
	} `json:"insight,omitempty"`
}

type projectQueryResult struct {
	Total    int       `json:"total"`
	Projects []project `json:"values"`
}

type issuePriority struct {
	StatusColor string `json:"statusColor"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

type issueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"iconUrl"`
	Subtask     bool   `json:"subtask"`
}

type customFieldQueryResult struct {
	ID   string `json:"id"`
	Key  string `json:"key"` // this is only on cloud and not server
	Name string `json:"name"`
}

type issueCreateMeta struct {
	// Expand   string `json:"expand"`
	Projects []projectIssueCreateMeta `json:"projects"`
}

type projectIssueCreateMeta struct {
	ID         string                 `json:"id"`
	Key        string                 `json:"key"`
	Name       string                 `json:"name"`
	Issuetypes []createMetaIssueTypes `json:"issuetypes"`
}

type createMetaIssueTypes struct {
	ID               string                    `json:"id"`
	Description      string                    `json:"description"`
	IconURL          string                    `json:"iconUrl"`
	Name             string                    `json:"name"`
	UntranslatedName string                    `json:"untranslatedName"`
	Subtask          bool                      `json:"subtask"`
	Expand           string                    `json:"expand"`
	Fields           map[string]issueTypeField `json:"fields"`
}

type issueTypeField struct {
	Required        bool                 `json:"required"`
	Schema          issueTypeFieldSchema `json:"schema"`
	Name            string               `json:"name"`
	Key             string               `json:"key"`
	HasDefaultValue bool                 `json:"hasDefaultValue"`
	AllowedValues   json.RawMessage      `json:"allowedValues,omitempty"`
}

type issueTypeFieldSchema struct {
	Type   string `json:"type"`
	System string `json:"system"`
}

type allowedValueComponent struct {
	RefID string `json:"id"`
	Name  string `json:"name"`
}
