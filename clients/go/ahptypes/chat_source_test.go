package ahptypes

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestChatSourceRoutesByKind(t *testing.T) {
	t.Run("fork", func(t *testing.T) {
		var value ChatSource
		if err := json.Unmarshal([]byte(`{"kind":"fork","chat":"ahp-chat:/main","turnId":"turn-12"}`), &value); err != nil {
			t.Fatalf("decode fork: %v", err)
		}
		fork, ok := value.Value.(*ForkChatSource)
		if !ok {
			t.Fatalf("expected *ForkChatSource, got %T", value.Value)
		}
		if fork.Kind != ChatSourceKindFork || fork.TurnId != "turn-12" {
			t.Fatalf("unexpected fork payload: %#v", fork)
		}
	})

	t.Run("sideChat", func(t *testing.T) {
		var value ChatSource
		if err := json.Unmarshal([]byte(`{"kind":"sideChat","chat":"ahp-chat:/main","turnId":"turn-active"}`), &value); err != nil {
			t.Fatalf("decode sideChat: %v", err)
		}
		sideChat, ok := value.Value.(*SideChatSource)
		if !ok {
			t.Fatalf("expected *SideChatSource, got %T", value.Value)
		}
		if sideChat.Kind != ChatSourceKindSideChat || sideChat.TurnId != "turn-active" {
			t.Fatalf("unexpected sideChat payload: %#v", sideChat)
		}
	})
}

func TestChatSourceRejectsMissingOrUnknownKind(t *testing.T) {
	for name, raw := range map[string]string{
		"missing kind": `{"chat":"ahp-chat:/main","turnId":"turn-12"}`,
		"unknown kind": `{"kind":"future","chat":"ahp-chat:/main","turnId":"turn-12"}`,
	} {
		t.Run(name, func(t *testing.T) {
			var value ChatSource
			err := json.Unmarshal([]byte(raw), &value)
			if err == nil {
				t.Fatalf("expected decode failure for %s", name)
			}
			if got := fmt.Sprint(err); got == "" {
				t.Fatalf("expected printable error for %s", name)
			}
		})
	}
}
