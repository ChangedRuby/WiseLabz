# Graph Report - WiseLabz  (2026-09-17)

## Corpus Check
- 517 files · ~630,753 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2951 nodes · 6492 edges · 95 communities detected
- Extraction: 54% EXTRACTED · 44% INFERRED · 0% AMBIGUOUS · INFERRED: 2885 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 58|Community 58]]
- [[_COMMUNITY_Community 59|Community 59]]
- [[_COMMUNITY_Community 60|Community 60]]
- [[_COMMUNITY_Community 63|Community 63]]
- [[_COMMUNITY_Community 64|Community 64]]
- [[_COMMUNITY_Community 66|Community 66]]
- [[_COMMUNITY_Community 67|Community 67]]
- [[_COMMUNITY_Community 76|Community 76]]
- [[_COMMUNITY_Community 80|Community 80]]
- [[_COMMUNITY_Community 81|Community 81]]
- [[_COMMUNITY_Community 82|Community 82]]
- [[_COMMUNITY_Community 83|Community 83]]
- [[_COMMUNITY_Community 84|Community 84]]
- [[_COMMUNITY_Community 85|Community 85]]
- [[_COMMUNITY_Community 86|Community 86]]
- [[_COMMUNITY_Community 119|Community 119]]
- [[_COMMUNITY_Community 120|Community 120]]
- [[_COMMUNITY_Community 121|Community 121]]
- [[_COMMUNITY_Community 122|Community 122]]
- [[_COMMUNITY_Community 123|Community 123]]
- [[_COMMUNITY_Community 124|Community 124]]
- [[_COMMUNITY_Community 125|Community 125]]
- [[_COMMUNITY_Community 126|Community 126]]
- [[_COMMUNITY_Community 127|Community 127]]
- [[_COMMUNITY_Community 128|Community 128]]
- [[_COMMUNITY_Community 129|Community 129]]
- [[_COMMUNITY_Community 130|Community 130]]
- [[_COMMUNITY_Community 131|Community 131]]
- [[_COMMUNITY_Community 132|Community 132]]
- [[_COMMUNITY_Community 133|Community 133]]
- [[_COMMUNITY_Community 134|Community 134]]
- [[_COMMUNITY_Community 197|Community 197]]
- [[_COMMUNITY_Community 198|Community 198]]
- [[_COMMUNITY_Community 199|Community 199]]
- [[_COMMUNITY_Community 200|Community 200]]
- [[_COMMUNITY_Community 201|Community 201]]
- [[_COMMUNITY_Community 202|Community 202]]
- [[_COMMUNITY_Community 203|Community 203]]
- [[_COMMUNITY_Community 204|Community 204]]
- [[_COMMUNITY_Community 205|Community 205]]
- [[_COMMUNITY_Community 206|Community 206]]

## God Nodes (most connected - your core abstractions)
1. `Errorf()` - 414 edges
2. `newTestApp()` - 195 edges
3. `Store` - 186 edges
4. `JSON()` - 127 edges
5. `New()` - 118 edges
6. `contains()` - 98 edges
7. `newDocTestStore()` - 96 edges
8. `Error()` - 85 edges
9. `UserIDFromContext()` - 63 edges
10. `New()` - 63 edges

## Surprising Connections (you probably didn't know these)
- `parseValue()` --calls--> `number()`  [INFERRED]
  web/src/features/compliance/RulesPage.tsx → backend/internal/compliance/engine.go
- `open()` --calls--> `newTestAppWithBackupDir()`  [INFERRED]
  web/src/features/services/ServiceDetailPage.tsx → backend/internal/api/testapp_test.go
- `open()` --calls--> `NewStore()`  [INFERRED]
  web/src/features/services/ServiceDetailPage.tsx → backend/internal/api/apitest/apitest.go
- `open()` --calls--> `newTestStore()`  [INFERRED]
  web/src/features/services/ServiceDetailPage.tsx → backend/internal/api/notifications/handlers_test.go
- `open()` --calls--> `newTestStore()`  [INFERRED]
  web/src/features/services/ServiceDetailPage.tsx → backend/internal/quality/checker_test.go

