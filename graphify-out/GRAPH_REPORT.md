# Graph Report - WiseLabz  (2026-09-30)

## Corpus Check
- 880 files · ~556,871 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 20 file(s) not represented in the graph (top: (none) 9, .toml 2, .tmpl 2)

## Summary
- 6734 nodes · 21079 edges · 257 communities (224 shown, 33 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1807 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0926eb38`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Client
- context.Context
- net/http.Request
- ServiceSnapshot
- NewStore
- backup/backup.go
- time.Duration
- Elector
- compliance/handlers.go
- Store
- diagnostics/diagnostics.go
- NewEngine
- newTestApp
- runbooks_test.go
- time.Time
- newDocTestStore
- DecodeKey
- Dispatcher
- WiseLabz Project
- SuggestRequest
- testing.T
- CommandPalette
- initial database schema
- rewritePlaceholders
- WebSocketProvider.tsx
- .OIDCCallback
- react
- Service
- Register
- Handler
- Product
- @tanstack/react-query
- newTestHandler
- RunSync
- chat/chat.go
- Engine
- Root
- go_pkg_context
- DashboardPage.tsx
- Bottom-dock shell
- RunMigrations
- ProfilePage.tsx
- WiseLabz
- fetch_test.go
- rowScanner
- ShareLinkPage.tsx
- Handler
- router.go
- App.tsx
- package.json
- createTestConnector
- ErrorWithDetails
- home_assistant/tables.go
- RunbookRecord
- cn
- fixtures.ts
- NewHTTPClient
- dispatcher_test.go
- icons.tsx
- lefthook Commit Hooks
- AppShell.tsx
- Get
- Registry
- go_pkg_strings
- routerDeps
- validate.go
- SnapshotEntity
- Manager
- AppShell — Bottom Dock Shell (single variant)
- NewUser
- export_test.go
- Connector
- ServiceDetailPage.tsx
- react-i18next
- newTestHandler
- dependencies
- NewMalformedResponseError
- DecodeJSON
- Config
- compliance/engine.go
- NewEngine
- Diff viewer
- Destructive connector confirmation
- single-instance deployment model
- Graphify Knowledge Graph Rules
- Topbar Notification Center (deferred from V1)
- Database: SQLite + PostgreSQL
- React + Vite Frontend
- portainer/tables.go
- ConnectorRecord
- home_assistant_test.go
- adguardhome/tables.go
- Checker
- traefik/tables_test.go
- Connector
- truenas_test.go
- unifi/tables_test.go
- WiseLabz — Deployment Guide
- Compare
- snapshotreport.go
- settings.mock.ts
- handlers.ts
- Sanitize
- api/mcp_test.go
- timeline.ts
- createUser
- go_pkg_reflect
- go_pkg_os
- unifi_test.go
- SystemPage.tsx
- nilToStr
- src/theme.ts
- VerifyBundleFile
- WiseLabz — Design Contract
- devDependencies
- response.go
- .runPerRevision
- portainer_test.go
- gitFixture
- pagination_contract_test.go
- go_pkg_net_http
- AuthMiddleware
- manifest.go
- notifications/handlers_test.go
- traefik_test.go
- Template catalog
- Change detail synthesizer
- Settings mock data
- DocTree
- pgPlaceholderDB
- safe application defaults
- Branch Naming Convention
- viper Config Loader
- chi HTTP Router
- GHCR Container Registry
- GET /api/version (undocumented ops endpoint)
- adguardhome_test.go
- scheduler/health_test.go
- connector_permission.go
- NewService
- Store
- config_test.go
- Connector
- AppearancePage.tsx
- httpx/retry_test.go
- Deps
- logging_test.go
- newDockerClient
- docker_test.go
- Hub
- ws.ts
- compilerOptions
- New
- changes/handlers_test.go
- WiseLabz Connector Guide
- docdiffmodel.ts
- middleware.go
- all.go
- connector/connector.go
- Backend test performance
- compilerOptions
- HashPassword
- Store
- NewRegistry
- reports/handlers.go
- vectorCache
- templatefuncs.go
- api/changes_test.go
- ListSchemas
- Handler
- api/auth/oidc.go
- connector_permission_test.go
- Contributor Covenant Code of Conduct
- cursor_pagination_test.go
- connectors_maintenance_test.go
- Store
- backup/backup_test.go
- Decision
- Decision
- scripts
- ExportToFile
- api/docs_test.go
- apikey_scopes_test.go
- Decision
- alerts/handlers_authz_test.go
- spec.go
- newTestHandler
- connectors_hardening_test.go
- dashboard/handlers_test.go
- Changelog
- mockServiceWorker.js
- WiseLabz WebSocket Contract (`/ws`)
- webAuthnUser
- gitAuth
- release-please-config.json
- ratelimit.go
- Decision
- Cache
- TimeAgo
- Panel
- store package
- OpenDB
- ErrNotFound
- ErrConflict
- slog (stdlib logging)
- Zustand State Management
- Tailwind CSS
- Docker Compose Deployment
- main
- httpx/retry.go
- vectorcache_test.go
- cloudflare/attributes_test.go
- doc_test.go
- scanMaintenanceWindow
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- keyset_test.go
- Mermaid.tsx
- store/mfa_test.go
- Security Policy
- NewWebAuthnService
- Notification Channels
- compose-smoke.sh
- OIDC Provider Config Decision (file-defined, app toggles only)
- go_pkg_testing
- retention/retention_test.go
- Saved Views
- WiseLabz — v2 Backlog
- tsconfig.json
- setup-env.sh
- CHANGE_PROVENANCE.md
- changes.sh
- zfsProp
- vite-env.d.ts
- github.com/WiseLabz/wiselabz
- test-shards.sh
- golden_snapshot_test.go
- computeNextRun
- notification_delivery_test.go
- main.tsx
- OpenDB
- PRODUCT.md
- seedScopeFixture
- New
- 0004 — PostgreSQL leader election for background workers
- Permissions & Step-Up for Mutating Actions
- Configuration & Documentation Backup (Export/Import)
- Backup Recovery: What Comes Back, and What Doesn't
- fakeDocRegenerator
- log/slog.Logger
- Connector Interface (Name/Fetch/Validate)
- 0001 — Lab-mutating operation boundaries
- BACKUP.md
- fakeQualityChecker
- coverage-parity.sh
- coverpkg.sh

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 231 edges
2. `Errorf()` - 184 edges
3. `newDocTestStore()` - 144 edges
4. `Store` - 142 edges
5. `UserIDFromContext()` - 86 edges
6. `SnapshotEntity` - 77 edges
7. `react` - 77 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
10. `@tanstack/react-query` - 64 edges

## Surprising Connections (you probably didn't know these)
- `Panel (`Panel.tsx`)` --references--> `Panel()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `Radii — rounded but tight. Soft-dark, not pill-everything.` --references--> `Panel()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `Surfaces — depth from lightness steps + shadow, never borders alone` --references--> `Panel()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `4. Typography` --references--> `PanelHeader()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `9. Anti-slop bans` --references--> `PanelHeader()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Contract-First API Codegen Pipeline** — docs_architecture_orval, docs_openapi_spec_document, docs_architecture_react_query, docs_architecture_diff_contract [EXTRACTED 0.85]
- **Dual Local/OIDC Auth Mode System** — docs_architecture_auth_design, docs_architecture_oidc_provider_config, docs_openapi_oidc_provider_schema, docs_openapi_auth_config_endpoint [EXTRACTED 0.90]
- **Step-Up Confirmation Flow for Destructive Actions** — docs_architecture_permissions_stepup, docs_architecture_destructive_confirm_pattern, docs_openapi_auth_elevate_endpoint, docs_openapi_removal_impact_endpoint [EXTRACTED 0.90]

## Communities (257 total, 33 thin omitted)

### Community 0 - "net/http.Client"
Cohesion: 0.03
Nodes (26): Connector, ollamaEmbedder, openAIEmbedder, NewServiceUnavailableError(), setHeaders(), tryParseEntities(), Connector, Connector (+18 more)

### Community 1 - "context.Context"
Cohesion: 0.02
Nodes (38): fakeStatusChecker, sanitizeSessions(), Connector, existingIDs(), Store, scanComplianceRule(), Store, docSearchWhere() (+30 more)

### Community 2 - "net/http.Request"
Cohesion: 0.05
Nodes (47): Handler, oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName() (+39 more)

### Community 3 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (32): noopValidatedConnector, ServiceDependency, ServiceSnapshot, WantsField(), Connector, environmentDependencies(), putMetadata(), unavailable() (+24 more)

### Community 4 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 5 - "backup/backup.go"
Cohesion: 0.20
Nodes (22): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+14 more)

