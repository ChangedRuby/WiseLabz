# Graph Report - wiselabz-wc-retry  (2026-09-23)

## Corpus Check
- 769 files · ~467,111 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 10, .toml 2, .example 1)

## Summary
- 5737 nodes · 18009 edges · 199 communities (184 shown, 15 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1429 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3c9c6805`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- SystemPage.tsx
- react
- go_pkg_context
- icons.tsx
- context.Context
- NewChecker
- ServiceDetailPage.tsx
- @tanstack/react-query
- go_pkg_testing
- connector/connector.go
- go_pkg_github_com_wiselabz_wiselabz_internal_store
- Errorf
- net/http.Request
- newTestHandler
- DashboardPage.tsx
- dispatcher_test.go
- App.tsx
- cn
- ServiceSnapshot
- UsersPage.tsx
- NewEngine
- package.json
- Runner
- export_test.go
- DecodeKey
- Store
- IsSecureRequest
- SnapshotEntity
- fixtures.ts
- Get
- WiseLabz — Architecture & Technical Decisions
- RunMigrations
- docker_test.go
- NewStore
- share_links_test.go
- paginatedQuery
- home_assistant/tables.go
- time.Time
- ErrorWithDetails
- dependencies
- RulesPage.tsx
- NewRegistry
- portainer/tables.go
- NewService
- adguardhome/tables.go
- git.go
- home_assistant_test.go
- lists.go
- response.go
- go_pkg_os
- traefik/tables.go
- GuardedDialer
- router.go
- Connector
- Connector
- truenas_test.go
- unifi/tables.go
- Dispatcher
- Configuration & Documentation Backup (Export/Import)
- SuggestRequest
- mountAPIRoutes
- settings.mock.ts
- DecodeJSON
- api/audit_test.go
- rewritePlaceholders
- newTestHandler
- nilToStr
- GetTypeSchema
- Register
- Connector
- src/theme.ts
- TemplatesPage.tsx
- handlers.ts
- Store
- unifi_test.go
- NotificationCenter.tsx
- timeline.ts
- VerifyBundleFile
- newTestHandler
- AuthedUser
- config_test.go
- WiseLabz — Design Contract
- AppearancePage.tsx
- devDependencies
- Compare
- logging.go
- ExportToFile
- NewMalformedResponseError
- net/http.Client
- portainer_test.go
- Hub
- adguardhome_test.go
- ConnectorRecord
- rowScanner
- middleware.go
- Service
- chat/chat.go
- Connector
- Config
- New
- Connector
- ws/ws_test.go
- log/slog.Logger
- Handler
- Config
- diagnostics/diagnostics.go
- keyset_test.go
- ws.ts
- compilerOptions
- Handler
- traefik_test.go
- docdiffmodel.ts
- main
- all.go
- Connector
- compilerOptions
- config_cmd_test.go
- changes/handlers_test.go
- Handler
- vectorCache
- DocRecord
- .batchDelete
- templates.fixtures.ts
- pagination_contract_test.go
- time.Duration
- changes_test.go
- connectors_health_test.go
- registry.go
- NewClient
- Contributing to WiseLabz
- Engine
- .call
- templatefuncs.go
- store/backup_test.go
- Store
- Decision
- WiseLabz Connector Guide
- main.tsx
- scripts
- channels.go
- Store
- Decision
- Product
- .call
- Store
- connectors_maintenance_test.go
- Changelog
- mockServiceWorker.js
- ComplianceRuleRecord
- ratelimit.go
- retention/retention_test.go
- RunbookRecord
- sync.Mutex
- release-please-config.json
- decodePaginated
- BackupSchedule
- ShareLink
- transform.go
- Cache
- Step by step
- WiseLabz
- TestComplianceRuleValidation
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- walkCursorPages
- APIKeyClaims
- dialSSHStdio
- Store
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- RequireConnectorRole
- testHandler
- routerOperations
- ref_test.go
- seedScopeFixture
- Enforcement Guidelines
- compose-smoke.sh
- ClassifyHealth
- RetentionSettings
- timeoutError
- MISSING — deferred & future frontend features
- Saved Views
- truenas/attributes_test.go
- dockerSSHAddr
- WiseLabz — v2 Backlog
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

## Communities (199 total, 15 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (162): templateBody, TestAPIKeyCreateRejectsInvalidExpiryAndEmptyName(), TestAPIKeyRoutesEndToEnd(), TestAttentionAuthenticatedAccess(), TestAttentionDaysWindow(), TestAttentionEmptyList(), TestAttentionHidesUngrantedConnectors(), TestAttentionMergesAlertsAndFindings() (+154 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.03
Nodes (132): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+124 more)

### Community 2 - "testing.T"
Cohesion: 0.03
Nodes (124): TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider(), TestOIDCCallbackRejectsMissingFlowCookie(), TestOIDCCallbackRejectsStateMismatchedWithCookie(), TestOIDCCallbackRejectsUnknownProvider(), TestWriteConfigRejection() (+116 more)

### Community 3 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (99): web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_me_me_deletemesessionssessionid, web_src_api_generated_me_me_getgetmequerykey, web_src_api_generated_me_me_getgetmesessionsquerykey, web_src_api_generated_me_me_getmesessions (+91 more)

### Community 4 - "react"
Cohesion: 0.05
Nodes (94): RFC-3339, react, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze, web_src_api_generated_alerts_alerts_postalertsbulksnooze (+86 more)

### Community 5 - "go_pkg_context"
Cohesion: 0.07
Nodes (26): contains(), searchString(), TestNormalizeFirewallRulesRewritesEnabledColumn(), go_pkg_bytes, go_pkg_context, go_pkg_crypto_hmac, go_pkg_crypto_sha256, go_pkg_crypto_tls (+18 more)

### Community 6 - "icons.tsx"
Cohesion: 0.03
Nodes (96): match-sorter, motion, @radix-ui/react-popover, zustand, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_postauthelevate (+88 more)

### Community 7 - "context.Context"
Cohesion: 0.04
Nodes (26): fakeStatusChecker, sanitizeSessions(), Handler, Connector, Connector, Connector, existingIDs(), Store (+18 more)

### Community 8 - "NewChecker"
Cohesion: 0.06
Nodes (71): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+63 more)

### Community 9 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (85): ADR-0001, ADR-0003, Frontend, 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice` (+77 more)

### Community 10 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (53): Client dispatch model, msw, react-router-dom, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_changes_changes (+45 more)

### Community 11 - "go_pkg_testing"
Cohesion: 0.06
Nodes (28): dashboardLayout, TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestRegisterClaudeDefaults(), Handler, newTestHandler() (+20 more)

### Community 12 - "connector/connector.go"
Cohesion: 0.04
Nodes (35): Connector, isTimeout(), NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL(), tryParseEntities() (+27 more)

### Community 13 - "go_pkg_github_com_wiselabz_wiselabz_internal_store"
Cohesion: 0.08
Nodes (31): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections(), TemplateVersionSection, bulkResolveItemResult (+23 more)

### Community 14 - "Errorf"
Cohesion: 0.07
Nodes (28): routerDeps, Handler, newToken(), setRefreshCookie(), Handler, diffToSpec(), NewHandler(), Handler (+20 more)

### Community 15 - "net/http.Request"
Cohesion: 0.07
Nodes (22): Handler, Handler, Handler, decodeBulkRequest(), Handler, Handler, Handler, Handler (+14 more)

### Community 16 - "newTestHandler"
Cohesion: 0.07
Nodes (64): AssertMatchesSpec(), loadSpec(), specPath(), actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews() (+56 more)

### Community 17 - "DashboardPage.tsx"
Cohesion: 0.05
Nodes (56): 4. `change.detected`, 5. `alert.created`, 8. `doc.generated`, web_src_api_generated_alerts_alerts_usegetalerts, web_src_api_generated_dashboard_dashboard_getdashboardlayout, web_src_api_generated_dashboard_dashboard_getdashboardlayoutadmindefault, web_src_api_generated_dashboard_dashboard_getgetdashboardlayoutadmindefaultquerykey, web_src_api_generated_dashboard_dashboard_postdashboardlayoutreset (+48 more)

### Community 18 - "dispatcher_test.go"
Cohesion: 0.12
Nodes (61): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), testLogger(), expireAlertsOnce() (+53 more)

