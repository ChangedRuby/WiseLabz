# Graph Report - feat-ws-filter-broadcast-events-by-per-connector  (2026-09-29)

## Corpus Check
- 879 files · ~552,320 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6597 nodes · 20944 edges · 241 communities (223 shown, 18 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1680 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `74e00975`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newDocTestStore
- react
- testing.T
- newTestApp
- context.Context
- DashboardPage.tsx
- go_pkg_net_http
- @tanstack/react-query
- NotificationsPage.tsx
- go_pkg_context
- go_pkg_strings
- newTestHandler
- Errorf
- go_pkg_testing
- net/http.Client
- cn
- App.tsx
- icons.tsx
- ServiceDetailPage.tsx
- ServiceSnapshot
- ErrorWithDetails
- ProfilePage.tsx
- export_test.go
- ShareLinkPage.tsx
- Hub
- net/http.ResponseWriter
- ConnectorRecord
- package.json
- DecodeKey
- Store
- Compare
- fixtures.ts
- Get
- dispatcher_test.go
- net/http.Request
- SnapshotEntity
- SystemPage.tsx
- RunMigrations
- DocRecord
- Connector
- api/audit_test.go
- relativeTime
- dependencies
- home_assistant_test.go
- home_assistant/tables.go
- WebSocketProvider.tsx
- NewChecker
- lists.go
- Manager
- portainer/tables.go
- Dispatcher
- Store
- NewUser
- connector/connector.go
- portainer_test.go
- adguardhome/tables.go
- newTestHandler
- GetTypeSchema
- docker_test.go
- traefik/tables.go
- unifi/tables.go
- main
- compliance/engine.go
- newTestHarness
- server/lifecycle.go
- config_test.go
- Connector
- New
- Connector
- NewStore
- mountAPIRoutes
- NewEngine
- settings.mock.ts
- SuggestRequest
- AuthMiddleware
- rewritePlaceholders
- handlers.ts
- runbooks_test.go
- response.go
- Service
- router.go
- Config
- middleware.go
- unifi_test.go
- src/theme.ts
- timeline.ts
- Store
- ExportToFile
- Checker
- Config
- compliance/handlers.go
- VerifyBundleFile
- NewMalformedResponseError
- Engine
- WiseLabz — Design Contract
- devDependencies
- NewHTTPClient
- git.go
- adguardhome_test.go
- scheduler/health_test.go
- Registry
- Handler
- NewService
- backup/backup.go
- Runner
- WiseLabz — Architecture & Technical Decisions
- ReportsPage.tsx
- ws.ts
- lifecycle_test.go
- .SnapshotDiff
- logging_test.go
- Handler
- truenas_test.go
- diagnostics/diagnostics.go
- keyset_test.go
- compilerOptions
- Register
- docs/handlers_test.go
- diagram.go
- sshStdioConn
- Connector
- traefik_test.go
- NewEngine
- RunbookRecord
- docdiffmodel.ts
- templates_test.go
- vectorCache
- all.go
- httpx/retry_test.go
- Store
- compilerOptions
- testApp
- AuthedUser
- handlers_contract_test.go
- compliance_rules_test.go
- api/mcp_test.go
- newDockerClient
- NotificationRecord
- data.go
- New
- log/slog.Logger
- apikey_scopes_test.go
- changes/handlers_test.go
- pagination_contract_test.go
- reports/handlers.go
- newTestHandler
- render_test.go
- AuditRecord
- Contributing to WiseLabz
- api/changes_test.go
- chat/chat.go
- Handler
- retention/retention_test.go
- Decision
- 0004 — PostgreSQL leader election for background workers
- scripts
- config_cmd_test.go
- NewRegistry
- time.Duration
- IsSecureRequest
- ReportData
- Store
- store/backup_test.go
- Engine
- Contributor Covenant Code of Conduct
- Decision
- Decision
- Backend test performance
- main.tsx
- runbooks/handlers_test.go
- httpx/retry.go
- Connector
- Decision
- WiseLabz Connector Guide
- Product
- .call
- cursor_pagination_test.go
- handlers_bulk_test.go
- connectors_maintenance_test.go
- .call
- Elector
- Changelog
- test-shards.sh
- mockServiceWorker.js
- openapi_contract_test.go
- time.Time
- BACKUP.md
- release-please-config.json
- auth/handlers_test.go
- snapshots_test.go
- RateLimit
- internal/auth/mfa.go
- validate.go
- ComputeWindow
- Store
- runbook_test.go
- Cache
- Step by step
- WiseLabz
- RequirePermission
- handlers_start_stop_configpush_test.go
- .UpdateAuthConfig
- scanMaintenanceWindow
- computeNextRun
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- spec.go
- DocVersionRecord
- Store
- engine_maintenance_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- WiseLabz — Deployment Guide
- Scheduled Doc Export
- Security Policy
- fakeDocRegenerator
- seedScopeFixture
- Authentication design
- Development workflow
- MfaEnrollDialog
- @vitejs/plugin-react
- compose-smoke.sh
- mountPublicSystemRoutes
- MISSING — deferred & future frontend features
- Saved Views
- Fixture reuse and lifecycle tests (#405–#407)
- internal/auth/oidc.go
- fields_test.go
- WiseLabz — v2 Backlog
- fakeEmbedder
- coverage-parity.sh
- tsconfig.json
- AGENTS.md
- versionSections
- setup-env.sh
- CHANGE_PROVENANCE.md
- coverpkg.sh
- vite-env.d.ts
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

## Communities (241 total, 18 thin omitted)

### Community 0 - "newDocTestStore"
Cohesion: 0.02
Nodes (151): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+143 more)

### Community 1 - "react"
Cohesion: 0.04
Nodes (123): RFC-3339, match-sorter, motion, @radix-ui/react-popover, react, react-i18next, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve (+115 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (143): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+135 more)

### Community 3 - "newTestApp"
Cohesion: 0.03
Nodes (130): TestAPIKeyCreateRejectsInvalidExpiryAndEmptyName(), TestAPIKeyRoutesEndToEnd(), testApp, loginRefreshCookie(), seedLocalUser(), TestChangePasswordWrongCurrentPassword(), TestDeleteSessionNotOwner(), TestDeleteSessionSuccess() (+122 more)

### Community 4 - "context.Context"
Cohesion: 0.03
Nodes (31): fakeStatusChecker, sanitizeSessions(), Connector, Connector, Store, scanComplianceRule(), SnapshotRecord, Store (+23 more)

### Community 5 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (102): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete` (+94 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.05
Nodes (54): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, updateUserRequest, newToken(), mustHashDummyPassword(), validHostPort(), changePromptData() (+46 more)

### Community 7 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (66): Frontend shell & theme (decided 2026-06), i18next, msw, react-error-boundary, sonner, @tanstack/react-query, @testing-library/react, vitest (+58 more)

### Community 8 - "NotificationsPage.tsx"
Cohesion: 0.03
Nodes (79): zustand, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridremovalimpact, web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey (+71 more)

### Community 9 - "go_pkg_context"
Cohesion: 0.08
Nodes (19): StatusError, ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), contains(), searchString(), go_pkg_context, go_pkg_database_sql (+11 more)

