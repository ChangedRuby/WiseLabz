# Graph Report - deps-dev  (2026-10-01)

## Corpus Check
- 880 files · ~557,022 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6615 nodes · 20979 edges · 225 communities (208 shown, 17 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1680 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `95d13374`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- testing.T
- newDocTestStore
- context.Context
- cn
- react
- go_pkg_strings
- ServiceDetailPage.tsx
- connector/connector.go
- icons.tsx
- DashboardPage.tsx
- @tanstack/react-query
- go_pkg_context
- Errorf
- newTestHandler
- dispatcher_test.go
- ProfilePage.tsx
- go_pkg_net_http
- ServiceSnapshot
- App.tsx
- New
- go_pkg_testing
- NewMalformedResponseError
- UsersPage.tsx
- NewService
- ConnectorRecord
- ErrorWithDetails
- WebSocketProvider.tsx
- HashToken
- SnapshotEntity
- go_pkg_os
- NewEngine
- package.json
- SystemPage.tsx
- GetTypeSchema
- net/http.Request
- fixtures.ts
- server/main.go
- docker_test.go
- NewUser
- Compare
- Store
- Get
- ExportToFile
- response.go
- nilToStr
- dependencies
- home_assistant/tables.go
- Connector
- home_assistant_test.go
- portainer/tables.go
- Dispatcher
- adguardhome/tables.go
- net/http.ResponseWriter
- newRouterDeps
- NewChecker
- traefik/tables.go
- unifi/tables.go
- connector_permission.go
- ws/ws_test.go
- traefik_test.go
- Store
- main
- NewStore
- rewritePlaceholders
- DecodeJSON
- log/slog.Logger
- settings.mock.ts
- SuggestRequest
- Config
- handlers.ts
- runbooks_test.go
- Service
- newTestHandler
- sync.Mutex
- unifi_test.go
- Hub
- timeline.ts
- src/theme.ts
- devDependencies
- Manager
- WiseLabz — Design Contract
- routerDeps
- compliance/engine.go
- Checker
- NewEngine
- New
- portainer_test.go
- Connector
- NewHTTPClient
- adguardhome_test.go
- .Fetch
- .runPerRevision
- router.go
- gitFixture
- Runner
- WiseLabz — Architecture & Technical Decisions
- Configuration & Documentation Backup (Export/Import)
- Registry
- connectors_health_test.go
- Handler
- Connector
- export_test.go
- NewClient
- time.Time
- Backend test performance
- ReportsPage.tsx
- ws.ts
- Handler
- Handler
- sshStdioConn
- truenas_test.go
- diagnostics/diagnostics.go
- Decision
- compilerOptions
- RunbookRecord
- docdiffmodel.ts
- handlers_actions_test.go
- vectorCache
- Register
- all.go
- httpx/retry_test.go
- Deps
- sync/engine_test.go
- compilerOptions
- testApp
- apikey_scopes_test.go
- handlers_contract_test.go
- newTestHarness
- pagination_contract_test.go
- Handler
- newTestHandler
- time.Duration
- render_test.go
- Contributing to WiseLabz
- Backup Recovery: What Comes Back, and What Doesn't
- runRestore
- NewRegistry
- api/changes_test.go
- mountAPIRoutes
- chat/chat.go
- retention/retention_test.go
- ARCHITECTURE.md
- scripts
- config_cmd_test.go
- newTestHandler
- IsSecureRequest
- ReportData
- Store
- JobHealthRecord
- Contributor Covenant Code of Conduct
- Decision
- Decision
- main.tsx
- net/http.Response
- compliance/engine_test.go
- Decision
- WiseLabz Connector Guide
- .call
- cursor_pagination_test.go
- RequireElevation
- RequireConnectorRole
- compliance_rules_test.go
- Engine
- connectors_maintenance_test.go
- NewHandler
- .call
- go_pkg_github_com_go_chi_chi_v5
- Changelog
- Connector
- Connector
- Product
- settings.ts
- changes.sh
- test-shards.sh
- mockServiceWorker.js
- ComplianceRuleRecord
- openapi_contract_test.go
- Elector
- release-please-config.json
- auth/handlers_test.go
- dashboard/handlers_test.go
- RateLimit
- validate.go
- ComputeWindow
- migrations.go
- Cache
- Audit Trail
- Step by step
- WiseLabz
- RequirePermission
- TestBulkReauth
- MigratedSQLite
- newTestHandler
- gitAuth
- scanMaintenanceWindow
- computeNextRun
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- Store
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- GuardedDialer
- Authentication design
- Development workflow
- MfaEnrollDialog
- @vitejs/plugin-react
- compose-smoke.sh
- mountPublicSystemRoutes
- ServiceDependency
- timeoutError
- fields_test.go
- WiseLabz — v2 Backlog
- coverage-parity.sh
- vite-env.d.ts
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- fakeNotifier
- snapshotChangingNotifier
- coverpkg.sh
- github.com/WiseLabz/wiselabz

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

## Communities (225 total, 17 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (172): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+164 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (170): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+162 more)

### Community 2 - "newDocTestStore"
Cohesion: 0.02
Nodes (157): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+149 more)

