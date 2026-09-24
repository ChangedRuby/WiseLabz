# Graph Report - wiselabz-we  (2026-09-24)

## Corpus Check
- 788 files · ~477,525 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 5868 nodes · 18429 edges · 214 communities (198 shown, 16 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1452 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a118412b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- cn
- testing.T
- react
- context.Context
- Button.tsx
- go_pkg_net_http
- DashboardPage.tsx
- @tanstack/react-query
- go_pkg_encoding_json
- ServiceDetailPage.tsx
- go_pkg_context
- App.tsx
- ServiceSnapshot
- go_pkg_testing
- icons.tsx
- RunMigrations
- Service
- Errorf
- net/http.Client
- rowScanner
- dispatcher_test.go
- MarshalConnectorConfig
- ErrorWithDetails
- Store
- newDocTestStore
- portainer/tables.go
- package.json
- net/http.Request
- fixtures.ts
- SuggestRequest
- WiseLabz — Architecture & Technical Decisions
- truenas/tables.go
- Connector
- lists.go
- home_assistant/tables.go
- DecodeJSON
- IsSecureRequest
- newTestHandler
- docker_test.go
- routerDeps
- traefik/tables_test.go
- nilToStr
- UsersPage.tsx
- NewChecker
- dependencies
- response.go
- NewStore
- middleware.go
- router.go
- share_links_test.go
- ShareLinkPage.tsx
- adguardhome/tables.go
- NewRegistry
- Register
- ThemeControls.tsx
- auth_test.go
- Dispatcher
- compliance/engine.go
- Config
- truenas_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- WebSocketProvider.tsx
- Connector
- config_test.go
- settings.mock.ts
- HashToken
- rewritePlaceholders
- git.go
- home_assistant_test.go
- GetTypeSchema
- handlers.ts
- ExportToFile
- unifi_test.go
- timeline.ts
- New
- Store
- time.Time
- createTestConnector
- src/theme.ts
- .call
- log/slog.Logger
- compliance/handlers.go
- Checker
- Manager
- WiseLabz — Design Contract
- Connector
- devDependencies
- newTestHandler
- Compare
- logging.go
- portainer_test.go
- export_test.go
- doc_test.go
- chat/chat.go
- adguardhome_test.go
- scheduler/health_test.go
- ConnectorRecord
- main
- APIKeyClaims
- diagnostics/diagnostics.go
- Runner
- gitAuth
- Hub
- ws/ws_test.go
- sshStdioConn
- ws.ts
- compilerOptions
- docs/handlers_test.go
- sync.Mutex
- Handler
- SnapshotEntity
- traefik_test.go
- NewEngine
- docdiffmodel.ts
- chat/handlers_test.go
- templates_test.go
- reports/handlers.go
- all.go
- .batchDelete
- templates.fixtures.ts
- compilerOptions
- notifications/handlers_test.go
- changes/handlers_test.go
- middleware_test.go
- vectorCache
- Engine
- templatefuncs.go
- NewMalformedResponseError
- system/handlers_test.go
- testApp
- pagination_contract_test.go
- newTestHandler
- Connector
- gitTarget
- render_test.go
- New
- connectors_health_test.go
- gitFixture
- NewClient
- maintenance_test.go
- Contributing to WiseLabz
- Engine
- Connector
- time.Duration
- ReportData
- quality_test.go
- transform.go
- Decision
- WiseLabz Connector Guide
- main.tsx
- scripts
- Handler
- fetch_test.go
- Store
- connector_permission_test.go
- keyset_test.go
- Decision
- Product
- cursor_pagination_test.go
- AuthMiddleware
- changes_test.go
- NewService
- .call
- Changelog
- mockServiceWorker.js
- connectors_hardening_test.go
- connector/connector.go
- retention/retention_test.go
- release-please-config.json
- NotificationRecord
- ComputeWindow
- notification_delivery_test.go
- isUniqueViolation
- openapi_contract_test.go
- Cache
- Step by step
- WiseLabz
- newRouterDeps
- .applyChannelSecrets
- responseWriter
- change_pattern_test.go
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- BackupSchedule
- Mermaid.tsx
- Security Policy
- .resolveConfigPusher
- golden_snapshot_test.go
- Enforcement Guidelines
- @vitejs/plugin-react
- compose-smoke.sh
- docker/snapshot.go
- timeoutError
- MISSING — deferred & future frontend features
- Store
- WiseLabz — v2 Backlog
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- vite-env.d.ts
- github.com/WiseLabz/wiselabz
- RequireConnectorRole
- .RunDigestSweep
- seedScopeFixture
- snapshot_attributes_test.go
- truenas/attributes_test.go

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 212 edges
2. `Errorf()` - 164 edges
3. `Store` - 133 edges
4. `newDocTestStore()` - 129 edges
5. `SnapshotEntity` - 74 edges
6. `react` - 73 edges
7. `UserIDFromContext()` - 69 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
10. `@tanstack/react-query` - 58 edges

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

## Communities (214 total, 16 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.03
Nodes (136): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+128 more)

