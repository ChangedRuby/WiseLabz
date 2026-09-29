# Graph Report - feat-ha-evaluate-cross-replica-websocket-event-d  (2026-09-29)

## Corpus Check
- 878 files · ~548,533 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6559 nodes · 20786 edges · 243 communities (217 shown, 26 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1666 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1d90ebf6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- react
- testing.T
- newDocTestStore
- go_pkg_testing
- net/http.Request
- cn
- context.Context
- DashboardPage.tsx
- package.json
- go_pkg_context
- go_pkg_github_com_wiselabz_wiselabz_internal_store
- @tanstack/react-query
- ServiceDetailPage.tsx
- newTestHandler
- App.tsx
- vitest
- go_pkg_os
- icons.tsx
- net/http.Client
- ServiceSnapshot
- NewMalformedResponseError
- ProfilePage.tsx
- git_test.go
- UsersPage.tsx
- net/http.ResponseWriter
- New
- NewEngine
- Store
- rowScanner
- Store
- Compare
- GetTypeSchema
- fixtures.ts
- response.go
- dispatcher_test.go
- SystemPage.tsx
- SnapshotEntity
- routerDeps
- DecodeKey
- home_assistant/tables.go
- dependencies
- docker_test.go
- store/backup_test.go
- RunMigrations
- NewChecker
- home_assistant_test.go
- router.go
- portainer/tables.go
- ErrorWithDetails
- adguardhome/tables.go
- compliance/engine.go
- Hub
- connector/connector.go
- ExportToFile
- NewEngine
- traefik/tables.go
- Dispatcher
- Connector
- unifi/tables.go
- Handler
- channels.go
- traefik_test.go
- WebSocketProvider.tsx
- handlers.ts
- AuthMiddleware
- User
- settings.mock.ts
- Config
- Connector
- nilToStr
- DecodeJSON
- NewUser
- config_test.go
- NewStore
- share_links_test.go
- rewritePlaceholders
- runbooks_test.go
- Service
- newTestHandler
- log/slog.Logger
- New
- unifi_test.go
- SuggestRequest
- timeline.ts
- Checker
- Manager
- WiseLabz — Design Contract
- devDependencies
- ConnectorRecord
- portainer_test.go
- adguardhome_test.go
- .Fetch
- time.Time
- lifecycleManager
- connector_permission.go
- NewService
- backup/backup.go
- NewHTTPClient
- httpx/retry_test.go
- Deps
- ReportsPage.tsx
- main
- ws.ts
- MarshalConnectorConfig
- Handler
- truenas_test.go
- diagnostics/diagnostics.go
- keyset_test.go
- compilerOptions
- net/http.Handler
- Handler
- Register
- Connector
- sshStdioConn
- Store
- RunbookRecord
- docdiffmodel.ts
- notifications/handlers_test.go
- system/handlers_test.go
- VerifyBundleFile
- vectorCache
- all.go
- time.Duration
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- testApp
- handlers_contract_test.go
- handlers_actions_test.go
- DocRecord
- SuggestWithFallback
- config_cmd_test.go
- newTestLifecycle
- changes/handlers_test.go
- pagination_contract_test.go
- Connector
- templatefuncs.go
- render_test.go
- Contributing to WiseLabz
- WiseLabz WebSocket Contract (`/ws`)
- templates.fixtures.ts
- NewRegistry
- api/changes_test.go
- connectors_health_test.go
- Handler
- newTestHandler
- internal/auth/mfa.go
- scripts
- IsSecureRequest
- NewClient
- TestListAttentionItems
- Store
- Store
- Contributor Covenant Code of Conduct
- Decision
- Decision
- Backend test performance
- Connector
- Handler
- apikey_scope.go
- ReportData
- 0001 — Lab-mutating operation boundaries
- Decision
- WiseLabz Connector Guide
- Product
- .call
- cursor_pagination_test.go
- api/auth/oidc.go
- handlers_bulk_test.go
- NewHandler
- .call
- Elector
- Changelog
- test-shards.sh
- mockServiceWorker.js
- ComplianceRuleRecord
- connectors_hardening_test.go
- newTestHandler
- retention/retention_test.go
- Engine
- BACKUP.md
- release-please-config.json
- WithAuth
- dashboard/handlers_test.go
- RateLimit
- validate.go
- pfsense.go
- ComputeWindow
- ShareLink
- Cache
- Decision
- Step by step
- WiseLabz
- TestComplianceRuleValidation
- serveSSHDockerConn
- Connector
- scanMaintenanceWindow
- computeNextRun
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- docker/snapshot.go
- .GetConnectorUptime
- Store
- engine_maintenance_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- WiseLabz — Deployment Guide
- Scheduled Doc Export
- Security Policy
- RequireConnectorRole
- APIKeyClaims
- routerOperations
- Connector
- ref_test.go
- Authentication design
- Development workflow
- compose-smoke.sh
- expireAlertsOnce
- ClassifyHealth
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- Saved Views
- Fixture reuse and lifecycle tests (#405–#407)
- fields_test.go
- truenas/attributes_test.go
- stubEmbedder
- dockerSSHAddr
- WiseLabz — v2 Backlog
- fakeEmbedder
- coverage-parity.sh
- fakeQualityChecker
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- coverpkg.sh
- QualityChecker
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

## Communities (243 total, 26 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (175): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+167 more)

### Community 1 - "react"
Cohesion: 0.04
Nodes (121): RFC-3339, motion, react, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze (+113 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (145): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), createKey(), testApp (+137 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.03
Nodes (136): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+128 more)

### Community 4 - "go_pkg_testing"
Cohesion: 0.05
Nodes (34): contextKey, elevationError, TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), normalizeEnabledColumn() (+26 more)

### Community 5 - "net/http.Request"
Cohesion: 0.05
Nodes (42): Handler, PermissionChecker, newToken(), sanitize(), setRefreshCookie(), factorJSON(), Handler, Handler (+34 more)

### Community 6 - "cn"
Cohesion: 0.03
Nodes (92): clsx, tailwind-merge, web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey, web_src_api_generated_settings_settings_getgetauthconfigquerykey, web_src_api_generated_settings_settings_getgetnotificationsconfigquerykey, web_src_api_generated_settings_settings_postaiconfigtest (+84 more)

### Community 7 - "context.Context"
Cohesion: 0.03
Nodes (25): fakeStatusChecker, sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, existingIDs(), Store, MFAFactor, Store (+17 more)

### Community 8 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (99): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete` (+91 more)

### Community 9 - "package.json"
Cohesion: 0.03
Nodes (87): codemirror, @codemirror/commands, @codemirror/state, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, @fontsource/ibm-plex-mono (+79 more)

### Community 10 - "go_pkg_context"
Cohesion: 0.07
Nodes (22): StatusError, IsTimeout(), TestIsTimeout(), contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors (+14 more)

### Community 11 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.05
Nodes (41): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection (+33 more)

### Community 12 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (72): Frontend shell & theme (decided 2026-06), Frontend, match-sorter, @radix-ui/react-popover, react-router-dom, @tanstack/react-query, zustand, web_src_api_generated_connectors_connectors (+64 more)

### Community 13 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (83): ADR-0001, ADR-0003, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest (+75 more)

### Community 14 - "newTestHandler"
Cohesion: 0.06
Nodes (71): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+63 more)

