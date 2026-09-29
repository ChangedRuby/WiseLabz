# Graph Report - test-reuse-a-migrated-sqlite-template-inside-int  (2026-09-28)

## Corpus Check
- 877 files · ~545,268 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6541 nodes · 20762 edges · 221 communities (206 shown, 15 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1665 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fb615e3f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- context.Context
- icons.tsx
- go_pkg_testing
- cn
- ServiceDetailPage.tsx
- Button.tsx
- @tanstack/react-query
- go_pkg_net_http
- DashboardPage.tsx
- newTestHandler
- NewChecker
- App.tsx
- ServiceSnapshot
- net/http.Client
- dispatcher_test.go
- go_pkg_context
- Errorf
- git_test.go
- UsersPage.tsx
- SystemPage.tsx
- Runner
- ProfilePage.tsx
- NewEngine
- rowScanner
- package.json
- NewUser
- TemplateEditorPage.tsx
- Compare
- net/http.Request
- fixtures.ts
- DecodeJSON
- DecodeKey
- RulesPage.tsx
- ExportToFile
- Connector
- SnapshotEntity
- routerDeps
- RunMigrations
- newTestHandler
- NewMalformedResponseError
- net/http.ResponseWriter
- docker_test.go
- home_assistant/tables.go
- time.Duration
- dependencies
- NewEngine
- Store
- Hub
- Register
- home_assistant_test.go
- portainer/tables.go
- go_pkg_github_com_wiselabz_wiselabz_internal_config
- Dispatcher
- .OIDCCallback
- router.go
- adguardhome/tables.go
- go_pkg_os
- GetTypeSchema
- response.go
- traefik/tables.go
- unifi/tables.go
- Store
- nilToStr
- connector_permission.go
- traefik_test.go
- New
- rewritePlaceholders
- settings.mock.ts
- HashPassword
- connector/connector.go
- Connector
- Manager
- NewStore
- SuggestRequest
- Service
- Connector
- src/theme.ts
- handlers.ts
- api/auth/oidc.go
- unifi_test.go
- ConnectorRecord
- timeline.ts
- backup/backup.go
- httpx/retry_test.go
- AppearancePage.tsx
- render_test.go
- WiseLabz — Design Contract
- devDependencies
- NewRegistry
- AuthMiddleware
- sync.Mutex
- newTestHandler
- portainer_test.go
- main
- compliance/handlers.go
- adguardhome_test.go
- .Fetch
- Deps
- ErrorWithDetails
- api/mcp_test.go
- NewService
- NotificationRecord
- RunbookRecord
- middleware.go
- channels.go
- Connector
- log/slog.Logger
- MarshalConnectorConfig
- ws/ws_test.go
- logging_test.go
- NewClient
- sshStdioConn
- truenas_test.go
- diagnostics/diagnostics.go
- time.Time
- ws.ts
- compilerOptions
- notifications/handlers_test.go
- Handler
- Get
- docdiffmodel.ts
- templates_test.go
- system/handlers_test.go
- all.go
- newTestHarness
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- chat/handlers_test.go
- handlers_contract_test.go
- changes/handlers_test.go
- pagination_contract_test.go
- reports/handlers.go
- templatefuncs.go
- Contributing to WiseLabz
- connectors_health_test.go
- Handler
- Engine
- Decision
- scripts
- Store
- Store
- Contributor Covenant Code of Conduct
- Decision
- Backend test performance
- main.tsx
- net/http.Handler
- findings/handlers_authz_test.go
- chat/chat.go
- vectorCache
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- .call
- cursor_pagination_test.go
- connectors_maintenance_test.go
- api/docs_test.go
- Elector
- store/mfa_test.go
- Changelog
- test-shards.sh
- mockServiceWorker.js
- config_cmd_test.go
- apikey_scopes_test.go
- openapi_contract_test.go
- WiseLabz — Deployment Guide
- release-please-config.json
- connectors_hardening_test.go
- snapshots_test.go
- dashboard/handlers_test.go
- RateLimit
- ComputeWindow
- ComplianceRuleRecord
- runbook_test.go
- ShareLink
- Cache
- Step by step
- WiseLabz
- webAuthnUser
- handlers_start_stop_configpush_test.go
- .UpdateAuthConfig
- Connector
- scanMaintenanceWindow
- computeNextRun
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- cloudflare/attributes_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Mermaid.tsx
- Security Policy
- APIKeyClaims
- ref_test.go
- Authentication design
- Development workflow
- MfaEnrollDialog
- @vitejs/plugin-react
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- Saved Views
- dockerSSHAddr
- WiseLabz — v2 Backlog
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

## Communities (221 total, 15 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (180): runbookResp, runbookStepResp, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation() (+172 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (151): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+143 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (147): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestWriteConfigRejection() (+139 more)

### Community 3 - "context.Context"
Cohesion: 0.02
Nodes (35): fakeStatusChecker, sanitizeSessions(), Connector, Connector, existingIDs(), SnapshotRecord, Store, docSearchWhere() (+27 more)

### Community 4 - "icons.tsx"
Cohesion: 0.03
Nodes (98): web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest, web_src_api_generated_docs_docs_postdocsdocidlock, web_src_api_generated_docs_docs_postdocsdocidlockrelease (+90 more)

### Community 5 - "go_pkg_testing"
Cohesion: 0.04
Nodes (46): complianceRule(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), Schema(), schemaFor(), TestSchemaMatchesConfig() (+38 more)

### Community 6 - "cn"
Cohesion: 0.04
Nodes (91): web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid, web_src_api_generated_docs_docs_postdocstopology (+83 more)

### Community 7 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (88): ADR-0001, ADR-0003, RFC-3339, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth (+80 more)

### Community 8 - "Button.tsx"
Cohesion: 0.03
Nodes (88): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss (+80 more)

### Community 9 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (66): 3. `sync.complete`, i18next, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_changes_changes, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_getgetchangesquerykey (+58 more)

### Community 10 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (50): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, updateUserRequest, newToken(), changePromptData(), diffToSpec(), stripPromptTags() (+42 more)

### Community 11 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (91): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 4. `change.detected` (+83 more)

### Community 12 - "newTestHandler"
Cohesion: 0.05
Nodes (82): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+74 more)