### Community 1 - "cn"
Cohesion: 0.02
Nodes (133): web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_auth_auth_usegetauthproviders, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey (+125 more)

### Community 2 - "testing.T"
Cohesion: 0.03
Nodes (124): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestWriteConfigRejection() (+116 more)

### Community 3 - "react"
Cohesion: 0.03
Nodes (102): react, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid (+94 more)

### Community 4 - "context.Context"
Cohesion: 0.03
Nodes (36): fakeStatusChecker, sanitize(), sanitizeSessions(), Connector, Connector, APIKey, Store, existingIDs() (+28 more)

### Community 5 - "Button.tsx"
Cohesion: 0.03
Nodes (94): Endpoints, Frontend, Saved Views, Scope, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover (+86 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (50): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, updateUserRequest, newToken(), changePromptData(), diffToSpec(), stripPromptTags() (+42 more)

### Community 7 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (78): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 5. `alert.created` (+70 more)

### Community 8 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (53): msw, @tanstack/react-query, @testing-library/react, vitest, zustand, web_src_api_generated_connectors_connectors, web_src_api_model_index_attentionpage, web_src_api_model_index_runbookpage (+45 more)

### Community 9 - "go_pkg_encoding_json"
Cohesion: 0.07
Nodes (38): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable() (+30 more)

### Community 10 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (64): ADR-0001, ADR-0003, RFC-3339, 1. `service.status`, 4. `change.detected`, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush (+56 more)

### Community 11 - "go_pkg_context"
Cohesion: 0.08
Nodes (21): StatusError, ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), contains(), searchString(), bulkRequest, go_pkg_context (+13 more)

### Community 12 - "App.tsx"
Cohesion: 0.05
Nodes (58): Frontend shell & theme (decided 2026-06), react-error-boundary, react-i18next, react-router-dom, sonner, setAccessToken(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_postauthlogin (+50 more)

### Community 13 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (20): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, Sanitize(), TestSanitize(), changePatternID() (+12 more)

### Community 14 - "go_pkg_testing"
Cohesion: 0.05
Nodes (30): TestDiagnosticsIncludesHealthVersionsAndFailures(), TestDiagnosticsRedactsSecrets(), TestDiagnosticsRoleBoundary(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), Schema() (+22 more)

### Community 15 - "icons.tsx"
Cohesion: 0.05
Nodes (56): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_runbooks_runbooks, web_src_api_generated_runbooks_runbooks_usegetrunbooks (+48 more)

### Community 16 - "RunMigrations"
Cohesion: 0.11
Nodes (37): main(), OpenDB(), GetMigrationStatus(), newMigrator(), collectColumns(), migrationFiles(), postgresSchemaColumns(), sqliteSchemaColumns() (+29 more)

### Community 17 - "Service"
Cohesion: 0.20
Nodes (10): Claims, ElevationClaims, ElevationToken, TokenPair, Service, hasAudience(), newTokenID(), go_pkg_github_com_golang_jwt_jwt_v5 (+2 more)

### Community 18 - "Errorf"
Cohesion: 0.08
Nodes (14): NewHandler(), Handler, Handler, Handler, Handler, Handler, Handler, UserIDFromContext() (+6 more)

### Community 19 - "net/http.Client"
Cohesion: 0.03
Nodes (42): Connector, ollamaEmbedder, openAIEmbedder, isTimeout(), NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), setHeaders() (+34 more)

### Community 20 - "rowScanner"
Cohesion: 0.06
Nodes (27): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), docSearchWhere(), escapeLike(), DocRecord (+19 more)

