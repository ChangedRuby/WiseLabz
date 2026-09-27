# Graph Report - docs-feat-document-single-instance-limits-option  (2026-09-27)

## Corpus Check
- 868 files · ~534,295 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6448 nodes · 20475 edges · 240 communities (217 shown, 23 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1654 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f15c66fc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- context.Context
- @tanstack/react-query
- go_pkg_github_com_wiselabz_wiselabz_internal_store
- DashboardPage.tsx
- go_pkg_testing
- SystemPage.tsx
- go_pkg_context
- cn
- ServiceDetailPage.tsx
- react-i18next
- newTestHandler
- App.tsx
- auth.ts
- react
- ProfilePage.tsx
- UserIDFromContext
- adguardhome/tables.go
- connector/connector.go
- net/http.Client
- UsersPage.tsx
- Runner
- icons.tsx
- dispatcher_test.go
- Compare
- fixtures.ts
- gitFixture
- ConnectorRecord
- Errorf
- SnapshotEntity
- rowScanner
- Store
- home_assistant/tables.go
- net/http.ResponseWriter
- dependencies
- traefik/tables.go
- NewEngine
- package.json
- NewService
- RunMigrations
- DecodeKey
- channels.go
- NewChecker
- newTestHandler
- docker_test.go
- NewUser
- GetTypeSchema
- portainer/tables.go
- Register
- NewMalformedResponseError
- NewEngine
- Dispatcher
- net/http.Request
- Store
- response.go
- New
- compliance/engine.go
- git.go
- ExportToFile
- NewRegistry
- mountAPIRoutes
- Connector
- Connector
- settings.mock.ts
- home_assistant_test.go
- rewritePlaceholders
- backup/main.go
- nilToStr
- Config
- SuggestRequest
- main
- HashPassword
- AuthedUser
- NewStore
- router.go
- config_test.go
- runbooks_test.go
- Service
- .OIDCCallback
- HashToken
- ContextWithUser
- log/slog.Logger
- handlers.ts
- newTestLifecycle
- ServiceSnapshot
- unifi_test.go
- timeline.ts
- api/auth/oidc.go
- fetch_test.go
- httpx/retry_test.go
- backup/backup.go
- Checker
- WiseLabz — Design Contract
- devDependencies
- middleware.go
- Store
- Manager
- WritePaginated
- portainer_test.go
- time.Time
- adguardhome_test.go
- src/theme.ts
- .Fetch
- DocRecord
- Deps
- Connector
- Handler
- truenas_test.go
- diagnostics/diagnostics.go
- keyset_test.go
- Hub
- ws/ws_test.go
- ws.ts
- compilerOptions
- time.Duration
- Handler
- docs/handlers_test.go
- traefik_test.go
- MarshalConnectorConfig
- Store
- RunbookRecord
- docdiffmodel.ts
- AppearancePage.tsx
- templates_test.go
- api/mcp_test.go
- all.go
- .batchDelete
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- testApp
- go_pkg_encoding_base64
- handlers_contract_test.go
- Handler
- Handler
- vectorCache
- Connector
- Connector
- doc/engine.go
- changes/handlers_test.go
- chat/chat.go
- net/http.Handler
- pagination_contract_test.go
- backup/backup_test.go
- sshStdioConn
- diagram.go
- render_test.go
- bulkFakeConnector
- store/theme.ts
- Handler
- IsSecureRequest
- newTestHandler
- Contributing to WiseLabz
- Decision
- scripts
- Config
- ReportData
- Store
- Decision
- main.tsx
- Handler
- newTCPDockerClient
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- .call
- cursor_pagination_test.go
- handlers_bulk_test.go
- connectors_maintenance_test.go
- .Fetch
- Elector
- Handler
- transform.go
- Changelog
- mockServiceWorker.js
- ThemeControls.tsx
- apikey_scopes_test.go
- ComplianceRuleRecord
- retention/retention_test.go
- WiseLabz — Deployment Guide
- useTheme
- release-please-config.json
- ComputeWindow
- Store
- ShareLink
- Engine
- Cache
- Step by step
- WiseLabz
- settings.ts
- RequirePermission
- webAuthnUser
- TestComplianceRuleValidation
- snapshotResponse
- .UpdateAuthConfig
- dialSSHStdio
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- fakeRefresherConnector
- noopValidatedConnector
- .GetConnectorUptime
- Store
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Security Policy
- RequireConnectorRole
- routerOperations
- seedScopeFixture
- Enforcement Guidelines
- Authentication design
- Development workflow
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- RetentionSettings
- transform_firewall.go
- Technology stack
- MISSING — deferred & future frontend features
- internal/auth/oidc.go
- fields_test.go
- truenas/attributes_test.go
- WiseLabz — v2 Backlog
- fakeEmbedder
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
5. `UserIDFromContext()` - 84 edges
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

## Communities (240 total, 23 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (176): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+168 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (147): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+139 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (149): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+141 more)

### Community 3 - "context.Context"
Cohesion: 0.03
Nodes (28): fakeStatusChecker, sanitizeSessions(), Connector, Connector, existingIDs(), Store, MFAFactor, Store (+20 more)

### Community 4 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (68): i18next, msw, react-error-boundary, react-router-dom, sonner, @tanstack/react-query, @testing-library/react, vitest (+60 more)

### Community 5 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.05
Nodes (45): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest, go_pkg_encoding_csv (+37 more)

### Community 6 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (92): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete` (+84 more)

### Community 7 - "go_pkg_testing"
Cohesion: 0.07
Nodes (23): TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), go_pkg_bytes, go_pkg_encoding_json, go_pkg_github_com_gorilla_websocket (+15 more)

