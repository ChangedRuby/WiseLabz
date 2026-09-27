# Graph Report - WiseLabz  (2026-09-27)

## Corpus Check
- 869 files · ~537,586 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 20 file(s) not represented in the graph (top: (none) 9, .toml 2, .tmpl 2)

## Summary
- 6599 nodes · 20725 edges · 250 communities (214 shown, 36 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1787 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9168dc80`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Client
- context.Context
- Errorf
- ServiceSnapshot
- provider_test.go
- backup/backup.go
- time.Duration
- Elector
- Handler
- RequireConnectorRole
- diagnostics/diagnostics.go
- Engine
- newTestApp
- runbooks_test.go
- time.Time
- newDocTestStore
- DecodeKey
- Dispatcher
- WiseLabz WebSocket Contract (`/ws`)
- Registry
- testing.T
- CommandPalette
- initial database schema
- rewritePlaceholders
- Hub
- .OIDCCallback
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
- NewEngine
- rowScanner
- UsersPage.tsx
- ErrorWithDetails
- router.go
- App.tsx
- package.json
- DecodeJSON
- HashToken
- home_assistant/tables.go
- RunbookRecord
- react
- fixtures.ts
- newSSHDockerClient
- dispatcher_test.go
- render_test.go
- lefthook Commit Hooks
- TemplateEditorPage.tsx
- net/http.Request
- Connector
- api/auth/oidc.go
- mountAPIRoutes
- config/validate_test.go
- truenas/tables.go
- nilToStr
- AppShell — Bottom Dock Shell (single variant)
- NewUser
- SnapshotEntity
- Connector
- ServiceDetailPage.tsx
- useRole.ts
- newTestHandler
- dependencies
- NewMalformedResponseError
- RulesPage.tsx
- NewStore
- GetTypeSchema
- NewEngine
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
- Connector
- truenas_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- main
- Compare
- settings.mock.ts
- handlers.ts
- NewHTTPClient
- ConnectorRecord
- timeline.ts
- apikey_scope.go
- middleware.go
- go_pkg_os
- unifi_test.go
- SnapshotsPage.tsx
- Manager
- src/theme.ts
- VerifyBundleFile
- WiseLabz — Design Contract
- devDependencies
- response.go
- MarshalConnectorConfig
- portainer_test.go
- export_test.go
- pagination_contract_test.go
- go_pkg_github_com_wiselabz_wiselabz_internal_store
- AuthMiddleware
- handlers_contract_test.go
- Store
- Register
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
- Runner
- connector_permission.go
- NewService
- .Fetch
- config_test.go
- Connector
- ThemeControls.tsx
- httpx/retry_test.go
- Deps
- Logger
- Handler
- docker_test.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- data.go
- newHandler
- WiseLabz Connector Guide
- docdiffmodel.ts
- ExportToFile
- all.go
- connector/connector.go
- manifest.go
- compilerOptions
- testApp
- changes/handlers_test.go
- newTestHandler
- Handler
- vectorCache
- templatefuncs.go
- Handler
- config_cmd_test.go
- api/changes_test.go
- net/http.Handler
- NewRegistry
- Contributor Covenant Code of Conduct
- Store
- ListSchemas
- Store
- Connector
- Decision
- Decision
- scripts
- New
- lifecycleManager
- backup/backup_test.go
- Decision
- .call
- Store
- TestComplianceRuleValidation
- handlers_bulk_test.go
- .call
- Changelog
- mockServiceWorker.js
- api/docs_test.go
- connectors_hardening_test.go
- ReportData
- release-please-config.json
- ratelimit.go
- net/http.ResponseWriter
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
- newDockerClient
- net/http.Response
- openapi_contract_test.go
- dashboard/handlers_test.go
- transform_test.go
- scanMaintenanceWindow
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- connectors_maintenance_test.go
- main.tsx
- ShareLink
- Security Policy
- RequirePermission
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
- webAuthnUser
- snapshotResponse
- .applyChannelSecrets
- golden_snapshot_test.go
- computeNextRun
- .GetConnectorUptime
- Store
- internal/auth/oidc.go
- engine_maintenance_test.go
- responseWriter
- seedScopeFixture
- RetentionSettings
- fields_test.go
- stubEmbedder

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

## Communities (250 total, 36 thin omitted)

### Community 0 - "net/http.Client"
Cohesion: 0.03
Nodes (35): Connector, ollamaEmbedder, openAIEmbedder, NewAuthError(), NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities() (+27 more)

### Community 1 - "context.Context"
Cohesion: 0.03
Nodes (34): sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, changeServiceIDs(), Store, existingIDs(), placeholders(), ChangeRecord (+26 more)

### Community 2 - "Errorf"
Cohesion: 0.04
Nodes (33): Handler, newToken(), sanitize(), Handler, Handler, Handler, Handler, decodeStoredSnapshot() (+25 more)

### Community 3 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (16): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, runTransformers(), TestRunTransformersAppliesInOrderAndStopsOnError(), TestRunTransformersUnknownCategoryIsNoop(), registryTestRefresher, actionConnector (+8 more)

### Community 4 - "provider_test.go"
Cohesion: 0.29
Nodes (6): testProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream()

### Community 5 - "backup/backup.go"
Cohesion: 0.19
Nodes (23): connectorIDs(), docIDs(), exportDocs(), exportWithin(), AIConfigSummary, Import(), importBundle(), importConnectors() (+15 more)

### Community 6 - "time.Duration"
Cohesion: 0.08
Nodes (27): newLogger(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+19 more)

### Community 7 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 8 - "Handler"
Cohesion: 0.25
Nodes (6): response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 9 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 10 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 11 - "Engine"
Cohesion: 0.17
Nodes (15): Engine, dedupKey(), matchEntities(), matchReason(), seedEngineConnectorWithEntities(), TestGenerateLabTopologyCreatesThenUpdatesInPlace(), TestMatchEntitiesDedupesExactExternalIDDuplicates(), TestMatchEntitiesExternalIDPrecedence() (+7 more)

### Community 12 - "newTestApp"
Cohesion: 0.02
Nodes (178): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+170 more)

### Community 13 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 14 - "time.Time"
Cohesion: 0.11
Nodes (13): closeQuietly(), digestDue(), formatDigest(), Dispatcher, TestDigestDue(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session (+5 more)

### Community 15 - "newDocTestStore"
Cohesion: 0.02
Nodes (160): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+152 more)

### Community 16 - "DecodeKey"
Cohesion: 0.07
Nodes (33): factorJSON(), Handler, testHandler, Handler, Handler, Handler, GenerateRecoveryCodes(), GenerateTOTPSecret() (+25 more)

### Community 17 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 18 - "WiseLabz WebSocket Contract (`/ws`)"
Cohesion: 0.06
Nodes (33): ADR 0001 — Monorepo, ADR Index (docs/adr/), AI Doc Generation Module (opt-in, provider-agnostic), API Design — REST + WebSocket split, Dual Auth Design (Local JWT + OIDC), Changes/Diff Contract (infra vs doc format), Change-Aware Diff Engine, Monorepo with Go Workspaces (+25 more)

### Community 19 - "Registry"
Cohesion: 0.08
Nodes (23): claudeProvider, openAICompatibleProvider, Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable() (+15 more)

### Community 20 - "testing.T"
Cohesion: 0.02
Nodes (156): cursorPage, newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestStandbyIsUnreadyAndRunsNoScheduler(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL() (+148 more)

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
Cohesion: 0.06
Nodes (22): Handler, isWritableField(), validateConfigPushRequest(), decodeBulkRequest(), Handler, loggablePath(), loggableQuery(), ConfigPusher (+14 more)

### Community 25 - ".OIDCCallback"
Cohesion: 0.16
Nodes (9): auditConnectorGrantDiffJSON(), Handler, newOIDCUser(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider, github.com/coreos/go-oidc/v3/oidc.Provider (+1 more)

### Community 26 - "cn"
Cohesion: 0.03
Nodes (99): web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_notifications_notifications, web_src_api_generated_notifications_notifications_getgetnotificationsquerykey, web_src_api_generated_notifications_notifications_postnotificationsnotificationidread, web_src_api_generated_notifications_notifications_postnotificationsreadall, web_src_api_generated_notifications_notifications_usegetnotifications, web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings (+91 more)

### Community 27 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 28 - "go_pkg_testing"
Cohesion: 0.06
Nodes (20): dashboardLayout, TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestBuildHostOverrideTableAttributes(), jsonType(), TestAttributeCatalogCoversEmittedKeys() (+12 more)

### Community 29 - "icons.tsx"
Cohesion: 0.04
Nodes (82): Live dashboard state, match-sorter, react-error-boundary, sonner, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey (+74 more)

### Community 30 - "Product"
Cohesion: 0.09
Nodes (24): Connector Guide (docs/connectors/CONNECTOR_GUIDE.md), Connector Interface (Name/Fetch/Validate), Connector Management via UI (full CRUD), Destructive-Action Pattern: Confirm + Blast Radius, Manager Actions (v1 scope), Permissions & Step-Up for Mutating Actions, Role Model — viewer/operator, ServiceSnapshot Data Structure (+16 more)

### Community 31 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (57): 7. `quality.finding.created` and `quality.findings.changed`, i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage (+49 more)

### Community 32 - "newTestHandler"
Cohesion: 0.06
Nodes (74): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+66 more)

### Community 33 - "RunSync"
Cohesion: 0.15
Nodes (14): Connector interface, connector schema registration, reverse proxy WebSocket support, OpenAPI REST contract, destructive-action step-up authentication, operational alerts, detected changes, Compare (+6 more)

### Community 34 - "chat/chat.go"
Cohesion: 0.19
Nodes (12): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+4 more)

### Community 35 - "Engine"
Cohesion: 0.25
Nodes (3): Engine, sync.Map, DocRegenerator

### Community 36 - "Application root"
Cohesion: 0.18
Nodes (12): OpenAPI client generation, Generated-code lint exclusions, Motion preference provider, Vite API and WebSocket proxy, Authentication and onboarding guards, Operator-only routes, Root(), router (+4 more)

### Community 37 - "go_pkg_context"
Cohesion: 0.08
Nodes (16): contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_go_webauthn_webauthn_protocol, go_pkg_github_com_go_webauthn_webauthn_webauthn (+8 more)

### Community 38 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (105): Connector category icon map, Dashboard widget frame, 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status` (+97 more)

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

### Community 43 - "NewEngine"
Cohesion: 0.21
Nodes (26): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+18 more)