### Community 13 - "NewChecker"
Cohesion: 0.06
Nodes (70): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+62 more)

### Community 14 - "App.tsx"
Cohesion: 0.04
Nodes (70): Frontend shell & theme (decided 2026-06), react, react-error-boundary, react-i18next, sonner, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth (+62 more)

### Community 15 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (26): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, changePatternID(), Engine, markError() (+18 more)

### Community 16 - "net/http.Client"
Cohesion: 0.03
Nodes (28): Connector, ollamaEmbedder, openAIEmbedder, setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector (+20 more)

### Community 17 - "dispatcher_test.go"
Cohesion: 0.07
Nodes (69): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext() (+61 more)

### Community 18 - "go_pkg_context"
Cohesion: 0.08
Nodes (15): StatusError, contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_coreos_go_oidc_v3_oidc (+7 more)

### Community 19 - "Errorf"
Cohesion: 0.06
Nodes (22): Handler, sanitize(), Handler, Handler, Handler, Handler, Handler, Handler (+14 more)

### Community 20 - "git_test.go"
Cohesion: 0.06
Nodes (48): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+40 more)

### Community 21 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (52): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+44 more)

### Community 22 - "SystemPage.tsx"
Cohesion: 0.05
Nodes (47): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey, web_src_api_generated_system_system_getsystembackupschedule (+39 more)

### Community 23 - "Runner"
Cohesion: 0.07
Nodes (31): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner, New() (+23 more)

### Community 24 - "ProfilePage.tsx"
Cohesion: 0.05
Nodes (43): web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin, web_src_api_generated_auth_auth_postauthloginmfawebauthnbegin, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_auth_auth_usegetauthelevatemethods (+35 more)

### Community 25 - "NewEngine"
Cohesion: 0.09
Nodes (42): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+34 more)

### Community 26 - "rowScanner"
Cohesion: 0.07
Nodes (26): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), BackupSchedule, Store, scanBackupRun() (+18 more)

