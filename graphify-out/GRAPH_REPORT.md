# Graph Report - wiselabz-wd-complete  (2026-09-23)

## Corpus Check
- 780 files · ~473,162 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 5794 nodes · 18170 edges · 201 communities (181 shown, 20 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1448 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `37e33414`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- react
- context.Context
- newTestApp
- newDocTestStore
- testing.T
- @tanstack/react-query
- go_pkg_net_http
- connector/connector.go
- go_pkg_context
- ServiceDetailPage.tsx
- go_pkg_testing
- net/http.Request
- DashboardPage.tsx
- ServiceSnapshot
- cn
- NewMalformedResponseError
- Store
- Runner
- AlertsPage.tsx
- net/http.ResponseWriter
- net/http.Client
- Connector
- ConnectorRecord
- icons.tsx
- export_test.go
- ServicesPage.tsx
- api/audit_test.go
- dispatcher_test.go
- NewEngine
- package.json
- Config
- App.tsx
- WebSocketProvider.tsx
- fixtures.ts
- WiseLabz — Architecture & Technical Decisions
- docker_test.go
- truenas/tables.go
- nilToStr
- IsSecureRequest
- share_links_test.go
- newTestHandler
- RunMigrations
- home_assistant/tables.go
- go_pkg_os
- dependencies
- GetTypeSchema
- response.go
- NewStore
- portainer/tables.go
- DecodeKey
- ExportToFile
- adguardhome/tables.go
- adguardhome_test.go
- NewChecker
- routerDeps
- traefik/tables.go
- traefik_test.go
- unifi/tables.go
- Dispatcher
- Configuration & Documentation Backup (Export/Import)
- AuthedUser
- Get
- settings.mock.ts
- ErrorWithDetails
- backup/backup.go
- rewritePlaceholders
- middleware.go
- testApp
- handlers.ts
- SystemPage.tsx
- NewEngine
- unifi_test.go
- AuditRecord
- timeline.ts
- compliance/engine.go
- Checker
- WiseLabz — Design Contract
- NewRegistry
- Connector
- devDependencies
- newTestHandler
- New
- Compare
- logging.go
- SnapshotEntity
- portainer_test.go
- useRole.ts
- AuthMiddleware
- Service
- src/theme.ts
- chat/chat.go
- router.go
- home_assistant_test.go
- middleware_test.go
- Register
- NewRouter
- Handler
- MarshalConnectorConfig
- GuardedDialer
- truenas_test.go
- diagnostics/diagnostics.go
- keyset_test.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- Handler
- sshStdioConn
- docdiffmodel.ts
- AppearancePage.tsx
- net/http.Handler
- templates_test.go
- main
- handlers_contract_test.go
- system/handlers_test.go
- all.go
- data.go
- log/slog.Logger
- compilerOptions
- vectorCache
- render_test.go
- DocRecord
- templates.fixtures.ts
- config_cmd_test.go
- changes/handlers_test.go
- pagination_contract_test.go
- newTestHandler
- pfsense.go
- store/theme.ts
- Handler
- changes_test.go
- Connector
- Contributing to WiseLabz
- net/http/httptest.ResponseRecorder
- NewService
- templatefuncs.go
- NotificationRecord
- store/backup_test.go
- Store
- transform_test.go
- Decision
- WiseLabz Connector Guide
- main.tsx
- scripts
- Store
- compliance/engine_test.go
- channels.go
- Engine
- Decision
- Product
- handlers_bulk_test.go
- connectors_maintenance_test.go
- Changelog
- mockServiceWorker.js
- ThemeControls.tsx
- ComplianceRuleRecord
- compliance_rules_test.go
- BackupSchedule
- useTheme
- release-please-config.json
- decodeBulkRequest
- dashboard/handlers_test.go
- RateLimit
- time.Time
- ComputeWindow
- ShareLink
- Step by step
- WiseLabz
- settings.ts
- .applyChannelSecrets
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- walkCursorPages
- Store
- engine_maintenance_test.go
- Security Policy
- auth/handlers_test.go
- CORS
- routerOperations
- Enforcement Guidelines
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- MISSING — deferred & future frontend features
- internal/auth/oidc.go
- fields_test.go
- dockerSSHAddr
- WiseLabz — v2 Backlog
- tsconfig.json
- AGENTS.md
- password.go
- setup-env.sh
- CHANGE_PROVENANCE.md
- fakeNotifier
- snapshotChangingNotifier
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 212 edges
2. `Errorf()` - 156 edges
3. `newDocTestStore()` - 129 edges
4. `Store` - 127 edges
5. `SnapshotEntity` - 74 edges
6. `react` - 72 edges
7. `cn()` - 69 edges
8. `UserIDFromContext()` - 68 edges
9. `NewStore()` - 65 edges
10. `@tanstack/react-query` - 57 edges

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

## Communities (201 total, 20 thin omitted)

### Community 0 - "react"
Cohesion: 0.03
Nodes (134): react, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_auth_auth_usegetauthproviders, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey (+126 more)

### Community 1 - "context.Context"
Cohesion: 0.03
Nodes (40): fakeStatusChecker, sanitizeSessions(), Connector, Sanitize(), TestSanitize(), Store, existingIDs(), placeholders() (+32 more)

### Community 2 - "newTestApp"
Cohesion: 0.03
Nodes (134): TestAPIKeyCreateRejectsInvalidExpiryAndEmptyName(), TestAPIKeyRoutesEndToEnd(), testApp, loginRefreshCookie(), seedLocalUser(), TestChangePasswordWrongCurrentPassword(), TestDeleteSessionNotOwner(), TestDeleteSessionSuccess() (+126 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.03
Nodes (130): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+122 more)

### Community 4 - "testing.T"
Cohesion: 0.02
Nodes (127): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+119 more)