### Community 44 - "rowScanner"
Cohesion: 0.05
Nodes (33): Dispatcher, Dispatcher, RunDeliveryRetries(), actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows() (+25 more)

### Community 45 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (57): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+49 more)

### Community 46 - "ErrorWithDetails"
Cohesion: 0.14
Nodes (15): oidcElevateFlow, Handler, readOIDCElevateFlowCookie(), randomOIDCToken(), parseScheduleUpdates(), validateRotationFields(), stepAuditDetail(), validTargetType() (+7 more)

### Community 47 - "router.go"
Cohesion: 0.10
Nodes (25): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+17 more)

### Community 48 - "App.tsx"
Cohesion: 0.04
Nodes (67): react-router-dom, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthloginmfa, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh (+59 more)

### Community 49 - "package.json"
Cohesion: 0.04
Nodes (49): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+41 more)

### Community 50 - "DecodeJSON"
Cohesion: 0.18
Nodes (8): webAuthnFlow, Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T, github.com/go-webauthn/webauthn/webauthn.SessionData

### Community 51 - "HashToken"
Cohesion: 0.10
Nodes (23): sanitizeUser(), setRefreshCookie(), Handler, mustHashDummyPassword(), instanceAdminRoleFor(), HashPassword(), TestHashPasswordRejectsOver72Bytes(), TestHashAndVerify() (+15 more)

