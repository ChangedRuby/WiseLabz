# Graph Report - WiseLabz  (2026-09-27)

## Corpus Check
- 868 files · ~536,619 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 20 file(s) not represented in the graph (top: (none) 9, .toml 2, .tmpl 2)

## Summary
- 6587 nodes · 20676 edges · 246 communities (213 shown, 33 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1786 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `dbd9ceac`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- connector/connector.go
- context.Context
- Errorf
- ServiceSnapshot
- SuggestRequest
- Store
- Config
- sync.Mutex
- Handler
- middleware.go
- diagnostics/diagnostics.go
- NewEngine
- newTestApp
- compliance/engine.go
- time.Time
- newDocTestStore
- DecodeKey
- Dispatcher
- WiseLabz WebSocket Contract (`/ws`)
- NewRegistry
- testing.T
- CommandPalette
- initial database schema
- rewritePlaceholders
- Hub
- User
- cn
- Service
- go_pkg_testing
- icons.tsx
- Product
- @tanstack/react-query
- newTestHandler
- RunSync
- chat/chat.go
- Engine
- Application root
- go_pkg_context
- DashboardPage.tsx
- Bottom-dock shell
- RunMigrations
- ProfilePage.tsx
- WiseLabz
- Button.tsx
- ConnectorRecord
- UsersPage.tsx
- Runner
- router.go
- App.tsx
- package.json
- net/http.ResponseWriter
- HashToken
- home_assistant/tables.go
- RunbookRecord
- ServicesPage.tsx
- fixtures.ts
- dialSSHStdio
- dispatcher_test.go
- render_test.go
- lefthook Commit Hooks
- TemplateEditorPage.tsx
- ErrorWithDetails
- .ConfigPush
- react
- routerDeps
- validate.go
- truenas/tables.go
- nilToStr
- AppShell — Bottom Dock Shell (single variant)
- NewUser
- SnapshotEntity
- Connector
- ServiceDetailPage.tsx
- ConnectorEditPage.tsx
- newTestHandler
- dependencies
- NewMalformedResponseError
- ReportsPage.tsx
- NewStore
- Get
- Register
- Diff viewer
- Destructive connector confirmation
- single-instance deployment model
- Graphify Knowledge Graph Rules
- Topbar Notification Center (deferred from V1)
- Database: SQLite + PostgreSQL
- React + Vite Frontend
- portainer/tables.go
- NewChecker
- home_assistant_test.go
- adguardhome/tables.go
- AuthedUser
- traefik/tables.go
- net/http.Client
- truenas_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- main
- Compare
- settings.mock.ts
- handlers.ts
- testhelpers_test.go
- Store
- timeline.ts
- newTestHarness
- channels.go
- go_pkg_os
- unifi_test.go
- SystemPage.tsx
- Manager
- src/theme.ts
- ExportToFile
- WiseLabz — Design Contract
- devDependencies
- response.go
- Checker
- portainer_test.go
- git_test.go
- NotificationRecord
- go_pkg_github_com_wiselabz_wiselabz_internal_store
- middleware_test.go
- handlers_contract_test.go
- createTestConnector
- fetch_test.go
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
- New
- connector_permission.go
- NewService
- .Fetch
- config_test.go
- Connector
- AppearancePage.tsx
- httpx/retry_test.go
- Deps
- logging.go
- net/http.Request
- docker_test.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- AuthMiddleware
- time.Duration
- WiseLabz Connector Guide
- docdiffmodel.ts
- apikey_scope.go
- all.go
- NewClient
- templates_test.go
- compilerOptions
- testApp
- changes/handlers_test.go
- newTestHandler
- Handler
- vectorCache
- api/auth/oidc.go
- connector_permission_test.go
- config_cmd_test.go
- api/changes_test.go
- NewRouter
- connectors_health_test.go
- Contributor Covenant Code of Conduct
- Handler
- newTestHandler
- Store
- Connector
- Decision
- Decision
- scripts
- New
- lifecycleManager
- maintenance_test.go
- Decision
- .call
- walkCursorPages
- TestComplianceRuleValidation
- handlers_bulk_test.go
- .call
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- connectors_hardening_test.go
- ReportData
- release-please-config.json
- RateLimit
- Handler
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
- mustCreateUser
- keyset_test.go
- openapi_contract_test.go
- store/mfa_test.go
- transform.go
- scanMaintenanceWindow
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- connectors_maintenance_test.go
- Mermaid.tsx
- main.tsx
- Security Policy
- seedQualityConnector
- Notification Channels
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- log/slog.Logger
- Saved Views
- WiseLabz — v2 Backlog
- tsconfig.json
- setup-env.sh
- CHANGE_PROVENANCE.md
- ComplianceRuleRecord
- ComputeWindow
- vite-env.d.ts
- github.com/WiseLabz/wiselabz
- notification_delivery_test.go
- change_pattern_test.go
- RunDocLockSweep
- golden_snapshot_test.go
- job_health_test.go
- store/oidc_test.go
- internal/auth/oidc.go
- truenas/attributes_test.go
- IsTimeout

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 229 edges
2. `Errorf()` - 184 edges
3. `newDocTestStore()` - 143 edges
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

## Communities (246 total, 33 thin omitted)

### Community 0 - "connector/connector.go"
Cohesion: 0.03
Nodes (34): Connector, TimeoutError, NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL(), tryParseEntities() (+26 more)

### Community 1 - "context.Context"
Cohesion: 0.02
Nodes (36): fakeStatusChecker, sanitizeSessions(), Connector, existingIDs(), Store, SnapshotRecord, Store, Store (+28 more)

### Community 2 - "Errorf"
Cohesion: 0.06
Nodes (28): Handler, newToken(), sanitize(), Handler, diffToSpec(), NewHandler(), Handler, Handler (+20 more)

### Community 3 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (19): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Sanitize(), changePatternID(), Engine, markError(), snapshotIDOrNil() (+11 more)

### Community 4 - "SuggestRequest"
Cohesion: 0.13
Nodes (11): claudeProvider, openAICompatibleProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList(), TestStubProviderName() (+3 more)

### Community 5 - "Store"
Cohesion: 0.19
Nodes (24): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+16 more)