### Community 5 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (87): axios, i18next, msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, AXIOS_INSTANCE (+79 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (45): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection (+37 more)

### Community 7 - "connector/connector.go"
Cohesion: 0.03
Nodes (42): Connector, isTimeout(), ServiceDependency, NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), WantsField(), setHeaders() (+34 more)

### Community 8 - "go_pkg_context"
Cohesion: 0.08
Nodes (18): StatusError, contains(), searchString(), go_pkg_context, go_pkg_crypto_tls, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt (+10 more)

### Community 9 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (78): ADR-0001, ADR-0003, RFC-3339, Frontend shell & theme (decided 2026-06), 1. `service.status`, react-i18next, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention (+70 more)

### Community 10 - "go_pkg_testing"
Cohesion: 0.05
Nodes (27): Handler, newTestHandler(), TestCreate(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable() (+19 more)

### Community 11 - "net/http.Request"
Cohesion: 0.06
Nodes (26): Handler, Handler, Handler, Handler, Handler, Handler, Handler, Handler (+18 more)

### Community 12 - "DashboardPage.tsx"
Cohesion: 0.05
Nodes (66): 4. `change.detected`, 5. `alert.created`, 8. `doc.generated`, web_src_api_generated_alerts_alerts_usegetalerts, web_src_api_generated_dashboard_dashboard_getdashboardlayout, web_src_api_generated_dashboard_dashboard_getdashboardlayoutadmindefault, web_src_api_generated_dashboard_dashboard_getgetdashboardlayoutadmindefaultquerykey, web_src_api_generated_dashboard_dashboard_postdashboardlayoutreset (+58 more)

### Community 13 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (17): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, runTransformers(), TestRunTransformersUnknownCategoryIsNoop() (+9 more)

### Community 14 - "cn"
Cohesion: 0.04
Nodes (57): clsx, tailwind-merge, web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview (+49 more)

### Community 15 - "NewMalformedResponseError"
Cohesion: 0.07
Nodes (43): SnapshotSection, NewMalformedResponseError(), ReadBody(), TestReadBodyLimit(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable() (+35 more)

### Community 16 - "Store"
Cohesion: 0.05
Nodes (33): Provider, Config, Handler, Embedder, EmbedRegistry, Registry, NewHandler(), NewHandler() (+25 more)

### Community 17 - "Runner"
Cohesion: 0.06
Nodes (41): newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart() (+33 more)

### Community 18 - "AlertsPage.tsx"
Cohesion: 0.05
Nodes (56): Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss (+48 more)

### Community 19 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (18): Handler, Handler, Handler, oidcProviderJSON(), boolToInt(), Handler, stripLogControlChars(), cron.EntryID (+10 more)

### Community 20 - "net/http.Client"
Cohesion: 0.04
Nodes (27): claudeProvider, ollamaEmbedder, openAICompatibleProvider, openAIEmbedder, StubProvider, testProvider, SuggestChunk, SuggestRequest (+19 more)

### Community 21 - "Connector"
Cohesion: 0.05
Nodes (19): isTimeout(), ConfigField, Connector, buildRouteTable(), isTimeout(), buildGatewayTable(), isTimeout(), primaryGatewayName() (+11 more)

### Community 22 - "ConnectorRecord"
Cohesion: 0.06
Nodes (33): changeFilterClause(), scanAlert(), scanChange(), ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr() (+25 more)

### Community 23 - "icons.tsx"
Cohesion: 0.07
Nodes (50): web_src_api_generated_connectors_connectors, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidversionsrevrestore, web_src_api_generated_docs_docs_usegetdocsdocid, web_src_api_generated_docs_docs_usegetdocsdocidversions (+42 more)

### Community 24 - "export_test.go"
Cohesion: 0.06
Nodes (44): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+36 more)

### Community 25 - "ServicesPage.tsx"
Cohesion: 0.05
Nodes (47): match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow, web_src_api_generated_connectors_connectors_getgetconnectorsmaintenancewindowsquerykey (+39 more)

### Community 26 - "api/audit_test.go"
Cohesion: 0.05
Nodes (53): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+45 more)

