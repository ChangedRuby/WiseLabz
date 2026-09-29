# Graph Report - feat-ws-add-id-and-ts-to-the-websocket-envelope  (2026-09-29)

## Corpus Check
- 878 files · ~549,129 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6566 nodes · 20813 edges · 244 communities (224 shown, 20 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1669 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2a1d3e1d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- react
- testing.T
- newDocTestStore
- icons.tsx
- @tanstack/react-query
- go_pkg_context
- context.Context
- NewChecker
- go_pkg_net_http
- DashboardPage.tsx
- newTestHandler
- dispatcher_test.go
- TemplateEditorPage.tsx
- go_pkg_testing
- App.tsx
- go_pkg_github_com_wiselabz_wiselabz_internal_connector
- UserIDFromContext
- ServiceSnapshot
- ServiceDetailPage.tsx
- ConnectorRecord
- ProfilePage.tsx
- ShareLinkPage.tsx
- react-i18next
- package.json
- Errorf
- go_pkg_os
- NewEngine
- MapTransportError
- SystemPage.tsx
- net/http.Request
- SnapshotEntity
- fixtures.ts
- Compare
- Store
- net/http.Client
- DecodeKey
- Store
- RunMigrations
- routerDeps
- ws.ts
- home_assistant/tables.go
- nilToStr
- dependencies
- net/http.ResponseWriter
- traefik/tables.go
- router.go
- NewUser
- time.Duration
- .OIDCCallback
- home_assistant_test.go
- Connector
- portainer/tables.go
- ThemeControls.tsx
- DiffViewer.tsx
- adguardhome/tables.go
- main
- docker_test.go
- store/backup_test.go
- newTestHandler
- ExportToFile
- Dispatcher
- connector_permission.go
- NewEngine
- unifi/tables.go
- NewRegistry
- Connector
- traefik_test.go
- lists.go
- AuthMiddleware
- response.go
- connector/connector.go
- HashPassword
- New
- rewritePlaceholders
- config_test.go
- settings.mock.ts
- NewStore
- GetTypeSchema
- handlers.ts
- Service
- Connector
- Register
- unifi_test.go
- motion
- .enrollmentResult
- Sanitize
- newTestHandler
- NewMalformedResponseError
- Manager
- ws/ws_test.go
- WiseLabz — Design Contract
- devDependencies
- portainer_test.go
- .runPerRevision
- Store
- backup/backup.go
- NewHTTPClient
- adguardhome_test.go
- scheduler/health_test.go
- Handler
- NewService
- gitFixture
- Runner
- export_test.go
- httpx/retry_test.go
- log/slog.Logger
- MarshalConnectorConfig
- Handler
- truenas_test.go
- Connector
- diagnostics/diagnostics.go
- newTestHarness
- keyset_test.go
- compilerOptions
- Handler
- .SnapshotDiff
- mcp/mcp_test.go
- Store
- RunbookRecord
- docdiffmodel.ts
- newHandler
- notifications/handlers_test.go
- vectorCache
- all.go
- .batchDelete
- WiseLabz — Architecture & Technical Decisions
- templates.fixtures.ts
- compilerOptions
- Handler
- handlers_contract_test.go
- handlers_actions_test.go
- newDockerClient
- Deps
- data.go
- New
- DocRecord
- config_cmd_test.go
- changes/handlers_test.go
- chat/chat.go
- ContextWithUser
- Logger
- pagination_contract_test.go
- backup/backup_test.go
- Connector
- NotificationRecord
- render_test.go
- Contributing to WiseLabz
- WiseLabz WebSocket Contract (`/ws`)
- api/changes_test.go
- connectors_health_test.go
- Handler
- Handler
- sshStdioConn
- Engine
- scripts
- NewHandler
- templatefuncs.go
- transform_test.go
- Contributor Covenant Code of Conduct
- Decision
- Decision
- Backend test performance
- Connector
- main.tsx
- net/http.Handler
- 0001 — Lab-mutating operation boundaries
- Decision
- WiseLabz Connector Guide
- Product
- EventRoutingTable.tsx
- .call
- cursor_pagination_test.go
- handlers_bulk_test.go
- connectors_maintenance_test.go
- Elector
- ReportData
- Changelog
- test-shards.sh
- mockServiceWorker.js
- ComplianceRuleRecord
- connectors_hardening_test.go
- BACKUP.md
- sync.Mutex
- release-please-config.json
- auth/handlers_test.go
- dashboard/handlers_test.go
- RateLimit
- httpx/retry.go
- time.Time
- ComputeWindow
- Store
- ShareLink
- Cache
- Decision
- Step by step
- WiseLabz
- webAuthnUser
- TestComplianceRuleValidation
- snapshotResponse
- .UpdateAuthConfig
- dialSSHStdio
- gitAuth
- scanMaintenanceWindow
- computeNextRun
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- .doRequestV5
- .GetConnectorUptime
- Store
- Backup Recovery: What Comes Back, and What Doesn't
- WiseLabz — Deployment Guide
- Scheduled Doc Export
- Mermaid.tsx
- Security Policy
- RequireConnectorRole
- all_test.go
- seedScopeFixture
- Authentication design
- Development workflow
- MfaEnrollDialog
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- Saved Views
- Fixture reuse and lifecycle tests (#405–#407)
- internal/auth/oidc.go
- dockerSSHAddr
- WiseLabz — v2 Backlog
- fakeEmbedder
- coverage-parity.sh
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- coverpkg.sh
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

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

## Communities (244 total, 20 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (192): runbookResp, runbookStepResp, templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary() (+184 more)

### Community 1 - "react"
Cohesion: 0.03
Nodes (155): Frontend, react, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze, web_src_api_generated_alerts_alerts_postalertsbulksnooze, web_src_api_generated_attention_attention_getgetattentionquerykey, web_src_api_generated_changes_changes_postchangesbulkresolve (+147 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (142): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), createKey(), testApp (+134 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.03
Nodes (134): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+126 more)

### Community 4 - "icons.tsx"
Cohesion: 0.04
Nodes (89): match-sorter, @radix-ui/react-popover, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow, web_src_api_generated_connectors_connectors_getgetconnectorsmaintenancewindowsquerykey, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors (+81 more)

### Community 5 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (66): RFC-3339, i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid (+58 more)

### Community 6 - "go_pkg_context"
Cohesion: 0.07
Nodes (20): StatusError, contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_go_webauthn_webauthn_protocol (+12 more)

### Community 7 - "context.Context"
Cohesion: 0.04
Nodes (25): sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, existingIDs(), Store, MFAFactor, Store, Store (+17 more)

### Community 8 - "NewChecker"
Cohesion: 0.06
Nodes (73): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+65 more)

### Community 9 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (37): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection (+29 more)

### Community 10 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (78): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+70 more)

### Community 11 - "newTestHandler"
Cohesion: 0.06
Nodes (71): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+63 more)