### Community 15 - "App.tsx"
Cohesion: 0.04
Nodes (65): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken(), setMfaEnrollmentRequiredHandler() (+57 more)

### Community 16 - "vitest"
Cohesion: 0.04
Nodes (46): i18next, msw, @testing-library/jest-dom, @testing-library/react, vitest, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage, ConfirmDestructive() (+38 more)

### Community 17 - "go_pkg_os"
Cohesion: 0.05
Nodes (48): writeExportState(), migrationFiles(), TestMigrationVersionParity(), exportCursor, exportState, go_pkg_crypto_ed25519, go_pkg_encoding_csv, go_pkg_encoding_pem (+40 more)

### Community 18 - "icons.tsx"
Cohesion: 0.05
Nodes (63): @codemirror/lang-markdown, @codemirror/view, @uiw/react-codemirror, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey (+55 more)

### Community 19 - "net/http.Client"
Cohesion: 0.04
Nodes (25): Connector, ollamaEmbedder, openAIEmbedder, setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector (+17 more)

### Community 20 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (19): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, init(), RegisterTransformer(), runTransformers(), TestRunTransformersAppliesInOrderAndStopsOnError(), TestRunTransformersUnknownCategoryIsNoop() (+11 more)

### Community 21 - "NewMalformedResponseError"
Cohesion: 0.07
Nodes (47): unavailable(), SnapshotSection, NewMalformedResponseError(), TestBuildHostsTableAttributes(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable() (+39 more)

### Community 22 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (54): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+46 more)

### Community 23 - "git_test.go"
Cohesion: 0.07
Nodes (46): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+38 more)