### Community 6 - "Config"
Cohesion: 0.14
Nodes (20): NewHandler(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+12 more)

### Community 7 - "sync.Mutex"
Cohesion: 0.12
Nodes (7): database/sql.Conn, sync.Mutex, Elector, Noop, fakeDocRegenerator, fakeNotifier, fakeQualityChecker

### Community 8 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 9 - "middleware.go"
Cohesion: 0.15
Nodes (18): AuditRecorder, ConnectorRoleChecker, contextKey, elevationError, PermissionChecker, TreatAsSafeMethod(), elevationFailureReason(), MFAEnrollOnlyFromContext() (+10 more)

### Community 10 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 11 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 12 - "newTestApp"
Cohesion: 0.02
Nodes (178): runbookResp, runbookStepResp, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation() (+170 more)

### Community 13 - "compliance/engine.go"
Cohesion: 0.11
Nodes (30): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+22 more)

### Community 14 - "time.Time"
Cohesion: 0.10
Nodes (25): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session, io.WriteCloser (+17 more)

### Community 15 - "newDocTestStore"
Cohesion: 0.06
Nodes (47): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+39 more)

### Community 16 - "DecodeKey"
Cohesion: 0.05
Nodes (46): ProviderConfig, factorJSON(), Handler, TestDiagnosticsRedactsSecrets(), Handler, primaryProviderConfig(), Handler, GenerateTOTPSecret() (+38 more)

### Community 17 - "Dispatcher"
Cohesion: 0.13
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 18 - "WiseLabz WebSocket Contract (`/ws`)"
Cohesion: 0.06
Nodes (33): ADR 0001 — Monorepo, ADR Index (docs/adr/), AI Doc Generation Module (opt-in, provider-agnostic), API Design — REST + WebSocket split, Dual Auth Design (Local JWT + OIDC), Changes/Diff Contract (infra vs doc format), Change-Aware Diff Engine, Monorepo with Go Workspaces (+25 more)