## Hyperedges (group relationships)
- **Step-Up Confirmation Flow for Destructive Actions** — architecture_permissions_stepup, architecture_destructive_confirm_pattern, openapi_auth_elevate_endpoint, openapi_removal_impact_endpoint [EXTRACTED 0.90]
- **Contract-First API Codegen Pipeline** — architecture_orval, openapi_spec_document, architecture_react_query, architecture_diff_contract [EXTRACTED 0.85]
- **Dual Local/OIDC Auth Mode System** — architecture_auth_design, architecture_oidc_provider_config, openapi_oidc_provider_schema, openapi_auth_config_endpoint [EXTRACTED 0.90]

## Communities

### Community 0 - "Community 0"
Cohesion: 0.01
Nodes (95): newTestHandler(), ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), join(), Errorf(), RequestID(), requestIDKey (+87 more)

### Community 1 - "Community 1"
Cohesion: 0.02
Nodes (77): Handler, NewRouter(), spaHandler(), wsRoleLabel(), Handler, newToken(), sanitize(), Handler (+69 more)

### Community 2 - "Community 2"
Cohesion: 0.01
Nodes (121): claudeProvider, openAICompatibleProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream(), testProvider (+113 more)

### Community 3 - "Community 3"
Cohesion: 0.02
Nodes (213): seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess(), TestAlertsResolveRoleBoundary() (+205 more)

### Community 4 - "Community 4"
Cohesion: 0.03
Nodes (183): bulkSnoozeItemResult, bulkSnoozeRequest, NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound() (+175 more)

### Community 5 - "Community 5"
Cohesion: 0.03
Nodes (98): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+90 more)

### Community 6 - "Community 6"
Cohesion: 0.03
Nodes (79): TestAllConnectorImplementationsRegister(), AIConfigSummary, connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), ExportToFile() (+71 more)

### Community 7 - "Community 7"
Cohesion: 0.04
Nodes (52): RegisterClaude(), TestClaudeSuggestErrors(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, NewEmbedRegistry(), EmbedRegistry (+44 more)

### Community 8 - "Community 8"
Cohesion: 0.04
Nodes (45): AttributeSpec, Catalog, Condition, contains(), equal(), Evaluate(), findAttribute(), matches() (+37 more)

### Community 9 - "Community 9"
Cohesion: 0.07
Nodes (53): RequestedFields(), Register(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists() (+45 more)

### Community 10 - "Community 10"
Cohesion: 0.05
Nodes (53): Config, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), APIKeyChecker, AuditRecorder, ConnectorRoleChecker, contextKey (+45 more)

### Community 11 - "Community 11"
Cohesion: 0.04
Nodes (37): renderAttention(), renderAuthGuard(), renderRoleGuard(), renderChat(), renderPage(), renderPage(), renderTab(), renderDashboard() (+29 more)

### Community 12 - "Community 12"
Cohesion: 0.04
Nodes (27): TestFetchMalformedResponse(), TestFetchNonJSONResponse(), TestFetchPlainJSONWithoutEntities(), TestFetchUsesConfiguredMethodAndHeaders(), TestIsTimeoutDetectsContextDeadlineExceeded(), TestIsTimeoutDetectsNetTimeout(), TestValidateCustomURL(), timeoutError (+19 more)

### Community 13 - "Community 13"
Cohesion: 0.06
Nodes (21): submit(), t(), buildCommands(), onKeyDown(), run(), registeredCommands(), commit(), move() (+13 more)

### Community 14 - "Community 14"
Cohesion: 0.11
Nodes (30): Checker, complianceDescription(), complianceRule(), NewChecker(), RunStaleSweepOnce(), createComplianceRule(), createComplianceSnapshot(), createConnector() (+22 more)

### Community 15 - "Community 15"
Cohesion: 0.14
Nodes (29): channelCfg, Dispatcher, discordPayload(), findChannel(), NewDispatcher(), sendWebhook(), slackPayload(), deliveriesFor() (+21 more)

### Community 16 - "Community 16"
Cohesion: 0.08
Nodes (29): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), clearOIDCFlowCookie(), emailDomainAllowed(), oidcFlowCookieName(), oidcRoleForGroups(), randomOIDCToken(), readOIDCFlowCookie() (+21 more)

