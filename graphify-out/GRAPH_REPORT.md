# Graph Report - wiselabz-pr387  (2026-09-24)

## Corpus Check
- 799 files · ~480,539 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 5915 nodes · 18584 edges · 204 communities (187 shown, 17 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1481 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e3dcb886`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- react
- newDocTestStore
- testing.T
- context.Context
- package.json
- App.tsx
- dispatcher_test.go
- Button.tsx
- net/http.ResponseWriter
- DashboardPage.tsx
- net/http.Client
- go_pkg_context
- icons.tsx
- go_pkg_testing
- go_pkg_strings
- states.tsx
- go_pkg_github_com_wiselabz_wiselabz_internal_store
- web_src_api_model_index
- Store
- Config
- diagnostics/diagnostics.go
- ServiceSnapshot
- SnapshotEntity
- Store
- UsersPage.tsx
- net/http.Request
- NewEngine
- rowScanner
- DocRecord
- cn
- newTestHandler
- fixtures.ts
- WiseLabz — Architecture & Technical Decisions
- ServiceDetailPage.tsx
- ErrorWithDetails
- IsSecureRequest
- NewMalformedResponseError
- home_assistant/tables.go
- share_links_test.go
- docker_test.go
- NewChecker
- dependencies
- SuggestWithFallback
- response.go
- portainer/tables.go
- portainer_test.go
- adguardhome/tables.go
- WebSocketProvider.tsx
- routerDeps
- compliance/engine.go
- Connector
- traefik/tables.go
- devDependencies
- Dispatcher
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- lists.go
- config_test.go
- git.go
- settings.mock.ts
- rewritePlaceholders
- home_assistant_test.go
- AuthMiddleware
- go_pkg_os
- NewStore
- ExportToFile
- Connector
- NewEngine
- newTestHandler
- GetTypeSchema
- handlers.ts
- unifi_test.go
- router.go
- AppearancePage.tsx
- main
- ConnectorRecord
- Checker
- WiseLabz — Design Contract
- .Fetch
- Manager
- newTestHandler
- Compare
- logging.go
- Connector
- gitTarget
- timeline.ts
- migrations.go
- NewHTTPClient
- adguardhome_test.go
- scheduler/health_test.go
- RunMigrations
- VerifyBundleFile
- export_test.go
- channels.go
- time.Time
- ReportsPage.tsx
- Service
- Handler
- Handler
- truenas_test.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- net/http.Handler
- Register
- NewService
- Connector
- traefik_test.go
- migrations_test.go
- docdiffmodel.ts
- notifications/handlers_test.go
- chat/chat.go
- ComplianceRuleRecord
- system/handlers_test.go
- all.go
- Store
- main.tsx
- compilerOptions
- middleware.go
- handlers_contract_test.go
- handlers_actions_test.go
- Handler
- vectorCache
- data.go
- apikey_scope.go
- changes/handlers_test.go
- pagination_contract_test.go
- registry.go
- gitFixture
- render_test.go
- New
- templates.fixtures.ts
- Handler
- changes_test.go
- connectors_health_test.go
- dialSSHStdio
- retention/retention_test.go
- Contributing to WiseLabz
- NewRegistry
- Store
- Decision
- runRestore
- .resolveShareLinkNode
- Store
- keyset_test.go
- Connector
- Decision
- WiseLabz Connector Guide
- Product
- provider_test.go
- .call
- connectors_maintenance_test.go
- NewHandler
- docs_test.go
- .call
- ReportData
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- handlers_bulk_test.go
- release-please-config.json
- dashboard/handlers_test.go
- RateLimit
- ComputeWindow
- Store
- DeliveryRecord
- Handler
- Step by step
- WiseLabz
- TestComplianceRuleValidation
- runbooks/handlers_test.go
- .applyChannelSecrets
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- walkCursorPages
- .UpdateAuthConfig
- .Fetch
- Store
- engine_maintenance_test.go
- Security Policy
- RequireConnectorRole
- auth/handlers_test.go
- CORS
- routerOperations
- schema.go
- seedScopeFixture
- Enforcement Guidelines
- compose-smoke.sh
- timeoutError
- MISSING — deferred & future frontend features
- stubEmbedder
- WiseLabz — v2 Backlog
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- AlertNotifier
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

## Communities (204 total, 17 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (174): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+166 more)

### Community 1 - "react"
Cohesion: 0.02
Nodes (134): react, sonner, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid (+126 more)

### Community 2 - "newDocTestStore"
Cohesion: 0.02
Nodes (137): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+129 more)

### Community 3 - "testing.T"
Cohesion: 0.02
Nodes (128): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+120 more)

### Community 4 - "context.Context"
Cohesion: 0.03
Nodes (34): fakeStatusChecker, sanitizeSessions(), Connector, Connector, existingIDs(), Store, SnapshotRecord, Store (+26 more)

### Community 5 - "package.json"
Cohesion: 0.03
Nodes (94): clsx, codemirror, @codemirror/commands, @codemirror/state, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker (+86 more)

### Community 6 - "App.tsx"
Cohesion: 0.03
Nodes (78): @codemirror/lang-markdown, @codemirror/view, @uiw/react-codemirror, AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), RefreshFn (+70 more)

### Community 7 - "dispatcher_test.go"
Cohesion: 0.06
Nodes (77): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), testLogger(), expireAlertsOnce() (+69 more)

### Community 8 - "Button.tsx"
Cohesion: 0.04
Nodes (72): RFC-3339, motion, @radix-ui/react-popover, react-i18next, @tanstack/react-query, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow (+64 more)

### Community 9 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (31): Handler, newToken(), sanitize(), Handler, diffToSpec(), Handler, Handler, Handler (+23 more)

### Community 10 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (82): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+74 more)

### Community 11 - "net/http.Client"
Cohesion: 0.04
Nodes (27): Connector, ollamaEmbedder, openAIEmbedder, NewServiceUnavailableError(), setHeaders(), tryParseEntities(), validateCustomURL(), Connector (+19 more)

### Community 12 - "go_pkg_context"
Cohesion: 0.08
Nodes (17): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors (+9 more)

### Community 13 - "icons.tsx"
Cohesion: 0.05
Nodes (65): Frontend shell & theme (decided 2026-06), react-router-dom, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_postdocstopology (+57 more)

### Community 14 - "go_pkg_testing"
Cohesion: 0.05
Nodes (21): TestLoggablePathMasksShareTokenUnderV1(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), IsTimeout(), TestIsTimeout() (+13 more)

### Community 15 - "go_pkg_strings"
Cohesion: 0.06
Nodes (31): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+23 more)

### Community 16 - "states.tsx"
Cohesion: 0.05
Nodes (67): Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve (+59 more)

### Community 17 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.06
Nodes (36): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection (+28 more)

### Community 18 - "web_src_api_model_index"
Cohesion: 0.04
Nodes (48): i18next, msw, @testing-library/jest-dom, @testing-library/react, vitest, web_src_api_generated_runbooks_runbooks, web_src_api_generated_runbooks_runbooks_usegetrunbooks, web_src_api_generated_runbooks_runbooks_usegetrunbooksrunbookid (+40 more)

### Community 19 - "Store"
Cohesion: 0.04
Nodes (39): Config, Handler, TestEmbeddedSPAWithoutFrontendBuild(), Embedder, EmbedRegistry, Registry, NewHandler(), NewHandler() (+31 more)

### Community 20 - "Config"
Cohesion: 0.05
Nodes (47): healthFakeConnector, Config, LogSettings, IsSSHRemote(), idempotent(), retryable(), RetryTransport(), sleep() (+39 more)

### Community 21 - "diagnostics/diagnostics.go"
Cohesion: 0.05
Nodes (56): ProviderConfig, Handler, primaryProviderConfig(), Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN() (+48 more)

### Community 22 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (19): noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, init(), RegisterTransformer(), runTransformers(), TestRunTransformersAppliesInOrderAndStopsOnError() (+11 more)

### Community 23 - "SnapshotEntity"
Cohesion: 0.09
Nodes (54): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable() (+46 more)

### Community 24 - "Store"
Cohesion: 0.06
Nodes (15): Sanitize(), TestSanitize(), placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert() (+7 more)

### Community 25 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (47): axios, customInstance(), web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetpassword (+39 more)

### Community 26 - "net/http.Request"
Cohesion: 0.08
Nodes (22): Handler, isWritableField(), validateConfigPushRequest(), applyConnectorScalarUpdates(), Handler, capitalize(), Handler, decodeBulkRequest() (+14 more)

### Community 27 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 28 - "rowScanner"
Cohesion: 0.07
Nodes (25): Dispatcher, actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), scanChange(), NotificationRecord (+17 more)

### Community 29 - "DocRecord"
Cohesion: 0.06
Nodes (29): docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, importBundle(), importDocVersions(), importTemplates() (+21 more)

### Community 30 - "cn"
Cohesion: 0.06
Nodes (40): web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateidversions (+32 more)

### Community 31 - "newTestHandler"
Cohesion: 0.06
Nodes (39): mustHashDummyPassword(), testHandler, Handler, instanceAdminRoleFor(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), spaHandler() (+31 more)

### Community 32 - "fixtures.ts"
Cohesion: 0.06
Nodes (41): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+33 more)

### Community 33 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 34 - "ServiceDetailPage.tsx"
Cohesion: 0.06
Nodes (39): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+31 more)

### Community 35 - "ErrorWithDetails"
Cohesion: 0.09
Nodes (21): updateUserRequest, Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, configRequestField(), parseScheduleUpdates() (+13 more)

### Community 36 - "IsSecureRequest"
Cohesion: 0.08
Nodes (27): clearOIDCFlowCookie(), Handler, newOIDCUser(), oidcFlowCookieName(), randomOIDCToken(), readOIDCFlowCookie(), setOIDCFlowCookie(), validHostPort() (+19 more)

### Community 37 - "NewMalformedResponseError"
Cohesion: 0.06
Nodes (25): TimeoutError, GuardedDialer(), IsDangerousIP(), NewAuthError(), NewMalformedResponseError(), NewTimeoutError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap() (+17 more)

### Community 38 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 39 - "share_links_test.go"
Cohesion: 0.17
Nodes (41): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld() (+33 more)

### Community 40 - "docker_test.go"
Cohesion: 0.07
Nodes (39): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+31 more)

### Community 41 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 42 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 43 - "SuggestWithFallback"
Cohesion: 0.09
Nodes (18): claudeProvider, openAICompatibleProvider, Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable() (+10 more)

### Community 44 - "response.go"
Cohesion: 0.08
Nodes (26): NewHandler(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes() (+18 more)

### Community 45 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 46 - "portainer_test.go"
Cohesion: 0.10
Nodes (35): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+27 more)

### Community 47 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 48 - "WebSocketProvider.tsx"
Cohesion: 0.08
Nodes (29): Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`) (+21 more)