### Community 6 - "time.Duration"
Cohesion: 0.10
Nodes (22): Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings, BackupSettings, Database, DocExportGitSettings (+14 more)

### Community 8 - "compliance/handlers.go"
Cohesion: 0.22
Nodes (10): catalog(), changedFields(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord, Handler (+2 more)

### Community 9 - "Store"
Cohesion: 0.07
Nodes (13): changeServiceIDs(), placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanChange(), changePatternID() (+5 more)

### Community 10 - "diagnostics/diagnostics.go"
Cohesion: 0.32
Nodes (12): CheckHealth(), Collect(), collectVersions(), AuthProviders, Bundle, Component, Health, OIDCProviderSummary (+4 more)

### Community 11 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 12 - "newTestApp"
Cohesion: 0.02
Nodes (168): healthFakeConnector, templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation() (+160 more)

### Community 13 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 14 - "time.Time"
Cohesion: 0.06
Nodes (49): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), connectorFilter(), RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles() (+41 more)

### Community 15 - "newDocTestStore"
Cohesion: 0.05
Nodes (58): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+50 more)

### Community 16 - "DecodeKey"
Cohesion: 0.06
Nodes (35): factorJSON(), Handler, testHandler, Handler, Handler, Handler, Handler, GenerateRecoveryCodes() (+27 more)