### Community 52 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 53 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 54 - "react"
Cohesion: 0.04
Nodes (106): RFC-3339, Frontend, motion, @radix-ui/react-popover, react, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_postalertsalertiddismiss (+98 more)

### Community 55 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 56 - "newSSHDockerClient"
Cohesion: 0.19
Nodes (14): generateSSHHostKey(), startSSHDockerServer(), TestDialSSHStdioHonorsContextCancel(), TestNewSSHDockerClientDialsAndExecutesDialStdio(), TestNewSSHDockerClientRejectsMissingHostKey(), TestNewSSHDockerClientRejectsWrongCredentials(), TestNewSSHDockerClientRejectsWrongHostKey(), TestNewSSHDockerClientSupportsSequentialRequests() (+6 more)

### Community 57 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 58 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 59 - "lefthook Commit Hooks"
Cohesion: 0.50
Nodes (5): commit-msg Hook, Conventional Commits Policy, lefthook Commit Hooks, pre-commit Hook, Commit Conventions & Hook Enforcement (dev workflow)

### Community 60 - "TemplateEditorPage.tsx"
Cohesion: 0.06
Nodes (39): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid (+31 more)

### Community 61 - "net/http.Request"
Cohesion: 0.08
Nodes (25): applyConnectorScalarUpdates(), configRequestField(), Handler, validateConnectorConfig(), writeConfigRejection(), capitalize(), Handler, Handler (+17 more)