### Community 10 - "go_pkg_strings"
Cohesion: 0.05
Nodes (41): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), dateFormat(), filterByTitle(), join(), TestDateFormat() (+33 more)

### Community 11 - "newTestHandler"
Cohesion: 0.06
Nodes (73): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+65 more)

### Community 12 - "Errorf"
Cohesion: 0.05
Nodes (29): Handler, sanitize(), Handler, Handler, Handler, Handler, Handler, Handler (+21 more)

### Community 13 - "go_pkg_testing"
Cohesion: 0.04
Nodes (36): CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), TestSchemaMatchesConfig(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6() (+28 more)

### Community 14 - "net/http.Client"
Cohesion: 0.04
Nodes (29): Connector, ollamaEmbedder, openAIEmbedder, NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), tryParseEntities(), validateCustomURL() (+21 more)

### Community 15 - "cn"
Cohesion: 0.04
Nodes (63): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid (+55 more)

### Community 16 - "App.tsx"
Cohesion: 0.04
Nodes (61): react-router-dom, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthloginmfa, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh (+53 more)

### Community 17 - "icons.tsx"
Cohesion: 0.05
Nodes (65): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey (+57 more)

### Community 18 - "ServiceDetailPage.tsx"
Cohesion: 0.04
Nodes (67): ADR-0001, ADR-0003, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+59 more)