### Community 8 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (79): web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey, web_src_api_generated_settings_settings_getgetauthconfigquerykey, web_src_api_generated_settings_settings_getgetnotificationsconfigquerykey, web_src_api_generated_settings_settings_postaiconfigtest, web_src_api_generated_settings_settings_postnotificationsconfigtest, web_src_api_generated_settings_settings_putaiconfig (+71 more)

### Community 9 - "go_pkg_context"
Cohesion: 0.07
Nodes (16): StatusError, dashboardLayout, contains(), searchString(), shareLinkContextKey, go_pkg_context, go_pkg_database_sql, go_pkg_errors (+8 more)

### Community 10 - "cn"
Cohesion: 0.03
Nodes (76): web_src_api_generated_docs_docs_usegetdocsdocid, web_src_api_generated_docs_docs_usegetdocstree, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore (+68 more)

### Community 11 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (82): ADR-0001, ADR-0003, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid (+74 more)

### Community 12 - "react-i18next"
Cohesion: 0.04
Nodes (76): Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, motion, @radix-ui/react-popover, react-i18next (+68 more)

### Community 13 - "newTestHandler"
Cohesion: 0.06
Nodes (74): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+66 more)

### Community 14 - "App.tsx"
Cohesion: 0.04
Nodes (63): RFC-3339, setMfaEnrollmentRequiredHandler(), web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectors (+55 more)

### Community 15 - "auth.ts"
Cohesion: 0.04
Nodes (71): Frontend shell & theme (decided 2026-06), Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport (+63 more)

### Community 16 - "react"
Cohesion: 0.04
Nodes (73): react, react-markdown, remark-gfm, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain (+65 more)

### Community 17 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (56): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+48 more)

### Community 18 - "UserIDFromContext"
Cohesion: 0.06
Nodes (31): Handler, newToken(), sanitize(), sanitizeUser(), setRefreshCookie(), Handler, Handler, Handler (+23 more)

### Community 19 - "adguardhome/tables.go"
Cohesion: 0.07
Nodes (61): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+53 more)

### Community 20 - "connector/connector.go"
Cohesion: 0.05
Nodes (26): ServiceDependency, TimeoutError, NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), Connector, TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap() (+18 more)

### Community 21 - "net/http.Client"
Cohesion: 0.05
Nodes (18): Connector, ollamaEmbedder, openAIEmbedder, Connector, LimitedBody(), ReadBody(), TestReadBodyLimit(), CheckStatus() (+10 more)

### Community 22 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (46): axios, customInstance(), web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa (+38 more)

### Community 23 - "Runner"
Cohesion: 0.07
Nodes (31): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+23 more)

### Community 24 - "icons.tsx"
Cohesion: 0.07
Nodes (47): @codemirror/lang-markdown, @codemirror/view, @uiw/react-codemirror, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest (+39 more)