### Community 19 - "App.tsx"
Cohesion: 0.05
Nodes (49): setAccessToken(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh, web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_model_index_authsession, web_src_api_model_index_oidccallbackrequest (+41 more)

### Community 20 - "cn"
Cohesion: 0.05
Nodes (48): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateidversions (+40 more)

### Community 21 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (18): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, runTransformers(), TestRunTransformersAppliesInOrderAndStopsOnError() (+10 more)

### Community 22 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (49): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), RefreshFn, setRefreshHandler() (+41 more)

### Community 23 - "NewEngine"
Cohesion: 0.09
Nodes (44): NewHandler(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine (+36 more)

### Community 24 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 25 - "Runner"
Cohesion: 0.08
Nodes (30): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner, New() (+22 more)

### Community 26 - "export_test.go"
Cohesion: 0.08
Nodes (36): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+28 more)

### Community 27 - "DecodeKey"
Cohesion: 0.08
Nodes (28): ProviderConfig, Handler, Handler, primaryProviderConfig(), Handler, Config, mask(), redactDSN() (+20 more)

### Community 28 - "Store"
Cohesion: 0.07
Nodes (13): Store, placeholders(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange(), changePatternID() (+5 more)

### Community 29 - "IsSecureRequest"
Cohesion: 0.08
Nodes (29): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), clearOIDCFlowCookie(), emailDomainAllowed(), Handler, newOIDCUser(), oidcFlowCookieName(), oidcRoleForGroups() (+21 more)

### Community 30 - "SnapshotEntity"
Cohesion: 0.12
Nodes (44): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools() (+36 more)

### Community 31 - "fixtures.ts"
Cohesion: 0.06
Nodes (41): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+33 more)

