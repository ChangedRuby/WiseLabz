# Graph Report - ci-retain-an-on-demand-shuffled-race-stress-gate  (2026-09-28)

## Corpus Check
- 877 files · ~547,196 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6546 nodes · 20772 edges · 226 communities (203 shown, 23 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1666 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `df748df0`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- testing.T
- cn
- newDocTestStore
- @tanstack/react-query
- react
- context.Context
- net/http.ResponseWriter
- go_pkg_net_http
- DashboardPage.tsx
- go_pkg_context
- icons.tsx
- go_pkg_testing
- newTestHandler
- NewChecker
- go_pkg_encoding_json
- net/http.Client
- ServiceSnapshot
- App.tsx
- Store
- time.Time
- ProfilePage.tsx
- ConnectorRecord
- ServiceDetailPage.tsx
- gitFixture
- package.json
- ErrorWithDetails
- New
- DecodeKey
- NewEngine
- ThemeControls.tsx
- net/http.Request
- SnapshotEntity
- ConnectorEditPage.tsx
- dispatcher_test.go
- fixtures.ts
- routerDeps
- Compare
- ExportToFile
- store/backup_test.go
- RunMigrations
- traefik_test.go
- NewMalformedResponseError
- home_assistant/tables.go
- dependencies
- time.Duration
- SystemPage.tsx
- docker_test.go
- WebSocketProvider.tsx
- home_assistant_test.go
- portainer/tables.go
- adguardhome/tables.go
- response.go
- Dispatcher
- Manager
- NewUser
- NewEngine
- .Fetch
- Store
- HashPassword
- connector/connector.go
- unifi/tables_test.go
- NewService
- settings.mock.ts
- Connector
- rewritePlaceholders
- DecodeJSON
- New
- NewStore
- nilToStr
- SuggestRequest
- connector_permission.go
- log/slog.Logger
- router.go
- runbooks_test.go
- Service
- main
- Register
- handlers.ts
- newTestHandler
- Connector
- unifi_test.go
- timeline.ts
- Connector
- NewRegistry
- notifications/handlers_test.go
- compliance/handlers.go
- Connector
- WiseLabz — Design Contract
- devDependencies
- AppearancePage.tsx
- backup/backup.go
- portainer_test.go
- git.go
- Store
- NewHTTPClient
- adguardhome_test.go
- Deps
- middleware_test.go
- lifecycleManager
- .Fetch
- connectors_health_test.go
- middleware.go
- lifecycle_test.go
- handlers_actions_test.go
- httpx/retry_test.go
- ws/ws_test.go
- ReportsPage.tsx
- DiffViewer.tsx
- truenas_test.go
- diagnostics/diagnostics.go
- ws.ts
- compilerOptions
- Handler
- docs/handlers_test.go
- api/mcp_test.go
- Handler
- createUser
- MarshalConnectorConfig
- RunbookRecord
- docdiffmodel.ts
- templates.fixtures.ts
- net/http.Handler
- all.go
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- handlers_contract_test.go
- Connector
- fetch_test.go
- DocRecord
- changes/handlers_test.go
- IsSecureRequest
- pagination_contract_test.go
- reports/handlers.go
- Contributing to WiseLabz
- api/changes_test.go
- compliance_rules_test.go
- newTestHandler
- dialSSHStdio
- Decision
- scripts
- Handler
- config/validate_test.go
- templatefuncs.go
- Engine
- Contributor Covenant Code of Conduct
- Decision
- Backend test performance
- main.tsx
- Handler
- findings/handlers_authz_test.go
- runbooks/handlers_test.go
- chat/chat.go
- vectorCache
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- .call
- cursor_pagination_test.go
- httpx/retry.go
- Handler
- Changelog
- test-shards.sh
- mockServiceWorker.js
- apikey_scopes_test.go
- connectors_hardening_test.go
- ListSchemas
- retention/retention_test.go
- BackupSchedule
- WiseLabz — Deployment Guide
- Elector
- release-please-config.json
- auth/handlers_test.go
- chat/handlers_test.go
- snapshots_test.go
- dashboard/handlers_test.go
- ShareLink
- Cache
- Step by step
- WiseLabz
- runConfigCommand
- TestBulkReauth
- openapi_contract_test.go
- internal/auth/mfa.go
- pfsense.go
- scanMaintenanceWindow
- computeNextRun
- registryTestRefresher
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- .GetConnectorUptime
- Store
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Security Policy
- golden_snapshot_test.go
- Authentication design
- Development workflow
- compose-smoke.sh
- .Redacted
- ClassifyHealth
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- Saved Views
- Fixture reuse and lifecycle tests (#405–#407)
- internal/auth/oidc.go
- stubEmbedder
- WiseLabz — v2 Backlog
- fakeEmbedder
- coverage-parity.sh
- fakeDocRegenerator
- fakeQualityChecker
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

## Communities (226 total, 23 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (176): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+168 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (158): TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider(), TestGetOrInitOIDCProvider(), TestOIDCCallbackRejectsMissingFlowCookie(), TestOIDCCallbackRejectsStateMismatchedWithCookie(), TestOIDCCallbackRejectsUnknownProvider() (+150 more)

### Community 2 - "cn"
Cohesion: 0.03
Nodes (126): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze, web_src_api_generated_alerts_alerts_postalertsbulksnooze (+118 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.03
Nodes (139): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+131 more)

### Community 4 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (80): axios, msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, customInstance(), web_src_api_model_index_attentionpage (+72 more)

### Community 5 - "react"
Cohesion: 0.03
Nodes (113): react, react-i18next, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid (+105 more)

### Community 6 - "context.Context"
Cohesion: 0.03
Nodes (27): sanitizeSessions(), Connector, existingIDs(), Store, scanComplianceRule(), SnapshotRecord, Store, Store (+19 more)

### Community 7 - "net/http.ResponseWriter"
Cohesion: 0.05
Nodes (40): Handler, sanitize(), Handler, isWritableField(), Handler, applyConnectorScalarUpdates(), Handler, validateConnectorConfig() (+32 more)

### Community 8 - "go_pkg_net_http"
Cohesion: 0.05
Nodes (53): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, updateUserRequest, newToken(), validHostPort(), changePromptData(), diffToSpec() (+45 more)

### Community 9 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (95): Frontend shell & theme (decided 2026-06), 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete` (+87 more)

### Community 10 - "go_pkg_context"
Cohesion: 0.07
Nodes (26): StatusError, IsGeneratedName(), pruneStale(), TestIsGeneratedName(), contains(), searchString(), go_pkg_bufio, go_pkg_context (+18 more)

### Community 11 - "icons.tsx"
Cohesion: 0.04
Nodes (80): i18next, match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow (+72 more)

### Community 12 - "go_pkg_testing"
Cohesion: 0.04
Nodes (39): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders() (+31 more)

### Community 13 - "newTestHandler"
Cohesion: 0.06
Nodes (76): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+68 more)

### Community 14 - "NewChecker"
Cohesion: 0.06
Nodes (70): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+62 more)