### Community 17 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 18 - "WiseLabz Project"
Cohesion: 0.12
Nodes (16): ADR 0001 — Monorepo, ADR Index (docs/adr/), AI Doc Generation Module (opt-in, provider-agnostic), Changes/Diff Contract (infra vs doc format), Change-Aware Diff Engine, Monorepo with Go Workspaces, AI-Suggestion Token Streaming (dropped), Soft-Lock on Concurrent Doc Editing (deferred) (+8 more)

### Community 19 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 20 - "testing.T"
Cohesion: 0.03
Nodes (126): TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider(), TestGetOrInitOIDCProvider(), TestOIDCCallbackRejectsMissingFlowCookie(), TestOIDCCallbackRejectsStateMismatchedWithCookie(), TestOIDCCallbackRejectsUnknownProvider() (+118 more)

### Community 21 - "CommandPalette"
Cohesion: 0.12
Nodes (20): Axios API client, Button and IconButton, CommandPalette, theme cycling command, ConfirmDialog, Dialog, ElevationConfirm, English translation catalog (+12 more)

### Community 22 - "initial database schema"
Cohesion: 0.16
Nodes (20): alerts, changes, connector config JSON, connectors, dashboard layouts, doc versions, docs, HashToken (+12 more)

### Community 23 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 24 - "WebSocketProvider.tsx"
Cohesion: 0.04
Nodes (66): Live dashboard state, RFC-3339, 15. `system.resync`, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_changes_changes, web_src_api_generated_changes_changes_getgetchangesquerykey, web_src_api_generated_connectors_connectors (+58 more)

### Community 25 - ".OIDCCallback"
Cohesion: 0.13
Nodes (10): Handler, Handler, newOIDCUser(), randomOIDCToken(), readOIDCFlowCookie(), OIDCClaims, OIDCProvider, OIDCProvider (+2 more)

### Community 26 - "react"
Cohesion: 0.03
Nodes (138): react, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules (+130 more)

### Community 27 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 28 - "Register"
Cohesion: 0.14
Nodes (22): init(), init(), init(), init(), init(), init(), init(), init() (+14 more)

### Community 29 - "Handler"
Cohesion: 0.26
Nodes (4): createVersion(), templateResponse(), templateVersionResponse(), Handler

### Community 30 - "Product"
Cohesion: 0.20
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 31 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (58): msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete, web_src_api_model_index_attentionpage (+50 more)

### Community 32 - "newTestHandler"
Cohesion: 0.06
Nodes (72): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+64 more)

### Community 33 - "RunSync"
Cohesion: 0.15
Nodes (14): Connector interface, connector schema registration, reverse proxy WebSocket support, OpenAPI REST contract, destructive-action step-up authentication, operational alerts, detected changes, Compare (+6 more)

### Community 34 - "chat/chat.go"
Cohesion: 0.15
Nodes (15): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), SplitSections(), SyncDocEmbeddings() (+7 more)

### Community 35 - "Engine"
Cohesion: 0.15
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 36 - "Root"
Cohesion: 0.18
Nodes (12): OpenAPI client generation, Generated-code lint exclusions, Motion preference provider, Vite API and WebSocket proxy, Authentication and onboarding guards, Operator-only routes, Root(), router (+4 more)

### Community 37 - "go_pkg_context"
Cohesion: 0.08
Nodes (22): dashboardLayout, versionSections(), ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), digestDue(), TestDigestDue(), TemplateVersionSection (+14 more)

### Community 38 - "DashboardPage.tsx"
Cohesion: 0.05
Nodes (65): Connector category icon map, Dashboard widget frame, 4. `change.detected`, 5. `alert.created`, 8. `doc.generated`, Shared SVG icon family, web_src_api_generated_dashboard_dashboard_getdashboardlayout, web_src_api_generated_dashboard_dashboard_getdashboardlayoutadmindefault (+57 more)

### Community 39 - "Bottom-dock shell"
Cohesion: 0.40
Nodes (5): Authenticated app frame, Bottom-dock shell, Primary navigation, Non-React navigation bridge, Floating dock navigation

### Community 40 - "RunMigrations"
Cohesion: 0.11
Nodes (36): main(), GetMigrationStatus(), newMigrator(), collectColumns(), migrationFiles(), postgresSchemaColumns(), sqliteSchemaColumns(), TestMigrationSchemaParity() (+28 more)

### Community 41 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (60): setAccessToken(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin, web_src_api_generated_auth_auth_postauthlogin (+52 more)

### Community 42 - "WiseLabz"
Cohesion: 0.20
Nodes (10): WCAG 2.2 AA accessibility, Docs-first information architecture, technical homelabbers, machine-honest interface, v1 narrow manager scope, trustworthy live documentation, WiseLabz, commit quality gates (+2 more)

### Community 43 - "fetch_test.go"
Cohesion: 0.19
Nodes (14): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+6 more)

### Community 44 - "rowScanner"
Cohesion: 0.07
Nodes (24): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), Store, scanBackupRun(), scanAlert() (+16 more)

