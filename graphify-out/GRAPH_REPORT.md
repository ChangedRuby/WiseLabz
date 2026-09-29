# Graph Report - test-remove-real-cron-tick-polling-from-schedule  (2026-09-28)

## Corpus Check
- 877 files · ~546,100 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6544 nodes · 20770 edges · 242 communities (225 shown, 17 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1666 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `91d073af`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- testing.T
- context.Context
- newDocTestStore
- react
- net/http.Request
- @tanstack/react-query
- go_pkg_context
- go_pkg_net_http
- NewChecker
- DashboardPage.tsx
- go_pkg_strings
- ServiceDetailPage.tsx
- ServiceSnapshot
- RulesPage.tsx
- newTestHandler
- net/http.Client
- App.tsx
- WebSocketProvider.tsx
- ProfilePage.tsx
- ConnectorRecord
- UsersPage.tsx
- gitFixture
- cn
- UserIDFromContext
- package.json
- go_pkg_os
- go_pkg_testing
- dispatcher_test.go
- Store
- router.go
- ConnectorEditPage.tsx
- fixtures.ts
- Compare
- routerDeps
- Errorf
- docker_test.go
- store/backup_test.go
- icons.tsx
- RunMigrations
- Connector
- truenas/tables.go
- ws.ts
- Manager
- dependencies
- DecodeKey
- home_assistant/tables.go
- NewUser
- Store
- HashPassword
- newTestHandler
- SystemPage.tsx
- Hub
- NewMalformedResponseError
- portainer/tables.go
- User
- connector/connector.go
- Dispatcher
- adguardhome/tables.go
- Handler
- traefik_test.go
- NewEngine
- response.go
- truenas_test.go
- traefik/tables.go
- unifi/tables.go
- GetTypeSchema
- Connector
- NewEngine
- config_test.go
- settings.mock.ts
- connector_permission.go
- New
- log/slog.Logger
- rewritePlaceholders
- handlers.ts
- Service
- HashToken
- channels.go
- Config
- backup/backup.go
- .Fetch
- Connector
- unifi_test.go
- NewStore
- WiseLabz — Design Contract
- devDependencies
- SuggestRequest
- portainer_test.go
- runbooks_test.go
- Register
- VerifyBundleFile
- NewHTTPClient
- adguardhome_test.go
- time.Time
- lifecycleManager
- home_assistant_test.go
- Runner
- src/theme.ts
- AuthMiddleware
- Connector
- ReportsPage.tsx
- middleware_test.go
- MarshalConnectorConfig
- SnapshotEntity
- diagnostics/diagnostics.go
- keyset_test.go
- ws/ws_test.go
- compilerOptions
- Store
- Handler
- Handler
- sshStdioConn
- New
- RunbookRecord
- docdiffmodel.ts
- net/http.Handler
- vectorCache
- all.go
- httpx/retry_test.go
- .batchDelete
- WiseLabz — Architecture & Technical Decisions
- templates.fixtures.ts
- compilerOptions
- testApp
- AuthedUser
- handlers_contract_test.go
- newTestHarness
- NewClient
- NewService
- manifest.go
- Deps
- DocRecord
- time.Duration
- AppearancePage.tsx
- newTestLifecycle
- Logger
- pagination_contract_test.go
- newTestHandler
- diagram.go
- templatefuncs.go
- render_test.go
- Contributing to WiseLabz
- Backend test performance
- NewRegistry
- api/changes_test.go
- connectors_health_test.go
- Handler
- Engine
- Decision
- scripts
- ThemeControls.tsx
- main
- newTestHandler
- Engine
- matchEntities
- Store
- JobHealthRecord
- Contributor Covenant Code of Conduct
- Decision
- main.tsx
- store/theme.ts
- SuggestWithFallback
- Handler
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- provider_test.go
- .call
- api/auth/oidc.go
- connectors_maintenance_test.go
- .call
- Elector
- Changelog
- test-shards.sh
- mockServiceWorker.js
- apikey_scopes_test.go
- ComplianceRuleRecord
- openapi_contract_test.go
- net/http.Response
- ReportData
- retention/retention_test.go
- scheduler/health_test.go
- BackupSchedule
- WiseLabz — Deployment Guide
- useTheme
- release-please-config.json
- settings.ts
- connectors_hardening_test.go
- dashboard/handlers_test.go
- RateLimit
- validate.go
- ComputeWindow
- ShareLink
- Cache
- Step by step
- WiseLabz
- runConfigCommand
- newHandler
- changes/handlers.go
- TestBulkReauth
- snapshotResponse
- .UpdateAuthConfig
- serveSSHDockerConn
- scanMaintenanceWindow
- computeNextRun
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- ReadyState
- createDNSResolverConnector
- engine_maintenance_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Security Policy
- NewHandler
- seedScopeFixture
- golden_snapshot_test.go
- Authentication design
- Development workflow
- compose-smoke.sh
- expireAlertsOnce
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- Saved Views
- testHarness
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

## Communities (242 total, 17 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (180): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+172 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (159): cursorPage, TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL() (+151 more)

### Community 2 - "context.Context"
Cohesion: 0.03
Nodes (35): fakeStatusChecker, sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, Connector, existingIDs(), SnapshotRecord, Store (+27 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.03
Nodes (123): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+115 more)

### Community 4 - "react"
Cohesion: 0.04
Nodes (102): Frontend, motion, @radix-ui/react-popover, react, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve (+94 more)

### Community 5 - "net/http.Request"
Cohesion: 0.05
Nodes (39): oidcElevateFlow, updateUserRequest, webAuthnFlow, Handler, writeUserWriteError(), clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie() (+31 more)

### Community 6 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (59): i18next, msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, web_src_api_model_index_attentionpage, web_src_api_model_index_docnode (+51 more)

### Community 7 - "go_pkg_context"
Cohesion: 0.07
Nodes (20): StatusError, dashboardLayout, ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), IsTimeout(), TestIsTimeout(), contains() (+12 more)

