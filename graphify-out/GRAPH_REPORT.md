# Graph Report - wiselabz-we  (2026-09-23)

## Corpus Check
- 788 files · ~477,192 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 5868 nodes · 18428 edges · 209 communities (190 shown, 19 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1452 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `be04041f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- cn
- testing.T
- react
- context.Context
- Button.tsx
- go_pkg_net_http
- DashboardPage.tsx
- @tanstack/react-query
- go_pkg_strings
- ServiceDetailPage.tsx
- go_pkg_context
- App.tsx
- ServiceSnapshot
- go_pkg_testing
- icons.tsx
- RunMigrations
- NewService
- net/http.ResponseWriter
- connector/connector.go
- rowScanner
- dispatcher_test.go
- MarshalConnectorConfig
- ErrorWithDetails
- Store
- newDocTestStore
- portainer/tables.go
- package.json
- Errorf
- fixtures.ts
- net/http.Client
- WiseLabz — Architecture & Technical Decisions
- SnapshotEntity
- Connector
- NewMalformedResponseError
- buildEntities
- net/http.Request
- IsSecureRequest
- newTestHandler
- docker_test.go
- routerDeps
- traefik/tables_test.go
- nilToStr
- UsersPage.tsx
- NewChecker
- dependencies
- response.go
- NewStore
- middleware.go
- router.go
- share_links_test.go
- ShareLinkPage.tsx
- adguardhome/tables.go
- NewRegistry
- Register
- ThemeControls.tsx
- NewEngine
- Dispatcher
- compliance/engine.go
- Config
- truenas_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- WebSocketProvider.tsx
- Connector
- config_test.go
- settings.mock.ts
- api/audit_test.go
- rewritePlaceholders
- go_pkg_os
- home_assistant_test.go
- GetTypeSchema
- handlers.ts
- VerifyBundleFile
- unifi_test.go
- timeline.ts
- ExportToFile
- Store
- time.Time
- createTestConnector
- src/theme.ts
- net/http/httptest.ResponseRecorder
- log/slog.Logger
- compliance/handlers.go
- Checker
- Manager
- WiseLabz — Design Contract
- Connector
- devDependencies
- newTestHandler
- Compare
- logging.go
- portainer_test.go
- export_test.go
- doc_test.go
- chat/chat.go
- adguardhome_test.go
- scheduler/health_test.go
- ConnectorRecord
- newTestAppWithBackupDir
- Store
- diagnostics/diagnostics.go
- Runner
- git_internal_test.go
- Hub
- ws/ws_test.go
- sshStdioConn
- ws.ts
- compilerOptions
- docs/handlers_test.go
- sync.Mutex
- Handler
- diagram.go
- traefik_test.go
- NewEngine
- docdiffmodel.ts
- chat/handlers_test.go
- handlers_contract_test.go
- reports/handlers.go
- all.go
- .batchDelete
- templates.fixtures.ts
- compilerOptions
- notifications/handlers_test.go
- changes/handlers_test.go
- backup/backup_test.go
- vectorCache
- Engine
- templatefuncs.go
- channels_test.go
- DocRecord
- testApp
- pagination_contract_test.go
- newTestHandler
- Connector
- gitTarget
- render_test.go
- New
- connectors_health_test.go
- gitFixture
- NewClient
- maintenance_test.go
- Contributing to WiseLabz
- Engine
- Connector
- registry.go
- ReportData
- quality_test.go
- transform_test.go
- Decision
- WiseLabz Connector Guide
- main.tsx
- scripts
- Handler
- Connector
- Store
- connector_permission_test.go
- keyset_test.go
- Decision
- Product
- cursor_pagination_test.go
- handlers_bulk_test.go
- connectors_maintenance_test.go
- docs_test.go
- Handler
- Changelog
- mockServiceWorker.js
- connectors_hardening_test.go
- matchEntities
- retention/retention_test.go
- release-please-config.json
- dashboard/handlers_test.go
- ComputeWindow
- notification_delivery_test.go
- RunbookRecord
- ShareLink
- Cache
- Step by step
- WiseLabz
- runbooks/handlers_test.go
- .applyChannelSecrets
- serveSSHDockerConn
- change_pattern_test.go
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- auth/handlers_test.go
- golden_snapshot_test.go
- Enforcement Guidelines
- @vitejs/plugin-react
- compose-smoke.sh
- RetentionSettings
- timeoutError
- MISSING — deferred & future frontend features
- stubEmbedder
- WiseLabz — v2 Backlog
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 212 edges
2. `Errorf()` - 164 edges
3. `Store` - 133 edges
4. `newDocTestStore()` - 129 edges
5. `SnapshotEntity` - 74 edges
6. `react` - 73 edges
7. `UserIDFromContext()` - 69 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
10. `@tanstack/react-query` - 58 edges

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

## Communities (209 total, 19 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (152): templateBody, TestAPIKeyCreateRejectsInvalidExpiryAndEmptyName(), TestAPIKeyRoutesEndToEnd(), TestAttentionAuthenticatedAccess(), TestAttentionDaysWindow(), TestAttentionEmptyList(), TestAttentionHidesUngrantedConnectors(), TestAttentionMergesAlertsAndFindings() (+144 more)

### Community 1 - "cn"
Cohesion: 0.02
Nodes (133): web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_auth_auth_usegetauthproviders, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey (+125 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (131): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+123 more)

### Community 3 - "react"
Cohesion: 0.03
Nodes (102): react, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid (+94 more)

### Community 4 - "context.Context"
Cohesion: 0.03
Nodes (31): sanitizeSessions(), Connector, Connector, Connector, existingIDs(), BackupSchedule, Store, scanComplianceRule() (+23 more)

### Community 5 - "Button.tsx"
Cohesion: 0.03
Nodes (94): Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover (+86 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (47): bulkSnoozeItemResult, bulkSnoozeRequest, updateUserRequest, newToken(), changePromptData(), diffToSpec(), stripPromptTags(), truncateUTF8() (+39 more)

### Community 7 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (78): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 5. `alert.created` (+70 more)

### Community 8 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (53): msw, @tanstack/react-query, @testing-library/react, vitest, zustand, web_src_api_generated_connectors_connectors, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage (+45 more)

### Community 9 - "go_pkg_strings"
Cohesion: 0.07
Nodes (30): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), unitSystemSummary(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent() (+22 more)

### Community 10 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (64): ADR-0001, ADR-0003, RFC-3339, 1. `service.status`, 4. `change.detected`, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush (+56 more)

### Community 11 - "go_pkg_context"
Cohesion: 0.09
Nodes (17): StatusError, dashboardLayout, ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), contains(), searchString(), go_pkg_context (+9 more)

### Community 12 - "App.tsx"
Cohesion: 0.05
Nodes (58): Frontend shell & theme (decided 2026-06), react-error-boundary, react-i18next, react-router-dom, sonner, setAccessToken(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_postauthlogin (+50 more)

### Community 13 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (18): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, changePatternID(), runTransformers() (+10 more)

### Community 14 - "go_pkg_testing"
Cohesion: 0.04
Nodes (34): TestDiagnosticsIncludesHealthVersionsAndFailures(), TestDiagnosticsRedactsSecrets(), TestDiagnosticsRoleBoundary(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), TestLoggablePathMasksShareTokenUnderV1() (+26 more)

### Community 15 - "icons.tsx"
Cohesion: 0.05
Nodes (56): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_runbooks_runbooks, web_src_api_generated_runbooks_runbooks_usegetrunbooks (+48 more)

### Community 16 - "RunMigrations"
Cohesion: 0.06
Nodes (62): main(), main(), newLogger(), runHealthcheck(), splitOrigins(), newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists() (+54 more)

### Community 17 - "NewService"
Cohesion: 0.06
Nodes (49): APIKeyChecker, Claims, ElevationClaims, ElevationToken, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, TokenPair (+41 more)

### Community 18 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (16): Handler, Handler, NewHandler(), Handler, Handler, Handler, Handler, Handler (+8 more)

### Community 19 - "connector/connector.go"
Cohesion: 0.05
Nodes (31): Connector, isTimeout(), NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL(), tryParseEntities() (+23 more)

### Community 20 - "rowScanner"
Cohesion: 0.06
Nodes (29): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), Store, scanBackupRun(), changeFilterClause() (+21 more)

### Community 21 - "dispatcher_test.go"
Cohesion: 0.16
Nodes (53): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), testLogger(), expireAlertsOnce(), NewDispatcher() (+45 more)

### Community 22 - "MarshalConnectorConfig"
Cohesion: 0.07
Nodes (40): ProviderConfig, Handler, primaryProviderConfig(), Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN() (+32 more)

### Community 23 - "ErrorWithDetails"
Cohesion: 0.08
Nodes (25): Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), contextWithShareLink(), Handler (+17 more)

### Community 24 - "Store"
Cohesion: 0.07
Nodes (13): Sanitize(), TestSanitize(), Store, placeholders(), AlertRecord, ChangeRecord, Store, scanAlert() (+5 more)

### Community 25 - "newDocTestStore"
Cohesion: 0.06
Nodes (47): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+39 more)

### Community 26 - "portainer/tables.go"
Cohesion: 0.08
Nodes (38): WantsField(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), isTimeout(), putMetadata(), buildEnvironmentTable(), buildStackTable(), cell() (+30 more)

### Community 27 - "package.json"
Cohesion: 0.04
Nodes (45): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+37 more)

### Community 28 - "Errorf"
Cohesion: 0.07
Nodes (19): Handler, applyConnectorScalarUpdates(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection(), Handler (+11 more)

### Community 29 - "fixtures.ts"
Cohesion: 0.06
Nodes (42): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+34 more)

### Community 30 - "net/http.Client"
Cohesion: 0.06
Nodes (19): claudeProvider, ollamaEmbedder, openAICompatibleProvider, openAIEmbedder, StubProvider, testProvider, SuggestChunk, SuggestRequest (+11 more)

### Community 31 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 32 - "SnapshotEntity"
Cohesion: 0.13
Nodes (42): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+34 more)

### Community 33 - "Connector"
Cohesion: 0.06
Nodes (15): init(), ConfigField, Connector, buildRouteTable(), isTimeout(), TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable() (+7 more)

### Community 34 - "NewMalformedResponseError"
Cohesion: 0.09
Nodes (41): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+33 more)

### Community 35 - "buildEntities"
Cohesion: 0.07
Nodes (37): jsonType(), TestAttributeCatalogCoversEmittedKeys(), isTimeout(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+29 more)

### Community 36 - "net/http.Request"
Cohesion: 0.08
Nodes (13): Handler, Handler, decodeBulkRequest(), Handler, Handler, Handler, Handler, oidcProviderJSON() (+5 more)

### Community 37 - "IsSecureRequest"
Cohesion: 0.09
Nodes (25): clearOIDCFlowCookie(), Handler, newOIDCUser(), oidcFlowCookieName(), randomOIDCToken(), readOIDCFlowCookie(), setOIDCFlowCookie(), validHostPort() (+17 more)

### Community 38 - "newTestHandler"
Cohesion: 0.11
Nodes (40): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+32 more)

### Community 39 - "docker_test.go"
Cohesion: 0.07
Nodes (41): IsDangerousIP(), buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer() (+33 more)

### Community 40 - "routerDeps"
Cohesion: 0.10
Nodes (34): routerDeps, TestEmbeddedSPAWithoutFrontendBuild(), chi.Router, NewRouter(), wsRoleLabel(), chi.Router, mountAuthRoutes(), mountMeRoutes() (+26 more)

### Community 41 - "traefik/tables_test.go"
Cohesion: 0.10
Nodes (32): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+24 more)

### Community 42 - "nilToStr"
Cohesion: 0.07
Nodes (14): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+6 more)

### Community 43 - "UsersPage.tsx"
Cohesion: 0.08
Nodes (34): AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), RefreshFn, setRefreshHandler(), server (+26 more)

### Community 44 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 45 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 46 - "response.go"
Cohesion: 0.08
Nodes (27): Handler, Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes() (+19 more)

### Community 47 - "NewStore"
Cohesion: 0.11
Nodes (34): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+26 more)

### Community 48 - "middleware.go"
Cohesion: 0.09
Nodes (26): AuditRecorder, ConnectorRoleChecker, contextKey, elevationError, PermissionChecker, Handler, isWritableField(), validateConfigPushRequest() (+18 more)

### Community 49 - "router.go"
Cohesion: 0.10
Nodes (28): TestEmbeddedFrontendEntryPoint(), go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat (+20 more)

### Community 50 - "share_links_test.go"
Cohesion: 0.23
Nodes (35): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser(), createTestShareLink() (+27 more)

### Community 51 - "ShareLinkPage.tsx"
Cohesion: 0.09
Nodes (27): axios, TODO: fold into docs/openapi.yaml and regenerate via `npm run gen:api`, ShareDoc, ShareLink, ShareLinkCreated, shareLinkErrorCode, shareLinksQueryKey, ShareTreeNode (+19 more)

### Community 52 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 53 - "NewRegistry"
Cohesion: 0.13
Nodes (27): Provider, SuggestResult, RegisterClaude(), TestRegisterClaudeDefaults(), registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail() (+19 more)

### Community 54 - "Register"
Cohesion: 0.10
Nodes (30): init(), newConnector(), Connector, GuardedDialer(), init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback() (+22 more)

### Community 55 - "ThemeControls.tsx"
Cohesion: 0.09
Nodes (28): MotionProvider(), AppearancePage(), ChoiceGroup(), AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented() (+20 more)

### Community 56 - "NewEngine"
Cohesion: 0.13
Nodes (26): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+18 more)

### Community 57 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 58 - "compliance/engine.go"
Cohesion: 0.12
Nodes (29): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+21 more)

### Community 59 - "Config"
Cohesion: 0.10
Nodes (20): Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings, BackupSettings, Database, DocExportGitSettings (+12 more)

### Community 60 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 61 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 62 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 63 - "WebSocketProvider.tsx"
Cohesion: 0.09
Nodes (26): Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`), web_src_api_generated_changes_changes (+18 more)

