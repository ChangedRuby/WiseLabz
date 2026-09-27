# Graph Report - refactor-connector-capabilities-descriptor-and-s  (2026-09-27)

## Corpus Check
- 869 files · ~537,456 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6474 nodes · 20598 edges · 221 communities (206 shown, 15 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1659 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8567eaff`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- context.Context
- App.tsx
- SystemPage.tsx
- net/http.Request
- icons.tsx
- go_pkg_net_http
- go_pkg_strings
- @tanstack/react-query
- NewChecker
- go_pkg_context
- newTestHandler
- cn
- ServiceDetailPage.tsx
- react
- ServiceSnapshot
- DashboardPage.tsx
- go_pkg_testing
- net/http.Client
- Store
- UsersPage.tsx
- export_test.go
- DecodeKey
- SnapshotEntity
- AlertsPage.tsx
- WebSocketProvider.tsx
- UserIDFromContext
- package.json
- portainer_test.go
- Runner
- NewEngine
- dispatcher_test.go
- fixtures.ts
- HashToken
- ExportToFile
- home_assistant/tables.go
- RunMigrations
- routerDeps
- truenas/tables.go
- dependencies
- User
- docker_test.go
- ConnectorRecord
- DBTX
- router.go
- portainer/tables.go
- config_test.go
- Connector
- adguardhome/tables.go
- NewStore
- home_assistant_test.go
- Dispatcher
- ThemeControls.tsx
- rowScanner
- traefik/tables.go
- SuggestWithFallback
- unifi/tables.go
- response.go
- MarshalConnectorConfig
- NewUser
- Connector
- settings.mock.ts
- .OIDCCallback
- share_links_test.go
- system/backup.go
- rewritePlaceholders
- ReportsPage.tsx
- Service
- git.go
- handlers.ts
- Handler
- Connector
- connector/connector.go
- NewEngine
- unifi_test.go
- timeline.ts
- newTestHandler
- snapshotdiff.go
- AppearancePage.tsx
- main
- Config
- nilToStr
- WiseLabz — Design Contract
- devDependencies
- Manager
- runbooks_test.go
- apikey_scope.go
- DecodeJSON
- New
- Compare
- diagnostics/diagnostics.go
- time.Time
- middleware_test.go
- Store
- NewRegistry
- handlers_contract_test.go
- NewHTTPClient
- adguardhome_test.go
- connector_permission.go
- chat/chat.go
- backup/backup.go
- .Fetch
- Deps
- log/slog.Logger
- Handler
- sshStdioConn
- truenas_test.go
- ws.ts
- compilerOptions
- Register
- Connector
- traefik_test.go
- Store
- docdiffmodel.ts
- system/handlers_test.go
- all.go
- httpx/retry_test.go
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- net/http.Handler
- testApp
- AuthMiddleware
- lifecycleManager
- handlers_actions_test.go
- Handler
- NewClient
- NewService
- vectorCache
- fetch_test.go
- Connector
- DocRecord
- changes/handlers_test.go
- pagination_contract_test.go
- newTestHandler
- registry.go
- NotificationRecord
- render_test.go
- time.Duration
- api/changes_test.go
- connectors_health_test.go
- Contributing to WiseLabz
- Decision
- scripts
- templatefuncs.go
- ReportData
- Engine
- Decision
- main.tsx
- Handler
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- provider_test.go
- .call
- api/docs_test.go
- .call
- DeliveryRecord
- Elector
- Handler
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- ComplianceRuleRecord
- handlers_bulk_test.go
- connectors_hardening_test.go
- net/http.Response
- WiseLabz — Deployment Guide
- release-please-config.json
- alerts/handlers_test.go
- dashboard/handlers_test.go
- RateLimit
- validate.go
- ComputeWindow
- runbook_test.go
- Cache
- Step by step
- WiseLabz
- webAuthnUser
- openapi_contract_test.go
- newTestHandler
- serveSSHDockerConn
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- walkCursorPages
- cloudflare/attributes_test.go
- Store
- engine_maintenance_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Mermaid.tsx
- Security Policy
- withGrant
- seedScopeFixture
- Enforcement Guidelines
- Authentication design
- Development workflow
- MfaEnrollDialog
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- WiseLabz — v2 Backlog
- fakeQualityChecker
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 229 edges
2. `Errorf()` - 183 edges
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
Nodes (174): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+166 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (161): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+153 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (155): newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestStandbyIsUnreadyAndRunsNoScheduler(), TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks() (+147 more)

### Community 3 - "context.Context"
Cohesion: 0.03
Nodes (36): MFAEnrollOnlyFromContext(), Connector, Store, existingIDs(), Store, placeholders(), scanBackupRun(), Store (+28 more)

### Community 4 - "App.tsx"
Cohesion: 0.02
Nodes (105): react-router-dom, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate (+97 more)

### Community 5 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (104): react-i18next, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules (+96 more)

### Community 6 - "net/http.Request"
Cohesion: 0.05
Nodes (43): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), setFlowCookie() (+35 more)

### Community 7 - "icons.tsx"
Cohesion: 0.04
Nodes (89): web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid, web_src_api_generated_docs_docs (+81 more)

### Community 8 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (41): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection (+33 more)

### Community 9 - "go_pkg_strings"
Cohesion: 0.04
Nodes (51): contextKey, elevationError, TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups() (+43 more)

### Community 10 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (59): i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage, CommandPalette() (+51 more)

### Community 11 - "NewChecker"
Cohesion: 0.06
Nodes (73): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+65 more)

### Community 12 - "go_pkg_context"
Cohesion: 0.08
Nodes (16): contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_go_webauthn_webauthn_protocol, go_pkg_github_com_go_webauthn_webauthn_webauthn (+8 more)

### Community 13 - "newTestHandler"
Cohesion: 0.05
Nodes (77): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+69 more)

### Community 14 - "cn"
Cohesion: 0.04
Nodes (72): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_templates_templates (+64 more)

### Community 15 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (74): ADR-0001, ADR-0003, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridconfigfields (+66 more)

### Community 16 - "react"
Cohesion: 0.04
Nodes (66): RFC-3339, match-sorter, motion, @radix-ui/react-popover, react, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow (+58 more)

### Community 17 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (23): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, init(), RegisterTransformer() (+15 more)

### Community 18 - "DashboardPage.tsx"
Cohesion: 0.05
Nodes (67): 4. `change.detected`, 5. `alert.created`, 8. `doc.generated`, web_src_api_generated_dashboard_dashboard_getdashboardlayout, web_src_api_generated_dashboard_dashboard_getdashboardlayoutadmindefault, web_src_api_generated_dashboard_dashboard_getgetdashboardlayoutadmindefaultquerykey, web_src_api_generated_dashboard_dashboard_postdashboardlayoutreset, web_src_api_generated_dashboard_dashboard_putdashboardlayout (+59 more)

### Community 19 - "go_pkg_testing"
Cohesion: 0.05
Nodes (27): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), TestLoggablePathMasksShareTokenUnderV1() (+19 more)

### Community 20 - "net/http.Client"
Cohesion: 0.04
Nodes (28): Connector, ollamaEmbedder, openAIEmbedder, NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL() (+20 more)

### Community 21 - "Store"
Cohesion: 0.05
Nodes (34): Config, Handler, Registry, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+26 more)

### Community 22 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (52): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+44 more)

### Community 23 - "export_test.go"
Cohesion: 0.07
Nodes (44): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+36 more)

### Community 24 - "DecodeKey"
Cohesion: 0.07
Nodes (33): factorJSON(), Handler, Handler, Handler, Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode() (+25 more)

### Community 25 - "SnapshotEntity"
Cohesion: 0.08
Nodes (54): SnapshotEntity, SnapshotSection, NewMalformedResponseError(), TestBuildContainerTableAttributes(), buildContainerTable(), unavailable(), TestBuildInterfaceTableAttributes(), buildInterfaceTable() (+46 more)

### Community 26 - "AlertsPage.tsx"
Cohesion: 0.05
Nodes (53): Frontend shell & theme (decided 2026-06), Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, react-error-boundary, sonner (+45 more)

### Community 27 - "WebSocketProvider.tsx"
Cohesion: 0.05
Nodes (53): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete` (+45 more)

### Community 28 - "UserIDFromContext"
Cohesion: 0.06
Nodes (28): Handler, PermissionChecker, newToken(), sanitize(), Handler, configRequestField(), parseScheduleUpdates(), validateConnectorConfig() (+20 more)

### Community 29 - "package.json"
Cohesion: 0.04
Nodes (49): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+41 more)