### Community 62 - "Connector"
Cohesion: 0.21
Nodes (4): SnapshotSection, Connector, unavailable(), session

### Community 63 - "api/auth/oidc.go"
Cohesion: 0.12
Nodes (14): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), shareLinkContextKey (+6 more)

### Community 64 - "mountAPIRoutes"
Cohesion: 0.10
Nodes (26): AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+18 more)

### Community 65 - "config/validate_test.go"
Cohesion: 0.17
Nodes (13): Config, mask(), redactDSN(), redactKVPassword(), Config, TestEveryKeyEnvOverridable(), TestRedactDSN(), TestRedacted() (+5 more)

### Community 66 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 67 - "nilToStr"
Cohesion: 0.09
Nodes (13): exportTemplates(), nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery() (+5 more)

### Community 68 - "AppShell — Bottom Dock Shell (single variant)"
Cohesion: 0.50
Nodes (4): AppShell — Bottom Dock Shell (single variant), Theme Engine — Code Default, User-Overridable, Per-User Dashboard Layout with Admin Default (v2), DashboardLayout Schema (per-user widget layout)

### Community 69 - "NewUser"
Cohesion: 0.18
Nodes (40): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld(), TestGetRootIsSynthetic() (+32 more)

### Community 70 - "SnapshotEntity"
Cohesion: 0.12
Nodes (26): ServiceDependency, SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildInterfaceTableAttributes(), buildInterfaceTable(), buildHostsTable(), parseHosts() (+18 more)

### Community 71 - "Connector"
Cohesion: 0.05
Nodes (18): init(), ConfigField, Connector, TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable() (+10 more)

### Community 72 - "ServiceDetailPage.tsx"
Cohesion: 0.04
Nodes (64): ADR-0001, ADR-0003, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart (+56 more)

### Community 73 - "useRole.ts"
Cohesion: 0.13
Nodes (18): web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, RoleGate(), RoleGateProps, ConnectorEditPage(), daysAgo(), DocHistory(), LinkedDocPanel() (+10 more)

### Community 74 - "newTestHandler"
Cohesion: 0.11
Nodes (41): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+33 more)

### Community 75 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 76 - "NewMalformedResponseError"
Cohesion: 0.12
Nodes (36): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+28 more)

