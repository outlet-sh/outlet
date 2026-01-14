You’re right, and this is the **correct instinct**. If MCP is going to be the control plane, **partial CRUD is worse than no CRUD** because agents assume completeness.

Right now you have **creation and mutation**, but not a closed loop. Below is what’s missing, what to add, and how to structure it so agents can _reliably_ manage email and sequences end to end.

I’ll be direct.

---

## The rule

> **Every MCP-managed domain object must support full lifecycle control by tools alone.**
> No UI-only operations. No hidden capabilities.

That means **Create, Read, Update, Delete, plus List** where it makes sense.

---

## Current state vs required state

### Lists

You currently have:

- `email_list_create`
- `email_list_list`

You are missing:

- ❌ `list_get`
- ❌ `list_update`
- ❌ `list_delete`

### Sequences

You currently have:

- `sequence_create`
- `sequence_list`
- `sequence_get`

You are missing:

- ❌ `sequence_update`
- ❌ `sequence_delete`
- ❌ `sequence_toggle` (soft enable/disable)

### Sequence Emails (Templates)

You currently have:

- `sequence_email_add`
- `sequence_email_update`

You are missing:

- ❌ `sequence_email_get`
- ❌ `sequence_email_delete`
- ❌ `sequence_email_list`
- ❌ `sequence_email_move` (reorder safely)

This is the **biggest hole**.

---

## What a complete MCP surface should look like

### Lists

```
list_create
list_list
list_get
list_update
list_delete
```

Optional but high value:

```
list_stats
list_subscribers_list
```

---

### Sequences

```
sequence_create
sequence_list
sequence_get
sequence_update
sequence_delete
sequence_toggle
```

You already store `IsActive`. Expose it cleanly.

---

### Sequence Emails (Templates)

Treat these as first-class resources.

```
sequence_email_create      (rename from add)
sequence_email_list
sequence_email_get
sequence_email_update
sequence_email_delete
sequence_email_move
```

Why `move` matters:

- Reordering is not just updating `position`
- Agents will break ordering without guardrails
- You want atomic repositioning logic

---

## Why this matters specifically for MCP (not REST)

Agents do not infer UI affordances. They assume tools are _truth_.

Right now an agent can:

- Create a sequence
- Add emails
- Update an email

But it **cannot**:

- Inspect a single email fully
- Delete a broken email
- Safely reorder emails
- Disable a sequence without guessing
- Remove a list cleanly

This forces hallucinated workflows like:

> “I will just update position to 999 and hope it goes last.”

That is how systems rot.

---

## Strong recommendation on semantics

### 1. Make delete explicit and irreversible

Do not overload update with “active=false”.

Good:

```
sequence_email_delete
sequence_delete
list_delete
```

Agents respect verbs.

---

### 2. Separate “toggle” from “update”

You did this correctly for rules.

Repeat the pattern:

```
sequence_toggle
sequence_email_toggle (optional)
```

This makes intent clear and auditable.

---

### 3. Never expose raw `position` mutation alone

Bad:

```
position: 2
```

Good:

```
sequence_email_move
  email_id
  after_email_id | before_email_id | position
```

This lets you normalize ordering safely.

---

## Naming cleanup you should do now

Since you already agreed on `list`:

Rename:

- `email_list_*` → `list_*`
- `sequence_email_add` → `sequence_email_create`

“Add” is UI language. MCP is API language.

---

## Minimal additions that unlock everything

If you want the **smallest possible set** that makes the system complete:

Add these 8 tools:

```
list_get
list_update
list_delete

sequence_update
sequence_delete
sequence_toggle

sequence_email_list
sequence_email_delete
```

Everything else is optimization.

---

## The bigger picture

You’re building something closer to **Sendy + ActiveCampaign + rules + MCP-native automation**.

The moment agents can:

- Create
- Inspect
- Modify
- Remove

…without guessing, your system becomes **composable infrastructure**, not just “email with AI”.

You’re very close. This is the last structural gap.

If you want, next I can:

- Propose exact input/output structs for the missing tools
- Or help you design the `sequence_email_move` logic so it’s agent-proof
