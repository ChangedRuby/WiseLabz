# Graph Report - test-profile-coverage-overhead-and-safely-parall  (2026-09-27)

## Corpus Check
- 875 files · ~543,524 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6513 nodes · 20685 edges · 238 communities (218 shown, 20 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1664 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4b396ad5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- SystemPage.tsx
- go_pkg_net_http
- context.Context
- AlertsPage.tsx
- @tanstack/react-query
- package.json
- NewChecker
- App.tsx
- ProfilePage.tsx
- go_pkg_context
- AuthMiddleware
- DashboardPage.tsx
- cn
- icons.tsx
- traefik_test.go
- newTestHandler
- UsersPage.tsx
- ServiceSnapshot
- NewMalformedResponseError
- go_pkg_testing
- .runPerRevision
- ServiceDetailPage.tsx
- GetTypeSchema
- Store
- rowScanner
- routerDeps
- NewEngine
- dispatcher_test.go
- Store
- Errorf
- Compare
- fixtures.ts
- home_assistant/tables.go
- truenas/tables.go
- data.go
- NewEngine
- docker_test.go
- RunMigrations
- dependencies
- traefik/tables.go
- Handler
- DecodeKey
- net/http.Request
- response.go
- ErrorWithDetails
- portainer/tables.go
- config_test.go
- NewUser
- handlers_contract_test.go
- Dispatcher
- Get
- adguardhome/tables.go
- main
- ExportToFile
- .SnapshotDiff
- unifi/tables.go
- testApp
- newTestHarness
- go_pkg_os
- AuthedUser
- settings.mock.ts
- config_cmd_test.go
- home_assistant_test.go
- rewritePlaceholders
- Registry
- channels.go
- .OIDCCallback
- runbooks_test.go
- NewRegistry
- NewStore
- connector_permission.go
- net/http.Client
- Service
- newTestHandler
- handlers.ts
- Connector
- unifi_test.go
- timeline.ts
- SuggestRequest
- Hub
- connector/connector.go
- AppearancePage.tsx
- New
- router.go
- WiseLabz — Design Contract
- devDependencies
- Sanitize
- time.Duration
- .Fetch
- Manager
- src/theme.ts
- Connector
- portainer_test.go
- adguardhome_test.go
- Connector
- chat/chat.go
- ConnectorRecord
- Deps
- NewService
- net/http.Handler
- RunDeliveryRetries
- diagnostics/diagnostics.go
- .PostMFATOTPConfirm
- Handler
- Handler
- Register
- logging.go
- Handler
- sshStdioConn
- truenas_test.go
- ReportsPage.tsx
- keyset_test.go
- Store
- ws/ws_test.go
- connectors_health_test.go
- compilerOptions
- middleware_test.go
- DocRecord
- NewHTTPClient
- RunbookRecord
- docdiffmodel.ts
- gitFixture
- fetch_test.go
- vectorCache
- all.go
- httpx/retry_test.go
- export_test.go
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- RunDocLockSweep
- newTestLifecycle
- cursor_pagination_test.go
- backup/backup.go
- Connector
- handlers_bulk_test.go
- ws.ts
- changes/handlers_test.go
- pagination_contract_test.go
- render_test.go
- docs/handlers_test.go
- Contributing to WiseLabz
- templates.fixtures.ts
- MarshalConnectorConfig
- api/changes_test.go
- ShareLink
- Handler
- newTestHandler
- Decision
- scripts
- templates_test.go
- newDockerClient
- IsSecureRequest
- findings.go
- go_pkg_github_com_wiselabz_wiselabz_internal_connector
- Store
- New
- Decision
- Connector
- SnapshotEntity
- newTestHandler
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- .call
- ReportData
- .call
- ratelimit.go
- Elector
- JobHealthRecord
- Changelog
- test-shards.sh
- mockServiceWorker.js
- AppShell.tsx
- ComplianceRuleRecord
- openapi_contract_test.go
- log/slog.Logger
- Handler
- Engine
- WiseLabz — Deployment Guide
- release-please-config.json
- ListSchemas
- apikey_scopes_test.go
- middleware.go
- provider_test.go
- Cache
- Step by step
- WiseLabz
- Handler
- TestComplianceRuleValidation
- Backend test performance
- time.Time
- Contributor Covenant Code of Conduct
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- nilToStr
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Security Policy
- main.tsx
- gitAuth
- NewWebAuthnService
- ClassifyHealth
- Authentication design
- Development workflow
- compose-smoke.sh
- ComputeWindow
- timeoutError
- BackupSchedule
- Technology stack
- MISSING — deferred & future frontend features
- Saved Views
- WiseLabz — v2 Backlog
- coverage-parity.sh
- tsconfig.json
- AGENTS.md
- fields_test.go
- setup-env.sh
- CHANGE_PROVENANCE.md
- coverpkg.sh
- vite-env.d.ts
- github.com/WiseLabz/wiselabz
- computeNextRun
- registryTestRefresher
- fakeEmbedder
- Mermaid.tsx
- fakeDocRegenerator
- versionSections
- MfaEnrollDialog
- @vitejs/plugin-react
- RetentionSettings

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

## Communities (238 total, 20 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (166): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+158 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (147): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+139 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (149): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+141 more)

### Community 3 - "SystemPage.tsx"
Cohesion: 0.02
Nodes (125): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+117 more)