### Community 8 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (34): bulkSnoozeItemResult, bulkSnoozeRequest, versionSections(), TemplateVersionSection, shareLinkContextKey, go_pkg_crypto_rand, go_pkg_encoding_base64, go_pkg_github_com_go_webauthn_webauthn_protocol (+26 more)

### Community 9 - "NewChecker"
Cohesion: 0.06
Nodes (73): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+65 more)

### Community 10 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (84): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+76 more)

### Community 11 - "go_pkg_strings"
Cohesion: 0.04
Nodes (37): contextKey, elevationError, Schema(), schemaFor(), TestSchemaMatchesConfig(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6() (+29 more)

### Community 12 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (79): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush (+71 more)

### Community 13 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (24): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, changePatternID(), Engine, markError(), snapshotIDOrNil(), init() (+16 more)

### Community 14 - "RulesPage.tsx"
Cohesion: 0.04
Nodes (74): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+66 more)

### Community 15 - "newTestHandler"
Cohesion: 0.06
Nodes (71): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+63 more)

### Community 16 - "net/http.Client"
Cohesion: 0.04
Nodes (25): Connector, ollamaEmbedder, NewServiceUnavailableError(), setHeaders(), tryParseEntities(), validateCustomURL(), Connector, Connector (+17 more)

### Community 17 - "App.tsx"
Cohesion: 0.04
Nodes (64): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken(), setMfaEnrollmentRequiredHandler() (+56 more)

### Community 18 - "WebSocketProvider.tsx"
Cohesion: 0.04
Nodes (64): Frontend shell & theme (decided 2026-06), Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport (+56 more)

### Community 19 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (54): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+46 more)

### Community 20 - "ConnectorRecord"
Cohesion: 0.05
Nodes (35): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), ConnectorRecord, Store, scanConnector() (+27 more)

### Community 21 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (51): axios, customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers (+43 more)