### Community 3 - "context.Context"
Cohesion: 0.02
Nodes (40): fakeConnectorRoleChecker, fakeStatusChecker, sanitizeSessions(), treatAsSafeFromContext(), APIKeyIDFromContext(), MFAEnrollOnlyFromContext(), Connector, existingIDs() (+32 more)

### Community 4 - "cn"
Cohesion: 0.02
Nodes (121): web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_notifications_notifications, web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey, web_src_api_generated_settings_settings_getgetauthconfigquerykey, web_src_api_generated_settings_settings_getgetnotificationsconfigquerykey, web_src_api_generated_settings_settings_postaiconfigtest (+113 more)

### Community 5 - "react"
Cohesion: 0.03
Nodes (121): Endpoints, Frontend, Saved Views, Scope, match-sorter, motion, @radix-ui/react-popover, react (+113 more)

### Community 6 - "go_pkg_strings"
Cohesion: 0.04
Nodes (43): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+35 more)

### Community 7 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (92): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules (+84 more)

### Community 8 - "connector/connector.go"
Cohesion: 0.03
Nodes (39): Connector, ollamaEmbedder, openAIEmbedder, TimeoutError, NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders() (+31 more)

### Community 9 - "icons.tsx"
Cohesion: 0.04
Nodes (80): Frontend shell & theme (decided 2026-06), i18next, react-error-boundary, react-i18next, sonner, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync (+72 more)

### Community 10 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (85): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+77 more)

### Community 11 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (50): msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage, ConfirmDestructive() (+42 more)

### Community 12 - "go_pkg_context"
Cohesion: 0.07
Nodes (19): StatusError, dashboardLayout, ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), contains(), searchString(), go_pkg_context (+11 more)

### Community 13 - "Errorf"
Cohesion: 0.05
Nodes (32): Handler, updateUserRequest, newToken(), sanitize(), Handler, writeUserWriteError(), Handler, Handler (+24 more)

### Community 14 - "newTestHandler"
Cohesion: 0.06
Nodes (71): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+63 more)

### Community 15 - "dispatcher_test.go"
Cohesion: 0.07
Nodes (69): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext() (+61 more)

### Community 16 - "ProfilePage.tsx"
Cohesion: 0.03
Nodes (69): @simplewebauthn/browser, setAccessToken(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin (+61 more)

### Community 17 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (29): bulkSnoozeItemResult, bulkSnoozeRequest, contextKey, elevationError, versionSections(), TemplateVersionSection, shareLinkContextKey, go_pkg_crypto_rand (+21 more)

### Community 18 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (23): noopValidatedConnector, ServiceSnapshot, changePatternID(), Engine, markError(), snapshotIDOrNil(), init(), RegisterTransformer() (+15 more)

### Community 19 - "App.tsx"
Cohesion: 0.04
Nodes (64): setMfaEnrollmentRequiredHandler(), web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest, web_src_api_generated_docs_docs_postdocsdocidlock, web_src_api_generated_docs_docs_postdocsdocidlockrelease (+56 more)

### Community 20 - "New"
Cohesion: 0.05
Nodes (70): main(), New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists() (+62 more)

### Community 21 - "go_pkg_testing"
Cohesion: 0.04
Nodes (23): TestSchemaMatchesConfig(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), IsTimeout(), TestIsTimeout() (+15 more)

### Community 22 - "NewMalformedResponseError"
Cohesion: 0.07
Nodes (44): SnapshotSection, NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+36 more)