### Community 32 - "Get"
Cohesion: 0.07
Nodes (28): Handler, isWritableField(), validateConfigPushRequest(), capitalize(), Handler, WriteElevationError(), ConfigPusher, ValidateCompositeRef() (+20 more)

### Community 33 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 34 - "RunMigrations"
Cohesion: 0.10
Nodes (38): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+30 more)

### Community 35 - "docker_test.go"
Cohesion: 0.06
Nodes (43): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+35 more)

### Community 36 - "NewStore"
Cohesion: 0.09
Nodes (41): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+33 more)

### Community 37 - "share_links_test.go"
Cohesion: 0.17
Nodes (41): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld() (+33 more)

### Community 38 - "paginatedQuery"
Cohesion: 0.09
Nodes (19): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), changeFilterClause(), countRows(), T (+11 more)

### Community 39 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 40 - "time.Time"
Cohesion: 0.05
Nodes (17): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), Store, registryTestRefresher, sshStdioConn, golang.org/x/crypto/ssh.Client (+9 more)

### Community 41 - "ErrorWithDetails"
Cohesion: 0.09
Nodes (22): updateUserRequest, Handler, sanitizeUser(), writeUserWriteError(), Handler, mustHashDummyPassword(), configRequestField(), parseScheduleUpdates() (+14 more)

### Community 42 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 43 - "RulesPage.tsx"
Cohesion: 0.06
Nodes (37): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+29 more)

### Community 44 - "NewRegistry"
Cohesion: 0.12
Nodes (26): Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail() (+18 more)

### Community 45 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 46 - "NewService"
Cohesion: 0.11
Nodes (32): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService() (+24 more)

### Community 47 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 48 - "git.go"
Cohesion: 0.08
Nodes (28): gitAuth(), Exporter, installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), TestGitAuthSSH(), writeTestKey(), keys() (+20 more)

### Community 49 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (35): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+27 more)

### Community 50 - "lists.go"
Cohesion: 0.12
Nodes (34): buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames(), parseAdlistsV6() (+26 more)

### Community 51 - "response.go"
Cohesion: 0.08
Nodes (26): Handler, Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes() (+18 more)

### Community 52 - "go_pkg_os"
Cohesion: 0.09
Nodes (19): TestSnapshotAttributesRoundTripPostgres(), TestSnapshotAttributesRoundTripSQLite(), testSnapshotWithAttributes(), go_pkg_bufio, go_pkg_flag, go_pkg_github_com_jackc_pgx_v5_stdlib, go_pkg_github_com_robfig_cron_v3, go_pkg_github_com_wiselabz_wiselabz_internal_backup (+11 more)