### Community 12 - "dispatcher_test.go"
Cohesion: 0.08
Nodes (68): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext() (+60 more)

### Community 13 - "TemplateEditorPage.tsx"
Cohesion: 0.04
Nodes (66): web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest, web_src_api_generated_docs_docs_postdocsdocidlock, web_src_api_generated_docs_docs_postdocsdocidlockrelease, web_src_api_generated_docs_docs_postdocsdocidversionsrevrestore (+58 more)

### Community 14 - "go_pkg_testing"
Cohesion: 0.06
Nodes (14): IsTimeout(), TestIsTimeout(), go_pkg_encoding_csv, go_pkg_encoding_json, go_pkg_github_com_gorilla_websocket, go_pkg_github_com_wiselabz_wiselabz_internal_api_apitest, go_pkg_github_com_wiselabz_wiselabz_internal_connector_all, go_pkg_net_http_httptest (+6 more)

### Community 15 - "App.tsx"
Cohesion: 0.05
Nodes (56): react-router-dom, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthloginmfa, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh (+48 more)

### Community 16 - "go_pkg_github_com_wiselabz_wiselabz_internal_connector"
Cohesion: 0.05
Nodes (39): contextKey, elevationError, TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), TestBuildInterfaceTableAttributes() (+31 more)

### Community 17 - "UserIDFromContext"
Cohesion: 0.06
Nodes (29): Handler, PermissionChecker, newToken(), sanitize(), sanitizeUser(), setRefreshCookie(), Handler, Handler (+21 more)