### Community 19 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (19): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, init(), RegisterTransformer(), runTransformers(), TestRunTransformersAppliesInOrderAndStopsOnError(), TestRunTransformersUnknownCategoryIsNoop() (+11 more)

### Community 20 - "ErrorWithDetails"
Cohesion: 0.07
Nodes (26): Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, Handler, Handler, newOIDCUser() (+18 more)

### Community 21 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (53): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+45 more)

### Community 22 - "export_test.go"
Cohesion: 0.07
Nodes (45): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+37 more)

### Community 23 - "ShareLinkPage.tsx"
Cohesion: 0.06
Nodes (46): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+38 more)

### Community 24 - "Hub"
Cohesion: 0.07
Nodes (43): TestBroadcastDocEventScoping(), Envelope, Hub, newEnvelope(), newHeartbeat(), NewHub(), normalizeOrigin(), assertEnvelope() (+35 more)

### Community 25 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (22): Handler, isWritableField(), validateConfigPushRequest(), Handler, decodeBulkRequest(), Handler, Handler, loggablePath() (+14 more)

### Community 26 - "ConnectorRecord"
Cohesion: 0.07
Nodes (29): enableFakeEmbedding(), changeFilterClause(), ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr() (+21 more)

### Community 27 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 28 - "DecodeKey"
Cohesion: 0.08
Nodes (28): ProviderConfig, factorJSON(), Handler, Handler, Handler, primaryProviderConfig(), GenerateTOTPSecret(), NormalizeRecoveryCode() (+20 more)

### Community 29 - "Store"
Cohesion: 0.07
Nodes (15): seedAlert(), changeServiceIDs(), seedChange(), placeholders(), AlertRecord, ChangeRecord, Store, scanAlert() (+7 more)

### Community 30 - "Compare"
Cohesion: 0.07
Nodes (45): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+37 more)

### Community 31 - "fixtures.ts"
Cohesion: 0.06
Nodes (44): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+36 more)

### Community 32 - "Get"
Cohesion: 0.07
Nodes (28): applyConnectorScalarUpdates(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection(), LifecycleOp(), Get() (+20 more)

### Community 33 - "dispatcher_test.go"
Cohesion: 0.21
Nodes (45): NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore(), setChannelAndRoutingConfig(), setChannelConfig(), setChannelConfigJSON() (+37 more)

### Community 34 - "net/http.Request"
Cohesion: 0.09
Nodes (19): Handler, oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName() (+11 more)

### Community 35 - "SnapshotEntity"
Cohesion: 0.13
Nodes (42): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+34 more)

### Community 36 - "SystemPage.tsx"
Cohesion: 0.06
Nodes (36): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey (+28 more)

### Community 37 - "RunMigrations"
Cohesion: 0.10
Nodes (39): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns() (+31 more)

### Community 38 - "DocRecord"
Cohesion: 0.06
Nodes (18): seedDelivery(), existingIDs(), nilToStr(), docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc() (+10 more)

### Community 39 - "Connector"
Cohesion: 0.06
Nodes (13): init(), ConfigField, Connector, buildRouteTable(), TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName() (+5 more)

### Community 40 - "api/audit_test.go"
Cohesion: 0.07
Nodes (40): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+32 more)

### Community 41 - "relativeTime"
Cohesion: 0.07
Nodes (37): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, RoleGate(), RoleGateProps, TimeAgo(), TimeAgoProps (+29 more)