### Community 45 - "ShareLinkPage.tsx"
Cohesion: 0.07
Nodes (43): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+35 more)

### Community 46 - "Handler"
Cohesion: 0.22
Nodes (7): validateConfigPushRequest(), stepAuditDetail(), validTargetType(), ValidateCompositeRef(), Handler, runbookResponse, stepResponse

### Community 47 - "router.go"
Cohesion: 0.10
Nodes (30): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), validVerb(), chi.Routes, go_pkg_github_com_go_chi_chi_v5, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts (+22 more)

### Community 48 - "App.tsx"
Cohesion: 0.04
Nodes (62): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 6. `alert.resolved`, 7a. `finding.created`, 7b. `digest.summary` (+54 more)

### Community 49 - "package.json"
Cohesion: 0.04
Nodes (50): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+42 more)

### Community 50 - "createTestConnector"
Cohesion: 0.11
Nodes (24): TestSnapshotKeysetSummaryAndConnectorOwnership(), TestDeleteOldHealthChecks(), TestGetConnectorUptimeDeterministicOutage(), TestGetConnectorUptimeNoData(), TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR(), TestRecordHealthCheckDefaults(), assertQueryPlanUsesIndex(), Store (+16 more)

### Community 51 - "ErrorWithDetails"
Cohesion: 0.10
Nodes (20): Handler, sanitize(), Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, ClientIP() (+12 more)

### Community 52 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 53 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 54 - "cn"
Cohesion: 0.03
Nodes (83): web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid, web_src_api_generated_docs_docs (+75 more)

### Community 55 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 56 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 57 - "dispatcher_test.go"
Cohesion: 0.07
Nodes (70): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext() (+62 more)

### Community 58 - "icons.tsx"
Cohesion: 0.07
Nodes (49): web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidversionsrevrestore, web_src_api_generated_docs_docs_usegetdocsdocid, web_src_api_generated_docs_docs_usegetdocsdocidversions, web_src_api_generated_docs_docs_usegetdocsdocidversionsrev, web_src_api_generated_docs_docs_usegetdocstree (+41 more)

### Community 59 - "lefthook Commit Hooks"
Cohesion: 0.50
Nodes (5): commit-msg Hook, Conventional Commits Policy, lefthook Commit Hooks, pre-commit Hook, Commit Conventions & Hook Enforcement (dev workflow)

### Community 60 - "AppShell.tsx"
Cohesion: 0.09
Nodes (22): i18next, react-error-boundary, sonner, AppShell, CommandPalette(), Command, CommandCtx, CommandFactory (+14 more)

### Community 61 - "Get"
Cohesion: 0.08
Nodes (28): applyConnectorScalarUpdates(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection(), LifecycleOp(), Get() (+20 more)

### Community 62 - "Registry"
Cohesion: 0.14
Nodes (16): Provider, StatusError, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds() (+8 more)

### Community 63 - "go_pkg_strings"
Cohesion: 0.07
Nodes (31): TestAttributeCatalogCoversEmittedKeys(), TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), NoRedirect(), buildEmailMessage(), sendSMTPChannel() (+23 more)

### Community 64 - "routerDeps"
Cohesion: 0.10
Nodes (33): routerDeps, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+25 more)

### Community 65 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 66 - "SnapshotEntity"
Cohesion: 0.14
Nodes (41): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+33 more)

### Community 67 - "Manager"
Cohesion: 0.16
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 68 - "AppShell — Bottom Dock Shell (single variant)"
Cohesion: 0.50
Nodes (4): AppShell — Bottom Dock Shell (single variant), Theme Engine — Code Default, User-Overridable, Per-User Dashboard Layout with Admin Default (v2), DashboardLayout Schema (per-user widget layout)

### Community 69 - "NewUser"
Cohesion: 0.07
Nodes (83): TestList(), TestRevoke(), GrantConnectorRole(), instanceAdminRole(), JWTService(), NewUser(), Token(), WithAuth() (+75 more)

### Community 70 - "export_test.go"
Cohesion: 0.19
Nodes (20): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+12 more)

### Community 71 - "Connector"
Cohesion: 0.09
Nodes (9): init(), ConfigField, Connector, buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment(), ValidateRefSegment() (+1 more)

### Community 72 - "ServiceDetailPage.tsx"
Cohesion: 0.04
Nodes (62): ADR-0001, ADR-0003, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart (+54 more)

### Community 73 - "react-i18next"
Cohesion: 0.04
Nodes (79): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover, react-i18next, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve (+71 more)

### Community 74 - "newTestHandler"
Cohesion: 0.07
Nodes (59): AssertMatchesSpec(), actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions() (+51 more)

### Community 75 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 76 - "NewMalformedResponseError"
Cohesion: 0.07
Nodes (43): SnapshotSection, NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+35 more)

### Community 77 - "DecodeJSON"
Cohesion: 0.18
Nodes (8): webAuthnFlow, Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T, github.com/go-webauthn/webauthn/webauthn.SessionData