### Community 18 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (16): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, runTransformers(), TestRunTransformersUnknownCategoryIsNoop() (+8 more)

### Community 19 - "ServiceDetailPage.tsx"
Cohesion: 0.04
Nodes (61): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart (+53 more)

### Community 20 - "ConnectorRecord"
Cohesion: 0.06
Nodes (34): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), ConnectorRecord, Store, scanConnector() (+26 more)

### Community 21 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (51): web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin (+43 more)

### Community 22 - "ShareLinkPage.tsx"
Cohesion: 0.05
Nodes (50): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+42 more)

### Community 23 - "react-i18next"
Cohesion: 0.06
Nodes (47): Frontend shell & theme (decided 2026-06), react-error-boundary, react-i18next, sonner, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention (+39 more)

### Community 24 - "package.json"
Cohesion: 0.04
Nodes (50): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+42 more)

### Community 25 - "Errorf"
Cohesion: 0.06
Nodes (16): diffToSpec(), Handler, Handler, Handler, Handler, Handler, Handler, Handler (+8 more)

### Community 26 - "go_pkg_os"
Cohesion: 0.05
Nodes (41): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), keys(), writeExportState(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas() (+33 more)

### Community 27 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 28 - "MapTransportError"
Cohesion: 0.06
Nodes (26): Connector, NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL() (+18 more)

### Community 29 - "SystemPage.tsx"
Cohesion: 0.06
Nodes (43): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey (+35 more)

### Community 30 - "net/http.Request"
Cohesion: 0.08
Nodes (26): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+18 more)

### Community 31 - "SnapshotEntity"
Cohesion: 0.10
Nodes (48): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), buildDatasets() (+40 more)

### Community 32 - "fixtures.ts"
Cohesion: 0.06
Nodes (44): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+36 more)

### Community 33 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 34 - "Store"
Cohesion: 0.07
Nodes (14): changeServiceIDs(), placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+6 more)

### Community 35 - "net/http.Client"
Cohesion: 0.06
Nodes (18): claudeProvider, ollamaEmbedder, openAICompatibleProvider, openAIEmbedder, StubProvider, testProvider, SuggestChunk, SuggestRequest (+10 more)

### Community 36 - "DecodeKey"
Cohesion: 0.08
Nodes (26): ProviderConfig, Handler, Handler, primaryProviderConfig(), Handler, Config, mask(), redactDSN() (+18 more)

### Community 37 - "Store"
Cohesion: 0.07
Nodes (26): Config, Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+18 more)

### Community 38 - "RunMigrations"
Cohesion: 0.10
Nodes (38): main(), OpenDB(), Store, newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+30 more)

### Community 39 - "routerDeps"
Cohesion: 0.10
Nodes (35): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+27 more)

### Community 40 - "ws.ts"
Cohesion: 0.07
Nodes (34): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+26 more)

### Community 41 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 42 - "nilToStr"
Cohesion: 0.07
Nodes (14): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+6 more)

### Community 43 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 44 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (9): Handler, webAuthnFlow, Handler, Handler, WriteDataPaginated(), SinceFromDays(), Handler, github.com/go-webauthn/webauthn/webauthn.SessionData (+1 more)

### Community 45 - "traefik/tables.go"
Cohesion: 0.11
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 46 - "router.go"
Cohesion: 0.09
Nodes (31): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+23 more)

### Community 47 - "NewUser"
Cohesion: 0.19
Nodes (38): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld(), TestGetRootIsSynthetic() (+30 more)

