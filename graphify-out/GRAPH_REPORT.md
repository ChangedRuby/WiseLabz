# Graph Report - store-newmigrator-leaks-a-dedicated-postgres-con  (2026-09-27)

## Corpus Check
- 869 files · ~537,696 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6476 nodes · 20616 edges · 237 communities (221 shown, 16 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1661 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ef81482d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- testing.T
- newDocTestStore
- App.tsx
- SystemPage.tsx
- go_pkg_testing
- context.Context
- net/http.Request
- newTestHandler
- go_pkg_net_http
- go_pkg_context
- icons.tsx
- ServiceDetailPage.tsx
- Button.tsx
- DashboardPage.tsx
- @tanstack/react-query
- cn
- git_test.go
- ServiceSnapshot
- UsersPage.tsx
- net/http.Client
- ErrorWithDetails
- package.json
- react-router-dom
- ConnectorRecord
- react-i18next
- portainer/tables.go
- dispatcher_test.go
- fixtures.ts
- net/http.ResponseWriter
- router.go
- routerDeps
- Compare
- NewMalformedResponseError
- SnapshotEntity
- DecodeKey
- traefik_test.go
- ExportToFile
- home_assistant/tables.go
- newTestHandler
- dependencies
- config_test.go
- RunMigrations
- time.Duration
- NewChecker
- Get
- response.go
- NewUser
- git.go
- Connector
- Dispatcher
- RunbookRecord
- Store
- adguardhome_test.go
- home_assistant_test.go
- nilToStr
- ThemeControls.tsx
- adguardhome/tables.go
- compliance/engine.go
- apikey_scope.go
- traefik/tables.go
- unifi/tables.go
- AuthMiddleware
- Connector
- WebSocketProvider.tsx
- SuggestRequest
- New
- rewritePlaceholders
- Registry
- NewStore
- Sanitize
- Store
- connector/connector.go
- handlers.ts
- settings.mock.ts
- runbooks_test.go
- Service
- Connector
- NewEngine
- unifi_test.go
- timeline.ts
- AuthedUser
- Manager
- AppearancePage.tsx
- .OIDCCallback
- Engine
- Checker
- WiseLabz — Design Contract
- devDependencies
- go_pkg_os
- Register
- New
- newRouterDeps
- portainer_test.go
- Store
- Handler
- handlers_contract_test.go
- NewHTTPClient
- Deps
- oidc_elevate_test.go
- NewService
- docker_test.go
- Runner
- ReportsPage.tsx
- MarshalConnectorConfig
- chat/chat.go
- mcp/mcp_test.go
- Handler
- logging_test.go
- truenas_test.go
- diagnostics/diagnostics.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- main
- docs/handlers_test.go
- Handler
- sshStdioConn
- NewEngine
- docdiffmodel.ts
- system/handlers_test.go
- all.go
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- testApp
- lifecycleManager
- changes/handlers_test.go
- InstanceAdminFromContext
- vectorCache
- Connector
- Connector
- httpx/retry_test.go
- NotificationRecord
- data.go
- DocRecord
- pagination_contract_test.go
- Handler
- newTestHandler
- diagram.go
- render_test.go
- AuditRecord
- Hub
- NewRegistry
- api/changes_test.go
- Engine
- connectors_health_test.go
- Handler
- ListSchemas
- Contributing to WiseLabz
- Decision
- scripts
- Handler
- templatefuncs.go
- IsSecureRequest
- ReportData
- store/backup_test.go
- JobHealthRecord
- transform_test.go
- Decision
- Connector
- main.tsx
- runbooks/handlers_test.go
- ServiceDependency
- newDockerClient
- newSSHDockerClient
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- compliance_rules_test.go
- handlers_bulk_test.go
- .call
- ratelimit.go
- Elector
- log/slog.Logger
- Changelog
- Connector
- mockServiceWorker.js
- apikey_scopes_test.go
- ComplianceRuleRecord
- connectors_hardening_test.go
- openapi_contract_test.go
- net/http.Response
- time.Time
- retention/retention_test.go
- WiseLabz — Deployment Guide
- release-please-config.json
- pfsense.go
- ComputeWindow
- ShareLink
- Cache
- Step by step
- sync.Mutex
- WiseLabz
- webAuthnUser
- snapshotResponse
- .UpdateAuthConfig
- serveSSHDockerConn
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- registryTestRefresher
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- fakeRefresherConnector
- walkCursorPages
- newHandler
- .Fetch
- .GetConnectorUptime
- Store
- engine_maintenance_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Mermaid.tsx
- Security Policy
- RequireConnectorRole
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
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
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

## Communities (237 total, 16 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (175): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+167 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (160): newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestStandbyIsUnreadyAndRunsNoScheduler(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+152 more)

### Community 2 - "newDocTestStore"
Cohesion: 0.02
Nodes (155): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+147 more)

