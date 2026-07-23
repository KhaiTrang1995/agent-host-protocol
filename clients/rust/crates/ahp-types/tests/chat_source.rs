use ahp_types::commands::ChatSource;

#[test]
fn chat_source_routes_by_kind() {
    let fork = serde_json::from_str::<ChatSource>(
        r#"{"kind":"fork","chat":"ahp-chat:/main","turnId":"turn-12"}"#,
    )
    .expect("decode fork source");
    let side_chat = serde_json::from_str::<ChatSource>(
        r#"{"kind":"sideChat","chat":"ahp-chat:/main","turnId":"turn-active"}"#,
    )
    .expect("decode side chat source");

    match fork {
        ChatSource::Fork(value) => {
            assert_eq!(value.chat, "ahp-chat:/main");
            assert_eq!(value.turn_id, "turn-12");
        }
        other => panic!("expected fork variant, got {other:?}"),
    }

    match side_chat {
        ChatSource::SideChat(value) => {
            assert_eq!(value.chat, "ahp-chat:/main");
            assert_eq!(value.turn_id, "turn-active");
        }
        other => panic!("expected sideChat variant, got {other:?}"),
    }
}

#[test]
fn chat_source_rejects_missing_or_unknown_kind() {
    for raw in [
        r#"{"chat":"ahp-chat:/main","turnId":"turn-12"}"#,
        r#"{"kind":"future","chat":"ahp-chat:/main","turnId":"turn-12"}"#,
    ] {
        assert!(
            serde_json::from_str::<ChatSource>(raw).is_err(),
            "expected decode failure for {raw}"
        );
    }
}