### Community 24 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (49): axios, customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers (+41 more)

### Community 25 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (28): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+20 more)

### Community 26 - "New"
Cohesion: 0.07
Nodes (33): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+25 more)

### Community 27 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 28 - "Store"
Cohesion: 0.06
Nodes (26): Provider, Config, Handler, Registry, NewHandler(), NewHandler(), NewHandler(), NewHandler() (+18 more)

### Community 29 - "rowScanner"
Cohesion: 0.07
Nodes (25): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), Store, scanBackupRun(), changeFilterClause() (+17 more)

### Community 30 - "Store"
Cohesion: 0.06
Nodes (16): changeServiceIDs(), Store, placeholders(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+8 more)

### Community 31 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 32 - "GetTypeSchema"
Cohesion: 0.06
Nodes (44): TestRegisteredSchema(), TestSchemaConfigValidation(), countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), TestFailureContract(), supportedLifecycleVerbs() (+36 more)

### Community 33 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 34 - "response.go"
Cohesion: 0.07
Nodes (29): decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+21 more)

### Community 35 - "dispatcher_test.go"
Cohesion: 0.20
Nodes (46): NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore(), setChannelAndRoutingConfig(), setChannelConfig(), setChannelConfigJSON() (+38 more)

### Community 36 - "SystemPage.tsx"
Cohesion: 0.06
Nodes (38): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey (+30 more)

### Community 37 - "SnapshotEntity"
Cohesion: 0.13
Nodes (42): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+34 more)

### Community 38 - "routerDeps"
Cohesion: 0.10
Nodes (35): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+27 more)

### Community 39 - "DecodeKey"
Cohesion: 0.10
Nodes (23): ProviderConfig, testHandler, Handler, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey() (+15 more)

### Community 40 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 41 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 42 - "docker_test.go"
Cohesion: 0.07
Nodes (40): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+32 more)

### Community 43 - "store/backup_test.go"
Cohesion: 0.07
Nodes (37): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+29 more)

### Community 44 - "RunMigrations"
Cohesion: 0.11
Nodes (34): main(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns(), TestMigrationSchemaParity(), RunMigrations() (+26 more)

### Community 45 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 46 - "home_assistant_test.go"
Cohesion: 0.09
Nodes (38): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+30 more)

### Community 47 - "router.go"
Cohesion: 0.08
Nodes (33): TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+25 more)

### Community 48 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 49 - "ErrorWithDetails"
Cohesion: 0.09
Nodes (17): applyConnectorScalarUpdates(), configRequestField(), Handler, parseScheduleUpdates(), validateConnectorConfig(), validateRotationFields(), writeConfigRejection(), stepAuditDetail() (+9 more)

### Community 50 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 51 - "compliance/engine.go"
Cohesion: 0.11
Nodes (31): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+23 more)

### Community 52 - "Hub"
Cohesion: 0.09
Nodes (16): Handler, isWritableField(), loggablePath(), loggableQuery(), ConfigPusher, Err(), Sanitize(), TestErr() (+8 more)

### Community 53 - "connector/connector.go"
Cohesion: 0.07
Nodes (22): Capabilities(), CapabilityDescriptor, ServiceDependency, TimeoutError, GuardedDialer(), IsDangerousIP(), NewServiceUnavailableError(), NewTimeoutError() (+14 more)

### Community 54 - "ExportToFile"
Cohesion: 0.14
Nodes (32): Export(), ExportToFile(), Import(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile() (+24 more)

### Community 55 - "NewEngine"
Cohesion: 0.13
Nodes (27): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+19 more)