### Community 23 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (56): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+48 more)

### Community 24 - "NewService"
Cohesion: 0.05
Nodes (64): APIKeyChecker, testAuditCall, testAuditRecorder, UserStatusChecker, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders() (+56 more)

### Community 25 - "ConnectorRecord"
Cohesion: 0.05
Nodes (37): enableFakeEmbedding(), actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), ConnectorRecord, Store (+29 more)

### Community 26 - "ErrorWithDetails"
Cohesion: 0.06
Nodes (31): sanitizeUser(), setRefreshCookie(), Handler, mustHashDummyPassword(), Handler, Handler, newOIDCUser(), randomOIDCToken() (+23 more)

### Community 27 - "WebSocketProvider.tsx"
Cohesion: 0.05
Nodes (49): RFC-3339, 15. `system.resync`, Client dispatch model, Delivery, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior (+41 more)

### Community 28 - "HashToken"
Cohesion: 0.07
Nodes (35): ProviderConfig, factorJSON(), Handler, testHandler, Handler, Handler, Handler, primaryProviderConfig() (+27 more)

### Community 29 - "SnapshotEntity"
Cohesion: 0.09
Nodes (54): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildPeerTableAttributes() (+46 more)

### Community 30 - "go_pkg_os"
Cohesion: 0.06
Nodes (35): keys(), writeExportState(), exportCursor, exportState, go_pkg_bufio, go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_flag (+27 more)

### Community 31 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 32 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 33 - "SystemPage.tsx"
Cohesion: 0.05
Nodes (40): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_getreportsreportiddownload, web_src_api_generated_reports_reports_usegetreportsreportid (+32 more)

### Community 34 - "GetTypeSchema"
Cohesion: 0.06
Nodes (46): TestRegisteredSchema(), TestSchemaConfigValidation(), countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), TestFailureContract(), Capabilities() (+38 more)

### Community 35 - "net/http.Request"
Cohesion: 0.08
Nodes (21): Handler, isWritableField(), capitalize(), Handler, decodeBulkRequest(), Handler, Handler, loggablePath() (+13 more)

### Community 36 - "fixtures.ts"
Cohesion: 0.06
Nodes (44): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+36 more)

### Community 37 - "server/main.go"
Cohesion: 0.07
Nodes (30): changePromptData(), stripPromptTags(), truncateUTF8(), buildPrompt(), TestBuildPrompt(), bulkResolveItemResult, bulkResolveRequest, go_pkg_github_com_mark3labs_mcp_go_client (+22 more)

### Community 38 - "docker_test.go"
Cohesion: 0.05
Nodes (48): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+40 more)

### Community 39 - "NewUser"
Cohesion: 0.15
Nodes (46): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestGetRequiresGrantEvenForInstanceAdmin(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON() (+38 more)

### Community 40 - "Compare"
Cohesion: 0.07
Nodes (42): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+34 more)

### Community 41 - "Store"
Cohesion: 0.06
Nodes (17): seedAlert(), changeServiceIDs(), seedChange(), Store, Store, placeholders(), scanBackupRun(), changeFilterClause() (+9 more)

### Community 42 - "Get"
Cohesion: 0.08
Nodes (25): applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), TestDiagnosticsRedactsSecrets(), LifecycleOp(), Get(), Connector, SupportsLifecycleVerb() (+17 more)

### Community 43 - "ExportToFile"
Cohesion: 0.10
Nodes (42): Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+34 more)

### Community 44 - "response.go"
Cohesion: 0.07
Nodes (29): Handler, decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor() (+21 more)

### Community 45 - "nilToStr"
Cohesion: 0.07
Nodes (14): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+6 more)