### Community 64 - "Connector"
Cohesion: 0.15
Nodes (11): unavailable(), SnapshotSection, unavailable(), Connector, isTimeout(), parseHosts(), unavailable(), unavailable() (+3 more)

### Community 65 - "config_test.go"
Cohesion: 0.09
Nodes (28): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Load(), TestAccessTokenTTLDuration() (+20 more)

### Community 66 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 67 - "api/audit_test.go"
Cohesion: 0.10
Nodes (28): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+20 more)

### Community 68 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 69 - "go_pkg_os"
Cohesion: 0.10
Nodes (22): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter, go_pkg_github_com_getkin_kin_openapi_routers (+14 more)

### Community 70 - "home_assistant_test.go"
Cohesion: 0.14
Nodes (28): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+20 more)

### Community 71 - "GetTypeSchema"
Cohesion: 0.11
Nodes (26): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion() (+18 more)

### Community 72 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 73 - "VerifyBundleFile"
Cohesion: 0.14
Nodes (25): BuildManifest(), BundleCounts(), ChecksumBytes(), ReadManifest(), WriteManifest(), failVerification(), LatestBundle(), newScratchStore() (+17 more)

### Community 74 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 75 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 76 - "ExportToFile"
Cohesion: 0.16
Nodes (22): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+14 more)