### Community 15 - "go_pkg_encoding_json"
Cohesion: 0.04
Nodes (39): healthSummary(), hostFromRule(), TestHealthSummary(), TestHostFromRule(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEmailMessage(), sendSMTPChannel() (+31 more)

### Community 16 - "net/http.Client"
Cohesion: 0.04
Nodes (28): Connector, ollamaEmbedder, openAIEmbedder, setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector (+20 more)

### Community 17 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (25): noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, changePatternID(), Engine, markError() (+17 more)

### Community 18 - "App.tsx"
Cohesion: 0.04
Nodes (58): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken(), setMfaEnrollmentRequiredHandler() (+50 more)

### Community 19 - "Store"
Cohesion: 0.04
Nodes (39): Provider, Config, Handler, Embedder, EmbedRegistry, Registry, NewHandler(), NewHandler() (+31 more)

### Community 20 - "time.Time"
Cohesion: 0.05
Nodes (54): TestRateLimit(), RateLimit(), digestDue(), TestDigestDue(), connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), RenderHTML() (+46 more)

### Community 21 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (55): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+47 more)

### Community 22 - "ConnectorRecord"
Cohesion: 0.06
Nodes (34): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), ConnectorRecord, Store, scanConnector() (+26 more)

### Community 23 - "ServiceDetailPage.tsx"
Cohesion: 0.05
Nodes (57): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+49 more)