### Community 42 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 43 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (39): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+31 more)

### Community 44 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 45 - "WebSocketProvider.tsx"
Cohesion: 0.08
Nodes (35): 15. `system.resync`, Client dispatch model, Delivery, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport (+27 more)

### Community 46 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 47 - "lists.go"
Cohesion: 0.11
Nodes (35): buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames(), parseAdlistsV6() (+27 more)

### Community 48 - "Manager"
Cohesion: 0.09
Nodes (13): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, Store, ReportDefinitionRecord (+5 more)

### Community 49 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 50 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 51 - "Store"
Cohesion: 0.10
Nodes (20): routerDeps, Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+12 more)

### Community 52 - "NewUser"
Cohesion: 0.22
Nodes (36): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestCreateConversationDocVisibility(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser() (+28 more)

### Community 53 - "connector/connector.go"
Cohesion: 0.06
Nodes (24): countLifecycle(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP() (+16 more)

### Community 54 - "portainer_test.go"
Cohesion: 0.10
Nodes (35): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+27 more)

### Community 55 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 56 - "newTestHandler"
Cohesion: 0.13
Nodes (34): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+26 more)

### Community 57 - "GetTypeSchema"
Cohesion: 0.09
Nodes (31): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion() (+23 more)

### Community 58 - "docker_test.go"
Cohesion: 0.08
Nodes (34): generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn(), startSSHDockerServer(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout() (+26 more)

### Community 59 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 60 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 61 - "main"
Cohesion: 0.09
Nodes (26): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+18 more)

### Community 62 - "compliance/engine.go"
Cohesion: 0.12
Nodes (28): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+20 more)

### Community 63 - "newTestHarness"
Cohesion: 0.15
Nodes (30): registerListAttentionItems(), TestListAttentionItems(), registerListChanges(), TestListChanges(), jsonResult(), registerListConnectors(), TestListConnectors(), registerSearchDocs() (+22 more)

### Community 64 - "server/lifecycle.go"
Cohesion: 0.08
Nodes (14): newLifecycleManager(), Election, go_pkg_github_com_wiselabz_wiselabz_internal_leader, go_pkg_github_com_wiselabz_wiselabz_internal_notifications, go_pkg_golang_org_x_sync_errgroup, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server (+6 more)

### Community 65 - "config_test.go"
Cohesion: 0.09
Nodes (27): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+19 more)

### Community 66 - "Connector"
Cohesion: 0.11
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 67 - "New"
Cohesion: 0.10
Nodes (28): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+20 more)

### Community 68 - "Connector"
Cohesion: 0.16
Nodes (10): SnapshotSection, TestBuildHostsTableV5(), buildHostsTable(), Connector, parseHosts(), TestBuildHostsTableMalformedCases(), TestBuildHostsTableValidRecords(), unavailable() (+2 more)

### Community 69 - "NewStore"
Cohesion: 0.15
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 70 - "mountAPIRoutes"
Cohesion: 0.12
Nodes (23): chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router, mountChatRoutes() (+15 more)

### Community 71 - "NewEngine"
Cohesion: 0.17
Nodes (25): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+17 more)

### Community 72 - "settings.mock.ts"
Cohesion: 0.09
Nodes (25): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+17 more)

### Community 73 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 74 - "AuthMiddleware"
Cohesion: 0.11
Nodes (22): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, AuthMiddleware(), RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel() (+14 more)

### Community 75 - "rewritePlaceholders"
Cohesion: 0.11
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 76 - "handlers.ts"
Cohesion: 0.07
Nodes (27): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+19 more)

### Community 77 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 78 - "response.go"
Cohesion: 0.10
Nodes (17): Handler, TestWritePaginatedOmitsNextCursor(), Error(), DataPaginatedResponse, HandleStoreError(), intQuery(), JSON(), Logger() (+9 more)

### Community 79 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 80 - "router.go"
Cohesion: 0.13
Nodes (24): validVerb(), go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+16 more)