### Community 27 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 28 - "NewUser"
Cohesion: 0.13
Nodes (51): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), NewHandler(), Handler, newTestHandler(), TestAISuggestInvalidJSON() (+43 more)

### Community 29 - "TemplateEditorPage.tsx"
Cohesion: 0.05
Nodes (45): web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid (+37 more)

### Community 30 - "Compare"
Cohesion: 0.07
Nodes (45): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+37 more)

### Community 31 - "net/http.Request"
Cohesion: 0.08
Nodes (14): Handler, applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), Handler, decodeBulkRequest(), Handler, Handler (+6 more)

### Community 32 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 33 - "DecodeJSON"
Cohesion: 0.09
Nodes (25): sanitizeUser(), setRefreshCookie(), Handler, factorJSON(), Handler, newOIDCUser(), SecurityHeaders(), TestSecurityHeaders() (+17 more)

### Community 34 - "DecodeKey"
Cohesion: 0.09
Nodes (27): ProviderConfig, Handler, Handler, primaryProviderConfig(), Handler, GenerateTOTPSecret(), TestGenerateTOTPSecretProducesScannableURL(), TestValidateTOTPAcceptsCurrentAndSkewedCodes() (+19 more)

### Community 35 - "RulesPage.tsx"
Cohesion: 0.05
Nodes (43): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+35 more)

### Community 36 - "ExportToFile"
Cohesion: 0.11
Nodes (42): Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+34 more)

### Community 37 - "Connector"
Cohesion: 0.06
Nodes (15): init(), ConfigField, Connector, TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable() (+7 more)

### Community 38 - "SnapshotEntity"
Cohesion: 0.13
Nodes (42): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+34 more)

### Community 39 - "routerDeps"
Cohesion: 0.09
Nodes (37): routerDeps, chi.Router, NewRouter(), wsRoleLabel(), chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes() (+29 more)

### Community 40 - "RunMigrations"
Cohesion: 0.10
Nodes (39): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns() (+31 more)

### Community 41 - "newTestHandler"
Cohesion: 0.11
Nodes (43): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+35 more)

### Community 42 - "NewMalformedResponseError"
Cohesion: 0.10
Nodes (42): NewMalformedResponseError(), TestBuildHostsTableAttributes(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+34 more)

### Community 43 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (11): Handler, webAuthnFlow, Handler, decodeStoredSnapshot(), Handler, snapshotStoreError(), WriteDataPaginated(), SinceFromDays() (+3 more)

### Community 44 - "docker_test.go"
Cohesion: 0.07
Nodes (42): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+34 more)

### Community 45 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 46 - "time.Duration"
Cohesion: 0.08
Nodes (27): newLogger(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+19 more)

### Community 47 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 48 - "NewEngine"
Cohesion: 0.11
Nodes (32): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+24 more)

### Community 49 - "Store"
Cohesion: 0.09
Nodes (24): Provider, Config, Handler, Registry, NewHandler(), NewHandler(), NewHandler(), NewHandler() (+16 more)

### Community 50 - "Hub"
Cohesion: 0.08
Nodes (19): Handler, isWritableField(), validateConfigPushRequest(), loggablePath(), loggableQuery(), WriteElevationError(), ConfigPusher, ValidateCompositeRef() (+11 more)

### Community 51 - "Register"
Cohesion: 0.09
Nodes (33): init(), init(), newConnector(), Connector, init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback() (+25 more)

### Community 52 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (37): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+29 more)

### Community 53 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 54 - "go_pkg_github_com_wiselabz_wiselabz_internal_config"
Cohesion: 0.09
Nodes (25): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+17 more)

### Community 55 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 56 - ".OIDCCallback"
Cohesion: 0.10
Nodes (17): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), Handler, readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), Handler (+9 more)

### Community 57 - "router.go"
Cohesion: 0.09
Nodes (28): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+20 more)

### Community 58 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 59 - "go_pkg_os"
Cohesion: 0.08
Nodes (26): main(), usage(), writeExportState(), exportCursor, exportState, go_pkg_bufio, go_pkg_crypto_ed25519, go_pkg_encoding_pem (+18 more)

### Community 60 - "GetTypeSchema"
Cohesion: 0.10
Nodes (29): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword(), TestRegisteredSchema() (+21 more)