### Community 49 - "routerDeps"
Cohesion: 0.12
Nodes (28): routerDeps, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+20 more)

### Community 50 - "compliance/engine.go"
Cohesion: 0.11
Nodes (31): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+23 more)

### Community 51 - "Connector"
Cohesion: 0.08
Nodes (10): init(), ConfigField, Connector, buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment(), ValidateRefSegment() (+2 more)

### Community 52 - "traefik/tables.go"
Cohesion: 0.14
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 53 - "devDependencies"
Cohesion: 0.06
Nodes (35): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+27 more)

### Community 54 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 55 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 56 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 57 - "lists.go"
Cohesion: 0.14
Nodes (30): buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), parseAdlistsV6(), parseClientsV6() (+22 more)

### Community 58 - "config_test.go"
Cohesion: 0.09
Nodes (29): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), runHealthcheck(), Load() (+21 more)

### Community 59 - "git.go"
Cohesion: 0.08
Nodes (25): fileName(), slugify(), TestCommitMessage(), keys(), go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter (+17 more)

### Community 60 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 61 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 62 - "home_assistant_test.go"
Cohesion: 0.14
Nodes (28): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+20 more)

### Community 63 - "AuthMiddleware"
Cohesion: 0.11
Nodes (23): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, AuthMiddleware(), extractBearerToken(), hashToken(), RequireInstanceAdmin() (+15 more)

