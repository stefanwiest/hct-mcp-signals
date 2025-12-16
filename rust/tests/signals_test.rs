//! Integration tests for hct-mcp-signals crate.
//!
//! These tests verify the public API and end-to-end behavior.

use hct_mcp_signals::{
    attacca, caesura, cue, downbeat, fermata, tacet, vamp,
    embed_signal, extract_signal, get_json_schema, get_json_schema_string,
    get_mcp_extension_schema, HCTSignal, HoldType, McpTaskSend, Performance,
    SignalType, Tempo,
};
use serde_json::json;

// ============================================================
// Signal Creation Tests
// ============================================================

#[test]
fn test_all_signal_types_exist() {
    // Verify all 7 signal types
    let cue_sig = cue("src", ["target"]).build();
    let fermata_sig = fermata("src", "reason").build();
    let attacca_sig = attacca("src", ["target"]).build();
    let vamp_sig = vamp("src", "condition").build();
    let caesura_sig = caesura("src", "reason").build();
    let tacet_sig = tacet("src").build();
    let downbeat_sig = downbeat("src", "sync").build();

    assert_eq!(cue_sig.signal_type, SignalType::Cue);
    assert_eq!(fermata_sig.signal_type, SignalType::Fermata);
    assert_eq!(attacca_sig.signal_type, SignalType::Attacca);
    assert_eq!(vamp_sig.signal_type, SignalType::Vamp);
    assert_eq!(caesura_sig.signal_type, SignalType::Caesura);
    assert_eq!(tacet_sig.signal_type, SignalType::Tacet);
    assert_eq!(downbeat_sig.signal_type, SignalType::Downbeat);
}

#[test]
fn test_cue_with_all_options() {
    let signal = cue("orchestrator", ["analyst", "synthesizer"])
        .with_urgency(9)
        .with_tempo(Tempo::Presto)
        .with_timeout_ms(5000)
        .with_payload_entry("task", json!("Analyze Q4"))
        .with_payload_entry("priority", json!("high"))
        .build();

    assert_eq!(signal.source, "orchestrator");
    assert_eq!(signal.targets, vec!["analyst", "synthesizer"]);
    assert_eq!(signal.performance.as_ref().unwrap().urgency, 9);
    assert_eq!(signal.performance.as_ref().unwrap().tempo, Tempo::Presto);
    assert_eq!(signal.performance.as_ref().unwrap().timeout_ms, Some(5000));
    assert_eq!(signal.payload.get("task"), Some(&json!("Analyze Q4")));
    assert_eq!(signal.payload.get("priority"), Some(&json!("high")));
}

#[test]
fn test_fermata_defaults() {
    let signal = fermata("reporter", "Ready for review").build();
    
    assert_eq!(signal.targets, vec!["governance"]);
    assert!(signal.conditions.is_some());
    let conditions = signal.conditions.unwrap();
    assert_eq!(conditions.hold_type, Some(HoldType::Human));
}

#[test]
fn test_attacca_is_urgent() {
    let signal = attacca("agent", ["next"]).build();
    
    assert_eq!(signal.performance.as_ref().unwrap().urgency, 10);
    assert_eq!(signal.performance.as_ref().unwrap().tempo, Tempo::Presto);
}

#[test]
fn test_vamp_conditions() {
    let signal = vamp("verifier", "score > 0.9").build();
    
    assert!(signal.conditions.is_some());
    let cond = signal.conditions.unwrap();
    assert_eq!(cond.repeat_until, Some("score > 0.9".to_string()));
    assert_eq!(cond.quality_threshold, Some(0.9));
}

#[test]
fn test_caesura_broadcasts() {
    let signal = caesura("governance", "Budget exceeded").build();
    
    assert_eq!(signal.targets, vec!["*"]);
    assert_eq!(signal.performance.as_ref().unwrap().urgency, 10);
}

#[test]
fn test_tacet_no_targets() {
    let signal = tacet("sleeping_agent").build();
    
    assert!(signal.targets.is_empty());
}

#[test]
fn test_downbeat_sync_point() {
    let signal = downbeat("conductor", "daily_standup").build();
    
    assert_eq!(signal.targets, vec!["*"]);
    assert_eq!(signal.payload.get("sync_point"), Some(&json!("daily_standup")));
}

// ============================================================
// Serialization/Deserialization Tests
// ============================================================

#[test]
fn test_signal_json_roundtrip() {
    let original = cue("orch", ["analyst"])
        .with_urgency(8)
        .with_tempo(Tempo::Allegro)
        .with_payload_entry("task", json!("test"))
        .build();
    
    let json_str = serde_json::to_string(&original).unwrap();
    let parsed: HCTSignal = serde_json::from_str(&json_str).unwrap();
    
    assert_eq!(original.signal_type, parsed.signal_type);
    assert_eq!(original.source, parsed.source);
    assert_eq!(original.targets, parsed.targets);
}