### Community 53 - "traefik/tables.go"
Cohesion: 0.14
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 54 - "GuardedDialer"
Cohesion: 0.08
Nodes (29): init(), newConnector(), Connector, GuardedDialer(), IsDangerousIP(), init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal() (+21 more)

### Community 55 - "router.go"
Cohesion: 0.10
Nodes (26): go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+18 more)

### Community 56 - "Connector"
Cohesion: 0.09
Nodes (12): init(), ConfigField, Connector, TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), isTimeout(), primaryGatewayName() (+4 more)

### Community 57 - "Connector"
Cohesion: 0.14
Nodes (12): SnapshotSection, TestBuildHostsTableAttributes(), TestBuildHostsTableV5(), buildHostsTable(), Connector, isTimeout(), parseHosts(), TestBuildHostsTableMalformedCases() (+4 more)

### Community 58 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 59 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 60 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 61 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 62 - "SuggestRequest"
Cohesion: 0.09
Nodes (14): claudeProvider, ollamaEmbedder, openAICompatibleProvider, openAIEmbedder, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet() (+6 more)

### Community 63 - "mountAPIRoutes"
Cohesion: 0.11
Nodes (25): chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router, mountChatRoutes() (+17 more)

### Community 64 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 65 - "DecodeJSON"
Cohesion: 0.11
Nodes (8): Handler, Handler, oidcProviderJSON(), boolToInt(), DecodeJSON(), T, SinceFromDays(), Handler

### Community 66 - "api/audit_test.go"
Cohesion: 0.10
Nodes (28): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+20 more)

### Community 67 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 68 - "newTestHandler"
Cohesion: 0.11
Nodes (22): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+14 more)

### Community 69 - "nilToStr"
Cohesion: 0.10
Nodes (10): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+2 more)

### Community 70 - "GetTypeSchema"
Cohesion: 0.11
Nodes (26): catalog(), TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion() (+18 more)

### Community 71 - "Register"
Cohesion: 0.20
Nodes (25): init(), RequestedFields(), Register(), init(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect() (+17 more)

### Community 72 - "Connector"
Cohesion: 0.12
Nodes (14): TestRateLimit(), ServiceDependency, environmentDependencies(), poolDependencies(), apiMessage(), controllerName(), countByKind(), isTimeout() (+6 more)

### Community 73 - "src/theme.ts"
Cohesion: 0.13
Nodes (26): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ThemeControls(), ColorMode, commit(), load(), Persisted (+18 more)

### Community 74 - "TemplatesPage.tsx"
Cohesion: 0.09
Nodes (20): i18next, web_src_api_generated_templates_templates_deletetemplatestemplateid, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_posttemplates, web_src_api_generated_templates_templates_usegettemplates, web_src_api_model_index_template, TemplatesPage, Command (+12 more)

### Community 75 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 76 - "Store"
Cohesion: 0.22
Nodes (22): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+14 more)

### Community 77 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 78 - "NotificationCenter.tsx"
Cohesion: 0.10
Nodes (20): Frontend shell & theme (decided 2026-06), react-error-boundary, sonner, web_src_api_generated_notifications_notifications, web_src_api_generated_notifications_notifications_getgetnotificationsquerykey, web_src_api_generated_notifications_notifications_postnotificationsnotificationidread, web_src_api_generated_notifications_notifications_postnotificationsreadall, web_src_api_generated_notifications_notifications_usegetnotifications (+12 more)

### Community 79 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 80 - "VerifyBundleFile"
Cohesion: 0.16
Nodes (23): AppVersion(), BuildManifest(), BundleCounts(), ChecksumBytes(), ReadManifest(), WriteManifest(), failVerification(), LatestBundle() (+15 more)

### Community 81 - "newTestHandler"
Cohesion: 0.18
Nodes (20): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+12 more)

### Community 82 - "AuthedUser"
Cohesion: 0.14
Nodes (24): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+16 more)

### Community 83 - "config_test.go"
Cohesion: 0.12
Nodes (22): runHealthcheck(), Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+14 more)