### Community 25 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 26 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 27 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 28 - "gitFixture"
Cohesion: 0.07
Nodes (34): Exporter, NewExporter(), RunExportOnce(), newTestStore(), readFile(), TestExportAllEmptyDirName(), TestExportAllKeepsOperatorFiles(), TestExportAllPrunesStaleFilesOnRerun() (+26 more)

### Community 29 - "ConnectorRecord"
Cohesion: 0.08
Nodes (22): Sanitize(), TestSanitize(), Store, ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr() (+14 more)

### Community 30 - "Errorf"
Cohesion: 0.09
Nodes (10): diffToSpec(), Handler, Handler, Handler, Handler, Handler, Handler, Errorf() (+2 more)

### Community 31 - "SnapshotEntity"
Cohesion: 0.13
Nodes (42): SnapshotEntity, buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices() (+34 more)

### Community 32 - "rowScanner"
Cohesion: 0.08
Nodes (22): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), countRows(), T, keysetQuery() (+14 more)

### Community 33 - "Store"
Cohesion: 0.08
Nodes (27): routerDeps, Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+19 more)

### Community 34 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 35 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (19): Handler, isWritableField(), validateConfigPushRequest(), capitalize(), Handler, decodeBulkRequest(), Handler, WriteElevationError() (+11 more)

### Community 36 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 37 - "traefik/tables.go"
Cohesion: 0.11
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 38 - "NewEngine"
Cohesion: 0.13
Nodes (31): Engine, NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists() (+23 more)

### Community 39 - "package.json"
Cohesion: 0.05
Nodes (36): clsx, codemirror, @codemirror/commands, @codemirror/state, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker (+28 more)

### Community 40 - "NewService"
Cohesion: 0.10
Nodes (35): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService() (+27 more)

### Community 41 - "RunMigrations"
Cohesion: 0.11
Nodes (34): main(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns(), TestMigrationSchemaParity(), RunMigrations() (+26 more)

### Community 42 - "DecodeKey"
Cohesion: 0.11
Nodes (22): ProviderConfig, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey(), Decrypt(), DeriveKey() (+14 more)

### Community 43 - "channels.go"
Cohesion: 0.08
Nodes (30): isSafeMethod(), treatAsSafeFromContext(), TreatAsSafeMethod(), IsTimeout(), TestIsTimeout(), IsSafeMethod(), idempotent(), buildEmailMessage() (+22 more)

### Community 44 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 45 - "newTestHandler"
Cohesion: 0.11
Nodes (37): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+29 more)

### Community 46 - "docker_test.go"
Cohesion: 0.07
Nodes (38): newDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn(), startSSHDockerServer(), TestConfigPush() (+30 more)

### Community 47 - "NewUser"
Cohesion: 0.21
Nodes (37): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestCreateConversationDocVisibility(), TestGetRequiresGrantEvenForInstanceAdmin(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler() (+29 more)

### Community 48 - "GetTypeSchema"
Cohesion: 0.09
Nodes (34): validateConnectorConfig(), TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), supportedLifecycleVerbs(), TestRegisteredSchema(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+26 more)

### Community 49 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 50 - "Register"
Cohesion: 0.09
Nodes (32): init(), init(), newConnector(), Connector, init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback() (+24 more)

### Community 51 - "NewMalformedResponseError"
Cohesion: 0.13
Nodes (35): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+27 more)

### Community 52 - "NewEngine"
Cohesion: 0.12
Nodes (31): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+23 more)

### Community 53 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 54 - "net/http.Request"
Cohesion: 0.10
Nodes (19): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+11 more)

### Community 55 - "Store"
Cohesion: 0.08
Nodes (10): changeServiceIDs(), placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+2 more)

### Community 56 - "response.go"
Cohesion: 0.08
Nodes (17): updateUserRequest, Handler, writeUserWriteError(), Error(), DataPaginatedResponse, HandleStoreError(), intQuery(), JSON() (+9 more)

### Community 57 - "New"
Cohesion: 0.07
Nodes (34): confirm(), formatCounts(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag() (+26 more)

### Community 58 - "compliance/engine.go"
Cohesion: 0.11
Nodes (31): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+23 more)

### Community 59 - "git.go"
Cohesion: 0.07
Nodes (24): IsGeneratedName(), pruneStale(), TestIsGeneratedName(), TestCommitMessage(), dockerSSHAddr, go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_getkin_kin_openapi_openapi3 (+16 more)

