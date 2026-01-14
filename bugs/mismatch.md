This is a solid “v1” spec. It is big, but it is coherent: clear separation of public auth, admin CRUD, and external SDK with API key auth.

That said, there are a handful of sharp edges that will absolutely bite you in implementation and client UX. Here’s what I’d tighten.

## 1) Path param types are inconsistent (string vs int64)

You mix `Id string \`path:"id"``with backing fields that are clearly`int64` (email lists, campaigns, designs, queue items, etc.).

Examples:

- `EmailListInfo.Id int64` but `GetEmailListRequest { Id string \`path:"id"` }`
- `CampaignInfo.Id int64` but `GetCampaignRequest { Id string \`path:"id"` }`
- `EmailDesignInfo.Id int64` but `GetEmailDesignRequest { Id string \`path:"id"` }`

Pick one:

- If IDs are numeric in DB and URLs, make request path fields `int64`.
- If you want opaque IDs in URLs, make the resource IDs strings everywhere (and stop exposing int64 IDs in response types).

Right now you’re forcing conversions and inviting 400s for “id not an int” bugs.

## 2) You’re mixing goctl-style “optional” tags with standard `omitempty`

You use both:

- `json:"field,optional"` (goctl convention)
- `json:"field,omitempty"` (Go convention)

Example:

- `ValidateInvitationResponse.Message string \`json:"message,omitempty"``
- Most other optional fields use `,optional`.

This becomes a consistency tax across the codebase and clients.

Recommendation:

- For go-zero API definitions, stick to `,optional` consistently.
- Only use `omitempty` if you’re deliberately bypassing goctl conventions (I would not).

## 3) Some delete endpoints return different “empty” shapes

You return:

- `AnalyticsResponse` in many deletes
- `Response` in some deletes
- `EmptyResponse {}` in some places
- And one delete has no `returns(...)` at all

Examples:

- `delete /organizations/:id (DeleteOrgRequest)` has no returns.
- `delete /rules/:id` returns `EmptyResponse`
- `delete /users/:id` returns `AnalyticsResponse`
- `delete /webhooks/:id` returns `Response`

Pick a standard:

- If you want “always return JSON”, use `Response { success, message }` for all non-resource deletes.
- If you want REST-pure, return 204 with no body (but goctl likes explicit returns, so most people standardize on `Response`).

Right now, client code has to special-case.

## 4) Public vs admin route prefixes are a little confusing

You have:

- Public health at `/api/v1/health`
- Public auth at `/api/v1/auth/...`
- Admin at `/api/admin/...`
- SDK at `/sdk/v1/...`

This is fine, but naming “admin/auth” group vs “auth” group can confuse middleware routing if you’re not careful.

My suggestion:

- Keep prefixes, but standardize group names by purpose not path:

  - `group: public`
  - `group: auth_public`
  - `group: admin`
  - `group: sdk`
  - `group: internal_tracking` (pixels/redirects)

That makes logs, metrics, rate limits, and middleware mapping way cleaner.

## 5) SDK endpoints need idempotency or you will create duplicates

You have endpoints that will commonly be retried by clients:

- `POST /sdk/v1/contacts`
- `POST /sdk/v1/emails`
- `POST /sdk/v1/sequences/enroll`
- webhook delivery retries on your side will trigger duplicate “send” calls if the client retries poorly

Add an idempotency mechanism now:

- Support `Idempotency-Key` header for POSTs that create or enqueue.
- Store key + request hash for a short TTL (24h is plenty) and return the original response on repeat.

Without this, you’ll get duplicate contacts, duplicate enrollments, and duplicate emails when someone’s network hiccups.

## 6) Webhook API is missing two practical things

You’ve got great coverage (CRUD + logs + test), but real webhook systems need:

1. **Replay a delivery**

- `POST /webhooks/:id/logs/:logId/replay` (or `/deliveries/:deliveryId/replay`)
- This is gold for support and for customer debugging.

2. **Disable on repeated failures**

- You track success/fail counts, but the spec doesn’t show auto-disable rules or the fields needed to explain status (ex: disabled_reason, disabled_at).
- Even if implemented internally, expose it in `WebhookInfo`.

## 7) Campaign scheduling needs timezone clarity

`ScheduleCampaignRequest.ScheduledAt string // ISO8601 timestamp`

ISO8601 can be timezone-aware or naive. Make it explicit:

- Require RFC3339 with timezone offset (or Z).
- Or allow “local time + timezone” fields.

Right now, someone will schedule “9am” and it will go out at 9am server time, then blame you.

## 8) Sequence timing model has a hidden sharp edge

You have both:

- Per-template `DelayHours`
- Sequence-level `SendHour` + `SendTimezone` (optional)

This is a good model, but define precedence in behavior (or clients will be confused):

- If `SendHour` is set, does `DelayHours` round up to the next `SendHour`?
- If `SendHour` is null, do you send exactly after delay?
- What happens if user joins at 23:55 and SendHour=9?

I’d codify it in the API docs, even briefly.

## 9) Contact and subscriber concepts overlap

You have:

- Subscriber admin actions by `/subscribers/:id/...`
- SDK contacts by `/contacts` with tags/lists/etc
- Email list subscribers endpoints return `contact_id`

This is fine if “subscriber is a contact in a list”, but the API surface implies two different entities.

If they’re the same underlying entity, I’d unify naming:

- Admin: prefer `/contacts` for lifecycle actions, and `/lists/:id/subscribers` for list membership actions.
- Or keep `/subscribers` but then rename SDK “contacts” to “subscribers” for conceptual match.

As-is, new users will ask “what is the difference between a contact and a subscriber?”

## 10) A couple of small correctness nits

- `DeleteSequence` uses `(GetSequenceRequest)` not a `DeleteSequenceRequest`. Works, but reads wrong.
- `GetTransactionalEmailStats` uses `(GetTransactionalEmailRequest)` which has `Id` but the name implies “get template”, not “get stats”. Minor, but these papercuts add up.
- `VerifyEmailRequest` uses `form:"token"` with `GET /auth/verify-email` which is correct for query param, but you also have separate confirm endpoints registered manually. Make sure you don’t end up with two ways to confirm email (or be explicit: one is account email verification, one is list double opt-in).

---

If you want the highest-leverage cleanup with minimal churn, do these 5 things:

1. Normalize ID types (string vs int64) across requests and responses.
2. Standardize delete responses (pick `Response` or `EmptyResponse` everywhere).
3. Add Idempotency-Key support to SDK POST endpoints.
4. Clarify schedule timestamps must be RFC3339 with timezone.
5. Write precedence rules for `DelayHours` vs `SendHour`.

If you tell me which direction you want IDs to go (opaque string IDs everywhere vs numeric IDs everywhere), I’ll give you a concrete diff-style set of type/request changes to make the whole spec consistent.