### Community 46 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 47 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 48 - "Connector"
Cohesion: 0.07
Nodes (11): init(), ConfigField, Connector, buildRouteTable(), buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment() (+3 more)

### Community 49 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (37): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+29 more)

### Community 50 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 51 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 52 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (14): Handler, oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName() (+6 more)

### Community 54 - "newRouterDeps"
Cohesion: 0.10
Nodes (36): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+28 more)

### Community 55 - "NewChecker"
Cohesion: 0.21
Nodes (33): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+25 more)

### Community 56 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 57 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 58 - "connector_permission.go"
Cohesion: 0.11
Nodes (17): APIKeyRestriction, testAPIKeyChecker, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), APIKeyClaims (+9 more)

### Community 59 - "ws/ws_test.go"
Cohesion: 0.13
Nodes (31): TestBroadcastDocEventScoping(), Envelope, newHeartbeat(), NewHub(), normalizeOrigin(), assertEnvelope(), assertNoFrame(), connectorHub() (+23 more)

### Community 60 - "traefik_test.go"
Cohesion: 0.11
Nodes (30): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+22 more)

### Community 61 - "Store"
Cohesion: 0.17
Nodes (27): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+19 more)

### Community 62 - "main"
Cohesion: 0.10
Nodes (25): Config, main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+17 more)

### Community 63 - "NewStore"
Cohesion: 0.15
Nodes (27): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+19 more)

### Community 64 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 65 - "DecodeJSON"
Cohesion: 0.12
Nodes (11): webAuthnFlow, webAuthnUser, Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T (+3 more)

### Community 66 - "log/slog.Logger"
Cohesion: 0.10
Nodes (15): newLogger(), formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep() (+7 more)

### Community 67 - "settings.mock.ts"
Cohesion: 0.09
Nodes (25): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+17 more)

### Community 68 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 69 - "Config"
Cohesion: 0.11
Nodes (21): NewHandler(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+13 more)

### Community 70 - "handlers.ts"
Cohesion: 0.07
Nodes (27): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+19 more)

### Community 71 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 72 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 73 - "newTestHandler"
Cohesion: 0.13
Nodes (25): TestConnectorStoreErrorPaths(), TestDataSnapshotPaths(), createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler() (+17 more)

### Community 74 - "sync.Mutex"
Cohesion: 0.10
Nodes (9): cron.EntryID, Handler, createVersion(), templateResponse(), templateVersionResponse(), sync/atomic.Bool, sync.Mutex, ReadyState (+1 more)

### Community 75 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 76 - "Hub"
Cohesion: 0.13
Nodes (12): Hub, newEnvelope(), setupWSConnection(), github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, Audience, broadcastMsg, Client (+4 more)

### Community 77 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 78 - "src/theme.ts"
Cohesion: 0.15
Nodes (23): @fontsource/space-mono, @fontsource-variable/space-grotesk, ColorMode, commit(), load(), Persisted, PRESETS_FONTS, ThemeState (+15 more)

### Community 79 - "devDependencies"
Cohesion: 0.08
Nodes (25): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+17 more)

### Community 80 - "Manager"
Cohesion: 0.16
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 81 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 82 - "routerDeps"
Cohesion: 0.15
Nodes (8): routerDeps, diffToSpec(), NewHandler(), NewHandler(), Config, Handler, Handler, Handler

### Community 83 - "compliance/engine.go"
Cohesion: 0.16
Nodes (20): catalog(), contains(), equal(), findAttribute(), Catalog, Condition, Entity, Rule (+12 more)

### Community 84 - "Checker"
Cohesion: 0.22
Nodes (7): Snapshot, complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, FindingNotifier, RotationConfig

### Community 85 - "NewEngine"
Cohesion: 0.17
Nodes (19): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+11 more)

### Community 86 - "New"
Cohesion: 0.21
Nodes (20): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires() (+12 more)

### Community 87 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 88 - "Connector"
Cohesion: 0.16
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

### Community 89 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 90 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 91 - ".Fetch"
Cohesion: 0.15
Nodes (8): WantsField(), Connector, putMetadata(), unavailable(), agentEnabled(), Connector, Connector, dockerSectionSpec

