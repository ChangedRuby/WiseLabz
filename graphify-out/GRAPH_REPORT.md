# Graph Report - agent-a627b7e821e91d6a1  (2026-09-23)

## Corpus Check
- 769 files · ~466,815 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 10, .toml 2, .example 1)

## Summary
- 5732 nodes · 18004 edges · 197 communities (178 shown, 19 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1429 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e8948908`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- cn
- context.Context
- DashboardPage.tsx
- go_pkg_net_http
- go_pkg_strings
- SystemPage.tsx
- ServiceSnapshot
- Button.tsx
- Errorf
- go_pkg_context
- icons.tsx
- App.tsx
- connector/connector.go
- vitest
- newTestHandler
- go_pkg_testing
- RulesPage.tsx
- UsersPage.tsx
- net/http.Request
- net/http.ResponseWriter
- NewEngine
- go_pkg_os
- MarshalConnectorConfig
- SnapshotEntity
- dispatcher_test.go
- package.json
- portainer_test.go
- ChangeDetailPage.tsx
- home_assistant/tables.go
- fixtures.ts
- ServiceDetailPage.tsx
- main
- DocRecord
- WiseLabz — Architecture & Technical Decisions
- NewMalformedResponseError
- IsSecureRequest
- ServicesPage.tsx
- ExportToFile
- rowScanner
- response.go
- RunMigrations
- share_links_test.go
- Hub
- NewChecker
- useRole.ts
- dependencies
- NewStore
- routerDeps
- portainer/tables.go
- Store
- newRouterDeps
- home_assistant_test.go
- router.go
- adguardhome/tables.go
- NewRegistry
- compliance/engine.go
- Connector
- traefik/tables.go
- truenas_test.go
- unifi/tables.go
- Dispatcher
- Configuration & Documentation Backup (Export/Import)
- Config
- .Fetch
- rewritePlaceholders
- Store
- Register
- Connector
- settings.mock.ts
- SuggestRequest
- boolToInt
- handlers.ts
- middleware_test.go
- config_test.go
- testApp
- NewEngine
- unifi_test.go
- timeline.ts
- src/theme.ts
- motion
- Checker
- WiseLabz — Design Contract
- Handler
- net/http.Client
- Connector
- devDependencies
- newTestHandler
- Compare
- export_test.go
- logging.go
- adguardhome_test.go
- scheduler/health_test.go
- ConnectorRecord
- Service
- Store
- chat/chat.go
- Connector
- Runner
- New
- Handler
- Handler
- Handler
- GuardedDialer
- docker_test.go
- diagnostics/diagnostics.go
- ws.ts
- compilerOptions
- AuthMiddleware
- traefik_test.go
- docdiffmodel.ts
- templates_test.go
- system/handlers_test.go
- all.go
- Connector
- gitTarget
- compilerOptions
- changes/handlers_test.go
- vectorCache
- gitFixture
- NotificationRecord
- pagination_contract_test.go
- New
- changes_test.go
- connectors_health_test.go
- newTestHandler
- Contributing to WiseLabz
- config_cmd_test.go
- Engine
- NewService
- Connector
- Store
- transform_test.go
- Decision
- main.tsx
- scripts
- Handler
- newSSHDockerClient
- Store
- sshStdioConn
- Decision
- WiseLabz Connector Guide
- Product
- cursor_pagination_test.go
- .call
- ratelimit.go
- Handler
- Changelog
- mockServiceWorker.js
- ComplianceRuleRecord
- newDockerClient
- retention/retention_test.go
- release-please-config.json
- dashboard/handlers_test.go
- Options
- Cache
- Step by step
- WiseLabz
- newHandler
- dialSSHStdio
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- noopValidatedConnector
- .UpdateAuthConfig
- .GetConnectorUptime
- Store
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- testHandler
- CORS
- routerOperations
- Enforcement Guidelines
- @vitejs/plugin-react
- compose-smoke.sh
- buildDockerTLSConfig
- RetentionSettings
- timeoutError
- MISSING — deferred & future frontend features
- Saved Views
- stubEmbedder
- WiseLabz — v2 Backlog
- fakeDocRegenerator
- fakeQualityChecker
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
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

## Communities (197 total, 19 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (181): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+173 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (149): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+141 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (152): newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider(), TestGetOrInitOIDCProvider() (+144 more)

### Community 3 - "cn"
Cohesion: 0.03
Nodes (84): web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_docs_docs_usegetdocstree, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_deletetemplatestemplateid, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplates (+76 more)

### Community 4 - "context.Context"
Cohesion: 0.04
Nodes (22): sanitizeSessions(), Handler, Connector, Connector, Connector, SnapshotRecord, Store, Store (+14 more)

### Community 5 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (86): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+78 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (39): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection (+31 more)

### Community 7 - "go_pkg_strings"
Cohesion: 0.06
Nodes (38): contextKey, elevationError, dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin() (+30 more)

### Community 8 - "SystemPage.tsx"
Cohesion: 0.04
Nodes (71): web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey, web_src_api_generated_settings_settings_getgetauthconfigquerykey, web_src_api_generated_settings_settings_getgetnotificationsconfigquerykey, web_src_api_generated_settings_settings_postaiconfigtest, web_src_api_generated_settings_settings_postnotificationsconfigtest, web_src_api_generated_settings_settings_putaiconfig (+63 more)

### Community 9 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (24): healthFakeConnector, ServiceSnapshot, agentEnabled(), Connector, Sanitize(), TestSanitize(), changePatternID(), Engine (+16 more)

### Community 10 - "Button.tsx"
Cohesion: 0.05
Nodes (68): web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze, web_src_api_generated_alerts_alerts_postalertsbulksnooze, web_src_api_generated_alerts_alerts_usegetalerts, web_src_api_generated_attention_attention_getgetattentionquerykey (+60 more)

### Community 11 - "Errorf"
Cohesion: 0.07
Nodes (31): Handler, newToken(), sanitizeUser(), setRefreshCookie(), Handler, Handler, NewHandler(), configRequestField() (+23 more)

### Community 12 - "go_pkg_context"
Cohesion: 0.09
Nodes (16): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors (+8 more)

### Community 13 - "icons.tsx"
Cohesion: 0.05
Nodes (57): RFC-3339, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest (+49 more)

### Community 14 - "App.tsx"
Cohesion: 0.05
Nodes (56): Frontend shell & theme (decided 2026-06), react, react-error-boundary, react-i18next, react-router-dom, sonner, setAccessToken(), web_src_api_generated_auth_auth (+48 more)

### Community 15 - "connector/connector.go"
Cohesion: 0.04
Nodes (32): Connector, ollamaEmbedder, openAIEmbedder, isTimeout(), NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders() (+24 more)

### Community 16 - "vitest"
Cohesion: 0.04
Nodes (40): msw, @testing-library/react, vitest, zustand, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage, CommandPalette(), DocDiff() (+32 more)

### Community 17 - "newTestHandler"
Cohesion: 0.07
Nodes (64): AssertMatchesSpec(), loadSpec(), specPath(), actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews() (+56 more)

### Community 18 - "go_pkg_testing"
Cohesion: 0.05
Nodes (24): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), TestLoggablePathMasksShareTokenUnderV1() (+16 more)

### Community 19 - "RulesPage.tsx"
Cohesion: 0.04
Nodes (56): web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules (+48 more)

### Community 20 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (51): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), RefreshFn, setRefreshHandler() (+43 more)

### Community 21 - "net/http.Request"
Cohesion: 0.05
Nodes (21): Handler, Handler, Handler, Handler, Handler, Handler, Handler, Handler (+13 more)

### Community 22 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (26): AuditRecorder, Handler, isWritableField(), validateConfigPushRequest(), applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), capitalize() (+18 more)

### Community 23 - "NewEngine"
Cohesion: 0.09
Nodes (44): NewHandler(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine (+36 more)

### Community 24 - "go_pkg_os"
Cohesion: 0.06
Nodes (37): gitAuth(), Exporter, installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), TestGitAuthSSH(), writeTestKey(), TestSnapshotAttributesRoundTripPostgres() (+29 more)

### Community 25 - "MarshalConnectorConfig"
Cohesion: 0.07
Nodes (37): ProviderConfig, TestDiagnosticsRedactsSecrets(), Handler, primaryProviderConfig(), Config, mask(), DecodeKey(), Decrypt() (+29 more)

### Community 26 - "SnapshotEntity"
Cohesion: 0.10
Nodes (48): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), buildDatasets() (+40 more)

### Community 27 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 28 - "package.json"
Cohesion: 0.04
Nodes (44): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+36 more)

### Community 29 - "portainer_test.go"
Cohesion: 0.07
Nodes (47): catalog(), TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion() (+39 more)

### Community 30 - "ChangeDetailPage.tsx"
Cohesion: 0.06
Nodes (41): Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`) (+33 more)

