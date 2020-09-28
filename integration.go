package main

import (
	"github.com/pinpt/jira/internal"
	"github.com/pinpt/agent/runner"
)

// Integration is used to export the integration
var Integration internal.JiraIntegration

func main() {
	runner.Main(&Integration)
}