### Community 27 - "dispatcher_test.go"
Cohesion: 0.15
Nodes (52): NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore(), setChannelAndRoutingConfig(), setChannelConfig(), setChannelConfigJSON() (+44 more)

### Community 28 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 29 - "package.json"
Cohesion: 0.04
Nodes (48): codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+40 more)

### Community 30 - "Config"
Cohesion: 0.06
Nodes (33): Config, LogSettings, IsSSHRemote(), clientTimeout(), NewTransport(), TestNewTransportTLS(), Cache, New() (+25 more)

### Community 31 - "App.tsx"
Cohesion: 0.06
Nodes (41): setAccessToken(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh, web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_model_index_authsession, web_src_api_model_index_oidccallbackrequest (+33 more)

### Community 32 - "WebSocketProvider.tsx"
Cohesion: 0.06
Nodes (42): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 6. `alert.resolved` (+34 more)

### Community 33 - "fixtures.ts"
Cohesion: 0.06
Nodes (41): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+33 more)

### Community 34 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 35 - "docker_test.go"
Cohesion: 0.06
Nodes (44): IsDangerousIP(), buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+36 more)

### Community 36 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 37 - "nilToStr"
Cohesion: 0.08
Nodes (14): Store, nilToStr(), DocVersionRecord, Store, Store, Store, RunbookRecord, Store (+6 more)

### Community 38 - "IsSecureRequest"
Cohesion: 0.09
Nodes (25): clearOIDCFlowCookie(), Handler, newOIDCUser(), oidcFlowCookieName(), randomOIDCToken(), readOIDCFlowCookie(), setOIDCFlowCookie(), validHostPort() (+17 more)

### Community 39 - "share_links_test.go"
Cohesion: 0.17
Nodes (41): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld() (+33 more)

### Community 40 - "newTestHandler"
Cohesion: 0.11
Nodes (40): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+32 more)

### Community 41 - "RunMigrations"
Cohesion: 0.11
Nodes (36): main(), OpenDB(), newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns() (+28 more)

### Community 42 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 43 - "go_pkg_os"
Cohesion: 0.08
Nodes (26): TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults(), TestWithinTransactionRollsBack(), go_pkg_bufio, go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_flag (+18 more)

### Community 44 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 45 - "GetTypeSchema"
Cohesion: 0.09
Nodes (35): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion() (+27 more)

### Community 46 - "response.go"
Cohesion: 0.08
Nodes (25): Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes(), TestCursorRoundTrip() (+17 more)

### Community 47 - "NewStore"
Cohesion: 0.11
Nodes (36): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+28 more)