### Community 31 - "home_assistant/tables.go"
Cohesion: 0.08
Nodes (40): jsonType(), TestAttributeCatalogCoversEmittedKeys(), isTimeout(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+32 more)

### Community 32 - "fixtures.ts"
Cohesion: 0.06
Nodes (42): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+34 more)

### Community 33 - "ServiceDetailPage.tsx"
Cohesion: 0.06
Nodes (40): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+32 more)

### Community 34 - "main"
Cohesion: 0.06
Nodes (39): Config, main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens() (+31 more)

### Community 35 - "DocRecord"
Cohesion: 0.06
Nodes (17): existingIDs(), nilToStr(), docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary() (+9 more)

### Community 36 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 37 - "NewMalformedResponseError"
Cohesion: 0.09
Nodes (38): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+30 more)

### Community 38 - "IsSecureRequest"
Cohesion: 0.08
Nodes (27): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), clearOIDCFlowCookie(), emailDomainAllowed(), Handler, newOIDCUser(), oidcFlowCookieName(), oidcRoleForGroups() (+19 more)

### Community 39 - "ServicesPage.tsx"
Cohesion: 0.06
Nodes (36): match-sorter, @radix-ui/react-popover, @tanstack/react-query, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow, web_src_api_generated_connectors_connectors_getgetconnectorsmaintenancewindowsquerykey, web_src_api_generated_connectors_connectors_postconnectorsbulkreauth (+28 more)