### Community 24 - "gitFixture"
Cohesion: 0.07
Nodes (41): fetchAllDocs(), fileName(), Exporter, NewExporter(), RunExportOnce(), slugify(), newTestStore(), readFile() (+33 more)

### Community 25 - "package.json"
Cohesion: 0.04
Nodes (54): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+46 more)

### Community 26 - "ErrorWithDetails"
Cohesion: 0.07
Nodes (22): Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, Handler, Handler, newOIDCUser() (+14 more)

### Community 27 - "New"
Cohesion: 0.07
Nodes (33): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+25 more)

### Community 28 - "DecodeKey"
Cohesion: 0.08
Nodes (30): ProviderConfig, factorJSON(), Handler, testHandler, Handler, Handler, Handler, primaryProviderConfig() (+22 more)

### Community 29 - "NewEngine"
Cohesion: 0.09
Nodes (42): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+34 more)

### Community 30 - "ThemeControls.tsx"
Cohesion: 0.07
Nodes (43): @fontsource/ibm-plex-mono, @fontsource/ibm-plex-sans, @fontsource/space-mono, @fontsource-variable/big-shoulders-text, @fontsource-variable/geist, @fontsource-variable/geist-mono, @fontsource-variable/inter-tight, @fontsource-variable/jetbrains-mono (+35 more)

### Community 31 - "net/http.Request"
Cohesion: 0.08
Nodes (24): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+16 more)

### Community 32 - "SnapshotEntity"
Cohesion: 0.10
Nodes (48): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildDatasets() (+40 more)

### Community 33 - "ConnectorEditPage.tsx"
Cohesion: 0.06
Nodes (38): RFC-3339, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsschema (+30 more)

### Community 34 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (49): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore() (+41 more)

### Community 35 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 36 - "routerDeps"
Cohesion: 0.08
Nodes (39): routerDeps, chi.Router, NewRouter(), spaHandler(), wsRoleLabel(), chi.Router, mountAuthRoutes(), mountMeRoutes() (+31 more)

### Community 37 - "Compare"
Cohesion: 0.07
Nodes (42): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+34 more)

### Community 38 - "ExportToFile"
Cohesion: 0.10
Nodes (44): Export(), ExportToFile(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory() (+36 more)

### Community 39 - "store/backup_test.go"
Cohesion: 0.07
Nodes (41): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+33 more)

### Community 40 - "RunMigrations"
Cohesion: 0.10
Nodes (37): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+29 more)

### Community 41 - "traefik_test.go"
Cohesion: 0.09
Nodes (40): TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword(), TestRegisteredSchema() (+32 more)

### Community 42 - "NewMalformedResponseError"
Cohesion: 0.10
Nodes (41): NewMalformedResponseError(), TestBuildHostsTableAttributes(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+33 more)

### Community 43 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 44 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 45 - "time.Duration"
Cohesion: 0.08
Nodes (26): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+18 more)

### Community 46 - "SystemPage.tsx"
Cohesion: 0.07
Nodes (32): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey (+24 more)

### Community 47 - "docker_test.go"
Cohesion: 0.07
Nodes (38): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+30 more)

