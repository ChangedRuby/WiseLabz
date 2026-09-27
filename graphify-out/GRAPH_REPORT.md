# Graph Report - docs-feat-document-single-instance-limits-option  (2026-09-27)

## Corpus Check
- 868 files · ~536,489 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6462 nodes · 20549 edges · 223 communities (206 shown, 17 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1658 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1a003e5f`
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
- react
- go_pkg_time
- cn
- ServiceDetailPage.tsx
- Button.tsx
- newTestHandler
- ConnectorEditPage.tsx
- App.tsx
- unifi/tables.go
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
- Get
- dependencies
- traefik/tables.go
- NewEngine
- package.json
- AuthMiddleware
- RunMigrations
- DecodeKey
- go_pkg_context
- NewChecker
- newTestHandler
- docker_test.go
- NewUser
- GetTypeSchema
- portainer/tables.go
- NewHTTPClient
- NewMalformedResponseError
- NewEngine
- Dispatcher
- net/http.Request
- Store
- Handler
- New
- Connector
- git.go
- ExportToFile
- InstanceAdminFromContext
- mountAPIRoutes
- Connector
- Connector
- settings.mock.ts
- home_assistant_test.go
- rewritePlaceholders
- logging.go
- nilToStr
- Config
- Registry
- main
- .call
- AuthedUser
- NewStore
- router.go
- NewService
- runbooks_test.go
- Service
- .OIDCCallback
- HashToken
- NewRegistry
- NotificationRecord
- handlers.ts
- lifecycleManager
- ServiceSnapshot
- unifi_test.go
- timeline.ts
- go_pkg_reflect
- ValidateConfig
- httpx/retry_test.go
- backup/backup.go
- system/handlers_test.go
- WiseLabz — Design Contract
- devDependencies
- middleware.go
- connector_permission.go
- Manager
- cursor_test.go
- portainer_test.go
- time.Time
- adguardhome_test.go
- Register
- Engine
- templates.fixtures.ts
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
- api/changes_test.go
- MarshalConnectorConfig
- Store
- RunbookRecord
- docdiffmodel.ts
- AppearancePage.tsx
- net/http.Response
- Store
- all.go
- api/docs_test.go
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- testApp
- newTestHandler
- handlers_contract_test.go
- Handler
- Handler
- vectorCache
- Connector
- Connector
- templatefuncs.go
- changes/handlers_test.go
- backup/main.go
- RateLimit
- pagination_contract_test.go
- backup/backup_test.go
- sshStdioConn
- diagram.go
- render_test.go
- connectors_hardening_test.go
- runbook_test.go
- net/http.ResponseWriter
- IsSecureRequest
- newTestHandler
- Contributing to WiseLabz
- Decision
- scripts
- net/http.Handler
- registryTestRefresher
- Store
- Decision
- browser.ts
- Handler
- newTCPDockerClient
- Decision
- 0004 — PostgreSQL leader election for background workers
- WiseLabz Connector Guide
- Product
- .call
- walkCursorPages
- handlers_bulk_test.go
- connectors_maintenance_test.go
- ReadyState
- Elector
- engine_maintenance_test.go
- Saved Views
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- ComplianceRuleRecord
- log/slog.Logger
- WiseLabz — Deployment Guide
- release-please-config.json
- ComputeWindow
- Store
- Cache
- Step by step
- WiseLabz
- webAuthnUser
- snapshotResponse
- .UpdateAuthConfig
- serveSSHDockerConn
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- fakeRefresherConnector
- .GetConnectorUptime
- Store
- Backup Recovery: What Comes Back, and What Doesn't
- Diagnostics Bundle
- Scheduled Doc Export
- Security Policy
- seedScopeFixture
- Enforcement Guidelines
- Authentication design
- Development workflow
- compose-smoke.sh
- timeoutError
- RetentionSettings
- Technology stack
- MISSING — deferred & future frontend features
- truenas/attributes_test.go
- WiseLabz — v2 Backlog
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

## Communities (223 total, 17 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (172): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+164 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (144): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+136 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (178): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown() (+170 more)

### Community 3 - "context.Context"
Cohesion: 0.03
Nodes (27): fakeStatusChecker, sanitizeSessions(), Connector, Connector, Connector, MFAFactor, Store, Store (+19 more)

### Community 4 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (54): msw, react-router-dom, @tanstack/react-query, @testing-library/jest-dom, @testing-library/react, vitest, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridremovalimpact (+46 more)

### Community 5 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.06
Nodes (39): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest, shareLinkContextKey (+31 more)

### Community 6 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (82): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete` (+74 more)

### Community 7 - "go_pkg_testing"
Cohesion: 0.05
Nodes (21): dashboardLayout, badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails() (+13 more)

### Community 8 - "react"
Cohesion: 0.03
Nodes (123): react, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid (+115 more)

### Community 9 - "go_pkg_time"
Cohesion: 0.08
Nodes (24): versionSections(), ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), TemplateVersionSection, contains(), searchString(), go_pkg_database_sql (+16 more)