### Community 30 - "portainer_test.go"
Cohesion: 0.07
Nodes (53): TestRegisteredSchema(), TestSchemaConfigValidation(), countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys() (+45 more)

### Community 31 - "Runner"
Cohesion: 0.08
Nodes (30): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner, New() (+22 more)

### Community 32 - "NewEngine"
Cohesion: 0.09
Nodes (41): TestTreeEmpty(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), NewEngine() (+33 more)

### Community 33 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 34 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 35 - "HashToken"
Cohesion: 0.08
Nodes (27): updateUserRequest, Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), HashPassword() (+19 more)

### Community 36 - "ExportToFile"
Cohesion: 0.10
Nodes (46): Export(), ExportToFile(), Import(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile() (+38 more)

### Community 37 - "home_assistant/tables.go"
Cohesion: 0.08
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+31 more)

### Community 38 - "RunMigrations"
Cohesion: 0.09
Nodes (40): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), Store, newPostgresTestStore(), GetMigrationStatus(), newMigrator() (+32 more)

### Community 39 - "routerDeps"
Cohesion: 0.10
Nodes (35): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+27 more)

### Community 40 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 41 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 42 - "User"
Cohesion: 0.07
Nodes (10): sanitizeSessions(), Store, Store, boolToInt(), Session, Store, User, isUniqueViolation() (+2 more)

