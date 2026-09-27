# Graph Report - feat-runbook-execution-with-approvals  (2026-09-27)

## Corpus Check
- 861 files · ~529,476 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6385 nodes · 20300 edges · 212 communities (198 shown, 14 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1650 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5cf3c733`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- newTestApp
- react
- testing.T
- newDocTestStore
- context.Context
- DashboardPage.tsx
- go_pkg_testing
- NewChecker
- net/http.Request
- SystemPage.tsx
- go_pkg_context
- cn
- @tanstack/react-query
- react-router-dom
- net/http.Client
- RunMigrations
- go_pkg_net_http
- Errorf
- ServiceSnapshot
- ConnectorRecord
- ServiceDetailPage.tsx
- App.tsx
- go_pkg_os
- ProfilePage.tsx
- icons.tsx
- Store
- ShareLinkPage.tsx
- Get
- NewEngine
- dispatcher_test.go
- SnapshotEntity
- Connector
- Manager
- package.json
- ExportToFile
- NewUser
- Compare
- fixtures.ts
- WiseLabz — Architecture & Technical Decisions
- routerDeps
- ErrorWithDetails
- docker_test.go
- traefik_test.go
- home_assistant/tables.go
- dependencies
- DecodeKey
- NewStore
- portainer/tables.go
- home_assistant_test.go
- NewMalformedResponseError
- Dispatcher
- SnapshotsPage.tsx
- adguardhome/tables.go
- main
- git.go
- NewRegistry
- connector_permission.go
- response.go
- channels.go
- traefik/tables.go
- SuggestRequest
- Store
- truenas_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- middleware.go
- Connector
- AuthMiddleware
- newTestHandler
- settings.mock.ts
- User
- AuthedUser
- NewEngine
- rewritePlaceholders
- api/auth/webauthn_test.go
- Config
- Service
- src/theme.ts
- handlers.ts
- connector/connector.go
- unifi_test.go
- motion
- timeline.ts
- Handler
- newTestHandler
- WiseLabz — Design Contract
- devDependencies
- time.Time
- portainer_test.go
- Connector
- log/slog.Logger
- logging.go
- adguardhome_test.go
- sshStdioConn
- Handler
- NewService
- NewClient
- NewHTTPClient
- Runner
- Connector
- httpx/retry_test.go
- Deps
- Handler
- ContextWithUser
- export_test.go
- time.Duration
- Hub
- ws/ws_test.go
- ws.ts
- compilerOptions
- Handler
- Register
- .Fetch
- RunbookRecord
- docdiffmodel.ts
- templates_test.go
- api/mcp_test.go
- system/handlers_test.go
- all.go
- DocRecord
- .batchDelete
- compilerOptions
- testApp
- handlers_contract_test.go
- changes/handlers_test.go
- handlers_actions_test.go
- Handler
- Store
- vectorCache
- registry.go
- doJSON
- gitTarget
- gitFixture
- render_test.go
- New
- config_cmd_test.go
- api/changes_test.go
- export.go
- connectors_health_test.go
- chat/chat.go
- createUser
- NotificationRecord
- Contributing to WiseLabz
- scripts
- templatefuncs.go
- IsSecureRequest
- Store
- Store
- JobHealthRecord
- Decision
- main.tsx
- Handler
- Decision
- WiseLabz Connector Guide
- Product
- .call
- cursor_pagination_test.go
- TestElevateOIDC
- handlers_bulk_test.go
- .SnapshotDiff
- .call
- diagnostics/diagnostics.go
- Changelog
- mockServiceWorker.js
- ComplianceRuleRecord
- connectors_hardening_test.go
- scheduler/health_test.go
- BackupSchedule
- release-please-config.json
- auth/handlers_test.go
- Engine
- RateLimit
- gitAuth
- cursor_test.go
- ComputeWindow
- Cache
- Step by step
- WiseLabz
- lifecycleDeps
- session_test.go
- openapi_contract_test.go
- scanMaintenanceWindow
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- testHandler
- Store
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- RequireConnectorRole
- Enforcement Guidelines
- MfaEnrollDialog
- @vitejs/plugin-react
- compose-smoke.sh
- internal/auth/oidc.go
- ClassifyHealth
- timeoutError
- RetentionSettings
- MISSING — deferred & future frontend features
- Saved Views
- dockerSSHAddr
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
4. `Store` - 141 edges
5. `UserIDFromContext()` - 84 edges
6. `SnapshotEntity` - 77 edges
7. `react` - 75 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
10. `@tanstack/react-query` - 62 edges

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

## Communities (212 total, 14 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (194): runbookResp, runbookStepResp, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation() (+186 more)

### Community 1 - "react"
Cohesion: 0.03
Nodes (140): match-sorter, @radix-ui/react-popover, react, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze (+132 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (161): runHealthcheck(), TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL() (+153 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.02
Nodes (154): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+146 more)

### Community 4 - "context.Context"
Cohesion: 0.03
Nodes (33): fakeStatusChecker, sanitizeSessions(), Handler, MFAEnrollOnlyFromContext(), Connector, Connector, Connector, existingIDs() (+25 more)

### Community 5 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (98): Frontend, 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 4. `change.detected` (+90 more)

### Community 6 - "go_pkg_testing"
Cohesion: 0.05
Nodes (29): TestLoggablePathMasksShareTokenUnderV1(), TestExtractGroups(), TestSchemaMatchesConfig(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides() (+21 more)

### Community 7 - "NewChecker"
Cohesion: 0.06
Nodes (73): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+65 more)

### Community 8 - "net/http.Request"
Cohesion: 0.05
Nodes (39): updateUserRequest, Handler, writeUserWriteError(), clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName() (+31 more)

### Community 9 - "SystemPage.tsx"
Cohesion: 0.03
Nodes (75): web_src_api_generated_connectors_connectors_usegetconnectors, web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey, web_src_api_generated_settings_settings_getgetauthconfigquerykey, web_src_api_generated_settings_settings_getgetnotificationsconfigquerykey, web_src_api_generated_settings_settings_postaiconfigtest, web_src_api_generated_settings_settings_postnotificationsconfigtest (+67 more)

### Community 10 - "go_pkg_context"
Cohesion: 0.08
Nodes (16): StatusError, dashboardLayout, contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt (+8 more)

### Community 11 - "cn"
Cohesion: 0.04
Nodes (73): 1. `service.status`, web_src_api_generated_docs_docs_usegetdocstemplateschema, web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore (+65 more)

### Community 12 - "@tanstack/react-query"
Cohesion: 0.04
Nodes (53): 3. `sync.complete`, i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_changes_changes (+45 more)

### Community 13 - "react-router-dom"
Cohesion: 0.05
Nodes (70): Frontend shell & theme (decided 2026-06), react-router-dom, sonner, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey (+62 more)

### Community 14 - "net/http.Client"
Cohesion: 0.04
Nodes (29): Connector, ollamaEmbedder, NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), Connector (+21 more)

### Community 15 - "RunMigrations"
Cohesion: 0.04
Nodes (79): confirm(), formatCounts(), runRestore(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag(), main() (+71 more)

### Community 16 - "go_pkg_net_http"
Cohesion: 0.08
Nodes (27): bulkSnoozeItemResult, bulkSnoozeRequest, shareLinkContextKey, go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_encoding_base64, go_pkg_github_com_go_webauthn_webauthn_protocol, go_pkg_github_com_go_webauthn_webauthn_webauthn (+19 more)

### Community 17 - "Errorf"
Cohesion: 0.06
Nodes (25): Handler, PermissionChecker, newToken(), sanitize(), diffToSpec(), Handler, Handler, Handler (+17 more)

### Community 18 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (22): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, init(), RegisterTransformer() (+14 more)

### Community 19 - "ConnectorRecord"
Cohesion: 0.05
Nodes (40): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), changeFilterClause(), ConnectorRecord, Store (+32 more)

### Community 20 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (59): ADR-0001, ADR-0003, RFC-3339, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush (+51 more)

### Community 21 - "App.tsx"
Cohesion: 0.05
Nodes (56): react-error-boundary, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthloginmfa, web_src_api_generated_auth_auth_postauthlogout, web_src_api_generated_auth_auth_postauthoidccallback, web_src_api_generated_auth_auth_postauthrefresh (+48 more)

### Community 22 - "go_pkg_os"
Cohesion: 0.05
Nodes (47): main(), usage(), changePromptData(), stripPromptTags(), truncateUTF8(), buildPrompt(), TestBuildPrompt(), bulkResolveItemResult (+39 more)

### Community 23 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (52): web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin (+44 more)

### Community 24 - "icons.tsx"
Cohesion: 0.04
Nodes (53): web_src_api_generated_attention_attention, web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_chat_chat (+45 more)

### Community 25 - "Store"
Cohesion: 0.06
Nodes (19): Sanitize(), TestSanitize(), changeServiceIDs(), Store, placeholders(), AlertRecord, ChangeRecord, Store (+11 more)

### Community 26 - "ShareLinkPage.tsx"
Cohesion: 0.06
Nodes (46): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+38 more)

### Community 27 - "Get"
Cohesion: 0.06
Nodes (33): Handler, isWritableField(), validateConfigPushRequest(), applyConnectorScalarUpdates(), configRequestField(), Handler, parseScheduleUpdates(), validateConnectorConfig() (+25 more)

### Community 28 - "NewEngine"
Cohesion: 0.09
Nodes (44): NewHandler(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine (+36 more)

### Community 29 - "dispatcher_test.go"
Cohesion: 0.16
Nodes (53): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), testLogger(), expireAlertsOnce(), NewDispatcher() (+45 more)

### Community 30 - "SnapshotEntity"
Cohesion: 0.10
Nodes (50): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildInterfaceTableAttributes() (+42 more)

### Community 31 - "Connector"
Cohesion: 0.06
Nodes (19): ConfigField, Connector, TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), buildRouteTable(), buildGatewayTable() (+11 more)

### Community 32 - "Manager"
Cohesion: 0.07
Nodes (20): connectorFilter(), cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), Store, nilToStr() (+12 more)

### Community 33 - "package.json"
Cohesion: 0.04
Nodes (46): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+38 more)

### Community 34 - "ExportToFile"
Cohesion: 0.09
Nodes (48): runVerify(), TestRunVerifyDefaultsToLatestInDir(), TestRunVerifyPassAndFail(), TestRunVerifyRequiresABundle(), Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage() (+40 more)

### Community 35 - "NewUser"
Cohesion: 0.13
Nodes (49): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestGetRequiresGrantEvenForInstanceAdmin(), TestListFiltersGrantsBeforePagination(), Handler, snapshotFixture(), snapshotRequest() (+41 more)

### Community 36 - "Compare"
Cohesion: 0.07
Nodes (44): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+36 more)

### Community 37 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 38 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (46): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+38 more)

### Community 39 - "routerDeps"
Cohesion: 0.09
Nodes (38): routerDeps, AuditRecorder, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+30 more)

### Community 40 - "ErrorWithDetails"
Cohesion: 0.08
Nodes (23): oidcElevateFlow, sanitizeUser(), setRefreshCookie(), Handler, mustHashDummyPassword(), Handler, readOIDCElevateFlowCookie(), randomOIDCToken() (+15 more)

### Community 41 - "docker_test.go"
Cohesion: 0.06
Nodes (44): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+36 more)

### Community 42 - "traefik_test.go"
Cohesion: 0.08
Nodes (41): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+33 more)

### Community 43 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 44 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 45 - "DecodeKey"
Cohesion: 0.11
Nodes (22): ProviderConfig, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey(), Decrypt(), DeriveKey() (+14 more)

### Community 46 - "NewStore"
Cohesion: 0.11
Nodes (36): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+28 more)

### Community 47 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 48 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (36): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+28 more)

### Community 49 - "NewMalformedResponseError"
Cohesion: 0.13
Nodes (35): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+27 more)

### Community 50 - "Dispatcher"
Cohesion: 0.12
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 51 - "SnapshotsPage.tsx"
Cohesion: 0.07
Nodes (31): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotssnapshotid, web_src_api_generated_notifications_notifications, web_src_api_generated_notifications_notifications_getgetnotificationsquerykey (+23 more)

### Community 52 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 53 - "main"
Cohesion: 0.08
Nodes (27): Config, main(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder (+19 more)

### Community 54 - "git.go"
Cohesion: 0.08
Nodes (30): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), TestCommitMessage() (+22 more)

### Community 55 - "NewRegistry"
Cohesion: 0.13
Nodes (27): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+19 more)

### Community 56 - "connector_permission.go"
Cohesion: 0.10
Nodes (20): APIKeyRestriction, testAPIKeyChecker, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), treatAsSafeFromContext() (+12 more)

### Community 57 - "response.go"
Cohesion: 0.08
Nodes (19): Handler, Handler, TestWritePaginatedOmitsNextCursor(), Error(), DataPaginatedResponse, HandleStoreError(), intQuery(), JSON() (+11 more)

### Community 58 - "channels.go"
Cohesion: 0.08
Nodes (25): Config, mask(), redactDSN(), redactKVPassword(), Config, TestEveryKeyEnvOverridable(), TestRedactDSN(), TestRedacted() (+17 more)

### Community 59 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 60 - "SuggestRequest"
Cohesion: 0.09
Nodes (14): claudeProvider, openAICompatibleProvider, openAIEmbedder, StubProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet() (+6 more)

### Community 61 - "Store"
Cohesion: 0.16
Nodes (29): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+21 more)

### Community 62 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 63 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 64 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 65 - "middleware.go"
Cohesion: 0.09
Nodes (14): contextKey, elevationError, go_pkg_bytes, go_pkg_crypto_hmac, go_pkg_crypto_sha256, go_pkg_encoding_hex, go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter (+6 more)

### Community 66 - "Connector"
Cohesion: 0.15
Nodes (11): SnapshotSection, TestBuildHostsTableAttributes(), TestBuildHostsTableV5(), buildHostsTable(), Connector, parseHosts(), TestBuildHostsTableMalformedCases(), TestBuildHostsTableValidRecords() (+3 more)

### Community 67 - "AuthMiddleware"
Cohesion: 0.09
Nodes (25): APIKeyChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, UserStatusChecker, AuthMiddleware(), extractBearerToken(), hashToken() (+17 more)

### Community 68 - "newTestHandler"
Cohesion: 0.15
Nodes (24): testHandler, TestConfirmingEnrollmentUpgradesSession(), TestDeleteFactorBlockedByPolicyWhenLast(), TestElevateMethodsReflectsMFA(), TestElevatePasswordRejectedWhenMFAEnabled(), TestElevateRejectsOIDCUsers(), TestElevateWithTOTPAndRecoveryCode(), TestEnrollmentOnlySessionBlockedOutsideAllowlist() (+16 more)

### Community 69 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 70 - "User"
Cohesion: 0.13
Nodes (11): Handler, newOIDCUser(), readOIDCFlowCookie(), validHostPort(), instanceAdminRoleFor(), OIDCClaims, OIDCProvider, OIDCProvider (+3 more)

### Community 71 - "AuthedUser"
Cohesion: 0.12
Nodes (30): TestEmbeddedSPAWithoutFrontendBuild(), NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token() (+22 more)

### Community 72 - "NewEngine"
Cohesion: 0.14
Nodes (25): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+17 more)

### Community 73 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 74 - "api/auth/webauthn_test.go"
Cohesion: 0.19
Nodes (24): secondFactorInput, virtualAuthenticator, flowCookieFrom(), beginWebAuthnLogin(), credentialResponse(), finishWebAuthnLogin(), testHandler, newVirtualAuthenticator() (+16 more)

### Community 75 - "Config"
Cohesion: 0.11
Nodes (23): newLogger(), NewHandler(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote() (+15 more)

### Community 76 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 77 - "src/theme.ts"
Cohesion: 0.14
Nodes (25): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+17 more)

### Community 78 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 79 - "connector/connector.go"
Cohesion: 0.09
Nodes (12): TimeoutError, NewAuthError(), NewTimeoutError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap(), AuthError, CredentialRefresher, MalformedResponseError (+4 more)

### Community 80 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 81 - "motion"
Cohesion: 0.14
Nodes (22): motion, zustand, MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast (+14 more)

### Community 82 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 83 - "Handler"
Cohesion: 0.14
Nodes (14): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+6 more)

### Community 84 - "newTestHandler"
Cohesion: 0.14
Nodes (23): createDNSResolverConnector(), Handler, itoa(), TestConfigFieldsHandler(), TestConfigPushHandler(), TestStartStopHandler(), Handler, newTestHandler() (+15 more)

### Community 85 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 86 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 87 - "time.Time"
Cohesion: 0.16
Nodes (21): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift (+13 more)

### Community 88 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 89 - "Connector"
Cohesion: 0.16
Nodes (8): apiMessage(), controllerName(), countByKind(), statusError(), unavailable(), Connector, sectionFetch, session

### Community 90 - "log/slog.Logger"
Cohesion: 0.18
Nodes (16): Dispatcher, RunDeliveryRetries(), RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure() (+8 more)

### Community 91 - "logging.go"
Cohesion: 0.16
Nodes (17): loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing() (+9 more)

### Community 92 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 93 - "sshStdioConn"
Cohesion: 0.10
Nodes (12): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, bufio.ReadWriter, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session (+4 more)

### Community 94 - "Handler"
Cohesion: 0.18
Nodes (6): webAuthnFlow, webAuthnUser, Handler, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/go-webauthn/webauthn/webauthn.SessionData, github.com/google/uuid.UUID

### Community 95 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 96 - "NewClient"
Cohesion: 0.13
Nodes (19): isSafeMethod(), GuardedDialer(), IsDangerousIP(), clientTimeout(), IsSafeMethod(), NewClient(), NewTransport(), NoRedirect() (+11 more)

### Community 97 - "NewHTTPClient"
Cohesion: 0.13
Nodes (15): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+7 more)

### Community 98 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 99 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 100 - "httpx/retry_test.go"
Cohesion: 0.29
Nodes (16): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+8 more)

### Community 101 - "Deps"
Cohesion: 0.23
Nodes (18): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+10 more)

### Community 102 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 103 - "ContextWithUser"
Cohesion: 0.13
Nodes (18): TestConnectorStoreErrorPaths(), TestByServiceNoDocsYet(), TestGenerate(), TestGetUnknownIDFallsBackToServicePlaceholder(), TestTree(), withGrant(), NewHandler(), Handler (+10 more)

### Community 104 - "export_test.go"
Cohesion: 0.22
Nodes (16): Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile(), TestDocExportDefaultCronExprIsValid() (+8 more)

### Community 105 - "time.Duration"
Cohesion: 0.15
Nodes (8): retryable(), sleep(), Database, Server, time.Duration, RetryPolicy, retryTransport, PoolConfig

### Community 106 - "Hub"
Cohesion: 0.14
Nodes (7): Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, broadcastMsg, Client, Revalidator, ticket

### Community 107 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 108 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 109 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 110 - "Handler"
Cohesion: 0.19
Nodes (3): Handler, stripLogControlChars(), Handler

### Community 111 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 112 - ".Fetch"
Cohesion: 0.18
Nodes (8): ServiceDependency, WantsField(), environmentDependencies(), putMetadata(), unavailable(), networkDependencies(), Connector, dockerSectionSpec

### Community 113 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 114 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 115 - "templates_test.go"
Cohesion: 0.26
Nodes (15): templateBody, TestTemplateMutationRoleMatrix(), testApp, seedPreviewConnector(), seedTemplate(), TestTemplatesConcurrentUpdatesCreateDistinctVersions(), TestTemplatesPreviewAffectedConnectors(), TestTemplatesPreviewCapturesMissingSnapshot() (+7 more)

### Community 116 - "api/mcp_test.go"
Cohesion: 0.17
Nodes (8): go_pkg_github_com_mark3labs_mcp_go_client, go_pkg_github_com_mark3labs_mcp_go_client_transport, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server, go_pkg_github_com_wiselabz_wiselabz_internal_chat, changeSummary, connectorSummary, findingSummary

### Community 117 - "system/handlers_test.go"
Cohesion: 0.23
Nodes (15): Handler, newTestHandler(), TestDiagnostics(), TestExportAudit(), TestExportImportBackupRoundTrip(), TestGetBackupScheduleDefault(), TestGetRetentionSettingsDefault(), TestHealth() (+7 more)

### Community 118 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 119 - "DocRecord"
Cohesion: 0.18
Nodes (7): fetchAllDocs(), docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc(), scanDocSummary()

### Community 121 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 122 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 123 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 124 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 125 - "handlers_actions_test.go"
Cohesion: 0.32
Nodes (14): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+6 more)

### Community 126 - "Handler"
Cohesion: 0.24
Nodes (6): stepAuditDetail(), validTargetType(), validVerb(), Handler, runbookResponse, stepResponse

### Community 127 - "Store"
Cohesion: 0.16
Nodes (5): versionSections(), Store, TemplateSectionRecord, TemplateVersionRecord, TemplateVersionSection

### Community 128 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 129 - "registry.go"
Cohesion: 0.17
Nodes (12): LifecycleOp(), supportedLifecycleVerbs(), IsCredentialRefresherType(), ListSchemas(), SupportsLifecycleVerb(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec (+4 more)

### Community 130 - "doJSON"
Cohesion: 0.27
Nodes (11): doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser(), TestLogin() (+3 more)

### Community 131 - "gitTarget"
Cohesion: 0.23
Nodes (7): commitMessage(), commitResult, gitTarget, git.Repository, github.com/go-git/go-git/v5/plumbing.Hash, github.com/go-git/go-git/v5/plumbing/object.Signature, github.com/go-git/go-git/v5/plumbing.ReferenceName

### Community 132 - "gitFixture"
Cohesion: 0.32
Nodes (8): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, Exporter, github.com/go-git/go-git/v5/plumbing/object.Commit

### Community 133 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 134 - "New"
Cohesion: 0.36
Nodes (13): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+5 more)

### Community 135 - "config_cmd_test.go"
Cohesion: 0.23
Nodes (10): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), Schema(), schemaFor() (+2 more)

### Community 136 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 137 - "export.go"
Cohesion: 0.23
Nodes (10): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails(), fileName() (+2 more)

### Community 138 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 139 - "chat/chat.go"
Cohesion: 0.24
Nodes (11): cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips(), TestSplitSections() (+3 more)

### Community 140 - "createUser"
Cohesion: 0.28
Nodes (13): seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding(), TestSearchDocs(), seedFinding() (+5 more)

### Community 141 - "NotificationRecord"
Cohesion: 0.26
Nodes (4): Dispatcher, NotificationRecord, Store, scanNotification()

### Community 142 - "Contributing to WiseLabz"
Cohesion: 0.15
Nodes (13): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+5 more)

### Community 143 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 144 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 145 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 146 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 147 - "Store"
Cohesion: 0.23
Nodes (4): ChatConversationRecord, Store, ChatMessageRecord, DocSectionEmbeddingRecord

### Community 148 - "JobHealthRecord"
Cohesion: 0.29
Nodes (4): JobHealthRecord, Store, scanJobHealth(), fakeHealthStore

### Community 149 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 150 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 152 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 153 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 154 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 155 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 156 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 157 - "TestElevateOIDC"
Cohesion: 0.27
Nodes (9): mockElevateOIDCServer, beginElevate(), containsCode(), defaultClaims(), elevateOIDCTestSetup(), testHandler, newMockElevateOIDCServer(), TestElevateOIDC() (+1 more)

### Community 158 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 159 - ".SnapshotDiff"
Cohesion: 0.29
Nodes (6): decodeStoredSnapshot(), Handler, snapshotStoreError(), Cursor(), T, NextCursor()

### Community 160 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 161 - "diagnostics/diagnostics.go"
Cohesion: 0.38
Nodes (9): collectVersions(), AuthProviders, Bundle, Component, Health, OIDCProviderSummary, SanitizedConfig, SanitizedConnector (+1 more)

### Community 162 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 163 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 164 - "ComplianceRuleRecord"
Cohesion: 0.36
Nodes (4): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule()

### Community 165 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 166 - "scheduler/health_test.go"
Cohesion: 0.44
Nodes (6): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), fakeNotifier, notifyCall

### Community 167 - "BackupSchedule"
Cohesion: 0.31
Nodes (4): BackupSchedule, Store, scanBackupRun(), BackupRun

### Community 168 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 169 - "auth/handlers_test.go"
Cohesion: 0.25
Nodes (7): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups()

### Community 170 - "Engine"
Cohesion: 0.29
Nodes (4): NewHandler(), Engine, sync.Map, DocRegenerator

### Community 171 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 172 - "gitAuth"
Cohesion: 0.29
Nodes (7): gitAuth(), Exporter, installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), GitOptions, github.com/go-git/go-git/v5/plumbing/transport.AuthMethod

### Community 173 - "cursor_test.go"
Cohesion: 0.43
Nodes (6): DecodeCursor(), EncodeCursor(), TestCursorRequestModes(), TestCursorRoundTrip(), TestDecodeCursorRejectsGarbage(), TestNextCursorStopsOnShortPage()

### Community 174 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 175 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 176 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 177 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 178 - "lifecycleDeps"
Cohesion: 0.33
Nodes (6): newLifecycleManager(), context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, lifecycleDeps, lifecycleManager

### Community 179 - "session_test.go"
Cohesion: 0.38
Nodes (6): refreshCookie(), TestDeleteSession(), TestElevate(), TestElevateMethods(), TestLogout(), TestRefresh()

### Community 180 - "openapi_contract_test.go"
Cohesion: 0.48
Nodes (6): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

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

### Community 187 - "testHandler"
Cohesion: 0.33
Nodes (4): testHandler, Handler, spaHandler(), net/http.HandlerFunc

### Community 189 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 190 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 191 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 192 - "RequireConnectorRole"
Cohesion: 0.50
Nodes (4): ConnectorRoleChecker, RequireConnectorRole(), TestRequireConnectorRole(), TestRequireConnectorRoleCheckerError()

### Community 193 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 194 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 195 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 196 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 197 - "internal/auth/oidc.go"
Cohesion: 0.50
Nodes (3): extractGroups(), go_pkg_github_com_coreos_go_oidc_v3_oidc, go_pkg_golang_org_x_oauth2

### Community 198 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 202 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 203 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Knowledge Gaps
- **564 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+559 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1279 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `testing.T`, `gitFixture`, `NewChecker`, `net/http.Request`, `go_pkg_context`, `chat/chat.go`, `createUser`, `RunMigrations`, `Errorf`, `ConnectorRecord`, `Handler`, `Store`, `Get`, `.call`, `dispatcher_test.go`, `NewEngine`, `Manager`, `.call`, `ExportToFile`, `NewUser`, `Engine`, `NewStore`, `Dispatcher`, `lifecycleDeps`, `main`, `NewRegistry`, `response.go`, `testHandler`, `engine_maintenance_test.go`, `AuthedUser`, `NewEngine`, `rewritePlaceholders`, `Config`, `time.Time`, `log/slog.Logger`, `Deps`, `Handler`, `ContextWithUser`, `export_test.go`, `DocRecord`, `testApp`, `changes/handlers_test.go`, `Handler`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `Hub` connect `Hub` to `log/slog.Logger`, `NewChecker`, `NewEngine`, `net/http.Request`, `Engine`, `time.Duration`, `ws/ws_test.go`, `NewStore`, `go_pkg_net_http`, `Errorf`, `Dispatcher`, `lifecycleDeps`, `main`, `NewRegistry`, `testApp`, `Get`, `changes/handlers_test.go`, `dispatcher_test.go`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `gitFixture` to `testing.T`, `context.Context`, `export_test.go`, `git.go`, `log/slog.Logger`, `Store`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _564 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.01964482885535517 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.029753694581280788 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019068500431864994 - nodes in this community are weakly interconnected._