### Community 3 - "App.tsx"
Cohesion: 0.02
Nodes (107): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken(), setMfaEnrollmentRequiredHandler() (+99 more)

### Community 4 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (114): sonner, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules (+106 more)

### Community 5 - "go_pkg_testing"
Cohesion: 0.04
Nodes (53): RegisterClaude(), TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestRegisterClaudeDefaults(), Schema(), schemaFor() (+45 more)

### Community 6 - "context.Context"
Cohesion: 0.03
Nodes (25): sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, existingIDs(), SnapshotRecord, Store, Store, MFAFactor (+17 more)

### Community 7 - "net/http.Request"
Cohesion: 0.06
Nodes (27): Handler, Handler, diffToSpec(), NewHandler(), Handler, Handler, Handler, Handler (+19 more)

### Community 8 - "newTestHandler"
Cohesion: 0.05
Nodes (78): fixture, mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz() (+70 more)

### Community 9 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (40): bulkSnoozeItemResult, bulkSnoozeRequest, contextKey, elevationError, TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess() (+32 more)

### Community 10 - "go_pkg_context"
Cohesion: 0.08
Nodes (18): dashboardLayout, contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_coreos_go_oidc_v3_oidc (+10 more)

### Community 11 - "icons.tsx"
Cohesion: 0.04
Nodes (77): react, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_chat_chat (+69 more)

### Community 12 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (81): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+73 more)

### Community 13 - "Button.tsx"
Cohesion: 0.04
Nodes (77): Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover (+69 more)

### Community 14 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (80): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+72 more)

### Community 15 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (47): i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_attention_attention_getgetattentionquerykey, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage (+39 more)

### Community 16 - "cn"
Cohesion: 0.04
Nodes (59): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid (+51 more)

### Community 17 - "git_test.go"
Cohesion: 0.07
Nodes (46): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+38 more)

### Community 18 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (17): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, changePatternID(), Engine, markError(), snapshotIDOrNil(), runTransformers() (+9 more)

### Community 19 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (51): axios, customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers (+43 more)

### Community 20 - "net/http.Client"
Cohesion: 0.05
Nodes (23): Connector, ollamaEmbedder, NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), tryParseEntities(), validateCustomURL(), Connector (+15 more)

### Community 21 - "ErrorWithDetails"
Cohesion: 0.07
Nodes (30): Handler, oidcElevateFlow, newToken(), sanitize(), sanitizeUser(), setRefreshCookie(), Handler, mustHashDummyPassword() (+22 more)

### Community 22 - "package.json"
Cohesion: 0.04
Nodes (49): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+41 more)

### Community 23 - "react-router-dom"
Cohesion: 0.07
Nodes (42): react-router-dom, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs, web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, buildCommands(), CommandPalette() (+34 more)

### Community 24 - "ConnectorRecord"
Cohesion: 0.07
Nodes (30): changeFilterClause(), scanAlert(), ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr() (+22 more)

### Community 25 - "react-i18next"
Cohesion: 0.05
Nodes (41): RFC-3339, Frontend shell & theme (decided 2026-06), react-error-boundary, react-i18next, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors (+33 more)

### Community 26 - "portainer/tables.go"
Cohesion: 0.09
Nodes (38): WantsField(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), environmentDependencies(), putMetadata(), buildEnvironmentTable(), buildStackTable(), cell() (+30 more)

### Community 27 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 28 - "fixtures.ts"
Cohesion: 0.06
Nodes (44): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+36 more)

