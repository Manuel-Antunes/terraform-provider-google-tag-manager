package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/tagmanager/v2"
)

var testClientInWorkspaceOptions = &ClientInWorkspaceOptions{
	WorkspaceName: "test-client-in-workspace",
	ClientOptions: testClientOptions,
}

type ClientInWorkspaceTestSuite struct {
	suite.Suite
	client *ClientInWorkspace
}

func (suite *ClientInWorkspaceTestSuite) SetupSuite() {
	client, err := NewClientInWorkspace(testClientInWorkspaceOptions)
	if err != nil {
		suite.T().Fatalf("Failed to create client in workspace: %v", err)
	}
	suite.client = client
}

func (suite *ClientInWorkspaceTestSuite) TearDownSuite() {
	if suite.client != nil && suite.client.Options != nil && suite.client.Options.WorkspaceId != "" {
		err := suite.client.DeleteWorkspace(suite.client.Options.WorkspaceId)
		if err != nil {
			suite.T().Errorf("Failed to delete workspace: %v", err)
		}
	}
}

func (suite *ClientInWorkspaceTestSuite) TestNewClientInWorkspace() {
	assert.NotNil(suite.T(), suite.client)
	assert.NotNil(suite.T(), suite.client.Options)
	assert.NotZero(suite.T(), suite.client.Options.WorkspaceId)
}