### Community 78 - "Config"
Cohesion: 0.15
Nodes (13): Config, TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders() (+5 more)

### Community 79 - "compliance/engine.go"
Cohesion: 0.12
Nodes (29): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+21 more)

### Community 80 - "NewEngine"
Cohesion: 0.16
Nodes (28): RequestedFields(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+20 more)

### Community 81 - "Diff viewer"
Cohesion: 0.67
Nodes (3): Document diff model, Diff layout preference, Diff viewer

### Community 82 - "Destructive connector confirmation"
Cohesion: 0.67
Nodes (3): Destructive connector confirmation, Connector removal impact, Step-up reauthentication

### Community 83 - "single-instance deployment model"
Cohesion: 0.67
Nodes (3): PostgreSQL compose deployment, single-instance deployment model, SQLite compose deployment

### Community 84 - "Graphify Knowledge Graph Rules"
Cohesion: 0.67
Nodes (3): GRAPH_REPORT.md, Graphify Knowledge Graph Rules, Graphify Wiki Index

### Community 85 - "Topbar Notification Center (deferred from V1)"
Cohesion: 0.67
Nodes (3): Topbar Notification Center (deferred from V1), NotificationDelivery Schema (per-channel delivery/retry), Notification / NotificationPage Schemas

### Community 86 - "Database: SQLite + PostgreSQL"
Cohesion: 0.67
Nodes (3): Database: SQLite + PostgreSQL, golang-migrate, sqlc (type-safe SQL codegen)

### Community 87 - "React + Vite Frontend"
Cohesion: 0.67
Nodes (3): Frontend Testing Policy (deferred until rewrite), go:embed SPA Embedding, React + Vite Frontend

### Community 88 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 89 - "ConnectorRecord"
Cohesion: 0.09
Nodes (48): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+40 more)

### Community 90 - "home_assistant_test.go"
Cohesion: 0.09
Nodes (40): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+32 more)

### Community 91 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 92 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 93 - "traefik/tables_test.go"
Cohesion: 0.10
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 94 - "Connector"
Cohesion: 0.10
Nodes (13): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap(), Connector, apiMessage(), controllerName(), countByKind(), statusError() (+5 more)

### Community 95 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 96 - "unifi/tables_test.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 97 - "WiseLabz — Deployment Guide"
Cohesion: 0.33
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 98 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 99 - "snapshotreport.go"
Cohesion: 0.14
Nodes (26): BuildSnapshotDiff(), CompareDependencies(), CompareEntities(), entityKey(), entityMap(), TestCompareDependenciesDeterministic(), TestCompareEntitiesKeepsExternalIDAndFallbackNameDistinct(), TestSnapshotDiffCSVNeutralizesFormulaCellsOnly() (+18 more)

### Community 100 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 101 - "handlers.ts"
Cohesion: 0.07
Nodes (29): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+21 more)

### Community 102 - "Sanitize"
Cohesion: 0.09
Nodes (18): Handler, isWritableField(), loggablePath(), loggableQuery(), WriteElevationError(), ConfigPusher, Err(), Sanitize() (+10 more)

### Community 103 - "api/mcp_test.go"
Cohesion: 0.15
Nodes (12): testApp, mintAPIKey(), TestMCPConnectorRestrictedKey(), TestMCPEndToEnd(), go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server (+4 more)

### Community 104 - "timeline.ts"
Cohesion: 0.14
Nodes (16): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+8 more)

### Community 105 - "createUser"
Cohesion: 0.17
Nodes (15): ContextWithAPIKeyRestriction(), TestClampConnectorRole(), seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding() (+7 more)

### Community 106 - "go_pkg_reflect"
Cohesion: 0.05
Nodes (31): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+23 more)

### Community 107 - "go_pkg_os"
Cohesion: 0.06
Nodes (38): newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets(), keys(), writeExportState(), New() (+30 more)

### Community 108 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 109 - "SystemPage.tsx"
Cohesion: 0.07
Nodes (32): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey, web_src_api_generated_system_system_getsystembackupschedule (+24 more)

### Community 110 - "nilToStr"
Cohesion: 0.07
Nodes (14): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+6 more)

### Community 111 - "src/theme.ts"
Cohesion: 0.14
Nodes (25): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+17 more)

### Community 112 - "VerifyBundleFile"
Cohesion: 0.22
Nodes (17): failVerification(), LatestBundle(), RunVerifyOnce(), ListVerifications(), RecordVerification(), seedOneDoc(), TestLatestBundleNoBundles(), TestLatestBundlePicksMostRecent() (+9 more)

### Community 113 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 114 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 115 - "response.go"
Cohesion: 0.08
Nodes (27): decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+19 more)

### Community 116 - ".runPerRevision"
Cohesion: 0.22
Nodes (9): commitMessage(), Exporter, TestCommitMessage(), commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash, github.com/go-git/go-git/v5/plumbing/object.Signature (+1 more)