### Community 81 - "Config"
Cohesion: 0.12
Nodes (22): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+14 more)

### Community 82 - "middleware.go"
Cohesion: 0.10
Nodes (19): AuditRecorder, ConnectorRoleChecker, contextKey, elevationError, testAPIKeyChecker, UserStatusChecker, APIKeyClaims, APIKeyIDFromContext() (+11 more)

### Community 83 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 84 - "src/theme.ts"
Cohesion: 0.14
Nodes (24): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+16 more)

### Community 85 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 86 - "Store"
Cohesion: 0.13
Nodes (13): APIKeyRestriction, auditConnectorGrantDiffJSON(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), getConnectorGrant(), ConnectorGrantDiff, Store (+5 more)

### Community 87 - "ExportToFile"
Cohesion: 0.16
Nodes (23): Export(), ExportToFile(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory() (+15 more)

### Community 88 - "Checker"
Cohesion: 0.19
Nodes (8): Snapshot, complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 89 - "Config"
Cohesion: 0.12
Nodes (24): Config, TestEmbeddedSPAWithoutFrontendBuild(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), Handler (+16 more)

### Community 90 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 91 - "VerifyBundleFile"
Cohesion: 0.16
Nodes (23): BuildManifest(), BundleCounts(), ReadManifest(), WriteManifest(), failVerification(), LatestBundle(), newScratchStore(), RunVerifyOnce() (+15 more)

### Community 92 - "NewMalformedResponseError"
Cohesion: 0.14
Nodes (11): NewMalformedResponseError(), WantsField(), TestBuildContainerTableAttributes(), buildContainerTable(), Connector, putMetadata(), unavailable(), agentEnabled() (+3 more)

### Community 93 - "Engine"
Cohesion: 0.16
Nodes (15): Engine, dedupKey(), matchEntities(), matchReason(), seedEngineConnectorWithEntities(), TestGenerateLabTopologyCreatesThenUpdatesInPlace(), TestMatchEntitiesDedupesExactExternalIDDuplicates(), TestMatchEntitiesExternalIDPrecedence() (+7 more)

### Community 94 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 95 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 96 - "NewHTTPClient"
Cohesion: 0.11
Nodes (17): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), unavailable() (+9 more)

### Community 97 - "git.go"
Cohesion: 0.12
Nodes (19): TestCommitMessage(), keys(), writeExportState(), exportCursor, exportState, go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_go_git_go_git_v5 (+11 more)

### Community 98 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 99 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 100 - "Registry"
Cohesion: 0.16
Nodes (13): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+5 more)

### Community 101 - "Handler"
Cohesion: 0.18
Nodes (6): webAuthnFlow, webAuthnUser, Handler, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/go-webauthn/webauthn/webauthn.SessionData, github.com/google/uuid.UUID

### Community 102 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 103 - "backup/backup.go"
Cohesion: 0.25
Nodes (19): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+11 more)

### Community 104 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 105 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.10
Nodes (20): ADR index, AI module, API design, Backend, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27) (+12 more)

### Community 106 - "ReportsPage.tsx"
Cohesion: 0.11
Nodes (19): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports (+11 more)

### Community 107 - "ws.ts"
Cohesion: 0.11
Nodes (18): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+10 more)

### Community 108 - "lifecycle_test.go"
Cohesion: 0.24
Nodes (10): newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestLifecycleManagerWaitsForSchedulerBeforeCancelAndDBClose(), TestStandbyIsUnreadyAndRunsNoScheduler(), waitForLifecycleSignal() (+2 more)

### Community 109 - ".SnapshotDiff"
Cohesion: 0.20
Nodes (12): decodeStoredSnapshot(), Handler, snapshotStoreError(), Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor() (+4 more)