### Community 48 - "time.Duration"
Cohesion: 0.09
Nodes (25): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+17 more)

### Community 49 - ".OIDCCallback"
Cohesion: 0.10
Nodes (18): Handler, newOIDCUser(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider, ClientIP(), hostOnly() (+10 more)

### Community 50 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (37): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+29 more)

### Community 51 - "Connector"
Cohesion: 0.08
Nodes (5): ConfigField, Connector, Connector, failingPushConnector, Connector

### Community 52 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 53 - "ThemeControls.tsx"
Cohesion: 0.11
Nodes (31): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls() (+23 more)

### Community 54 - "DiffViewer.tsx"
Cohesion: 0.07
Nodes (25): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_model_index_diff, ChangeDetailPage (+17 more)

### Community 55 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 56 - "main"
Cohesion: 0.09
Nodes (29): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+21 more)

### Community 57 - "docker_test.go"
Cohesion: 0.08
Nodes (34): generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn(), startSSHDockerServer(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout() (+26 more)

### Community 58 - "store/backup_test.go"
Cohesion: 0.09
Nodes (32): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+24 more)

### Community 59 - "newTestHandler"
Cohesion: 0.09
Nodes (26): testHandler, Handler, instanceAdminRoleFor(), Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz() (+18 more)

### Community 60 - "ExportToFile"
Cohesion: 0.14
Nodes (32): ExportToFile(), ImportFromFile(), AppVersion(), BuildManifest(), BundleCounts(), ChecksumBytes(), ManifestPath(), ReadManifest() (+24 more)

### Community 61 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 62 - "connector_permission.go"
Cohesion: 0.11
Nodes (19): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), treatAsSafeFromContext(), TreatAsSafeMethod() (+11 more)

### Community 63 - "NewEngine"
Cohesion: 0.14
Nodes (29): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+21 more)

### Community 64 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 65 - "NewRegistry"
Cohesion: 0.14
Nodes (24): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+16 more)

### Community 66 - "Connector"
Cohesion: 0.15
Nodes (11): unavailable(), SnapshotSection, groupNames(), parseGroupsV6(), Connector, unavailable(), parseGroupsV5(), unavailable() (+3 more)

### Community 67 - "traefik_test.go"
Cohesion: 0.11
Nodes (30): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+22 more)

### Community 68 - "lists.go"
Cohesion: 0.14
Nodes (30): buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), parseAdlistsV6(), parseClientsV6() (+22 more)

### Community 69 - "AuthMiddleware"
Cohesion: 0.09
Nodes (25): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, AuthMiddleware(), extractBearerToken(), hashToken() (+17 more)

### Community 70 - "response.go"
Cohesion: 0.10
Nodes (17): Handler, TestWritePaginatedOmitsNextCursor(), Error(), DataPaginatedResponse, HandleStoreError(), intQuery(), JSON(), Logger() (+9 more)

### Community 71 - "connector/connector.go"
Cohesion: 0.07
Nodes (18): Capabilities(), CapabilityDescriptor, ServiceDependency, TimeoutError, GuardedDialer(), IsDangerousIP(), environmentDependencies(), networkDependencies() (+10 more)

### Community 72 - "HashPassword"
Cohesion: 0.10
Nodes (21): mustHashDummyPassword(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty() (+13 more)

### Community 73 - "New"
Cohesion: 0.09
Nodes (29): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+21 more)

### Community 74 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 75 - "config_test.go"
Cohesion: 0.10
Nodes (27): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+19 more)

### Community 76 - "settings.mock.ts"
Cohesion: 0.09
Nodes (25): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+17 more)

### Community 77 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 78 - "GetTypeSchema"
Cohesion: 0.11
Nodes (25): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword(), TestRegisteredSchema() (+17 more)

### Community 79 - "handlers.ts"
Cohesion: 0.07
Nodes (27): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+19 more)

### Community 80 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 81 - "Connector"
Cohesion: 0.11
Nodes (10): init(), Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName(), PathSegment() (+2 more)

