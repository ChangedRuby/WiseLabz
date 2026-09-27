# Graph Report - feat-runbook-execution-with-approvals  (2026-09-27)

## Corpus Check
- 865 files · ~532,253 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: (none) 10, .toml 2, .tmpl 2)

## Summary
- 6405 nodes · 20369 edges · 207 communities (191 shown, 16 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1650 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f16b8b03`
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
- TemplateEditorPage.tsx
- @tanstack/react-query
- icons.tsx
- MapTransportError
- RunMigrations
- go_pkg_net_http
- Errorf
- ServiceSnapshot
- ConnectorRecord
- ConnectorEditPage.tsx
- App.tsx
- router.go
- ProfilePage.tsx
- ChatPage.tsx
- Store
- UsersPage.tsx
- net/http.ResponseWriter
- NewEngine
- dispatcher_test.go
- SnapshotEntity
- Connector
- Config
- package.json
- ExportToFile
- NewUser
- snapshotdiff.go
- fixtures.ts
- WiseLabz — Architecture & Technical Decisions
- routerDeps
- ServicesPage.tsx
- docker_test.go
- truenas_test.go
- home_assistant/tables.go
- dependencies
- HashToken
- NewStore
- portainer/tables.go
- home_assistant_test.go
- NewMalformedResponseError
- Dispatcher
- ServiceDetailPage.tsx
- adguardhome/tables.go
- main
- pagination_contract_test.go
- NewRegistry
- connector_permission.go
- response.go
- config/validate_test.go
- traefik/tables.go
- net/http.Client
- Store
- traefik_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- go_pkg_strings
- compliance/engine.go
- middleware_test.go
- newTestHandler
- settings.mock.ts
- ErrorWithDetails
- AuthedUser
- NewEngine
- rewritePlaceholders
- runbooks_test.go
- Config
- Service
- src/theme.ts
- handlers.ts
- connector/connector.go
- unifi_test.go
- AppearancePage.tsx
- timeline.ts
- internal/auth/mfa.go
- newTestHandler
- WiseLabz — Design Contract
- devDependencies
- time.Time
- portainer_test.go
- Connector
- retention/retention_test.go
- backup/main.go
- adguardhome_test.go
- sshStdioConn
- Handler
- NewService
- New
- NewHTTPClient
- Runner
- Connector
- httpx/retry_test.go
- Deps
- Handler
- newTestHandler
- git.go
- Checker
- Hub
- ws/ws_test.go
- ws.ts
- compilerOptions
- Handler
- Register
- .Fetch
- RunbookRecord
- docdiffmodel.ts
- Compare
- config_test.go
- MarshalConnectorConfig
- all.go
- DocRecord
- ReportsPage.tsx
- compilerOptions
- HashPassword
- handlers_contract_test.go
- changes/handlers_test.go
- keyset_test.go
- Handler
- nilToStr
- vectorCache
- registry.go
- Handler
- git_internal_test.go
- gitFixture
- render_test.go
- NewRouter
- config_cmd_test.go
- .CreateShareLink
- TestComplianceRuleValidation
- connectors_health_test.go
- chat/chat.go
- newTestHarness
- NotificationRecord
- Contributing to WiseLabz
- scripts
- templatefuncs.go
- IsSecureRequest
- Store
- Handler
- newTestHandler
- Decision
- main.tsx
- Handler
- Decision
- WiseLabz Connector Guide
- Product
- .call
- store/backup_test.go
- Connector
- connectors_maintenance_test.go
- .SnapshotDiff
- .call
- diagnostics/diagnostics.go
- Changelog
- mockServiceWorker.js
- ComplianceRuleRecord
- connectors_hardening_test.go
- api/docs_test.go
- dashboard/handlers_test.go
- release-please-config.json
- api/auth/oidc.go
- ShareLink
- net/http.Handler
- snapshotResponse
- .applyChannelSecrets
- ComputeWindow
- Cache
- Step by step
- WiseLabz
- newTestLifecycle
- cloudflare/attributes_test.go
- openapi_contract_test.go
- AuditRecorder
- computeNextRun
- Contributor Covenant Code of Conduct
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- fields_test.go
- engine_maintenance_test.go
- Mermaid.tsx
- Security Policy
- Enforcement Guidelines
- MfaEnrollDialog
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- RetentionSettings
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
1. `newTestApp()` - 229 edges
2. `Errorf()` - 183 edges
3. `newDocTestStore()` - 143 edges
4. `Store` - 141 edges
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

## Communities (207 total, 16 thin omitted)

### Community 0 - "newTestApp"
Cohesion: 0.02
Nodes (171): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+163 more)

### Community 1 - "react"
Cohesion: 0.03
Nodes (127): react, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules (+119 more)

### Community 2 - "testing.T"
Cohesion: 0.02
Nodes (143): cursorPage, TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), createKey() (+135 more)

### Community 3 - "newDocTestStore"
Cohesion: 0.02
Nodes (146): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+138 more)

### Community 4 - "context.Context"
Cohesion: 0.02
Nodes (38): fakeStatusChecker, sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, existingIDs(), Store, SnapshotRecord, Store (+30 more)

### Community 5 - "DashboardPage.tsx"
Cohesion: 0.03
Nodes (93): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress`, 3. `sync.complete`, 4. `change.detected` (+85 more)

### Community 6 - "go_pkg_testing"
Cohesion: 0.05
Nodes (27): TestLoggablePathMasksShareTokenUnderV1(), go_pkg_crypto_rsa, go_pkg_crypto_tls, go_pkg_encoding_json, go_pkg_github_com_go_jose_go_jose_v4, go_pkg_github_com_gorilla_websocket, go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_apitest (+19 more)

### Community 7 - "NewChecker"
Cohesion: 0.17
Nodes (35): NewChecker(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves(), TestCheckFailingDetectsAndAutoResolves() (+27 more)

### Community 8 - "net/http.Request"
Cohesion: 0.08
Nodes (23): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+15 more)

### Community 9 - "SystemPage.tsx"
Cohesion: 0.05
Nodes (39): AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setRefreshHandler(), server (+31 more)

### Community 10 - "go_pkg_context"
Cohesion: 0.07
Nodes (18): StatusError, dashboardLayout, contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt (+10 more)

### Community 11 - "TemplateEditorPage.tsx"
Cohesion: 0.04
Nodes (51): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatesquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid (+43 more)

### Community 12 - "@tanstack/react-query"
Cohesion: 0.02
Nodes (91): Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`), axios (+83 more)

### Community 13 - "icons.tsx"
Cohesion: 0.04
Nodes (98): Frontend shell & theme (decided 2026-06), i18next, react-i18next, zustand, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey (+90 more)

### Community 14 - "MapTransportError"
Cohesion: 0.04
Nodes (29): Connector, NewServiceUnavailableError(), setHeaders(), TestValidateCustomURL(), tryParseEntities(), validateCustomURL(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable() (+21 more)

### Community 15 - "RunMigrations"
Cohesion: 0.09
Nodes (43): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), newPostgresTestStore(), GetMigrationStatus(), newMigrator(), collectColumns() (+35 more)

### Community 16 - "go_pkg_net_http"
Cohesion: 0.05
Nodes (45): bulkSnoozeItemResult, bulkSnoozeRequest, contextKey, elevationError, changePromptData(), stripPromptTags(), truncateUTF8(), versionSections() (+37 more)

### Community 17 - "Errorf"
Cohesion: 0.05
Nodes (27): Handler, newToken(), sanitize(), Handler, diffToSpec(), NewHandler(), Handler, Handler (+19 more)

### Community 18 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (24): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Sanitize(), TestSanitize(), changePatternID(), Engine, markError() (+16 more)

### Community 19 - "ConnectorRecord"
Cohesion: 0.05
Nodes (39): enableFakeEmbedding(), actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), ConnectorRecord, Store (+31 more)

### Community 20 - "ConnectorEditPage.tsx"
Cohesion: 0.07
Nodes (28): RFC-3339, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_putconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsschema (+20 more)

### Community 21 - "App.tsx"
Cohesion: 0.05
Nodes (57): react-error-boundary, react-router-dom, sonner, setAccessToken(), setMfaEnrollmentRequiredHandler(), web_src_api_generated_auth_auth_postauthlogin, web_src_api_generated_auth_auth_postauthloginmfa, web_src_api_generated_auth_auth_postauthlogout (+49 more)

### Community 22 - "router.go"
Cohesion: 0.14
Nodes (22): go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth, go_pkg_github_com_wiselabz_wiselabz_internal_api_changes, go_pkg_github_com_wiselabz_wiselabz_internal_api_chat, go_pkg_github_com_wiselabz_wiselabz_internal_api_compliance, go_pkg_github_com_wiselabz_wiselabz_internal_api_connectors (+14 more)

### Community 23 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (51): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+43 more)

### Community 24 - "ChatPage.tsx"
Cohesion: 0.06
Nodes (34): web_src_api_generated_changes_changes_getgetchangeschangeidquerykey, web_src_api_generated_changes_changes_postchangeschangeidack, web_src_api_generated_changes_changes_postchangeschangeidaiupdate, web_src_api_generated_changes_changes_postchangeschangeiddismiss, web_src_api_generated_changes_changes_postchangeschangeidexplain, web_src_api_generated_changes_changes_usegetchangeschangeid, web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey (+26 more)

### Community 25 - "Store"
Cohesion: 0.06
Nodes (18): seedAlert(), changeServiceIDs(), seedChange(), Store, Store, placeholders(), scanBackupRun(), changeFilterClause() (+10 more)

### Community 26 - "UsersPage.tsx"
Cohesion: 0.14
Nodes (23): customInstance(), web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa, web_src_api_generated_users_users_postusersuseridresetpassword, web_src_api_generated_users_users_usegetusers, ConnectorGrant (+15 more)

### Community 27 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (26): Handler, isWritableField(), validateConfigPushRequest(), applyConnectorScalarUpdates(), configRequestField(), Handler, parseScheduleUpdates(), validateConnectorConfig() (+18 more)

### Community 28 - "NewEngine"
Cohesion: 0.09
Nodes (44): NewHandler(), entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine (+36 more)

### Community 29 - "dispatcher_test.go"
Cohesion: 0.18
Nodes (50): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher (+42 more)

### Community 30 - "SnapshotEntity"
Cohesion: 0.09
Nodes (52): SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable(), TestBuildPeerTableAttributes(), TestBuildPolicyTableAttributes(), buildPeerTable(), buildPolicyTable(), TestBuildInterfaceTableAttributes() (+44 more)

### Community 31 - "Connector"
Cohesion: 0.07
Nodes (11): init(), ConfigField, Connector, buildRouteTable(), buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment() (+3 more)

### Community 32 - "Config"
Cohesion: 0.17
Nodes (11): Config, cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), migrationFiles(), TestMigrationVersionParity() (+3 more)

### Community 33 - "package.json"
Cohesion: 0.04
Nodes (48): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+40 more)

### Community 34 - "ExportToFile"
Cohesion: 0.10
Nodes (43): ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions() (+35 more)

### Community 35 - "NewUser"
Cohesion: 0.17
Nodes (41): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON(), TestGetLockNoneHeld() (+33 more)

### Community 36 - "snapshotdiff.go"
Cohesion: 0.11
Nodes (28): ServiceDependency, environmentDependencies(), poolDependencies(), networkDependencies(), BuildSnapshotDiff(), CompareDependencies(), CompareEntities(), entityKey() (+20 more)

### Community 37 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 38 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.04
Nodes (46): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+38 more)

### Community 39 - "routerDeps"
Cohesion: 0.12
Nodes (28): routerDeps, chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router (+20 more)

### Community 40 - "ServicesPage.tsx"
Cohesion: 0.04
Nodes (66): Frontend, match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve (+58 more)

### Community 41 - "docker_test.go"
Cohesion: 0.05
Nodes (48): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+40 more)

### Community 42 - "truenas_test.go"
Cohesion: 0.08
Nodes (42): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+34 more)

### Community 43 - "home_assistant/tables.go"
Cohesion: 0.10
Nodes (38): jsonType(), TestAttributeCatalogCoversEmittedKeys(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations(), buildOverview() (+30 more)

### Community 44 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 45 - "HashToken"
Cohesion: 0.07
Nodes (32): ProviderConfig, factorJSON(), Handler, Handler, Handler, primaryProviderConfig(), Handler, oidcProviderJSON() (+24 more)

### Community 46 - "NewStore"
Cohesion: 0.10
Nodes (40): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+32 more)

### Community 47 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 48 - "home_assistant_test.go"
Cohesion: 0.13
Nodes (29): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+21 more)

### Community 49 - "NewMalformedResponseError"
Cohesion: 0.07
Nodes (48): SnapshotSection, NewMalformedResponseError(), TestBuildHostsTableAttributes(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell() (+40 more)

### Community 50 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 51 - "ServiceDetailPage.tsx"
Cohesion: 0.05
Nodes (47): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+39 more)

### Community 52 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 53 - "main"
Cohesion: 0.08
Nodes (30): main(), newLogger(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder() (+22 more)

### Community 54 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 55 - "NewRegistry"
Cohesion: 0.14
Nodes (24): Provider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+16 more)

### Community 56 - "connector_permission.go"
Cohesion: 0.10
Nodes (20): APIKeyRestriction, testAPIKeyChecker, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), TestClampConnectorRole(), treatAsSafeFromContext() (+12 more)

### Community 57 - "response.go"
Cohesion: 0.10
Nodes (18): Handler, Handler, TestWritePaginatedOmitsNextCursor(), Error(), DataPaginatedResponse, HandleStoreError(), intQuery(), JSON() (+10 more)

### Community 58 - "config/validate_test.go"
Cohesion: 0.18
Nodes (11): Config, mask(), redactDSN(), redactKVPassword(), Config, TestEveryKeyEnvOverridable(), TestRedactDSN(), TestRedacted() (+3 more)

### Community 59 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 60 - "net/http.Client"
Cohesion: 0.06
Nodes (18): claudeProvider, ollamaEmbedder, openAICompatibleProvider, openAIEmbedder, StubProvider, testProvider, SuggestChunk, SuggestRequest (+10 more)

### Community 61 - "Store"
Cohesion: 0.16
Nodes (30): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+22 more)

### Community 62 - "traefik_test.go"
Cohesion: 0.11
Nodes (30): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+22 more)

### Community 63 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 64 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 65 - "go_pkg_strings"
Cohesion: 0.04
Nodes (46): Schema(), schemaFor(), TestSchemaMatchesConfig(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEmailMessage(), sendNtfyChannel(), sendSMTPChannel() (+38 more)

### Community 66 - "compliance/engine.go"
Cohesion: 0.11
Nodes (31): catalog(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+23 more)

### Community 67 - "middleware_test.go"
Cohesion: 0.16
Nodes (15): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser() (+7 more)

### Community 68 - "newTestHandler"
Cohesion: 0.06
Nodes (73): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, NewHandler(), doJSON(), testHandler, req(), TestChangePassword() (+65 more)

### Community 69 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 70 - "ErrorWithDetails"
Cohesion: 0.08
Nodes (21): sanitizeUser(), setRefreshCookie(), Handler, Handler, Handler, newOIDCUser(), randomOIDCToken(), validHostPort() (+13 more)

### Community 71 - "AuthedUser"
Cohesion: 0.11
Nodes (31): NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth() (+23 more)

### Community 72 - "NewEngine"
Cohesion: 0.18
Nodes (23): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+15 more)

### Community 73 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 74 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 75 - "Config"
Cohesion: 0.09
Nodes (25): NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings, AuthSettings (+17 more)

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
Cohesion: 0.08
Nodes (14): TimeoutError, GuardedDialer(), IsDangerousIP(), NewTimeoutError(), newWebhookClient(), AuthError, CredentialRefresher, MalformedResponseError (+6 more)

### Community 80 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 81 - "AppearancePage.tsx"
Cohesion: 0.15
Nodes (20): MotionProvider(), AppearancePage(), ChoiceGroup(), AppearanceState, apply(), Contrast, css(), DEFAULTS (+12 more)

### Community 82 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 83 - "internal/auth/mfa.go"
Cohesion: 0.38
Nodes (5): GenerateRecoveryCodes(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), go_pkg_github_com_pquerna_otp, go_pkg_github_com_pquerna_otp_totp

### Community 84 - "newTestHandler"
Cohesion: 0.09
Nodes (49): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+41 more)

### Community 85 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 86 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 87 - "time.Time"
Cohesion: 0.15
Nodes (22): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), golang.org/x/time/rate.Limiter, time.Time, visitor, ChangeEntry (+14 more)

### Community 88 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 89 - "Connector"
Cohesion: 0.11
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), TestTypedErrorsWrapAndUnwrap(), Connector, apiMessage(), controllerName(), countByKind(), statusError() (+3 more)

### Community 90 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 91 - "backup/main.go"
Cohesion: 0.13
Nodes (20): main(), usage(), loggablePath(), loggableQuery(), Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken() (+12 more)

### Community 92 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 93 - "sshStdioConn"
Cohesion: 0.09
Nodes (13): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, bufio.ReadWriter, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session (+5 more)

### Community 94 - "Handler"
Cohesion: 0.18
Nodes (6): webAuthnFlow, webAuthnUser, Handler, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/go-webauthn/webauthn/webauthn.SessionData, github.com/google/uuid.UUID

### Community 95 - "NewService"
Cohesion: 0.10
Nodes (28): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs() (+20 more)

### Community 96 - "New"
Cohesion: 0.10
Nodes (24): confirm(), formatCounts(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag() (+16 more)

### Community 97 - "NewHTTPClient"
Cohesion: 0.13
Nodes (15): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+7 more)

### Community 98 - "Runner"
Cohesion: 0.07
Nodes (31): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+23 more)

### Community 99 - "Connector"
Cohesion: 0.16
Nodes (5): buildGatewayTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 100 - "httpx/retry_test.go"
Cohesion: 0.20
Nodes (20): retryable(), RetryTransport(), sleep(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry() (+12 more)

### Community 101 - "Deps"
Cohesion: 0.31
Nodes (16): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), findingConnectorIDs(), registerListFindings() (+8 more)

### Community 102 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 103 - "newTestHandler"
Cohesion: 0.24
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 104 - "git.go"
Cohesion: 0.10
Nodes (18): IsGeneratedName(), pruneStale(), TestIsGeneratedName(), go_pkg_github_com_getkin_kin_openapi_openapi3, go_pkg_github_com_getkin_kin_openapi_openapi3filter, go_pkg_github_com_getkin_kin_openapi_routers, go_pkg_github_com_getkin_kin_openapi_routers_gorillamux, go_pkg_github_com_go_git_go_git_v5 (+10 more)

### Community 105 - "Checker"
Cohesion: 0.20
Nodes (7): complianceRule(), Checker, RunStaleSweepOnce(), QualityFindingRecord, scanQualityFinding(), FindingNotifier, RotationConfig

### Community 106 - "Hub"
Cohesion: 0.08
Nodes (13): NewHandler(), Engine, Hub, github.com/gorilla/websocket.Conn, github.com/gorilla/websocket.Upgrader, sync.Map, AlertNotifier, DocRegenerator (+5 more)

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
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 111 - "Register"
Cohesion: 0.21
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 112 - ".Fetch"
Cohesion: 0.15
Nodes (8): WantsField(), Connector, putMetadata(), unavailable(), agentEnabled(), Connector, Connector, dockerSectionSpec

### Community 113 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 114 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 115 - "Compare"
Cohesion: 0.15
Nodes (18): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+10 more)

### Community 116 - "config_test.go"
Cohesion: 0.14
Nodes (18): Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH(), TestLoadFromYAML() (+10 more)

### Community 117 - "MarshalConnectorConfig"
Cohesion: 0.18
Nodes (16): IsSecretFieldType(), MarshalConnectorConfig(), ParseConnectorConfig(), SecretFieldsChanged(), init(), TestConnectorRotationFieldsRoundTrip(), TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(), TestSecretFieldsChangedFalseOnRenameOnly() (+8 more)

### Community 118 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 119 - "DocRecord"
Cohesion: 0.14
Nodes (10): fetchAllDocs(), fileName(), slugify(), docSearchWhere(), escapeLike(), DocRecord, Store, scanDoc() (+2 more)

### Community 120 - "ReportsPage.tsx"
Cohesion: 0.11
Nodes (19): web_src_api_generated_reports_reports, web_src_api_generated_reports_reports_deletereportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_getgetreportsdefinitionsquerykey, web_src_api_generated_reports_reports_getgetreportsquerykey, web_src_api_generated_reports_reports_postreportsdefinitions, web_src_api_generated_reports_reports_postreportsdefinitionsreportdefinitionidrun, web_src_api_generated_reports_reports_putreportsdefinitionsreportdefinitionid, web_src_api_generated_reports_reports_usegetreports (+11 more)

### Community 121 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 122 - "HashPassword"
Cohesion: 0.10
Nodes (19): mustHashDummyPassword(), instanceAdminRoleFor(), testApp, TestDashboardAdminDefaultPermissionGate(), TestDashboardResetRestoresAdminDefault(), testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable() (+11 more)

### Community 123 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 124 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 125 - "keyset_test.go"
Cohesion: 0.19
Nodes (16): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern(), assertSameSet(), Store (+8 more)

### Community 126 - "Handler"
Cohesion: 0.31
Nodes (4): NewHandler(), Handler, runbookResponse, stepResponse

### Community 127 - "nilToStr"
Cohesion: 0.07
Nodes (14): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+6 more)

### Community 128 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 129 - "registry.go"
Cohesion: 0.19
Nodes (10): supportedLifecycleVerbs(), IsCredentialRefresherType(), ListSchemas(), TestIsCredentialRefresherType(), TestRegisterStubRoundTrips(), AttributeSpec, ConfigValidationError, Factory (+2 more)

### Community 130 - "Handler"
Cohesion: 0.27
Nodes (4): updateUserRequest, Handler, writeUserWriteError(), NoContent()

### Community 131 - "git_internal_test.go"
Cohesion: 0.09
Nodes (23): commitMessage(), gitAuth(), Exporter, installHTTPS(), SetBeforePushForTest(), TestCommitMessage(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken() (+15 more)

### Community 132 - "gitFixture"
Cohesion: 0.18
Nodes (18): Exporter, NewExporter(), RunExportOnce(), newTestStore(), readFile(), TestExportAllEmptyDirName(), TestExportAllKeepsOperatorFiles(), TestExportAllPrunesStaleFilesOnRerun() (+10 more)

### Community 133 - "render_test.go"
Cohesion: 0.20
Nodes (17): connectorFilter(), RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden() (+9 more)

### Community 134 - "NewRouter"
Cohesion: 0.19
Nodes (11): TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), chi.Router, NewRouter(), wsRoleLabel() (+3 more)

### Community 135 - "config_cmd_test.go"
Cohesion: 0.36
Nodes (7): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), io.Writer

### Community 136 - ".CreateShareLink"
Cohesion: 0.23
Nodes (6): contextWithShareLink(), Handler, newShareToken(), shareLinkFromContext(), shareLinkNode, shareLinkScope

### Community 137 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 138 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 139 - "chat/chat.go"
Cohesion: 0.14
Nodes (18): buildPrompt(), TestBuildPrompt(), Handler, cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections() (+10 more)

### Community 140 - "newTestHarness"
Cohesion: 0.27
Nodes (14): TestListAttentionItems(), TestListChanges(), TestListConnectors(), TestSearchDocs(), seedFinding(), TestListFindings(), createConnector(), createUser() (+6 more)

### Community 141 - "NotificationRecord"
Cohesion: 0.13
Nodes (10): Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep(), NotificationRecord, Store (+2 more)

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

### Community 147 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 148 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

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

### Community 156 - "store/backup_test.go"
Cohesion: 0.32
Nodes (11): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+3 more)

### Community 158 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 159 - ".SnapshotDiff"
Cohesion: 0.21
Nodes (12): decodeStoredSnapshot(), Handler, snapshotStoreError(), Cursor(), DecodeCursor(), EncodeCursor(), T, NextCursor() (+4 more)

### Community 160 - ".call"
Cohesion: 0.33
Nodes (8): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), spaHandler(), fixture, net/http.HandlerFunc

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
Cohesion: 0.29
Nodes (7): testApp, TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 166 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 167 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 168 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 169 - "api/auth/oidc.go"
Cohesion: 0.27
Nodes (8): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), go_pkg_crypto_subtle

### Community 171 - "net/http.Handler"
Cohesion: 0.14
Nodes (13): ConnectorRoleChecker, PermissionChecker, TestRateLimit(), RateLimit(), SecurityHeaders(), TestSecurityHeaders(), RequireConnectorRole(), RequirePermission() (+5 more)

### Community 172 - "snapshotResponse"
Cohesion: 0.48
Nodes (7): Handler, snapshotFixture(), snapshotRequest(), snapshotResponse(), TestSnapshotDiffValidationOwnershipAndAudit(), TestSnapshotOwnershipAndFullShape(), TestSnapshotsViewerAndCursor()

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

### Community 178 - "newTestLifecycle"
Cohesion: 0.18
Nodes (11): newLifecycleManager(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, sync/atomic.Bool (+3 more)

### Community 179 - "cloudflare/attributes_test.go"
Cohesion: 0.33
Nodes (5): TestAttributeCatalogCoversEmittedKeys(), TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable()

### Community 180 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

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

### Community 189 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 190 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 191 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 193 - "Enforcement Guidelines"
Cohesion: 0.40
Nodes (5): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Enforcement Guidelines

### Community 194 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 196 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

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
- **568 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+563 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1285 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `Handler`, `gitFixture`, `render_test.go`, `NewChecker`, `net/http.Request`, `go_pkg_context`, `chat/chat.go`, `newTestHarness`, `RunMigrations`, `Errorf`, `ServiceSnapshot`, `ConnectorRecord`, `Handler`, `Handler`, `Store`, `net/http.ResponseWriter`, `.call`, `dispatcher_test.go`, `NewEngine`, `store/backup_test.go`, `Config`, `.call`, `ExportToFile`, `NewUser`, `NewStore`, `Dispatcher`, `newTestLifecycle`, `main`, `NewRegistry`, `response.go`, `engine_maintenance_test.go`, `newTestHandler`, `AuthedUser`, `NewEngine`, `rewritePlaceholders`, `time.Time`, `retention/retention_test.go`, `New`, `Deps`, `Handler`, `Checker`, `Hub`, `DocRecord`, `HashPassword`, `changes/handlers_test.go`, `Handler`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `Hub` connect `Hub` to `Config`, `NewChecker`, `NewEngine`, `Checker`, `net/http.Request`, `Config`, `ws/ws_test.go`, `NotificationRecord`, `NewStore`, `go_pkg_net_http`, `Errorf`, `Dispatcher`, `newTestLifecycle`, `NewRegistry`, `HashPassword`, `net/http.ResponseWriter`, `changes/handlers_test.go`, `dispatcher_test.go`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Why does `DecodeKey()` connect `HashToken` to `config/validate_test.go`, `MarshalConnectorConfig`, `Handler`, `.applyChannelSecrets`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _568 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021560940841054883 - nodes in this community are weakly interconnected._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.027227131597802722 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022830254444057164 - nodes in this community are weakly interconnected._