### Community 17 - "Community 17"
Cohesion: 0.07
Nodes (20): AISettings, AuthSettings, BackupSettings, Config, Load(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields() (+12 more)

### Community 18 - "Community 18"
Cohesion: 0.14
Nodes (18): DecodeKey(), Decrypt(), DeriveKey(), Encrypt(), TestDecodeKey(), TestDecodeKeyUsableForEncryptDecrypt(), TestDecryptTampered(), TestDecryptWrongKey() (+10 more)

### Community 19 - "Community 19"
Cohesion: 0.08
Nodes (27): ADR 0001 — Monorepo, ADR Index (docs/adr/), AI Doc Generation Module (opt-in, provider-agnostic), API Design — REST + WebSocket split, Dual Auth Design (Local JWT + OIDC), Changes/Diff Contract (infra vs doc format), Change-Aware Diff Engine, Monorepo with Go Workspaces (+19 more)

### Community 20 - "Community 20"
Cohesion: 0.11
Nodes (12): customInstance(), getAccessToken(), setAccessToken(), getConnectorPermissionsQueryKey(), patchUserInstanceAdmin(), postUserInstanceAdmin(), useGetConnectorPermissions(), apply() (+4 more)

### Community 21 - "Community 21"
Cohesion: 0.13
Nodes (6): DBTX, pgPlaceholderDB, pgTransactionDB, rewritePlaceholders(), TestRewritePlaceholders(), transactionDB

### Community 22 - "Community 22"
Cohesion: 0.12
Nodes (20): Axios API client, Button and IconButton, CommandPalette, theme cycling command, ConfirmDialog, Dialog, ElevationConfirm, English translation catalog (+12 more)

### Community 23 - "Community 23"
Cohesion: 0.16
Nodes (20): alerts, changes, connector config JSON, connectors, dashboard layouts, doc versions, docs, HashToken (+12 more)

### Community 24 - "Community 24"
Cohesion: 0.17
Nodes (13): broadcastMsg, Client, Envelope, NewHub(), assertEnvelope(), setupWSConnection(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect() (+5 more)

### Community 25 - "Community 25"
Cohesion: 0.15
Nodes (8): extractGroups(), newMockOIDCServer(), TestAuthURLAfterInitialization(), TestExtractGroups(), TestInitializeSuccess(), TestIsInitializedBeforeAndAfter(), OIDCClaims, OIDCProvider

### Community 26 - "Community 26"
Cohesion: 0.19
Nodes (14): TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare(), lineCount(), relatedServiceIDs(), severityForChange(), sortedKeys() (+6 more)

### Community 27 - "Community 27"
Cohesion: 0.18
Nodes (8): APIKeyClaims, Claims, ElevationClaims, ElevationToken, hasAudience(), newTokenID(), Service, TokenPair

### Community 28 - "Community 28"
Cohesion: 0.23
Nodes (8): bulkFakeConnector, bulkReq(), bulkResults(), createBulkFakeConnector(), registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync()

### Community 29 - "Community 29"
Cohesion: 0.16
Nodes (4): dashboardOverview(), minsAgo(), serviceSnapshot(), syncRunsFor()

### Community 30 - "Community 30"
Cohesion: 0.16
Nodes (15): Connector Interface (Name/Fetch/Validate), Connector Management via UI (full CRUD), Destructive-Action Pattern: Confirm + Blast Radius, Manager Actions (v1 scope), Permissions & Step-Up for Mutating Actions, PRODUCT.md, Role Model — viewer/operator, ServiceSnapshot Data Structure (+7 more)

### Community 31 - "Community 31"
Cohesion: 0.19
Nodes (6): MockWebSocket, env(), heartbeat(), newId(), startTimeline(), syncProgress()

### Community 32 - "Community 32"
Cohesion: 0.22
Nodes (10): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), runTransformers(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), TestRunTransformersUnknownCategoryIsNoop() (+2 more)

### Community 33 - "Community 33"
Cohesion: 0.15
Nodes (14): Connector interface, connector schema registration, reverse proxy WebSocket support, OpenAPI REST contract, destructive-action step-up authentication, operational alerts, detected changes, Compare (+6 more)

### Community 34 - "Community 34"
Cohesion: 0.24
Nodes (11): cosineSimilarity(), packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips(), TestSplitSections() (+3 more)

### Community 35 - "Community 35"
Cohesion: 0.19
Nodes (11): AlertNotifier, TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), DocRegenerator, changePatternID(), computeNextRun() (+3 more)

### Community 36 - "Community 36"
Cohesion: 0.17
Nodes (13): OpenAPI client generation, Authentication and onboarding guards, Operator-only routes, Application root, Application router, Generated-code lint exclusions, Client bootstrap, Conditional mock bootstrap (+5 more)

