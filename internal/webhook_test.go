package internal

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/pinpt/agent.next/sdk"
	"github.com/pinpt/agent.next/sdk/sdktest"
	"github.com/pinpt/integration-sdk/agent"
	"github.com/stretchr/testify/assert"
)

func loadFile(fn string) []byte {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic("error reading test data file: " + err.Error())
	}
	return b
}

func TestWebhookJiraIssueDeleted(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookDeleteIssue("1234", "1", loadFile("testdata/jira:issue_deleted.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("false", update.Set["active"])
}

func TestWebhookJiraIssueCommentDeleted(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookDeleteComment("1234", "1", loadFile("testdata/comment_deleted.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("false", update.Set["active"])
}

func TestWebhookJiraProjectDeleted(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookDeleteProject("1234", "1", loadFile("testdata/project_deleted.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("false", update.Set["active"])
}

func TestWebhookJiraUserCreated(t *testing.T) {
	assert := assert.New(t)
	logger := sdk.NewNoOpTestLogger()
	um := &mockUserManager{}
	assert.NoError(webhookUpsertUser(logger, um, loadFile("testdata/user_created.json")))
	assert.Len(um.users, 1)
	assert.EqualValues("jhaynie+1", um.users[0].DisplayName)
	assert.EqualValues("5f03c8345ee2c300232945de", um.users[0].AccountID)
}

func TestWebhookJiraUserUpdated(t *testing.T) {
	assert := assert.New(t)
	logger := sdk.NewNoOpTestLogger()
	um := &mockUserManager{}
	assert.NoError(webhookUpsertUser(logger, um, loadFile("testdata/user_updated.json")))
	assert.Len(um.users, 1)
	assert.EqualValues("jeff haynie test", um.users[0].DisplayName)
	assert.EqualValues("5f03c8345ee2c300232945de", um.users[0].AccountID)
}

func TestWebhookJiraSprintDeleted(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookDeleteSprint("1234", "1", loadFile("testdata/sprint_deleted.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("false", update.Set["active"])
}

const sprintUpdateGoalAdded = `{
  "timestamp": 1596142226397,
  "webhookEvent": "sprint_updated",
  "sprint": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "future",
    "name": "TES Sprint 2",
    "originBoardId": 9,
    "goal": "take over the world 🌍🌎🌏"
  },
  "oldValue": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "future",
    "name": "TES Sprint 2",
    "originBoardId": 9
  }
}`

const sprintUpdateGoalUpdated = `{
  "timestamp": 1596142259617,
  "webhookEvent": "sprint_updated",
  "sprint": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "future",
    "name": "TES Sprint 2",
    "originBoardId": 9,
    "goal": "take over the world! 🌍🌎🌏"
  },
  "oldValue": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "future",
    "name": "TES Sprint 2",
    "originBoardId": 9,
    "goal": "take over the world 🌍🌎🌏"
  }
}
`

const sprintUpdatedEndDate = `{
  "timestamp": 1596144049198,
  "webhookEvent": "sprint_updated",
  "sprint": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "active",
    "name": "TES Sprint 2",
    "startDate": "2020-07-30T21:13:24.588Z",
    "endDate": "2020-08-14T21:13:00.000Z",
    "originBoardId": 9,
    "goal": "take over the world! 🌍🌎🌏"
  },
  "oldValue": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "active",
    "name": "TES Sprint 2",
    "startDate": "2020-07-30T21:13:24.588Z",
    "endDate": "2020-08-13T21:13:00.000Z",
    "originBoardId": 9,
    "goal": "take over the world! 🌍🌎🌏"
  }
}`

func TestWebhookBuildSprintUpdateChangeName(t *testing.T) {
	assert := assert.New(t)
	val, changed := buildSprintUpdate(sprintProjection{
		ID:   1,
		Name: "Sprint 1",
	}, sprintProjection{
		ID:   1,
		Name: "Sprint! 1",
	})
	assert.True(changed)
	assert.NotNil(val.Set.Name)
	assert.EqualValues("Sprint! 1", *val.Set.Name)
}

func TestWebhookBuildSprintUpdateChangeGoal(t *testing.T) {
	assert := assert.New(t)
	val, changed := buildSprintUpdate(sprintProjection{
		ID:   1,
		Name: "Sprint 1",
		Goal: sdk.StringPointer("hello"),
	}, sprintProjection{
		ID:   1,
		Name: "Sprint 1",
		Goal: sdk.StringPointer("hello world!"),
	})
	assert.True(changed)
	assert.Nil(val.Set.Name)
	assert.NotNil(val.Set.Goal)
	assert.EqualValues("hello world!", *val.Set.Goal)
}

func TestWebhookJiraSprintUpdateStarted(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookUpdateSprint("1234", "1", loadFile("testdata/sprint_updated.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("", update.Set["active"])
	assert.EqualValues("\"ACTIVE\"", update.Set["status"])
	assert.EqualValues("{\"epoch\":1596143604588,\"offset\":0,\"rfc3339\":\"2020-07-30T21:13:24.588+00:00\"}", update.Set["started_date"])
	assert.EqualValues("{\"epoch\":1597353180000,\"offset\":0,\"rfc3339\":\"2020-08-13T21:13:00+00:00\"}", update.Set["ended_date"])
}

func TestWebhookJiraSprintUpdateGoalSet(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookUpdateSprint("1234", "1", []byte(sprintUpdateGoalAdded), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("", update.Set["status"])
	assert.EqualValues("", update.Set["name"])
	assert.EqualValues("\"take over the world 🌍🌎🌏\"", update.Set["goal"])
}

func TestWebhookJiraSprintUpdateGoalUpdated(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookUpdateSprint("1234", "1", []byte(sprintUpdateGoalUpdated), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("", update.Set["status"])
	assert.EqualValues("", update.Set["name"])
	assert.EqualValues("", update.Set["ended_date"])
	assert.EqualValues("", update.Set["started_date"])
	assert.EqualValues("\"take over the world! 🌍🌎🌏\"", update.Set["goal"])
}

func TestWebhookJiraSprintUpdateEndDate(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookUpdateSprint("1234", "1", []byte(sprintUpdatedEndDate), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("", update.Set["status"])
	assert.EqualValues("", update.Set["name"])
	assert.EqualValues("", update.Set["goal"])
	assert.EqualValues("{\"epoch\":1597439580000,\"offset\":0,\"rfc3339\":\"2020-08-14T21:13:00+00:00\"}", update.Set["ended_date"])
	assert.EqualValues("", update.Set["started_date"])
}

const sprintUpdateNothing = `{
  "timestamp": 1596142226397,
  "webhookEvent": "sprint_updated",
  "sprint": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "future",
    "name": "TES Sprint 2",
    "originBoardId": 9
  },
  "oldValue": {
    "id": 196,
    "self": "https://pinpt-hq.atlassian.net/rest/agile/1.0/sprint/196",
    "state": "future",
    "name": "TES Sprint 2",
    "originBoardId": 9
  }
}`

func TestWebhookJiraSprintUpdateNothing(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookUpdateSprint("1234", "1", []byte(sprintUpdateNothing), pipe))
	assert.Len(pipe.Written, 0)
}

func TestWebhookJiraSprintClosed(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookCloseSprint("1234", "1", loadFile("testdata/sprint_closed.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("\"CLOSED\"", update.Set["status"])
	assert.EqualValues("", update.Set["name"])
	assert.EqualValues("", update.Set["ended_date"])
	assert.EqualValues("", update.Set["started_date"])
	assert.EqualValues("{\"epoch\":1594082275692,\"offset\":0,\"rfc3339\":\"2020-07-07T00:37:55.692+00:00\"}", update.Set["completed_date"])
	assert.EqualValues("", update.Set["goal"])
}

func TestWebhookBoardUpdated(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookUpdateBoard("1234", "1", loadFile("testdata/board_updated.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("\"Teamoji Board (updated)\"", update.Set["name"])
}

func TestWebhookBoardDeleted(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookDeleteBoard("1234", "1", loadFile("testdata/board_deleted.json"), pipe))
	assert.Len(pipe.Written, 1)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.EqualValues("false", update.Set["active"])
}

func TestWebhookCreateLinkedIssueBlocks(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookIssueLinkCreated("1234", "1", loadFile("testdata/issuelink_created.json"), pipe))
	assert.Len(pipe.Written, 2)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("af83c065adcd9a05", update.ID)
	var res []sdk.WorkIssueLinkedIssues
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("a5539aea796c83ed", res[0].IssueID)
	assert.EqualValues("20734", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeBlocks, res[0].LinkType)
	assert.EqualValues("22768", res[0].RefID)
	assert.EqualValues(false, res[0].ReverseDirection)

	update = pipe.Written[1].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("a5539aea796c83ed", update.ID)
	res = nil
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("af83c065adcd9a05", res[0].IssueID)
	assert.EqualValues("20192", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeBlocks, res[0].LinkType)
	assert.EqualValues("22768", res[0].RefID)
	assert.EqualValues(true, res[0].ReverseDirection)
}

const dupLink = `{
  "timestamp": 1596481635907,
  "webhookEvent": "issuelink_created",
  "issueLink": {
    "id": 23161,
    "sourceIssueId": 18715,
    "destinationIssueId": 11917,
    "issueLinkType": {
      "id": 10002,
      "name": "Duplicate",
      "outwardName": "duplicates",
      "inwardName": "is duplicated by",
      "isSubTaskLinkType": false,
      "isSystemLinkType": false
    },
    "systemLink": false
  }
}`

func TestWebhookCreateLinkedIssueDuplicates(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookIssueLinkCreated("1234", "1", []byte(dupLink), pipe))
	assert.Len(pipe.Written, 2)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("0d3454a12c41b1d4", update.ID)

	var res []sdk.WorkIssueLinkedIssues
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("91135726a7b2592f", res[0].IssueID)
	assert.EqualValues("11917", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeDuplicates, res[0].LinkType)
	assert.EqualValues("23161", res[0].RefID)
	assert.EqualValues(false, res[0].ReverseDirection)

	update = pipe.Written[1].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("91135726a7b2592f", update.ID)
	res = nil
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("0d3454a12c41b1d4", res[0].IssueID)
	assert.EqualValues("18715", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeDuplicates, res[0].LinkType)
	assert.EqualValues("23161", res[0].RefID)
	assert.EqualValues(true, res[0].ReverseDirection)
}

const cloneLink = `{
  "timestamp": 1596481538074,
  "webhookEvent": "issuelink_created",
  "issueLink": {
    "id": 23160,
    "sourceIssueId": 18715,
    "destinationIssueId": 11917,
    "issueLinkType": {
      "id": 10001,
      "name": "Cloners",
      "outwardName": "clones",
      "inwardName": "is cloned by",
      "isSubTaskLinkType": false,
      "isSystemLinkType": false
    },
    "systemLink": false
  }
}`

func TestWebhookCreateLinkedIssueClones(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookIssueLinkCreated("1234", "1", []byte(cloneLink), pipe))
	assert.Len(pipe.Written, 2)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("0d3454a12c41b1d4", update.ID)

	var res []sdk.WorkIssueLinkedIssues
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("91135726a7b2592f", res[0].IssueID)
	assert.EqualValues("11917", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeClones, res[0].LinkType)
	assert.EqualValues("23160", res[0].RefID)
	assert.EqualValues(false, res[0].ReverseDirection)

	update = pipe.Written[1].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("91135726a7b2592f", update.ID)
	res = nil
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("0d3454a12c41b1d4", res[0].IssueID)
	assert.EqualValues("18715", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeClones, res[0].LinkType)
	assert.EqualValues("23160", res[0].RefID)
	assert.EqualValues(true, res[0].ReverseDirection)
}

const relatesLink = `{
  "timestamp": 1596476927095,
  "webhookEvent": "issuelink_created",
  "issueLink": {
    "id": 23156,
    "sourceIssueId": 18715,
    "destinationIssueId": 11917,
    "issueLinkType": {
      "id": 10003,
      "name": "Relates",
      "outwardName": "relates to",
      "inwardName": "relates to",
      "isSubTaskLinkType": false,
      "isSystemLinkType": false
    },
    "systemLink": false
  }
}`

func TestWebhookCreateLinkedIssueRelates(t *testing.T) {
	assert := assert.New(t)
	pipe := &sdktest.MockPipe{}
	i := JiraIntegration{
		logger: sdk.NewNoOpTestLogger(),
	}
	assert.NoError(i.webhookIssueLinkCreated("1234", "1", []byte(relatesLink), pipe))
	assert.Len(pipe.Written, 2)
	update := pipe.Written[0].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("0d3454a12c41b1d4", update.ID)

	var res []sdk.WorkIssueLinkedIssues
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("91135726a7b2592f", res[0].IssueID)
	assert.EqualValues("11917", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeRelates, res[0].LinkType)
	assert.EqualValues("23156", res[0].RefID)
	assert.EqualValues(false, res[0].ReverseDirection)

	update = pipe.Written[1].(*agent.UpdateData)
	assert.Len(update.Unset, 0)
	assert.EqualValues("91135726a7b2592f", update.ID)
	res = nil
	json.Unmarshal([]byte(update.Push["linked_issues"]), &res)
	assert.Len(res, 1)
	assert.EqualValues("0d3454a12c41b1d4", res[0].IssueID)
	assert.EqualValues("18715", res[0].IssueRefID)
	assert.EqualValues(sdk.WorkIssueLinkedIssuesLinkTypeRelates, res[0].LinkType)
	assert.EqualValues("23156", res[0].RefID)
	assert.EqualValues(true, res[0].ReverseDirection)
}