### Community 61 - "response.go"
Cohesion: 0.09
Nodes (24): Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes(), TestCursorRoundTrip() (+16 more)

### Community 62 - "traefik/tables.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+21 more)

### Community 63 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 64 - "Store"
Cohesion: 0.09
Nodes (8): placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange(), ChangeSummary

### Community 65 - "nilToStr"
Cohesion: 0.09
Nodes (12): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+4 more)

### Community 66 - "connector_permission.go"
Cohesion: 0.11
Nodes (18): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), treatAsSafeFromContext(), Store (+10 more)

### Community 67 - "traefik_test.go"
Cohesion: 0.11
Nodes (30): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+22 more)

### Community 68 - "New"
Cohesion: 0.09
Nodes (30): confirm(), formatCounts(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag() (+22 more)

### Community 69 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 70 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 71 - "HashPassword"
Cohesion: 0.10
Nodes (21): mustHashDummyPassword(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty() (+13 more)

### Community 72 - "connector/connector.go"
Cohesion: 0.08
Nodes (17): TimeoutError, GuardedDialer(), IsDangerousIP(), NewServiceUnavailableError(), NewTimeoutError(), TestTypedErrorsWrapAndUnwrap(), newWebhookClient(), AuthError (+9 more)

### Community 73 - "Connector"
Cohesion: 0.12
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 74 - "Manager"
Cohesion: 0.12
Nodes (11): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, ReportDefinitionRecord, ReportRecord (+3 more)

### Community 75 - "NewStore"
Cohesion: 0.15
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 76 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 77 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 78 - "Connector"
Cohesion: 0.19
Nodes (6): unavailable(), SnapshotSection, Connector, unavailable(), unavailable(), session

### Community 79 - "src/theme.ts"
Cohesion: 0.14
Nodes (25): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+17 more)

### Community 80 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 81 - "api/auth/oidc.go"
Cohesion: 0.09
Nodes (16): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), validHostPort() (+8 more)

### Community 82 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 83 - "ConnectorRecord"
Cohesion: 0.14
Nodes (14): ConnectorRecord, Store, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), scanSyncRun() (+6 more)

### Community 84 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 85 - "backup/backup.go"
Cohesion: 0.19
Nodes (24): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+16 more)

### Community 86 - "httpx/retry_test.go"
Cohesion: 0.20
Nodes (20): retryable(), RetryTransport(), sleep(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry() (+12 more)

### Community 87 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 88 - "render_test.go"
Cohesion: 0.20
Nodes (17): connectorFilter(), RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden() (+9 more)

### Community 89 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 90 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 91 - "NewRegistry"
Cohesion: 0.22
Nodes (21): SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders(), TestSuggestWithFallbackStopsOnNonRetryableError() (+13 more)

### Community 92 - "AuthMiddleware"
Cohesion: 0.14
Nodes (18): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, AuthMiddleware(), assertElevationAuditCalls(), boolLabel(), requestWithUser() (+10 more)

### Community 93 - "sync.Mutex"
Cohesion: 0.12
Nodes (7): Handler, writeUserWriteError(), cron.EntryID, Handler, sync/atomic.Bool, sync.Mutex, ReadyState

### Community 94 - "newTestHandler"
Cohesion: 0.11
Nodes (20): testHandler, Handler, instanceAdminRoleFor(), spaHandler(), templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths() (+12 more)

### Community 95 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 96 - "main"
Cohesion: 0.14
Nodes (16): main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder (+8 more)

### Community 97 - "compliance/handlers.go"
Cohesion: 0.22
Nodes (9): catalog(), changedFields(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator (+1 more)

### Community 98 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 99 - ".Fetch"
Cohesion: 0.14
Nodes (11): ServiceDependency, WantsField(), TestRequestedFields(), TestWantsField(), environmentDependencies(), putMetadata(), unavailable(), networkDependencies() (+3 more)

### Community 100 - "Deps"
Cohesion: 0.20
Nodes (20): registerListAttentionItems(), changeServiceIDs(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs() (+12 more)

### Community 101 - "ErrorWithDetails"
Cohesion: 0.19
Nodes (10): parseScheduleUpdates(), validateRotationFields(), writeConfigRejection(), stepAuditDetail(), validTargetType(), ErrorWithDetails(), FieldError, Handler (+2 more)

### Community 102 - "api/mcp_test.go"
Cohesion: 0.14
Nodes (12): testApp, mintAPIKey(), TestMCPConnectorRestrictedKey(), TestMCPEndToEnd(), go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server (+4 more)

### Community 103 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 104 - "NotificationRecord"
Cohesion: 0.14
Nodes (10): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord (+2 more)

### Community 105 - "RunbookRecord"
Cohesion: 0.23
Nodes (7): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep(), isUniqueViolation(), TestIsUniqueViolation()

### Community 106 - "middleware.go"
Cohesion: 0.16
Nodes (14): AuditRecorder, contextKey, elevationError, PermissionChecker, UserStatusChecker, elevationFailureReason(), extractBearerToken(), hashToken() (+6 more)

### Community 107 - "channels.go"
Cohesion: 0.15
Nodes (16): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), buildEmailMessage(), sendSMTPChannel(), splitRecipients() (+8 more)

### Community 108 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 109 - "log/slog.Logger"
Cohesion: 0.22
Nodes (14): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+6 more)

### Community 110 - "MarshalConnectorConfig"
Cohesion: 0.17
Nodes (15): IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly(), TestSecretFieldsChangedFalseOnResubmittedUnchangedSecret() (+7 more)

### Community 111 - "ws/ws_test.go"
Cohesion: 0.19
Nodes (18): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+10 more)