### Community 92 - ".runPerRevision"
Cohesion: 0.20
Nodes (11): fileName(), slugify(), commitMessage(), Exporter, TestCommitMessage(), commitResult, gitTarget, git.Repository (+3 more)

### Community 93 - "router.go"
Cohesion: 0.18
Nodes (18): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+10 more)

### Community 94 - "gitFixture"
Cohesion: 0.31
Nodes (11): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), TestGitPerRevisionBootstrapReplayAndCap(), TestGitPerRevisionBotCatchUpAndRejectedPush(), TestGitPerRevisionMissingIntermediateAndEngineAuthor() (+3 more)

### Community 95 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 96 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.10
Nodes (20): ADR index, AI module, API design, Backend, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27) (+12 more)

### Community 97 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.10
Nodes (17): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included, Bundle format (+9 more)

### Community 98 - "Registry"
Cohesion: 0.17
Nodes (12): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+4 more)

### Community 99 - "connectors_health_test.go"
Cohesion: 0.19
Nodes (13): healthFakeConnector, testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline() (+5 more)

### Community 100 - "Handler"
Cohesion: 0.17
Nodes (10): validateConfigPushRequest(), NewHandler(), stepAuditDetail(), validTargetType(), validVerb(), ValidateCompositeRef(), configPushRequest, Handler (+2 more)

### Community 101 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 102 - "export_test.go"
Cohesion: 0.20
Nodes (17): Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile(), TestDocExportDefaultCronExprIsValid() (+9 more)

### Community 103 - "NewClient"
Cohesion: 0.13
Nodes (14): clientTimeout(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects(), TestNewClientTimeout(), TestNewClientUsesCustomDialContext() (+6 more)

### Community 104 - "time.Time"
Cohesion: 0.21
Nodes (18): digestDue(), TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection (+10 more)

### Community 105 - "Backend test performance"
Cohesion: 0.11
Nodes (19): Backend test performance, CI job times, CI measurements, Coverage strategy, Dependency updates and action pinning, Deterministic scheduled exports (#408), Fixture reuse and lifecycle tests (#405–#407), Follow-ups (+11 more)

### Community 106 - "ReportsPage.tsx"
Cohesion: 0.12
Nodes (18): web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports, web_src_api_generated_reports_reports_usegetreportsdefinitions (+10 more)

### Community 107 - "ws.ts"
Cohesion: 0.11
Nodes (18): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+10 more)

### Community 108 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 109 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 110 - "sshStdioConn"
Cohesion: 0.12
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 111 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 112 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 113 - "Decision"
Cohesion: 0.12
Nodes (14): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+6 more)

### Community 114 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 115 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 116 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 117 - "handlers_actions_test.go"
Cohesion: 0.30
Nodes (15): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+7 more)

### Community 118 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 119 - "Register"
Cohesion: 0.23
Nodes (16): init(), init(), init(), init(), init(), init(), init(), init() (+8 more)

### Community 120 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 121 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 122 - "Deps"
Cohesion: 0.31
Nodes (16): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs(), registerListFindings() (+8 more)

### Community 123 - "sync/engine_test.go"
Cohesion: 0.18
Nodes (10): TestRunSync_NoNotifyOnInfoOnlyChange(), TestRunSync_NotifiesOnEligibleAlert(), TestRunSyncBatchRollbackAndRetry(), TestRunSyncDocRegeneratorErrorIsNonFatal(), TestRunSyncInvokesDocRegeneratorAfterQualityChecker(), TestRunSyncInvokesQualityChecker(), TestRunSyncQualityCheckerErrorIsNonFatal(), fakeDocRegenerator (+2 more)

### Community 124 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 125 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 126 - "apikey_scopes_test.go"
Cohesion: 0.23
Nodes (14): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer() (+6 more)

### Community 127 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 128 - "newTestHarness"
Cohesion: 0.27
Nodes (14): TestListAttentionItems(), TestListChanges(), TestListConnectors(), TestSearchDocs(), seedFinding(), TestListFindings(), createConnector(), createUser() (+6 more)