### Community 60 - "ExportToFile"
Cohesion: 0.14
Nodes (32): ExportToFile(), ImportFromFile(), AppVersion(), BuildManifest(), BundleCounts(), ChecksumBytes(), ManifestPath(), ReadManifest() (+24 more)

### Community 61 - "NewRegistry"
Cohesion: 0.14
Nodes (25): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+17 more)

### Community 62 - "mountAPIRoutes"
Cohesion: 0.11
Nodes (26): chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router, mountChatRoutes() (+18 more)

### Community 63 - "Connector"
Cohesion: 0.09
Nodes (11): init(), ConfigField, Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName() (+3 more)

### Community 64 - "Connector"
Cohesion: 0.14
Nodes (12): SnapshotSection, MapTransportError(), TestMapTransportError(), TestBuildHostsTableV5(), buildHostsTable(), Connector, parseHosts(), TestBuildHostsTableMalformedCases() (+4 more)

### Community 65 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 66 - "home_assistant_test.go"
Cohesion: 0.13
Nodes (29): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+21 more)

### Community 67 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 68 - "backup/main.go"
Cohesion: 0.11
Nodes (23): main(), usage(), loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken() (+15 more)

### Community 69 - "nilToStr"
Cohesion: 0.10
Nodes (12): seedDelivery(), ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus (+4 more)

### Community 70 - "Config"
Cohesion: 0.11
Nodes (22): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+14 more)

### Community 71 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 72 - "main"
Cohesion: 0.11
Nodes (23): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+15 more)

### Community 73 - "HashPassword"
Cohesion: 0.11
Nodes (19): mustHashDummyPassword(), testHandler, Handler, instanceAdminRoleFor(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), Handler (+11 more)

### Community 74 - "AuthedUser"
Cohesion: 0.12
Nodes (28): TestEmbeddedSPAWithoutFrontendBuild(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth() (+20 more)

### Community 75 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 76 - "router.go"
Cohesion: 0.12
Nodes (24): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+16 more)

### Community 77 - "config_test.go"
Cohesion: 0.11
Nodes (26): Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH(), TestLoadFromYAML() (+18 more)

### Community 78 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 79 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 80 - ".OIDCCallback"
Cohesion: 0.15
Nodes (9): Handler, newOIDCUser(), randomOIDCToken(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider, github.com/coreos/go-oidc/v3/oidc.Provider (+1 more)

### Community 81 - "HashToken"
Cohesion: 0.15
Nodes (15): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+7 more)

### Community 82 - "ContextWithUser"
Cohesion: 0.11
Nodes (27): TestConnectorStoreErrorPaths(), Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound() (+19 more)

### Community 83 - "log/slog.Logger"
Cohesion: 0.11
Nodes (14): formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+6 more)

### Community 84 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 85 - "newTestLifecycle"
Cohesion: 0.09
Nodes (13): newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestStandbyIsUnreadyAndRunsNoScheduler(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group (+5 more)

### Community 86 - "ServiceSnapshot"
Cohesion: 0.09
Nodes (7): ServiceSnapshot, agentEnabled(), Connector, blockingConnector, fakeConnector, fieldsRecordingConnector, sequentialConnector

### Community 87 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 88 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 89 - "api/auth/oidc.go"
Cohesion: 0.10
Nodes (17): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), Config (+9 more)

### Community 90 - "fetch_test.go"
Cohesion: 0.11
Nodes (23): extractGroups(), newMockOIDCServer(), TestAuthURLAfterInitialization(), TestAuthURLBeforeInitialization(), TestExtractGroups(), TestInitializeFailure(), TestInitializeInvalidJSON(), TestInitializeSuccess() (+15 more)

### Community 91 - "httpx/retry_test.go"
Cohesion: 0.20
Nodes (20): retryable(), RetryTransport(), sleep(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry() (+12 more)

### Community 92 - "backup/backup.go"
Cohesion: 0.20
Nodes (23): versionSections(), connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+15 more)

### Community 93 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 94 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 95 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 96 - "middleware.go"
Cohesion: 0.11
Nodes (15): APIKeyChecker, AuditRecorder, contextKey, elevationError, testAPIKeyChecker, UserStatusChecker, APIKeyClaims, elevationFailureReason() (+7 more)