### Community 84 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 85 - "AppearancePage.tsx"
Cohesion: 0.15
Nodes (20): MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+12 more)

### Community 86 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 87 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 88 - "logging.go"
Cohesion: 0.15
Nodes (17): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestLoggablePathMasksShareTokenUnderV1() (+9 more)

### Community 89 - "ExportToFile"
Cohesion: 0.19
Nodes (21): Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+13 more)

### Community 90 - "NewMalformedResponseError"
Cohesion: 0.14
Nodes (10): NewMalformedResponseError(), WantsField(), TestRequestedFields(), TestWantsField(), isTimeout(), putMetadata(), unavailable(), MalformedResponseError (+2 more)

### Community 91 - "net/http.Client"
Cohesion: 0.11
Nodes (5): Connector, Connector, newWebhookClient(), Connector, net/http.Client

### Community 92 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 93 - "Hub"
Cohesion: 0.13
Nodes (7): Sanitize(), TestSanitize(), Hub, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator

### Community 94 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 95 - "ConnectorRecord"
Cohesion: 0.16
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 96 - "rowScanner"
Cohesion: 0.15
Nodes (10): Store, scanMaintenanceWindow(), NotificationRecord, Store, scanNotification(), scanQualityFinding(), scanTemplate(), scanUser() (+2 more)

### Community 97 - "middleware.go"
Cohesion: 0.15
Nodes (13): APIKeyChecker, AuditRecorder, contextKey, elevationError, PermissionChecker, UserStatusChecker, elevationFailureReason(), extractBearerToken() (+5 more)

### Community 98 - "Service"
Cohesion: 0.20
Nodes (10): Claims, ElevationClaims, ElevationToken, TokenPair, Service, hasAudience(), newTokenID(), go_pkg_github_com_golang_jwt_jwt_v5 (+2 more)

### Community 99 - "chat/chat.go"
Cohesion: 0.16
Nodes (17): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings() (+9 more)

### Community 100 - "Connector"
Cohesion: 0.14
Nodes (8): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), isTimeout(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 101 - "Config"
Cohesion: 0.13
Nodes (17): Config, CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+9 more)

### Community 102 - "New"
Cohesion: 0.14
Nodes (19): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+11 more)

### Community 103 - "Connector"
Cohesion: 0.16
Nodes (7): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), isTimeout(), Connector

### Community 104 - "ws/ws_test.go"
Cohesion: 0.19
Nodes (18): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+10 more)

### Community 105 - "log/slog.Logger"
Cohesion: 0.16
Nodes (11): newLogger(), WithLogger(), Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+3 more)

### Community 106 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 107 - "Config"
Cohesion: 0.18
Nodes (14): Config, LogSettings, IsSSHRemote(), AISettings, BackupSettings, DocExportGitSettings, DocExportSettings, EncryptionSettings (+6 more)

### Community 108 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 109 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 110 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 111 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 112 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 113 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 114 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 115 - "main"
Cohesion: 0.21
Nodes (11): main(), splitOrigins(), RegisterClaude(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry, NewEmbedRegistry() (+3 more)

### Community 116 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 117 - "Connector"
Cohesion: 0.16
Nodes (7): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), isTimeout(), Connector

### Community 118 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 119 - "config_cmd_test.go"
Cohesion: 0.19
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 120 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 121 - "Handler"
Cohesion: 0.20
Nodes (4): applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), updateConnectorRequest

### Community 122 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 123 - "DocRecord"
Cohesion: 0.20
Nodes (6): docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 125 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (12): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate(), resolveToken() (+4 more)

### Community 126 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 127 - "time.Duration"
Cohesion: 0.19
Nodes (5): AuthSettings, Database, Server, time.Duration, OIDCProvider

### Community 128 - "changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 129 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 130 - "registry.go"
Cohesion: 0.19
Nodes (10): IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError, Factory, SchemaField (+2 more)

### Community 131 - "NewClient"
Cohesion: 0.18
Nodes (12): clientTimeout(), NewClient(), NewTransport(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects(), TestNewClientTimeout(), TestNewClientUsesCustomDialContext(), TestNewTransportTLS() (+4 more)

### Community 132 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 133 - "Engine"
Cohesion: 0.18
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 134 - ".call"
Cohesion: 0.33
Nodes (8): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), spaHandler(), fixture, net/http.HandlerFunc

