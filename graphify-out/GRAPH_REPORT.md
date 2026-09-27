# Graph Report - gsaraiva2109-feat-docexport-per-revision-commits  (2026-09-27)

## Corpus Check
- 865 files · ~534,402 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6419 nodes · 20443 edges · 213 communities (193 shown, 20 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1654 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `94798326`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- newDocTestStore
- react
- App.tsx
- context.Context
- package.json
- go_pkg_net_http
- icons.tsx
- @tanstack/react-query
- DashboardPage.tsx
- Errorf
- Button.tsx
- ProfilePage.tsx
- newTestHandler
- go_pkg_context
- ServiceDetailPage.tsx
- NewStore
- cn
- ServiceSnapshot
- go_pkg_testing
- newTestApp
- net/http.Client
- go_pkg_github_com_wiselabz_wiselabz_internal_connector
- gitFixture
- UsersPage.tsx
- net/http.Request
- NewUser
- Connector
- rowScanner
- NewEngine
- Runner
- dispatcher_test.go
- newTestHandler
- net/http.ResponseWriter
- Compare
- fixtures.ts
- WiseLabz — Architecture & Technical Decisions
- RunMigrations
- SnapshotEntity
- routerDeps
- docker_test.go
- ErrorWithDetails
- go_pkg_os
- go_pkg_strings
- home_assistant/tables.go
- DecodeKey
- dependencies
- response.go
- User
- chat/chat.go
- NewChecker
- SuggestWithFallback
- Store
- portainer/tables.go
- compliance/engine.go
- Dispatcher
- ConnectorEditPage.tsx
- Connector
- adguardhome/tables.go
- ExportToFile
- traefik_test.go
- traefik/tables.go
- auth_test.go
- router.go
- NewMalformedResponseError
- unifi/tables.go
- Store
- Configuration & Documentation Backup (Export/Import)
- NewEngine
- settings.mock.ts
- log/slog.Logger
- nilToStr
- rewritePlaceholders
- api/audit_test.go
- httpx/retry_test.go
- AuthMiddleware
- connector_permission.go
- Config
- main
- runbooks_test.go
- Service
- HashToken
- GetTypeSchema
- connector/connector.go
- ConnectorRecord
- handlers.ts
- Hub
- unifi_test.go
- timeline.ts
- Checker
- Manager
- WiseLabz — Design Contract
- devDependencies
- newTestHandler
- portainer_test.go
- AppearancePage.tsx
- adguardhome_test.go
- NewClient
- Connector
- RunbookRecord
- time.Duration
- NewService
- backup/backup.go
- config_test.go
- NewHTTPClient
- home_assistant_test.go
- Handler
- MarshalConnectorConfig
- truenas_test.go
- diagnostics/diagnostics.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- Store
- logging_test.go
- Handler
- Connector
- time.Time
- Store
- docdiffmodel.ts
- templates_test.go
- system/handlers_test.go
- Register
- all.go
- data.go
- .batchDelete
- compilerOptions
- testApp
- OpenDB
- handlers_contract_test.go
- api/attention_test.go
- Handler
- Handler
- createUser
- VerifyBundleFile
- vectorCache
- config/validate_test.go
- pagination_contract_test.go
- Connector
- render_test.go
- change_pattern_test.go
- templates.fixtures.ts
- Handler
- Config
- runRestore
- config_cmd_test.go
- api/changes_test.go
- connectors_health_test.go
- Contributing to WiseLabz
- scripts
- templatefuncs.go
- IsSecureRequest
- Decision
- Connector
- lifecycleDeps
- registry.go
- server.go
- ReportData
- Connector
- Decision
- WiseLabz Connector Guide
- Product
- .call
- net/http.Handler
- cursor_pagination_test.go
- connectors_maintenance_test.go
- api/docs_test.go
- ratelimit.go
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- ComplianceRuleRecord
- connectors_hardening_test.go
- retention/retention_test.go
- BackupSchedule
- release-please-config.json
- ComputeWindow
- ShareLink
- Engine
- Step by step
- WiseLabz
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- fakeRefresherConnector
- .GetConnectorUptime
- Store
- engine_maintenance_test.go
- Security Policy
- browser.ts
- Enforcement Guidelines
- compose-smoke.sh
- .operatorWithPermission
- mintAPIKey
- ClassifyHealth
- timeoutError
- RetentionSettings
- MISSING — deferred & future frontend features
- Saved Views
- stubEmbedder
- WiseLabz — v2 Backlog
- fakeEmbedder
- tsconfig.json
- AGENTS.md
- versionSections
- setup-env.sh
- CHANGE_PROVENANCE.md
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 229 edges
2. `Errorf()` - 183 edges
3. `newDocTestStore()` - 143 edges
4. `Store` - 141 edges
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

## Communities (213 total, 20 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.02
Nodes (160): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+152 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.03
Nodes (143): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+135 more)