### Community 64 - "go_pkg_os"
Cohesion: 0.11
Nodes (15): main(), usage(), go_pkg_flag, go_pkg_github_com_robfig_cron_v3, go_pkg_github_com_wiselabz_wiselabz_internal_backup, go_pkg_github_com_wiselabz_wiselabz_internal_config, go_pkg_github_com_wiselabz_wiselabz_internal_connector_all, go_pkg_github_com_wiselabz_wiselabz_internal_diagnostics (+7 more)

### Community 65 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 66 - "ExportToFile"
Cohesion: 0.13
Nodes (28): Export(), ExportToFile(), Import(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile() (+20 more)

### Community 67 - "Connector"
Cohesion: 0.18
Nodes (8): SnapshotSection, groupNames(), parseGroupsV6(), Connector, unavailable(), parseGroupsV5(), groupRow, session

### Community 68 - "NewEngine"
Cohesion: 0.16
Nodes (24): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+16 more)

### Community 69 - "newTestHandler"
Cohesion: 0.12
Nodes (26): TestConnectorStoreErrorPaths(), TestDataSnapshotPaths(), TestSyncsLimitAndIsolation(), createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler() (+18 more)

### Community 70 - "GetTypeSchema"
Cohesion: 0.12
Nodes (25): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion() (+17 more)

