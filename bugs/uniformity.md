You’re right. **`emaillist` is noise**, and it leaks implementation instead of intent.

This is one of those naming decisions that quietly determines whether the API feels clean or clunky long-term. You should fix it now.

## Clear recommendation

Use **`list` everywhere** at the domain level.

- Users think in “lists”, not “email lists”
- You already imply email by context (this is an email platform)
- Shorter names matter a lot in URLs, SDKs, and mental models

## What I would standardize to

### Routes

Change:

```
/email-lists
/email-lists/:id
/email-lists/:id/subscribers
```

To:

```
/lists
/lists/:id
/lists/:id/subscribers
```

Admin grouping becomes:

```
group: admin/lists
```

SDK already does this correctly:

```
/sdk/v1/lists/:slug/subscribe
```

So this actually **reduces inconsistency** instead of creating it.

## Types and structs

Rename these types:

- `EmailListInfo` → `ListInfo`
- `ListEmailListsResponse` → `ListListsResponse` (or just `ListsResponse`)
- `GetEmailListRequest` → `GetListRequest`
- `CreateEmailListRequest` → `CreateListRequest`
- `UpdateEmailListRequest` → `UpdateListRequest`
- `EmailListStatsResponse` → `ListStatsResponse`

Handlers:

- `ListEmailLists` → `ListLists`
- `GetEmailList` → `GetList`
- `CreateEmailList` → `CreateList`
- `UpdateEmailList` → `UpdateList`
- `DeleteEmailList` → `DeleteList`

Yes, `ListLists` looks a little awkward in Go naming, but:

- It’s idiomatic
- It’s predictable
- It’s far better than baking “email” into everything

If it bothers you, `ListsIndex` or `GetLists` are acceptable alternatives.

## Why this matters more than it seems

1. **Future-proofing**
   If you ever add:

- SMS lists
- Push notification lists
- Segmented lists
- Dynamic lists

You don’t want `EmailEmailListEmailStats` syndrome.

2. **SDK ergonomics**
   Users will write things like:

```ts
client.lists.subscribe(...)
client.lists.get(id)
```

That reads clean.
`client.emailLists.subscribe(...)` already feels dated.

3. **Mental compression**
   Every extra word in an API name increases cognitive load.
   You’ve already done the hard work architecturally. Don’t lose polish on naming.

## One last alignment tweak I’d make

You already use:

- `ListSlug`
- `ListIds`
- `lists/:slug/subscribe`

So renaming the admin side to `/lists` actually **completes the model**. Right now it’s half-clean, half-legacy.

---

**Verdict:**
Yes. Rename `emaillist` to `list` everywhere.
This is a correct instinct, and it’s worth doing before anyone depends on the old names.

If you want, I can:

- Produce a clean “before → after” diff of the list-related sections
- Or help you pick final Go handler names that don’t feel awkward while staying idiomatic