### Community 43 - "docker_test.go"
Cohesion: 0.07
Nodes (40): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+32 more)

### Community 44 - "ConnectorRecord"
Cohesion: 0.09
Nodes (17): Sanitize(), TestSanitize(), changeServiceIDs(), AlertRecord, ChangeRecord, ConnectorRecord, Store, scanConnectorRows() (+9 more)

### Community 45 - "DBTX"
Cohesion: 0.09
Nodes (18): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), changeFilterClause(), countRows(), T (+10 more)

### Community 46 - "router.go"
Cohesion: 0.08
Nodes (31): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+23 more)

### Community 47 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 48 - "config_test.go"
Cohesion: 0.08
Nodes (35): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), runHealthcheck(), Load() (+27 more)

### Community 49 - "Connector"
Cohesion: 0.07
Nodes (9): ConfigField, TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), Connector, Connector (+1 more)

### Community 50 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 51 - "NewStore"
Cohesion: 0.12
Nodes (35): TestEmbeddedSPAWithoutFrontendBuild(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), NewStore(), Token() (+27 more)

### Community 52 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (35): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+27 more)

### Community 53 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 54 - "ThemeControls.tsx"
Cohesion: 0.11
Nodes (31): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls() (+23 more)

### Community 55 - "rowScanner"
Cohesion: 0.11
Nodes (18): scanAlert(), scanChange(), scanConnector(), nullInt64ToIntPtr(), nullStrToStr(), scanDelivery(), RunbookRecord, RunbookStepRecord (+10 more)

### Community 56 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 57 - "SuggestWithFallback"
Cohesion: 0.10
Nodes (14): claudeProvider, openAICompatibleProvider, Provider, StatusError, StubProvider, SuggestResult, ProviderConfig, SuggestChunk (+6 more)

### Community 58 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 59 - "response.go"
Cohesion: 0.08
Nodes (22): DecodeCursor(), EncodeCursor(), TestCursorRequestModes(), TestCursorRoundTrip(), TestDecodeCursorRejectsGarbage(), TestNextCursorStopsOnShortPage(), TestWritePaginatedOmitsNextCursor(), Error() (+14 more)