### Community 48 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 49 - "DecodeKey"
Cohesion: 0.11
Nodes (25): ProviderConfig, Handler, primaryProviderConfig(), Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN() (+17 more)

### Community 50 - "ExportToFile"
Cohesion: 0.12
Nodes (36): ExportToFile(), ImportFromFile(), newTestStore(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions(), AppVersion() (+28 more)

### Community 51 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 52 - "adguardhome_test.go"
Cohesion: 0.10
Nodes (34): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+26 more)

### Community 53 - "NewChecker"
Cohesion: 0.21
Nodes (33): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+25 more)

### Community 54 - "routerDeps"
Cohesion: 0.11
Nodes (29): routerDeps, PermissionChecker, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountChatRoutes() (+21 more)

### Community 55 - "traefik/tables.go"
Cohesion: 0.14
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 56 - "traefik_test.go"
Cohesion: 0.11
Nodes (32): AllowLoopbackForTest(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath() (+24 more)

### Community 57 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 58 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 59 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 60 - "AuthedUser"
Cohesion: 0.11
Nodes (31): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+23 more)

### Community 61 - "Get"
Cohesion: 0.11
Nodes (15): Handler, isWritableField(), validateConfigPushRequest(), capitalize(), Handler, elevationFailureReason(), ValidateElevationHeader(), WriteElevationError() (+7 more)

### Community 62 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 63 - "ErrorWithDetails"
Cohesion: 0.13
Nodes (12): updateUserRequest, newToken(), Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, validTargetType() (+4 more)

### Community 64 - "backup/backup.go"
Cohesion: 0.15
Nodes (29): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+21 more)

### Community 65 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (13): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Tx (+5 more)

### Community 66 - "middleware.go"
Cohesion: 0.11
Nodes (15): contextKey, elevationError, go_pkg_bytes, go_pkg_crypto_hmac, go_pkg_crypto_sha256, go_pkg_encoding_hex, go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter (+7 more)

### Community 67 - "testApp"
Cohesion: 0.11
Nodes (19): mustHashDummyPassword(), instanceAdminRoleFor(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable() (+11 more)

### Community 68 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 69 - "SystemPage.tsx"
Cohesion: 0.09
Nodes (20): web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey, web_src_api_generated_system_system_getsystembackupschedule, web_src_api_generated_system_system_postsystembackuprun, web_src_api_generated_system_system_putsystembackupschedule, web_src_api_generated_system_system_usegethealth, web_src_api_generated_system_system_usegetsystembackupruns, web_src_api_generated_system_system_usegetsysteminfo (+12 more)

### Community 70 - "NewEngine"
Cohesion: 0.18
Nodes (23): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+15 more)

### Community 71 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 72 - "AuditRecord"
Cohesion: 0.15
Nodes (10): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), Store, scanConnectorGrants(), database/sql.Rows (+2 more)

### Community 73 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 74 - "compliance/engine.go"
Cohesion: 0.15
Nodes (21): catalog(), contains(), equal(), findAttribute(), Catalog, Condition, Entity, Rule (+13 more)

### Community 75 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 76 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 77 - "NewRegistry"
Cohesion: 0.22
Nodes (21): SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders(), TestSuggestWithFallbackStopsOnNonRetryableError() (+13 more)

### Community 78 - "Connector"
Cohesion: 0.15
Nodes (9): apiMessage(), controllerName(), countByKind(), isTimeout(), statusError(), unavailable(), Connector, sectionFetch (+1 more)

### Community 79 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 80 - "newTestHandler"
Cohesion: 0.21
Nodes (18): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+10 more)

### Community 81 - "New"
Cohesion: 0.12
Nodes (22): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+14 more)

### Community 82 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 83 - "logging.go"
Cohesion: 0.15
Nodes (18): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestLoggablePathMasksShareTokenUnderV1() (+10 more)

### Community 84 - "SnapshotEntity"
Cohesion: 0.11
Nodes (21): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable() (+13 more)

### Community 85 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 86 - "useRole.ts"
Cohesion: 0.15
Nodes (16): web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, RoleGate(), RoleGateProps, DocHistory(), LinkedDocPanel(), SETTINGS_SECTIONS, SettingsLayout() (+8 more)