### Community 129 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 130 - "Handler"
Cohesion: 0.24
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 131 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 132 - "time.Duration"
Cohesion: 0.18
Nodes (7): sleep(), Database, HASettings, Server, SyncSettings, time.Duration, PoolConfig

### Community 133 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 134 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 135 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.14
Nodes (12): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order), Backups, PostgreSQL support (+4 more)

### Community 136 - "runRestore"
Cohesion: 0.21
Nodes (13): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+5 more)

### Community 137 - "NewRegistry"
Cohesion: 0.45
Nodes (12): NewRegistry(), TestExplain(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip() (+4 more)

### Community 138 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 139 - "mountAPIRoutes"
Cohesion: 0.28
Nodes (11): chi.Router, mountChatRoutes(), mountDocRoutes(), mountShareRoutes(), mountTemplateRoutes(), chi.Router, mountAPIRoutes(), chi.Router (+3 more)

### Community 140 - "chat/chat.go"
Cohesion: 0.24
Nodes (11): cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips(), TestSplitSections() (+3 more)

### Community 141 - "retention/retention_test.go"
Cohesion: 0.35
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 142 - "ARCHITECTURE.md"
Cohesion: 0.17
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 143 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 144 - "config_cmd_test.go"
Cohesion: 0.26
Nodes (10): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+2 more)

### Community 145 - "newTestHandler"
Cohesion: 0.23
Nodes (12): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+4 more)

### Community 146 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 147 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 148 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 149 - "JobHealthRecord"
Cohesion: 0.29
Nodes (4): JobHealthRecord, Store, scanJobHealth(), fakeHealthStore

### Community 150 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 151 - "Decision"
Cohesion: 0.17
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 152 - "Decision"
Cohesion: 0.17
Nodes (12): 0005 — Cross-replica WebSocket event relay for active/active, Authorization and secrets, Consequences, Context, Decision, Duplicate suppression and ordering, Event ownership: local first, then relay, Mechanism: PostgreSQL LISTEN/NOTIFY (+4 more)

### Community 153 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 154 - "net/http.Response"
Cohesion: 0.25
Nodes (8): isSafeMethod(), IsSafeMethod(), idempotent(), retryable(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 155 - "compliance/engine_test.go"
Cohesion: 0.29
Nodes (10): Evaluate(), testCatalog(), TestEvaluateAndKindAndOrder(), TestEvaluateEdgeCases(), TestEvaluateLargeSnapshot(), TestEvaluateOperators(), TestValidate(), TestValidateFieldAttribution() (+2 more)

### Community 156 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 157 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 158 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 159 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 160 - "RequireElevation"
Cohesion: 0.29
Nodes (8): AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), recordElevationAudit(), RequireElevation(), go_pkg_github_com_wiselabz_wiselabz_internal_api_middleware

### Community 161 - "RequireConnectorRole"
Cohesion: 0.33
Nodes (9): ConnectorRoleChecker, Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor() (+1 more)

### Community 162 - "compliance_rules_test.go"
Cohesion: 0.31
Nodes (8): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), go_pkg_regexp

### Community 163 - "Engine"
Cohesion: 0.22
Nodes (4): NewHandler(), Engine, sync.Map, DocRegenerator

### Community 164 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 165 - "NewHandler"
Cohesion: 0.24
Nodes (10): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+2 more)

### Community 166 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 167 - "go_pkg_github_com_go_chi_chi_v5"
Cohesion: 0.20
Nodes (7): chi.Router, mountConnectorRoutes(), chi.Router, mountMCPRoutes(), chi.Router, mountReportRoutes(), go_pkg_github_com_go_chi_chi_v5

### Community 168 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 171 - "Product"
Cohesion: 0.20
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 172 - "settings.ts"
Cohesion: 0.33
Nodes (8): zustand, MotionProvider(), apply(), framerReducedMotion(), MotionPref, seed(), SettingsState, useSettings

### Community 173 - "changes.sh"
Cohesion: 0.40
Nodes (9): classify_path(), cmd_check(), cmd_classify(), cmd_go_closure(), describe_flags(), die(), go_package_dirs(), package_version_only() (+1 more)