### Community 10 - "cn"
Cohesion: 0.04
Nodes (62): clsx, tailwind-merge, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_deletetemplatestemplateid, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplates (+54 more)

### Community 11 - "ServiceDetailPage.tsx"
Cohesion: 0.04
Nodes (75): ADR-0001, ADR-0003, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridconfigfields (+67 more)

### Community 12 - "Button.tsx"
Cohesion: 0.04
Nodes (77): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey (+69 more)

### Community 13 - "newTestHandler"
Cohesion: 0.05
Nodes (77): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+69 more)

### Community 14 - "ConnectorEditPage.tsx"
Cohesion: 0.07
Nodes (26): RFC-3339, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_connectors_connectors_putconnectorsconnectorid (+18 more)

### Community 15 - "App.tsx"
Cohesion: 0.02
Nodes (109): Frontend shell & theme (decided 2026-06), Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport (+101 more)

### Community 16 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 17 - "ProfilePage.tsx"
Cohesion: 0.03
Nodes (74): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+66 more)

### Community 18 - "UserIDFromContext"
Cohesion: 0.07
Nodes (30): newToken(), sanitize(), sanitizeUser(), setRefreshCookie(), Handler, mustHashDummyPassword(), Handler, randomOIDCToken() (+22 more)

### Community 19 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 20 - "connector/connector.go"
Cohesion: 0.09
Nodes (12): TimeoutError, NewAuthError(), NewTimeoutError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap(), AuthError, CredentialRefresher, MalformedResponseError (+4 more)

### Community 21 - "net/http.Client"
Cohesion: 0.04
Nodes (27): Connector, ollamaEmbedder, openAIEmbedder, NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL() (+19 more)

### Community 22 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (52): axios, customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers (+44 more)

### Community 23 - "Runner"
Cohesion: 0.07
Nodes (31): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+23 more)

### Community 24 - "icons.tsx"
Cohesion: 0.04
Nodes (82): @codemirror/lang-markdown, @codemirror/view, i18next, @uiw/react-codemirror, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest (+74 more)

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
Nodes (42): fetchAllDocs(), fileName(), Exporter, NewExporter(), RunExportOnce(), slugify(), newTestStore(), readFile() (+34 more)

### Community 29 - "ConnectorRecord"
Cohesion: 0.16
Nodes (13): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), scanSyncRun(), connectorWithRole (+5 more)

### Community 30 - "Errorf"
Cohesion: 0.06
Nodes (29): diffToSpec(), Handler, Handler, decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Handler (+21 more)

### Community 31 - "SnapshotEntity"
Cohesion: 0.12
Nodes (44): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools() (+36 more)

### Community 32 - "rowScanner"
Cohesion: 0.06
Nodes (27): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), existingIDs(), docSearchWhere(), escapeLike() (+19 more)

### Community 33 - "Store"
Cohesion: 0.10
Nodes (23): routerDeps, Handler, Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+15 more)

### Community 34 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 35 - "Get"
Cohesion: 0.10
Nodes (15): Handler, isWritableField(), validateConfigPushRequest(), capitalize(), Handler, WriteElevationError(), ConfigPusher, LifecycleOp() (+7 more)

### Community 36 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 37 - "traefik/tables.go"
Cohesion: 0.11
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 38 - "NewEngine"
Cohesion: 0.25
Nodes (24): NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists(), TestMatchingConnectorsEmptyAppliesToIsWildcard() (+16 more)

### Community 39 - "package.json"
Cohesion: 0.03
Nodes (84): codemirror, @codemirror/commands, @codemirror/state, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, @fontsource/ibm-plex-mono (+76 more)

### Community 40 - "AuthMiddleware"
Cohesion: 0.10
Nodes (24): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, AuthMiddleware(), extractBearerToken(), hashToken(), mfaEnrollmentAllowed() (+16 more)