### Community 77 - "Store"
Cohesion: 0.22
Nodes (21): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+13 more)

### Community 78 - "time.Time"
Cohesion: 0.14
Nodes (22): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), golang.org/x/time/rate.Limiter, time.Time, visitor, ChangeEntry (+14 more)

### Community 79 - "createTestConnector"
Cohesion: 0.12
Nodes (22): TestDeleteOldHealthChecks(), TestGetConnectorUptimeDeterministicOutage(), TestGetConnectorUptimeNoData(), TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR(), TestRecordHealthCheckDefaults(), assertQueryPlanUsesIndex(), Store, TestDeleteOldAlertsSkipsActiveStatuses() (+14 more)

### Community 80 - "src/theme.ts"
Cohesion: 0.15
Nodes (23): @fontsource/space-mono, @fontsource-variable/space-grotesk, ColorMode, commit(), load(), Persisted, PRESETS_FONTS, ThemeState (+15 more)

### Community 81 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.17
Nodes (15): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz(), Handler (+7 more)

### Community 82 - "log/slog.Logger"
Cohesion: 0.12
Nodes (16): newLifecycleManager(), WithLogger(), Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+8 more)

### Community 83 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 84 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 85 - "Manager"
Cohesion: 0.16
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 86 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 87 - "Connector"
Cohesion: 0.15
Nodes (9): TestRateLimit(), apiMessage(), controllerName(), countByKind(), isTimeout(), statusError(), Connector, sectionFetch (+1 more)