### Community 110 - "logging_test.go"
Cohesion: 0.18
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 111 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 112 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 113 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 114 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 115 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 116 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 117 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 118 - "diagram.go"
Cohesion: 0.20
Nodes (15): ServiceDependency, environmentDependencies(), networkDependencies(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid() (+7 more)

### Community 119 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 120 - "Connector"
Cohesion: 0.18
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 121 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 122 - "NewEngine"
Cohesion: 0.42
Nodes (16): NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists(), TestMatchingConnectorsEmptyAppliesToIsWildcard() (+8 more)

### Community 123 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 124 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 125 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 126 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 127 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 128 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 129 - "Store"
Cohesion: 0.16
Nodes (7): Store, ChatConversationRecord, Store, apiKeyConnectorFilter(), AttentionItem, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 130 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 131 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 132 - "AuthedUser"
Cohesion: 0.21
Nodes (15): TestCreate(), AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+7 more)

### Community 133 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 134 - "compliance_rules_test.go"
Cohesion: 0.17
Nodes (12): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), isSafeMethod(), treatAsSafeFromContext() (+4 more)

### Community 135 - "api/mcp_test.go"
Cohesion: 0.18
Nodes (8): go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, go_pkg_github_com_wiselabz_wiselabz_internal_chat, changeSummary, connectorSummary, findingSummary

### Community 136 - "newDockerClient"
Cohesion: 0.13
Nodes (14): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+6 more)

### Community 137 - "NotificationRecord"
Cohesion: 0.21
Nodes (6): Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord, Store, scanNotification()

### Community 138 - "data.go"
Cohesion: 0.27
Nodes (14): ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary, JobHealthEntry (+6 more)

### Community 139 - "New"
Cohesion: 0.34
Nodes (14): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+6 more)

### Community 140 - "log/slog.Logger"
Cohesion: 0.21
Nodes (10): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), Store, RunDocLockSweep(), runDocLockSweep(), Engine (+2 more)

### Community 141 - "apikey_scopes_test.go"
Cohesion: 0.25
Nodes (13): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer() (+5 more)

### Community 142 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 143 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 144 - "reports/handlers.go"
Cohesion: 0.30
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 145 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 146 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 147 - "AuditRecord"
Cohesion: 0.31
Nodes (6): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), AuditRecord

### Community 148 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 149 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 150 - "chat/chat.go"
Cohesion: 0.18
Nodes (12): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+4 more)

### Community 151 - "Handler"
Cohesion: 0.29
Nodes (5): stepAuditDetail(), validTargetType(), Handler, runbookResponse, stepResponse

### Community 152 - "retention/retention_test.go"
Cohesion: 0.35
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 153 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 154 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.15
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 155 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 156 - "config_cmd_test.go"
Cohesion: 0.26
Nodes (10): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+2 more)

### Community 157 - "NewRegistry"
Cohesion: 0.50
Nodes (11): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+3 more)

### Community 158 - "time.Duration"
Cohesion: 0.21
Nodes (5): wsRead(), Database, Server, time.Duration, PoolConfig

### Community 159 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 160 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 161 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 162 - "store/backup_test.go"
Cohesion: 0.32
Nodes (11): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+3 more)

### Community 163 - "Engine"
Cohesion: 0.18
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 164 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 165 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 166 - "Decision"
Cohesion: 0.17
Nodes (12): 0005 — Cross-replica WebSocket event relay for active/active, Authorization and secrets, Consequences, Context, Decision, Duplicate suppression and ordering, Event ownership: local first, then relay, Mechanism: PostgreSQL LISTEN/NOTIFY (+4 more)