### Community 117 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 118 - "gitFixture"
Cohesion: 0.31
Nodes (11): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), TestGitPerRevisionBootstrapReplayAndCap(), TestGitPerRevisionBotCatchUpAndRejectedPush(), TestGitPerRevisionMissingIntermediateAndEngineAuthor() (+3 more)

### Community 119 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 120 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (48): bulkSnoozeItemResult, bulkSnoozeRequest, updateUserRequest, newToken(), changePromptData(), diffToSpec(), stripPromptTags(), truncateUTF8() (+40 more)

### Community 121 - "AuthMiddleware"
Cohesion: 0.14
Nodes (18): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, AuthMiddleware(), assertElevationAuditCalls(), boolLabel(), requestWithUser() (+10 more)

### Community 122 - "manifest.go"
Cohesion: 0.29
Nodes (12): ImportFromFile(), BuildManifest(), BundleCounts(), ChecksumBytes(), ManifestPath(), ReadManifest(), corruptFile(), TestImportFromFileVerifiesChecksum() (+4 more)

### Community 123 - "notifications/handlers_test.go"
Cohesion: 0.26
Nodes (16): TestCreate(), AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty() (+8 more)

### Community 124 - "traefik_test.go"
Cohesion: 0.09
Nodes (40): RedactConnectorConfig(), TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+32 more)

### Community 136 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 137 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 138 - "connector_permission.go"
Cohesion: 0.11
Nodes (17): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), RejectRestrictedAPIKey(), treatAsSafeFromContext(), Store, apiKeyConnectorFilter() (+9 more)

### Community 139 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 140 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 141 - "config_test.go"
Cohesion: 0.09
Nodes (27): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+19 more)

### Community 142 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 143 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 144 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 145 - "Deps"
Cohesion: 0.22
Nodes (19): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs(), registerListFindings() (+11 more)

### Community 146 - "logging_test.go"
Cohesion: 0.18
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 147 - "newDockerClient"
Cohesion: 0.13
Nodes (14): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+6 more)

### Community 148 - "docker_test.go"
Cohesion: 0.08
Nodes (34): generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn(), startSSHDockerServer(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout() (+26 more)

### Community 149 - "Hub"
Cohesion: 0.06
Nodes (51): TestBroadcastDocEventScoping(), testApp, TestWebSocketConnectorEventsFilteredByGrant(), TestWebSocketRevokedKeyFailsRevalidation(), TestWebSocketTicketFlow(), waitForClients(), wsDial(), wsRead() (+43 more)

### Community 150 - "ws.ts"
Cohesion: 0.11
Nodes (18): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+10 more)

### Community 151 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 152 - "New"
Cohesion: 0.34
Nodes (14): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+6 more)

### Community 153 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 154 - "WiseLabz Connector Guide"
Cohesion: 0.08
Nodes (24): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Conventions (+16 more)

### Community 155 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 156 - "middleware.go"
Cohesion: 0.09
Nodes (22): AuditRecorder, ConnectorRoleChecker, contextKey, elevationError, PermissionChecker, UserStatusChecker, TreatAsSafeMethod(), APIKeyIDFromContext() (+14 more)

### Community 157 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 158 - "connector/connector.go"
Cohesion: 0.09
Nodes (15): Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP(), NewTimeoutError(), newWebhookClient(), CredentialRefresher (+7 more)

### Community 159 - "Backend test performance"
Cohesion: 0.10
Nodes (19): Backend test performance, CI job times, CI measurements, Coverage strategy, Dependency updates and action pinning, Deterministic scheduled exports (#408), Fixture reuse and lifecycle tests (#405–#407), Follow-ups (+11 more)

### Community 160 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 161 - "HashPassword"
Cohesion: 0.09
Nodes (22): mustHashDummyPassword(), instanceAdminRoleFor(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable() (+14 more)

### Community 162 - "Store"
Cohesion: 0.13
Nodes (16): Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+8 more)

### Community 163 - "NewRegistry"
Cohesion: 0.45
Nodes (12): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+4 more)

### Community 164 - "reports/handlers.go"
Cohesion: 0.28
Nodes (7): definition(), NewHandler(), record(), reportJSON(), valid(), Handler, input

### Community 165 - "vectorCache"
Cohesion: 0.25
Nodes (6): vectorCache, vectorEntry, vectorKey, go_pkg_container_list, container/list.Element, container/list.List

### Community 166 - "templatefuncs.go"
Cohesion: 0.17
Nodes (12): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+4 more)

### Community 167 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 168 - "ListSchemas"
Cohesion: 0.21
Nodes (10): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), supportedLifecycleVerbs(), IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType() (+2 more)

### Community 169 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 170 - "api/auth/oidc.go"
Cohesion: 0.24
Nodes (9): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), validHostPort() (+1 more)

### Community 171 - "connector_permission_test.go"
Cohesion: 0.33
Nodes (15): Store, newTestConnector(), newTestUser(), TestConnectorReaderIDs(), TestDeleteConnectorGrant(), TestFilterConnectorIDsByGrant(), TestGetUserConnectorRoleHighestAcrossSources(), TestListConnectorGrants() (+7 more)