### Community 2 - "react"
Cohesion: 0.03
Nodes (123): react, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid (+115 more)

### Community 3 - "App.tsx"
Cohesion: 0.02
Nodes (109): Frontend shell & theme (decided 2026-06), Sync flow, Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport (+101 more)

### Community 4 - "context.Context"
Cohesion: 0.03
Nodes (27): fakeStatusChecker, sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, Connector, Connector, Connector, existingIDs() (+19 more)

### Community 5 - "package.json"
Cohesion: 0.03
Nodes (84): codemirror, @codemirror/commands, @codemirror/state, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, @fontsource/ibm-plex-mono (+76 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (36): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest (+28 more)

### Community 7 - "icons.tsx"
Cohesion: 0.04
Nodes (82): @codemirror/lang-markdown, @codemirror/view, i18next, @uiw/react-codemirror, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest (+74 more)

### Community 8 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (54): msw, react-router-dom, @tanstack/react-query, @testing-library/jest-dom, @testing-library/react, vitest, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridremovalimpact (+46 more)

### Community 9 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (82): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 1. `service.status`, 2. `sync.progress`, 3. `sync.complete` (+74 more)

### Community 10 - "Errorf"
Cohesion: 0.05
Nodes (27): Handler, PermissionChecker, newToken(), sanitize(), Handler, diffToSpec(), Handler, Handler (+19 more)

### Community 11 - "Button.tsx"
Cohesion: 0.04
Nodes (77): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey (+69 more)

### Community 12 - "ProfilePage.tsx"
Cohesion: 0.03
Nodes (74): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+66 more)

### Community 13 - "newTestHandler"
Cohesion: 0.06
Nodes (73): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+65 more)

### Community 14 - "go_pkg_context"
Cohesion: 0.08
Nodes (13): StatusError, go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_go_webauthn_webauthn_protocol, go_pkg_github_com_google_uuid, go_pkg_github_com_jackc_pgx_v5_stdlib (+5 more)

### Community 15 - "ServiceDetailPage.tsx"
Cohesion: 0.04
Nodes (75): ADR-0001, ADR-0003, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridconfigfields (+67 more)

### Community 16 - "NewStore"
Cohesion: 0.06
Nodes (80): TestEmbeddedSPAWithoutFrontendBuild(), NewRegistry(), NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound() (+72 more)

### Community 17 - "cn"
Cohesion: 0.04
Nodes (62): clsx, tailwind-merge, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_deletetemplatestemplateid, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplates (+54 more)

### Community 18 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (24): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, changePatternID(), Engine, markError() (+16 more)

### Community 19 - "go_pkg_testing"
Cohesion: 0.05
Nodes (25): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), TestSchemaMatchesConfig() (+17 more)

### Community 20 - "newTestApp"
Cohesion: 0.04
Nodes (70): TestAPIKeyCreateRejectsInvalidExpiryAndEmptyName(), TestAPIKeyRoutesEndToEnd(), TestBackupExportImportRoundTrip(), TestBackupExportRedactsSecrets(), TestBackupExportRoleBoundary(), TestBackupImportBadVersion(), TestBackupImportMalformedJSON(), TestBackupImportRejectsOversizedBody() (+62 more)

### Community 21 - "net/http.Client"
Cohesion: 0.04
Nodes (29): Connector, ollamaEmbedder, openAIEmbedder, NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL() (+21 more)