### Community 22 - "gitFixture"
Cohesion: 0.07
Nodes (40): fetchAllDocs(), fileName(), Exporter, NewExporter(), RunExportOnce(), slugify(), newTestStore(), readFile() (+32 more)

### Community 23 - "cn"
Cohesion: 0.05
Nodes (47): clsx, tailwind-merge, web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview (+39 more)

### Community 24 - "UserIDFromContext"
Cohesion: 0.07
Nodes (22): Handler, newToken(), sanitize(), Handler, configRequestField(), validateConnectorConfig(), writeConfigRejection(), Handler (+14 more)

### Community 25 - "package.json"
Cohesion: 0.04
Nodes (51): codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+43 more)

### Community 26 - "go_pkg_os"
Cohesion: 0.06
Nodes (37): main(), usage(), IsGeneratedName(), pruneStale(), TestIsGeneratedName(), TestCommitMessage(), keys(), writeExportState() (+29 more)

### Community 27 - "go_pkg_testing"
Cohesion: 0.07
Nodes (14): go_pkg_crypto_rsa, go_pkg_crypto_tls, go_pkg_encoding_json, go_pkg_github_com_go_jose_go_jose_v4, go_pkg_github_com_gorilla_websocket, go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_wiselabz_wiselabz_internal_api_apitest (+6 more)

### Community 28 - "dispatcher_test.go"
Cohesion: 0.15
Nodes (51): seedDelivery(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore(), setChannelAndRoutingConfig(), setChannelConfig() (+43 more)

### Community 29 - "Store"
Cohesion: 0.06
Nodes (39): Provider, Config, Handler, TestEmbeddedSPAWithoutFrontendBuild(), Embedder, EmbedRegistry, Registry, NewHandler() (+31 more)

### Community 30 - "router.go"
Cohesion: 0.07
Nodes (36): buildPrompt(), TestBuildPrompt(), Match, go_pkg_github_com_wiselabz_wiselabz_internal_ai, go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention (+28 more)

### Community 31 - "ConnectorEditPage.tsx"
Cohesion: 0.06
Nodes (35): RFC-3339, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsschema, web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme (+27 more)

### Community 32 - "fixtures.ts"
Cohesion: 0.06
Nodes (44): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+36 more)

### Community 33 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 34 - "routerDeps"
Cohesion: 0.08
Nodes (39): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+31 more)

### Community 35 - "Errorf"
Cohesion: 0.08
Nodes (15): diffToSpec(), isWritableField(), applyConnectorScalarUpdates(), Handler, Handler, Handler, Handler, ConfigPusher (+7 more)

### Community 36 - "docker_test.go"
Cohesion: 0.06
Nodes (44): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+36 more)

### Community 37 - "store/backup_test.go"
Cohesion: 0.07
Nodes (41): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+33 more)

### Community 38 - "icons.tsx"
Cohesion: 0.10
Nodes (38): web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidversionsrevrestore, web_src_api_generated_docs_docs_usegetdocsdocidversions, web_src_api_generated_docs_docs_usegetdocsdocidversionsrev, web_src_api_model_index_docversionmetatrigger, categoryIcon (+30 more)

### Community 39 - "RunMigrations"
Cohesion: 0.10
Nodes (38): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+30 more)

### Community 40 - "Connector"
Cohesion: 0.07
Nodes (13): init(), ConfigField, Connector, buildRouteTable(), TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName() (+5 more)

### Community 41 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 42 - "ws.ts"
Cohesion: 0.07
Nodes (34): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+26 more)

### Community 43 - "Manager"
Cohesion: 0.08
Nodes (15): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, nilToStr(), Store (+7 more)

### Community 44 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 45 - "DecodeKey"
Cohesion: 0.11
Nodes (22): ProviderConfig, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey(), Decrypt(), DeriveKey() (+14 more)

### Community 46 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 47 - "NewUser"
Cohesion: 0.19
Nodes (39): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld(), TestGetRootIsSynthetic() (+31 more)