### Community 40 - "ExportToFile"
Cohesion: 0.11
Nodes (41): ExportToFile(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+33 more)

### Community 41 - "rowScanner"
Cohesion: 0.09
Nodes (21): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), countRows(), T, keysetQuery() (+13 more)

### Community 42 - "response.go"
Cohesion: 0.07
Nodes (23): Handler, parseScheduleUpdates(), validateRotationFields(), Handler, TestWritePaginatedOmitsNextCursor(), Error(), DataPaginatedResponse, FieldError (+15 more)

### Community 43 - "RunMigrations"
Cohesion: 0.11
Nodes (37): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+29 more)

### Community 44 - "share_links_test.go"
Cohesion: 0.17
Nodes (41): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld() (+33 more)

### Community 45 - "Hub"
Cohesion: 0.07
Nodes (21): newLifecycleManager(), Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep(), Engine (+13 more)

### Community 46 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 47 - "useRole.ts"
Cohesion: 0.08
Nodes (33): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, i18next, web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, Command, CommandCtx (+25 more)

### Community 48 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 49 - "NewStore"
Cohesion: 0.10
Nodes (37): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+29 more)

### Community 50 - "routerDeps"
Cohesion: 0.11
Nodes (31): routerDeps, PermissionChecker, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+23 more)

### Community 51 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 52 - "Store"
Cohesion: 0.08
Nodes (11): Store, placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+3 more)

### Community 53 - "newRouterDeps"
Cohesion: 0.09
Nodes (35): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+27 more)

### Community 54 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (35): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+27 more)

### Community 55 - "router.go"
Cohesion: 0.10
Nodes (26): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+18 more)

### Community 56 - "adguardhome/tables.go"
Cohesion: 0.15
Nodes (30): statusInfo, unavailable(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+22 more)

### Community 57 - "NewRegistry"
Cohesion: 0.14
Nodes (24): Provider, StatusError, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds() (+16 more)

### Community 58 - "compliance/engine.go"
Cohesion: 0.11
Nodes (30): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+22 more)

### Community 59 - "Connector"
Cohesion: 0.09
Nodes (12): init(), ConfigField, Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), isTimeout(), primaryGatewayName() (+4 more)

### Community 60 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 61 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 62 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 63 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 64 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 65 - "Config"
Cohesion: 0.10
Nodes (20): Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings, BackupSettings, Database, DocExportGitSettings (+12 more)

### Community 66 - ".Fetch"
Cohesion: 0.09
Nodes (15): upstreamDependencies(), TestUpstreamDependenciesAreDedupedAndSorted(), ServiceDependency, WantsField(), environmentDependencies(), isTimeout(), putMetadata(), unavailable() (+7 more)

### Community 67 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 68 - "Store"
Cohesion: 0.18
Nodes (25): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+17 more)

### Community 69 - "Register"
Cohesion: 0.12
Nodes (26): init(), init(), init(), init(), init(), init(), init(), init() (+18 more)

### Community 70 - "Connector"
Cohesion: 0.17
Nodes (8): SnapshotSection, unavailable(), Connector, isTimeout(), parseHosts(), unavailable(), unavailable(), session