### Community 56 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 57 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 58 - "Connector"
Cohesion: 0.09
Nodes (10): ConfigField, Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName(), PathSegment() (+2 more)

### Community 59 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 60 - "Handler"
Cohesion: 0.11
Nodes (15): updateUserRequest, Handler, sanitizeUser(), writeUserWriteError(), Handler, mustHashDummyPassword(), HashPassword(), TestBcryptCostOnlyLoweredUnderTest() (+7 more)

### Community 61 - "channels.go"
Cohesion: 0.08
Nodes (26): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), SplitSections(), SyncDocEmbeddings() (+18 more)

### Community 62 - "traefik_test.go"
Cohesion: 0.11
Nodes (30): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+22 more)

### Community 63 - "WebSocketProvider.tsx"
Cohesion: 0.09
Nodes (26): Sync flow, qrcode, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_changes_changes, web_src_api_generated_changes_changes_getgetchangesquerykey, web_src_api_generated_dashboard_dashboard, web_src_api_generated_dashboard_dashboard_getgetdashboardoverviewquerykey, web_src_api_generated_findings_findings_getgetfindingsquerykey (+18 more)

### Community 64 - "handlers.ts"
Cohesion: 0.07
Nodes (29): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+21 more)

### Community 65 - "AuthMiddleware"
Cohesion: 0.09
Nodes (25): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, AuthMiddleware(), extractBearerToken(), hashToken() (+17 more)

### Community 66 - "User"
Cohesion: 0.12
Nodes (12): Handler, Handler, newOIDCUser(), randomOIDCToken(), validHostPort(), instanceAdminRoleFor(), OIDCClaims, OIDCProvider (+4 more)

### Community 67 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 68 - "Config"
Cohesion: 0.10
Nodes (23): newLogger(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+15 more)

### Community 69 - "Connector"
Cohesion: 0.12
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 70 - "nilToStr"
Cohesion: 0.09
Nodes (10): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+2 more)

### Community 71 - "DecodeJSON"
Cohesion: 0.12
Nodes (11): webAuthnFlow, webAuthnUser, Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T (+3 more)

### Community 72 - "NewUser"
Cohesion: 0.15
Nodes (29): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationDocVisibility() (+21 more)

### Community 73 - "config_test.go"
Cohesion: 0.10
Nodes (27): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+19 more)

### Community 74 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 75 - "share_links_test.go"
Cohesion: 0.20
Nodes (27): Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestListEmpty(), TestSave(), TestVersionsOfUnknownDoc() (+19 more)

### Community 76 - "rewritePlaceholders"
Cohesion: 0.11
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 77 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 78 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 79 - "newTestHandler"
Cohesion: 0.12
Nodes (26): TestConnectorStoreErrorPaths(), TestDataSnapshotPaths(), TestSyncsLimitAndIsolation(), createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler() (+18 more)

### Community 80 - "log/slog.Logger"
Cohesion: 0.11
Nodes (14): formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+6 more)

### Community 81 - "New"
Cohesion: 0.10
Nodes (26): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+18 more)

### Community 82 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 83 - "SuggestRequest"
Cohesion: 0.13
Nodes (11): claudeProvider, openAICompatibleProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList(), TestStubProviderName() (+3 more)

### Community 84 - "timeline.ts"
Cohesion: 0.15
Nodes (16): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+8 more)

### Community 85 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 86 - "Manager"
Cohesion: 0.16
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 87 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 88 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 89 - "ConnectorRecord"
Cohesion: 0.15
Nodes (13): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), scanSyncRun(), connectorWithRole (+5 more)

### Community 90 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 91 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 92 - ".Fetch"
Cohesion: 0.15
Nodes (8): WantsField(), Connector, putMetadata(), unavailable(), agentEnabled(), Connector, Connector, dockerSectionSpec

### Community 93 - "time.Time"
Cohesion: 0.18
Nodes (19): digestDue(), TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift, DefinitionSummary, DocChangeEntry (+11 more)

### Community 94 - "lifecycleManager"
Cohesion: 0.12
Nodes (9): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, gatedStopScheduler, lifecycleDeps, lifecycleManager (+1 more)