### Community 48 - "WebSocketProvider.tsx"
Cohesion: 0.07
Nodes (32): Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`) (+24 more)

### Community 49 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (37): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+29 more)

### Community 50 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 51 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 52 - "response.go"
Cohesion: 0.08
Nodes (24): Handler, DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes(), TestCursorRoundTrip(), TestDecodeCursorRejectsGarbage() (+16 more)

### Community 53 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 54 - "Manager"
Cohesion: 0.10
Nodes (14): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, Store, ReportDefinitionRecord (+6 more)

### Community 55 - "NewUser"
Cohesion: 0.24
Nodes (34): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newTestHandler(), asUser(), createTestShareLink(), futureExpiry() (+26 more)

### Community 56 - "NewEngine"
Cohesion: 0.14
Nodes (30): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+22 more)

### Community 57 - ".Fetch"
Cohesion: 0.11
Nodes (25): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+17 more)

### Community 58 - "Store"
Cohesion: 0.09
Nodes (8): placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange(), ChangeSummary

### Community 59 - "HashPassword"
Cohesion: 0.09
Nodes (22): mustHashDummyPassword(), instanceAdminRoleFor(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable() (+14 more)

### Community 60 - "connector/connector.go"
Cohesion: 0.07
Nodes (19): Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP(), NewServiceUnavailableError(), NewTimeoutError(), TestTypedErrorsWrapAndUnwrap() (+11 more)

### Community 61 - "unifi/tables_test.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 62 - "NewService"
Cohesion: 0.12
Nodes (27): APIKeyChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired() (+19 more)

### Community 63 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 64 - "Connector"
Cohesion: 0.12
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 65 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 66 - "DecodeJSON"
Cohesion: 0.12
Nodes (11): webAuthnFlow, webAuthnUser, Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T (+3 more)

### Community 67 - "New"
Cohesion: 0.09
Nodes (28): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+20 more)

### Community 68 - "NewStore"
Cohesion: 0.15
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 69 - "nilToStr"
Cohesion: 0.10
Nodes (11): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+3 more)

### Community 70 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 71 - "connector_permission.go"
Cohesion: 0.13
Nodes (15): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), Store, apiKeyConnectorFilter(), getConnectorGrant(), ConnectorGrantDiff (+7 more)

### Community 72 - "log/slog.Logger"
Cohesion: 0.11
Nodes (15): newLogger(), formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep() (+7 more)

### Community 73 - "router.go"
Cohesion: 0.12
Nodes (24): validVerb(), go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+16 more)

### Community 74 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 75 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 76 - "main"
Cohesion: 0.11
Nodes (26): main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), NewEmbedRegistry() (+18 more)

### Community 77 - "Register"
Cohesion: 0.13
Nodes (24): init(), init(), init(), init(), init(), init(), init(), init() (+16 more)

### Community 78 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 79 - "newTestHandler"
Cohesion: 0.13
Nodes (25): TestDataSnapshotPaths(), TestSyncsLimitAndIsolation(), createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler() (+17 more)

### Community 80 - "Connector"
Cohesion: 0.20
Nodes (5): SnapshotSection, Connector, parseHosts(), unavailable(), session

### Community 81 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 82 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 83 - "Connector"
Cohesion: 0.12
Nodes (10): init(), Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName(), PathSegment() (+2 more)

### Community 84 - "NewRegistry"
Cohesion: 0.21
Nodes (22): SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders(), TestSuggestWithFallbackStopsOnNonRetryableError() (+14 more)

### Community 85 - "notifications/handlers_test.go"
Cohesion: 0.17
Nodes (23): TestEmbeddedSPAWithoutFrontendBuild(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth() (+15 more)

### Community 86 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 87 - "Connector"
Cohesion: 0.11
Nodes (3): ConfigField, Connector, Connector

### Community 88 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 89 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 90 - "AppearancePage.tsx"
Cohesion: 0.15
Nodes (20): MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+12 more)

### Community 91 - "backup/backup.go"
Cohesion: 0.21
Nodes (22): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+14 more)

### Community 92 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 93 - "git.go"
Cohesion: 0.13
Nodes (18): TestCommitMessage(), writeExportState(), exportCursor, exportState, go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_go_git_go_git_v5, go_pkg_github_com_go_git_go_git_v5_config (+10 more)

### Community 94 - "Store"
Cohesion: 0.13
Nodes (8): fakeStatusChecker, testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 95 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 96 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 97 - "Deps"
Cohesion: 0.20
Nodes (20): registerListAttentionItems(), changeServiceIDs(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs() (+12 more)

### Community 98 - "middleware_test.go"
Cohesion: 0.15
Nodes (16): ConnectorRoleChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireConnectorRole(), RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel() (+8 more)

### Community 99 - "lifecycleManager"
Cohesion: 0.12
Nodes (9): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, gatedStopScheduler, lifecycleDeps, lifecycleManager (+1 more)

### Community 100 - ".Fetch"
Cohesion: 0.15
Nodes (11): ServiceDependency, WantsField(), environmentDependencies(), putMetadata(), unavailable(), TestServiceDependencies(), serviceDependencies(), networkDependencies() (+3 more)

### Community 101 - "connectors_health_test.go"
Cohesion: 0.19
Nodes (13): healthFakeConnector, testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline() (+5 more)

### Community 102 - "middleware.go"
Cohesion: 0.16
Nodes (14): AuditRecorder, contextKey, elevationError, PermissionChecker, UserStatusChecker, elevationFailureReason(), extractBearerToken(), hashToken() (+6 more)

### Community 103 - "lifecycle_test.go"
Cohesion: 0.23
Nodes (11): newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestLifecycleManagerWaitsForSchedulerBeforeCancelAndDBClose(), testLogger(), TestStandbyIsUnreadyAndRunsNoScheduler() (+3 more)

### Community 104 - "handlers_actions_test.go"
Cohesion: 0.25
Nodes (18): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+10 more)

### Community 105 - "httpx/retry_test.go"
Cohesion: 0.29
Nodes (16): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+8 more)

### Community 106 - "ws/ws_test.go"
Cohesion: 0.19
Nodes (18): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+10 more)

### Community 107 - "ReportsPage.tsx"
Cohesion: 0.12
Nodes (18): web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports, web_src_api_generated_reports_reports_usegetreportsdefinitions (+10 more)

### Community 108 - "DiffViewer.tsx"
Cohesion: 0.12
Nodes (12): web_src_api_model_index_diff, DiffCell(), DiffViewer(), Gutter(), LayoutToggle(), SplitCell(), UnifiedLine(), unifiedRows() (+4 more)

### Community 109 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 110 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 111 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 112 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 113 - "Handler"
Cohesion: 0.20
Nodes (8): validateConfigPushRequest(), stepAuditDetail(), validTargetType(), ValidateCompositeRef(), SupportsLifecycleVerb(), Handler, runbookResponse, stepResponse

### Community 114 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 115 - "api/mcp_test.go"
Cohesion: 0.16
Nodes (11): testApp, mintAPIKey(), TestMCPConnectorRestrictedKey(), TestMCPEndToEnd(), go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server (+3 more)

### Community 116 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 117 - "createUser"
Cohesion: 0.21
Nodes (15): ContextWithAPIKeyRestriction(), TestClampConnectorRole(), seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding() (+7 more)

### Community 118 - "MarshalConnectorConfig"
Cohesion: 0.20
Nodes (14): IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly(), TestSecretFieldsChangedFalseOnResubmittedUnchangedSecret() (+6 more)

### Community 119 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 120 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 121 - "templates.fixtures.ts"
Cohesion: 0.16
Nodes (14): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+6 more)

### Community 122 - "net/http.Handler"
Cohesion: 0.17
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 123 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 124 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 125 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 126 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 127 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 128 - "fetch_test.go"
Cohesion: 0.19
Nodes (14): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+6 more)

### Community 129 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 130 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 131 - "IsSecureRequest"
Cohesion: 0.24
Nodes (12): SecurityHeaders(), TestSecurityHeaders(), ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor() (+4 more)

### Community 132 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 133 - "reports/handlers.go"
Cohesion: 0.30
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 134 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 135 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 136 - "compliance_rules_test.go"
Cohesion: 0.22
Nodes (10): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), treatAsSafeFromContext() (+2 more)

### Community 137 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 138 - "dialSSHStdio"
Cohesion: 0.15
Nodes (11): serveOneHTTPExchange(), serveSSHDockerConn(), TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ClientConfig (+3 more)

### Community 139 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 140 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 141 - "Handler"
Cohesion: 0.26
Nodes (4): createVersion(), templateResponse(), templateVersionResponse(), Handler

### Community 142 - "config/validate_test.go"
Cohesion: 0.24
Nodes (11): redactDSN(), redactKVPassword(), Config, TestEveryKeyEnvOverridable(), TestRedactDSN(), TestRedacted(), TestValidate(), TestValidateLeaderElectionRequiresPostgres() (+3 more)

### Community 143 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 144 - "Engine"
Cohesion: 0.23
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 145 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 146 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 147 - "Backend test performance"
Cohesion: 0.17
Nodes (12): Backend test performance, CI job times, Coverage strategy, Deterministic scheduled exports (#408), Follow-ups, Measurements and validation, Measuring, Rules for new tests (+4 more)

### Community 148 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 150 - "findings/handlers_authz_test.go"
Cohesion: 0.45
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 151 - "runbooks/handlers_test.go"
Cohesion: 0.36
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 152 - "chat/chat.go"
Cohesion: 0.25
Nodes (9): cosineSimilarity(), packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips(), TestSplitSections(), unpackVector(), Section (+1 more)

### Community 153 - "vectorCache"
Cohesion: 0.25
Nodes (6): vectorCache, vectorEntry, vectorKey, go_pkg_container_list, container/list.Element, container/list.List

### Community 154 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 155 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 156 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 157 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 158 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 159 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 160 - "httpx/retry.go"
Cohesion: 0.31
Nodes (7): isSafeMethod(), IsSafeMethod(), idempotent(), retryable(), sleep(), RetryPolicy, retryTransport

### Community 162 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 163 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 164 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 165 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 166 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 167 - "ListSchemas"
Cohesion: 0.28
Nodes (8): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), supportedLifecycleVerbs(), ListSchemas(), TestRegisterStubRoundTrips(), go_pkg_github_com_wiselabz_wiselabz_internal_connector_connectortest

### Community 168 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 169 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 170 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 171 - "Elector"
Cohesion: 0.25
Nodes (3): database/sql.Conn, Elector, Noop

### Community 172 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 173 - "auth/handlers_test.go"
Cohesion: 0.25
Nodes (7): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups()

### Community 174 - "chat/handlers_test.go"
Cohesion: 0.54
Nodes (7): Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationDocVisibility(), TestCreateConversationValidation(), TestPostMessageErrors()

### Community 175 - "snapshots_test.go"
Cohesion: 0.57
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 176 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 178 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 179 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 180 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 181 - "runConfigCommand"
Cohesion: 0.33
Nodes (7): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), io.Writer

### Community 182 - "TestBulkReauth"
Cohesion: 0.48
Nodes (7): bulkReq(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync()

### Community 183 - "openapi_contract_test.go"
Cohesion: 0.48
Nodes (6): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 184 - "internal/auth/mfa.go"
Cohesion: 0.38
Nodes (5): GenerateRecoveryCodes(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), go_pkg_github_com_pquerna_otp, go_pkg_github_com_pquerna_otp_totp

### Community 185 - "pfsense.go"
Cohesion: 0.48
Nodes (5): buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName()

### Community 186 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 187 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 189 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 190 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 191 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 192 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 193 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 195 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 196 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 197 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 198 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 199 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 201 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 202 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 203 - "compose-smoke.sh"
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

### Community 211 - "Fixture reuse and lifecycle tests (#405–#407)"
Cohesion: 0.50
Nodes (4): CI measurements, Fixture reuse and lifecycle tests (#405–#407), Local measurements, Validation and coverage

## Knowledge Gaps
- **587 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+582 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1314 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Hub` connect `Store` to `dispatcher_test.go`, `lifecycleManager`, `net/http.ResponseWriter`, `log/slog.Logger`, `go_pkg_net_http`, `ws/ws_test.go`, `time.Duration`, `NewChecker`, `Engine`, `docs/handlers_test.go`, `Dispatcher`, `NewEngine`, `HashPassword`, `net/http.Request`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Election` connect `lifecycleManager` to `go_pkg_context`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Errorf()` connect `net/http.ResponseWriter` to `Handler`, `DecodeJSON`, `routerDeps`, `reports/handlers.go`, `net/http.Handler`, `Handler`, `Handler`, `Store`, `response.go`, `Handler`, `compliance/handlers.go`, `Handler`, `ErrorWithDetails`, `DecodeKey`, `net/http.Request`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _587 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.020724292889241342 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.020675252349460493 - nodes in this community are weakly interconnected._
- **Should `cn` be split into smaller, more focused modules?**
  _Cohesion score 0.03250040829658664 - nodes in this community are weakly interconnected._