### Community 37 - "Community 37"
Cohesion: 0.18
Nodes (4): durationLabel(), relativeTime(), relativeTime(), cn()

### Community 38 - "Community 38"
Cohesion: 0.35
Nodes (8): a(), B(), D(), g(), i(), k(), Q(), y()

### Community 39 - "Community 39"
Cohesion: 0.18
Nodes (11): Connector category icon map, Live dashboard state, Dashboard widget frame, Dashboard widgets, Shared SVG icon family, Authenticated app frame, Bottom-dock shell, Floating dock navigation (+3 more)

### Community 40 - "Community 40"
Cohesion: 0.2
Nodes (10): WCAG 2.2 AA accessibility, Docs-first information architecture, technical homelabbers, machine-honest interface, v1 narrow manager scope, trustworthy live documentation, WiseLabz, commit quality gates (+2 more)

### Community 41 - "Community 41"
Cohesion: 0.33
Nodes (7): applyTokens(), loadFont(), makePalette(), presetOpts(), commit(), load(), tokensFor()

### Community 42 - "Community 42"
Cohesion: 0.33
Nodes (6): ShortcutsModal(), useCanMutate(), useConnectorRole(), useIsInstanceAdmin(), useOperatorConnectorIds(), RoleGate()

### Community 43 - "Community 43"
Cohesion: 0.31
Nodes (6): buildDocDiff(), fold(), toUnits(), diffStats(), lineDiff(), wordDiff()

### Community 44 - "Community 44"
Cohesion: 0.32
Nodes (1): Runner

### Community 45 - "Community 45"
Cohesion: 0.38
Nodes (4): navigateTo(), invalidate(), markAllRead(), markRead()

### Community 47 - "Community 47"
Cohesion: 0.38
Nodes (4): loadCache(), persistCache(), widgetsFromWire(), widgetsToWire()

### Community 48 - "Community 48"
Cohesion: 0.62
Nodes (6): getResponse(), handleRequest(), resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 49 - "Community 49"
Cohesion: 0.4
Nodes (2): findRoute(), rowSeverity()

### Community 50 - "Community 50"
Cohesion: 0.6
Nodes (3): fillBody(), generatePreview(), renderTemplate()

### Community 52 - "Community 52"
Cohesion: 0.5
Nodes (2): reportBulkResult(), t()

### Community 54 - "Community 54"
Cohesion: 0.7
Nodes (4): goToNext(), goToPrevious(), makeCurrent(), toggleClass()

### Community 55 - "Community 55"
Cohesion: 0.4
Nodes (1): responseWriter

### Community 56 - "Community 56"
Cohesion: 0.4
Nodes (2): AlertRecord, ChangeRecord

### Community 58 - "Community 58"
Cohesion: 0.5
Nodes (5): Commit Conventions & Hook Enforcement (dev workflow), commit-msg Hook, Conventional Commits Policy, lefthook Commit Hooks, pre-commit Hook

### Community 59 - "Community 59"
Cohesion: 0.83
Nodes (3): close(), onKey(), reset()

### Community 60 - "Community 60"
Cohesion: 0.83
Nodes (3): close(), onKey(), reset()

### Community 63 - "Community 63"
Cohesion: 0.5
Nodes (1): TestWebSocket

### Community 64 - "Community 64"
Cohesion: 0.5
Nodes (2): runSync(), triggerMockSync()

### Community 66 - "Community 66"
Cohesion: 0.5
Nodes (3): ChatConversationRecord, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 67 - "Community 67"
Cohesion: 0.5
Nodes (4): AppShell — Bottom Dock Shell (single variant), Theme Engine — Code Default, User-Overridable, Per-User Dashboard Layout with Admin Default (v2), DashboardLayout Schema (per-user widget layout)

### Community 76 - "Community 76"
Cohesion: 1.0
Nodes (2): focusable(), onKeyDown()

### Community 80 - "Community 80"
Cohesion: 0.67
Nodes (3): Document diff model, Diff layout preference, Diff viewer

### Community 81 - "Community 81"
Cohesion: 0.67
Nodes (3): Destructive connector confirmation, Connector removal impact, Step-up reauthentication

### Community 82 - "Community 82"
Cohesion: 0.67
Nodes (3): PostgreSQL compose deployment, single-instance deployment model, SQLite compose deployment