### Community 19 - "NewRegistry"
Cohesion: 0.12
Nodes (26): Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail() (+18 more)

### Community 20 - "testing.T"
Cohesion: 0.02
Nodes (163): newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestStandbyIsUnreadyAndRunsNoScheduler(), TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks() (+155 more)

### Community 21 - "CommandPalette"
Cohesion: 0.12
Nodes (20): Axios API client, Button and IconButton, CommandPalette, theme cycling command, ConfirmDialog, Dialog, ElevationConfirm, English translation catalog (+12 more)

### Community 22 - "initial database schema"
Cohesion: 0.16
Nodes (20): alerts, changes, connector config JSON, connectors, dashboard layouts, doc versions, docs, HashToken (+12 more)

### Community 23 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 24 - "Hub"
Cohesion: 0.14
Nodes (7): Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator, ticket

### Community 25 - "User"
Cohesion: 0.09
Nodes (15): webAuthnUser, Handler, Handler, newOIDCUser(), randomOIDCToken(), validHostPort(), instanceAdminRoleFor(), OIDCClaims (+7 more)

### Community 26 - "cn"
Cohesion: 0.03
Nodes (127): react-i18next, sonner, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations (+119 more)

### Community 27 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 28 - "go_pkg_testing"
Cohesion: 0.06
Nodes (26): TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), go_pkg_bytes, go_pkg_crypto_hmac, go_pkg_crypto_rsa (+18 more)

### Community 29 - "icons.tsx"
Cohesion: 0.04
Nodes (87): Live dashboard state, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest (+79 more)

### Community 30 - "Product"
Cohesion: 0.09
Nodes (24): Connector Guide (docs/connectors/CONNECTOR_GUIDE.md), Connector Interface (Name/Fetch/Validate), Connector Management via UI (full CRUD), Destructive-Action Pattern: Confirm + Blast Radius, Manager Actions (v1 scope), Permissions & Step-Up for Mutating Actions, Role Model — viewer/operator, ServiceSnapshot Data Structure (+16 more)

### Community 31 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (51): i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage, Command (+43 more)

### Community 32 - "newTestHandler"
Cohesion: 0.05
Nodes (78): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+70 more)

### Community 33 - "RunSync"
Cohesion: 0.15
Nodes (14): Connector interface, connector schema registration, reverse proxy WebSocket support, OpenAPI REST contract, destructive-action step-up authentication, operational alerts, detected changes, Compare (+6 more)

### Community 34 - "chat/chat.go"
Cohesion: 0.20
Nodes (13): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest() (+5 more)

### Community 35 - "Engine"
Cohesion: 0.14
Nodes (7): Config, NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 36 - "Application root"
Cohesion: 0.18
Nodes (12): OpenAPI client generation, Generated-code lint exclusions, Motion preference provider, Vite API and WebSocket proxy, Authentication and onboarding guards, Operator-only routes, Root(), router (+4 more)

### Community 37 - "go_pkg_context"
Cohesion: 0.07
Nodes (18): dashboardLayout, contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_go_webauthn_webauthn_protocol (+10 more)

### Community 38 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (85): Connector category icon map, Dashboard widget frame, 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status` (+77 more)

### Community 39 - "Bottom-dock shell"
Cohesion: 0.40
Nodes (5): Authenticated app frame, Bottom-dock shell, Primary navigation, Non-React navigation bridge, Floating dock navigation

### Community 40 - "RunMigrations"
Cohesion: 0.09
Nodes (40): main(), OpenDB(), newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns() (+32 more)

### Community 41 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (51): web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin (+43 more)

### Community 42 - "WiseLabz"
Cohesion: 0.20
Nodes (10): WCAG 2.2 AA accessibility, Docs-first information architecture, technical homelabbers, machine-honest interface, v1 narrow manager scope, trustworthy live documentation, WiseLabz, commit quality gates (+2 more)

### Community 43 - "Button.tsx"
Cohesion: 0.05
Nodes (58): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze, web_src_api_generated_alerts_alerts_postalertsbulksnooze (+50 more)