### Community 22 - "go_pkg_github_com_wiselabz_wiselabz_internal_connector"
Cohesion: 0.05
Nodes (39): contextKey, elevationError, buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), TestBuildInterfaceTableAttributes(), buildGatewayTable() (+31 more)

### Community 23 - "gitFixture"
Cohesion: 0.07
Nodes (44): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+36 more)

### Community 24 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (52): axios, customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers (+44 more)

### Community 25 - "net/http.Request"
Cohesion: 0.06
Nodes (30): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+22 more)

### Community 26 - "NewUser"
Cohesion: 0.10
Nodes (59): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestCreateConversationDocVisibility(), TestListFiltersGrantsBeforePagination(), Handler, snapshotFixture(), snapshotRequest() (+51 more)

### Community 27 - "Connector"
Cohesion: 0.09
Nodes (38): SnapshotSection, buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+30 more)

### Community 28 - "rowScanner"
Cohesion: 0.06
Nodes (27): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), docSearchWhere(), escapeLike(), DocRecord (+19 more)

### Community 29 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 30 - "Runner"
Cohesion: 0.07
Nodes (31): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+23 more)

### Community 31 - "dispatcher_test.go"
Cohesion: 0.16
Nodes (53): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), testLogger(), expireAlertsOnce(), NewDispatcher() (+45 more)

### Community 32 - "newTestHandler"
Cohesion: 0.09
Nodes (49): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+41 more)

### Community 33 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (15): webAuthnFlow, Handler, applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), Handler, Handler, oidcProviderJSON() (+7 more)

### Community 34 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 35 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 36 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (46): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+38 more)

### Community 37 - "RunMigrations"
Cohesion: 0.09
Nodes (43): newScratchStore(), TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(), TestUpsertBackupSchedulePostgresParity(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns(), sqliteSchemaColumns() (+35 more)

### Community 38 - "SnapshotEntity"
Cohesion: 0.12
Nodes (44): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools() (+36 more)

### Community 39 - "routerDeps"
Cohesion: 0.08
Nodes (38): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+30 more)

### Community 40 - "docker_test.go"
Cohesion: 0.06
Nodes (44): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+36 more)

### Community 41 - "ErrorWithDetails"
Cohesion: 0.09
Nodes (22): updateUserRequest, Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), configRequestField() (+14 more)

### Community 42 - "go_pkg_os"
Cohesion: 0.07
Nodes (28): TestCommitMessage(), writeExportState(), exportCursor, exportState, dockerSSHAddr, go_pkg_bufio, go_pkg_crypto_ed25519, go_pkg_encoding_csv (+20 more)

### Community 43 - "go_pkg_strings"
Cohesion: 0.06
Nodes (22): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), normalizeEnabledColumn() (+14 more)

### Community 44 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 45 - "DecodeKey"
Cohesion: 0.10
Nodes (24): ProviderConfig, testHandler, Handler, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey() (+16 more)

### Community 46 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 47 - "response.go"
Cohesion: 0.07
Nodes (28): decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+20 more)

### Community 48 - "User"
Cohesion: 0.09
Nodes (15): webAuthnUser, Handler, Handler, newOIDCUser(), randomOIDCToken(), validHostPort(), instanceAdminRoleFor(), OIDCClaims (+7 more)

### Community 49 - "chat/chat.go"
Cohesion: 0.09
Nodes (37): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections() (+29 more)

### Community 50 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 51 - "SuggestWithFallback"
Cohesion: 0.08
Nodes (22): claudeProvider, openAICompatibleProvider, StubProvider, SuggestResult, testProvider, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError() (+14 more)

### Community 52 - "Store"
Cohesion: 0.09
Nodes (23): Handler, NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler(), NewHandler() (+15 more)

### Community 53 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 54 - "compliance/engine.go"
Cohesion: 0.10
Nodes (33): catalog(), toRule(), validRecord(), contains(), equal(), Evaluate(), findAttribute(), Catalog (+25 more)

### Community 55 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 56 - "ConnectorEditPage.tsx"
Cohesion: 0.07
Nodes (26): RFC-3339, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_connectors_connectors_putconnectorsconnectorid (+18 more)