### Community 48 - "Store"
Cohesion: 0.07
Nodes (11): changeServiceIDs(), placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+3 more)

### Community 49 - "HashPassword"
Cohesion: 0.09
Nodes (25): sanitizeUser(), setRefreshCookie(), Handler, mustHashDummyPassword(), instanceAdminRoleFor(), SecurityHeaders(), TestSecurityHeaders(), HashPassword() (+17 more)

### Community 50 - "newTestHandler"
Cohesion: 0.12
Nodes (37): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+29 more)

### Community 51 - "SystemPage.tsx"
Cohesion: 0.07
Nodes (31): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey, web_src_api_generated_system_system_getsystembackupschedule (+23 more)

### Community 52 - "Hub"
Cohesion: 0.08
Nodes (17): Handler, decodeBulkRequest(), Handler, loggablePath(), loggableQuery(), Err(), Sanitize(), TestErr() (+9 more)

### Community 53 - "NewMalformedResponseError"
Cohesion: 0.12
Nodes (36): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+28 more)

### Community 54 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 55 - "User"
Cohesion: 0.10
Nodes (14): webAuthnUser, Handler, Handler, newOIDCUser(), randomOIDCToken(), validHostPort(), OIDCClaims, OIDCProvider (+6 more)

### Community 56 - "connector/connector.go"
Cohesion: 0.06
Nodes (23): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP() (+15 more)

### Community 57 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 58 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 59 - "Handler"
Cohesion: 0.09
Nodes (17): validateConfigPushRequest(), parseScheduleUpdates(), validateRotationFields(), capitalize(), Handler, stepAuditDetail(), validTargetType(), validVerb() (+9 more)

### Community 60 - "traefik_test.go"
Cohesion: 0.10
Nodes (34): AllowLoopbackForTest(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath() (+26 more)

### Community 61 - "NewEngine"
Cohesion: 0.13
Nodes (27): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+19 more)

### Community 62 - "response.go"
Cohesion: 0.09
Nodes (24): decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+16 more)

### Community 63 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 64 - "traefik/tables.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+21 more)

### Community 65 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 66 - "GetTypeSchema"
Cohesion: 0.12
Nodes (27): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestRegisteredSchema(), AttributeCatalog() (+19 more)

### Community 67 - "Connector"
Cohesion: 0.12
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 68 - "NewEngine"
Cohesion: 0.19
Nodes (28): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+20 more)

### Community 69 - "config_test.go"
Cohesion: 0.10
Nodes (27): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+19 more)

### Community 70 - "settings.mock.ts"
Cohesion: 0.09
Nodes (25): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+17 more)

### Community 71 - "connector_permission.go"
Cohesion: 0.13
Nodes (15): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), Store, apiKeyConnectorFilter(), getConnectorGrant(), ConnectorGrantDiff (+7 more)

### Community 72 - "New"
Cohesion: 0.11
Nodes (27): confirm(), runRestore(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag(), TestRunVerifyDefaultsToLatestInDir(), TestRunVerifyPassAndFail() (+19 more)

### Community 73 - "log/slog.Logger"
Cohesion: 0.11
Nodes (15): newLogger(), formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep() (+7 more)

### Community 74 - "rewritePlaceholders"
Cohesion: 0.11
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 75 - "handlers.ts"
Cohesion: 0.07
Nodes (27): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+19 more)

### Community 76 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 77 - "HashToken"
Cohesion: 0.15
Nodes (15): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+7 more)

### Community 78 - "channels.go"
Cohesion: 0.10
Nodes (21): cosineSimilarity(), packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips(), TestSplitSections(), unpackVector(), buildEmailMessage() (+13 more)

### Community 79 - "Config"
Cohesion: 0.11
Nodes (20): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+12 more)

### Community 80 - "backup/backup.go"
Cohesion: 0.17
Nodes (26): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+18 more)

### Community 81 - ".Fetch"
Cohesion: 0.11
Nodes (13): ServiceDependency, WantsField(), Connector, environmentDependencies(), putMetadata(), unavailable(), agentEnabled(), Connector (+5 more)