### Community 44 - "ConnectorRecord"
Cohesion: 0.06
Nodes (35): enableFakeEmbedding(), actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), ConnectorRecord, Store (+27 more)

### Community 45 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (46): axios, customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers (+38 more)

### Community 46 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 47 - "router.go"
Cohesion: 0.08
Nodes (35): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts (+27 more)

### Community 48 - "App.tsx"
Cohesion: 0.04
Nodes (58): react-router-dom, AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken() (+50 more)

### Community 49 - "package.json"
Cohesion: 0.04
Nodes (49): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+41 more)

### Community 50 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (19): webAuthnFlow, Handler, Handler, Handler, Handler, Handler, Handler, Handler (+11 more)

### Community 51 - "HashToken"
Cohesion: 0.08
Nodes (26): updateUserRequest, Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), HashPassword() (+18 more)

### Community 52 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 53 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 54 - "ServicesPage.tsx"
Cohesion: 0.06
Nodes (38): match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow, web_src_api_generated_connectors_connectors_getgetconnectorsmaintenancewindowsquerykey, web_src_api_generated_connectors_connectors_postconnectorsbulkreauth, web_src_api_generated_connectors_connectors_postconnectorsbulkrestart (+30 more)

### Community 55 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 56 - "dialSSHStdio"
Cohesion: 0.15
Nodes (11): serveOneHTTPExchange(), serveSSHDockerConn(), TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ClientConfig (+3 more)

### Community 57 - "dispatcher_test.go"
Cohesion: 0.20
Nodes (46): NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore(), setChannelAndRoutingConfig(), setChannelConfig(), setChannelConfigJSON() (+38 more)

### Community 58 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 59 - "lefthook Commit Hooks"
Cohesion: 0.50
Nodes (5): commit-msg Hook, Conventional Commits Policy, lefthook Commit Hooks, pre-commit Hook, Commit Conventions & Hook Enforcement (dev workflow)

### Community 60 - "TemplateEditorPage.tsx"
Cohesion: 0.04
Nodes (56): web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid (+48 more)

### Community 61 - "ErrorWithDetails"
Cohesion: 0.06
Nodes (24): validateConfigPushRequest(), applyConnectorScalarUpdates(), configRequestField(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection() (+16 more)

### Community 62 - ".ConfigPush"
Cohesion: 0.27
Nodes (5): configPushLanded(), Handler, isWritableField(), ConfigPusher, Err()

### Community 63 - "react"
Cohesion: 0.06
Nodes (39): 3. `sync.complete`, react, web_src_api_generated_changes_changes, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_getgetchangesquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss (+31 more)

### Community 64 - "routerDeps"
Cohesion: 0.11
Nodes (31): routerDeps, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+23 more)

### Community 65 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 66 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 67 - "nilToStr"
Cohesion: 0.06
Nodes (14): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+6 more)

### Community 68 - "AppShell — Bottom Dock Shell (single variant)"
Cohesion: 0.50
Nodes (4): AppShell — Bottom Dock Shell (single variant), Theme Engine — Code Default, User-Overridable, Per-User Dashboard Layout with Admin Default (v2), DashboardLayout Schema (per-user widget layout)

### Community 69 - "NewUser"
Cohesion: 0.12
Nodes (53): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestGetRequiresGrantEvenForInstanceAdmin(), TestListFiltersGrantsBeforePagination(), Handler, snapshotFixture(), snapshotRequest() (+45 more)

### Community 70 - "SnapshotEntity"
Cohesion: 0.11
Nodes (19): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), ServiceDependency, SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable() (+11 more)

### Community 71 - "Connector"
Cohesion: 0.08
Nodes (10): init(), Connector, buildRouteTable(), buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment(), ValidateRefSegment() (+2 more)

### Community 72 - "ServiceDetailPage.tsx"
Cohesion: 0.05
Nodes (47): ADR-0001, ADR-0003, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridconfigfields (+39 more)