### Community 97 - "Store"
Cohesion: 0.14
Nodes (12): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), getConnectorGrant(), ConnectorGrantDiff, Store, highestConnectorRole() (+4 more)

### Community 98 - "Manager"
Cohesion: 0.16
Nodes (8): cron.EntryID, Manager, LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store, Scheduler

### Community 99 - "WritePaginated"
Cohesion: 0.15
Nodes (16): decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+8 more)

### Community 100 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 101 - "time.Time"
Cohesion: 0.17
Nodes (19): digestDue(), TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection (+11 more)

### Community 102 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 103 - "src/theme.ts"
Cohesion: 0.10
Nodes (16): @fontsource/ibm-plex-mono, @fontsource/ibm-plex-sans, @fontsource/space-mono, @fontsource-variable/big-shoulders-text, @fontsource-variable/geist, @fontsource-variable/geist-mono, @fontsource-variable/inter-tight, @fontsource-variable/jetbrains-mono (+8 more)

### Community 104 - ".Fetch"
Cohesion: 0.16
Nodes (8): WantsField(), TestBuildContainerTableAttributes(), buildContainerTable(), Connector, putMetadata(), unavailable(), Connector, dockerSectionSpec

### Community 105 - "DocRecord"
Cohesion: 0.14
Nodes (10): fetchAllDocs(), fileName(), slugify(), docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc() (+2 more)

### Community 106 - "Deps"
Cohesion: 0.22
Nodes (19): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs(), registerListFindings() (+11 more)

### Community 107 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 108 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 109 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 110 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 111 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 112 - "Hub"
Cohesion: 0.15
Nodes (6): Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator

### Community 113 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 114 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 115 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 116 - "time.Duration"
Cohesion: 0.14
Nodes (5): healthFakeConnector, Database, Server, time.Duration, PoolConfig

### Community 117 - "Handler"
Cohesion: 0.25
Nodes (6): response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 118 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 119 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 120 - "MarshalConnectorConfig"
Cohesion: 0.20
Nodes (14): IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly(), TestSecretFieldsChangedFalseOnResubmittedUnchangedSecret() (+6 more)

### Community 121 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 122 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 123 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 124 - "AppearancePage.tsx"
Cohesion: 0.18
Nodes (15): zustand, AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+7 more)

### Community 125 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 126 - "api/mcp_test.go"
Cohesion: 0.17
Nodes (8): go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, go_pkg_github_com_wiselabz_wiselabz_internal_chat, changeSummary, connectorSummary, findingSummary

### Community 127 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 129 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 130 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 131 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 132 - "go_pkg_encoding_base64"
Cohesion: 0.19
Nodes (12): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+4 more)

### Community 133 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 134 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 135 - "Handler"
Cohesion: 0.24
Nodes (6): stepAuditDetail(), validTargetType(), validVerb(), Handler, runbookResponse, stepResponse

### Community 136 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 137 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 138 - "Connector"
Cohesion: 0.17
Nodes (6): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), Connector

### Community 139 - "doc/engine.go"
Cohesion: 0.17
Nodes (11): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+3 more)

### Community 140 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 141 - "chat/chat.go"
Cohesion: 0.19
Nodes (12): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+4 more)

### Community 142 - "net/http.Handler"
Cohesion: 0.15
Nodes (12): TestRateLimit(), RateLimit(), SecurityHeaders(), TestSecurityHeaders(), RequireInstanceAdmin(), contextWithInstanceAdmin(), TestRequireInstanceAdmin(), golang.org/x/time/rate.Limit (+4 more)

### Community 143 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 144 - "backup/backup_test.go"
Cohesion: 0.26
Nodes (13): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+5 more)

### Community 145 - "sshStdioConn"
Cohesion: 0.15
Nodes (7): closeQuietly(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session, io.Closer, io.WriteCloser, net.Addr

### Community 146 - "diagram.go"
Cohesion: 0.26
Nodes (12): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), relatedEntities(), EntityLink (+4 more)

### Community 147 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 148 - "bulkFakeConnector"
Cohesion: 0.14
Nodes (3): actionConnector, bulkFakeConnector, failingPushConnector

