# Graph Report - wiselabz-pr387  (2026-09-24)

## Corpus Check
- 799 files · ~480,462 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 5903 nodes · 18532 edges · 208 communities (188 shown, 20 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1472 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `643f4703`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- cn
- testing.T
- SystemPage.tsx
- ServiceDetailPage.tsx
- react
- context.Context
- go_pkg_net_http
- @tanstack/react-query
- NewChecker
- package.json
- go_pkg_context
- go_pkg_testing
- ServiceSnapshot
- newTestHandler
- Config
- Errorf
- UsersPage.tsx
- icons.tsx
- MapTransportError
- net/http.Request
- go_pkg_github_com_wiselabz_wiselabz_internal_connector
- rowScanner
- MarshalConnectorConfig
- App.tsx
- dispatcher_test.go
- go_pkg_os
- fixtures.ts
- Connector
- WiseLabz — Architecture & Technical Decisions
- ChatPage.tsx
- RulesPage.tsx
- IsSecureRequest
- NewMalformedResponseError
- buildEntities
- response.go
- truenas/tables.go
- net/http.Client
- net/http.ResponseWriter
- NewEngine
- channels.go
- docker_test.go
- dependencies
- Connector
- ConnectorRecord
- routerDeps
- ErrorWithDetails
- portainer/tables.go
- ChangeDetailPage.tsx
- NewRegistry
- share_links_test.go
- SnapshotEntity
- adguardhome/tables.go
- traefik/tables.go
- NewStore
- ExportToFile
- Store
- ConnectorEditPage.tsx
- Configuration & Documentation Backup (Export/Import)
- RunMigrations
- NewEngine
- settings.mock.ts
- rewritePlaceholders
- Register
- Dispatcher
- nilToStr
- testApp
- GetTypeSchema
- Connector
- connector/connector.go
- handlers.ts
- AuthMiddleware
- Manager
- unifi_test.go
- Store
- compliance/handlers.go
- config_test.go
- time.Time
- WiseLabz — Design Contract
- AppearancePage.tsx
- New
- devDependencies
- Compare
- router.go
- portainer_test.go
- timeline.ts
- main
- newTestHandler
- NewHTTPClient
- adguardhome_test.go
- Service
- Handler
- logging.go
- home_assistant_test.go
- Runner
- middleware.go
- New
- chat/handlers_test.go
- chat/chat.go
- export_test.go
- Hub
- ws/ws_test.go
- Config
- sshStdioConn
- truenas_test.go
- diagnostics/diagnostics.go
- ws.ts
- compilerOptions
- docs/handlers_test.go
- Handler
- NewService
- diagram.go
- .Fetch
- traefik_test.go
- templates_test.go
- system/handlers_test.go
- all.go
- .batchDelete
- docdiffmodel.ts
- main.tsx
- compilerOptions
- notifications/handlers_test.go
- changes/handlers_test.go
- backup/backup_test.go
- vectorCache
- Connector
- fetch_test.go
- gitTarget
- DocRecord
- newTestLifecycle
- pagination_contract_test.go
- newTestHandler
- gitFixture
- render_test.go
- bulkFakeConnector
- templates.fixtures.ts
- config_cmd_test.go
- connectors_health_test.go
- Handler
- log/slog.Logger
- quality_test.go
- Contributing to WiseLabz
- export.go
- Engine
- templatefuncs.go
- Store
- JobHealthRecord
- Decision
- WiseLabz Connector Guide
- scripts
- Handler
- findings/handlers_authz_test.go
- Store
- Handler
- Connector
- Decision
- Product
- .call
- cursor_pagination_test.go
- connectors_maintenance_test.go
- docs_test.go
- ReportData
- Handler
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- openapi_contract_test.go
- retention/retention_test.go
- BackupSchedule
- release-please-config.json
- apikey_scope.go
- dashboard/handlers_test.go
- Options
- ComputeWindow
- Step by step
- WiseLabz
- Connector
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- fakeRefresherConnector
- .GetConnectorUptime
- Store
- engine_maintenance_test.go
- Security Policy
- APIKeyClaims
- seedScopeFixture
- Enforcement Guidelines
- compose-smoke.sh
- testHandler
- timeoutError
- RetentionSettings
- MISSING — deferred & future frontend features
- Saved Views
- stubEmbedder
- WiseLabz — v2 Backlog
- fakeQualityChecker
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 217 edges
2. `Errorf()` - 165 edges
3. `Store` - 133 edges
4. `newDocTestStore()` - 129 edges
5. `SnapshotEntity` - 74 edges
6. `react` - 73 edges
7. `UserIDFromContext()` - 70 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
10. `@tanstack/react-query` - 58 edges

## Surprising Connections (you probably didn't know these)
- `5. `alert.created`` --references--> `AlertSummaryWidget()`  [INFERRED]
  docs/WS_CONTRACT.md → web/src/components/dashboard/widgets.tsx
- `Panel (`Panel.tsx`)` --references--> `Panel()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `Radii — rounded but tight. Soft-dark, not pill-everything.` --references--> `Panel()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `Surfaces — depth from lightness steps + shadow, never borders alone` --references--> `Panel()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx
- `4. Typography` --references--> `PanelHeader()`  [INFERRED]
  docs/DESIGN.md → web/src/components/ui/Panel.tsx

## Import Cycles
- None detected.

## Communities (208 total, 20 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (174): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+166 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.03
Nodes (139): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+131 more)

### Community 2 - "cn"
Cohesion: 0.03
Nodes (99): clsx, @codemirror/lang-markdown, @codemirror/view, react-i18next, tailwind-merge, @uiw/react-codemirror, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey (+91 more)

### Community 3 - "testing.T"
Cohesion: 0.03
Nodes (117): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+109 more)

### Community 4 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (96): sonner, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_me_me_deletemesessionssessionid, web_src_api_generated_me_me_getgetmequerykey (+88 more)

### Community 5 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (101): ADR-0001, ADR-0003, 1. `service.status`, 2. `sync.progress`, 4. `change.detected`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart (+93 more)

### Community 6 - "react"
Cohesion: 0.04
Nodes (93): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover, react, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey (+85 more)

### Community 7 - "context.Context"
Cohesion: 0.03
Nodes (32): fakeStatusChecker, sanitizeSessions(), Handler, Connector, Connector, existingIDs(), Store, scanComplianceRule() (+24 more)

### Community 8 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (51): bulkSnoozeItemResult, bulkSnoozeRequest, updateUserRequest, newToken(), changePromptData(), diffToSpec(), stripPromptTags(), truncateUTF8() (+43 more)

### Community 9 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (71): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 3. `sync.complete`, 5. `alert.created`, 6. `alert.resolved` (+63 more)

### Community 10 - "NewChecker"
Cohesion: 0.06
Nodes (70): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+62 more)

### Community 11 - "package.json"
Cohesion: 0.04
Nodes (74): codemirror, @codemirror/commands, @codemirror/state, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, @fontsource/ibm-plex-mono (+66 more)

### Community 12 - "go_pkg_context"
Cohesion: 0.07
Nodes (20): dashboardLayout, versionSections(), ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), TemplateVersionSection, contains(), searchString() (+12 more)

### Community 13 - "go_pkg_testing"
Cohesion: 0.05
Nodes (27): TestLoggablePathMasksShareTokenUnderV1(), NewHandler(), Handler, newTestHandler(), TestCreate(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete() (+19 more)

### Community 14 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (23): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, Sanitize(), TestSanitize(), changePatternID() (+15 more)

### Community 15 - "newTestHandler"
Cohesion: 0.07
Nodes (63): AssertMatchesSpec(), loadSpec(), specPath(), actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews() (+55 more)

### Community 16 - "Config"
Cohesion: 0.06
Nodes (45): Config, LogSettings, IsSSHRemote(), idempotent(), retryable(), RetryTransport(), sleep(), do() (+37 more)

### Community 17 - "Errorf"
Cohesion: 0.07
Nodes (21): PermissionChecker, Handler, NewHandler(), Handler, Handler, Handler, Handler, contextWithShareLink() (+13 more)

### Community 18 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (49): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), RefreshFn, setRefreshHandler() (+41 more)

### Community 19 - "icons.tsx"
Cohesion: 0.07
Nodes (51): i18next, zustand, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs, categoryIcon, buildCommands(), CommandPalette() (+43 more)

### Community 20 - "MapTransportError"
Cohesion: 0.05
Nodes (24): Connector, NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector, ReadBody() (+16 more)

### Community 21 - "net/http.Request"
Cohesion: 0.05
Nodes (23): Handler, Handler, Handler, Handler, TestRateLimit(), RateLimit(), Handler, oidcProviderJSON() (+15 more)

### Community 22 - "go_pkg_github_com_wiselabz_wiselabz_internal_connector"
Cohesion: 0.07
Nodes (25): TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable() (+17 more)

### Community 23 - "rowScanner"
Cohesion: 0.07
Nodes (24): Dispatcher, actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), NotificationRecord, Store (+16 more)

### Community 24 - "MarshalConnectorConfig"
Cohesion: 0.07
Nodes (39): ProviderConfig, Handler, primaryProviderConfig(), Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN() (+31 more)

### Community 25 - "App.tsx"
Cohesion: 0.07
Nodes (40): react-router-dom, setAccessToken(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh, web_src_api_generated_auth_auth_usegetauthproviders (+32 more)

### Community 26 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (49): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+41 more)

### Community 27 - "go_pkg_os"
Cohesion: 0.06
Nodes (38): gitAuth(), Exporter, installHTTPS(), TestCommitMessage(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), TestGitAuthSSH(), writeTestKey() (+30 more)

### Community 28 - "fixtures.ts"
Cohesion: 0.06
Nodes (41): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+33 more)

### Community 29 - "Connector"
Cohesion: 0.06
Nodes (17): init(), ConfigField, Connector, TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable() (+9 more)

### Community 30 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 31 - "ChatPage.tsx"
Cohesion: 0.05
Nodes (41): react-markdown, remark-gfm, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations (+33 more)

### Community 32 - "RulesPage.tsx"
Cohesion: 0.05
Nodes (42): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+34 more)

### Community 33 - "IsSecureRequest"
Cohesion: 0.08
Nodes (27): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), clearOIDCFlowCookie(), emailDomainAllowed(), Handler, newOIDCUser(), oidcFlowCookieName(), oidcRoleForGroups() (+19 more)

### Community 34 - "NewMalformedResponseError"
Cohesion: 0.09
Nodes (42): NewMalformedResponseError(), TestBuildHostsTableAttributes(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+34 more)

### Community 35 - "buildEntities"
Cohesion: 0.07
Nodes (37): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+29 more)

### Community 36 - "response.go"
Cohesion: 0.07
Nodes (28): Handler, Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes() (+20 more)

### Community 37 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 38 - "net/http.Client"
Cohesion: 0.07
Nodes (17): claudeProvider, ollamaEmbedder, openAICompatibleProvider, openAIEmbedder, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet() (+9 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (18): Handler, isWritableField(), validateConfigPushRequest(), Handler, Handler, decodeBulkRequest(), Handler, Handler (+10 more)

### Community 40 - "NewEngine"
Cohesion: 0.13
Nodes (32): NewHandler(), Engine, NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError() (+24 more)

### Community 41 - "channels.go"
Cohesion: 0.09
Nodes (37): AllowLoopbackForTest(), buildEmailMessage(), discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendNtfyChannel(), sendSlackChannel(), sendSMTPChannel() (+29 more)

### Community 42 - "docker_test.go"
Cohesion: 0.07
Nodes (39): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+31 more)

### Community 43 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 44 - "Connector"
Cohesion: 0.08
Nodes (15): TimeoutError, NewAuthError(), NewTimeoutError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap(), Connector, apiMessage(), controllerName() (+7 more)

### Community 45 - "ConnectorRecord"
Cohesion: 0.08
Nodes (20): Store, ChatConversationRecord, Store, ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr() (+12 more)

### Community 46 - "routerDeps"
Cohesion: 0.11
Nodes (31): routerDeps, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+23 more)

### Community 47 - "ErrorWithDetails"
Cohesion: 0.11
Nodes (15): Handler, sanitize(), Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, validTargetType() (+7 more)

### Community 48 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 49 - "ChangeDetailPage.tsx"
Cohesion: 0.07
Nodes (31): Frontend shell & theme (decided 2026-06), react-error-boundary, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid (+23 more)

### Community 50 - "NewRegistry"
Cohesion: 0.12
Nodes (25): Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail() (+17 more)

### Community 51 - "share_links_test.go"
Cohesion: 0.23
Nodes (35): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser(), createTestShareLink() (+27 more)

### Community 52 - "SnapshotEntity"
Cohesion: 0.13
Nodes (33): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable() (+25 more)

### Community 53 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 54 - "traefik/tables.go"
Cohesion: 0.14
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 55 - "NewStore"
Cohesion: 0.13
Nodes (30): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+22 more)

### Community 56 - "ExportToFile"
Cohesion: 0.14
Nodes (31): ExportToFile(), AppVersion(), BuildManifest(), BundleCounts(), ChecksumBytes(), ManifestPath(), ReadManifest(), corruptFile() (+23 more)

### Community 57 - "Store"
Cohesion: 0.09
Nodes (8): placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange(), ChangeSummary

### Community 58 - "ConnectorEditPage.tsx"
Cohesion: 0.08
Nodes (23): RFC-3339, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsschema (+15 more)

### Community 59 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 60 - "RunMigrations"
Cohesion: 0.11
Nodes (26): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+18 more)

### Community 61 - "NewEngine"
Cohesion: 0.15
Nodes (25): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+17 more)

### Community 62 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 63 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 64 - "Register"
Cohesion: 0.12
Nodes (26): init(), init(), init(), init(), init(), init(), init(), init() (+18 more)

### Community 65 - "Dispatcher"
Cohesion: 0.19
Nodes (8): findChannel(), findRoute(), Dispatcher, severityRank(), shouldSkipRoute(), sync.WaitGroup, channelCfg, routeCfg

### Community 66 - "nilToStr"
Cohesion: 0.10
Nodes (10): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+2 more)

### Community 67 - "testApp"
Cohesion: 0.11
Nodes (19): mustHashDummyPassword(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty() (+11 more)

### Community 68 - "GetTypeSchema"
Cohesion: 0.12
Nodes (25): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+17 more)

### Community 69 - "Connector"
Cohesion: 0.19
Nodes (6): unavailable(), SnapshotSection, Connector, parseHosts(), unavailable(), session

### Community 70 - "connector/connector.go"
Cohesion: 0.10
Nodes (20): GuardedDialer(), IsDangerousIP(), clientTimeout(), NewClient(), NewTransport(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects(), TestNewClientTimeout() (+12 more)

### Community 71 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 72 - "AuthMiddleware"
Cohesion: 0.12
Nodes (21): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, AuthMiddleware(), extractBearerToken(), assertElevationAuditCalls(), boolLabel() (+13 more)

### Community 73 - "Manager"
Cohesion: 0.14
Nodes (10): NewHandler(), cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord (+2 more)

### Community 74 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 75 - "Store"
Cohesion: 0.22
Nodes (21): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+13 more)

### Community 76 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 77 - "config_test.go"
Cohesion: 0.12
Nodes (21): Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH(), TestLoadFromYAML() (+13 more)

### Community 78 - "time.Time"
Cohesion: 0.15
Nodes (22): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), golang.org/x/time/rate.Limiter, time.Time, visitor, ChangeEntry (+14 more)

### Community 79 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 80 - "AppearancePage.tsx"
Cohesion: 0.15
Nodes (20): MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+12 more)

### Community 81 - "New"
Cohesion: 0.22
Nodes (19): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires() (+11 more)

### Community 82 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 83 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 84 - "router.go"
Cohesion: 0.16
Nodes (20): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+12 more)

### Community 85 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 86 - "timeline.ts"
Cohesion: 0.15
Nodes (14): installMockWebSocket(), Listenerish, MockWebSocket, Emit, env(), heartbeat(), newId(), ScheduledEvent (+6 more)

### Community 87 - "main"
Cohesion: 0.14
Nodes (16): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+8 more)

### Community 88 - "newTestHandler"
Cohesion: 0.22
Nodes (17): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+9 more)

### Community 89 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 90 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 91 - "Service"
Cohesion: 0.20
Nodes (10): Claims, ElevationClaims, ElevationToken, TokenPair, Service, hasAudience(), newTokenID(), go_pkg_github_com_golang_jwt_jwt_v5 (+2 more)

### Community 92 - "Handler"
Cohesion: 0.16
Nodes (8): applyConnectorScalarUpdates(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection(), FieldError, updateConnectorRequest

### Community 93 - "logging.go"
Cohesion: 0.17
Nodes (17): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+9 more)

### Community 94 - "home_assistant_test.go"
Cohesion: 0.22
Nodes (19): Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities(), TestFetchConfigIsRequestedOnce() (+11 more)

### Community 95 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 96 - "middleware.go"
Cohesion: 0.16
Nodes (14): AuditRecorder, ConnectorRoleChecker, contextKey, elevationError, UserStatusChecker, elevationFailureReason(), hashToken(), recordElevationAudit() (+6 more)

### Community 97 - "New"
Cohesion: 0.18
Nodes (18): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+10 more)

### Community 98 - "chat/handlers_test.go"
Cohesion: 0.19
Nodes (18): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth() (+10 more)

### Community 99 - "chat/chat.go"
Cohesion: 0.17
Nodes (17): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings() (+9 more)

### Community 100 - "export_test.go"
Cohesion: 0.20
Nodes (17): fetchAllDocs(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile() (+9 more)

### Community 101 - "Hub"
Cohesion: 0.13
Nodes (8): Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, sync.RWMutex, broadcastMsg, Client, Revalidator, ticket

### Community 102 - "ws/ws_test.go"
Cohesion: 0.19
Nodes (18): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+10 more)

### Community 103 - "Config"
Cohesion: 0.14
Nodes (14): Config, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+6 more)

### Community 104 - "sshStdioConn"
Cohesion: 0.12
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 105 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 106 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 107 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 108 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 109 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 110 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 111 - "NewService"
Cohesion: 0.21
Nodes (15): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+7 more)

### Community 112 - "diagram.go"
Cohesion: 0.20
Nodes (15): ServiceDependency, environmentDependencies(), networkDependencies(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid() (+7 more)

### Community 113 - ".Fetch"
Cohesion: 0.18
Nodes (7): WantsField(), TestRequestedFields(), TestWantsField(), putMetadata(), unavailable(), Connector, dockerSectionSpec

### Community 114 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 115 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 116 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 117 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 119 - "docdiffmodel.ts"
Cohesion: 0.23
Nodes (13): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, fold(), toUnits(), DiffLine, DiffLineType (+5 more)

### Community 120 - "main.tsx"
Cohesion: 0.17
Nodes (11): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers (+3 more)

### Community 121 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 122 - "notifications/handlers_test.go"
Cohesion: 0.31
Nodes (14): AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+6 more)

### Community 123 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 124 - "backup/backup_test.go"
Cohesion: 0.25
Nodes (14): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+6 more)

### Community 125 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 126 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 127 - "fetch_test.go"
Cohesion: 0.19
Nodes (14): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+6 more)

### Community 128 - "gitTarget"
Cohesion: 0.21
Nodes (8): commitMessage(), commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash, github.com/go-git/go-git/v5/plumbing/object.Signature, github.com/go-git/go-git/v5/plumbing.ReferenceName, github.com/go-git/go-git/v5/plumbing/transport.AuthMethod

### Community 129 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 130 - "newTestLifecycle"
Cohesion: 0.18
Nodes (11): newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, sync/atomic.Bool (+3 more)

### Community 131 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 132 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 133 - "gitFixture"
Cohesion: 0.32
Nodes (8): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, Exporter, github.com/go-git/go-git/v5/plumbing/object.Commit

### Community 134 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 135 - "bulkFakeConnector"
Cohesion: 0.14
Nodes (3): actionConnector, bulkFakeConnector, failingPushConnector

### Community 136 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (11): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, resolveToken(), Snapshot (+3 more)

### Community 137 - "config_cmd_test.go"
Cohesion: 0.23
Nodes (10): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+2 more)

### Community 138 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 139 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 140 - "log/slog.Logger"
Cohesion: 0.22
Nodes (8): Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep(), Engine, log/slog.Logger, DocLockRecord

### Community 141 - "quality_test.go"
Cohesion: 0.23
Nodes (11): TestComplianceFindingRuleDedupAndResolve(), TestComplianceRuleCRUD(), Store, newConcurrentQualityTestStore(), seedQualityConnector(), TestListQualityFindingsFilters(), TestResolveThenReopenCreatesFreshRow(), TestUpsertQualityFindingConcurrentDedup() (+3 more)

### Community 142 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 143 - "export.go"
Cohesion: 0.24
Nodes (10): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), fileName(), slugify() (+2 more)

### Community 144 - "Engine"
Cohesion: 0.18
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 145 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 146 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 147 - "JobHealthRecord"
Cohesion: 0.29
Nodes (4): JobHealthRecord, Store, scanJobHealth(), fakeHealthStore

### Community 148 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 149 - "WiseLabz Connector Guide"
Cohesion: 0.17
Nodes (12): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+4 more)

### Community 150 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 152 - "findings/handlers_authz_test.go"
Cohesion: 0.45
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 153 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 156 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 157 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 158 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 159 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 160 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 161 - "docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 162 - "ReportData"
Cohesion: 0.47
Nodes (4): connectorFilter(), DefinitionSummary, Generator, ReportData

### Community 164 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 165 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 166 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 167 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 168 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 169 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 170 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 171 - "apikey_scope.go"
Cohesion: 0.43
Nodes (6): APIKeyRestriction, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), isSafeMethod(), TestClampConnectorRole()

### Community 172 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 173 - "Options"
Cohesion: 0.25
Nodes (7): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn, Options

### Community 174 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 175 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 176 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 178 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 179 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 180 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 181 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 182 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 183 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 185 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 187 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 188 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 189 - "APIKeyClaims"
Cohesion: 0.50
Nodes (3): testAPIKeyChecker, APIKeyClaims, validAPIKey()

### Community 190 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 191 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 192 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 193 - "testHandler"
Cohesion: 0.50
Nodes (3): testHandler, Handler, instanceAdminRoleFor()

### Community 197 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 198 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **542 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+537 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1200 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `newTestLifecycle`, `testing.T`, `gitFixture`, `NewChecker`, `Handler`, `go_pkg_context`, `go_pkg_testing`, `ServiceSnapshot`, `Engine`, `Errorf`, `net/http.Request`, `Handler`, `rowScanner`, `findings/handlers_authz_test.go`, `dispatcher_test.go`, `Handler`, `.call`, `ReportData`, `Handler`, `response.go`, `NewEngine`, `retention/retention_test.go`, `ErrorWithDetails`, `NewRegistry`, `share_links_test.go`, `NewStore`, `ExportToFile`, `engine_maintenance_test.go`, `RunMigrations`, `NewEngine`, `rewritePlaceholders`, `testHandler`, `Dispatcher`, `testApp`, `Manager`, `compliance/handlers.go`, `time.Time`, `main`, `Handler`, `New`, `chat/handlers_test.go`, `chat/chat.go`, `export_test.go`, `Config`, `diagnostics/diagnostics.go`, `docs/handlers_test.go`, `notifications/handlers_test.go`, `changes/handlers_test.go`, `backup/backup_test.go`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Why does `Connector` connect `buildEntities` to `go_pkg_context`, `net/http.Client`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `gitFixture` to `testing.T`, `export_test.go`, `context.Context`, `Store`, `log/slog.Logger`, `go_pkg_os`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _542 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.020342495636998255 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.025889164598842017 - nodes in this community are weakly interconnected._
- **Should `cn` be split into smaller, more focused modules?**
  _Cohesion score 0.03014271653543307 - nodes in this community are weakly interconnected._