### Community 73 - "ConnectorEditPage.tsx"
Cohesion: 0.08
Nodes (22): RFC-3339, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsschema, web_src_api_model_index_connector (+14 more)

### Community 74 - "newTestHandler"
Cohesion: 0.11
Nodes (38): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+30 more)

### Community 75 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 76 - "NewMalformedResponseError"
Cohesion: 0.07
Nodes (43): SnapshotSection, NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+35 more)

### Community 77 - "ReportsPage.tsx"
Cohesion: 0.11
Nodes (19): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports (+11 more)

### Community 78 - "NewStore"
Cohesion: 0.09
Nodes (43): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+35 more)

### Community 79 - "Get"
Cohesion: 0.06
Nodes (54): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), supportedLifecycleVerbs(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion() (+46 more)

### Community 80 - "Register"
Cohesion: 0.06
Nodes (62): init(), newConnector(), Connector, RequestedFields(), init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback() (+54 more)

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
Cohesion: 0.11
Nodes (34): jsonType(), TestAttributeCatalogCoversEmittedKeys(), environmentDependencies(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames() (+26 more)

### Community 89 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 90 - "home_assistant_test.go"
Cohesion: 0.14
Nodes (28): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+20 more)

### Community 91 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 92 - "AuthedUser"
Cohesion: 0.12
Nodes (29): NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth() (+21 more)

### Community 93 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 94 - "net/http.Client"
Cohesion: 0.07
Nodes (14): ollamaEmbedder, openAIEmbedder, Connector, LimitedBody(), Connector, apiMessage(), controllerName(), countByKind() (+6 more)

### Community 95 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 96 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 97 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.05
Nodes (33): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior (+25 more)

### Community 98 - "main"
Cohesion: 0.15
Nodes (16): main(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry (+8 more)

### Community 99 - "Compare"
Cohesion: 0.08
Nodes (41): driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare(), DiffResult (+33 more)

### Community 100 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 101 - "handlers.ts"
Cohesion: 0.07
Nodes (29): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+21 more)

### Community 102 - "testhelpers_test.go"
Cohesion: 0.31
Nodes (6): GenerateRecoveryCodes(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestNormalizeRecoveryCode(), go_pkg_github_com_pquerna_otp, go_pkg_github_com_pquerna_otp_totp

### Community 103 - "Store"
Cohesion: 0.06
Nodes (15): Store, BackupSchedule, Store, placeholders(), scanBackupRun(), changeFilterClause(), AlertRecord, ChangeRecord (+7 more)

### Community 104 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 105 - "newTestHarness"
Cohesion: 0.22
Nodes (16): seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), TestSearchDocs(), seedFinding(), TestListFindings() (+8 more)

### Community 106 - "channels.go"
Cohesion: 0.18
Nodes (16): buildEmailMessage(), sendNtfyChannel(), sendSMTPChannel(), sendTelegramChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendNtfyChannel_DefaultsToNtfySh(), TestSendNtfyChannel_MissingTopic() (+8 more)

### Community 107 - "go_pkg_os"
Cohesion: 0.05
Nodes (34): IsGeneratedName(), pruneStale(), TestIsGeneratedName(), TestCommitMessage(), writeExportState(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults() (+26 more)

### Community 108 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 109 - "SystemPage.tsx"
Cohesion: 0.06
Nodes (33): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey, web_src_api_generated_system_system_getsystembackupschedule (+25 more)

### Community 110 - "Manager"
Cohesion: 0.15
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 111 - "src/theme.ts"
Cohesion: 0.15
Nodes (23): @fontsource/space-mono, @fontsource-variable/space-grotesk, ColorMode, commit(), load(), Persisted, PRESETS_FONTS, ThemeState (+15 more)

### Community 112 - "ExportToFile"
Cohesion: 0.10
Nodes (42): Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+34 more)

### Community 113 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 114 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 115 - "response.go"
Cohesion: 0.08
Nodes (29): Handler, decodeStoredSnapshot(), Handler, snapshotStoreError(), Cursor(), DecodeCursor(), EncodeCursor(), T (+21 more)

### Community 116 - "Checker"
Cohesion: 0.11
Nodes (14): Snapshot, complianceRule(), Checker, RunStaleSweepOnce(), docSearchWhere(), escapeLike(), DocRecord, Store (+6 more)

### Community 117 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 118 - "git_test.go"
Cohesion: 0.06
Nodes (47): fetchAllDocs(), fileName(), Exporter, NewExporter(), RunExportOnce(), slugify(), newTestStore(), readFile() (+39 more)

### Community 119 - "NotificationRecord"
Cohesion: 0.21
Nodes (6): Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord, Store, scanNotification()

### Community 120 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.05
Nodes (50): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection, bulkResolveItemResult (+42 more)

### Community 121 - "middleware_test.go"
Cohesion: 0.14
Nodes (17): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser(), TestAuthMiddlewareAllowsCurrentRoleClaim() (+9 more)

### Community 122 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 123 - "createTestConnector"
Cohesion: 0.10
Nodes (27): TestSnapshotKeysetSummaryAndConnectorOwnership(), skipOnPostgres(), TestTemplateVersionIndexAndCascade(), TestDeleteOldHealthChecks(), TestGetConnectorUptimeDeterministicOutage(), TestGetConnectorUptimeNoData(), TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR(), TestRecordHealthCheckDefaults() (+19 more)

### Community 124 - "fetch_test.go"
Cohesion: 0.11
Nodes (23): extractGroups(), newMockOIDCServer(), TestAuthURLAfterInitialization(), TestAuthURLBeforeInitialization(), TestExtractGroups(), TestInitializeFailure(), TestInitializeInvalidJSON(), TestInitializeSuccess() (+15 more)

### Community 136 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 137 - "New"
Cohesion: 0.13
Nodes (24): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression() (+16 more)

### Community 138 - "connector_permission.go"
Cohesion: 0.19
Nodes (9): auditConnectorGrantDiffJSON(), getConnectorGrant(), ConnectorGrantDiff, Store, highestConnectorRole(), listOIDCConnectorGrants(), scanConnectorGrants(), upsertConnectorGrant() (+1 more)

### Community 139 - "NewService"
Cohesion: 0.25
Nodes (14): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+6 more)

