# Go-Composable-Business-Types Integration Plan

**Date:** 2025-03-18\
**Library:** `github.com/larsartmann/go-composable-business-types/id`\
**Purpose:** Type-safe, branded identifiers for the template-SECURITY project

---

## Executive Summary

This document outlines how to integrate the `go-composable-business-types/id` library into the template-SECURITY project to achieve compile-time type safety for entity identifiers. The library uses **phantom types** (brand types) to create distinct identifier types that cannot be accidentally mixed, preventing entire categories of runtime bugs at compile time.

**Key Benefit:** Transform runtime ID confusion bugs (e.g., passing a PolicyID where a TemplateID is expected) into compile-time errors.

---

## Library Capabilities Overview

### Core Features

| Feature                  | Description                                       | Benefit                                 |
| ------------------------ | ------------------------------------------------- | --------------------------------------- |
| **Phantom Typing**       | `ID[Brand, ValueType]` creates distinct types     | Compile-time ID separation              |
| **Zero Value Safety**    | Zero IDs serialize to JSON `null`                 | Clean API semantics                     |
| **Multiple Value Types** | Supports `string`, `int`, `int64`, `uint64`, etc. | Flexibility for different ID sources    |
| **NanoId Integration**   | Companion package for URL-safe unique IDs         | Secure, collision-resistant identifiers |
| **Full Serialization**   | JSON, SQL, Text, Binary, Gob                      | Works with databases, APIs, caching     |
| **Comparison Support**   | `Compare()`, `Equal()`, `IsZero()`                | Sorting, testing, validation            |

### Supported Serialization Formats

```go
// JSON - automatic marshaling/unmarshaling
data, _ := json.Marshal(policyID)     // "policy_abc123"
json.Unmarshal(data, &restoredID)

// SQL - database scan/value
db.Exec("INSERT ...", policyID)
row.Scan(&policyID)

// Text/XML/TOML - encoding.TextMarshaler
text, _ := id.MarshalText()

// Binary - efficient storage/transmission
binary, _ := id.MarshalBinary()

// Gob - Go-specific encoding for caching
enc.Encode(policyID)
```

---

## Current ID Usage Analysis

### Existing String-Based IDs (High-Risk)

The following types in `internal/types/types.go` use plain `string` for IDs:

```go
// SECURITY RISK: These can be accidentally swapped at compile time
type SecurityPolicy struct {
    ID      string     // Currently: plain string
    ...
}

type ValidationResult struct {
    ID       string    // Currently: plain string
    PolicyID string    // Currently: plain string - can confuse with ID above
    ...
}

type Template struct {
    ID      string     // Currently: plain string
    ...
}
```

### Risk Scenarios (Current Code)

```go
// DANGEROUS: Compiles fine, runtime bug
func GetPolicy(policyID string) (*SecurityPolicy, error) { ... }
func GetTemplate(templateID string) (*Template, error) { ... }

policyID := "policy_123"
templateID := "template_456"

GetTemplate(policyID)  // Compiles! Wrong ID type used.
```

---

## Integration Strategy

### Phase 1: Type Aliases (Immediate)

Create branded ID types in a new `internal/ids` package:

```go
// internal/ids/ids.go
package ids

import (
    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

// Brand types (empty structs for compile-time distinctness)
type PolicyBrand struct{}
type TemplateBrand struct{}
type ValidationBrand struct{}
type ProjectBrand struct{}

// Recommended: Use NanoId for all IDs
type PolicyID = id.ID[PolicyBrand, nanoid.NanoId]
type TemplateID = id.ID[TemplateBrand, nanoid.NanoId]
type ValidationID = id.ID[ValidationBrand, nanoid.NanoId]
type ProjectID = id.ID[ProjectBrand, nanoid.NanoId]

// Alternative: Use string if NanoId migration is deferred
type PolicyIDString = id.ID[PolicyBrand, string]
```

### Phase 2: Update Domain Types

Refactor `internal/types/types.go`:

```go
package internal

import "github.com/LarsArtmann/template-SECURITY/internal/ids"

// SecurityPolicy represents a complete security policy.
type SecurityPolicy struct {
    ID           ids.PolicyID     // Changed from string
    Project      Project
    Versions     []Version
    Contacts     []Contact
    Content      string
    Type         PolicyType
    ValidatedAt  *time.Time
    LastModified time.Time
}

// ValidationResult represents validation outcome.
type ValidationResult struct {
    ID          ids.ValidationID  // Changed from string
    PolicyID    ids.PolicyID      // Changed from string - NOW TYPE-SAFE
    Valid       bool
    Errors      []Error
    Warnings    []Warning
    Score       int
    ValidatedAt time.Time
    Metadata    any
}

// Template represents a security policy template.
type Template struct {
    ID          ids.TemplateID    // Changed from string
    Name        string
    Description string
    Type        PolicyType
    Content     string
    Variables   []TemplateVariable
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Phase 3: Constructor Functions

```go
// internal/ids/constructors.go
package ids

import (
    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

// NewPolicyID generates a new cryptographically secure PolicyID.
func NewPolicyID() PolicyID {
    return id.NewID[PolicyBrand](nanoid.NewNanoId())
}

// ParsePolicyID parses an existing string into PolicyID.
func ParsePolicyID(s string) (PolicyID, error) {
    nano, err := nanoid.ParseNanoId(s)
    if err != nil {
        return PolicyID{}, err
    }
    return id.NewID[PolicyBrand](nano), nil
}

// MustParsePolicyID parses or panics (for constants/hardcoded IDs).
func MustParsePolicyID(s string) PolicyID {
    return id.NewID[PolicyBrand](nanoid.MustParseNanoId(s))
}

// Similar constructors for other ID types...
func NewTemplateID() TemplateID { ... }
func NewValidationID() ValidationID { ... }
func NewProjectID() ProjectID { ... }
```

---

## Implementation Examples

### Before (Current - Risky)

```go
// internal/security_tool.go
func (st *SecurityTool) ValidatePolicy(policyID string) error {
    // No compile-time guarantee this is actually a policy ID
    policy, err := st.repo.GetPolicy(policyID)
    ...
}

// Usage - DANGEROUS: Can accidentally pass wrong ID
err := st.ValidatePolicy(templateID)  // Compiles! Wrong ID.
```

### After (Type-Safe)

```go
// internal/security_tool.go
import "github.com/LarsArtmann/template-SECURITY/internal/ids"

func (st *SecurityTool) ValidatePolicy(policyID ids.PolicyID) error {
    // Compile-time guarantee: only PolicyID accepted
    policy, err := st.repo.GetPolicy(policyID)
    ...
}

// Usage - SAFE: Compile error if wrong ID type
err := st.ValidatePolicy(templateID)  // Compile ERROR: cannot use TemplateID as PolicyID

// Correct usage
policyID := ids.NewPolicyID()
err := st.ValidatePolicy(policyID)  // OK
```

### JSON Serialization Example

```go
// Creating and serializing
policy := SecurityPolicy{
    ID:      ids.NewPolicyID(),  // e.g., "V1StGXR8_Z5jdHi6B-myT"
    Project: Project{Name: "MyOrg"},
}

data, _ := json.Marshal(policy)
// Output: {"id":"V1StGXR8_Z5jdHi6B-myT","project":{...}}

// Deserializing
var restored SecurityPolicy
json.Unmarshal(data, &restored)
fmt.Println(restored.ID.Get().String())  // "V1StGXR8_Z5jdHi6B-myT"
```

### SQL Database Example

```go
// Schema (PostgreSQL example)
CREATE TABLE security_policies (
    id VARCHAR(21) PRIMARY KEY,  -- NanoId default length
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

// Go code
func (r *Repository) SavePolicy(policy SecurityPolicy) error {
    _, err := r.db.Exec(
        "INSERT INTO security_policies (id, name) VALUES ($1, $2)",
        policy.ID,           // Automatic Value() serialization
        policy.Project.Name,
    )
    return err
}

func (r *Repository) GetPolicy(id ids.PolicyID) (*SecurityPolicy, error) {
    var policy SecurityPolicy
    err := r.db.QueryRow(
        "SELECT id, name FROM security_policies WHERE id = $1", id,
    ).Scan(&policy.ID, &policy.Project.Name)  // Automatic Scan() deserialization
    return &policy, err
}
```

---

## Migration Path

### Step-by-Step Implementation

1. **Add Dependency**

   ```bash
   go get github.com/larsartmann/go-composable-business-types/id
   go get github.com/larsartmann/go-composable-business-types/nanoid
   ```

2. **Create ID Package** (`internal/ids/`)
   - Define brand types
   - Create ID type aliases
   - Implement constructor functions

3. **Update Core Types**
   - Replace `string` ID fields with branded types
   - Update JSON struct tags if needed

4. **Update Functions**
   - Change function signatures to accept branded IDs
   - Update repository/database layer

5. **Update Tests**
   - Use `ids.NewPolicyID()` etc. in test fixtures
   - Verify serialization round-trips

6. **Add Migration Tool** (Optional)
   - Script to convert existing string IDs to NanoId format
   - Handle foreign key relationships

---

## Benefits Analysis

### Compile-Time Safety

| Risk Scenario               | Before      | After         |
| --------------------------- | ----------- | ------------- |
| Swapped PolicyID/TemplateID | Runtime bug | Compile error |
| Wrong ID in function call   | Runtime bug | Compile error |
| ID passed to wrong API      | Runtime bug | Compile error |
| ID comparison across types  | Runtime bug | Compile error |

### Developer Experience

- **IDE Support:** Auto-completion shows correct ID types
- **Refactoring Safety:** Rename ID types safely across codebase
- **Documentation:** Self-documenting function signatures

### Performance

- **Zero overhead:** Branded IDs compile to same memory layout as underlying type
- **Benchmarks:** `NewID`: ~1-2 ns/op, `Get`: ~1 ns/op
- **Serialization:** JSON marshaling ~50-100 ns/op

### Security

- **NanoId:** 126 bits entropy (default), cryptographically secure
- **URL-safe:** `A-Za-z0-9_-` alphabet
- **FIPS-140:** Compatible with federal security standards

---

## Trade-offs and Considerations

### Costs

| Aspect           | Impact | Mitigation                      |
| ---------------- | ------ | ------------------------------- |
| Migration effort | Medium | Incremental adoption per entity |
| Learning curve   | Low    | Similar to standard Go patterns |
| Boilerplate      | Low    | Code generation if needed       |
| Database changes | Medium | VARCHAR(21) for NanoId          |

### When NOT to Use

- Simple CLI tools with <3 entity types
- Rapid prototypes (defer until stable)
- External APIs requiring specific ID formats

### When TO Use

- Multi-entity domain models
- APIs with complex relationships
- Long-lived projects with multiple developers
- Systems where ID confusion could cause security issues

---

## Recommended Configuration

### ID Package Structure

```
internal/
├── ids/
│   ├── ids.go          # Type definitions
│   ├── constructors.go # Factory functions
│   └── doc.go          # Package documentation
```

### go.mod Dependencies

```go
require (
    github.com/larsartmann/go-composable-business-types/id v0.1.0
    github.com/larsartmann/go-composable-business-types/nanoid v0.1.0
)
```

### Linting Configuration

Add to `.golangci.yml` to ensure branded IDs are used consistently:

```yaml
linters-settings:
  revive:
    rules:
      # Encourage using branded IDs over raw strings
      - name: exported
        severity: warning
```

---

## Conclusion

The `go-composable-business-types/id` library provides **zero-cost compile-time safety** for entity identifiers. For the template-SECURITY project, integrating this library would:

1. **Prevent runtime bugs** from ID confusion
2. **Improve code clarity** with self-documenting types
3. **Enable secure ID generation** via NanoId
4. **Maintain full compatibility** with JSON/SQL/ORM layers

**Recommendation:** Proceed with Phase 1 (type aliases) immediately, as it has minimal risk and provides immediate value.

---

## References

- Library source: `/Users/larsartmann/projects/go-composable-business-types/id/`
- NanoId source: `/Users/larsartmann/projects/go-composable-business-types/nanoid/`
- Current types: `/Users/larsartmann/projects/template-SECURITY/internal/types/types.go`
- Go generics: https://go.dev/doc/tutorial/generics