### Community 21 - "dispatcher_test.go"
Cohesion: 0.17
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 22 - "MarshalConnectorConfig"
Cohesion: 0.07
Nodes (38): ProviderConfig, Handler, primaryProviderConfig(), Config, mask(), DecodeKey(), Decrypt(), DeriveKey() (+30 more)

### Community 23 - "ErrorWithDetails"
Cohesion: 0.09
Nodes (20): Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), parseScheduleUpdates(), validateRotationFields() (+12 more)

### Community 24 - "Store"
Cohesion: 0.08
Nodes (10): Store, placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+2 more)

### Community 25 - "newDocTestStore"
Cohesion: 0.06
Nodes (47): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+39 more)

### Community 26 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 27 - "package.json"
Cohesion: 0.04
Nodes (45): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+37 more)

### Community 28 - "net/http.Request"
Cohesion: 0.06
Nodes (30): validateConfigPushRequest(), Handler, applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), writeConfigRejection(), Handler, Handler (+22 more)

### Community 29 - "fixtures.ts"
Cohesion: 0.06
Nodes (42): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+34 more)

### Community 30 - "SuggestRequest"
Cohesion: 0.11
Nodes (12): claudeProvider, openAICompatibleProvider, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList() (+4 more)

### Community 31 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 32 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 33 - "Connector"
Cohesion: 0.08
Nodes (11): init(), ConfigField, Connector, buildGatewayTable(), isTimeout(), primaryGatewayName(), wanInterfaceName(), Connector (+3 more)

### Community 34 - "lists.go"
Cohesion: 0.13
Nodes (31): buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), parseAdlistsV6(), parseClientsV6() (+23 more)

### Community 35 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 36 - "DecodeJSON"
Cohesion: 0.17
Nodes (6): Handler, Handler, oidcProviderJSON(), DecodeJSON(), T, Handler

### Community 37 - "IsSecureRequest"
Cohesion: 0.08
Nodes (29): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), clearOIDCFlowCookie(), emailDomainAllowed(), Handler, newOIDCUser(), oidcFlowCookieName(), oidcRoleForGroups() (+21 more)

### Community 38 - "newTestHandler"
Cohesion: 0.06
Nodes (69): AssertMatchesSpec(), loadSpec(), specPath(), testHandler, Handler, instanceAdminRoleFor(), actionRequest(), actionResponse() (+61 more)

### Community 39 - "docker_test.go"
Cohesion: 0.07
Nodes (41): IsDangerousIP(), buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer() (+33 more)

### Community 40 - "routerDeps"
Cohesion: 0.11
Nodes (30): routerDeps, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+22 more)

### Community 41 - "traefik/tables_test.go"
Cohesion: 0.09
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+25 more)

### Community 42 - "nilToStr"
Cohesion: 0.10
Nodes (11): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+3 more)

### Community 43 - "UsersPage.tsx"
Cohesion: 0.08
Nodes (34): AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), RefreshFn, setRefreshHandler(), server (+26 more)

### Community 44 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 45 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 46 - "response.go"
Cohesion: 0.08
Nodes (27): Handler, Handler, Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor(), TestCursorRequestModes() (+19 more)

### Community 47 - "NewStore"
Cohesion: 0.15
Nodes (27): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+19 more)

### Community 48 - "middleware.go"
Cohesion: 0.17
Nodes (15): AuditRecorder, contextKey, elevationError, PermissionChecker, SecurityHeaders(), TestSecurityHeaders(), elevationFailureReason(), recordElevationAudit() (+7 more)

### Community 49 - "router.go"
Cohesion: 0.14
Nodes (22): TestEmbeddedFrontendEntryPoint(), go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance (+14 more)

### Community 50 - "share_links_test.go"
Cohesion: 0.23
Nodes (35): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), asUser(), createTestShareLink() (+27 more)

### Community 51 - "ShareLinkPage.tsx"
Cohesion: 0.09
Nodes (27): axios, TODO: fold into docs/openapi.yaml and regenerate via `npm run gen:api`, ShareDoc, ShareLink, ShareLinkCreated, shareLinkErrorCode, shareLinksQueryKey, ShareTreeNode (+19 more)