### Community 140 - ".Fetch"
Cohesion: 0.15
Nodes (8): WantsField(), Connector, putMetadata(), unavailable(), agentEnabled(), Connector, Connector, dockerSectionSpec

### Community 141 - "config_test.go"
Cohesion: 0.09
Nodes (28): runHealthcheck(), Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields() (+20 more)

### Community 142 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 143 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 144 - "httpx/retry_test.go"
Cohesion: 0.20
Nodes (20): retryable(), RetryTransport(), sleep(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry() (+12 more)

### Community 145 - "Deps"
Cohesion: 0.28
Nodes (17): registerListAttentionItems(), changeServiceIDs(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs() (+9 more)

### Community 146 - "logging.go"
Cohesion: 0.13
Nodes (18): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestLoggablePathMasksShareTokenUnderV1() (+10 more)

### Community 147 - "net/http.Request"
Cohesion: 0.08
Nodes (21): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+13 more)

### Community 148 - "docker_test.go"
Cohesion: 0.07
Nodes (40): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+32 more)

### Community 149 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 150 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 151 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 152 - "AuthMiddleware"
Cohesion: 0.13
Nodes (14): APIKeyChecker, testAPIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), APIKeyClaims (+6 more)

### Community 153 - "time.Duration"
Cohesion: 0.15
Nodes (7): AuthSettings, Database, Server, WebAuthnSettings, time.Duration, OIDCProvider, PoolConfig

### Community 154 - "WiseLabz Connector Guide"
Cohesion: 0.08
Nodes (24): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Conventions (+16 more)

### Community 155 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 156 - "apikey_scope.go"
Cohesion: 0.24
Nodes (10): APIKeyRestriction, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), isSafeMethod(), TestClampConnectorRole(), treatAsSafeFromContext(), IsSafeMethod() (+2 more)