### Community 82 - "Register"
Cohesion: 0.14
Nodes (25): init(), supportedLifecycleVerbs(), init(), init(), init(), init(), init(), init() (+17 more)

### Community 83 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 84 - "motion"
Cohesion: 0.14
Nodes (22): motion, zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast (+14 more)

### Community 85 - ".enrollmentResult"
Cohesion: 0.15
Nodes (14): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+6 more)

### Community 86 - "Sanitize"
Cohesion: 0.13
Nodes (12): Handler, isWritableField(), decodeBulkRequest(), Handler, loggablePath(), loggableQuery(), ConfigPusher, Err() (+4 more)

### Community 87 - "newTestHandler"
Cohesion: 0.14
Nodes (23): createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler(), Handler, newTestHandler() (+15 more)

### Community 88 - "NewMalformedResponseError"
Cohesion: 0.13
Nodes (14): NewMalformedResponseError(), WantsField(), TestRequestedFields(), TestWantsField(), TestBuildHostsTableAttributes(), TestBuildHostsTableV5(), buildHostsTable(), parseHosts() (+6 more)

### Community 89 - "Manager"
Cohesion: 0.16
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 90 - "ws/ws_test.go"
Cohesion: 0.16
Nodes (23): newHeartbeat(), NewHub(), normalizeOrigin(), assertEnvelope(), decodeEnvelope(), mustMarshalEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock() (+15 more)

### Community 91 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 92 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 93 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 94 - ".runPerRevision"
Cohesion: 0.18
Nodes (12): fetchAllDocs(), fileName(), slugify(), commitMessage(), Exporter, TestCommitMessage(), commitResult, gitTarget (+4 more)

### Community 95 - "Store"
Cohesion: 0.13
Nodes (8): fakeStatusChecker, testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 96 - "backup/backup.go"
Cohesion: 0.25
Nodes (20): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+12 more)

### Community 97 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 98 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 99 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 100 - "Handler"
Cohesion: 0.16
Nodes (12): validateConfigPushRequest(), parseScheduleUpdates(), validateRotationFields(), stepAuditDetail(), validTargetType(), validVerb(), ValidateCompositeRef(), FieldError (+4 more)

### Community 101 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 102 - "gitFixture"
Cohesion: 0.31
Nodes (11): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), TestGitPerRevisionBootstrapReplayAndCap(), TestGitPerRevisionBotCatchUpAndRejectedPush(), TestGitPerRevisionMissingIntermediateAndEngineAuthor() (+3 more)

### Community 103 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 104 - "export_test.go"
Cohesion: 0.20
Nodes (17): Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile(), TestDocExportDefaultCronExprIsValid() (+9 more)

### Community 105 - "httpx/retry_test.go"
Cohesion: 0.29
Nodes (16): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+8 more)

### Community 106 - "log/slog.Logger"
Cohesion: 0.22
Nodes (14): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+6 more)

### Community 107 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 108 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 109 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 110 - "Connector"
Cohesion: 0.19
Nodes (6): controllerName(), countByKind(), unavailable(), Connector, sectionFetch, session

### Community 111 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 112 - "newTestHarness"
Cohesion: 0.20
Nodes (17): seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding(), TestSearchDocs(), seedFinding() (+9 more)

### Community 113 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 114 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 115 - "Handler"
Cohesion: 0.25
Nodes (6): response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 116 - ".SnapshotDiff"
Cohesion: 0.21
Nodes (12): decodeStoredSnapshot(), Handler, snapshotStoreError(), Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor() (+4 more)

### Community 117 - "mcp/mcp_test.go"
Cohesion: 0.17
Nodes (8): go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, go_pkg_github_com_wiselabz_wiselabz_internal_chat, changeSummary, connectorSummary, findingSummary

### Community 118 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 119 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 120 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 121 - "newHandler"
Cohesion: 0.21
Nodes (16): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth() (+8 more)