### Community 88 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 89 - "newTestHandler"
Cohesion: 0.21
Nodes (18): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+10 more)

### Community 90 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 91 - "logging.go"
Cohesion: 0.15
Nodes (17): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+9 more)

### Community 92 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 93 - "export_test.go"
Cohesion: 0.19
Nodes (19): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+11 more)

### Community 94 - "doc_test.go"
Cohesion: 0.13
Nodes (20): Store, mustCreateUser(), newPostgresTestStore(), skipOnPostgres(), TestDocLockAcquireAfterExpiry(), TestDocLockConflict(), TestDocLockReleaseOnlyByHolder(), TestDocLockRenewalByHolder() (+12 more)

### Community 95 - "chat/chat.go"
Cohesion: 0.13
Nodes (19): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections() (+11 more)

### Community 96 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 97 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 98 - "ConnectorRecord"
Cohesion: 0.17
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 99 - "newTestAppWithBackupDir"
Cohesion: 0.13
Nodes (12): Config, RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry, NewEmbedRegistry(), spaHandler(), newTestAppWithBackupDir() (+4 more)

### Community 100 - "Store"
Cohesion: 0.12
Nodes (7): fakeStatusChecker, testAPIKeyChecker, sanitize(), APIKeyClaims, validAPIKey(), APIKey, Store