### Community 57 - "Connector"
Cohesion: 0.08
Nodes (12): init(), ConfigField, Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName() (+4 more)

### Community 58 - "adguardhome/tables.go"
Cohesion: 0.15
Nodes (30): statusInfo, unavailable(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+22 more)

### Community 59 - "ExportToFile"
Cohesion: 0.14
Nodes (32): Export(), ExportToFile(), Import(), ImportFromFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile() (+24 more)

### Community 60 - "traefik_test.go"
Cohesion: 0.10
Nodes (33): AllowLoopbackForTest(), TestNewHTTPClientRetriesAndBlocksRedirects(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection() (+25 more)

### Community 61 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 62 - "auth_test.go"
Cohesion: 0.09
Nodes (32): testApp, loginRefreshCookie(), seedLocalUser(), TestChangePasswordWrongCurrentPassword(), TestDeleteSessionNotOwner(), TestDeleteSessionSuccess(), TestElevateSuccess(), TestElevateWrongPassword() (+24 more)

### Community 63 - "router.go"
Cohesion: 0.10
Nodes (29): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_github_com_go_chi_chi_v5, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys (+21 more)

### Community 64 - "NewMalformedResponseError"
Cohesion: 0.09
Nodes (22): upstreamDependencies(), TestUpstreamDependenciesAreDedupedAndSorted(), ServiceDependency, NewMalformedResponseError(), WantsField(), TestRequestedFields(), TestWantsField(), TestBuildHostsTableAttributes() (+14 more)

### Community 65 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 66 - "Store"
Cohesion: 0.09
Nodes (8): placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange(), ChangeSummary

### Community 67 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 68 - "NewEngine"
Cohesion: 0.14
Nodes (26): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+18 more)

### Community 69 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 70 - "log/slog.Logger"
Cohesion: 0.10
Nodes (16): newLogger(), WithLogger(), formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store (+8 more)

### Community 71 - "nilToStr"
Cohesion: 0.10
Nodes (12): seedDelivery(), ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus (+4 more)

### Community 72 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 73 - "api/audit_test.go"
Cohesion: 0.10
Nodes (27): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+19 more)

### Community 74 - "httpx/retry_test.go"
Cohesion: 0.17
Nodes (23): isSafeMethod(), IsSafeMethod(), idempotent(), retryable(), RetryTransport(), sleep(), do(), fail() (+15 more)

### Community 75 - "AuthMiddleware"
Cohesion: 0.11
Nodes (22): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, AuthMiddleware(), extractBearerToken(), hashToken() (+14 more)

### Community 76 - "connector_permission.go"
Cohesion: 0.13
Nodes (15): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), treatAsSafeFromContext(), TreatAsSafeMethod(), getConnectorGrant(), ConnectorGrantDiff (+7 more)

### Community 77 - "Config"
Cohesion: 0.11
Nodes (21): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+13 more)

### Community 78 - "main"
Cohesion: 0.11
Nodes (16): Provider, main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+8 more)

### Community 79 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 80 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 81 - "HashToken"
Cohesion: 0.15
Nodes (15): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+7 more)

### Community 82 - "GetTypeSchema"
Cohesion: 0.12
Nodes (25): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+17 more)

### Community 83 - "connector/connector.go"
Cohesion: 0.08
Nodes (14): TimeoutError, GuardedDialer(), IsDangerousIP(), supportedLifecycleVerbs(), newWebhookClient(), AuthError, CredentialRefresher, MalformedResponseError (+6 more)

### Community 84 - "ConnectorRecord"
Cohesion: 0.13
Nodes (16): Store, ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), apiKeyConnectorFilter() (+8 more)

### Community 85 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 86 - "Hub"
Cohesion: 0.10
Nodes (11): loggablePath(), loggableQuery(), Sanitize(), TestSanitize(), Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, broadcastMsg (+3 more)

### Community 87 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 88 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 89 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 90 - "Manager"
Cohesion: 0.16
Nodes (8): cron.EntryID, Manager, LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store, Scheduler

### Community 91 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 92 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 93 - "newTestHandler"
Cohesion: 0.15
Nodes (18): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), templateRequest(), TestListPagination(), TestPreviewDoesNotPersist() (+10 more)