// Test tag creation
func (suite *ClientInWorkspaceTestSuite) TestTagCreate() {
	t := suite.T()

	tag, err := suite.client.CreateTag(&tagmanager.Tag{
		Name:  testName("test-tag-create"),
		Notes: "created by integration test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "test-event", Type: "template"},
			{Key: "measurementIdOverride", Value: "G-XXXXXX", Type: "template"},
			{Key: "eventParameters", Type: "list", List: []*tagmanager.Parameter{
				{Type: "map", Map: []*tagmanager.Parameter{
					{Key: "name", Type: "template", Value: "param-name"},
					{Key: "value", Type: "template", Value: "param-value"},
				}},
			}},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.NotEmpty(t, tag.TagId)

	// Clean up
	defer func() {
		err := suite.client.DeleteTag(tag.TagId)
		assert.NoError(t, err)
	}()
}

// Test tag read
func (suite *ClientInWorkspaceTestSuite) TestTagRead() {
	t := suite.T()

	// First create a tag to read
	tag, err := suite.client.CreateTag(&tagmanager.Tag{
		Name:  testName("test-tag-read"),
		Notes: "created for read test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "test-event-read", Type: "template"},
			{Key: "measurementIdOverride", Value: "G-XXXXXX", Type: "template"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tag)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteTag(tag.TagId)
		assert.NoError(t, err)
	}()

	// Get tag by ID
	fetchedTag, err := suite.client.Tag(tag.TagId)
	assert.NoError(t, err)
	assert.Equal(t, tag.Name, fetchedTag.Name)
	assert.Equal(t, tag.Type, fetchedTag.Type)
	assert.Equal(t, "test-event-read", fetchedTag.Parameter[0].Value)
}

// Test tag list
func (suite *ClientInWorkspaceTestSuite) TestTagList() {
	t := suite.T()

	// Create a tag first to ensure there's at least one
	tag, err := suite.client.CreateTag(&tagmanager.Tag{
		Name:  testName("test-tag-list"),
		Notes: "created for list test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "test-event-list", Type: "template"},
			{Key: "measurementIdOverride", Value: "G-XXXXXX", Type: "template"},
		},
	})
	assert.NoError(t, err)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteTag(tag.TagId)
		assert.NoError(t, err)
	}()

	// List tags
	tags, err := suite.client.ListTags()
	assert.NoError(t, err)
	assert.Greater(t, len(tags), 0)

	// Find our tag in the list
	found := false
	for _, t := range tags {
		if t.TagId == tag.TagId {
			found = true
			break
		}
	}
	assert.True(t, found, "Created tag was not found in the list")
}

// Test tag update
func (suite *ClientInWorkspaceTestSuite) TestTagUpdate() {
	t := suite.T()

	// First create a tag to update
	tag, err := suite.client.CreateTag(&tagmanager.Tag{
		Name:  testName("test-tag-update"),
		Notes: "created for update test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "original-event", Type: "template"},
			{Key: "measurementIdOverride", Value: "G-XXXXXX", Type: "template"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tag)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteTag(tag.TagId)
		assert.NoError(t, err)
	}()

	// Update tag
	updatedTagName := testName("updated-tag")
	updatedTag, err := suite.client.UpdateTag(tag.TagId, &tagmanager.Tag{
		Name:  updatedTagName,
		Notes: "updated by integration test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "updated-event", Type: "template"},
			{Key: "measurementIdOverride", Value: "G-XXXXXX", Type: "template"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, updatedTagName, updatedTag.Name)
	assert.Equal(t, "updated-event", updatedTag.Parameter[0].Value)

	// Verify update by reading again
	fetchedTag, err := suite.client.Tag(tag.TagId)
	assert.NoError(t, err)
	assert.Equal(t, updatedTagName, fetchedTag.Name)
	assert.Equal(t, "updated-event", fetchedTag.Parameter[0].Value)
}

// Test tag delete
func (suite *ClientInWorkspaceTestSuite) TestTagDelete() {
	t := suite.T()

	// First create a tag to delete
	tag, err := suite.client.CreateTag(&tagmanager.Tag{
		Name:  testName("test-tag-delete"),
		Notes: "created for delete test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "test-event-delete", Type: "template"},
			{Key: "measurementIdOverride", Value: "G-XXXXXX", Type: "template"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tag)

	// Delete tag
	err = suite.client.DeleteTag(tag.TagId)
	assert.NoError(t, err)

	// Verify deletion
	_, err = suite.client.Tag(tag.TagId)
	assert.Equal(t, ErrNotExist, err)
}

// Test variable creation
func (suite *ClientInWorkspaceTestSuite) TestVariableCreate() {
	t := suite.T()

	variable, err := suite.client.CreateVariable(&tagmanager.Variable{
		Name: testName("test-variable-create"),
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "test-param"},
			{Key: "value", Type: "template", Value: "test-value"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, variable)
	assert.NotEmpty(t, variable.VariableId)

	// Clean up
	defer func() {
		err := suite.client.DeleteVariable(variable.VariableId)
		assert.NoError(t, err)
	}()
}

// Test variable read
func (suite *ClientInWorkspaceTestSuite) TestVariableRead() {
	t := suite.T()

	// First create a variable to read
	variable, err := suite.client.CreateVariable(&tagmanager.Variable{
		Name: testName("test-variable-read"),
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "read-param"},
			{Key: "value", Type: "template", Value: "read-value"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, variable)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteVariable(variable.VariableId)
		assert.NoError(t, err)
	}()

	// Get variable by ID
	fetchedVariable, err := suite.client.Variable(variable.VariableId)
	assert.NoError(t, err)
	assert.Equal(t, variable.Name, fetchedVariable.Name)
	assert.Equal(t, variable.Type, fetchedVariable.Type)
	assert.Equal(t, "read-param", fetchedVariable.Parameter[0].Value)
}

// Test variable list
func (suite *ClientInWorkspaceTestSuite) TestVariableList() {
	t := suite.T()

	// Create a variable first to ensure there's at least one
	variable, err := suite.client.CreateVariable(&tagmanager.Variable{
		Name: testName("test-variable-list"),
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "list-param"},
			{Key: "value", Type: "template", Value: "list-value"},
		},
	})
	assert.NoError(t, err)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteVariable(variable.VariableId)
		assert.NoError(t, err)
	}()

	// List variables
	variables, err := suite.client.ListVariables()
	assert.NoError(t, err)
	assert.Greater(t, len(variables), 0)

	// Find our variable in the list
	found := false
	for _, v := range variables {
		if v.VariableId == variable.VariableId {
			found = true
			break
		}
	}
	assert.True(t, found, "Created variable was not found in the list")
}

// Test variable update
func (suite *ClientInWorkspaceTestSuite) TestVariableUpdate() {
	t := suite.T()

	// First create a variable to update
	variable, err := suite.client.CreateVariable(&tagmanager.Variable{
		Name: testName("test-variable-update"),
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "original-param"},
			{Key: "value", Type: "template", Value: "original-value"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, variable)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteVariable(variable.VariableId)
		assert.NoError(t, err)
	}()

	// Update variable
	updatedVariableName := testName("updated-variable")
	updatedVariable, err := suite.client.UpdateVariable(variable.VariableId, &tagmanager.Variable{
		Name: updatedVariableName,
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "updated-param"},
			{Key: "value", Type: "template", Value: "updated-value"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, updatedVariableName, updatedVariable.Name)
	assert.Equal(t, "updated-param", updatedVariable.Parameter[0].Value)

	// Verify update by reading again
	fetchedVariable, err := suite.client.Variable(variable.VariableId)
	assert.NoError(t, err)
	assert.Equal(t, updatedVariableName, fetchedVariable.Name)
	assert.Equal(t, "updated-param", fetchedVariable.Parameter[0].Value)
}

// Test variable delete
func (suite *ClientInWorkspaceTestSuite) TestVariableDelete() {
	t := suite.T()

	// First create a variable to delete
	variable, err := suite.client.CreateVariable(&tagmanager.Variable{
		Name: testName("test-variable-delete"),
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "delete-param"},
			{Key: "value", Type: "template", Value: "delete-value"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, variable)

	// Delete variable
	err = suite.client.DeleteVariable(variable.VariableId)
	assert.NoError(t, err)

	// Verify deletion
	_, err = suite.client.Variable(variable.VariableId)
	assert.Equal(t, ErrNotExist, err)
}

// Test trigger creation
func (suite *ClientInWorkspaceTestSuite) TestTriggerCreate() {
	t := suite.T()

	trigger, err := suite.client.CreateTrigger(&tagmanager.Trigger{
		Name:  testName("test-trigger-create"),
		Type:  "customEvent",
		Notes: "Created by integration test",
		CustomEventFilter: []*tagmanager.Condition{
			{
				Type: "equals",
				Parameter: []*tagmanager.Parameter{
					{Key: "arg0", Value: "{{_event}}", Type: "template"},
					{Key: "arg1", Value: "test-event-create", Type: "template"},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, trigger)
	assert.NotEmpty(t, trigger.TriggerId)

	// Clean up
	defer func() {
		err := suite.client.DeleteTrigger(trigger.TriggerId)
		assert.NoError(t, err)
	}()
}

// Test trigger read
func (suite *ClientInWorkspaceTestSuite) TestTriggerRead() {
	t := suite.T()

	// First create a trigger to read
	trigger, err := suite.client.CreateTrigger(&tagmanager.Trigger{
		Name:  testName("test-trigger-read"),
		Type:  "customEvent",
		Notes: "Created for read test",
		CustomEventFilter: []*tagmanager.Condition{
			{
				Type: "equals",
				Parameter: []*tagmanager.Parameter{
					{Key: "arg0", Value: "{{_event}}", Type: "template"},
					{Key: "arg1", Value: "test-event-read", Type: "template"},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, trigger)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteTrigger(trigger.TriggerId)
		assert.NoError(t, err)
	}()

	// Get trigger by ID
	fetchedTrigger, err := suite.client.Trigger(trigger.TriggerId)
	assert.NoError(t, err)
	assert.Equal(t, trigger.Name, fetchedTrigger.Name)
	assert.Equal(t, trigger.Type, fetchedTrigger.Type)
	assert.Equal(t, "test-event-read", fetchedTrigger.CustomEventFilter[0].Parameter[1].Value)
}

// Test trigger list
func (suite *ClientInWorkspaceTestSuite) TestTriggerList() {
	t := suite.T()

	// Create a trigger first to ensure there's at least one
	trigger, err := suite.client.CreateTrigger(&tagmanager.Trigger{
		Name:  testName("test-trigger-list"),
		Type:  "customEvent",
		Notes: "Created for list test",
		CustomEventFilter: []*tagmanager.Condition{
			{
				Type: "equals",
				Parameter: []*tagmanager.Parameter{
					{Key: "arg0", Value: "{{_event}}", Type: "template"},
					{Key: "arg1", Value: "test-event-list", Type: "template"},
				},
			},
		},
	})
	assert.NoError(t, err)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteTrigger(trigger.TriggerId)
		assert.NoError(t, err)
	}()

	// List triggers
	triggers, err := suite.client.ListTriggers()
	assert.NoError(t, err)
	assert.Greater(t, len(triggers), 0)

	// Find our trigger in the list
	found := false
	for _, t := range triggers {
		if t.TriggerId == trigger.TriggerId {
			found = true
			break
		}
	}
	assert.True(t, found, "Created trigger was not found in the list")
}

// Test trigger update
func (suite *ClientInWorkspaceTestSuite) TestTriggerUpdate() {
	t := suite.T()

	// First create a trigger to update
	trigger, err := suite.client.CreateTrigger(&tagmanager.Trigger{
		Name:  testName("test-trigger-update"),
		Type:  "customEvent",
		Notes: "Created for update test",
		CustomEventFilter: []*tagmanager.Condition{
			{
				Type: "equals",
				Parameter: []*tagmanager.Parameter{
					{Key: "arg0", Value: "{{_event}}", Type: "template"},
					{Key: "arg1", Value: "original-event", Type: "template"},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, trigger)

	// Clean up when done
	defer func() {
		err := suite.client.DeleteTrigger(trigger.TriggerId)
		assert.NoError(t, err)
	}()

	// Update trigger
	updatedTriggerName := testName("updated-trigger")
	updatedTrigger, err := suite.client.UpdateTrigger(trigger.TriggerId, &tagmanager.Trigger{
		Name:  updatedTriggerName,
		Type:  "click",
		Notes: "Updated by integration test",
		Parameter: []*tagmanager.Parameter{
			{Key: "clickText", Value: "UpdatedButton", Type: "template"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, updatedTriggerName, updatedTrigger.Name)
	assert.Equal(t, "click", updatedTrigger.Type)

	// Verify update by reading again
	fetchedTrigger, err := suite.client.Trigger(trigger.TriggerId)
	assert.NoError(t, err)
	assert.Equal(t, updatedTriggerName, fetchedTrigger.Name)
	assert.Equal(t, "click", fetchedTrigger.Type)
}

// Test trigger delete
func (suite *ClientInWorkspaceTestSuite) TestTriggerDelete() {
	t := suite.T()

	// First create a trigger to delete
	trigger, err := suite.client.CreateTrigger(&tagmanager.Trigger{
		Name:  testName("test-trigger-delete"),
		Type:  "customEvent",
		Notes: "Created for delete test",
		CustomEventFilter: []*tagmanager.Condition{
			{
				Type: "equals",
				Parameter: []*tagmanager.Parameter{
					{Key: "arg0", Value: "{{_event}}", Type: "template"},
					{Key: "arg1", Value: "test-event-delete", Type: "template"},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, trigger)

	// Delete trigger
	err = suite.client.DeleteTrigger(trigger.TriggerId)
	assert.NoError(t, err)

	// Verify deletion
	_, err = suite.client.Trigger(trigger.TriggerId)
	assert.Equal(t, ErrNotExist, err)
}

// Helper function to generate a unique test name based on current time
func testName(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102-150405")
}

func TestClientInWorkspace(t *testing.T) {
	suite.Run(t, new(ClientInWorkspaceTestSuite))
}