### Community 83 - "Community 83"
Cohesion: 0.67
Nodes (3): GRAPH_REPORT.md, Graphify Knowledge Graph Rules, Graphify Wiki Index

### Community 84 - "Community 84"
Cohesion: 0.67
Nodes (3): go:embed SPA Embedding, React + Vite Frontend, Frontend Testing Policy (deferred until rewrite)

### Community 85 - "Community 85"
Cohesion: 0.67
Nodes (3): Database: SQLite + PostgreSQL, golang-migrate, sqlc (type-safe SQL codegen)

### Community 86 - "Community 86"
Cohesion: 0.67
Nodes (3): Topbar Notification Center (deferred from V1), NotificationDelivery Schema (per-channel delivery/retry), Notification / NotificationPage Schemas

### Community 119 - "Community 119"
Cohesion: 1.0
Nodes (1): RetentionSettingsRequest

### Community 120 - "Community 120"
Cohesion: 1.0
Nodes (1): APIKey

### Community 121 - "Community 121"
Cohesion: 1.0
Nodes (1): RetentionSettings

### Community 122 - "Community 122"
Cohesion: 1.0
Nodes (1): ShareLink

### Community 123 - "Community 123"
Cohesion: 1.0
Nodes (1): SavedView

### Community 124 - "Community 124"
Cohesion: 1.0
Nodes (2): Template catalog, Template preview generator

### Community 125 - "Community 125"
Cohesion: 1.0
Nodes (2): Change detail synthesizer, Homelab mock data

### Community 126 - "Community 126"
Cohesion: 1.0
Nodes (2): Settings mock data, Notification routing matrix

### Community 127 - "Community 127"
Cohesion: 1.0
Nodes (2): DocTree, Markdown

### Community 128 - "Community 128"
Cohesion: 1.0
Nodes (2): pgPlaceholderDB, rewritePlaceholders

### Community 129 - "Community 129"
Cohesion: 1.0
Nodes (2): safe application defaults, WISELABZ environment overrides

### Community 130 - "Community 130"
Cohesion: 1.0
Nodes (2): Branch Naming Convention, Pull Request Process

### Community 131 - "Community 131"
Cohesion: 1.0
Nodes (2): chi HTTP Router, gorilla/websocket

### Community 132 - "Community 132"
Cohesion: 1.0
Nodes (2): viper Config Loader, WISELABZ_ Env Var Config Override

### Community 133 - "Community 133"
Cohesion: 1.0
Nodes (2): GHCR Container Registry, GitHub Actions CI

### Community 134 - "Community 134"
Cohesion: 1.0
Nodes (2): GET /api/version (undocumented ops endpoint), GET /system/info

### Community 197 - "Community 197"
Cohesion: 1.0
Nodes (1): TimeAgo

### Community 198 - "Community 198"
Cohesion: 1.0
Nodes (1): Panel

### Community 199 - "Community 199"
Cohesion: 1.0
Nodes (1): store package

### Community 200 - "Community 200"
Cohesion: 1.0
Nodes (1): OpenDB

### Community 201 - "Community 201"
Cohesion: 1.0
Nodes (1): ErrNotFound

### Community 202 - "Community 202"
Cohesion: 1.0
Nodes (1): ErrConflict

### Community 203 - "Community 203"
Cohesion: 1.0
Nodes (1): slog (stdlib logging)

### Community 204 - "Community 204"
Cohesion: 1.0
Nodes (1): Zustand State Management

### Community 205 - "Community 205"
Cohesion: 1.0
Nodes (1): Tailwind CSS

### Community 206 - "Community 206"
Cohesion: 1.0
Nodes (1): Docker Compose Deployment

## Ambiguous Edges - Review These
- `Topbar Notification Center (deferred from V1)` → `Notification / NotificationPage Schemas`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to
- `Topbar Notification Center (deferred from V1)` → `NotificationDelivery Schema (per-channel delivery/retry)`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to