### Community 172 - "Contributor Covenant Code of Conduct"
Cohesion: 0.15
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 173 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 174 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 175 - "Store"
Cohesion: 0.18
Nodes (7): testAPIKeyChecker, APIKeyClaims, ValidAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 176 - "backup/backup_test.go"
Cohesion: 0.25
Nodes (14): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+6 more)

### Community 177 - "Decision"
Cohesion: 0.25
Nodes (8): Audit, Authorization, Confirmation / step-up, Decision, Dry-run, Eligible first operation: `service.restart` only, Out of scope, Rollback expectations

### Community 178 - "Decision"
Cohesion: 0.25
Nodes (8): Audit, Authorization, Confirmation / step-up, Decision, Dry-run, Eligible ops: `service.start`, `service.stop`, Out of scope, Rollback expectations

### Community 179 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 180 - "ExportToFile"
Cohesion: 0.18
Nodes (16): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+8 more)

### Community 181 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 182 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 183 - "Decision"
Cohesion: 0.20
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 184 - "alerts/handlers_authz_test.go"
Cohesion: 0.20
Nodes (14): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz(), Handler (+6 more)

### Community 185 - "spec.go"
Cohesion: 0.25
Nodes (8): loadSpec(), specPath(), go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter, go_pkg_github_com_getkin_kin_openapi_routers, go_pkg_github_com_getkin_kin_openapi_routers_gorillamux, go_pkg_runtime, github.com/getkin/kin-openapi/routers.Router

