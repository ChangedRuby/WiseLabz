# Graph Report - wiselabz-wd-implementation  (2026-09-23)

## Corpus Check
- 785 files · ~475,881 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 5841 nodes · 18361 edges · 211 communities (186 shown, 25 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1451 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `17c16978`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- newDocTestStore
- testing.T
- icons.tsx
- SystemPage.tsx
- go_pkg_net_http
- context.Context
- @tanstack/react-query
- go_pkg_context
- react
- DashboardPage.tsx
- go_pkg_strings
- cn
- ServiceDetailPage.tsx
- App.tsx
- ServiceSnapshot
- states.tsx
- net/http.Request
- go_pkg_testing
- ServicesPage.tsx
- ErrorWithDetails
- dispatcher_test.go
- connector/connector.go
- rowScanner
- NewEngine
- portainer/tables.go
- MarshalConnectorConfig
- net/http.ResponseWriter
- NewMalformedResponseError
- package.json
- ExportToFile
- fixtures.ts
- SnapshotEntity
- WiseLabz — Architecture & Technical Decisions
- IsSecureRequest
- RunMigrations
- home_assistant/tables.go
- traefik/tables.go
- RulesPage.tsx
- buildHostsTable
- NewChecker
- dependencies
- sync.Mutex
- routerDeps
- docker_test.go
- Manager
- adguardhome/tables.go
- Store
- response.go
- Dispatcher
- main
- SuggestWithFallback
- auth_test.go
- AppearancePage.tsx
- Config
- Connector
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- share_links_test.go
- newRouterDeps
- settings.mock.ts
- NewStore
- rewritePlaceholders
- home_assistant_test.go
- net/http.Client
- SuggestRequest
- newTestHandler
- git.go
- newTestHandler
- Register
- GetTypeSchema
- NewEngine
- ThemeControls.tsx
- handlers.ts
- ConnectorRecord
- unifi_test.go
- NotificationRecord
- timeline.ts
- Store
- compliance/handlers.go
- Connector
- config_test.go
- Encrypt
- WiseLabz — Design Contract
- channels_test.go
- Compare
- router.go
- devDependencies
- portainer_test.go
- AuditRecord
- export.go
- backup/main.go
- adguardhome_test.go
- Service
- bulkFakeConnector
- NewHandler
- data.go
- Runner
- diagnostics/diagnostics.go
- changes_test.go
- ws/ws_test.go
- NewService
- Handler
- GuardedDialer
- truenas_test.go
- ws.ts
- compilerOptions
- runRestore
- chat/chat.go
- traefik_test.go
- docdiffmodel.ts
- templates_test.go
- testApp
- handlers_contract_test.go
- Handler
- system/handlers_test.go
- all.go
- sshStdioConn
- nilToStr
- .batchDelete
- templates.fixtures.ts
- compilerOptions
- notifications/handlers_test.go
- NewRegistry
- handlers_actions_test.go
- vectorCache
- Connector
- Handler
- registry.go
- compliance_rules_test.go
- provider_test.go
- New
- RateLimit
- pagination_contract_test.go
- newTestHandler
- Connector
- .Fetch
- render_test.go
- DecodeKey
- Handler
- .doRequest
- Connector
- retention/retention_test.go
- Contributing to WiseLabz
- time.Time
- Engine
- ShareLink
- Store
- .doRequest
- Decision
- main.tsx
- scripts
- Handler
- .call
- Store
- runbooks/handlers_test.go
- Decision
- WiseLabz Connector Guide
- Product
- .call
- walkCursorPages
- handlers_bulk_test.go
- redactDSN
- serveSSHDockerConn
- .Fetch
- Changelog
- Connector
- mockServiceWorker.js
- config_cmd_test.go
- openapi_contract_test.go
- Provider
- .Fetch
- ReportData
- healthFakeConnector
- seedScopeFixture
- release-please-config.json
- crypto.go
- ComputeWindow
- BackupSchedule
- Cache
- Step by step
- WiseLabz
- .applyChannelSecrets
- stubEmbedder
- scanMaintenanceWindow
- computeNextRun
- fakeQualityChecker
- Contributor Covenant Code of Conduct
- registryTestRefresher
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- .doRequestV5
- Store
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- newHandler
- Enforcement Guidelines
- compose-smoke.sh
- timeoutError
- MISSING — deferred & future frontend features
- Saved Views
- WiseLabz — v2 Backlog
- tsconfig.json
- AGENTS.md
- setup-env.sh
- CHANGE_PROVENANCE.md
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 212 edges
2. `Errorf()` - 164 edges
3. `Store` - 133 edges
4. `newDocTestStore()` - 129 edges
5. `SnapshotEntity` - 74 edges
6. `react` - 72 edges
7. `UserIDFromContext()` - 69 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
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

## Communities (211 total, 25 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (164): testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow(), TestAlertsListSuccess() (+156 more)

### Community 1 - "newDocTestStore"
Cohesion: 0.02
Nodes (139): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+131 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (122): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+114 more)