### Community 82 - "Connector"
Cohesion: 0.20
Nodes (5): SnapshotSection, Connector, unavailable(), unavailable(), session

### Community 83 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 84 - "NewStore"
Cohesion: 0.19
Nodes (22): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+14 more)

### Community 85 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 86 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 87 - "SuggestRequest"
Cohesion: 0.13
Nodes (8): claudeProvider, openAICompatibleProvider, openAIEmbedder, StubProvider, SuggestChunk, SuggestRequest, LimitedBody(), countingProvider

### Community 88 - "portainer_test.go"
Cohesion: 0.19
Nodes (22): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestAPIKeyIsStoredAsPassword(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout() (+14 more)

### Community 89 - "runbooks_test.go"
Cohesion: 0.18
Nodes (21): runbookResp, runbookStepResp, createRunbookWithStep(), testApp, seedProxmoxConnector(), TestRunbookExecuteStepDryRunWithoutTokenWorks(), TestRunbookExecuteStepFailureCreatesAlert(), TestRunbookExecuteStepForbiddenWithoutOperatorGrant() (+13 more)

### Community 90 - "Register"
Cohesion: 0.15
Nodes (21): init(), init(), init(), init(), init(), init(), init(), init() (+13 more)

### Community 91 - "VerifyBundleFile"
Cohesion: 0.18
Nodes (20): formatCounts(), runVerify(), TestRunVerifyRequiresABundle(), failVerification(), LatestBundle(), newScratchStore(), RunVerifyOnce(), ListVerifications() (+12 more)

### Community 92 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 93 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 94 - "time.Time"
Cohesion: 0.18
Nodes (19): digestDue(), TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift, DefinitionSummary, DocChangeEntry (+11 more)

### Community 95 - "lifecycleManager"
Cohesion: 0.12
Nodes (9): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, gatedStopScheduler, lifecycleDeps, lifecycleManager (+1 more)

### Community 96 - "home_assistant_test.go"
Cohesion: 0.22
Nodes (19): Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities(), TestFetchConfigIsRequestedOnce() (+11 more)

### Community 97 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 98 - "src/theme.ts"
Cohesion: 0.10
Nodes (15): @fontsource/ibm-plex-mono, @fontsource/ibm-plex-sans, @fontsource/space-mono, @fontsource-variable/big-shoulders-text, @fontsource-variable/geist, @fontsource-variable/geist-mono, @fontsource-variable/inter-tight, @fontsource-variable/jetbrains-mono (+7 more)

### Community 99 - "AuthMiddleware"
Cohesion: 0.13
Nodes (15): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), treatAsSafeFromContext(), AuthMiddleware() (+7 more)

### Community 100 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 101 - "ReportsPage.tsx"
Cohesion: 0.12
Nodes (18): web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports, web_src_api_generated_reports_reports_usegetreportsdefinitions (+10 more)

### Community 102 - "middleware_test.go"
Cohesion: 0.16
Nodes (15): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser() (+7 more)

### Community 103 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestExportRedactsConnectorSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 104 - "SnapshotEntity"
Cohesion: 0.13
Nodes (18): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildPeerTableAttributes() (+10 more)

### Community 105 - "diagnostics/diagnostics.go"
Cohesion: 0.24
Nodes (16): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+8 more)

### Community 106 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 107 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 108 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 109 - "Store"
Cohesion: 0.18
Nodes (7): testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 110 - "Handler"
Cohesion: 0.25
Nodes (6): response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 111 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 112 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 113 - "New"
Cohesion: 0.28
Nodes (16): TestDocExportDefaultCronExprIsValid(), TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart() (+8 more)

### Community 114 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 115 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 116 - "net/http.Handler"
Cohesion: 0.16
Nodes (12): ConnectorRoleChecker, PermissionChecker, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), TreatAsSafeMethod(), RequireConnectorRole() (+4 more)