### Community 112 - "logging_test.go"
Cohesion: 0.18
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 113 - "NewClient"
Cohesion: 0.14
Nodes (16): isSafeMethod(), clientTimeout(), IsSafeMethod(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects() (+8 more)

### Community 114 - "sshStdioConn"
Cohesion: 0.12
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 115 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 116 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 117 - "time.Time"
Cohesion: 0.23
Nodes (16): time.Time, ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary (+8 more)

### Community 118 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 119 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 120 - "notifications/handlers_test.go"
Cohesion: 0.26
Nodes (16): TestCreate(), AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty() (+8 more)

### Community 121 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 122 - "Get"
Cohesion: 0.17
Nodes (16): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), Capabilities(), CapabilityDescriptor, LifecycleOp(), supportedLifecycleVerbs() (+8 more)

### Community 123 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 124 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 125 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 126 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 127 - "newTestHarness"
Cohesion: 0.26
Nodes (14): seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding(), TestSearchDocs(), seedFinding() (+6 more)

### Community 128 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 129 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 130 - "chat/handlers_test.go"
Cohesion: 0.27
Nodes (14): TestEmbeddedSPAWithoutFrontendBuild(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), Handler, newHandler() (+6 more)

### Community 131 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 132 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 133 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 134 - "reports/handlers.go"
Cohesion: 0.30
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 135 - "templatefuncs.go"
Cohesion: 0.19
Nodes (11): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+3 more)

### Community 136 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 137 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 138 - "Handler"
Cohesion: 0.23
Nodes (4): createVersion(), templateResponse(), templateVersionResponse(), Handler