### Community 71 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 72 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 73 - "router.go"
Cohesion: 0.14
Nodes (22): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+14 more)

### Community 74 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 75 - "main"
Cohesion: 0.13
Nodes (23): main(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), NewEmbedRegistry(), RegisterOpenAICompatible() (+15 more)

### Community 76 - "ConnectorRecord"
Cohesion: 0.14
Nodes (14): connectorIDs(), importConnectors(), ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr() (+6 more)

### Community 77 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 78 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 79 - ".Fetch"
Cohesion: 0.13
Nodes (12): ServiceDependency, WantsField(), TestRequestedFields(), TestWantsField(), environmentDependencies(), putMetadata(), unavailable(), poolDependencies() (+4 more)

### Community 80 - "Manager"
Cohesion: 0.16
Nodes (8): cron.EntryID, Manager, LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store, Scheduler

### Community 81 - "newTestHandler"
Cohesion: 0.21
Nodes (18): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+10 more)

### Community 82 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 83 - "logging.go"
Cohesion: 0.15
Nodes (17): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+9 more)

### Community 84 - "Connector"
Cohesion: 0.16
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

### Community 85 - "gitTarget"
Cohesion: 0.13
Nodes (14): commitMessage(), gitAuth(), Exporter, installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), commitResult, GitOptions (+6 more)

### Community 86 - "timeline.ts"
Cohesion: 0.15
Nodes (14): installMockWebSocket(), Listenerish, MockWebSocket, Emit, env(), heartbeat(), newId(), ScheduledEvent (+6 more)

### Community 87 - "migrations.go"
Cohesion: 0.11
Nodes (17): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), runSQLiteWithForeignKeysOff(), TestGetMigrationStatus() (+9 more)

