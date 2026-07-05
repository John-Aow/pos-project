# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`

**Created**: [DATE]

**Status**: Draft

**Input**: User description: "$ARGUMENTS"

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - [Brief Title] (Priority: P1)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently - e.g., "Can be fully tested by [specific action] and delivers [specific value]"]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]
2. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 2 - [Brief Title] (Priority: P2)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 3 - [Brief Title] (Priority: P3)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases. For POS features, explicitly consider
  checkout interruption, duplicate action, permission denial, catalog mismatch,
  inventory conflict, payment handoff failure, and reconciliation after offline use
  when applicable. For database-impacting work, explicitly identify whether schema
  changes, existing schema modification, or data deletion could be implicated. For
  frontend work, explicitly consider mobile, tablet, desktop, portrait, landscape,
  touch input, hidden browser chrome, and native WebView constraints.
-->

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST [specific capability, e.g., "allow users to create accounts"]
- **FR-002**: System MUST [specific capability, e.g., "validate email addresses"]
- **FR-003**: Users MUST be able to [key interaction, e.g., "reset their password"]
- **FR-004**: System MUST [data requirement, e.g., "persist user preferences"]
- **FR-005**: System MUST [behavior, e.g., "log all security events"]
- **FR-006**: System MUST define audit events for sensitive staff, customer, payment,
  transaction, catalog, inventory, and permission changes that this feature touches.
- **FR-007**: Backend API endpoints MUST validate input before calling usecase logic.
- **FR-008**: Implementation MUST NOT begin until acceptance criteria are defined.
- **FR-009**: Frontend interfaces MUST support mobile, tablet, desktop, and native
  WebView usage when this feature includes user-facing UI.

*Example of marking unclear requirements:*

- **FR-010**: System MUST authenticate users via [NEEDS CLARIFICATION: auth method not specified - PIN, email/password, SSO, OAuth?]
- **FR-011**: System MUST retain user data for [NEEDS CLARIFICATION: retention period not specified]

### POS Risk Requirements *(include when applicable)*

- **Transaction Integrity**: [Totals, rounding, tax, discount, tender, failure, retry, and idempotency requirements]
- **Catalog/Inventory Accuracy**: [Product, price, stock, modifier, tax, availability, snapshot, and traceability requirements]
- **Offline Operation**: [Network loss, local queue, sync, duplicate prevention, and reconciliation requirements]
- **Security/Audit**: [Permission, privacy, secret handling, and audit log requirements]
- **Database Safety**: [Confirm no existing schema/data deletion; identify new migrations if needed]
- **Monorepo Boundary**: [Affected side: frontend, backend, packages; confirm no direct cross-side implementation dependency]
- **Responsive/WebView**: [Mobile ≥360px, tablet ≥768px, desktop ≥1024px, portrait/landscape, touch targets, in-app navigation, viewport meta, WebView constraints]

### Architecture Requirements *(include when applicable)*

- **Frontend**: [Vue 3 Composition API, TypeScript, Pinia, `<script setup>`, component/store/service requirements]
- **Backend**: [Clean Architecture layer responsibilities and dependency direction]
- **API Contract**: [Request/response validation, errors, and shared contract location]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: [Measurable metric, e.g., "Users can complete account creation in under 2 minutes"]
- **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
- **SC-003**: [User satisfaction metric, e.g., "90% of users successfully complete primary task on first attempt"]
- **SC-004**: [Business metric, e.g., "Reduce support tickets related to [X] by 50%"]
- **SC-005**: [POS reliability metric, e.g., "Cashiers can complete the primary checkout flow within [N] seconds at p95"]
- **SC-006**: [Quality metric, e.g., "Frontend and backend coverage remain at or above 80%"]
- **SC-007**: [Responsive metric, e.g., "All interactive UI works at one mobile, one tablet, and one desktop breakpoint"]

## Assumptions

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right assumptions based on reasonable defaults
  chosen when the feature description did not specify certain details.
-->

- [Assumption about target users, e.g., "Users have stable internet connectivity"]
- [Assumption about scope boundaries, e.g., "Mobile support is out of scope for v1"]
- [Assumption about data/environment, e.g., "Existing authentication system will be reused"]
- [Dependency on existing system/service, e.g., "Requires access to the existing user profile API"]