### Community 41 - "RunMigrations"
Cohesion: 0.11
Nodes (35): main(), OpenDB(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns(), TestMigrationSchemaParity() (+27 more)

### Community 42 - "DecodeKey"
Cohesion: 0.09
Nodes (23): testHandler, Handler, Handler, Handler, Handler, Config, mask(), DecodeKey() (+15 more)

### Community 43 - "go_pkg_context"
Cohesion: 0.05
Nodes (40): buildEmailMessage(), sendNtfyChannel(), sendSMTPChannel(), sendTelegramChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendNtfyChannel_DefaultsToNtfySh(), TestSendNtfyChannel_MissingTopic() (+32 more)

### Community 44 - "NewChecker"
Cohesion: 0.06
Nodes (73): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+65 more)

### Community 45 - "newTestHandler"
Cohesion: 0.11
Nodes (40): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+32 more)

### Community 46 - "docker_test.go"
Cohesion: 0.08
Nodes (32): newDockerClient(), init(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout(), TestDoRequestErrorCases() (+24 more)

### Community 47 - "NewUser"
Cohesion: 0.22
Nodes (36): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestCreateConversationDocVisibility(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser() (+28 more)

### Community 48 - "GetTypeSchema"
Cohesion: 0.14
Nodes (20): TestRegisteredSchema(), TestAllConnectorImplementationsRegister(), supportedLifecycleVerbs(), TestRegisteredSchema(), TestAPIKeyIsStoredAsPassword(), TestRegisteredSchema(), AttributeCatalog(), GetTypeSchema() (+12 more)

### Community 49 - "portainer/tables.go"
Cohesion: 0.08
Nodes (38): WantsField(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), putMetadata(), unavailable(), buildEnvironmentTable(), buildStackTable(), cell() (+30 more)

### Community 50 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 51 - "NewMalformedResponseError"
Cohesion: 0.10
Nodes (42): NewMalformedResponseError(), TestBuildHostsTableAttributes(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+34 more)

### Community 52 - "NewEngine"
Cohesion: 0.14
Nodes (26): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+18 more)

### Community 53 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 54 - "net/http.Request"
Cohesion: 0.08
Nodes (23): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+15 more)

### Community 55 - "Store"
Cohesion: 0.08
Nodes (10): changeServiceIDs(), placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+2 more)

### Community 56 - "Handler"
Cohesion: 0.22
Nodes (4): updateUserRequest, Handler, writeUserWriteError(), NoContent()

### Community 57 - "New"
Cohesion: 0.14
Nodes (14): Handler, SyncDocEmbeddings(), TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(), TestRetrieveUsesCacheAndSyncInvalidates(), TestUpsertBackupSchedulePostgresParity(), Store, newPostgresTestStore(), TestUpsertQualityFindingPostgresDedup() (+6 more)

### Community 58 - "Connector"
Cohesion: 0.16
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

### Community 59 - "git.go"
Cohesion: 0.07
Nodes (31): IsGeneratedName(), pruneStale(), TestIsGeneratedName(), writeExportState(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults(), TestWithinTransactionRollsBack() (+23 more)

### Community 60 - "ExportToFile"
Cohesion: 0.09
Nodes (45): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+37 more)

### Community 61 - "InstanceAdminFromContext"
Cohesion: 0.13
Nodes (11): Handler, contextWithShareLink(), Handler, newShareToken(), shareLinkFromContext(), chi.Router, mountWSRoutes(), RejectRestrictedAPIKey() (+3 more)

### Community 62 - "mountAPIRoutes"
Cohesion: 0.12
Nodes (23): chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountChatRoutes(), mountDocRoutes(), mountShareRoutes() (+15 more)

### Community 63 - "Connector"
Cohesion: 0.07
Nodes (12): init(), ConfigField, Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName() (+4 more)

### Community 64 - "Connector"
Cohesion: 0.19
Nodes (6): unavailable(), SnapshotSection, Connector, unavailable(), unavailable(), session

### Community 65 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 66 - "home_assistant_test.go"
Cohesion: 0.09
Nodes (44): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+36 more)

### Community 67 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 68 - "logging.go"
Cohesion: 0.13
Nodes (19): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestLoggablePathMasksShareTokenUnderV1() (+11 more)

### Community 69 - "nilToStr"
Cohesion: 0.08
Nodes (12): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), FailedSyncRun (+4 more)

### Community 70 - "Config"
Cohesion: 0.13
Nodes (19): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, BackupSettings (+11 more)

### Community 71 - "Registry"
Cohesion: 0.06
Nodes (29): claudeProvider, openAICompatibleProvider, Provider, StatusError, StubProvider, SuggestResult, testProvider, registerFailThenSucceed() (+21 more)