### Community 60 - "MarshalConnectorConfig"
Cohesion: 0.11
Nodes (20): capitalize(), Handler, TestDiagnosticsRedactsSecrets(), WriteElevationError(), LifecycleOp(), IsSecretFieldType(), MarshalConnectorConfig(), ParseConnectorConfig() (+12 more)

### Community 61 - "NewUser"
Cohesion: 0.14
Nodes (31): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationDocVisibility() (+23 more)

### Community 62 - "Connector"
Cohesion: 0.11
Nodes (12): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap(), Connector, apiMessage(), controllerName(), countByKind(), statusError() (+4 more)

### Community 63 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 64 - ".OIDCCallback"
Cohesion: 0.13
Nodes (11): Handler, Handler, newOIDCUser(), randomOIDCToken(), readOIDCFlowCookie(), validHostPort(), OIDCClaims, OIDCProvider (+3 more)

### Community 65 - "share_links_test.go"
Cohesion: 0.19
Nodes (29): Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestListEmpty(), TestSave(), TestVersionsOfUnknownDoc() (+21 more)

### Community 66 - "system/backup.go"
Cohesion: 0.10
Nodes (22): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+14 more)

### Community 67 - "rewritePlaceholders"
Cohesion: 0.11
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 68 - "ReportsPage.tsx"
Cohesion: 0.08
Nodes (25): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_getreportsreportiddownload, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid (+17 more)

### Community 69 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 70 - "git.go"
Cohesion: 0.10
Nodes (20): TestCommitMessage(), keys(), writeExportState(), exportCursor, exportState, dockerSSHAddr, go_pkg_crypto_ed25519, go_pkg_encoding_pem (+12 more)

### Community 71 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 72 - "Handler"
Cohesion: 0.12
Nodes (12): Handler, isWritableField(), validateConfigPushRequest(), stepAuditDetail(), validTargetType(), validVerb(), ConfigPusher, ValidateCompositeRef() (+4 more)

### Community 73 - "Connector"
Cohesion: 0.12
Nodes (8): init(), Connector, buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment(), ValidateRefSegment(), Connector

### Community 74 - "connector/connector.go"
Cohesion: 0.09
Nodes (15): Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP(), NewTimeoutError(), newWebhookClient(), AuthError (+7 more)

### Community 75 - "NewEngine"
Cohesion: 0.18
Nodes (23): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+15 more)

### Community 76 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 77 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 78 - "newTestHandler"
Cohesion: 0.14
Nodes (24): TestConnectorStoreErrorPaths(), createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler(), Handler (+16 more)

### Community 79 - "snapshotdiff.go"
Cohesion: 0.14
Nodes (24): BuildSnapshotDiff(), CompareDependencies(), CompareEntities(), entityKey(), entityMap(), TestCompareDependenciesDeterministic(), TestCompareEntitiesKeepsExternalIDAndFallbackNameDistinct(), TestSnapshotDiffCSVNeutralizesFormulaCellsOnly() (+16 more)

### Community 80 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 81 - "main"
Cohesion: 0.12
Nodes (19): main(), newLogger(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder (+11 more)

### Community 82 - "Config"
Cohesion: 0.13
Nodes (19): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, BackupSettings (+11 more)

### Community 83 - "nilToStr"
Cohesion: 0.11
Nodes (8): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 84 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 85 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 86 - "Manager"
Cohesion: 0.16
Nodes (8): cron.EntryID, Manager, LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store, Scheduler

### Community 87 - "runbooks_test.go"
Cohesion: 0.18
Nodes (21): runbookResp, runbookStepResp, createRunbookWithStep(), testApp, seedProxmoxConnector(), TestRunbookExecuteStepDryRunWithoutTokenWorks(), TestRunbookExecuteStepFailureCreatesAlert(), TestRunbookExecuteStepForbiddenWithoutOperatorGrant() (+13 more)

### Community 88 - "apikey_scope.go"
Cohesion: 0.16
Nodes (20): APIKeyRestriction, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), treatAsSafeFromContext(), seedAlert(), TestListAttentionItems() (+12 more)