### Community 88 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 89 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 90 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 91 - "RunMigrations"
Cohesion: 0.14
Nodes (19): Handler, newScratchStore(), SyncDocEmbeddings(), TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(), TestRetrieveUsesCacheAndSyncInvalidates(), TestUpsertBackupSchedulePostgresParity(), Store, newPostgresTestStore() (+11 more)

### Community 92 - "VerifyBundleFile"
Cohesion: 0.18
Nodes (19): TestValidateBundleRejectsOrphanDocVersion(), TestValidateBundleRejectsWrongVersion(), ValidateBundle(), failVerification(), LatestBundle(), RunVerifyOnce(), ListVerifications(), RecordVerification() (+11 more)

### Community 93 - "export_test.go"
Cohesion: 0.20
Nodes (17): fetchAllDocs(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile() (+9 more)

### Community 94 - "channels.go"
Cohesion: 0.17
Nodes (17): buildEmailMessage(), sendNtfyChannel(), sendSMTPChannel(), sendTelegramChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendNtfyChannel_DefaultsToNtfySh(), TestSendNtfyChannel_MissingTopic() (+9 more)

### Community 95 - "time.Time"
Cohesion: 0.13
Nodes (11): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session, io.WriteCloser (+3 more)

### Community 96 - "ReportsPage.tsx"
Cohesion: 0.12
Nodes (18): web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports, web_src_api_generated_reports_reports_usegetreportsdefinitions (+10 more)

### Community 97 - "Service"
Cohesion: 0.20
Nodes (9): Claims, ElevationClaims, ElevationToken, TokenPair, Service, hasAudience(), newTokenID(), jwt.ClaimStrings (+1 more)

### Community 98 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 99 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 100 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 101 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 102 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 103 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 104 - "net/http.Handler"
Cohesion: 0.19
Nodes (16): NewHandler(), TestCreate(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), Handler (+8 more)

### Community 105 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 106 - "NewService"
Cohesion: 0.21
Nodes (15): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+7 more)

### Community 107 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 108 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 109 - "migrations_test.go"
Cohesion: 0.30
Nodes (16): collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns(), TestMigrationSchemaParity(), RunMigrationsDown(), assertGoldenSnapshotsExists(), assertKeysetPaginationIndexes(), assertNtfyTelegramChannelsAllowed() (+8 more)

### Community 110 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 111 - "notifications/handlers_test.go"
Cohesion: 0.28
Nodes (15): AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+7 more)

### Community 112 - "chat/chat.go"
Cohesion: 0.16
Nodes (14): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+6 more)

### Community 113 - "ComplianceRuleRecord"
Cohesion: 0.19
Nodes (7): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule(), Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 114 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 115 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 116 - "Store"
Cohesion: 0.16
Nodes (7): Store, ChatConversationRecord, Store, apiKeyConnectorFilter(), AttentionItem, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 117 - "main.tsx"
Cohesion: 0.17
Nodes (11): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers (+3 more)

### Community 118 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 119 - "middleware.go"
Cohesion: 0.22
Nodes (10): AuditRecorder, contextKey, elevationError, PermissionChecker, UserStatusChecker, elevationFailureReason(), recordElevationAudit(), RequireElevation() (+2 more)

### Community 120 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 121 - "handlers_actions_test.go"
Cohesion: 0.35
Nodes (14): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+6 more)

### Community 122 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 123 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 124 - "data.go"
Cohesion: 0.27
Nodes (14): ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary, JobHealthEntry (+6 more)

### Community 125 - "apikey_scope.go"
Cohesion: 0.21
Nodes (10): APIKeyRestriction, testAPIKeyChecker, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), isSafeMethod(), TestClampConnectorRole(), APIKeyClaims (+2 more)