### Community 149 - "store/theme.ts"
Cohesion: 0.25
Nodes (13): ColorMode, commit(), Persisted, PRESETS_FONTS, ThemeState, tokensFor(), ACTIVE, applyTokens() (+5 more)

### Community 150 - "Handler"
Cohesion: 0.35
Nodes (3): webAuthnFlow, Handler, github.com/go-webauthn/webauthn/webauthn.SessionData

### Community 151 - "IsSecureRequest"
Cohesion: 0.27
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 152 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 153 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 154 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 155 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 156 - "Config"
Cohesion: 0.21
Nodes (10): Config, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), chi.Router, NewRouter(), spaHandler() (+2 more)

### Community 157 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 158 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 159 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 160 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 162 - "newTCPDockerClient"
Cohesion: 0.18
Nodes (11): GuardedDialer(), IsDangerousIP(), buildDockerTLSConfig(), newTCPDockerClient(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair(), Unwrap(), newWebhookClient() (+3 more)

### Community 163 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 164 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.18
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 165 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 166 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 167 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 168 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 169 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 170 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 171 - ".Fetch"
Cohesion: 0.24
Nodes (5): setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector

### Community 172 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 174 - "transform.go"
Cohesion: 0.31
Nodes (7): init(), RegisterTransformer(), runTransformers(), TestRunTransformersAppliesInOrderAndStopsOnError(), TestRunTransformersUnknownCategoryIsNoop(), Transformer, TransformerFunc

### Community 175 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 176 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 177 - "ThemeControls.tsx"
Cohesion: 0.24
Nodes (7): AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls(), makePalette()

### Community 178 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 179 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 180 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 181 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 182 - "useTheme"
Cohesion: 0.33
Nodes (7): mermaid, cssVar(), Mermaid(), resolveColor(), load(), useTheme, presetOpts()

### Community 183 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 184 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 185 - "Store"
Cohesion: 0.32
Nodes (3): Store, scanBackupRun(), BackupRun

### Community 187 - "Engine"
Cohesion: 0.29
Nodes (3): Engine, sync.Map, DocRegenerator

### Community 188 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 189 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 190 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 191 - "settings.ts"
Cohesion: 0.46
Nodes (6): MotionProvider(), apply(), framerReducedMotion(), seed(), SettingsState, useSettings

### Community 192 - "RequirePermission"
Cohesion: 0.38
Nodes (5): PermissionChecker, chi.Router, mountDashboardRoutes(), mountWorkflowRoutes(), RequirePermission()

### Community 193 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 194 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 195 - "snapshotResponse"
Cohesion: 0.48
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

### Community 196 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 197 - "dialSSHStdio"
Cohesion: 0.29
Nodes (6): TestDialSSHStdioHonorsContextCancel(), dialSSHStdio(), bufio.ReadWriter, golang.org/x/crypto/ssh.ClientConfig, net.Conn, Options

### Community 198 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 199 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 200 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 201 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 202 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 203 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 204 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 207 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 209 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 210 - "Diagnostics Bundle"
Cohesion: 0.33
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 211 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 212 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 213 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 214 - "routerOperations"
Cohesion: 0.50
Nodes (5): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes

### Community 215 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 216 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 217 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 218 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 219 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 221 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 224 - "transform_firewall.go"
Cohesion: 0.67
Nodes (3): normalizeEnabledColumn(), normalizeFirewallRules(), TestNormalizeFirewallRulesRewritesEnabledColumn()

### Community 226 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 227 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **572 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+567 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1297 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `gitFixture` connect `gitFixture` to `Store`, `testing.T`, `context.Context`, `log/slog.Logger`, `git.go`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Why does `Connector` connect `Connector` to `Connector`, `net/http.Client`, `go_pkg_testing`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Why does `UserIDFromContext()` connect `UserIDFromContext` to `context.Context`, `Handler`, `Handler`, `Handler`, `IsSecureRequest`, `Errorf`, `rowScanner`, `Handler`, `Store`, `net/http.ResponseWriter`, `NewService`, `Handler`, `net/http.Request`, `response.go`, `mountAPIRoutes`, `RequirePermission`, `HashToken`, `RequireConnectorRole`, `middleware.go`, `WritePaginated`, `Deps`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _572 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021052631578947368 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.024764327397875204 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.021654326884802167 - nodes in this community are weakly interconnected._