### Community 77 - "RulesPage.tsx"
Cohesion: 0.05
Nodes (45): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+37 more)

### Community 78 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 79 - "GetTypeSchema"
Cohesion: 0.10
Nodes (30): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+22 more)

### Community 80 - "NewEngine"
Cohesion: 0.14
Nodes (26): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+18 more)

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
Cohesion: 0.06
Nodes (72): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+64 more)

### Community 90 - "home_assistant_test.go"
Cohesion: 0.09
Nodes (44): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+36 more)

### Community 91 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 92 - "AuthedUser"
Cohesion: 0.20
Nodes (16): TestCreate(), AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty() (+8 more)

### Community 93 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 94 - "Connector"
Cohesion: 0.16
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

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
Cohesion: 0.11
Nodes (22): main(), runHealthcheck(), splitOrigins(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry, NewEmbedRegistry() (+14 more)

### Community 99 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 100 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 101 - "handlers.ts"
Cohesion: 0.07
Nodes (28): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+20 more)

### Community 102 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 103 - "ConnectorRecord"
Cohesion: 0.09
Nodes (19): AlertRecord, ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), changePatternID() (+11 more)

### Community 104 - "timeline.ts"
Cohesion: 0.14
Nodes (16): enableMocks(), installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env() (+8 more)

### Community 105 - "apikey_scope.go"
Cohesion: 0.12
Nodes (24): APIKeyRestriction, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), isSafeMethod(), TestClampConnectorRole(), treatAsSafeFromContext(), TreatAsSafeMethod() (+16 more)

### Community 106 - "middleware.go"
Cohesion: 0.05
Nodes (40): contextKey, elevationError, RegisterClaude(), TestRegisterClaudeDefaults(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides() (+32 more)

### Community 107 - "go_pkg_os"
Cohesion: 0.07
Nodes (29): keys(), writeExportState(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults(), TestWithinTransactionRollsBack(), exportCursor, exportState (+21 more)

### Community 108 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 109 - "SnapshotsPage.tsx"
Cohesion: 0.12
Nodes (21): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotssnapshotid, web_src_api_model_index_entitychange, web_src_api_model_index_snapshotdiff (+13 more)

### Community 110 - "Manager"
Cohesion: 0.09
Nodes (12): cron.EntryID, Manager, LogPartial(), NewManager(), Store, Store, ReportDefinitionRecord, ReportRecord (+4 more)

### Community 111 - "src/theme.ts"
Cohesion: 0.14
Nodes (25): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+17 more)

### Community 112 - "VerifyBundleFile"
Cohesion: 0.23
Nodes (16): failVerification(), LatestBundle(), RunVerifyOnce(), ListVerifications(), RecordVerification(), seedOneDoc(), TestLatestBundleNoBundles(), TestLatestBundlePicksMostRecent() (+8 more)

### Community 113 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 114 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 115 - "response.go"
Cohesion: 0.12
Nodes (18): DecodeCursor(), EncodeCursor(), TestCursorRequestModes(), TestCursorRoundTrip(), TestDecodeCursorRejectsGarbage(), TestNextCursorStopsOnShortPage(), TestWritePaginatedOmitsNextCursor(), Error() (+10 more)

### Community 116 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 117 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 118 - "export_test.go"
Cohesion: 0.07
Nodes (45): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+37 more)

### Community 119 - "pagination_contract_test.go"
Cohesion: 0.16
Nodes (14): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+6 more)

### Community 120 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.05
Nodes (46): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection, bulkResolveItemResult (+38 more)

### Community 121 - "AuthMiddleware"
Cohesion: 0.09
Nodes (25): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, AuthMiddleware(), extractBearerToken(), hashToken() (+17 more)

### Community 122 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 123 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 124 - "Register"
Cohesion: 0.23
Nodes (16): init(), init(), init(), init(), init(), init(), init(), init() (+8 more)

### Community 136 - "adguardhome_test.go"
Cohesion: 0.10
Nodes (34): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+26 more)