### Community 101 - "diagnostics/diagnostics.go"
Cohesion: 0.21
Nodes (18): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+10 more)

### Community 102 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 103 - "git_internal_test.go"
Cohesion: 0.13
Nodes (17): gitAuth(), Exporter, installHTTPS(), SetBeforePushForTest(), TestCommitMessage(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), TestGitAuthSSH() (+9 more)

### Community 104 - "Hub"
Cohesion: 0.14
Nodes (7): Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator, ticket

### Community 105 - "ws/ws_test.go"
Cohesion: 0.19
Nodes (18): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+10 more)

### Community 106 - "sshStdioConn"
Cohesion: 0.12
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 107 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 108 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 109 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 110 - "sync.Mutex"
Cohesion: 0.15
Nodes (8): RateLimit(), createVersion(), templateResponse(), templateVersionResponse(), golang.org/x/time/rate.Limit, sync.Mutex, limiterStore, Handler

### Community 111 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 112 - "diagram.go"
Cohesion: 0.20
Nodes (15): ServiceDependency, environmentDependencies(), networkDependencies(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid() (+7 more)

### Community 113 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 114 - "NewEngine"
Cohesion: 0.42
Nodes (16): NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists(), TestMatchingConnectorsEmptyAppliesToIsWildcard() (+8 more)

### Community 115 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 116 - "chat/handlers_test.go"
Cohesion: 0.25
Nodes (15): NewHandler(), TestCreate(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), Handler (+7 more)

### Community 117 - "handlers_contract_test.go"
Cohesion: 0.25
Nodes (15): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+7 more)

