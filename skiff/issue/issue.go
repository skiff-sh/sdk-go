package issue

import (
	"fmt"

	"github.com/skiff-sh/api/go/skiff/plugin/v1alpha1"
)

var _ error = (*Issue)(nil)
var _ PluginIssue = (*Issue)(nil)

// PluginIssue is a generic interface to surface issues to the user of this plugin.
type PluginIssue interface {
	Issue() *v1alpha1.Issue
}

// Issue represents issues the plugin is surfacing to the user. You can join multiple Issue via errors.Join.
type Issue struct {
	IssueLevel v1alpha1.IssueLevel
	Message    string
}

func (i *Issue) Issue() *v1alpha1.Issue {
	if i == nil {
		return nil
	}
	return &v1alpha1.Issue{
		Level:   i.IssueLevel,
		Message: i.Message,
	}
}

func (i *Issue) Error() string {
	if i == nil {
		return ""
	}
	return i.Message
}

// Error returns an error level Issue.
func Error(msg string) *Issue {
	return &Issue{
		IssueLevel: v1alpha1.IssueLevel_LEVEL_ERROR,
		Message:    msg,
	}
}

// Errorf same as Error but with string formatting.
func Errorf(msg string, args ...any) *Issue {
	return &Issue{
		IssueLevel: v1alpha1.IssueLevel_LEVEL_ERROR,
		Message:    fmt.Sprintf(msg, args...),
	}
}

// Warn returns an error level Issue.
func Warn(msg string) *Issue {
	return &Issue{
		IssueLevel: v1alpha1.IssueLevel_LEVEL_WARN,
		Message:    msg,
	}
}

// Warnf same as Warn but with string formatting.
func Warnf(msg string, args ...any) *Issue {
	return &Issue{
		IssueLevel: v1alpha1.IssueLevel_LEVEL_WARN,
		Message:    fmt.Sprintf(msg, args...),
	}
}
