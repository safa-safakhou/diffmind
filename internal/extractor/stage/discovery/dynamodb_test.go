package discovery

import (
	"encoding/json"
	"testing"
)

func TestDynamoDBOperationRetainsAllCallSitesDeterministically(t *testing.T) {
	idx := buildAgentsIndex(t, map[string]string{
		"application.yml": "services:\n  routing:\n    dynamodb-table: traffic-info\n",
		"Controller.java": `class Controller { private DynamoDbTemplate dynamoDbTemplate; void update() { dynamoDbTemplate.save(item); } }`,
		"Job.java":        `class Job { private DynamoDbTemplate dynamoDbTemplate; void refresh() { dynamoDbTemplate.save(item); } }`,
	})
	var baseline string
	for i := 0; i < 100; i++ {
		got := DeterministicDynamoDBOperations(idx)
		if len(got) != 1 || len(got[0].Locations) != 2 || len(got[0].Evidence) != 2 {
			t.Fatalf("call site lost: %+v", got)
		}
		if got[0].Locations[0].File != "Controller.java" || got[0].Locations[1].File != "Job.java" {
			t.Fatal("locations are not sorted")
		}
		raw, _ := json.Marshal(got)
		if i > 0 && string(raw) != baseline {
			t.Fatal("map iteration changes extraction")
		}
		baseline = string(raw)
	}
}

func TestRepositoryOperationRetainsAllCallSitesDeterministically(t *testing.T) {
	idx := buildAgentsIndex(t, map[string]string{
		"Controller.java":      `class Controller { private OrderRepository orderRepository; void first() { orderRepository.findById(id); } }`,
		"Job.java":             `class Job { private OrderRepository orderRepository; void second() { orderRepository.findById(id); } }`,
		"OrderRepository.java": `interface OrderRepository { Order findById(String id); }`,
	})
	var baseline string
	for i := 0; i < 100; i++ {
		got := DeterministicDBOperations(idx)
		if len(got) != 1 || len(got[0].Locations) != 2 || len(got[0].Evidence) != 2 {
			t.Fatalf("call site lost: %+v", got)
		}
		raw, _ := json.Marshal(got)
		if i > 0 && string(raw) != baseline {
			t.Fatal("repository call-site aggregation is nondeterministic")
		}
		baseline = string(raw)
	}
}