### Community 135 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 136 - "store/backup_test.go"
Cohesion: 0.32
Nodes (11): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+3 more)

### Community 137 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 138 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 139 - "WiseLabz Connector Guide"
Cohesion: 0.17
Nodes (12): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+4 more)

### Community 140 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 141 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 142 - "channels.go"
Cohesion: 0.25
Nodes (9): buildEmailMessage(), sendSMTPChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendSMTPChannel_MissingConfig(), TestSplitRecipients(), redactURLError(), go_pkg_net_smtp (+1 more)

### Community 143 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 144 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 145 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 146 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 147 - "Store"
Cohesion: 0.27
Nodes (3): sanitize(), APIKey, Store

### Community 148 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 149 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 150 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 151 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 152 - "ratelimit.go"
Cohesion: 0.31
Nodes (6): RateLimit(), go_pkg_golang_org_x_time_rate, golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 153 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 154 - "RunbookRecord"
Cohesion: 0.44
Nodes (3): RunbookRecord, Store, scanRunbook()

### Community 155 - "sync.Mutex"
Cohesion: 0.22
Nodes (4): sync.Mutex, fakeDocRegenerator, fakeNotifier, fakeQualityChecker

### Community 156 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 157 - "decodePaginated"
Cohesion: 0.36
Nodes (8): decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty(), TestListDeliveriesNoFilter(), TestListDeliveriesStatusFilter(), PaginatedResponse

### Community 158 - "BackupSchedule"
Cohesion: 0.32
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 160 - "transform.go"
Cohesion: 0.32
Nodes (5): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), Transformer

### Community 161 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 162 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 163 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 164 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 165 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 166 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 167 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 168 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 169 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 170 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 171 - "APIKeyClaims"
Cohesion: 0.40
Nodes (3): testAPIKeyChecker, APIKeyClaims, validAPIKey()

### Community 172 - "dialSSHStdio"
Cohesion: 0.33
Nodes (5): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), golang.org/x/crypto/ssh.ClientConfig, io.Closer

### Community 174 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 175 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 176 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 177 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 178 - "testHandler"
Cohesion: 0.40
Nodes (3): testHandler, Handler, instanceAdminRoleFor()

### Community 179 - "routerOperations"
Cohesion: 0.50
Nodes (5): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes

### Community 180 - "ref_test.go"
Cohesion: 0.40
Nodes (4): TestValidateCompositeRef(), TestValidateRefSegment(), TestValidateUnixSocketPath(), ValidateUnixSocketPath()

### Community 181 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 182 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 183 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 184 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 188 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 189 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **540 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+535 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1181 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `UserIDFromContext()` connect `Errorf` to `DecodeJSON`, `middleware.go`, `paginatedQuery`, `context.Context`, `ErrorWithDetails`, `NewService`, `net/http.Request`, `RequireConnectorRole`, `response.go`, `Handler`, `mountAPIRoutes`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `Store` connect `Store` to `Engine`, `go_pkg_context`, `.call`, `NewChecker`, `store/backup_test.go`, `Errorf`, `net/http.Request`, `.call`, `dispatcher_test.go`, `NewEngine`, `retention/retention_test.go`, `export_test.go`, `sync.Mutex`, `Store`, `decodePaginated`, `RunMigrations`, `NewStore`, `share_links_test.go`, `paginatedQuery`, `time.Time`, `ErrorWithDetails`, `NewRegistry`, `engine_maintenance_test.go`, `testHandler`, `response.go`, `Dispatcher`, `DecodeJSON`, `rewritePlaceholders`, `newTestHandler`, `Register`, `VerifyBundleFile`, `AuthedUser`, `ExportToFile`, `chat/chat.go`, `Config`, `New`, `Handler`, `diagnostics/diagnostics.go`, `changes/handlers_test.go`, `Handler`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `WebSocketProvider()` connect `@tanstack/react-query` to `App.tsx`, `ServiceDetailPage.tsx`, `WiseLabz Connector Guide`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _540 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021883116883116883 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.025347125151430436 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0260235947258848 - nodes in this community are weakly interconnected._