### Community 29 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (22): webAuthnFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie(), setFlowCookie(), setOIDCFlowCookie() (+14 more)

### Community 30 - "router.go"
Cohesion: 0.07
Nodes (33): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+25 more)

### Community 31 - "routerDeps"
Cohesion: 0.09
Nodes (38): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+30 more)

### Community 32 - "Compare"
Cohesion: 0.07
Nodes (42): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+34 more)

### Community 33 - "NewMalformedResponseError"
Cohesion: 0.09
Nodes (43): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+35 more)

### Community 34 - "SnapshotEntity"
Cohesion: 0.13
Nodes (42): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+34 more)

### Community 35 - "DecodeKey"
Cohesion: 0.09
Nodes (23): testHandler, Handler, Handler, Handler, Handler, Config, mask(), DecodeKey() (+15 more)

### Community 36 - "traefik_test.go"
Cohesion: 0.08
Nodes (41): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+33 more)

### Community 37 - "ExportToFile"
Cohesion: 0.11
Nodes (41): ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+33 more)

### Community 38 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 39 - "newTestHandler"
Cohesion: 0.11
Nodes (40): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+32 more)

### Community 40 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 41 - "config_test.go"
Cohesion: 0.07
Nodes (37): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Load(), TestAccessTokenTTLDuration() (+29 more)

### Community 42 - "RunMigrations"
Cohesion: 0.11
Nodes (34): main(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns(), TestMigrationSchemaParity(), RunMigrations() (+26 more)

### Community 43 - "time.Duration"
Cohesion: 0.08
Nodes (27): NewHandler(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+19 more)

### Community 44 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 45 - "Get"
Cohesion: 0.08
Nodes (16): applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), NewHandler(), stepAuditDetail(), validTargetType(), validVerb(), Get() (+8 more)

### Community 46 - "response.go"
Cohesion: 0.08
Nodes (27): Handler, decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor() (+19 more)

### Community 47 - "NewUser"
Cohesion: 0.22
Nodes (36): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestCreateConversationDocVisibility(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser() (+28 more)

### Community 48 - "git.go"
Cohesion: 0.07
Nodes (23): gitAuth(), installHTTPS(), TestCommitMessage(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), TestGitAuthSSH(), writeTestKey(), writeExportState() (+15 more)

### Community 49 - "Connector"
Cohesion: 0.08
Nodes (10): init(), ConfigField, Connector, buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), Connector, PathSegment() (+2 more)

### Community 50 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 51 - "RunbookRecord"
Cohesion: 0.10
Nodes (14): BackupSchedule, Store, scanBackupRun(), Store, Store, RunbookRecord, RunbookStepRecord, Store (+6 more)

### Community 52 - "Store"
Cohesion: 0.08
Nodes (8): placeholders(), AlertRecord, ChangeRecord, Store, scanChange(), DocVersionRecord, Store, ChangeSummary

### Community 53 - "adguardhome_test.go"
Cohesion: 0.10
Nodes (34): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+26 more)

### Community 54 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (35): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+27 more)

### Community 55 - "nilToStr"
Cohesion: 0.08
Nodes (12): ChatConversationRecord, Store, nilToStr(), DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+4 more)

### Community 56 - "ThemeControls.tsx"
Cohesion: 0.11
Nodes (31): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls() (+23 more)

### Community 57 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 58 - "compliance/engine.go"
Cohesion: 0.11
Nodes (31): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+23 more)

### Community 59 - "apikey_scope.go"
Cohesion: 0.09
Nodes (21): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), treatAsSafeFromContext(), TreatAsSafeMethod() (+13 more)

### Community 60 - "traefik/tables.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+21 more)

### Community 61 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 62 - "AuthMiddleware"
Cohesion: 0.09
Nodes (25): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, AuthMiddleware(), extractBearerToken(), hashToken() (+17 more)