### Community 126 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 127 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 128 - "registry.go"
Cohesion: 0.19
Nodes (10): IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType(), TestRegisterDefaultsToNonStub(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError, Factory (+2 more)

### Community 129 - "gitFixture"
Cohesion: 0.32
Nodes (8): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, Exporter, github.com/go-git/go-git/v5/plumbing/object.Commit

### Community 130 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 131 - "New"
Cohesion: 0.36
Nodes (13): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+5 more)

### Community 132 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (11): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, resolveToken(), Snapshot (+3 more)

### Community 134 - "changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 135 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 136 - "dialSSHStdio"
Cohesion: 0.15
Nodes (11): serveOneHTTPExchange(), serveSSHDockerConn(), TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ClientConfig (+3 more)

### Community 137 - "retention/retention_test.go"
Cohesion: 0.35
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 138 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 139 - "NewRegistry"
Cohesion: 0.50
Nodes (11): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+3 more)

### Community 140 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 141 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 142 - "runRestore"
Cohesion: 0.24
Nodes (11): confirm(), formatCounts(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag() (+3 more)

### Community 143 - ".resolveShareLinkNode"
Cohesion: 0.25
Nodes (5): contextWithShareLink(), Handler, shareLinkFromContext(), shareLinkNode, shareLinkScope

### Community 144 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 145 - "keyset_test.go"
Cohesion: 0.33
Nodes (10): assertSameSet(), Store, T, queryPlan(), TestKeysetQueriesUseCoveringIndexes(), TestListAuditRecordsKeysetHonoursFilters(), TestListAuditRecordsKeysetTraversal(), TestListChangesKeysetTraversal() (+2 more)

### Community 147 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 148 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 149 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 150 - "provider_test.go"
Cohesion: 0.29
Nodes (6): testProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream()

### Community 151 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 152 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 153 - "NewHandler"
Cohesion: 0.24
Nodes (10): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+2 more)

### Community 154 - "docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 155 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 156 - "ReportData"
Cohesion: 0.47
Nodes (4): connectorFilter(), DefinitionSummary, Generator, ReportData

### Community 157 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 158 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 159 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 160 - "handlers_bulk_test.go"
Cohesion: 0.56
Nodes (8): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync()

### Community 161 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 162 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 163 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 164 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 165 - "Store"
Cohesion: 0.32
Nodes (3): Store, scanBackupRun(), BackupRun

### Community 166 - "DeliveryRecord"
Cohesion: 0.39
Nodes (4): DeliveryRecord, DeliveryStatus, Store, scanDelivery()

### Community 168 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 169 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 170 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 171 - "runbooks/handlers_test.go"
Cohesion: 0.48
Nodes (6): Handler, newTestHandler(), TestCreate(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete()

### Community 173 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 174 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 175 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 176 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 177 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 178 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 180 - ".Fetch"
Cohesion: 0.47
Nodes (3): TestBuildContainerTableAttributes(), buildContainerTable(), Connector

### Community 182 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 183 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 184 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 185 - "auth/handlers_test.go"
Cohesion: 0.40
Nodes (4): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), emailDomainAllowed(), oidcRoleForGroups()

### Community 186 - "CORS"
Cohesion: 0.60
Nodes (4): CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders()

### Community 187 - "routerOperations"
Cohesion: 0.50
Nodes (5): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes

### Community 188 - "schema.go"
Cohesion: 0.50
Nodes (4): Schema(), schemaFor(), TestSchemaMatchesConfig(), reflect.Type

### Community 189 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 190 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 191 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 194 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **542 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+537 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1200 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Connector` connect `net/http.Client` to `go_pkg_context`, `home_assistant/tables.go`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `gitFixture` to `testing.T`, `context.Context`, `dispatcher_test.go`, `Store`, `git.go`, `export_test.go`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Why does `createTestConnector()` connect `newDocTestStore` to `keyset_test.go`, `testing.T`, `context.Context`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _542 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.02104903857094095 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.0213857998289136 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.024337805297557618 - nodes in this community are weakly interconnected._