### Community 71 - "settings.mock.ts"
Cohesion: 0.09
Nodes (25): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+17 more)

### Community 72 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 73 - "boolToInt"
Cohesion: 0.11
Nodes (10): Store, scanBackupRun(), Store, Store, RunbookRecord, Store, scanRunbook(), boolToInt() (+2 more)

### Community 74 - "handlers.ts"
Cohesion: 0.07
Nodes (27): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+19 more)

### Community 75 - "middleware_test.go"
Cohesion: 0.11
Nodes (23): ConnectorRoleChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, SecurityHeaders(), TestSecurityHeaders(), RequireConnectorRole(), RequireElevation() (+15 more)

### Community 76 - "config_test.go"
Cohesion: 0.10
Nodes (24): Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH(), TestLoadFromYAML() (+16 more)

### Community 77 - "testApp"
Cohesion: 0.15
Nodes (17): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz(), testApp (+9 more)

### Community 78 - "NewEngine"
Cohesion: 0.18
Nodes (23): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+15 more)

### Community 79 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 80 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 81 - "src/theme.ts"
Cohesion: 0.15
Nodes (23): @fontsource/space-mono, @fontsource-variable/space-grotesk, ColorMode, commit(), load(), Persisted, PRESETS_FONTS, ThemeState (+15 more)

### Community 82 - "motion"
Cohesion: 0.14
Nodes (21): motion, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 83 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 84 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 85 - "Handler"
Cohesion: 0.17
Nodes (11): diffToSpec(), Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes(), TestCursorRoundTrip() (+3 more)

### Community 86 - "net/http.Client"
Cohesion: 0.11
Nodes (6): isTimeout(), Connector, Connector, newWebhookClient(), Connector, net/http.Client

### Community 87 - "Connector"
Cohesion: 0.15
Nodes (9): apiMessage(), controllerName(), countByKind(), isTimeout(), statusError(), unavailable(), Connector, sectionFetch (+1 more)

### Community 88 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 89 - "newTestHandler"
Cohesion: 0.21
Nodes (18): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+10 more)

### Community 90 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 91 - "export_test.go"
Cohesion: 0.19
Nodes (19): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+11 more)

### Community 92 - "logging.go"
Cohesion: 0.16
Nodes (17): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+9 more)

### Community 93 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 94 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 95 - "ConnectorRecord"
Cohesion: 0.16
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 96 - "Service"
Cohesion: 0.20
Nodes (10): Claims, ElevationClaims, ElevationToken, TokenPair, Service, hasAudience(), newTokenID(), go_pkg_github_com_golang_jwt_jwt_v5 (+2 more)

### Community 97 - "Store"
Cohesion: 0.12
Nodes (7): fakeStatusChecker, testAPIKeyChecker, sanitize(), APIKeyClaims, validAPIKey(), APIKey, Store

### Community 98 - "chat/chat.go"
Cohesion: 0.16
Nodes (18): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings() (+10 more)

### Community 99 - "Connector"
Cohesion: 0.14
Nodes (8): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), isTimeout(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 100 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 101 - "New"
Cohesion: 0.14
Nodes (19): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+11 more)

### Community 102 - "Handler"
Cohesion: 0.17
Nodes (11): updateUserRequest, Handler, writeUserWriteError(), mustHashDummyPassword(), HashPassword(), TestHashPasswordRejectsOver72Bytes(), TestHashAndVerify(), TestHashTooShort() (+3 more)

### Community 103 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 104 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 105 - "GuardedDialer"
Cohesion: 0.14
Nodes (13): newConnector(), Connector, GuardedDialer(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector() (+5 more)

### Community 106 - "docker_test.go"
Cohesion: 0.11
Nodes (17): TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout(), TestDoRequestErrorCases(), TestFetchBuildsSectionsFromEndpoints(), TestFetchSurfacesMalformedSystemResponse(), TestFetchToleratesEndpointFailure(), TestFetchWithFieldsHintSkipsUnrequestedCalls() (+9 more)

### Community 107 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 108 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 109 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 110 - "AuthMiddleware"
Cohesion: 0.15
Nodes (13): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), AuthMiddleware(), extractBearerToken() (+5 more)

### Community 111 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 112 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 113 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 114 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 115 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 116 - "Connector"
Cohesion: 0.16
Nodes (7): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), isTimeout(), Connector

### Community 117 - "gitTarget"
Cohesion: 0.19
Nodes (9): commitMessage(), TestCommitMessage(), commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash, github.com/go-git/go-git/v5/plumbing/object.Signature, github.com/go-git/go-git/v5/plumbing.ReferenceName (+1 more)