### Community 72 - "main"
Cohesion: 0.09
Nodes (27): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+19 more)

### Community 73 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 74 - "AuthedUser"
Cohesion: 0.12
Nodes (28): TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth(), Handler (+20 more)

### Community 75 - "NewStore"
Cohesion: 0.16
Nodes (26): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+18 more)

### Community 76 - "router.go"
Cohesion: 0.08
Nodes (32): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+24 more)

### Community 77 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 78 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 79 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 80 - ".OIDCCallback"
Cohesion: 0.16
Nodes (8): Handler, newOIDCUser(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider, github.com/coreos/go-oidc/v3/oidc.Provider, golang.org/x/oauth2.Config

### Community 81 - "HashToken"
Cohesion: 0.16
Nodes (13): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+5 more)

### Community 82 - "NewRegistry"
Cohesion: 0.17
Nodes (27): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+19 more)

### Community 83 - "NotificationRecord"
Cohesion: 0.21
Nodes (6): Dispatcher, Dispatcher, RunDeliveryRetries(), NotificationRecord, Store, scanNotification()

### Community 84 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 85 - "lifecycleManager"
Cohesion: 0.16
Nodes (7): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, lifecycleDeps, lifecycleManager

### Community 86 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (24): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, changePatternID(), Engine, markError() (+16 more)

### Community 87 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 88 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 89 - "go_pkg_reflect"
Cohesion: 0.07
Nodes (22): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), Schema() (+14 more)

### Community 90 - "ValidateConfig"
Cohesion: 0.16
Nodes (15): TestSchemaConfigValidation(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), ValidateConfig(), TestSchemaConfigValidation(), TestSchemaConfigValidation(), isConfigValidationError() (+7 more)

### Community 91 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 92 - "backup/backup.go"
Cohesion: 0.24
Nodes (21): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+13 more)

### Community 93 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 94 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 95 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 96 - "middleware.go"
Cohesion: 0.13
Nodes (17): AuditRecorder, ConnectorRoleChecker, contextKey, elevationError, PermissionChecker, UserStatusChecker, chi.Router, mountConnectorRoutes() (+9 more)

### Community 97 - "connector_permission.go"
Cohesion: 0.11
Nodes (18): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), treatAsSafeFromContext(), TreatAsSafeMethod(), Store, apiKeyConnectorFilter() (+10 more)

### Community 98 - "Manager"
Cohesion: 0.09
Nodes (13): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, Store, ReportDefinitionRecord (+5 more)

### Community 99 - "cursor_test.go"
Cohesion: 0.43
Nodes (6): DecodeCursor(), EncodeCursor(), TestCursorRequestModes(), TestCursorRoundTrip(), TestDecodeCursorRejectsGarbage(), TestNextCursorStopsOnShortPage()

### Community 100 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 101 - "time.Time"
Cohesion: 0.15
Nodes (20): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift (+12 more)

### Community 102 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 103 - "Register"
Cohesion: 0.23
Nodes (16): init(), init(), init(), init(), init(), init(), init(), init() (+8 more)

### Community 104 - "Engine"
Cohesion: 0.21
Nodes (7): Engine, dedupKey(), matchReason(), TemplateFuncs(), GenerateResult, renderResult, text/template.FuncMap

### Community 105 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (11): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, resolveToken(), Snapshot (+3 more)

### Community 106 - "Deps"
Cohesion: 0.23
Nodes (18): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+10 more)

### Community 107 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 108 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 109 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 110 - "diagnostics/diagnostics.go"
Cohesion: 0.24
Nodes (16): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+8 more)

### Community 111 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 112 - "Hub"
Cohesion: 0.07
Nodes (17): Config, Sanitize(), TestSanitize(), Store, RunDocLockSweep(), runDocLockSweep(), Engine, Hub (+9 more)

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
Cohesion: 0.16
Nodes (7): AuthSettings, Database, Server, WebAuthnSettings, time.Duration, OIDCProvider, PoolConfig

### Community 117 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 118 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 119 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 120 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

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
Cohesion: 0.16
Nodes (19): zustand, AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+11 more)

### Community 125 - "net/http.Response"
Cohesion: 0.23
Nodes (9): isSafeMethod(), IsSafeMethod(), idempotent(), retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport (+1 more)

### Community 126 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 127 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 128 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 129 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.12
Nodes (16): ADR index, AI module, API design, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27), Data retention (decided 2026-09-05) (+8 more)