### Community 63 - "Connector"
Cohesion: 0.11
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 64 - "WebSocketProvider.tsx"
Cohesion: 0.09
Nodes (25): Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`), web_src_api_generated_changes_changes (+17 more)

### Community 65 - "SuggestRequest"
Cohesion: 0.10
Nodes (13): claudeProvider, openAICompatibleProvider, openAIEmbedder, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+5 more)

### Community 66 - "New"
Cohesion: 0.09
Nodes (29): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+21 more)

### Community 67 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 68 - "Registry"
Cohesion: 0.12
Nodes (18): Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail() (+10 more)

### Community 69 - "NewStore"
Cohesion: 0.15
Nodes (27): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+19 more)

### Community 70 - "Sanitize"
Cohesion: 0.10
Nodes (14): Handler, isWritableField(), decodeBulkRequest(), Handler, loggablePath(), loggableQuery(), ConfigPusher, Err() (+6 more)

### Community 71 - "Store"
Cohesion: 0.19
Nodes (25): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+17 more)

### Community 72 - "connector/connector.go"
Cohesion: 0.08
Nodes (15): TimeoutError, GuardedDialer(), IsDangerousIP(), LifecycleOp(), supportedLifecycleVerbs(), newWebhookClient(), AuthError, CredentialRefresher (+7 more)

### Community 73 - "handlers.ts"
Cohesion: 0.07
Nodes (27): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+19 more)

### Community 74 - "settings.mock.ts"
Cohesion: 0.09
Nodes (24): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session, web_src_api_model_index_systeminfo (+16 more)

### Community 75 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 76 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 77 - "Connector"
Cohesion: 0.17
Nodes (7): unavailable(), SnapshotSection, Connector, unavailable(), unavailable(), unavailable(), session

### Community 78 - "NewEngine"
Cohesion: 0.17
Nodes (24): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+16 more)

### Community 79 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 80 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 81 - "AuthedUser"
Cohesion: 0.14
Nodes (25): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+17 more)

### Community 82 - "Manager"
Cohesion: 0.15
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 83 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 84 - ".OIDCCallback"
Cohesion: 0.17
Nodes (8): Handler, newOIDCUser(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider, github.com/coreos/go-oidc/v3/oidc.Provider, golang.org/x/oauth2.Config

### Community 85 - "Engine"
Cohesion: 0.16
Nodes (15): Engine, dedupKey(), matchEntities(), matchReason(), seedEngineConnectorWithEntities(), TestGenerateLabTopologyCreatesThenUpdatesInPlace(), TestMatchEntitiesDedupesExactExternalIDDuplicates(), TestMatchEntitiesExternalIDPrecedence() (+7 more)

### Community 86 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 87 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 88 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 89 - "go_pkg_os"
Cohesion: 0.11
Nodes (14): TestSnapshotAttributesRoundTripPostgres(), TestSnapshotAttributesRoundTripSQLite(), testSnapshotWithAttributes(), go_pkg_bufio, go_pkg_flag, go_pkg_github_com_pquerna_otp, go_pkg_github_com_pquerna_otp_totp, go_pkg_github_com_robfig_cron_v3 (+6 more)

### Community 90 - "Register"
Cohesion: 0.15
Nodes (21): init(), init(), init(), init(), init(), init(), init(), init() (+13 more)

### Community 91 - "New"
Cohesion: 0.22
Nodes (19): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires() (+11 more)

### Community 92 - "newRouterDeps"
Cohesion: 0.12
Nodes (19): Config, PermissionChecker, NewHandler(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders() (+11 more)

### Community 93 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 94 - "Store"
Cohesion: 0.13
Nodes (8): fakeStatusChecker, testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 95 - "Handler"
Cohesion: 0.16
Nodes (12): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+4 more)

### Community 96 - "handlers_contract_test.go"
Cohesion: 0.18
Nodes (19): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+11 more)

### Community 97 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 98 - "Deps"
Cohesion: 0.20
Nodes (20): registerListAttentionItems(), changeServiceIDs(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs() (+12 more)

### Community 99 - "oidc_elevate_test.go"
Cohesion: 0.15
Nodes (17): isSafeMethod(), clientTimeout(), IsSafeMethod(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects() (+9 more)

### Community 100 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 101 - "docker_test.go"
Cohesion: 0.11
Nodes (19): generateSelfSignedCert(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout(), TestDoRequestErrorCases(), TestFetchBuildsSectionsFromEndpoints(), TestFetchSurfacesMalformedSystemResponse(), TestFetchToleratesEndpointFailure() (+11 more)

### Community 102 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 103 - "ReportsPage.tsx"
Cohesion: 0.11
Nodes (19): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports (+11 more)

### Community 104 - "MarshalConnectorConfig"
Cohesion: 0.17
Nodes (16): TestBackupExportRedactsSecrets(), TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt() (+8 more)

### Community 105 - "chat/chat.go"
Cohesion: 0.15
Nodes (17): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections() (+9 more)

### Community 106 - "mcp/mcp_test.go"
Cohesion: 0.16
Nodes (8): go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, go_pkg_github_com_wiselabz_wiselabz_internal_chat, changeSummary, connectorSummary, findingSummary

### Community 107 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 108 - "logging_test.go"
Cohesion: 0.18
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 109 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 110 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 111 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 112 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 113 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 114 - "main"
Cohesion: 0.18
Nodes (12): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry (+4 more)

### Community 115 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 116 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 117 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 118 - "NewEngine"
Cohesion: 0.42
Nodes (16): NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists(), TestMatchingConnectorsEmptyAppliesToIsWildcard() (+8 more)

### Community 119 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 120 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 121 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 122 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 123 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 124 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 125 - "lifecycleManager"
Cohesion: 0.16
Nodes (7): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, lifecycleDeps, lifecycleManager

### Community 126 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 127 - "InstanceAdminFromContext"
Cohesion: 0.20
Nodes (7): contextWithShareLink(), Handler, newShareToken(), shareLinkFromContext(), InstanceAdminFromContext(), shareLinkNode, shareLinkScope

### Community 128 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 129 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 130 - "Connector"
Cohesion: 0.17
Nodes (6): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), Connector

### Community 131 - "httpx/retry_test.go"
Cohesion: 0.36
Nodes (13): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+5 more)

### Community 132 - "NotificationRecord"
Cohesion: 0.21
Nodes (6): Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord, Store, scanNotification()

### Community 133 - "data.go"
Cohesion: 0.27
Nodes (14): ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary, JobHealthEntry (+6 more)

### Community 134 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 135 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 136 - "Handler"
Cohesion: 0.24
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 137 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 138 - "diagram.go"
Cohesion: 0.26
Nodes (12): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), relatedEntities(), EntityLink (+4 more)

### Community 139 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 140 - "AuditRecord"
Cohesion: 0.31
Nodes (6): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), AuditRecord

### Community 141 - "Hub"
Cohesion: 0.18
Nodes (5): Hub, github.com/gorilla/websocket.Upgrader, broadcastMsg, Revalidator, ticket

### Community 142 - "NewRegistry"
Cohesion: 0.45
Nodes (12): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+4 more)

### Community 143 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 144 - "Engine"
Cohesion: 0.17
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 145 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 146 - "Handler"
Cohesion: 0.18
Nodes (4): cron.EntryID, Handler, sync/atomic.Bool, ReadyState

### Community 147 - "ListSchemas"
Cohesion: 0.21
Nodes (11): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), Capabilities(), CapabilityDescriptor, ListSchemas(), TestRegisterStubRoundTrips() (+3 more)

### Community 148 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 149 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 150 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 151 - "Handler"
Cohesion: 0.30
Nodes (4): updateUserRequest, Handler, writeUserWriteError(), NoContent()

### Community 152 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 153 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 154 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 155 - "store/backup_test.go"
Cohesion: 0.32
Nodes (11): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+3 more)

### Community 156 - "JobHealthRecord"
Cohesion: 0.29
Nodes (4): JobHealthRecord, Store, scanJobHealth(), fakeHealthStore

### Community 157 - "transform_test.go"
Cohesion: 0.24
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 158 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 160 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 161 - "runbooks/handlers_test.go"
Cohesion: 0.36
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 162 - "ServiceDependency"
Cohesion: 0.22
Nodes (5): ServiceDependency, poolDependencies(), unavailable(), networkDependencies(), Connector

### Community 163 - "newDockerClient"
Cohesion: 0.18
Nodes (11): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+3 more)

### Community 164 - "newSSHDockerClient"
Cohesion: 0.25
Nodes (11): generateSSHHostKey(), startSSHDockerServer(), TestNewSSHDockerClientDialsAndExecutesDialStdio(), TestNewSSHDockerClientRejectsMissingHostKey(), TestNewSSHDockerClientRejectsWrongCredentials(), TestNewSSHDockerClientRejectsWrongHostKey(), TestNewSSHDockerClientSupportsSequentialRequests(), newSSHDockerClient() (+3 more)

### Community 165 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 166 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 167 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 168 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 169 - "compliance_rules_test.go"
Cohesion: 0.31
Nodes (8): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), go_pkg_regexp

### Community 170 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 171 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 172 - "ratelimit.go"
Cohesion: 0.27
Nodes (7): TestRateLimit(), RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 173 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 174 - "log/slog.Logger"
Cohesion: 0.29
Nodes (6): Store, RunDocLockSweep(), runDocLockSweep(), Engine, log/slog.Logger, DocLockRecord

### Community 175 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 177 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 178 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 179 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 180 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 181 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 182 - "net/http.Response"
Cohesion: 0.33
Nodes (6): retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 183 - "time.Time"
Cohesion: 0.25
Nodes (6): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, userStatus

### Community 184 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 185 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 186 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 187 - "pfsense.go"
Cohesion: 0.39
Nodes (6): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName()

### Community 188 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 190 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 191 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 192 - "sync.Mutex"
Cohesion: 0.25
Nodes (4): sync.Mutex, fakeDocRegenerator, fakeNotifier, fakeQualityChecker

### Community 193 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 194 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 195 - "snapshotResponse"
Cohesion: 0.48
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 196 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 197 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 198 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 199 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 200 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 202 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 203 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 204 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 205 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 207 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 208 - "newHandler"
Cohesion: 0.47
Nodes (6): Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationValidation(), TestPostMessageErrors()

### Community 209 - ".Fetch"
Cohesion: 0.47
Nodes (3): TestBuildContainerTableAttributes(), buildContainerTable(), Connector

### Community 210 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 212 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 213 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 214 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 215 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 216 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 217 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 218 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 219 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 220 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 221 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 222 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 223 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 224 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 228 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 229 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **572 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+567 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1295 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `testing.T`, `net/http.Request`, `newTestHandler`, `Handler`, `go_pkg_context`, `NewRegistry`, `Engine`, `git_test.go`, `Handler`, `ServiceSnapshot`, `ErrorWithDetails`, `Handler`, `ConnectorRecord`, `ReportData`, `dispatcher_test.go`, `store/backup_test.go`, `net/http.ResponseWriter`, `DecodeKey`, `ExportToFile`, `RunMigrations`, `time.Duration`, `NewChecker`, `Get`, `response.go`, `NewUser`, `.call`, `Dispatcher`, `time.Time`, `retention/retention_test.go`, `sync.Mutex`, `New`, `rewritePlaceholders`, `NewStore`, `NewEngine`, `newHandler`, `AuthedUser`, `Manager`, `engine_maintenance_test.go`, `Engine`, `Checker`, `newRouterDeps`, `Deps`, `chat/chat.go`, `Handler`, `diagnostics/diagnostics.go`, `docs/handlers_test.go`, `NewEngine`, `testApp`, `lifecycleManager`, `changes/handlers_test.go`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `net/http.Request` to `DecodeKey`, `.UpdateAuthConfig`, `Sanitize`, `Handler`, `Handler`, `logging_test.go`, `Get`, `response.go`, `InstanceAdminFromContext`, `.OIDCCallback`, `ErrorWithDetails`, `Handler`, `Handler`, `routerDeps`, `net/http.ResponseWriter`, `Handler`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Why does `WiseLabz — Architecture & Technical Decisions` connect `WiseLabz — Architecture & Technical Decisions` to `Technology stack`, `0004 — PostgreSQL leader election for background workers`, `react-i18next`, `Authentication design`, `Development workflow`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _572 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021219715956558062 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01937984496124031 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.02311981504147967 - nodes in this community are weakly interconnected._