### Community 157 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 158 - "NewClient"
Cohesion: 0.15
Nodes (16): GuardedDialer(), IsDangerousIP(), clientTimeout(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects() (+8 more)

### Community 159 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 160 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 161 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 162 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 163 - "newTestHandler"
Cohesion: 0.29
Nodes (7): Handler, newTestHandler(), TestCreate(), TestExecuteStepNotFound(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete()

### Community 164 - "Handler"
Cohesion: 0.24
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 165 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 166 - "api/auth/oidc.go"
Cohesion: 0.08
Nodes (24): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), dateFormat() (+16 more)

### Community 167 - "connector_permission_test.go"
Cohesion: 0.34
Nodes (14): Store, newTestConnector(), newTestUser(), TestDeleteConnectorGrant(), TestFilterConnectorIDsByGrant(), TestGetUserConnectorRoleHighestAcrossSources(), TestListConnectorGrants(), TestListConnectorIDs() (+6 more)

### Community 168 - "config_cmd_test.go"
Cohesion: 0.21
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 169 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 170 - "NewRouter"
Cohesion: 0.22
Nodes (10): TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+2 more)

### Community 171 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 172 - "Contributor Covenant Code of Conduct"
Cohesion: 0.15
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 173 - "Handler"
Cohesion: 0.18
Nodes (4): cron.EntryID, Handler, sync/atomic.Bool, ReadyState

### Community 174 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 175 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 176 - "Connector"
Cohesion: 0.11
Nodes (3): ConfigField, Connector, Connector

### Community 177 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 178 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 179 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 180 - "New"
Cohesion: 0.09
Nodes (26): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+18 more)

### Community 181 - "lifecycleManager"
Cohesion: 0.16
Nodes (7): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, lifecycleDeps, lifecycleManager

### Community 182 - "maintenance_test.go"
Cohesion: 0.23
Nodes (12): TestListConnectorNames(), TestListDocsGroupedByService(), Store, mustCreateMaintenanceConnector(), TestCloseMaintenanceWindow(), TestCloseMaintenanceWindowAlreadyClosed(), TestCloseMaintenanceWindowNotFound(), TestCreateAndGetActiveMaintenanceWindow() (+4 more)

### Community 183 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 184 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 185 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 186 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 187 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 188 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 189 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 190 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 191 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 192 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 193 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 194 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 195 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 197 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 208 - "mustCreateUser"
Cohesion: 0.24
Nodes (10): Store, mustCreateUser(), TestDocLockAcquireAfterExpiry(), TestDocLockConflict(), TestDocLockReleaseOnlyByHolder(), TestDocLockRenewalByHolder(), TestMarkAllNotificationsRead(), TestMarkNotificationReadIdempotentAndScoped() (+2 more)

### Community 209 - "keyset_test.go"
Cohesion: 0.33
Nodes (10): assertSameSet(), Store, T, queryPlan(), TestKeysetQueriesUseCoveringIndexes(), TestListAuditRecordsKeysetHonoursFilters(), TestListAuditRecordsKeysetTraversal(), TestListChangesKeysetTraversal() (+2 more)

### Community 210 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 211 - "store/mfa_test.go"
Cohesion: 0.31
Nodes (10): Store, mfaTestUser(), TestConfirmFactorRejectsSecondConfirmedTOTP(), TestConsumeTOTPStepReplayGuard(), TestDeleteFactorLastOneAlsoLeavesRecoveryCodesForCallerToWipe(), TestDeleteUserCascadesMFAFactorsAndRecoveryCodes(), TestDeleteUserFactorsForAdminReset(), TestGetRequire2FADefaultsToNone() (+2 more)

### Community 212 - "transform.go"
Cohesion: 0.25
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 213 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 214 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 215 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 216 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 217 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 218 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 219 - "main.tsx"
Cohesion: 0.40
Nodes (4): react-dom, App(), USE_MOCKS, web_src_index