### Community 174 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 175 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 176 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 177 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 178 - "Elector"
Cohesion: 0.25
Nodes (3): database/sql.Conn, Elector, Noop

### Community 179 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 180 - "auth/handlers_test.go"
Cohesion: 0.25
Nodes (7): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups()

### Community 181 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 182 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 183 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 184 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 185 - "migrations.go"
Cohesion: 0.25
Nodes (6): go_pkg_embed, go_pkg_github_com_golang_migrate_migrate_v4, go_pkg_github_com_golang_migrate_migrate_v4_database_postgres, go_pkg_github_com_golang_migrate_migrate_v4_database_sqlite3, go_pkg_github_com_golang_migrate_migrate_v4_source_file, go_pkg_github_com_golang_migrate_migrate_v4_source_iofs

### Community 186 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 187 - "Audit Trail"
Cohesion: 0.25
Nodes (7): Audit Trail, Endpoint, Filtering and export, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 188 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 189 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 190 - "RequirePermission"
Cohesion: 0.38
Nodes (5): PermissionChecker, chi.Router, mountDashboardRoutes(), mountWorkflowRoutes(), RequirePermission()

### Community 191 - "TestBulkReauth"
Cohesion: 0.48
Nodes (7): bulkReq(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync()

### Community 192 - "MigratedSQLite"
Cohesion: 0.33
Nodes (6): Handler, SyncDocEmbeddings(), TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(), TestRetrieveUsesCacheAndSyncInvalidates(), buildTemplate(), MigratedSQLite()

### Community 193 - "newTestHandler"
Cohesion: 0.29
Nodes (7): Handler, newTestHandler(), TestCreate(), TestExecuteStepNotFound(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete()

### Community 194 - "gitAuth"
Cohesion: 0.33
Nodes (6): gitAuth(), installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), GitOptions, github.com/go-git/go-git/v5/plumbing/transport.AuthMethod

### Community 195 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 196 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 197 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 198 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 200 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 201 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 202 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 203 - "GuardedDialer"
Cohesion: 0.40
Nodes (5): GuardedDialer(), IsDangerousIP(), newWebhookClient(), net.Dialer, net.IP

### Community 205 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 206 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 207 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 208 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 209 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 210 - "mountPublicSystemRoutes"
Cohesion: 0.67
Nodes (3): chi.Router, mountPublicSystemRoutes(), mountRootRoutes()

### Community 211 - "ServiceDependency"
Cohesion: 0.50
Nodes (4): ServiceDependency, environmentDependencies(), poolDependencies(), networkDependencies()

## Knowledge Gaps
- **602 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+597 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1321 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `WiseLabz — Architecture & Technical Decisions` connect `WiseLabz — Architecture & Technical Decisions` to `icons.tsx`, `Development workflow`, `Authentication design`, `ARCHITECTURE.md`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Why does `Frontend shell & theme (decided 2026-06)` connect `icons.tsx` to `WiseLabz — Architecture & Technical Decisions`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `newTestHarness`, `Handler`, `runRestore`, `NewRegistry`, `chat/chat.go`, `Errorf`, `retention/retention_test.go`, `dispatcher_test.go`, `go_pkg_context`, `ServiceSnapshot`, `ReportData`, `New`, `ConnectorRecord`, `HashToken`, `.call`, `NewEngine`, `Engine`, `NewHandler`, `.call`, `NewUser`, `Store`, `Get`, `ExportToFile`, `response.go`, `Dispatcher`, `net/http.ResponseWriter`, `newRouterDeps`, `NewChecker`, `main`, `NewStore`, `MigratedSQLite`, `rewritePlaceholders`, `Config`, `engine_maintenance_test.go`, `sync.Mutex`, `Manager`, `routerDeps`, `Checker`, `NewEngine`, `gitFixture`, `Handler`, `export_test.go`, `time.Time`, `Handler`, `diagnostics/diagnostics.go`, `Deps`, `testApp`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _602 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.0213903743315508 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019337016574585635 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.023448275862068966 - nodes in this community are weakly interconnected._