### Community 186 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 187 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 188 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 189 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 190 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 191 - "WiseLabz WebSocket Contract (`/ws`)"
Cohesion: 0.25
Nodes (8): Client dispatch model, Delivery, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`)

### Community 192 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 193 - "gitAuth"
Cohesion: 0.33
Nodes (6): gitAuth(), installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), GitOptions, github.com/go-git/go-git/v5/plumbing/transport.AuthMethod

### Community 194 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 195 - "ratelimit.go"
Cohesion: 0.27
Nodes (7): TestRateLimit(), RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 196 - "Decision"
Cohesion: 0.17
Nodes (12): 0005 — Cross-replica WebSocket event relay for active/active, Authorization and secrets, Consequences, Context, Decision, Duplicate suppression and ordering, Event ownership: local first, then relay, Mechanism: PostgreSQL LISTEN/NOTIFY (+4 more)

### Community 197 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 208 - "main"
Cohesion: 0.08
Nodes (29): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+21 more)

### Community 209 - "httpx/retry.go"
Cohesion: 0.26
Nodes (9): isSafeMethod(), IsSafeMethod(), idempotent(), retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport (+1 more)

### Community 210 - "vectorcache_test.go"
Cohesion: 0.53
Nodes (5): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace()

### Community 211 - "cloudflare/attributes_test.go"
Cohesion: 0.33
Nodes (5): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable()

### Community 212 - "doc_test.go"
Cohesion: 0.13
Nodes (20): Store, mustCreateUser(), newPostgresTestStore(), skipOnPostgres(), TestDocLockAcquireAfterExpiry(), TestDocLockConflict(), TestDocLockReleaseOnlyByHolder(), TestDocLockRenewalByHolder() (+12 more)

### Community 213 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 214 - "Audit Trail"
Cohesion: 0.25
Nodes (7): Audit Trail, Endpoint, Filtering and export, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 215 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 216 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 217 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 218 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 219 - "store/mfa_test.go"
Cohesion: 0.31
Nodes (10): Store, mfaTestUser(), TestConfirmFactorRejectsSecondConfirmedTOTP(), TestConsumeTOTPStepReplayGuard(), TestDeleteFactorLastOneAlsoLeavesRecoveryCodesForCallerToWipe(), TestDeleteUserCascadesMFAFactorsAndRecoveryCodes(), TestDeleteUserFactorsForAdminReset(), TestGetRequire2FADefaultsToNone() (+2 more)

### Community 220 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 221 - "NewWebAuthnService"
Cohesion: 0.50
Nodes (4): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), github.com/go-webauthn/webauthn/webauthn.WebAuthn

### Community 222 - "Notification Channels"
Cohesion: 0.40
Nodes (4): Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 223 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 224 - "OIDC Provider Config Decision (file-defined, app toggles only)"
Cohesion: 0.22
Nodes (10): API Design — REST + WebSocket split, Dual Auth Design (Local JWT + OIDC), OIDC Provider Config Decision (file-defined, app toggles only), orval Generated API Client, React Query (server state), http.Flusher on middleware.responseWriter (TODO), GET/PUT /auth/config, docs/FRONTEND_PLAN.md (+2 more)

### Community 225 - "go_pkg_testing"
Cohesion: 0.04
Nodes (36): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation() (+28 more)

### Community 226 - "retention/retention_test.go"
Cohesion: 0.35
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 227 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 232 - "changes.sh"
Cohesion: 0.40
Nodes (9): classify_path(), cmd_check(), cmd_classify(), cmd_go_closure(), describe_flags(), die(), go_package_dirs(), package_version_only() (+1 more)

### Community 236 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 239 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 240 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 245 - "notification_delivery_test.go"
Cohesion: 0.39
Nodes (7): createTestNotification(), Store, TestDeliveryCreateAndList(), TestListDeliveriesStatusFilterAndPagination(), TestListDueDeliveries(), TestUpdateDeliveryResultNotFound(), TestUpdateDeliveryResultTransitionsAndClearsNextAttempt()

### Community 246 - "main.tsx"
Cohesion: 0.40
Nodes (4): react-dom, App(), USE_MOCKS, web_src_index

### Community 248 - "OpenDB"
Cohesion: 0.18
Nodes (12): OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults(), TestWithinTransactionRollsBack(), Store, newCascadeTestStore(), TestGetRunbookByTarget() (+4 more)

### Community 249 - "PRODUCT.md"
Cohesion: 0.25
Nodes (5): 0002 — Start/stop lab-mutating operations, Consequences, Context, Manager Actions (v1 scope), Lab-Mutating Manager Ops (out of v1 scope)

### Community 250 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 252 - "New"
Cohesion: 0.10
Nodes (25): newScratchStore(), TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(), TestRetrieveUsesCacheAndSyncInvalidates(), TestUpsertBackupSchedulePostgresParity(), TestComplianceFindingRuleDedupAndResolve(), TestComplianceRuleCRUD(), Store, newConcurrentQualityTestStore() (+17 more)

### Community 254 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.29
Nodes (4): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision

### Community 255 - "Permissions & Step-Up for Mutating Actions"
Cohesion: 0.33
Nodes (7): Connector Management via UI (full CRUD), Destructive-Action Pattern: Confirm + Blast Radius, Permissions & Step-Up for Mutating Actions, Role Model — viewer/operator, POST /auth/elevate (step-up elevation), GET /connectors/{id}/removal-impact (blast radius), Role Schema (viewer/operator)

### Community 256 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 259 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 261 - "log/slog.Logger"
Cohesion: 0.11
Nodes (14): formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+6 more)

### Community 262 - "Connector Interface (Name/Fetch/Validate)"
Cohesion: 0.60
Nodes (5): Connector Guide (docs/connectors/CONNECTOR_GUIDE.md), Connector Interface (Name/Fetch/Validate), ServiceSnapshot Data Structure, Connector / ConnectorCreate / ConnectorUpdate / ServiceSnapshot Schemas, Supported Service Connectors (Proxmox, Docker/Portainer, pfSense/OPNsense)

### Community 263 - "0001 — Lab-mutating operation boundaries"
Cohesion: 0.40
Nodes (4): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Consequences, Context

### Community 265 - "BACKUP.md"
Cohesion: 0.14
Nodes (10): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included, Behavior, Configuration, Failure notifications (+2 more)

## Ambiguous Edges - Review These
- `Topbar Notification Center (deferred from V1)` → `NotificationDelivery Schema (per-channel delivery/retry)`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to
- `Topbar Notification Center (deferred from V1)` → `Notification / NotificationPage Schemas`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to

## Knowledge Gaps
- **610 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+605 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1374 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `NotificationDelivery Schema (per-channel delivery/retry)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `Notification / NotificationPage Schemas`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `Store` connect `Store` to `net/http.Request`, `NewStore`, `backup/backup.go`, `compliance/handlers.go`, `Store`, `diagnostics/diagnostics.go`, `NewEngine`, `time.Time`, `DecodeKey`, `Deps`, `Dispatcher`, `testing.T`, `rewritePlaceholders`, `Handler`, `HashPassword`, `chat/chat.go`, `NewRegistry`, `reports/handlers.go`, `go_pkg_context`, `Engine`, `RunMigrations`, `rowScanner`, `Handler`, `backup/backup_test.go`, `ErrorWithDetails`, `ExportToFile`, `alerts/handlers_authz_test.go`, `dispatcher_test.go`, `Get`, `Manager`, `NewUser`, `export_test.go`, `Config`, `main`, `NewEngine`, `ConnectorRecord`, `Checker`, `retention/retention_test.go`, `Sanitize`, `createUser`, `go_pkg_os`, `VerifyBundleFile`, `gitFixture`, `manifest.go`, `notifications/handlers_test.go`, `New`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Why does `WiseLabz WebSocket Contract (`/ws`)` connect `WiseLabz WebSocket Contract (`/ws`)` to `App.tsx`, `0004 — PostgreSQL leader election for background workers`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `net/http.Request` to `routerDeps`, `Store`, `reports/handlers.go`, `Sanitize`, `compliance/handlers.go`, `Handler`, `DecodeJSON`, `Handler`, `DecodeKey`, `logging_test.go`, `ErrorWithDetails`, `response.go`, `Handler`, `.OIDCCallback`, `Get`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _610 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Client` be split into smaller, more focused modules?**
  _Cohesion score 0.03412809724170173 - nodes in this community are weakly interconnected._