#[test]
fn test_mcp_json_format() {
    let signal = fermata("reporter", "test").build();
    let mcp_json = signal.to_mcp_json().unwrap();
    
    assert!(mcp_json.contains(r#""hct_signal""#));
    assert!(mcp_json.contains(r#""fermata""#));
}

#[test]
fn test_from_mcp_json() {
    let signal = cue("test", ["target"]).build();
    let mcp_json = signal.to_mcp_json().unwrap();
    
    let parsed = HCTSignal::from_mcp_json(&mcp_json).unwrap();
    assert_eq!(parsed.signal_type, SignalType::Cue);
    assert_eq!(parsed.source, "test");
}

// ============================================================
// MCP Integration Tests
// ============================================================

#[test]
fn test_mcp_task_send() {
    let signal = cue("orch", ["analyst"]).with_urgency(8).build();
    let msg = McpTaskSend::new("task-123", "Analyze Q4", signal);
    
    assert_eq!(msg.jsonrpc, "2.0");
    assert_eq!(msg.method, "tasks/send");
    assert_eq!(msg.params.id, "task-123");
    assert_eq!(msg.params.message.content, "Analyze Q4");
    assert!(msg.params.hct_signal.is_some());
}

#[test]
fn test_mcp_task_send_json() {
    let signal = cue("orch", ["analyst"]).build();
    let msg = McpTaskSend::new("task-123", "test", signal);
    let json = msg.to_json().unwrap();
    
    assert!(json.contains(r#""jsonrpc": "2.0""#));
    assert!(json.contains(r#""method": "tasks/send""#));
}

#[test]
fn test_embed_signal() {
    let signal = cue("orch", ["analyst"]).build();
    let params = r#"{"id": "task-123", "other": "data"}"#;
    
    let result = embed_signal(params, &signal).unwrap();
    assert!(result.contains("hct_signal"));
    assert!(result.contains("task-123"));
}

#[test]
fn test_extract_signal() {
    let signal = fermata("test", "reason").build();
    let json = signal.to_mcp_json().unwrap();
    
    let extracted = extract_signal(&json).unwrap();
    assert!(extracted.is_some());
    assert_eq!(extracted.unwrap().signal_type, SignalType::Fermata);
}

#[test]
fn test_extract_signal_none() {
    let json = r#"{"id": "task-123"}"#;
    let extracted = extract_signal(json).unwrap();
    assert!(extracted.is_none());
}

// ============================================================
// Schema Tests
// ============================================================

#[test]
fn test_json_schema_structure() {
    let schema = get_json_schema();
    
    assert!(schema["$schema"].is_string());
    assert!(schema["title"].is_string());
    assert!(schema["required"].is_array());
    assert!(schema["properties"].is_object());
}

#[test]
fn test_json_schema_required_fields() {
    let schema = get_json_schema();
    let required = schema["required"].as_array().unwrap();
    
    let required_strs: Vec<&str> = required.iter()
        .filter_map(|v| v.as_str())
        .collect();
    
    assert!(required_strs.contains(&"type"));
    assert!(required_strs.contains(&"source"));
    assert!(required_strs.contains(&"targets"));
    assert!(required_strs.contains(&"payload"));
    assert!(required_strs.contains(&"timestamp"));
}

#[test]
fn test_json_schema_signal_types() {
    let schema = get_json_schema();
    let type_enum = &schema["properties"]["type"]["enum"];
    
    assert!(type_enum.as_array().unwrap().contains(&json!("cue")));
    assert!(type_enum.as_array().unwrap().contains(&json!("fermata")));
    assert!(type_enum.as_array().unwrap().contains(&json!("downbeat")));
}

#[test]
fn test_json_schema_string() {
    let schema_str = get_json_schema_string().unwrap();
    
    assert!(schema_str.contains("HCTSignal"));
    assert!(schema_str.contains("json-schema.org"));
}

#[test]
fn test_mcp_extension_schema() {
    let schema = get_mcp_extension_schema();
    
    assert!(schema["properties"]["hct_signal"].is_object());
    assert_eq!(schema["title"], "HCT-MCP Extension");
}

// ============================================================
// Performance Builder Tests
// ============================================================

#[test]
fn test_performance_builder() {
    let perf = Performance::new()
        .with_urgency(7)
        .with_tempo(Tempo::Andante)
        .with_timeout_ms(10000);
    
    assert_eq!(perf.urgency, 7);
    assert_eq!(perf.tempo, Tempo::Andante);
    assert_eq!(perf.timeout_ms, Some(10000));
}

#[test]
fn test_urgency_clamping() {
    let too_high = Performance::new().with_urgency(15);
    let too_low = Performance::new().with_urgency(0);
    
    assert_eq!(too_high.urgency, 10);
    assert_eq!(too_low.urgency, 1);
}

#[test]
fn test_all_tempos() {
    let tempos = vec![
        Tempo::Largo,
        Tempo::Andante,
        Tempo::Moderato,
        Tempo::Allegro,
        Tempo::Presto,
    ];
    
    for tempo in tempos {
        let signal = cue("test", ["target"]).with_tempo(tempo.clone()).build();
        assert_eq!(signal.performance.as_ref().unwrap().tempo, tempo);
    }
}

#[test]
fn test_all_hold_types() {
    // Verify all hold types exist
    let types = vec![
        HoldType::Human,
        HoldType::Governance,
        HoldType::Resource,
        HoldType::Quality,
    ];
    
    for ht in types {
        let json = serde_json::to_string(&ht).unwrap();
        assert!(!json.is_empty());
    }
}

// ============================================================
// Signal Builder Edge Cases
// ============================================================

#[test]
fn test_empty_targets() {
    let signal = cue("orch", Vec::<String>::new()).build();
    assert!(signal.targets.is_empty());
}

#[test]
fn test_string_targets() {
    let signal = cue("orch", vec!["a".to_string(), "b".to_string()]).build();
    assert_eq!(signal.targets.len(), 2);
}

#[test]
fn test_timestamp_is_set() {
    let signal = cue("orch", ["target"]).build();
    // Timestamp should be recent (within last minute)
    let now = chrono::Utc::now();
    let diff = now - signal.timestamp;
    assert!(diff.num_seconds() < 60);
}