### Community 3 - "icons.tsx"
Cohesion: 0.06
Nodes (44): web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid (+36 more)

### Community 4 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (95): react-i18next, sonner, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_usegetauthapikeys, web_src_api_generated_auth_auth_usegetauthproviders, web_src_api_generated_me_me_deletemesessionssessionid (+87 more)

### Community 5 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (49): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest (+41 more)

### Community 6 - "context.Context"
Cohesion: 0.03
Nodes (32): fakeStatusChecker, testAPIKeyChecker, sanitize(), sanitizeSessions(), Handler, APIKeyClaims, validAPIKey(), Connector (+24 more)

### Community 7 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (57): Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`), i18next (+49 more)

### Community 8 - "go_pkg_context"
Cohesion: 0.08
Nodes (18): versionSections(), ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), TemplateVersionSection, contains(), searchString(), go_pkg_context (+10 more)

### Community 9 - "react"
Cohesion: 0.04
Nodes (69): Frontend shell & theme (decided 2026-06), react, react-router-dom, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey (+61 more)

### Community 10 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (91): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+83 more)

### Community 11 - "go_pkg_strings"
Cohesion: 0.06
Nodes (36): contextKey, elevationError, dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin() (+28 more)

### Community 12 - "cn"
Cohesion: 0.06
Nodes (40): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid (+32 more)

### Community 13 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (79): ADR-0001, ADR-0003, RFC-3339, 1. `service.status`, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations (+71 more)

### Community 14 - "App.tsx"
Cohesion: 0.05
Nodes (47): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), RefreshFn, setAccessToken(), setRefreshHandler(), server (+39 more)

### Community 15 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (22): noopValidatedConnector, ServiceSnapshot, agentEnabled(), Connector, Sanitize(), TestSanitize(), changePatternID(), Engine (+14 more)

### Community 16 - "states.tsx"
Cohesion: 0.06
Nodes (53): axios, customInstance(), web_src_api_generated_docs_docs_postdocstopology, web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetpassword (+45 more)

### Community 17 - "net/http.Request"
Cohesion: 0.06
Nodes (27): Handler, Registry, Handler, diffToSpec(), NewHandler(), NewHandler(), Handler, Handler (+19 more)

### Community 18 - "go_pkg_testing"
Cohesion: 0.06
Nodes (17): TestLoggablePathMasksShareTokenUnderV1(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), jsonType(), TestAttributeCatalogCoversEmittedKeys() (+9 more)

### Community 19 - "ServicesPage.tsx"
Cohesion: 0.04
Nodes (83): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss (+75 more)

### Community 20 - "ErrorWithDetails"
Cohesion: 0.07
Nodes (29): updateUserRequest, newToken(), Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword() (+21 more)

### Community 21 - "dispatcher_test.go"
Cohesion: 0.06
Nodes (78): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), testLogger(), expireAlertsOnce() (+70 more)

### Community 22 - "connector/connector.go"
Cohesion: 0.06
Nodes (26): isTimeout(), IsDangerousIP(), NewAuthError(), NewServiceUnavailableError(), NewTimeoutError(), Connector, isTimeout(), TestTypedErrorsAreDistinguishableByType() (+18 more)

### Community 23 - "rowScanner"
Cohesion: 0.06
Nodes (22): existingIDs(), docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary(), countRows() (+14 more)