### Community 122 - "notifications/handlers_test.go"
Cohesion: 0.28
Nodes (15): AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+7 more)

### Community 123 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 124 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 126 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 127 - "templates.fixtures.ts"
Cohesion: 0.17
Nodes (13): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+5 more)

### Community 128 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 129 - "Handler"
Cohesion: 0.22
Nodes (4): updateUserRequest, Handler, writeUserWriteError(), NoContent()

### Community 130 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 131 - "handlers_actions_test.go"
Cohesion: 0.32
Nodes (14): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+6 more)

### Community 132 - "newDockerClient"
Cohesion: 0.13
Nodes (14): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+6 more)

### Community 133 - "Deps"
Cohesion: 0.34
Nodes (15): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+7 more)

### Community 134 - "data.go"
Cohesion: 0.27
Nodes (14): ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary, JobHealthEntry (+6 more)

### Community 135 - "New"
Cohesion: 0.34
Nodes (14): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+6 more)

### Community 136 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 137 - "config_cmd_test.go"
Cohesion: 0.21
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 138 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 139 - "chat/chat.go"
Cohesion: 0.19
Nodes (12): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+4 more)

### Community 140 - "ContextWithUser"
Cohesion: 0.18
Nodes (14): TestConnectorStoreErrorPaths(), TestGetRequiresGrantEvenForInstanceAdmin(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation() (+6 more)

### Community 141 - "Logger"
Cohesion: 0.19
Nodes (14): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+6 more)

### Community 142 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 143 - "backup/backup_test.go"
Cohesion: 0.26
Nodes (13): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+5 more)

### Community 144 - "Connector"
Cohesion: 0.16
Nodes (6): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), Connector

### Community 145 - "NotificationRecord"
Cohesion: 0.23
Nodes (6): Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord, Store, scanNotification()

### Community 146 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 147 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 148 - "WiseLabz WebSocket Contract (`/ws`)"
Cohesion: 0.14
Nodes (11): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention (+3 more)

### Community 149 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 150 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 151 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 152 - "Handler"
Cohesion: 0.18
Nodes (4): cron.EntryID, Handler, sync/atomic.Bool, ReadyState

### Community 153 - "sshStdioConn"
Cohesion: 0.17
Nodes (7): closeQuietly(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session, io.Closer, io.WriteCloser, net.Addr

### Community 154 - "Engine"
Cohesion: 0.17
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 155 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 156 - "NewHandler"
Cohesion: 0.20
Nodes (12): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+4 more)

### Community 157 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 158 - "transform_test.go"
Cohesion: 0.24
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 159 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 160 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 161 - "Decision"
Cohesion: 0.17
Nodes (12): 0005 — Cross-replica WebSocket event relay for active/active, Authorization and secrets, Consequences, Context, Decision, Duplicate suppression and ordering, Event ownership: local first, then relay, Mechanism: PostgreSQL LISTEN/NOTIFY (+4 more)

