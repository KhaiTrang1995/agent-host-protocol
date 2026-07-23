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
		if err := json.Unmarshal([]byte(`{"kind":"sideChat","chat":"ahp-chat:/main","turnId":"turn-active","selection":{"text":"const value = compute()","responsePartId":"part-7"}}`), &value); err != nil {
			t.Fatalf("decode sideChat: %v", err)
		}
		sideChat, ok := value.Value.(*SideChatSource)
		if !ok {
			t.Fatalf("expected *SideChatSource, got %T", value.Value)
		}
		if sideChat.Kind != ChatSourceKindSideChat || sideChat.TurnId != "turn-active" {
			t.Fatalf("unexpected sideChat payload: %#v", sideChat)
		}
		if sideChat.Selection == nil || sideChat.Selection.Text != "const value = compute()" || sideChat.Selection.ResponsePartId == nil || *sideChat.Selection.ResponsePartId != "part-7" {
			t.Fatalf("unexpected sideChat selection: %#v", sideChat.Selection)
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

func TestChatSourceSerializationForcesExactKinds(t *testing.T) {
	t.Run("fork branch and union ignore contradictory kind", func(t *testing.T) {
		branch := ForkChatSource{
			Kind:   ChatSourceKindSideChat,
			Chat:   "ahp-chat:/main",
			TurnId: "turn-12",
		}

		branchRaw, err := json.Marshal(branch)
		if err != nil {
			t.Fatalf("marshal fork branch: %v", err)
		}
		unionRaw, err := json.Marshal(ChatSource{Value: &branch})
		if err != nil {
			t.Fatalf("marshal fork union: %v", err)
		}

		for _, raw := range [][]byte{branchRaw, unionRaw} {
			var decoded map[string]any
			if err := json.Unmarshal(raw, &decoded); err != nil {
				t.Fatalf("decode serialized fork payload: %v", err)
			}
			if got := decoded["kind"]; got != "fork" {
				t.Fatalf("expected fork kind, got %#v in %s", got, raw)
			}
		}
	})

	t.Run("sideChat branch and union ignore zero kind", func(t *testing.T) {
		branch := SideChatSource{
			Chat:   "ahp-chat:/main",
			TurnId: "turn-active",
			Selection: &SideChatSelection{
				Text:           "const value = compute()",
				ResponsePartId: stringPtr("part-7"),
			},
		}

		branchRaw, err := json.Marshal(branch)
		if err != nil {
			t.Fatalf("marshal sideChat branch: %v", err)
		}
		unionRaw, err := json.Marshal(ChatSource{Value: &branch})
		if err != nil {
			t.Fatalf("marshal sideChat union: %v", err)
		}

		for _, raw := range [][]byte{branchRaw, unionRaw} {
			var decoded map[string]any
			if err := json.Unmarshal(raw, &decoded); err != nil {
				t.Fatalf("decode serialized sideChat payload: %v", err)
			}
			if got := decoded["kind"]; got != "sideChat" {
				t.Fatalf("expected sideChat kind, got %#v in %s", got, raw)
			}
			selection, ok := decoded["selection"].(map[string]any)
			if !ok || selection["text"] != "const value = compute()" || selection["responsePartId"] != "part-7" {
				t.Fatalf("expected sideChat selection to round-trip, got %#v in %s", decoded["selection"], raw)
			}
		}
	})
}

func stringPtr(value string) *string {
	return &value
}