### Community 24 - "NewEngine"
Cohesion: 0.09
Nodes (44): NewHandler(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine (+36 more)

### Community 25 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 26 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (15): TestDiagnosticsRedactsSecrets(), IsSecretFieldType(), MarshalConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+7 more)

### Community 27 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (25): Handler, isWritableField(), validateConfigPushRequest(), applyConnectorScalarUpdates(), configRequestField(), Handler, validateConnectorConfig(), writeConfigRejection() (+17 more)

### Community 28 - "NewMalformedResponseError"
Cohesion: 0.08
Nodes (41): unavailable(), SnapshotSection, NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell() (+33 more)

### Community 29 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 30 - "ExportToFile"
Cohesion: 0.12
Nodes (35): ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+27 more)

### Community 31 - "fixtures.ts"
Cohesion: 0.06
Nodes (41): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+33 more)

### Community 32 - "SnapshotEntity"
Cohesion: 0.11
Nodes (46): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildInterfaceTableAttributes(), buildInterfaceTable(), buildDatasets(), buildDisks(), buildInterfaces() (+38 more)

### Community 33 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (44): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+36 more)

### Community 34 - "IsSecureRequest"
Cohesion: 0.08
Nodes (27): TestEmailDomainAllowed(), TestOIDCRoleForGroups(), clearOIDCFlowCookie(), emailDomainAllowed(), Handler, newOIDCUser(), oidcFlowCookieName(), oidcRoleForGroups() (+19 more)

### Community 35 - "RunMigrations"
Cohesion: 0.10
Nodes (38): main(), newScratchStore(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns() (+30 more)

### Community 36 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 37 - "traefik/tables.go"
Cohesion: 0.14
Nodes (31): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+23 more)

### Community 38 - "RulesPage.tsx"
Cohesion: 0.05
Nodes (42): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+34 more)

### Community 39 - "buildHostsTable"
Cohesion: 0.33
Nodes (7): TestBuildHostsTableAttributes(), TestBuildHostsTableV5(), buildHostsTable(), parseHosts(), TestBuildHostsTableMalformedCases(), TestBuildHostsTableValidRecords(), buildHostsTableV5()

### Community 40 - "NewChecker"
Cohesion: 0.06
Nodes (70): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+62 more)

### Community 41 - "dependencies"
Cohesion: 0.05
Nodes (40): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+32 more)

### Community 42 - "sync.Mutex"
Cohesion: 0.13
Nodes (7): cron.EntryID, Handler, createVersion(), templateResponse(), templateVersionResponse(), sync.Mutex, Handler

### Community 43 - "routerDeps"
Cohesion: 0.09
Nodes (37): routerDeps, AuditRecorder, PermissionChecker, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router (+29 more)

### Community 44 - "docker_test.go"
Cohesion: 0.06
Nodes (43): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), startSSHDockerServer(), TestConfigPush() (+35 more)

### Community 45 - "Manager"
Cohesion: 0.08
Nodes (17): NewHandler(), cron.EntryID, Manager, LogPartial(), NewManager(), Store, Store, ReportDefinitionRecord (+9 more)

### Community 46 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 47 - "Store"
Cohesion: 0.08
Nodes (11): Store, placeholders(), changeFilterClause(), AlertRecord, ChangeRecord, Store, scanAlert(), scanChange() (+3 more)

### Community 48 - "response.go"
Cohesion: 0.08
Nodes (29): Handler, parseScheduleUpdates(), validateRotationFields(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+21 more)

### Community 49 - "Dispatcher"
Cohesion: 0.14
Nodes (14): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+6 more)