### Community 220 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 221 - "seedQualityConnector"
Cohesion: 0.22
Nodes (9): TestComplianceFindingRuleDedupAndResolve(), TestComplianceRuleCRUD(), Store, newConcurrentQualityTestStore(), seedQualityConnector(), TestListQualityFindingsFilters(), TestResolveThenReopenCreatesFreshRow(), TestUpsertQualityFindingConcurrentDedup() (+1 more)

### Community 222 - "Notification Channels"
Cohesion: 0.40
Nodes (4): Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 223 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 224 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 226 - "log/slog.Logger"
Cohesion: 0.17
Nodes (18): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), newLogger(), WithLogger(), RunCleanupOnce(), newTestStore() (+10 more)

### Community 227 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 232 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 233 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 236 - "notification_delivery_test.go"
Cohesion: 0.39
Nodes (7): createTestNotification(), Store, TestDeliveryCreateAndList(), TestListDeliveriesStatusFilterAndPagination(), TestListDueDeliveries(), TestUpdateDeliveryResultNotFound(), TestUpdateDeliveryResultTransitionsAndClearsNextAttempt()

### Community 237 - "change_pattern_test.go"
Cohesion: 0.48
Nodes (6): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern()

### Community 238 - "RunDocLockSweep"
Cohesion: 0.43
Nodes (4): Store, RunDocLockSweep(), runDocLockSweep(), DocLockRecord

### Community 239 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 240 - "job_health_test.go"
Cohesion: 0.40
Nodes (4): TestJobHealthDelete(), TestJobHealthGetMissingReturnsNoRows(), TestJobHealthListOrderedByName(), TestJobHealthUpsertInsertsThenUpdates()

### Community 242 - "store/oidc_test.go"
Cohesion: 0.50
Nodes (3): TestGetOIDCProviderFlagsErrorHandling(), TestOIDCProviderFlags(), TestSetOIDCProviderEnabledErrorHandling()

## Ambiguous Edges - Review These
- `Topbar Notification Center (deferred from V1)` → `NotificationDelivery Schema (per-channel delivery/retry)`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to
- `Topbar Notification Center (deferred from V1)` → `Notification / NotificationPage Schemas`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to

## Knowledge Gaps
- **583 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+578 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1347 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `NotificationDelivery Schema (per-channel delivery/retry)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `Notification / NotificationPage Schemas`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `Store` connect `Store` to `Errorf`, `ServiceSnapshot`, `Config`, `sync.Mutex`, `Handler`, `diagnostics/diagnostics.go`, `NewEngine`, `time.Time`, `Deps`, `Dispatcher`, `NewRegistry`, `testing.T`, `rewritePlaceholders`, `newTestHandler`, `testApp`, `changes/handlers_test.go`, `Engine`, `chat/chat.go`, `go_pkg_context`, `Handler`, `RunMigrations`, `ConnectorRecord`, `Handler`, `net/http.ResponseWriter`, `HashToken`, `New`, `lifecycleManager`, `.call`, `dispatcher_test.go`, `.call`, `ErrorWithDetails`, `ReportData`, `Handler`, `NewUser`, `NewStore`, `Register`, `NewChecker`, `AuthedUser`, `log/slog.Logger`, `newTestHarness`, `Manager`, `ExportToFile`, `response.go`, `Checker`, `git_test.go`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Why does `Runner` connect `Runner` to `testApp`, `context.Context`, `Engine`, `log/slog.Logger`, `go_pkg_context`, `sync.Mutex`, `New`, `Handler`, `NewStore`, `lifecycleManager`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `@tanstack/react-query` connect `@tanstack/react-query` to `DashboardPage.tsx`, `ServiceDetailPage.tsx`, `ProfilePage.tsx`, `ConnectorEditPage.tsx`, `Button.tsx`, `UsersPage.tsx`, `ReportsPage.tsx`, `SystemPage.tsx`, `App.tsx`, `package.json`, `ServicesPage.tsx`, `cn`, `TemplateEditorPage.tsx`, `icons.tsx`, `react`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _583 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `connector/connector.go` be split into smaller, more focused modules?**
  _Cohesion score 0.03474399164054336 - nodes in this community are weakly interconnected._