### Community 89 - "DecodeJSON"
Cohesion: 0.18
Nodes (8): webAuthnFlow, Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T, github.com/go-webauthn/webauthn/webauthn.SessionData

### Community 90 - "New"
Cohesion: 0.12
Nodes (22): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+14 more)

### Community 91 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 92 - "diagnostics/diagnostics.go"
Cohesion: 0.18
Nodes (19): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+11 more)

### Community 93 - "time.Time"
Cohesion: 0.17
Nodes (20): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift (+12 more)

### Community 94 - "middleware_test.go"
Cohesion: 0.13
Nodes (18): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser() (+10 more)

### Community 95 - "Store"
Cohesion: 0.13
Nodes (8): fakeStatusChecker, testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 96 - "NewRegistry"
Cohesion: 0.23
Nodes (19): registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders(), TestSuggestWithFallbackStopsOnNonRetryableError(), NewRegistry() (+11 more)

### Community 97 - "handlers_contract_test.go"
Cohesion: 0.18
Nodes (19): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+11 more)

### Community 98 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 99 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 100 - "connector_permission.go"
Cohesion: 0.19
Nodes (9): auditConnectorGrantDiffJSON(), getConnectorGrant(), ConnectorGrantDiff, Store, highestConnectorRole(), listOIDCConnectorGrants(), scanConnectorGrants(), upsertConnectorGrant() (+1 more)

### Community 101 - "chat/chat.go"
Cohesion: 0.14
Nodes (18): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections() (+10 more)

### Community 102 - "backup/backup.go"
Cohesion: 0.24
Nodes (19): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, importBundle(), importConnectors() (+11 more)

### Community 103 - ".Fetch"
Cohesion: 0.14
Nodes (11): ServiceDependency, WantsField(), TestRequestedFields(), TestWantsField(), environmentDependencies(), putMetadata(), unavailable(), networkDependencies() (+3 more)

### Community 104 - "Deps"
Cohesion: 0.23
Nodes (18): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+10 more)

### Community 105 - "log/slog.Logger"
Cohesion: 0.22
Nodes (14): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+6 more)

### Community 106 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 107 - "sshStdioConn"
Cohesion: 0.12
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 108 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 109 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 110 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 111 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 112 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 113 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 114 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 115 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 116 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 117 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 118 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 119 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 120 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 121 - "net/http.Handler"
Cohesion: 0.16
Nodes (12): ConnectorRoleChecker, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), TreatAsSafeMethod() (+4 more)

### Community 122 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 123 - "AuthMiddleware"
Cohesion: 0.17
Nodes (11): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), AuthMiddleware(), extractBearerToken() (+3 more)

### Community 124 - "lifecycleManager"
Cohesion: 0.16
Nodes (7): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, lifecycleDeps, lifecycleManager

### Community 125 - "handlers_actions_test.go"
Cohesion: 0.32
Nodes (14): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+6 more)

### Community 126 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 127 - "NewClient"
Cohesion: 0.18
Nodes (14): isSafeMethod(), clientTimeout(), IsSafeMethod(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects() (+6 more)

### Community 128 - "NewService"
Cohesion: 0.25
Nodes (14): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+6 more)

### Community 129 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 130 - "fetch_test.go"
Cohesion: 0.19
Nodes (14): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+6 more)