### Community 137 - "Runner"
Cohesion: 0.08
Nodes (30): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner, New() (+22 more)

### Community 138 - "connector_permission.go"
Cohesion: 0.21
Nodes (8): getConnectorGrant(), ConnectorGrantDiff, Store, highestConnectorRole(), listOIDCConnectorGrants(), scanConnectorGrants(), upsertConnectorGrant(), ConnectorGrant

### Community 139 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 140 - ".Fetch"
Cohesion: 0.15
Nodes (8): WantsField(), Connector, putMetadata(), unavailable(), agentEnabled(), Connector, Connector, dockerSectionSpec

### Community 141 - "config_test.go"
Cohesion: 0.15
Nodes (19): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+11 more)

### Community 142 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 143 - "ThemeControls.tsx"
Cohesion: 0.09
Nodes (27): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented() (+19 more)

### Community 144 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 145 - "Deps"
Cohesion: 0.22
Nodes (19): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs(), registerListFindings() (+11 more)

### Community 146 - "Logger"
Cohesion: 0.18
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRateLimit(), TestRecovererPassThrough() (+7 more)

### Community 147 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 148 - "docker_test.go"
Cohesion: 0.09
Nodes (23): generateSelfSignedCert(), serveOneHTTPExchange(), serveSSHDockerConn(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout(), TestDoRequestErrorCases(), TestFetchBuildsSectionsFromEndpoints() (+15 more)

### Community 149 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 150 - "ws.ts"
Cohesion: 0.11
Nodes (18): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+10 more)

### Community 151 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 152 - "data.go"
Cohesion: 0.24
Nodes (15): ChangeEntry, ComplianceSection, ConnectorDrift, DefinitionSummary, DocChangeEntry, DocsSection, DriftSection, FindingSummary (+7 more)

### Community 153 - "newHandler"
Cohesion: 0.22
Nodes (14): TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), Handler, newHandler(), serve() (+6 more)

### Community 154 - "WiseLabz Connector Guide"
Cohesion: 0.08
Nodes (24): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Conventions (+16 more)

### Community 155 - "docdiffmodel.ts"
Cohesion: 0.23
Nodes (13): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, fold(), toUnits(), DiffLine, DiffLineType (+5 more)

### Community 156 - "ExportToFile"
Cohesion: 0.20
Nodes (15): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+7 more)

### Community 157 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 158 - "connector/connector.go"
Cohesion: 0.09
Nodes (13): TimeoutError, GuardedDialer(), IsDangerousIP(), NewTimeoutError(), newWebhookClient(), AuthError, CredentialRefresher, MalformedResponseError (+5 more)

### Community 159 - "manifest.go"
Cohesion: 0.25
Nodes (14): ImportFromFile(), AppVersion(), BuildManifest(), BundleCounts(), ChecksumBytes(), ManifestPath(), ReadManifest(), corruptFile() (+6 more)

### Community 160 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 161 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 162 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 163 - "newTestHandler"
Cohesion: 0.24
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 164 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 165 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 166 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 167 - "Handler"
Cohesion: 0.27
Nodes (4): updateUserRequest, Handler, writeUserWriteError(), NoContent()

### Community 168 - "config_cmd_test.go"
Cohesion: 0.21
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 169 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 170 - "net/http.Handler"
Cohesion: 0.13
Nodes (17): Config, TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders() (+9 more)

### Community 171 - "NewRegistry"
Cohesion: 0.45
Nodes (12): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+4 more)

### Community 172 - "Contributor Covenant Code of Conduct"
Cohesion: 0.15
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 173 - "Store"
Cohesion: 0.07
Nodes (27): routerDeps, Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+19 more)

### Community 174 - "ListSchemas"
Cohesion: 0.19
Nodes (12): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), Capabilities(), CapabilityDescriptor, supportedLifecycleVerbs(), IsCredentialRefresherType() (+4 more)