### Community 167 - "Backend test performance"
Cohesion: 0.17
Nodes (12): Backend test performance, CI job times, Coverage strategy, Deterministic scheduled exports (#408), Follow-ups, Measurements and validation, Measuring, Rules for new tests (+4 more)

### Community 168 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 169 - "runbooks/handlers_test.go"
Cohesion: 0.36
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 170 - "httpx/retry.go"
Cohesion: 0.33
Nodes (7): idempotent(), retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 172 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 173 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 174 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 175 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 176 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 177 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 178 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 179 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 180 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 181 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 182 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 183 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 184 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 185 - "time.Time"
Cohesion: 0.25
Nodes (6): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, userStatus

### Community 186 - "BACKUP.md"
Cohesion: 0.25
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 187 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 188 - "auth/handlers_test.go"
Cohesion: 0.25
Nodes (7): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups()

### Community 189 - "snapshots_test.go"
Cohesion: 0.57
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 190 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 191 - "internal/auth/mfa.go"
Cohesion: 0.32
Nodes (6): GenerateRecoveryCodes(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestNormalizeRecoveryCode(), go_pkg_github_com_pquerna_otp, go_pkg_github_com_pquerna_otp_totp

### Community 192 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 193 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 194 - "Store"
Cohesion: 0.32
Nodes (3): Store, scanBackupRun(), BackupRun

### Community 195 - "runbook_test.go"
Cohesion: 0.32
Nodes (7): Store, newCascadeTestStore(), TestGetRunbookByTarget(), TestRunbookRoundTrip(), TestRunbookStepsCascadeOnConnectorDelete(), TestRunbookStepsCascadeOnRunbookDelete(), TestRunbookStepsRoundTrip()

### Community 196 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 197 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 198 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 199 - "RequirePermission"
Cohesion: 0.38
Nodes (5): PermissionChecker, chi.Router, mountDashboardRoutes(), mountWorkflowRoutes(), RequirePermission()

### Community 200 - "handlers_start_stop_configpush_test.go"
Cohesion: 0.48
Nodes (6): createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler()

### Community 201 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 202 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 203 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 204 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 205 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 206 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 207 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 208 - "spec.go"
Cohesion: 0.33
Nodes (5): go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter, go_pkg_github_com_getkin_kin_openapi_routers, go_pkg_github_com_getkin_kin_openapi_routers_gorillamux, go_pkg_runtime

### Community 211 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 212 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 213 - "WiseLabz — Deployment Guide"
Cohesion: 0.33
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 214 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 215 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 217 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 219 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 220 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 221 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 222 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 223 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 224 - "mountPublicSystemRoutes"
Cohesion: 0.67
Nodes (3): chi.Router, mountPublicSystemRoutes(), mountRootRoutes()

### Community 225 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 226 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 227 - "Fixture reuse and lifecycle tests (#405–#407)"
Cohesion: 0.50
Nodes (4): CI measurements, Fixture reuse and lifecycle tests (#405–#407), Local measurements, Validation and coverage

## Knowledge Gaps
- **596 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+591 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1319 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `WebSocketProvider()` connect `WebSocketProvider.tsx` to `App.tsx`, `DashboardPage.tsx`, `MfaEnrollDialog`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `testApp`, `AuthedUser`, `go_pkg_context`, `newTestHandler`, `Errorf`, `log/slog.Logger`, `reports/handlers.go`, `ErrorWithDetails`, `export_test.go`, `Handler`, `retention/retention_test.go`, `net/http.ResponseWriter`, `ConnectorRecord`, `NewRegistry`, `Store`, `Get`, `dispatcher_test.go`, `net/http.Request`, `ReportData`, `store/backup_test.go`, `RunMigrations`, `DocRecord`, `Engine`, `NewChecker`, `.call`, `Manager`, `Dispatcher`, `.call`, `NewUser`, `time.Time`, `newTestHarness`, `server/lifecycle.go`, `New`, `NewStore`, `NewEngine`, `rewritePlaceholders`, `response.go`, `engine_maintenance_test.go`, `ExportToFile`, `Checker`, `Config`, `compliance/handlers.go`, `VerifyBundleFile`, `Engine`, `backup/backup.go`, `diagnostics/diagnostics.go`, `docs/handlers_test.go`, `NewEngine`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `WiseLabz Connector Guide` connect `WiseLabz Connector Guide` to `MfaEnrollDialog`, `CONTRIBUTING.md`, `Step by step`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _596 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.023374936873241468 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.03569182389937107 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022307424189613106 - nodes in this community are weakly interconnected._