### Community 139 - "Engine"
Cohesion: 0.17
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 140 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 141 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 142 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 143 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 144 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 145 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 146 - "Backend test performance"
Cohesion: 0.17
Nodes (12): Backend test performance, CI job times, CI measurements, Coverage strategy, Fixture reuse and lifecycle tests (#405–#407), Follow-ups, Local measurements, Measuring (+4 more)

### Community 147 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 148 - "net/http.Handler"
Cohesion: 0.20
Nodes (9): ConnectorRoleChecker, TreatAsSafeMethod(), RequireConnectorRole(), RequireInstanceAdmin(), contextWithInstanceAdmin(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError(), TestRequireInstanceAdmin() (+1 more)

### Community 149 - "findings/handlers_authz_test.go"
Cohesion: 0.45
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 150 - "chat/chat.go"
Cohesion: 0.25
Nodes (9): cosineSimilarity(), packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips(), TestSplitSections(), unpackVector(), Section (+1 more)

### Community 151 - "vectorCache"
Cohesion: 0.25
Nodes (6): vectorCache, vectorEntry, vectorKey, go_pkg_container_list, container/list.Element, container/list.List

### Community 152 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 153 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 154 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 155 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 156 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 157 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 158 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 159 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 160 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 161 - "store/mfa_test.go"
Cohesion: 0.33
Nodes (9): Store, mfaTestUser(), TestConfirmFactorRejectsSecondConfirmedTOTP(), TestConsumeTOTPStepReplayGuard(), TestDeleteFactorLastOneAlsoLeavesRecoveryCodesForCallerToWipe(), TestDeleteUserFactorsForAdminReset(), TestGetRequire2FADefaultsToNone(), TestRecoveryCodesSingleUse() (+1 more)

### Community 162 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 163 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 164 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 165 - "config_cmd_test.go"
Cohesion: 0.36
Nodes (7): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), io.Writer

### Community 166 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 167 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 168 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 169 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 170 - "connectors_hardening_test.go"
Cohesion: 0.29
Nodes (7): testApp, TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 171 - "snapshots_test.go"
Cohesion: 0.57
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 172 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 173 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 174 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 175 - "ComplianceRuleRecord"
Cohesion: 0.43
Nodes (3): ComplianceRuleRecord, Store, scanComplianceRule()

### Community 176 - "runbook_test.go"
Cohesion: 0.32
Nodes (7): Store, newCascadeTestStore(), TestGetRunbookByTarget(), TestRunbookRoundTrip(), TestRunbookStepsCascadeOnConnectorDelete(), TestRunbookStepsCascadeOnRunbookDelete(), TestRunbookStepsRoundTrip()

### Community 178 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 179 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 180 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 181 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 182 - "handlers_start_stop_configpush_test.go"
Cohesion: 0.48
Nodes (6): createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler()

### Community 183 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 185 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 186 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 187 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 188 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 189 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 190 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 191 - "cloudflare/attributes_test.go"
Cohesion: 0.33
Nodes (5): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable()

### Community 192 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 193 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 194 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 195 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 196 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 197 - "APIKeyClaims"
Cohesion: 0.50
Nodes (3): testAPIKeyChecker, APIKeyClaims, validAPIKey()

### Community 198 - "ref_test.go"
Cohesion: 0.40
Nodes (4): TestValidateCompositeRef(), TestValidateRefSegment(), TestValidateUnixSocketPath(), ValidateUnixSocketPath()

### Community 200 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 201 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 202 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 203 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 204 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 205 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 208 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 209 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 210 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **585 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+580 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1313 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Election` connect `dispatcher_test.go` to `go_pkg_context`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `Errorf` to `DecodeJSON`, `compliance/handlers.go`, `DecodeKey`, `ErrorWithDetails`, `reports/handlers.go`, `routerDeps`, `Handler`, `net/http.ResponseWriter`, `logging_test.go`, `Store`, `Hub`, `response.go`, `.UpdateAuthConfig`, `.OIDCCallback`, `Handler`, `sync.Mutex`, `net/http.Request`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `chat/handlers_test.go`, `testing.T`, `reports/handlers.go`, `Handler`, `Engine`, `NewChecker`, `ServiceSnapshot`, `dispatcher_test.go`, `go_pkg_context`, `Errorf`, `git_test.go`, `findings/handlers_authz_test.go`, `NewEngine`, `rowScanner`, `NewUser`, `.call`, `net/http.Request`, `ExportToFile`, `RunMigrations`, `net/http.ResponseWriter`, `NewEngine`, `Dispatcher`, `New`, `rewritePlaceholders`, `HashPassword`, `Manager`, `NewStore`, `backup/backup.go`, `render_test.go`, `NewRegistry`, `sync.Mutex`, `newTestHandler`, `main`, `compliance/handlers.go`, `Deps`, `ErrorWithDetails`, `log/slog.Logger`, `diagnostics/diagnostics.go`, `time.Time`, `notifications/handlers_test.go`, `newTestHarness`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _585 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.020926136952242828 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.023527754296985066 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022663358147229116 - nodes in this community are weakly interconnected._