### Community 118 - "reports/handlers.go"
Cohesion: 0.18
Nodes (13): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), NewHandler() (+5 more)

### Community 119 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 121 - "templates.fixtures.ts"
Cohesion: 0.17
Nodes (13): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+5 more)

### Community 122 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 123 - "notifications/handlers_test.go"
Cohesion: 0.31
Nodes (14): AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+6 more)

### Community 124 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 125 - "backup/backup_test.go"
Cohesion: 0.25
Nodes (14): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+6 more)

### Community 126 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 127 - "Engine"
Cohesion: 0.21
Nodes (7): Engine, dedupKey(), matchReason(), TemplateFuncs(), GenerateResult, renderResult, text/template.FuncMap

### Community 128 - "templatefuncs.go"
Cohesion: 0.17
Nodes (12): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+4 more)

### Community 129 - "channels_test.go"
Cohesion: 0.20
Nodes (14): buildEmailMessage(), sendNtfyChannel(), sendSMTPChannel(), sendTelegramChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendNtfyChannel_DefaultsToNtfySh(), TestSendNtfyChannel_MissingTopic() (+6 more)

### Community 130 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 131 - "testApp"
Cohesion: 0.24
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 132 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 133 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 135 - "gitTarget"
Cohesion: 0.23
Nodes (7): commitMessage(), commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash, github.com/go-git/go-git/v5/plumbing/object.Signature, github.com/go-git/go-git/v5/plumbing.ReferenceName

### Community 136 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 137 - "New"
Cohesion: 0.36
Nodes (13): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+5 more)