### Community 50 - "main"
Cohesion: 0.11
Nodes (24): main(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder, EmbedRegistry (+16 more)

### Community 51 - "SuggestWithFallback"
Cohesion: 0.20
Nodes (12): StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds() (+4 more)

### Community 52 - "auth_test.go"
Cohesion: 0.08
Nodes (33): testApp, loginRefreshCookie(), seedLocalUser(), TestChangePasswordWrongCurrentPassword(), TestDeleteSessionNotOwner(), TestDeleteSessionSuccess(), TestElevateSuccess(), TestElevateWrongPassword() (+25 more)

### Community 53 - "AppearancePage.tsx"
Cohesion: 0.14
Nodes (21): zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css() (+13 more)

### Community 54 - "Config"
Cohesion: 0.09
Nodes (20): Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings, BackupSettings, Database, DocExportGitSettings (+12 more)

### Community 55 - "Connector"
Cohesion: 0.10
Nodes (9): ConfigField, Connector, buildGatewayTable(), isTimeout(), primaryGatewayName(), wanInterfaceName(), PathSegment(), ValidateRefSegment() (+1 more)

### Community 56 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 57 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 58 - "share_links_test.go"
Cohesion: 0.19
Nodes (38): GrantConnectorRole(), instanceAdminRole(), NewUser(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld(), TestGetRootIsSynthetic() (+30 more)

### Community 59 - "newRouterDeps"
Cohesion: 0.14
Nodes (17): Config, NewHandler(), NewHandler(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders() (+9 more)

### Community 60 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 61 - "NewStore"
Cohesion: 0.15
Nodes (27): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+19 more)

### Community 62 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 63 - "home_assistant_test.go"
Cohesion: 0.14
Nodes (28): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+20 more)

### Community 64 - "net/http.Client"
Cohesion: 0.11
Nodes (6): ollamaEmbedder, openAIEmbedder, Connector, Connector, newWebhookClient(), net/http.Client

### Community 65 - "SuggestRequest"
Cohesion: 0.21
Nodes (5): claudeProvider, openAICompatibleProvider, SuggestChunk, SuggestRequest, countingProvider

### Community 66 - "newTestHandler"
Cohesion: 0.22
Nodes (17): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+9 more)

### Community 67 - "git.go"
Cohesion: 0.10
Nodes (20): keys(), go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter, go_pkg_github_com_getkin_kin_openapi_routers, go_pkg_github_com_getkin_kin_openapi_routers_gorillamux, go_pkg_github_com_go_git_go_git_v5 (+12 more)

### Community 68 - "newTestHandler"
Cohesion: 0.12
Nodes (27): TestConnectorStoreErrorPaths(), TestGetRequiresGrantEvenForInstanceAdmin(), TestListFiltersGrantsBeforePagination(), createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler() (+19 more)

### Community 69 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 70 - "GetTypeSchema"
Cohesion: 0.12
Nodes (25): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+17 more)

### Community 71 - "NewEngine"
Cohesion: 0.18
Nodes (23): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+15 more)

### Community 72 - "ThemeControls.tsx"
Cohesion: 0.11
Nodes (31): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), FONT_KEYS, OPT_KEYS, PRESET_KEYS, Segmented(), ThemeControls() (+23 more)

### Community 73 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 74 - "ConnectorRecord"
Cohesion: 0.17
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 75 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 76 - "NotificationRecord"
Cohesion: 0.29
Nodes (4): Dispatcher, NotificationRecord, Store, scanNotification()

### Community 77 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 78 - "Store"
Cohesion: 0.22
Nodes (22): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+14 more)

### Community 79 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 80 - "Connector"
Cohesion: 0.15
Nodes (9): apiMessage(), controllerName(), countByKind(), isTimeout(), statusError(), unavailable(), Connector, sectionFetch (+1 more)

### Community 81 - "config_test.go"
Cohesion: 0.12
Nodes (22): runHealthcheck(), Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH() (+14 more)

### Community 82 - "Encrypt"
Cohesion: 0.29
Nodes (12): Decrypt(), DeriveKey(), Encrypt(), TestDecodeKeyUsableForEncryptDecrypt(), TestDecryptTampered(), TestDecryptWrongKey(), TestDeriveKeyDeterministic(), TestDeriveKeyDifferent() (+4 more)

### Community 83 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 84 - "channels_test.go"
Cohesion: 0.20
Nodes (14): buildEmailMessage(), sendNtfyChannel(), sendSMTPChannel(), sendTelegramChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendNtfyChannel_DefaultsToNtfySh(), TestSendNtfyChannel_MissingTopic() (+6 more)

### Community 85 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 86 - "router.go"
Cohesion: 0.17
Nodes (19): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+11 more)

### Community 87 - "devDependencies"
Cohesion: 0.09
Nodes (23): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+15 more)

### Community 88 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 89 - "AuditRecord"
Cohesion: 0.31
Nodes (6): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), AuditRecord