### Community 95 - "connector_permission.go"
Cohesion: 0.19
Nodes (9): auditConnectorGrantDiffJSON(), getConnectorGrant(), ConnectorGrantDiff, Store, highestConnectorRole(), listOIDCConnectorGrants(), scanConnectorGrants(), upsertConnectorGrant() (+1 more)

### Community 96 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 97 - "backup/backup.go"
Cohesion: 0.24
Nodes (19): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, importBundle(), importConnectors() (+11 more)

### Community 98 - "NewHTTPClient"
Cohesion: 0.13
Nodes (15): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+7 more)

### Community 99 - "httpx/retry_test.go"
Cohesion: 0.26
Nodes (17): retryable(), RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries() (+9 more)

### Community 100 - "Deps"
Cohesion: 0.22
Nodes (19): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs(), registerListFindings() (+11 more)

### Community 101 - "ReportsPage.tsx"
Cohesion: 0.11
Nodes (19): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports (+11 more)

### Community 102 - "main"
Cohesion: 0.16
Nodes (14): main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder (+6 more)

### Community 103 - "ws.ts"
Cohesion: 0.11
Nodes (18): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+10 more)

### Community 104 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 105 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 106 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 107 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 108 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 109 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 110 - "net/http.Handler"
Cohesion: 0.16
Nodes (16): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+8 more)

### Community 111 - "Handler"
Cohesion: 0.25
Nodes (6): response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 112 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 113 - "Connector"
Cohesion: 0.18
Nodes (6): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 114 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 115 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 116 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 117 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 118 - "notifications/handlers_test.go"
Cohesion: 0.28
Nodes (15): AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+7 more)

### Community 119 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 120 - "VerifyBundleFile"
Cohesion: 0.23
Nodes (15): failVerification(), LatestBundle(), newScratchStore(), RunVerifyOnce(), ListVerifications(), RecordVerification(), seedOneDoc(), TestLatestBundleNoBundles() (+7 more)

### Community 121 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 122 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 123 - "time.Duration"
Cohesion: 0.17
Nodes (7): sleep(), Database, Server, time.Duration, RetryPolicy, retryTransport, PoolConfig

### Community 124 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 125 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 126 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 127 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 128 - "handlers_actions_test.go"
Cohesion: 0.35
Nodes (14): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+6 more)

### Community 129 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 130 - "SuggestWithFallback"
Cohesion: 0.23
Nodes (11): StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+3 more)

### Community 131 - "config_cmd_test.go"
Cohesion: 0.21
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 132 - "newTestLifecycle"
Cohesion: 0.26
Nodes (9): newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestLifecycleManagerWaitsForSchedulerBeforeCancelAndDBClose(), TestStandbyIsUnreadyAndRunsNoScheduler(), waitForLifecycleSignal() (+1 more)

### Community 133 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 134 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 135 - "Connector"
Cohesion: 0.19
Nodes (6): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), Connector

### Community 136 - "templatefuncs.go"
Cohesion: 0.19
Nodes (11): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+3 more)

### Community 137 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 138 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 139 - "WiseLabz WebSocket Contract (`/ws`)"
Cohesion: 0.14
Nodes (11): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention (+3 more)

### Community 140 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (11): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, resolveToken(), Snapshot (+3 more)

### Community 141 - "NewRegistry"
Cohesion: 0.45
Nodes (12): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+4 more)