### Community 87 - "AuthMiddleware"
Cohesion: 0.12
Nodes (15): APIKeyChecker, testAPIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), APIKeyClaims (+7 more)

### Community 88 - "Service"
Cohesion: 0.16
Nodes (11): Claims, ElevationClaims, ElevationToken, TokenPair, testHandler, Handler, Service, hasAudience() (+3 more)

### Community 89 - "src/theme.ts"
Cohesion: 0.10
Nodes (16): @fontsource/ibm-plex-mono, @fontsource/ibm-plex-sans, @fontsource/space-mono, @fontsource-variable/big-shoulders-text, @fontsource-variable/geist, @fontsource-variable/geist-mono, @fontsource-variable/inter-tight, @fontsource-variable/jetbrains-mono (+8 more)

### Community 90 - "chat/chat.go"
Cohesion: 0.14
Nodes (18): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections() (+10 more)

### Community 91 - "router.go"
Cohesion: 0.18
Nodes (18): go_pkg_github_com_go_chi_chi_v5, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+10 more)

### Community 92 - "home_assistant_test.go"
Cohesion: 0.22
Nodes (19): Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities(), TestFetchConfigIsRequestedOnce() (+11 more)

### Community 93 - "middleware_test.go"
Cohesion: 0.15
Nodes (16): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser() (+8 more)

### Community 94 - "Register"
Cohesion: 0.19
Nodes (19): init(), init(), init(), init(), init(), init(), init(), init() (+11 more)

### Community 95 - "NewRouter"
Cohesion: 0.19
Nodes (15): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz(), chi.Router (+7 more)

### Community 96 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 97 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 98 - "GuardedDialer"
Cohesion: 0.14
Nodes (13): newConnector(), Connector, GuardedDialer(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector() (+5 more)

### Community 99 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 100 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 101 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 102 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 103 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 104 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 105 - "Handler"
Cohesion: 0.18
Nodes (9): applyConnectorScalarUpdates(), configRequestField(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection(), FieldError (+1 more)

### Community 106 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 107 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 108 - "AppearancePage.tsx"
Cohesion: 0.18
Nodes (15): zustand, AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+7 more)

### Community 109 - "net/http.Handler"
Cohesion: 0.15
Nodes (12): AuditRecorder, ConnectorRoleChecker, SecurityHeaders(), TestSecurityHeaders(), chi.Router, mountConnectorRoutes(), recordElevationAudit(), RequireConnectorRole() (+4 more)

### Community 110 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 111 - "main"
Cohesion: 0.17
Nodes (16): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), main(), newLogger(), runHealthcheck(), splitOrigins() (+8 more)

### Community 112 - "handlers_contract_test.go"
Cohesion: 0.25
Nodes (15): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+7 more)

### Community 113 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 114 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 115 - "data.go"
Cohesion: 0.24
Nodes (15): ChangeEntry, ComplianceSection, ConnectorDrift, DefinitionSummary, DocChangeEntry, DocsSection, DriftSection, FindingSummary (+7 more)

### Community 116 - "log/slog.Logger"
Cohesion: 0.26
Nodes (12): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+4 more)

### Community 117 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 118 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 119 - "render_test.go"
Cohesion: 0.30
Nodes (14): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+6 more)

### Community 120 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 121 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (12): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate(), resolveToken() (+4 more)

### Community 122 - "config_cmd_test.go"
Cohesion: 0.21
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 123 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 124 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 125 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 126 - "pfsense.go"
Cohesion: 0.21
Nodes (11): TestAttributeCatalogCoversEmittedKeys(), TestBuildInterfaceTableAttributes(), TestBuildRuleTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildRuleTable(), buildSystemContent(), primaryGatewayName() (+3 more)

### Community 127 - "store/theme.ts"
Cohesion: 0.25
Nodes (13): ColorMode, commit(), Persisted, PRESETS_FONTS, ThemeState, tokensFor(), ACTIVE, applyTokens() (+5 more)

### Community 129 - "changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 131 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 132 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.33
Nodes (7): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture, net/http/httptest.ResponseRecorder

### Community 133 - "NewService"
Cohesion: 0.30
Nodes (11): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+3 more)