### Community 94 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 95 - "AppearancePage.tsx"
Cohesion: 0.16
Nodes (19): zustand, AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+11 more)

### Community 96 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 97 - "NewClient"
Cohesion: 0.12
Nodes (18): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), clientTimeout(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects() (+10 more)

### Community 98 - "Connector"
Cohesion: 0.17
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

### Community 99 - "RunbookRecord"
Cohesion: 0.21
Nodes (9): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep(), contains(), isUniqueViolation(), searchString() (+1 more)

### Community 100 - "time.Duration"
Cohesion: 0.14
Nodes (9): Cache, New(), Database, Server, time.Duration, PoolConfig, Cache[V], entry (+1 more)

### Community 101 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 102 - "backup/backup.go"
Cohesion: 0.24
Nodes (19): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, importBundle(), importConnectors() (+11 more)

### Community 103 - "config_test.go"
Cohesion: 0.15
Nodes (19): Load(), TestAccessTokenTTLDuration(), TestDocExportGitCommitModeValidation(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+11 more)

### Community 104 - "NewHTTPClient"
Cohesion: 0.13
Nodes (15): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+7 more)

### Community 105 - "home_assistant_test.go"
Cohesion: 0.22
Nodes (19): Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities(), TestFetchConfigIsRequestedOnce() (+11 more)

### Community 106 - "Handler"
Cohesion: 0.18
Nodes (9): validateConfigPushRequest(), stepAuditDetail(), validTargetType(), validVerb(), ValidateCompositeRef(), configPushRequest, Handler, runbookResponse (+1 more)

### Community 107 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 108 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 109 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 110 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 111 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 112 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 113 - "Store"
Cohesion: 0.18
Nodes (7): testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 114 - "logging_test.go"
Cohesion: 0.20
Nodes (14): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+6 more)

### Community 115 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 116 - "Connector"
Cohesion: 0.17
Nodes (6): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 117 - "time.Time"
Cohesion: 0.14
Nodes (9): digestDue(), TestDigestDue(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.Session, io.WriteCloser, net.Addr, time.Time (+1 more)

### Community 118 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 119 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 120 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 121 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 122 - "Register"
Cohesion: 0.23
Nodes (16): init(), init(), init(), init(), init(), init(), init(), init() (+8 more)

### Community 123 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 124 - "data.go"
Cohesion: 0.24
Nodes (15): ChangeEntry, ComplianceSection, ConnectorDrift, DefinitionSummary, DocChangeEntry, DocsSection, DriftSection, FindingSummary (+7 more)

### Community 126 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 127 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 128 - "OpenDB"
Cohesion: 0.15
Nodes (14): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), Store, newConcurrentQualityTestStore(), TestUpsertQualityFindingConcurrentDedup(), Store (+6 more)

### Community 129 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 130 - "api/attention_test.go"
Cohesion: 0.18
Nodes (13): TestAttentionAuthenticatedAccess(), TestAttentionDaysWindow(), TestAttentionEmptyList(), TestAttentionHidesUngrantedConnectors(), TestAttentionMergesAlertsAndFindings(), TestAttentionSeverityOrdering(), TestFindingResolveProducesAuditRecord(), testApp (+5 more)

### Community 131 - "Handler"
Cohesion: 0.27
Nodes (4): response(), writeRuleRejection(), Handler, RuleEvaluator

### Community 132 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 133 - "createUser"
Cohesion: 0.23
Nodes (15): ContextWithAPIKeyRestriction(), TestClampConnectorRole(), seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding() (+7 more)

### Community 134 - "VerifyBundleFile"
Cohesion: 0.25
Nodes (14): failVerification(), LatestBundle(), RunVerifyOnce(), ListVerifications(), RecordVerification(), seedOneDoc(), TestLatestBundleNoBundles(), TestListVerificationsNewestFirstAndLimit() (+6 more)

### Community 135 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 136 - "config/validate_test.go"
Cohesion: 0.18
Nodes (11): Config, mask(), redactDSN(), redactKVPassword(), Config, TestEveryKeyEnvOverridable(), TestRedactDSN(), TestRedacted() (+3 more)

