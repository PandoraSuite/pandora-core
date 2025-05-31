package enums

type SensitiveAction string

const (
	SensitiveActionNull         SensitiveAction = ""
	SensitiveActionRevealAPIKey SensitiveAction = "REVEAL_API_KEY"
)

func ParseSensitiveAction(action string) (SensitiveAction, bool) {
	switch s := SensitiveAction(action); s {
	case SensitiveActionRevealAPIKey:
		return s, true
	default:
		return SensitiveActionNull, false
	}
}

type Scope string

const (
	ScopeNull         Scope = ""
	ScopeRevealAPIKey Scope = "api_key:reveal"
)

func ParseScope(action string) (Scope, bool) {
	switch s := Scope(action); s {
	case ScopeRevealAPIKey:
		return s, true
	default:
		return ScopeNull, false
	}
}