### Community 175 - "Store"
Cohesion: 0.13
Nodes (8): fakeStatusChecker, testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 176 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

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
Cohesion: 0.11
Nodes (24): Handler, newScratchStore(), SyncDocEmbeddings(), TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(), TestRetrieveUsesCacheAndSyncInvalidates(), newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists() (+16 more)

### Community 181 - "lifecycleManager"
Cohesion: 0.10
Nodes (9): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, sync/atomic.Bool, lifecycleDeps, lifecycleManager (+1 more)

### Community 182 - "backup/backup_test.go"
Cohesion: 0.33
Nodes (11): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+3 more)

### Community 183 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 184 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 185 - "Store"
Cohesion: 0.24
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

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

### Community 191 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 192 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 193 - "ReportData"
Cohesion: 0.40
Nodes (5): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), Generator, ReportData

### Community 194 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 195 - "ratelimit.go"
Cohesion: 0.31
Nodes (6): RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 196 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (15): Handler, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie(), setFlowCookie() (+7 more)

### Community 197 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 208 - "newDockerClient"
Cohesion: 0.20
Nodes (10): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+2 more)

### Community 209 - "net/http.Response"
Cohesion: 0.33
Nodes (6): retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 210 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 211 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 212 - "transform_test.go"
Cohesion: 0.24
Nodes (6): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), Transformer

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

### Community 218 - "main.tsx"
Cohesion: 0.20
Nodes (8): mermaid, react-dom, App(), cssVar(), Mermaid(), resolveColor(), USE_MOCKS, web_src_index

### Community 220 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 221 - "RequirePermission"
Cohesion: 0.38
Nodes (5): PermissionChecker, chi.Router, mountDashboardRoutes(), mountWorkflowRoutes(), RequirePermission()

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
Cohesion: 0.22
Nodes (14): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+6 more)

### Community 227 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 232 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 233 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 236 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 237 - "snapshotResponse"
Cohesion: 0.48
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 239 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 240 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 241 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 244 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 246 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

## Ambiguous Edges - Review These
- `Topbar Notification Center (deferred from V1)` → `NotificationDelivery Schema (per-channel delivery/retry)`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to
- `Topbar Notification Center (deferred from V1)` → `Notification / NotificationPage Schemas`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to

## Knowledge Gaps
- **583 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+578 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1347 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **36 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `NotificationDelivery Schema (per-channel delivery/retry)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `Notification / NotificationPage Schemas`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `Store` connect `Store` to `Errorf`, `backup/backup.go`, `Handler`, `diagnostics/diagnostics.go`, `Engine`, `time.Time`, `DecodeKey`, `Deps`, `Dispatcher`, `rewritePlaceholders`, `newHandler`, `ExportToFile`, `manifest.go`, `testApp`, `Engine`, `Handler`, `go_pkg_context`, `Handler`, `RunMigrations`, `net/http.Handler`, `NewEngine`, `NewRegistry`, `rowScanner`, `ErrorWithDetails`, `New`, `lifecycleManager`, `backup/backup_test.go`, `.call`, `dispatcher_test.go`, `.call`, `net/http.Request`, `ReportData`, `nilToStr`, `net/http.ResponseWriter`, `NewUser`, `NewStore`, `NewEngine`, `NewChecker`, `AuthedUser`, `log/slog.Logger`, `ConnectorRecord`, `apikey_scope.go`, `Manager`, `VerifyBundleFile`, `engine_maintenance_test.go`, `export_test.go`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `export_test.go` to `context.Context`, `log/slog.Logger`, `go_pkg_os`, `Store`, `testing.T`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Why does `newTestHandler()` connect `newTestHandler` to `home_assistant_test.go`, `NewService`, `Store`, `NewStore`, `snapshotResponse`, `NewEngine`, `testing.T`, `handlers_contract_test.go`, `handlers_bulk_test.go`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _583 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Client` be split into smaller, more focused modules?**
  _Cohesion score 0.03176940879444561 - nodes in this community are weakly interconnected._