### Community 118 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 119 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 120 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 121 - "gitFixture"
Cohesion: 0.29
Nodes (9): SetBeforePushForTest(), keys(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, Exporter (+1 more)

### Community 122 - "NotificationRecord"
Cohesion: 0.20
Nodes (7): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), NotificationRecord, Store, scanNotification()

### Community 123 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 124 - "New"
Cohesion: 0.36
Nodes (13): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+5 more)

### Community 125 - "changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 126 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 127 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 128 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 129 - "config_cmd_test.go"
Cohesion: 0.26
Nodes (10): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+2 more)

### Community 130 - "Engine"
Cohesion: 0.18
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 131 - "NewService"
Cohesion: 0.30
Nodes (11): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+3 more)

### Community 133 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 134 - "transform_test.go"
Cohesion: 0.24
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 135 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 136 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 137 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 139 - "newSSHDockerClient"
Cohesion: 0.25
Nodes (11): generateSSHHostKey(), startSSHDockerServer(), TestNewSSHDockerClientDialsAndExecutesDialStdio(), TestNewSSHDockerClientRejectsMissingHostKey(), TestNewSSHDockerClientRejectsWrongCredentials(), TestNewSSHDockerClientRejectsWrongHostKey(), TestNewSSHDockerClientSupportsSequentialRequests(), newSSHDockerClient() (+3 more)

### Community 140 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 141 - "sshStdioConn"
Cohesion: 0.20
Nodes (5): sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session, io.WriteCloser, net.Addr

### Community 142 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 143 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 144 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 145 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 146 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 147 - "ratelimit.go"
Cohesion: 0.27
Nodes (7): TestRateLimit(), RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 149 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 150 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 151 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 152 - "newDockerClient"
Cohesion: 0.22
Nodes (9): IsDangerousIP(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+1 more)

### Community 153 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 154 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 155 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 156 - "Options"
Cohesion: 0.25
Nodes (7): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn, Options

### Community 157 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 158 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 159 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 160 - "newHandler"
Cohesion: 0.43
Nodes (7): Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationDocVisibility(), TestCreateConversationValidation(), TestPostMessageErrors()

### Community 161 - "dialSSHStdio"
Cohesion: 0.29
Nodes (5): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), golang.org/x/crypto/ssh.ClientConfig, io.Closer

### Community 162 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 163 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 164 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 165 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 166 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 167 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 170 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 172 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 173 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 174 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 175 - "testHandler"
Cohesion: 0.40
Nodes (3): testHandler, Handler, instanceAdminRoleFor()

### Community 176 - "CORS"
Cohesion: 0.60
Nodes (4): CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders()

### Community 177 - "routerOperations"
Cohesion: 0.50
Nodes (5): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes

### Community 178 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 179 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 180 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 181 - "buildDockerTLSConfig"
Cohesion: 0.50
Nodes (4): buildDockerTLSConfig(), generateSelfSignedCert(), TestNewTCPDockerClientMutualTLS(), crypto/tls.Config

### Community 185 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 186 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **540 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+535 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1176 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `UserIDFromContext()` connect `Errorf` to `context.Context`, `go_pkg_strings`, `rowScanner`, `Handler`, `middleware_test.go`, `response.go`, `routerDeps`, `Handler`, `Handler`, `net/http.ResponseWriter`, `net/http.Request`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `Engine`, `testing.T`, `ServiceSnapshot`, `Handler`, `Errorf`, `go_pkg_context`, `.call`, `Handler`, `net/http.Request`, `net/http.ResponseWriter`, `NewEngine`, `retention/retention_test.go`, `dispatcher_test.go`, `newHandler`, `main`, `ExportToFile`, `rowScanner`, `response.go`, `RunMigrations`, `share_links_test.go`, `engine_maintenance_test.go`, `NewChecker`, `testHandler`, `Hub`, `NewStore`, `newRouterDeps`, `NewRegistry`, `Dispatcher`, `rewritePlaceholders`, `testApp`, `NewEngine`, `Checker`, `Handler`, `export_test.go`, `chat/chat.go`, `New`, `Handler`, `Handler`, `diagnostics/diagnostics.go`, `changes/handlers_test.go`, `gitFixture`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `gitFixture` to `testing.T`, `Store`, `context.Context`, `Hub`, `go_pkg_os`, `export_test.go`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _540 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.01983839220967575 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.022782037239868564 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02147239263803681 - nodes in this community are weakly interconnected._