### Community 52 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 53 - "NewRegistry"
Cohesion: 0.14
Nodes (24): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+16 more)

### Community 54 - "Register"
Cohesion: 0.05
Nodes (69): init(), newConnector(), Connector, GuardedDialer(), RequestedFields(), init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal() (+61 more)

### Community 55 - "ThemeControls.tsx"
Cohesion: 0.09
Nodes (28): MotionProvider(), AppearancePage(), ChoiceGroup(), AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented() (+20 more)

### Community 56 - "auth_test.go"
Cohesion: 0.08
Nodes (33): testApp, loginRefreshCookie(), seedLocalUser(), TestChangePasswordWrongCurrentPassword(), TestDeleteSessionNotOwner(), TestDeleteSessionSuccess(), TestElevateSuccess(), TestElevateWrongPassword() (+25 more)

### Community 57 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 58 - "compliance/engine.go"
Cohesion: 0.12
Nodes (29): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+21 more)

### Community 59 - "Config"
Cohesion: 0.20
Nodes (14): Config, LogSettings, IsSSHRemote(), AISettings, BackupSettings, DocExportGitSettings, DocExportSettings, EncryptionSettings (+6 more)

### Community 60 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 61 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 62 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 63 - "WebSocketProvider.tsx"
Cohesion: 0.09
Nodes (26): Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`), web_src_api_generated_changes_changes (+18 more)

### Community 64 - "Connector"
Cohesion: 0.17
Nodes (9): SnapshotSection, groupNames(), parseGroupsV6(), Connector, isTimeout(), unavailable(), parseGroupsV5(), groupRow (+1 more)

### Community 65 - "config_test.go"
Cohesion: 0.08
Nodes (31): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Load(), TestAccessTokenTTLDuration() (+23 more)

### Community 66 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 67 - "HashToken"
Cohesion: 0.12
Nodes (10): Handler, Handler, contextWithShareLink(), Handler, shareLinkFromContext(), InstanceAdminFromContext(), NoContent(), HashToken() (+2 more)

### Community 68 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 69 - "git.go"
Cohesion: 0.08
Nodes (25): TestCommitMessage(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), TestPoolConfigWithDefaults(), TestWithinTransactionRollsBack(), go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_getkin_kin_openapi_openapi3 (+17 more)

### Community 70 - "home_assistant_test.go"
Cohesion: 0.14
Nodes (28): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+20 more)

### Community 71 - "GetTypeSchema"
Cohesion: 0.11
Nodes (26): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestBuildHostsTableAttributes(), TestSchemaExposesAPIVersion() (+18 more)

### Community 72 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 73 - "ExportToFile"
Cohesion: 0.09
Nodes (46): Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+38 more)

### Community 74 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 75 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 76 - "New"
Cohesion: 0.18
Nodes (19): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+11 more)

### Community 77 - "Store"
Cohesion: 0.22
Nodes (22): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+14 more)

### Community 78 - "time.Time"
Cohesion: 0.18
Nodes (18): golang.org/x/time/rate.Limiter, time.Time, visitor, ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection (+10 more)

### Community 79 - "createTestConnector"
Cohesion: 0.12
Nodes (22): TestDeleteOldHealthChecks(), TestGetConnectorUptimeDeterministicOutage(), TestGetConnectorUptimeNoData(), TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR(), TestRecordHealthCheckDefaults(), assertQueryPlanUsesIndex(), Store, TestDeleteOldAlertsSkipsActiveStatuses() (+14 more)

### Community 80 - "src/theme.ts"
Cohesion: 0.15
Nodes (23): @fontsource/space-mono, @fontsource-variable/space-grotesk, ColorMode, commit(), load(), Persisted, PRESETS_FONTS, ThemeState (+15 more)

### Community 81 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 82 - "log/slog.Logger"
Cohesion: 0.12
Nodes (16): newLifecycleManager(), WithLogger(), Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+8 more)

### Community 83 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 84 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 85 - "Manager"
Cohesion: 0.15
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 86 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 87 - "Connector"
Cohesion: 0.14
Nodes (10): TestRateLimit(), apiMessage(), controllerName(), countByKind(), isTimeout(), statusError(), unavailable(), Connector (+2 more)

### Community 88 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 89 - "newTestHandler"
Cohesion: 0.15
Nodes (23): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+15 more)

### Community 90 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 91 - "logging.go"
Cohesion: 0.15
Nodes (18): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestLoggablePathMasksShareTokenUnderV1() (+10 more)

### Community 92 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 93 - "export_test.go"
Cohesion: 0.19
Nodes (19): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+11 more)

### Community 94 - "doc_test.go"
Cohesion: 0.13
Nodes (20): Store, mustCreateUser(), newPostgresTestStore(), skipOnPostgres(), TestDocLockAcquireAfterExpiry(), TestDocLockConflict(), TestDocLockReleaseOnlyByHolder(), TestDocLockRenewalByHolder() (+12 more)

### Community 95 - "chat/chat.go"
Cohesion: 0.17
Nodes (17): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings() (+9 more)

### Community 96 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 97 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 98 - "ConnectorRecord"
Cohesion: 0.17
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 99 - "main"
Cohesion: 0.14
Nodes (16): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+8 more)

### Community 100 - "APIKeyClaims"
Cohesion: 0.50
Nodes (3): testAPIKeyChecker, APIKeyClaims, validAPIKey()

### Community 101 - "diagnostics/diagnostics.go"
Cohesion: 0.21
Nodes (18): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+10 more)

### Community 102 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 103 - "gitAuth"
Cohesion: 0.40
Nodes (5): gitAuth(), installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), GitOptions

### Community 104 - "Hub"
Cohesion: 0.14
Nodes (7): Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator, ticket

### Community 105 - "ws/ws_test.go"
Cohesion: 0.15
Nodes (21): newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock() (+13 more)

### Community 106 - "sshStdioConn"
Cohesion: 0.13
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 107 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 108 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 109 - "docs/handlers_test.go"
Cohesion: 0.19
Nodes (16): NewHandler(), TestAISuggestInvalidJSON(), TestByServiceNoDocsYet(), TestGenerate(), TestGetLockNoneHeld(), TestGetRootIsSynthetic(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestListEmpty() (+8 more)

### Community 110 - "sync.Mutex"
Cohesion: 0.16
Nodes (8): RateLimit(), createVersion(), templateResponse(), templateVersionResponse(), golang.org/x/time/rate.Limit, sync.Mutex, limiterStore, Handler

### Community 111 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 112 - "SnapshotEntity"
Cohesion: 0.14
Nodes (23): ServiceDependency, SnapshotEntity, buildHostsTable(), parseHosts(), TestBuildHostsTableMalformedCases(), TestBuildHostsTableValidRecords(), buildHostsTableV5(), environmentDependencies() (+15 more)

### Community 113 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 114 - "NewEngine"
Cohesion: 0.25
Nodes (24): NewEngine(), newEngineTestStore(), seedEngineConnector(), seedEngineTemplate(), TestGenerateFromSnapshotIncludesDependencies(), TestGenerateFromTemplateReturnsVersionPersistenceError(), TestGenerateFromTemplateStillPersists(), TestMatchingConnectorsEmptyAppliesToIsWildcard() (+16 more)

### Community 115 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 116 - "chat/handlers_test.go"
Cohesion: 0.54
Nodes (7): Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationDocVisibility(), TestCreateConversationValidation(), TestPostMessageErrors()

### Community 117 - "templates_test.go"
Cohesion: 0.22
Nodes (16): templateBody, TestConnectorGrantRouteMatrix(), TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors() (+8 more)

### Community 118 - "reports/handlers.go"
Cohesion: 0.17
Nodes (14): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), NewHandler() (+6 more)

### Community 119 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 121 - "templates.fixtures.ts"
Cohesion: 0.17
Nodes (13): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+5 more)

### Community 122 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 123 - "notifications/handlers_test.go"
Cohesion: 0.17
Nodes (24): NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth() (+16 more)

### Community 124 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 125 - "middleware_test.go"
Cohesion: 0.15
Nodes (15): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, assertElevationAuditCalls(), boolLabel(), requestWithUser(), TestAuthMiddlewareAllowsCurrentRoleClaim(), TestAuthMiddlewareInvalidToken() (+7 more)

### Community 126 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 127 - "Engine"
Cohesion: 0.21
Nodes (7): Engine, dedupKey(), matchReason(), TemplateFuncs(), GenerateResult, renderResult, text/template.FuncMap

### Community 128 - "templatefuncs.go"
Cohesion: 0.17
Nodes (12): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+4 more)

### Community 129 - "NewMalformedResponseError"
Cohesion: 0.18
Nodes (7): NewMalformedResponseError(), WantsField(), isTimeout(), putMetadata(), unavailable(), MalformedResponseError, Connector

### Community 130 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 131 - "testApp"
Cohesion: 0.24
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 132 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 133 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 134 - "Connector"
Cohesion: 0.15
Nodes (6): buildGatewayTable(), buildSystemContent(), isTimeout(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 135 - "gitTarget"
Cohesion: 0.19
Nodes (9): commitMessage(), Exporter, commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash, github.com/go-git/go-git/v5/plumbing/object.Signature, github.com/go-git/go-git/v5/plumbing.ReferenceName (+1 more)

### Community 136 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 137 - "New"
Cohesion: 0.36
Nodes (13): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+5 more)

### Community 138 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 139 - "gitFixture"
Cohesion: 0.29
Nodes (9): SetBeforePushForTest(), keys(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, Exporter (+1 more)

### Community 140 - "NewClient"
Cohesion: 0.24
Nodes (11): clientTimeout(), NewClient(), NewTransport(), NoRedirect(), TestNewClientDoesNotFollowRedirects(), TestNewClientInsecureSkipVerifyConnects(), TestNewClientTimeout(), TestNewClientUsesCustomDialContext() (+3 more)

### Community 141 - "maintenance_test.go"
Cohesion: 0.23
Nodes (12): TestListConnectorNames(), TestListDocsGroupedByService(), Store, mustCreateMaintenanceConnector(), TestCloseMaintenanceWindow(), TestCloseMaintenanceWindowAlreadyClosed(), TestCloseMaintenanceWindowNotFound(), TestCreateAndGetActiveMaintenanceWindow() (+4 more)

### Community 142 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 143 - "Engine"
Cohesion: 0.18
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 145 - "time.Duration"
Cohesion: 0.16
Nodes (6): AuthSettings, Database, Server, time.Duration, OIDCProvider, PoolConfig

### Community 146 - "ReportData"
Cohesion: 0.47
Nodes (4): connectorFilter(), DefinitionSummary, Generator, ReportData

### Community 147 - "quality_test.go"
Cohesion: 0.26
Nodes (10): TestComplianceFindingRuleDedupAndResolve(), TestComplianceRuleCRUD(), Store, newConcurrentQualityTestStore(), seedQualityConnector(), TestListQualityFindingsFilters(), TestResolveThenReopenCreatesFreshRow(), TestUpsertQualityFindingConcurrentDedup() (+2 more)

### Community 148 - "transform.go"
Cohesion: 0.25
Nodes (8): init(), normalizeEnabledColumn(), normalizeFirewallRules(), RegisterTransformer(), TestNormalizeFirewallRulesRewritesEnabledColumn(), TestRunTransformersAppliesInOrderAndStopsOnError(), Transformer, TransformerFunc

### Community 149 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 150 - "WiseLabz Connector Guide"
Cohesion: 0.17
Nodes (12): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+4 more)

### Community 151 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 152 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 154 - "fetch_test.go"
Cohesion: 0.19
Nodes (14): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+6 more)

### Community 155 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 156 - "connector_permission_test.go"
Cohesion: 0.45
Nodes (10): Store, newTestConnector(), newTestUser(), TestDeleteConnectorGrant(), TestFilterConnectorIDsByGrant(), TestListConnectorGrants(), TestUpsertConnectorGrantUpdatesExistingRole(), TestUserHasConnectorRoleDefaultDeny() (+2 more)

### Community 157 - "keyset_test.go"
Cohesion: 0.33
Nodes (10): assertSameSet(), Store, T, queryPlan(), TestKeysetQueriesUseCoveringIndexes(), TestListAuditRecordsKeysetHonoursFilters(), TestListAuditRecordsKeysetTraversal(), TestListChangesKeysetTraversal() (+2 more)

### Community 158 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 159 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 160 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 161 - "AuthMiddleware"
Cohesion: 0.19
Nodes (10): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), AuthMiddleware(), extractBearerToken() (+2 more)

### Community 162 - "changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 163 - "NewService"
Cohesion: 0.30
Nodes (11): NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner(), TestElevationWrongAction(), TestExpiredAccessToken(), TestIssueAndValidateAccess(), TestIssueAndValidateElevation() (+3 more)

### Community 164 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 165 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 166 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 167 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 168 - "connector/connector.go"
Cohesion: 0.20
Nodes (5): CredentialRefresher, Restarter, Starter, Stopper, go_pkg_syscall

### Community 169 - "retention/retention_test.go"
Cohesion: 0.35
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 170 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 171 - "NotificationRecord"
Cohesion: 0.36
Nodes (3): NotificationRecord, Store, scanNotification()

### Community 172 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 173 - "notification_delivery_test.go"
Cohesion: 0.39
Nodes (7): createTestNotification(), Store, TestDeliveryCreateAndList(), TestListDeliveriesStatusFilterAndPagination(), TestListDueDeliveries(), TestUpdateDeliveryResultNotFound(), TestUpdateDeliveryResultTransitionsAndClearsNextAttempt()

### Community 174 - "isUniqueViolation"
Cohesion: 0.35
Nodes (5): RunbookRecord, Store, scanRunbook(), isUniqueViolation(), TestIsUniqueViolation()

### Community 175 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 176 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 177 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 178 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 179 - "newRouterDeps"
Cohesion: 0.11
Nodes (20): Config, TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), NewHandler(), chi.Router, NewRouter(), newRouterDeps(), spaHandler() (+12 more)

### Community 181 - "responseWriter"
Cohesion: 0.18
Nodes (7): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn, responseWriter

### Community 182 - "change_pattern_test.go"
Cohesion: 0.48
Nodes (6): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern()

### Community 183 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 184 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 185 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 186 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 187 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 188 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 189 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 190 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 191 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 192 - ".resolveConfigPusher"
Cohesion: 0.39
Nodes (3): Handler, isWritableField(), ConfigPusher

### Community 193 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 194 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 195 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 196 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 200 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 209 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 210 - ".RunDigestSweep"
Cohesion: 0.40
Nodes (4): digestDue(), formatDigest(), Dispatcher, TestDigestDue()

### Community 211 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 212 - "snapshot_attributes_test.go"
Cohesion: 0.83
Nodes (3): TestSnapshotAttributesRoundTripPostgres(), TestSnapshotAttributesRoundTripSQLite(), testSnapshotWithAttributes()

## Knowledge Gaps
- **542 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+537 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1200 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `testing.T`, `testApp`, `go_pkg_context`, `gitFixture`, `ServiceSnapshot`, `Engine`, `RunMigrations`, `Errorf`, `ReportData`, `rowScanner`, `dispatcher_test.go`, `ErrorWithDetails`, `Handler`, `net/http.Request`, `DecodeJSON`, `.call`, `newTestHandler`, `retention/retention_test.go`, `NewChecker`, `response.go`, `NewStore`, `share_links_test.go`, `newRouterDeps`, `NewRegistry`, `Register`, `Dispatcher`, `HashToken`, `rewritePlaceholders`, `ExportToFile`, `New`, `time.Time`, `.call`, `log/slog.Logger`, `compliance/handlers.go`, `Checker`, `Manager`, `export_test.go`, `chat/chat.go`, `main`, `diagnostics/diagnostics.go`, `docs/handlers_test.go`, `sync.Mutex`, `NewEngine`, `chat/handlers_test.go`, `reports/handlers.go`, `notifications/handlers_test.go`, `changes/handlers_test.go`, `Engine`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `UserIDFromContext()` connect `Errorf` to `HashToken`, `context.Context`, `DecodeJSON`, `routerDeps`, `response.go`, `sync.Mutex`, `middleware.go`, `RequireConnectorRole`, `rowScanner`, `ErrorWithDetails`, `Handler`, `net/http.Request`, `middleware_test.go`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Why does `Runner` connect `Runner` to `testApp`, `context.Context`, `New`, `go_pkg_context`, `sync.Mutex`, `NewStore`, `log/slog.Logger`, `newRouterDeps`, `net/http.Request`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _542 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.0256186824677588 - nodes in this community are weakly interconnected._
- **Should `cn` be split into smaller, more focused modules?**
  _Cohesion score 0.021643460067816176 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025249691392660756 - nodes in this community are weakly interconnected._