### Community 138 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 139 - "gitFixture"
Cohesion: 0.36
Nodes (7): keys(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, github.com/go-git/go-git/v5/plumbing/object.Commit

### Community 140 - "NewClient"
Cohesion: 0.27
Nodes (11): clientTimeout(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects(), TestNewClientTimeout(), TestNewClientUsesCustomDialContext() (+3 more)

### Community 141 - "maintenance_test.go"
Cohesion: 0.23
Nodes (12): TestListConnectorNames(), TestListDocsGroupedByService(), Store, mustCreateMaintenanceConnector(), TestCloseMaintenanceWindow(), TestCloseMaintenanceWindowAlreadyClosed(), TestCloseMaintenanceWindowNotFound(), TestCreateAndGetActiveMaintenanceWindow() (+4 more)

### Community 142 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 143 - "Engine"
Cohesion: 0.18
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 145 - "registry.go"
Cohesion: 0.21
Nodes (9): IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError, Factory, SchemaField (+1 more)

### Community 146 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 147 - "quality_test.go"
Cohesion: 0.26
Nodes (10): TestComplianceFindingRuleDedupAndResolve(), TestComplianceRuleCRUD(), Store, newConcurrentQualityTestStore(), seedQualityConnector(), TestListQualityFindingsFilters(), TestResolveThenReopenCreatesFreshRow(), TestUpsertQualityFindingConcurrentDedup() (+2 more)

### Community 148 - "transform_test.go"
Cohesion: 0.24
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 149 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 150 - "WiseLabz Connector Guide"
Cohesion: 0.17
Nodes (12): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+4 more)

### Community 151 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 152 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 155 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 156 - "connector_permission_test.go"
Cohesion: 0.45
Nodes (10): Store, newTestConnector(), newTestUser(), TestDeleteConnectorGrant(), TestFilterConnectorIDsByGrant(), TestListConnectorGrants(), TestUpsertConnectorGrantUpdatesExistingRole(), TestUserHasConnectorRoleDefaultDeny() (+2 more)

### Community 157 - "keyset_test.go"
Cohesion: 0.33
Nodes (10): assertSameSet(), Store, T, queryPlan(), TestKeysetQueriesUseCoveringIndexes(), TestListAuditRecordsKeysetHonoursFilters(), TestListAuditRecordsKeysetTraversal(), TestListChangesKeysetTraversal() (+2 more)

### Community 158 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 159 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 160 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 161 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 162 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 163 - "docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 165 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 166 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 167 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 168 - "matchEntities"
Cohesion: 0.50
Nodes (8): matchEntities(), seedEngineConnectorWithEntities(), TestGenerateLabTopologyCreatesThenUpdatesInPlace(), TestMatchEntitiesDedupesExactExternalIDDuplicates(), TestMatchEntitiesExternalIDPrecedence(), TestMatchEntitiesFansOutAcrossConnectors(), TestMatchEntitiesHostnamePrecedenceCaseInsensitive(), TestMatchEntitiesIPPrecedence()

### Community 169 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 170 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 171 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 172 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 173 - "notification_delivery_test.go"
Cohesion: 0.39
Nodes (7): createTestNotification(), Store, TestDeliveryCreateAndList(), TestListDeliveriesStatusFilterAndPagination(), TestListDueDeliveries(), TestUpdateDeliveryResultNotFound(), TestUpdateDeliveryResultTransitionsAndClearsNextAttempt()

### Community 174 - "RunbookRecord"
Cohesion: 0.50
Nodes (3): RunbookRecord, Store, scanRunbook()

### Community 176 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 177 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 178 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 179 - "runbooks/handlers_test.go"
Cohesion: 0.48
Nodes (6): Handler, newTestHandler(), TestCreate(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete()

### Community 181 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 182 - "change_pattern_test.go"
Cohesion: 0.48
Nodes (6): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern()

### Community 183 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 184 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 185 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 186 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 187 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 188 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 189 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 190 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 191 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 192 - "auth/handlers_test.go"
Cohesion: 0.40
Nodes (4): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), emailDomainAllowed(), oidcRoleForGroups()

### Community 193 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 194 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 195 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 196 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 200 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **542 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+537 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1200 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `testApp`, `go_pkg_context`, `gitFixture`, `ServiceSnapshot`, `Engine`, `RunMigrations`, `NewService`, `net/http.ResponseWriter`, `ReportData`, `rowScanner`, `dispatcher_test.go`, `ErrorWithDetails`, `Handler`, `Errorf`, `net/http.Request`, `Handler`, `matchEntities`, `retention/retention_test.go`, `NewChecker`, `response.go`, `NewStore`, `share_links_test.go`, `NewRegistry`, `NewEngine`, `Dispatcher`, `engine_maintenance_test.go`, `rewritePlaceholders`, `VerifyBundleFile`, `ExportToFile`, `time.Time`, `net/http/httptest.ResponseRecorder`, `log/slog.Logger`, `compliance/handlers.go`, `Checker`, `Manager`, `export_test.go`, `chat/chat.go`, `newTestAppWithBackupDir`, `diagnostics/diagnostics.go`, `docs/handlers_test.go`, `sync.Mutex`, `NewEngine`, `chat/handlers_test.go`, `reports/handlers.go`, `notifications/handlers_test.go`, `changes/handlers_test.go`, `backup/backup_test.go`, `Engine`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `UserIDFromContext()` connect `net/http.ResponseWriter` to `context.Context`, `net/http.Request`, `Handler`, `routerDeps`, `sync.Mutex`, `response.go`, `middleware.go`, `NewService`, `rowScanner`, `ErrorWithDetails`, `Handler`, `Errorf`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Why does `Runner` connect `Runner` to `newTestAppWithBackupDir`, `net/http.Request`, `testApp`, `context.Context`, `New`, `go_pkg_context`, `sync.Mutex`, `NewStore`, `log/slog.Logger`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _542 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.02352437981180496 - nodes in this community are weakly interconnected._
- **Should `cn` be split into smaller, more focused modules?**
  _Cohesion score 0.021643460067816176 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.023773848766357006 - nodes in this community are weakly interconnected._