### Community 4 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (40): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), buildPrompt(), TestBuildPrompt(), bulkResolveItemResult (+32 more)

### Community 5 - "context.Context"
Cohesion: 0.03
Nodes (30): sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, Connector, existingIDs(), SnapshotRecord, Store, Store (+22 more)

### Community 6 - "AlertsPage.tsx"
Cohesion: 0.04
Nodes (75): Frontend, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`) (+67 more)

### Community 7 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (84): axios, i18next, msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, web_src_api_model_index_attentionpage (+76 more)

### Community 8 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 9 - "NewChecker"
Cohesion: 0.05
Nodes (74): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+66 more)

### Community 10 - "App.tsx"
Cohesion: 0.04
Nodes (65): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken(), setMfaEnrollmentRequiredHandler() (+57 more)

### Community 11 - "ProfilePage.tsx"
Cohesion: 0.05
Nodes (45): web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin, web_src_api_generated_auth_auth_postauthloginmfawebauthnbegin, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_auth_auth_usegetauthelevatemethods (+37 more)

### Community 12 - "go_pkg_context"
Cohesion: 0.07
Nodes (18): dashboardLayout, contains(), searchString(), go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_fmt, go_pkg_github_com_coreos_go_oidc_v3_oidc (+10 more)

### Community 13 - "AuthMiddleware"
Cohesion: 0.14
Nodes (14): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), treatAsSafeFromContext(), AuthMiddleware() (+6 more)

### Community 14 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (81): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+73 more)

### Community 15 - "cn"
Cohesion: 0.04
Nodes (78): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey (+70 more)

### Community 16 - "icons.tsx"
Cohesion: 0.05
Nodes (65): web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest, web_src_api_generated_docs_docs_postdocsdocidlock, web_src_api_generated_docs_docs_postdocsdocidlockrelease, web_src_api_generated_docs_docs_postdocsdocidversionsrevrestore, web_src_api_generated_docs_docs_putdocsdocid (+57 more)

### Community 17 - "traefik_test.go"
Cohesion: 0.10
Nodes (34): AllowLoopbackForTest(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath() (+26 more)

### Community 18 - "newTestHandler"
Cohesion: 0.11
Nodes (40): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+32 more)

### Community 19 - "UsersPage.tsx"
Cohesion: 0.14
Nodes (23): customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers, ConnectorGrant (+15 more)

### Community 20 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (20): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, changePatternID(), Engine, markError(), snapshotIDOrNil(), runTransformers() (+12 more)

### Community 21 - "NewMalformedResponseError"
Cohesion: 0.12
Nodes (36): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+28 more)

### Community 22 - "go_pkg_testing"
Cohesion: 0.05
Nodes (19): TestLoggablePathMasksShareTokenUnderV1(), IsTimeout(), TestIsTimeout(), go_pkg_crypto_rsa, go_pkg_encoding_csv, go_pkg_encoding_json, go_pkg_github_com_go_jose_go_jose_v4, go_pkg_github_com_gorilla_websocket (+11 more)

### Community 23 - ".runPerRevision"
Cohesion: 0.21
Nodes (10): fileName(), slugify(), commitMessage(), Exporter, commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash (+2 more)

### Community 24 - "ServiceDetailPage.tsx"
Cohesion: 0.02
Nodes (138): ADR-0001, ADR-0003, RFC-3339, 1. `service.status`, match-sorter, motion, @radix-ui/react-popover, react (+130 more)

### Community 25 - "GetTypeSchema"
Cohesion: 0.10
Nodes (30): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+22 more)

### Community 26 - "Store"
Cohesion: 0.13
Nodes (15): Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), Handler (+7 more)

### Community 27 - "rowScanner"
Cohesion: 0.07
Nodes (26): enableFakeEmbedding(), actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), NotificationRecord, Store (+18 more)

### Community 28 - "routerDeps"
Cohesion: 0.10
Nodes (35): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+27 more)

### Community 29 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 30 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 31 - "Store"
Cohesion: 0.07
Nodes (15): seedAlert(), changeServiceIDs(), seedChange(), Store, placeholders(), changeFilterClause(), AlertRecord, ChangeRecord (+7 more)

### Community 32 - "Errorf"
Cohesion: 0.06
Nodes (28): Handler, PermissionChecker, newToken(), sanitize(), Handler, diffToSpec(), Handler, Handler (+20 more)

### Community 33 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 34 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 35 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 36 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 37 - "data.go"
Cohesion: 0.27
Nodes (14): ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary, JobHealthEntry (+6 more)

### Community 38 - "NewEngine"
Cohesion: 0.14
Nodes (29): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+21 more)

### Community 39 - "docker_test.go"
Cohesion: 0.08
Nodes (34): generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn(), startSSHDockerServer(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout() (+26 more)

### Community 40 - "RunMigrations"
Cohesion: 0.11
Nodes (36): main(), OpenDB(), newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns() (+28 more)

### Community 41 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 42 - "traefik/tables.go"
Cohesion: 0.11
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 43 - "Handler"
Cohesion: 0.18
Nodes (6): webAuthnFlow, webAuthnUser, Handler, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/go-webauthn/webauthn/webauthn.SessionData, github.com/google/uuid.UUID

### Community 44 - "DecodeKey"
Cohesion: 0.10
Nodes (23): ProviderConfig, testHandler, Handler, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey() (+15 more)

### Community 45 - "net/http.Request"
Cohesion: 0.06
Nodes (33): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+25 more)

### Community 46 - "response.go"
Cohesion: 0.10
Nodes (16): Handler, TestWritePaginatedOmitsNextCursor(), Error(), HandleStoreError(), intQuery(), JSON(), Logger(), Paginate() (+8 more)

### Community 47 - "ErrorWithDetails"
Cohesion: 0.07
Nodes (30): updateUserRequest, Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), instanceAdminRoleFor() (+22 more)

### Community 48 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 49 - "config_test.go"
Cohesion: 0.10
Nodes (28): runHealthcheck(), Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields() (+20 more)

### Community 50 - "NewUser"
Cohesion: 0.23
Nodes (35): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser(), createTestShareLink() (+27 more)

### Community 51 - "handlers_contract_test.go"
Cohesion: 0.18
Nodes (19): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+11 more)

### Community 52 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 53 - "Get"
Cohesion: 0.12
Nodes (12): Handler, isWritableField(), capitalize(), Handler, WriteElevationError(), ConfigPusher, LifecycleOp(), Get() (+4 more)

### Community 54 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 55 - "main"
Cohesion: 0.09
Nodes (27): main(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry (+19 more)

### Community 56 - "ExportToFile"
Cohesion: 0.11
Nodes (42): Export(), ExportToFile(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory() (+34 more)

### Community 57 - ".SnapshotDiff"
Cohesion: 0.20
Nodes (12): decodeStoredSnapshot(), Handler, snapshotStoreError(), Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor() (+4 more)

### Community 58 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 59 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 60 - "newTestHarness"
Cohesion: 0.30
Nodes (13): TestListAttentionItems(), TestListChanges(), TestListConnectors(), TestSearchDocs(), TestListFindings(), createConnector(), createUser(), newTestHarness() (+5 more)

### Community 61 - "go_pkg_os"
Cohesion: 0.05
Nodes (39): main(), usage(), TestCommitMessage(), keys(), writeExportState(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults() (+31 more)

### Community 62 - "AuthedUser"
Cohesion: 0.10
Nodes (34): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+26 more)

### Community 63 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 64 - "config_cmd_test.go"
Cohesion: 0.36
Nodes (7): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), io.Writer

### Community 65 - "home_assistant_test.go"
Cohesion: 0.22
Nodes (19): Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities(), TestFetchConfigIsRequestedOnce() (+11 more)

### Community 66 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 67 - "Registry"
Cohesion: 0.16
Nodes (13): Provider, StatusError, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds() (+5 more)

### Community 68 - "channels.go"
Cohesion: 0.13
Nodes (14): isSafeMethod(), IsSafeMethod(), idempotent(), buildEmailMessage(), sendSMTPChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendSMTPChannel_MissingConfig() (+6 more)

### Community 69 - ".OIDCCallback"
Cohesion: 0.14
Nodes (10): Handler, Handler, newOIDCUser(), randomOIDCToken(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider (+2 more)

### Community 70 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 71 - "NewRegistry"
Cohesion: 0.50
Nodes (11): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+3 more)

### Community 72 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 73 - "connector_permission.go"
Cohesion: 0.09
Nodes (18): APIKeyRestriction, fakeStatusChecker, testAPIKeyChecker, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole() (+10 more)

### Community 74 - "net/http.Client"
Cohesion: 0.04
Nodes (31): Connector, ollamaEmbedder, openAIEmbedder, NewAuthError(), NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities() (+23 more)

### Community 75 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 76 - "newTestHandler"
Cohesion: 0.06
Nodes (70): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+62 more)

### Community 77 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 78 - "Connector"
Cohesion: 0.06
Nodes (13): init(), ConfigField, Connector, buildRouteTable(), TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName() (+5 more)

### Community 79 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 80 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 81 - "SuggestRequest"
Cohesion: 0.17
Nodes (6): claudeProvider, openAICompatibleProvider, StubProvider, SuggestChunk, SuggestRequest, countingProvider

### Community 82 - "Hub"
Cohesion: 0.15
Nodes (6): Hub, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator, ticket

### Community 83 - "connector/connector.go"
Cohesion: 0.09
Nodes (15): Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP(), NewTimeoutError(), newWebhookClient(), AuthError (+7 more)

### Community 84 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 85 - "New"
Cohesion: 0.10
Nodes (24): confirm(), formatCounts(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag() (+16 more)

### Community 86 - "router.go"
Cohesion: 0.16
Nodes (20): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+12 more)

### Community 87 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 88 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 89 - "Sanitize"
Cohesion: 0.12
Nodes (11): Err(), Sanitize(), TestErr(), TestSanitize(), cron.EntryID, Runner, cron.Cron, HealthStore (+3 more)

### Community 90 - "time.Duration"
Cohesion: 0.09
Nodes (23): newLogger(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings, BackupSettings, Database (+15 more)

### Community 91 - ".Fetch"
Cohesion: 0.12
Nodes (12): ServiceDependency, WantsField(), Connector, environmentDependencies(), putMetadata(), unavailable(), agentEnabled(), Connector (+4 more)

### Community 92 - "Manager"
Cohesion: 0.10
Nodes (14): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, Store, ReportDefinitionRecord (+6 more)

### Community 93 - "src/theme.ts"
Cohesion: 0.13
Nodes (26): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+18 more)

### Community 94 - "Connector"
Cohesion: 0.20
Nodes (5): unavailable(), SnapshotSection, Connector, unavailable(), session

### Community 95 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 96 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 97 - "Connector"
Cohesion: 0.16
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

### Community 98 - "chat/chat.go"
Cohesion: 0.16
Nodes (16): Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings(), TestCosineSimilarityRanksClosestVectorHighest() (+8 more)

### Community 99 - "ConnectorRecord"
Cohesion: 0.15
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 100 - "Deps"
Cohesion: 0.34
Nodes (15): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+7 more)

### Community 101 - "NewService"
Cohesion: 0.25
Nodes (14): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+6 more)

### Community 102 - "net/http.Handler"
Cohesion: 0.10
Nodes (25): Config, ConnectorRoleChecker, Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape() (+17 more)

### Community 103 - "RunDeliveryRetries"
Cohesion: 0.50
Nodes (3): Dispatcher, Dispatcher, RunDeliveryRetries()

### Community 104 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 105 - ".PostMFATOTPConfirm"
Cohesion: 0.15
Nodes (14): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+6 more)

### Community 106 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 107 - "Handler"
Cohesion: 0.17
Nodes (10): validateConfigPushRequest(), NewHandler(), stepAuditDetail(), validTargetType(), validVerb(), ValidateCompositeRef(), configPushRequest, Handler (+2 more)

### Community 108 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 109 - "logging.go"
Cohesion: 0.14
Nodes (18): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+10 more)

### Community 110 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 111 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 112 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 113 - "ReportsPage.tsx"
Cohesion: 0.11
Nodes (19): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports (+11 more)

### Community 114 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 115 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 116 - "ws/ws_test.go"
Cohesion: 0.19
Nodes (18): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+10 more)

### Community 117 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 118 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 119 - "middleware_test.go"
Cohesion: 0.15
Nodes (16): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser() (+8 more)

### Community 120 - "DocRecord"
Cohesion: 0.13
Nodes (9): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary(), Store, scanMaintenanceWindow() (+1 more)

### Community 121 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 122 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 123 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 124 - "gitFixture"
Cohesion: 0.31
Nodes (11): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), TestGitPerRevisionBootstrapReplayAndCap(), TestGitPerRevisionBotCatchUpAndRejectedPush(), TestGitPerRevisionMissingIntermediateAndEngineAuthor() (+3 more)

### Community 125 - "fetch_test.go"
Cohesion: 0.11
Nodes (23): extractGroups(), newMockOIDCServer(), TestAuthURLAfterInitialization(), TestAuthURLBeforeInitialization(), TestExtractGroups(), TestInitializeFailure(), TestInitializeInvalidJSON(), TestInitializeSuccess() (+15 more)

### Community 126 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 127 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 128 - "httpx/retry_test.go"
Cohesion: 0.20
Nodes (20): retryable(), RetryTransport(), sleep(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry() (+12 more)

### Community 129 - "export_test.go"
Cohesion: 0.20
Nodes (17): fetchAllDocs(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile() (+9 more)

### Community 130 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 131 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 132 - "RunDocLockSweep"
Cohesion: 0.43
Nodes (4): Store, RunDocLockSweep(), runDocLockSweep(), DocLockRecord

### Community 133 - "newTestLifecycle"
Cohesion: 0.09
Nodes (13): newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestStandbyIsUnreadyAndRunsNoScheduler(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group (+5 more)

### Community 134 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 135 - "backup/backup.go"
Cohesion: 0.20
Nodes (23): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+15 more)

### Community 137 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 138 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 139 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 140 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 141 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 142 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 143 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 144 - "templates.fixtures.ts"
Cohesion: 0.16
Nodes (14): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+6 more)

### Community 145 - "MarshalConnectorConfig"
Cohesion: 0.20
Nodes (14): IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly(), TestSecretFieldsChangedFalseOnResubmittedUnchangedSecret() (+6 more)

### Community 146 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 147 - "ShareLink"
Cohesion: 0.19
Nodes (4): Store, ShareLink, Store, SavedView

### Community 148 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 149 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 150 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 151 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 152 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 153 - "newDockerClient"
Cohesion: 0.13
Nodes (14): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+6 more)

### Community 154 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 155 - "findings.go"
Cohesion: 0.24
Nodes (6): go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, go_pkg_github_com_wiselabz_wiselabz_internal_chat, changeSummary, connectorSummary, findingSummary

### Community 156 - "go_pkg_github_com_wiselabz_wiselabz_internal_connector"
Cohesion: 0.06
Nodes (29): Schema(), schemaFor(), TestSchemaMatchesConfig(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides() (+21 more)

### Community 157 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 158 - "New"
Cohesion: 0.20
Nodes (20): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires() (+12 more)

### Community 159 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 160 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 161 - "SnapshotEntity"
Cohesion: 0.11
Nodes (20): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable() (+12 more)

### Community 162 - "newTestHandler"
Cohesion: 0.24
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 163 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 164 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 165 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 166 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 167 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 168 - "ReportData"
Cohesion: 0.47
Nodes (4): connectorFilter(), DefinitionSummary, Generator, ReportData

### Community 169 - ".call"
Cohesion: 0.33
Nodes (7): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture, net/http.HandlerFunc

### Community 170 - "ratelimit.go"
Cohesion: 0.27
Nodes (7): TestRateLimit(), RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 171 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 172 - "JobHealthRecord"
Cohesion: 0.29
Nodes (4): JobHealthRecord, Store, scanJobHealth(), fakeHealthStore

### Community 173 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 174 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 175 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 176 - "AppShell.tsx"
Cohesion: 0.23
Nodes (9): Frontend shell & theme (decided 2026-06), react-error-boundary, sonner, AppShell, AppShell(), NavigatorBridge(), Dock(), ShellDock() (+1 more)

### Community 177 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 178 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 179 - "log/slog.Logger"
Cohesion: 0.39
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 181 - "Engine"
Cohesion: 0.17
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 182 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 183 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 184 - "ListSchemas"
Cohesion: 0.22
Nodes (10): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), supportedLifecycleVerbs(), IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType() (+2 more)

### Community 185 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 186 - "middleware.go"
Cohesion: 0.06
Nodes (27): contextKey, elevationError, TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups() (+19 more)

### Community 187 - "provider_test.go"
Cohesion: 0.29
Nodes (6): testProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream()

### Community 188 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 189 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 190 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 192 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 193 - "Backend test performance"
Cohesion: 0.25
Nodes (8): Backend test performance, CI job times, Coverage strategy, Follow-ups, Measuring, Rules for new tests, What changed, Where the time went

### Community 194 - "time.Time"
Cohesion: 0.29
Nodes (6): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, userStatus

### Community 195 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 196 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 197 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 198 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 199 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 200 - "nilToStr"
Cohesion: 0.09
Nodes (10): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+2 more)

### Community 201 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 202 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 203 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 204 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 205 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 206 - "gitAuth"
Cohesion: 0.33
Nodes (6): gitAuth(), installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), GitOptions, github.com/go-git/go-git/v5/plumbing/transport.AuthMethod

### Community 207 - "NewWebAuthnService"
Cohesion: 0.50
Nodes (4): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), github.com/go-webauthn/webauthn/webauthn.WebAuthn

### Community 208 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 209 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 210 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 211 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 212 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 214 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 216 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 217 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 218 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 229 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 232 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 235 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 236 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

## Knowledge Gaps
- **582 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+577 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1309 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `export_test.go`, `testing.T`, `newTestLifecycle`, `backup/backup.go`, `NewChecker`, `go_pkg_context`, `docs/handlers_test.go`, `Handler`, `ServiceSnapshot`, `rowScanner`, `NewEngine`, `dispatcher_test.go`, `Store`, `Errorf`, `NewEngine`, `.call`, `RunMigrations`, `.call`, `ReportData`, `DecodeKey`, `net/http.Request`, `response.go`, `ErrorWithDetails`, `NewUser`, `log/slog.Logger`, `Handler`, `Dispatcher`, `Engine`, `main`, `ExportToFile`, `testApp`, `newTestHarness`, `AuthedUser`, `Handler`, `rewritePlaceholders`, `time.Time`, `NewRegistry`, `NewStore`, `New`, `Manager`, `chat/chat.go`, `Deps`, `net/http.Handler`, `diagnostics/diagnostics.go`, `Handler`, `Handler`, `gitFixture`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `gitFixture` to `export_test.go`, `testing.T`, `context.Context`, `log/slog.Logger`, `Store`, `go_pkg_os`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `Errorf` to `Sanitize`, `.OIDCCallback`, `.PostMFATOTPConfirm`, `Handler`, `Handler`, `Handler`, `net/http.Request`, `response.go`, `ErrorWithDetails`, `logging.go`, `DecodeKey`, `Handler`, `Handler`, `Get`, `Handler`, `.SnapshotDiff`, `Store`, `Handler`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _582 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021539416511483552 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.02378247368022419 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022052384364302204 - nodes in this community are weakly interconnected._