### Community 130 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 131 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 132 - "newTestHandler"
Cohesion: 0.24
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 133 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 134 - "Handler"
Cohesion: 0.24
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 135 - "Handler"
Cohesion: 0.26
Nodes (6): stepAuditDetail(), validTargetType(), validVerb(), Handler, runbookResponse, stepResponse

### Community 136 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 137 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 139 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 140 - "changes/handlers_test.go"
Cohesion: 0.32
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 141 - "backup/main.go"
Cohesion: 0.13
Nodes (15): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+7 more)

### Community 142 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 143 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 144 - "backup/backup_test.go"
Cohesion: 0.26
Nodes (13): Export(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+5 more)

### Community 145 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 146 - "diagram.go"
Cohesion: 0.20
Nodes (15): ServiceDependency, environmentDependencies(), networkDependencies(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid() (+7 more)

### Community 147 - "render_test.go"
Cohesion: 0.20
Nodes (17): connectorFilter(), RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden() (+9 more)

### Community 148 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 149 - "runbook_test.go"
Cohesion: 0.32
Nodes (7): Store, newCascadeTestStore(), TestGetRunbookByTarget(), TestRunbookRoundTrip(), TestRunbookStepsCascadeOnConnectorDelete(), TestRunbookStepsCascadeOnRunbookDelete(), TestRunbookStepsRoundTrip()

### Community 150 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (12): webAuthnFlow, Handler, Handler, decodeBulkRequest(), Handler, Handler, WriteDataPaginated(), SinceFromDays() (+4 more)

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

### Community 156 - "net/http.Handler"
Cohesion: 0.16
Nodes (14): TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+6 more)

### Community 158 - "Store"
Cohesion: 0.18
Nodes (7): testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 159 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 160 - "browser.ts"
Cohesion: 0.40
Nodes (4): bootstrap(), worker, enableMocks(), handlers

### Community 162 - "newTCPDockerClient"
Cohesion: 0.15
Nodes (13): GuardedDialer(), IsDangerousIP(), buildDockerTLSConfig(), newTCPDockerClient(), generateSelfSignedCert(), TestNewTCPDockerClientMutualTLS(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair() (+5 more)

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

### Community 168 - "walkCursorPages"
Cohesion: 0.20
Nodes (11): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages(), normalizeParams(), routerOperations() (+3 more)

### Community 169 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 170 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 172 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 173 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 174 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 175 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 176 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 178 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 179 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 180 - "log/slog.Logger"
Cohesion: 0.39
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 181 - "WiseLabz — Deployment Guide"
Cohesion: 0.25
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 183 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 184 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 185 - "Store"
Cohesion: 0.32
Nodes (3): Store, scanBackupRun(), BackupRun

### Community 188 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 189 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 190 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 193 - "webAuthnUser"
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

### Community 226 - "Technology stack"
Cohesion: 0.50
Nodes (4): Backend, Frontend, Infrastructure, Technology stack

### Community 227 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

## Knowledge Gaps
- **572 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+567 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1295 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `testing.T`, `testApp`, `Handler`, `Handler`, `go_pkg_time`, `backup/backup_test.go`, `UserIDFromContext`, `render_test.go`, `net/http.ResponseWriter`, `dispatcher_test.go`, `gitFixture`, `Errorf`, `rowScanner`, `Handler`, `NewEngine`, `.call`, `RunMigrations`, `DecodeKey`, `NewChecker`, `engine_maintenance_test.go`, `NewUser`, `log/slog.Logger`, `Dispatcher`, `net/http.Request`, `NewEngine`, `Handler`, `New`, `ExportToFile`, `rewritePlaceholders`, `main`, `.call`, `AuthedUser`, `NewStore`, `NewRegistry`, `lifecycleManager`, `ServiceSnapshot`, `backup/backup.go`, `Manager`, `time.Time`, `Engine`, `Deps`, `diagnostics/diagnostics.go`, `Hub`, `Handler`, `docs/handlers_test.go`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `WiseLabz — Architecture & Technical Decisions` connect `WiseLabz — Architecture & Technical Decisions` to `Technology stack`, `0004 — PostgreSQL leader election for background workers`, `App.tsx`, `Authentication design`, `Development workflow`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Why does `Frontend shell & theme (decided 2026-06)` connect `App.tsx` to `WiseLabz — Architecture & Technical Decisions`, `Button.tsx`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _572 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021737866899157222 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.02437106918238994 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.018323586744639377 - nodes in this community are weakly interconnected._