### Community 134 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 135 - "NotificationRecord"
Cohesion: 0.29
Nodes (4): Dispatcher, NotificationRecord, Store, scanNotification()

### Community 136 - "store/backup_test.go"
Cohesion: 0.32
Nodes (11): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+3 more)

### Community 137 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 138 - "transform_test.go"
Cohesion: 0.24
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 139 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 140 - "WiseLabz Connector Guide"
Cohesion: 0.17
Nodes (12): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+4 more)

### Community 141 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 142 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 143 - "Store"
Cohesion: 0.24
Nodes (3): sanitize(), APIKey, Store

### Community 144 - "compliance/engine_test.go"
Cohesion: 0.29
Nodes (10): Evaluate(), testCatalog(), TestEvaluateAndKindAndOrder(), TestEvaluateEdgeCases(), TestEvaluateLargeSnapshot(), TestEvaluateOperators(), TestValidate(), TestValidateFieldAttribution() (+2 more)

### Community 145 - "channels.go"
Cohesion: 0.25
Nodes (9): buildEmailMessage(), sendSMTPChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendSMTPChannel_MissingConfig(), TestSplitRecipients(), redactURLError(), go_pkg_net_smtp (+1 more)

### Community 146 - "Engine"
Cohesion: 0.20
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 147 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 148 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 149 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 150 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 151 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 152 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 153 - "ThemeControls.tsx"
Cohesion: 0.24
Nodes (7): AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls(), makePalette()

### Community 154 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 155 - "compliance_rules_test.go"
Cohesion: 0.33
Nodes (8): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), go_pkg_regexp, go_pkg_slices

### Community 156 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 157 - "useTheme"
Cohesion: 0.33
Nodes (7): mermaid, cssVar(), Mermaid(), resolveColor(), load(), useTheme, presetOpts()

### Community 158 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 159 - "decodeBulkRequest"
Cohesion: 0.36
Nodes (3): decodeBulkRequest(), Handler, bulkRequest

### Community 160 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 161 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 162 - "time.Time"
Cohesion: 0.29
Nodes (6): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, userStatus

### Community 163 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 165 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 166 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 167 - "settings.ts"
Cohesion: 0.46
Nodes (6): MotionProvider(), apply(), framerReducedMotion(), seed(), SettingsState, useSettings

### Community 169 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 170 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 171 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 172 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 173 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 174 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 175 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 177 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 178 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 179 - "auth/handlers_test.go"
Cohesion: 0.40
Nodes (4): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), emailDomainAllowed(), oidcRoleForGroups()

### Community 180 - "CORS"
Cohesion: 0.60
Nodes (4): CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders()

### Community 181 - "routerOperations"
Cohesion: 0.50
Nodes (5): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes

### Community 182 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 183 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 184 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 187 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **540 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+535 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1182 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `UserIDFromContext()` connect `net/http.Request` to `Handler`, `context.Context`, `middleware.go`, `decodeBulkRequest`, `AuditRecord`, `Handler`, `net/http.Handler`, `response.go`, `Store`, `net/http.ResponseWriter`, `Get`, `routerDeps`, `middleware_test.go`, `ErrorWithDetails`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `Handler`, `context.Context`, `net/http/httptest.ResponseRecorder`, `store/backup_test.go`, `go_pkg_context`, `net/http.Request`, `Runner`, `Engine`, `net/http.ResponseWriter`, `ConnectorRecord`, `export_test.go`, `dispatcher_test.go`, `NewEngine`, `time.Time`, `share_links_test.go`, `RunMigrations`, `response.go`, `NewStore`, `engine_maintenance_test.go`, `ExportToFile`, `NewChecker`, `Dispatcher`, `AuthedUser`, `ErrorWithDetails`, `backup/backup.go`, `rewritePlaceholders`, `testApp`, `NewEngine`, `Checker`, `NewRegistry`, `New`, `Service`, `chat/chat.go`, `NewRouter`, `Handler`, `diagnostics/diagnostics.go`, `Handler`, `main`, `log/slog.Logger`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `queryPlan()` connect `keyset_test.go` to `context.Context`, `testing.T`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _540 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.02559167379526661 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.0262839635572039 - nodes in this community are weakly interconnected._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.026053639846743294 - nodes in this community are weakly interconnected._