### Community 142 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 143 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 144 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 145 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 146 - "internal/auth/mfa.go"
Cohesion: 0.24
Nodes (11): GenerateRecoveryCodes(), GenerateTOTPSecret(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL(), TestNormalizeRecoveryCode(), TestValidateTOTPAcceptsCurrentAndSkewedCodes(), TestValidateTOTPRejectsEmptyCode() (+3 more)

### Community 147 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 148 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 149 - "NewClient"
Cohesion: 0.24
Nodes (11): clientTimeout(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects(), TestNewClientTimeout(), TestNewClientUsesCustomDialContext() (+3 more)

### Community 150 - "TestListAttentionItems"
Cohesion: 0.26
Nodes (12): seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding(), TestSearchDocs(), seedFinding() (+4 more)

### Community 151 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 152 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 153 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 154 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 155 - "Decision"
Cohesion: 0.17
Nodes (12): 0005 — Cross-replica WebSocket event relay for active/active, Authorization and secrets, Consequences, Context, Decision, Duplicate suppression and ordering, Event ownership: local first, then relay, Mechanism: PostgreSQL LISTEN/NOTIFY (+4 more)

### Community 156 - "Backend test performance"
Cohesion: 0.17
Nodes (12): Backend test performance, CI job times, Coverage strategy, Deterministic scheduled exports (#408), Follow-ups, Measurements and validation, Measuring, Rules for new tests (+4 more)

### Community 159 - "apikey_scope.go"
Cohesion: 0.27
Nodes (9): APIKeyRestriction, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), isSafeMethod(), TestClampConnectorRole(), treatAsSafeFromContext(), IsSafeMethod() (+1 more)

### Community 160 - "ReportData"
Cohesion: 0.40
Nodes (5): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), Generator, ReportData

### Community 161 - "0001 — Lab-mutating operation boundaries"
Cohesion: 0.18
Nodes (8): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Consequences, Context, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 162 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 163 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 164 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 165 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 166 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 167 - "api/auth/oidc.go"
Cohesion: 0.27
Nodes (8): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), go_pkg_crypto_subtle

### Community 168 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 169 - "NewHandler"
Cohesion: 0.24
Nodes (10): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+2 more)

### Community 170 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 171 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 172 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 173 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 174 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 175 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 176 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 177 - "newTestHandler"
Cohesion: 0.25
Nodes (9): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound(), TestListMutuallyExclusiveFilters() (+1 more)

### Community 178 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 179 - "Engine"
Cohesion: 0.25
Nodes (3): Engine, sync.Map, DocRegenerator

### Community 180 - "BACKUP.md"
Cohesion: 0.25
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 181 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 182 - "WithAuth"
Cohesion: 0.43
Nodes (8): NewHandler(), TestCreate(), TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), TestDelete()

### Community 183 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 184 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 185 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 186 - "pfsense.go"
Cohesion: 0.39
Nodes (6): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName()

### Community 187 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 189 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 190 - "Decision"
Cohesion: 0.25
Nodes (8): Audit, Authorization, Confirmation / step-up, Decision, Dry-run, Eligible first operation: `service.restart` only, Out of scope, Rollback expectations

### Community 191 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 192 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 193 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 194 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 196 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 197 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 198 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 199 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 200 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 201 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 203 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 205 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 206 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 207 - "WiseLabz — Deployment Guide"
Cohesion: 0.33
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 208 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 209 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 210 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 211 - "APIKeyClaims"
Cohesion: 0.50
Nodes (3): testAPIKeyChecker, APIKeyClaims, validAPIKey()

### Community 212 - "routerOperations"
Cohesion: 0.50
Nodes (5): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes

### Community 214 - "ref_test.go"
Cohesion: 0.40
Nodes (4): TestValidateCompositeRef(), TestValidateRefSegment(), TestValidateUnixSocketPath(), ValidateUnixSocketPath()

### Community 216 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 217 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 218 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 219 - "expireAlertsOnce"
Cohesion: 0.67
Nodes (4): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce()

### Community 220 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 223 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 224 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 225 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 226 - "Fixture reuse and lifecycle tests (#405–#407)"
Cohesion: 0.50
Nodes (4): CI measurements, Fixture reuse and lifecycle tests (#405–#407), Local measurements, Validation and coverage

## Knowledge Gaps
- **597 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+592 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1323 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Hub` connect `Hub` to `testing.T`, `dispatcher_test.go`, `net/http.Request`, `main`, `NewHandler`, `go_pkg_github_com_wiselabz_wiselabz_internal_store`, `NewChecker`, `log/slog.Logger`, `ErrorWithDetails`, `Engine`, `lifecycleManager`, `Checker`, `NewEngine`, `Dispatcher`, `time.Duration`, `Store`, `testApp`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Election` connect `lifecycleManager` to `go_pkg_context`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `net/http.Request` to `response.go`, `User`, `DecodeJSON`, `DecodeKey`, `Handler`, `Store`, `net/http.Handler`, `Handler`, `Handler`, `ErrorWithDetails`, `Hub`, `net/http.ResponseWriter`, `Handler`, `Handler`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _597 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.02127541074909496 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.03555591389180037 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022287390029325515 - nodes in this community are weakly interconnected._