### Community 137 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 138 - "Connector"
Cohesion: 0.16
Nodes (6): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), Connector

### Community 139 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 140 - "change_pattern_test.go"
Cohesion: 0.21
Nodes (13): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), T (+5 more)

### Community 141 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (11): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, resolveToken(), Snapshot (+3 more)

### Community 143 - "Config"
Cohesion: 0.19
Nodes (12): Config, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+4 more)

### Community 144 - "runRestore"
Cohesion: 0.21
Nodes (13): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+5 more)

### Community 145 - "config_cmd_test.go"
Cohesion: 0.23
Nodes (10): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+2 more)

### Community 146 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 147 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 148 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 149 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 150 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 151 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 152 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 154 - "lifecycleDeps"
Cohesion: 0.20
Nodes (8): newLifecycleManager(), context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, sync/atomic.Bool, lifecycleDeps, lifecycleManager, ReadyState

### Community 155 - "registry.go"
Cohesion: 0.24
Nodes (9): IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError, Factory, SchemaField (+1 more)

### Community 156 - "server.go"
Cohesion: 0.27
Nodes (5): go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, changeSummary, connectorSummary, findingSummary

### Community 157 - "ReportData"
Cohesion: 0.40
Nodes (5): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), Generator, ReportData

### Community 159 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 160 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 161 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 162 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 163 - "net/http.Handler"
Cohesion: 0.22
Nodes (8): ConnectorRoleChecker, RequireConnectorRole(), RequireInstanceAdmin(), contextWithInstanceAdmin(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError(), TestRequireInstanceAdmin(), net/http.Handler

### Community 164 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 165 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 166 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 167 - "ratelimit.go"
Cohesion: 0.27
Nodes (7): TestRateLimit(), RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 168 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 169 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 170 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 171 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 172 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 173 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 174 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 175 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 176 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 178 - "Engine"
Cohesion: 0.29
Nodes (3): Engine, sync.Map, DocRegenerator

### Community 179 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 180 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 181 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 182 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 183 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 184 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 185 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 186 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 188 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 190 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 191 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 192 - "browser.ts"
Cohesion: 0.40
Nodes (4): bootstrap(), worker, enableMocks(), handlers

### Community 193 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 194 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 195 - ".operatorWithPermission"
Cohesion: 0.50
Nodes (3): testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault()

### Community 196 - "mintAPIKey"
Cohesion: 0.50
Nodes (4): testApp, mintAPIKey(), TestMCPConnectorRestrictedKey(), TestMCPEndToEnd()

### Community 197 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 201 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 202 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **568 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+563 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1283 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `testing.T`, `Handler`, `Handler`, `createUser`, `VerifyBundleFile`, `Errorf`, `Handler`, `Config`, `runRestore`, `NewStore`, `go_pkg_context`, `ServiceSnapshot`, `gitFixture`, `NewUser`, `lifecycleDeps`, `rowScanner`, `NewEngine`, `ReportData`, `dispatcher_test.go`, `net/http.ResponseWriter`, `.call`, `RunMigrations`, `ErrorWithDetails`, `DecodeKey`, `retention/retention_test.go`, `chat/chat.go`, `NewChecker`, `Engine`, `Dispatcher`, `ExportToFile`, `engine_maintenance_test.go`, `NewEngine`, `nilToStr`, `rewritePlaceholders`, `Checker`, `Manager`, `newTestHandler`, `backup/backup.go`, `Handler`, `diagnostics/diagnostics.go`, `time.Time`, `testApp`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `gitFixture` to `testing.T`, `context.Context`, `log/slog.Logger`, `go_pkg_os`, `Store`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Why does `UserIDFromContext()` connect `Errorf` to `net/http.ResponseWriter`, `net/http.Handler`, `context.Context`, `Handler`, `routerDeps`, `ErrorWithDetails`, `Handler`, `AuthMiddleware`, `Handler`, `response.go`, `User`, `HashToken`, `chat/chat.go`, `Store`, `go_pkg_github_com_wiselabz_wiselabz_internal_connector`, `net/http.Request`, `NewUser`, `rowScanner`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _568 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02018795683954055 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.025707547169811322 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.025804261098378745 - nodes in this community are weakly interconnected._