### Community 90 - "export.go"
Cohesion: 0.06
Nodes (41): fetchAllDocs(), fileName(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), slugify() (+33 more)

### Community 91 - "backup/main.go"
Cohesion: 0.14
Nodes (18): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+10 more)

### Community 92 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 93 - "Service"
Cohesion: 0.20
Nodes (10): Claims, ElevationClaims, ElevationToken, TokenPair, Service, hasAudience(), newTokenID(), go_pkg_github_com_golang_jwt_jwt_v5 (+2 more)

### Community 94 - "bulkFakeConnector"
Cohesion: 0.14
Nodes (3): actionConnector, bulkFakeConnector, failingPushConnector

### Community 95 - "NewHandler"
Cohesion: 0.20
Nodes (12): NewHandler(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestRestore(), TestTemplateSchema(), TestTree(), TestTreeEmpty() (+4 more)

### Community 96 - "data.go"
Cohesion: 0.27
Nodes (14): ChangeEntry, ComplianceSection, ConnectorDrift, DocChangeEntry, DocsSection, DriftSection, FindingSummary, JobHealthEntry (+6 more)

### Community 97 - "Runner"
Cohesion: 0.07
Nodes (31): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+23 more)

### Community 98 - "diagnostics/diagnostics.go"
Cohesion: 0.22
Nodes (17): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+9 more)

### Community 99 - "changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 100 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 101 - "NewService"
Cohesion: 0.07
Nodes (44): APIKeyChecker, ConnectorRoleChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle() (+36 more)

### Community 102 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 103 - "GuardedDialer"
Cohesion: 0.14
Nodes (13): newConnector(), Connector, GuardedDialer(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector() (+5 more)

### Community 104 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 105 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 106 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 107 - "runRestore"
Cohesion: 0.12
Nodes (22): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+14 more)

### Community 108 - "chat/chat.go"
Cohesion: 0.17
Nodes (17): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings() (+9 more)

### Community 109 - "traefik_test.go"
Cohesion: 0.11
Nodes (30): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+22 more)

### Community 110 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 111 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 112 - "testApp"
Cohesion: 0.24
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 113 - "handlers_contract_test.go"
Cohesion: 0.25
Nodes (15): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+7 more)

### Community 114 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 115 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 116 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 117 - "sshStdioConn"
Cohesion: 0.12
Nodes (10): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session, io.Closer (+2 more)

### Community 118 - "nilToStr"
Cohesion: 0.10
Nodes (10): nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store, scanDelivery(), Store (+2 more)

### Community 120 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (12): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate(), resolveToken() (+4 more)

### Community 121 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 122 - "notifications/handlers_test.go"
Cohesion: 0.16
Nodes (25): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+17 more)

### Community 123 - "NewRegistry"
Cohesion: 0.18
Nodes (24): NewRegistry(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+16 more)

### Community 124 - "handlers_actions_test.go"
Cohesion: 0.32
Nodes (14): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+6 more)

### Community 125 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 126 - "Connector"
Cohesion: 0.16
Nodes (7): TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), isTimeout(), Connector

### Community 127 - "Handler"
Cohesion: 0.23
Nodes (4): Handler, Handler, oidcProviderJSON(), boolToInt()

### Community 128 - "registry.go"
Cohesion: 0.21
Nodes (9): IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError, Factory, SchemaField (+1 more)

### Community 129 - "compliance_rules_test.go"
Cohesion: 0.27
Nodes (9): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), go_pkg_regexp (+1 more)

### Community 130 - "provider_test.go"
Cohesion: 0.29
Nodes (6): testProvider, TestRegistryGet(), TestRegistryList(), TestStubProviderName(), TestStubProviderSuggest(), TestStubProviderSuggestStream()

### Community 131 - "New"
Cohesion: 0.09
Nodes (28): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+20 more)

### Community 132 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 133 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 134 - "newTestHandler"
Cohesion: 0.26
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 135 - "Connector"
Cohesion: 0.21
Nodes (5): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), Connector

### Community 136 - ".Fetch"
Cohesion: 0.24
Nodes (5): setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector

### Community 137 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 138 - "DecodeKey"
Cohesion: 0.33
Nodes (5): ProviderConfig, primaryProviderConfig(), DecodeKey(), TestDecodeKey(), AIConfigValues