### Community 117 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 118 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 119 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 121 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 122 - "templates.fixtures.ts"
Cohesion: 0.17
Nodes (13): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+5 more)

### Community 123 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 124 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 125 - "AuthedUser"
Cohesion: 0.21
Nodes (15): TestCreate(), AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter() (+7 more)

### Community 126 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 127 - "newTestHarness"
Cohesion: 0.23
Nodes (15): ContextWithAPIKeyRestriction(), TestClampConnectorRole(), seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding() (+7 more)

### Community 128 - "NewClient"
Cohesion: 0.18
Nodes (14): isSafeMethod(), clientTimeout(), IsSafeMethod(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects() (+6 more)

### Community 129 - "NewService"
Cohesion: 0.25
Nodes (14): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+6 more)

### Community 130 - "manifest.go"
Cohesion: 0.25
Nodes (14): ImportFromFile(), AppVersion(), BuildManifest(), BundleCounts(), ChecksumBytes(), ManifestPath(), ReadManifest(), corruptFile() (+6 more)

### Community 131 - "Deps"
Cohesion: 0.34
Nodes (15): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+7 more)

### Community 132 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 133 - "time.Duration"
Cohesion: 0.16
Nodes (6): Database, HASettings, Server, SyncSettings, time.Duration, PoolConfig

### Community 134 - "AppearancePage.tsx"
Cohesion: 0.22
Nodes (13): AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS, Density (+5 more)

### Community 135 - "newTestLifecycle"
Cohesion: 0.26
Nodes (9): newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestLifecycleManagerWaitsForSchedulerBeforeCancelAndDBClose(), TestStandbyIsUnreadyAndRunsNoScheduler(), waitForLifecycleSignal() (+1 more)

### Community 136 - "Logger"
Cohesion: 0.19
Nodes (14): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+6 more)

### Community 137 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 138 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 139 - "diagram.go"
Cohesion: 0.26
Nodes (12): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), relatedEntities(), EntityLink (+4 more)

### Community 140 - "templatefuncs.go"
Cohesion: 0.19
Nodes (11): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+3 more)

### Community 141 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 142 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 143 - "Backend test performance"
Cohesion: 0.14
Nodes (14): Backend test performance, CI job times, CI measurements, Coverage strategy, Deterministic scheduled exports (#408), Fixture reuse and lifecycle tests (#405–#407), Follow-ups, Local measurements (+6 more)

### Community 144 - "NewRegistry"
Cohesion: 0.45
Nodes (12): NewRegistry(), TestExplain(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip() (+4 more)

### Community 145 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 146 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 147 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 148 - "Engine"
Cohesion: 0.17
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 149 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 150 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 151 - "ThemeControls.tsx"
Cohesion: 0.18
Nodes (10): AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls(), FONT_SETS, makePalette() (+2 more)

### Community 152 - "main"
Cohesion: 0.21
Nodes (12): main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), RegisterOpenAICompatible() (+4 more)

### Community 153 - "newTestHandler"
Cohesion: 0.23
Nodes (12): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+4 more)

### Community 154 - "Engine"
Cohesion: 0.27
Nodes (5): Engine, TemplateFuncs(), GenerateResult, renderResult, text/template.FuncMap

### Community 155 - "matchEntities"
Cohesion: 0.33
Nodes (10): dedupKey(), matchEntities(), matchReason(), seedEngineConnectorWithEntities(), TestGenerateLabTopologyCreatesThenUpdatesInPlace(), TestMatchEntitiesDedupesExactExternalIDDuplicates(), TestMatchEntitiesExternalIDPrecedence(), TestMatchEntitiesFansOutAcrossConnectors() (+2 more)

### Community 156 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 157 - "JobHealthRecord"
Cohesion: 0.29
Nodes (4): JobHealthRecord, Store, scanJobHealth(), fakeHealthStore

### Community 158 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 159 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 160 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 161 - "store/theme.ts"
Cohesion: 0.32
Nodes (11): ColorMode, commit(), Persisted, PRESETS_FONTS, ThemeState, tokensFor(), applyTokens(), FontSetName (+3 more)