### Community 132 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 133 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 134 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 135 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 136 - "registry.go"
Cohesion: 0.18
Nodes (11): supportedLifecycleVerbs(), IsCredentialRefresherType(), ListSchemas(), SupportsLifecycleVerb(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError (+3 more)

### Community 137 - "NotificationRecord"
Cohesion: 0.23
Nodes (6): Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord, Store, scanNotification()

### Community 138 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 139 - "time.Duration"
Cohesion: 0.19
Nodes (6): AuthSettings, Database, Server, WebAuthnSettings, time.Duration, OIDCProvider

### Community 140 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 141 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 142 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 143 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 144 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 145 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 146 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 147 - "Engine"
Cohesion: 0.18
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 148 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 149 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 151 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 152 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 153 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 154 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 155 - "provider_test.go"
Cohesion: 0.29
Nodes (6): testProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream()

### Community 156 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 157 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 158 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 159 - "DeliveryRecord"
Cohesion: 0.29
Nodes (10): decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter(), TestListDeliveriesStatusFilter(), PaginatedResponse (+2 more)

### Community 160 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 162 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 163 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 164 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 165 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 166 - "handlers_bulk_test.go"
Cohesion: 0.56
Nodes (8): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync()

### Community 167 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 168 - "net/http.Response"
Cohesion: 0.33
Nodes (6): retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 169 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 170 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 171 - "alerts/handlers_test.go"
Cohesion: 0.43
Nodes (7): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze()

### Community 172 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 173 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 174 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 175 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 176 - "runbook_test.go"
Cohesion: 0.32
Nodes (7): Store, newCascadeTestStore(), TestGetRunbookByTarget(), TestRunbookRoundTrip(), TestRunbookStepsCascadeOnConnectorDelete(), TestRunbookStepsCascadeOnRunbookDelete(), TestRunbookStepsRoundTrip()

### Community 177 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 178 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 179 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 180 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 181 - "openapi_contract_test.go"
Cohesion: 0.48
Nodes (6): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 182 - "newTestHandler"
Cohesion: 0.29
Nodes (7): Handler, newTestHandler(), TestCreate(), TestExecuteStepNotFound(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete()

### Community 183 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 184 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 185 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 186 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

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

### Community 191 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 192 - "cloudflare/attributes_test.go"
Cohesion: 0.33
Nodes (5): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable()

### Community 194 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 195 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 196 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 197 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 198 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 199 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 200 - "withGrant"
Cohesion: 0.40
Nodes (5): TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestTree(), withGrant()

### Community 201 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 202 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 203 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 204 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 205 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 206 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 207 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 211 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 212 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **572 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+567 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1295 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `UserIDFromContext()` connect `UserIDFromContext` to `.OIDCCallback`, `net/http.Handler`, `Handler`, `HashToken`, `context.Context`, `net/http.Request`, `routerDeps`, `Deps`, `go_pkg_strings`, `Handler`, `DBTX`, `Store`, `Handler`, `DecodeKey`, `DecodeJSON`, `Handler`, `response.go`, `middleware_test.go`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `testing.T`, `net/http.Request`, `NewChecker`, `go_pkg_context`, `newTestHandler`, `ReportData`, `Engine`, `Handler`, `export_test.go`, `UserIDFromContext`, `.call`, `.call`, `DeliveryRecord`, `NewEngine`, `dispatcher_test.go`, `Handler`, `HashToken`, `ExportToFile`, `RunMigrations`, `alerts/handlers_test.go`, `ConnectorRecord`, `DBTX`, `NewStore`, `Dispatcher`, `response.go`, `NewUser`, `share_links_test.go`, `engine_maintenance_test.go`, `rewritePlaceholders`, `withGrant`, `Handler`, `NewEngine`, `Manager`, `apikey_scope.go`, `New`, `diagnostics/diagnostics.go`, `time.Time`, `NewRegistry`, `chat/chat.go`, `backup/backup.go`, `Deps`, `log/slog.Logger`, `Handler`, `testApp`, `lifecycleManager`, `Handler`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Why does `newTestHandler()` connect `newTestHandler` to `NewService`, `handlers_contract_test.go`, `testing.T`, `handlers_bulk_test.go`, `NewEngine`, `NewStore`, `home_assistant_test.go`, `handlers_actions_test.go`, `Store`, `NewUser`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _572 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.02127659574468085 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.021969744523256545 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.020948675744426156 - nodes in this community are weakly interconnected._