## Knowledge Gaps
- **201 isolated node(s):** `testAuditCall`, `contextKey`, `APIKeyChecker`, `UserStatusChecker`, `ConnectorRoleChecker` (+196 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 44`** (8 nodes): `scheduler.go`, `Runner`, `.AddJob()`, `.EntryCount()`, `.jobContext()`, `.RemoveJob()`, `.Start()`, `.Stop()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 49`** (6 nodes): `eventLabel()`, `findRoute()`, `rowSeverity()`, `setCell()`, `setRowSeverity()`, `EventRoutingTable.tsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 52`** (5 nodes): `async()`, `reportBulkResult()`, `t()`, `toggleSelected()`, `ServicesPage.tsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 55`** (5 nodes): `logging.go`, `responseWriter`, `.Hijack()`, `.Write()`, `.WriteHeader()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 56`** (5 nodes): `change.go`, `AlertRecord`, `scanAlert()`, `scanChange()`, `ChangeRecord`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 63`** (4 nodes): `WebSocketProvider.test.tsx`, `TestWebSocket`, `.close()`, `.constructor()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 64`** (4 nodes): `runSync()`, `runSync.ts`, `triggerSync.ts`, `triggerMockSync()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 76`** (3 nodes): `focusable()`, `onKeyDown()`, `DocsPage.tsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 119`** (2 nodes): `retention.go`, `RetentionSettingsRequest`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 120`** (2 nodes): `api_key.go`, `APIKey`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 121`** (2 nodes): `retention_settings.go`, `RetentionSettings`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 122`** (2 nodes): `share_link.go`, `ShareLink`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 123`** (2 nodes): `saved_view.go`, `SavedView`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 124`** (2 nodes): `Template catalog`, `Template preview generator`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 125`** (2 nodes): `Change detail synthesizer`, `Homelab mock data`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 126`** (2 nodes): `Settings mock data`, `Notification routing matrix`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 127`** (2 nodes): `DocTree`, `Markdown`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 128`** (2 nodes): `pgPlaceholderDB`, `rewritePlaceholders`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 129`** (2 nodes): `safe application defaults`, `WISELABZ environment overrides`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 130`** (2 nodes): `Branch Naming Convention`, `Pull Request Process`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 131`** (2 nodes): `chi HTTP Router`, `gorilla/websocket`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 132`** (2 nodes): `viper Config Loader`, `WISELABZ_ Env Var Config Override`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 133`** (2 nodes): `GHCR Container Registry`, `GitHub Actions CI`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 134`** (2 nodes): `GET /api/version (undocumented ops endpoint)`, `GET /system/info`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 197`** (1 nodes): `TimeAgo`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 198`** (1 nodes): `Panel`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 199`** (1 nodes): `store package`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 200`** (1 nodes): `OpenDB`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 201`** (1 nodes): `ErrNotFound`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 202`** (1 nodes): `ErrConflict`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 203`** (1 nodes): `slog (stdlib logging)`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 204`** (1 nodes): `Zustand State Management`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 205`** (1 nodes): `Tailwind CSS`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 206`** (1 nodes): `Docker Compose Deployment`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `Notification / NotificationPage Schemas`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `NotificationDelivery Schema (per-channel delivery/retry)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `Errorf()` connect `Community 0` to `Community 1`, `Community 2`, `Community 34`, `Community 6`, `Community 7`, `Community 8`, `Community 9`, `Community 11`, `Community 44`, `Community 14`, `Community 15`, `Community 16`, `Community 17`, `Community 18`, `Community 25`, `Community 27`?**
  _High betweenness centrality (0.252) - this node is a cross-community bridge._
- **Why does `newTestApp()` connect `Community 3` to `Community 0`, `Community 1`, `Community 7`, `Community 9`, `Community 10`?**
  _High betweenness centrality (0.127) - this node is a cross-community bridge._
- **Why does `contains()` connect `Community 4` to `Community 0`, `Community 32`, `Community 2`, `Community 3`, `Community 5`, `Community 6`, `Community 7`, `Community 8`, `Community 9`, `Community 10`, `Community 12`, `Community 14`, `Community 15`, `Community 21`, `Community 25`, `Community 28`?**
  _High betweenness centrality (0.078) - this node is a cross-community bridge._
- **Are the 412 inferred relationships involving `Errorf()` (e.g. with `WriteElevationError()` and `.IssuePair()`) actually correct?**
  _`Errorf()` has 412 INFERRED edges - model-reasoned connections that need verification._
- **Are the 186 inferred relationships involving `newTestApp()` (e.g. with `TestAlertsListSuccess()` and `TestAlertsResolveRoleBoundary()`) actually correct?**
  _`newTestApp()` has 186 INFERRED edges - model-reasoned connections that need verification._