### Community 162 - "SuggestWithFallback"
Cohesion: 0.31
Nodes (10): SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders(), TestSuggestWithFallbackStopsOnNonRetryableError() (+2 more)

### Community 164 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 165 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 166 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 167 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 168 - "provider_test.go"
Cohesion: 0.29
Nodes (6): testProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream()

### Community 169 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 170 - "api/auth/oidc.go"
Cohesion: 0.27
Nodes (8): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), go_pkg_crypto_subtle

### Community 171 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 172 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 173 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 174 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 175 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 176 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 177 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 178 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 179 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 180 - "net/http.Response"
Cohesion: 0.33
Nodes (6): retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 181 - "ReportData"
Cohesion: 0.56
Nodes (3): connectorFilter(), Generator, ReportData

### Community 182 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 183 - "scheduler/health_test.go"
Cohesion: 0.44
Nodes (6): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), fakeNotifier, notifyCall

### Community 184 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 185 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 186 - "useTheme"
Cohesion: 0.33
Nodes (7): mermaid, cssVar(), Mermaid(), resolveColor(), load(), useTheme, presetOpts()

### Community 187 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 188 - "settings.ts"
Cohesion: 0.39
Nodes (7): MotionProvider(), apply(), framerReducedMotion(), MotionPref, seed(), SettingsState, useSettings

### Community 189 - "connectors_hardening_test.go"
Cohesion: 0.29
Nodes (7): testApp, TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 190 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 191 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 192 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 193 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 195 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 196 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 197 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 198 - "runConfigCommand"
Cohesion: 0.33
Nodes (7): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), io.Writer

### Community 199 - "newHandler"
Cohesion: 0.29
Nodes (7): NewEmbedRegistry(), Handler, newHandler(), TestConversationOwnership(), TestCreateConversationDocVisibility(), TestCreateConversationValidation(), TestPostMessageErrors()

### Community 200 - "changes/handlers.go"
Cohesion: 0.38
Nodes (6): changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest, go_pkg_unicode_utf8

### Community 201 - "TestBulkReauth"
Cohesion: 0.48
Nodes (7): bulkReq(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync()

### Community 202 - "snapshotResponse"
Cohesion: 0.48
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 203 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 204 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 205 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 206 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 207 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 208 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 209 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 210 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 212 - "createDNSResolverConnector"
Cohesion: 0.33
Nodes (6): createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler()

### Community 213 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 214 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 215 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 216 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 217 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 218 - "NewHandler"
Cohesion: 0.60
Nodes (4): NewHandler(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound()

### Community 219 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 220 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 222 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 223 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 224 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 225 - "expireAlertsOnce"
Cohesion: 0.67
Nodes (4): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce()

### Community 228 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 229 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 230 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 231 - "testHarness"
Cohesion: 0.50
Nodes (3): github.com/mark3labs/mcp-go/client.Client, github.com/mark3labs/mcp-go/mcp.CallToolResult, testHarness

## Knowledge Gaps
- **586 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+581 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1313 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Hub` connect `Hub` to `Errorf`, `NewEngine`, `net/http.Request`, `time.Duration`, `go_pkg_net_http`, `NewChecker`, `log/slog.Logger`, `ws/ws_test.go`, `dispatcher_test.go`, `NewEngine`, `Engine`, `Dispatcher`, `testApp`, `Store`, `lifecycleManager`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Election` connect `lifecycleManager` to `go_pkg_context`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `Errorf` to `routerDeps`, `Handler`, `net/http.Request`, `Logger`, `.UpdateAuthConfig`, `HashToken`, `Handler`, `DecodeKey`, `Handler`, `HashPassword`, `Handler`, `Hub`, `User`, `UserIDFromContext`, `Handler`, `response.go`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _586 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.020401554404145077 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019746221842030225 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.025127723907830257 - nodes in this community are weakly interconnected._