### Community 141 - "Connector"
Cohesion: 0.14
Nodes (8): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), isTimeout(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 142 - "retention/retention_test.go"
Cohesion: 0.35
Nodes (10): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories() (+2 more)

### Community 143 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 144 - "time.Time"
Cohesion: 0.29
Nodes (6): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, userStatus

### Community 145 - "Engine"
Cohesion: 0.18
Nodes (6): NewHandler(), Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 147 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 149 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 150 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 151 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 153 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 154 - "Store"
Cohesion: 0.33
Nodes (3): Store, scanConnectorGrants(), ConnectorGrant

### Community 155 - "runbooks/handlers_test.go"
Cohesion: 0.48
Nodes (6): Handler, newTestHandler(), TestCreate(), TestGetNotFound(), TestListMutuallyExclusiveFilters(), TestUpdateAndDelete()

### Community 156 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 157 - "WiseLabz Connector Guide"
Cohesion: 0.17
Nodes (12): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+4 more)

### Community 158 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 159 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 160 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 161 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 162 - "redactDSN"
Cohesion: 0.29
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 163 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 164 - ".Fetch"
Cohesion: 0.12
Nodes (13): ServiceDependency, WantsField(), TestRequestedFields(), TestWantsField(), environmentDependencies(), isTimeout(), putMetadata(), unavailable() (+5 more)

### Community 165 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 167 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 168 - "config_cmd_test.go"
Cohesion: 0.23
Nodes (11): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+3 more)

### Community 169 - "openapi_contract_test.go"
Cohesion: 0.48
Nodes (6): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 171 - ".Fetch"
Cohesion: 0.28
Nodes (3): isTimeout(), unavailable(), Connector

### Community 172 - "ReportData"
Cohesion: 0.35
Nodes (6): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), DefinitionSummary, Generator, ReportData

### Community 174 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 175 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 177 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 178 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 179 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 180 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 181 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 184 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 185 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 187 - "Contributor Covenant Code of Conduct"
Cohesion: 0.29
Nodes (7): Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Responsibilities, Our Pledge, Our Standards, Scope

### Community 189 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 190 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 191 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 197 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 198 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 199 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 200 - "newHandler"
Cohesion: 0.18
Nodes (11): testHandler, Handler, instanceAdminRoleFor(), Handler, newHandler(), serve(), TestConversationOwnership(), TestCreateConversationDocVisibility() (+3 more)

### Community 202 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 204 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 208 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 209 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **540 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+535 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1184 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `New`, `go_pkg_context`, `Handler`, `retention/retention_test.go`, `ServiceSnapshot`, `time.Time`, `net/http.Request`, `Engine`, `ErrorWithDetails`, `dispatcher_test.go`, `rowScanner`, `Handler`, `NewEngine`, `.call`, `net/http.ResponseWriter`, `ExportToFile`, `.call`, `RunMigrations`, `NewChecker`, `sync.Mutex`, `ReportData`, `Manager`, `response.go`, `Dispatcher`, `share_links_test.go`, `newRouterDeps`, `NewStore`, `rewritePlaceholders`, `engine_maintenance_test.go`, `NewEngine`, `newHandler`, `compliance/handlers.go`, `export.go`, `NewHandler`, `diagnostics/diagnostics.go`, `runRestore`, `chat/chat.go`, `testApp`, `Handler`, `notifications/handlers_test.go`, `NewRegistry`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `UserIDFromContext()` connect `net/http.Request` to `NewService`, `context.Context`, `sync.Mutex`, `routerDeps`, `go_pkg_strings`, `Handler`, `response.go`, `Handler`, `ErrorWithDetails`, `Handler`, `AuditRecord`, `net/http.ResponseWriter`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Why does `Runner` connect `Runner` to `context.Context`, `go_pkg_context`, `sync.Mutex`, `testApp`, `dispatcher_test.go`, `newRouterDeps`, `NewStore`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _540 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021477343265052764 - nodes in this community are weakly interconnected._
- **Should `newDocTestStore` be split into smaller, more focused modules?**
  _Cohesion score 0.024801005446166736 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024448419797257006 - nodes in this community are weakly interconnected._