### Community 162 - "Backend test performance"
Cohesion: 0.17
Nodes (12): Backend test performance, CI job times, Coverage strategy, Deterministic scheduled exports (#408), Follow-ups, Measurements and validation, Measuring, Rules for new tests (+4 more)

### Community 164 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 165 - "net/http.Handler"
Cohesion: 0.24
Nodes (10): CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router, NewRouter() (+2 more)

### Community 166 - "0001 — Lab-mutating operation boundaries"
Cohesion: 0.18
Nodes (8): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Consequences, Context, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 167 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 168 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 169 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 170 - "EventRoutingTable.tsx"
Cohesion: 0.24
Nodes (9): web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_model_index_connectorcategory, web_src_api_model_index_notificationchanneltype, web_src_api_model_index_notificationconfig, CHANNEL_LABELS, eventLabel(), EventRoutingTable(), KNOWN_EVENT_TYPES (+1 more)

### Community 171 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 172 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 173 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 174 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 175 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 176 - "ReportData"
Cohesion: 0.47
Nodes (4): connectorFilter(), DefinitionSummary, Generator, ReportData

### Community 177 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 178 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 179 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 180 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 181 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 182 - "BACKUP.md"
Cohesion: 0.25
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 183 - "sync.Mutex"
Cohesion: 0.22
Nodes (4): sync.Mutex, fakeDocRegenerator, fakeNotifier, fakeQualityChecker

### Community 184 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 185 - "auth/handlers_test.go"
Cohesion: 0.25
Nodes (7): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups()

### Community 186 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 187 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 188 - "httpx/retry.go"
Cohesion: 0.43
Nodes (5): idempotent(), retryable(), sleep(), RetryPolicy, retryTransport

### Community 189 - "time.Time"
Cohesion: 0.29
Nodes (6): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, userStatus

### Community 190 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 191 - "Store"
Cohesion: 0.32
Nodes (3): Store, scanBackupRun(), BackupRun

### Community 193 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 194 - "Decision"
Cohesion: 0.25
Nodes (8): Audit, Authorization, Confirmation / step-up, Decision, Dry-run, Eligible first operation: `service.restart` only, Out of scope, Rollback expectations

### Community 195 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 196 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 197 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 198 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 199 - "snapshotResponse"
Cohesion: 0.48
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 200 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 201 - "dialSSHStdio"
Cohesion: 0.29
Nodes (5): TestDialSSHStdioHonorsContextCancel(), dialSSHStdio(), bufio.ReadWriter, golang.org/x/crypto/ssh.ClientConfig, net.Conn

### Community 202 - "gitAuth"
Cohesion: 0.33
Nodes (6): gitAuth(), installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), GitOptions, github.com/go-git/go-git/v5/plumbing/transport.AuthMethod

### Community 203 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 204 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 205 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 206 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 207 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 208 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 210 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 212 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 213 - "WiseLabz — Deployment Guide"
Cohesion: 0.33
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 214 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 215 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 216 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 217 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 218 - "all_test.go"
Cohesion: 0.50
Nodes (4): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), go_pkg_github_com_wiselabz_wiselabz_internal_connector_connectortest

### Community 219 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 221 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 222 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 223 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 224 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 225 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 228 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 229 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 230 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 231 - "Fixture reuse and lifecycle tests (#405–#407)"
Cohesion: 0.50
Nodes (4): CI measurements, Fixture reuse and lifecycle tests (#405–#407), Local measurements, Validation and coverage

## Knowledge Gaps
- **595 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+590 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1320 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Election` connect `dispatcher_test.go` to `go_pkg_context`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `Errorf` to `Handler`, `Handler`, `Store`, `response.go`, `DecodeKey`, `.UpdateAuthConfig`, `net/http.ResponseWriter`, `Logger`, `Handler`, `UserIDFromContext`, `.OIDCCallback`, `Handler`, `.SnapshotDiff`, `.enrollmentResult`, `Sanitize`, `Handler`, `net/http.Request`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `Handler`, `Deps`, `go_pkg_context`, `NewChecker`, `dispatcher_test.go`, `backup/backup_test.go`, `UserIDFromContext`, `ConnectorRecord`, `Handler`, `Handler`, `Errorf`, `Engine`, `NewEngine`, `NewHandler`, `net/http.Request`, `Store`, `RunMigrations`, `.call`, `net/http.ResponseWriter`, `NewUser`, `ReportData`, `sync.Mutex`, `main`, `store/backup_test.go`, `newTestHandler`, `ExportToFile`, `Dispatcher`, `time.Time`, `NewEngine`, `NewRegistry`, `response.go`, `HashPassword`, `New`, `rewritePlaceholders`, `NewStore`, `Manager`, `.runPerRevision`, `backup/backup.go`, `Handler`, `gitFixture`, `export_test.go`, `log/slog.Logger`, `diagnostics/diagnostics.go`, `newTestHarness`, `Handler`, `newHandler`, `notifications/handlers_test.go`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _595 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.020308615918577928 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.02681480690134074 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022997762863534676 - nodes in this community are weakly interconnected._