# Graph Report - agent-ad13b247f7307abdb  (2026-09-29)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 6569 nodes · 20826 edges · 238 communities (217 shown, 21 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1671 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7328fce2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- react
- newTestApp
- context.Context
- testing.T
- @tanstack/react-query
- DashboardPage.tsx
- net/http.Request
- cn
- go_pkg_context
- newTestHandler
- net/http.Client
- go_pkg_testing
- icons.tsx
- go_pkg_net_http
- auth/handlers_test.go
- New
- lists.go
- ServiceSnapshot
- templates.fixtures.ts
- ProfilePage.tsx
- git_test.go
- DBTX
- net/http.ResponseWriter
- ServiceDetailPage.tsx
- App.tsx
- ErrorWithDetails
- package.json
- docker_test.go
- Connector
- NewEngine
- SystemPage.tsx
- RunMigrations
- go_pkg_strings
- newDocTestStore
- fixtures.ts
- UsersPage.tsx
- truenas/tables.go
- home_assistant/tables.go
- Compare
- dispatcher_test.go
- ConnectorRecord
- mountAPIRoutes
- DecodeKey
- NotificationCenter.tsx
- go_pkg_os
- nilToStr
- dependencies
- NewUser
- ConnectorForm.tsx
- home_assistant_test.go
- portainer/tables.go
- response.go
- RulesPage.tsx
- NewChecker
- Saved Views
- adguardhome/tables.go
- Sanitize
- src/theme.ts
- traefik/tables.go
- NewStore
- unifi/tables_test.go
- main
- Store
- connector_permission.go
- Connector
- traefik_test.go
- middleware_test.go
- .UpdateAuthConfig
- config_test.go
- Config
- settings.mock.ts
- SuggestRequest
- rewritePlaceholders
- handlers.ts
- runbooks_test.go
- Service
- NewMalformedResponseError
- createTestConnector
- newTestHandler
- GetTypeSchema
- unifi_test.go
- timeline.ts
- newTestHandler
- router.go
- HashToken
- Store
- Manager
- WiseLabz — Design Contract
- devDependencies
- notifications/handlers_test.go
- backup/backup.go
- store/mfa_test.go
- portainer_test.go
- Handler
- doc_test.go
- NewHTTPClient
- adguardhome_test.go
- httpx/retry_test.go
- webAuthnUser
- lifecycleManager
- NewService
- SuggestWithFallback
- net/http.Handler
- validate.go
- NewEngine
- Connector
- sshStdioConn
- ws.ts
- Handler
- ExportToFile
- truenas_test.go
- diagnostics/diagnostics.go
- ws/ws_test.go
- compilerOptions
- Store
- Register
- newTestHarness
- Get
- RunbookRecord
- docdiffmodel.ts
- .batchDelete
- vectorCache
- all.go
- fakeDocRegenerator
- WiseLabz — Architecture & Technical Decisions
- compilerOptions
- api/auth/webauthn_test.go
- New
- handlers_contract_test.go
- Handler
- logging_test.go
- Deps
- time.Time
- connector_permission_test.go
- newTestLifecycle
- pagination_contract_test.go
- Handler
- render_test.go
- Contributing to WiseLabz
- 0004 — PostgreSQL leader election for background workers
- newTestHandler
- api/changes_test.go
- chat/chat.go
- connectors_health_test.go
- NewRegistry
- .CreateShareLink
- newTestHandler
- net/http.Response
- ListSchemas
- connectors_maintenance_test.go
- maintenance_test.go
- Engine
- scripts
- api/docs_test.go
- templatefuncs.go
- IsSecureRequest
- store/backup_test.go
- Contributor Covenant Code of Conduct
- Decision
- Decision
- Backend test performance
- main.tsx
- Elector
- retention/retention_test.go
- keyset_test.go
- time.Duration
- runbook_test.go
- Decision
- WiseLabz Connector Guide
- Product
- .call
- EmbedRegistry
- .call
- ReportData
- seedQualityConnector
- RequirePermission
- Changelog
- test-shards.sh
- mockServiceWorker.js
- apikey_scopes_test.go
- rowScanner
- openapi_contract_test.go
- BACKUP.md
- release-please-config.json
- @vitejs/plugin-react
- TestComplianceRuleValidation
- RetentionSettings
- dashboard/handlers_test.go
- RateLimit
- ComputeWindow
- notification_delivery_test.go
- AuditRecorder
- Cache
- Decision
- Step by step
- WiseLabz
- newHandler
- dockerSSHAddr
- change_pattern_test.go
- log/slog.Logger
- fakeQualityChecker
- computeNextRun
- Audit Trail
- Configuration & Documentation Backup (Export/Import)
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- AlertNotifier
- walkCursorPages
- SnapshotEntity
- .GetConnectorUptime
- Store
- engine_maintenance_test.go
- Backup Recovery: What Comes Back, and What Doesn't
- WiseLabz — Deployment Guide
- Scheduled Doc Export
- QualityChecker
- Security Policy
- connector/connector.go
- seedScopeFixture
- golden_snapshot_test.go
- Authentication design
- Development workflow
- MfaEnrollDialog
- compose-smoke.sh
- timeoutError
- MISSING — deferred & future frontend features
- Fixture reuse and lifecycle tests (#405–#407)
- fields_test.go
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
6. `react` - 77 edges
7. `SnapshotEntity` - 77 edges
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

## Communities (238 total, 21 thin omitted)

### Community 0 - "react"
Cohesion: 0.03
Nodes (139): RFC-3339, Frontend, react, react-i18next, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze, web_src_api_generated_alerts_alerts_postalertsbulksnooze (+131 more)

### Community 1 - "newTestApp"
Cohesion: 0.02
Nodes (166): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+158 more)

### Community 2 - "context.Context"
Cohesion: 0.03
Nodes (35): fakeStatusChecker, sanitizeSessions(), MFAEnrollOnlyFromContext(), Connector, existingIDs(), Store, SnapshotRecord, Store (+27 more)

### Community 3 - "testing.T"
Cohesion: 0.02
Nodes (138): TestClaudeSuggest(), TestClaudeSuggestDefaultMaxTokens(), TestClaudeSuggestErrors(), TestClaudeSuggestMultipleContentBlocks(), TestOpenAICompatibleSuggest(), TestOpenAICompatibleSuggestErrors(), TestOIDCRedirectURL(), TestFindOIDCProvider() (+130 more)

### Community 4 - "@tanstack/react-query"
Cohesion: 0.03
Nodes (76): i18next, msw, @tanstack/react-query, @testing-library/react, vitest, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_changes_changes (+68 more)

### Community 5 - "DashboardPage.tsx"
Cohesion: 0.05
Nodes (56): 4. `change.detected`, 5. `alert.created`, 8. `doc.generated`, web_src_api_generated_alerts_alerts_usegetalerts, web_src_api_generated_dashboard_dashboard_getdashboardlayout, web_src_api_generated_dashboard_dashboard_getdashboardlayoutadmindefault, web_src_api_generated_dashboard_dashboard_getgetdashboardlayoutadmindefaultquerykey, web_src_api_generated_dashboard_dashboard_postdashboardlayoutreset (+48 more)

### Community 6 - "net/http.Request"
Cohesion: 0.05
Nodes (32): Handler, Handler, newToken(), sanitize(), Handler, diffToSpec(), Handler, Handler (+24 more)

### Community 7 - "cn"
Cohesion: 0.03
Nodes (106): zustand, web_src_api_generated_notifications_notifications_usegetnotificationsdeliveries, web_src_api_generated_settings_settings_getgetaiconfigfallbackprovidersquerykey, web_src_api_generated_settings_settings_getgetaiconfigquerykey, web_src_api_generated_settings_settings_getgetauthconfigquerykey, web_src_api_generated_settings_settings_getgetnotificationsconfigquerykey, web_src_api_generated_settings_settings_postaiconfigtest, web_src_api_generated_settings_settings_postnotificationsconfigtest (+98 more)

### Community 8 - "go_pkg_context"
Cohesion: 0.07
Nodes (19): StatusError, dashboardLayout, versionSections(), ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold(), TemplateVersionSection, contains() (+11 more)

### Community 9 - "newTestHandler"
Cohesion: 0.08
Nodes (47): mockElevateOIDCServer, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys(), TestCreateUser(), TestDeleteUser() (+39 more)

### Community 10 - "net/http.Client"
Cohesion: 0.03
Nodes (31): Connector, ollamaEmbedder, openAIEmbedder, NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL(), tryParseEntities() (+23 more)

### Community 11 - "go_pkg_testing"
Cohesion: 0.05
Nodes (17): Schema(), schemaFor(), TestSchemaMatchesConfig(), IsTimeout(), TestIsTimeout(), go_pkg_crypto_rsa, go_pkg_crypto_tls, go_pkg_github_com_go_jose_go_jose_v4 (+9 more)

### Community 12 - "icons.tsx"
Cohesion: 0.05
Nodes (64): match-sorter, motion, @radix-ui/react-popover, web_src_api_generated_attention_attention, web_src_api_generated_attention_attention_usegetattention, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow, web_src_api_generated_connectors_connectors_getgetconnectorsmaintenancewindowsquerykey (+56 more)

### Community 13 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (42): bulkSnoozeItemResult, bulkSnoozeRequest, changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest, shareLinkContextKey (+34 more)

### Community 14 - "auth/handlers_test.go"
Cohesion: 0.25
Nodes (7): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups()

### Community 15 - "New"
Cohesion: 0.07
Nodes (33): TestDocExportDefaultCronExprIsValid(), newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), TestJobHealthWithoutStoreDoesNothing(), cron.EntryID, Runner (+25 more)

### Community 16 - "lists.go"
Cohesion: 0.10
Nodes (41): unavailable(), SnapshotSection, buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP() (+33 more)

### Community 17 - "ServiceSnapshot"
Cohesion: 0.03
Nodes (19): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, changePatternID(), Engine, markError(), snapshotIDOrNil(), runTransformers() (+11 more)

### Community 18 - "templates.fixtures.ts"
Cohesion: 0.18
Nodes (12): web_src_api_model_index_docversion, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate(), resolveToken() (+4 more)

### Community 19 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (51): @simplewebauthn/browser, web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete (+43 more)

### Community 20 - "git_test.go"
Cohesion: 0.06
Nodes (48): fetchAllDocs(), fileName(), Exporter, NewExporter(), RunExportOnce(), slugify(), newTestStore(), readFile() (+40 more)

### Community 21 - "DBTX"
Cohesion: 0.08
Nodes (20): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), changeFilterClause(), countRows(), T (+12 more)

### Community 22 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (25): oidcElevateFlow, webAuthnFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName() (+17 more)

### Community 23 - "ServiceDetailPage.tsx"
Cohesion: 0.03
Nodes (74): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_changes_changes_usegetchanges, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart (+66 more)

### Community 24 - "App.tsx"
Cohesion: 0.04
Nodes (70): react-router-dom, AXIOS_INSTANCE, BodyType, ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn, setAccessToken() (+62 more)

### Community 25 - "ErrorWithDetails"
Cohesion: 0.06
Nodes (36): updateUserRequest, Handler, sanitizeUser(), setRefreshCookie(), writeUserWriteError(), Handler, mustHashDummyPassword(), Handler (+28 more)

### Community 26 - "package.json"
Cohesion: 0.04
Nodes (47): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+39 more)

### Community 27 - "docker_test.go"
Cohesion: 0.05
Nodes (48): buildDockerTLSConfig(), newDockerClient(), newTCPDockerClient(), init(), generateSelfSignedCert(), generateSSHHostKey(), serveOneHTTPExchange(), serveSSHDockerConn() (+40 more)

### Community 28 - "Connector"
Cohesion: 0.12
Nodes (11): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+3 more)

### Community 29 - "NewEngine"
Cohesion: 0.09
Nodes (42): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+34 more)

### Community 30 - "SystemPage.tsx"
Cohesion: 0.06
Nodes (34): web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_getgetconnectorsconnectoridgoldensnapshotquerykey, web_src_api_generated_connectors_connectors_postconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectorid, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridgoldensnapshot, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotsdiff, web_src_api_generated_connectors_connectors_usegetconnectorsconnectoridsnapshotssnapshotid, web_src_api_generated_system_system_getgetsystembackuprunsquerykey (+26 more)

### Community 31 - "RunMigrations"
Cohesion: 0.10
Nodes (37): main(), OpenDB(), TestOpenDBEnablesSQLiteForeignKeys(), TestOpenDBSetsSQLiteDurabilityPragmas(), GetMigrationStatus(), newMigrator(), collectColumns(), postgresSchemaColumns() (+29 more)

### Community 32 - "go_pkg_strings"
Cohesion: 0.04
Nodes (51): contextKey, elevationError, TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6(), TestBuildHostOverrideTableMalformedCases(), TestBuildHostOverrideTableValidOverrides(), jsonType() (+43 more)

### Community 33 - "newDocTestStore"
Cohesion: 0.06
Nodes (47): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+39 more)

### Community 34 - "fixtures.ts"
Cohesion: 0.06
Nodes (43): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+35 more)

### Community 35 - "UsersPage.tsx"
Cohesion: 0.06
Nodes (46): axios, customInstance(), web_src_api_generated_me_me, web_src_api_generated_me_me_usegetme, web_src_api_generated_users_users, web_src_api_generated_users_users_deleteusersuserid, web_src_api_generated_users_users_getgetusersquerykey, web_src_api_generated_users_users_postusersuseridresetmfa (+38 more)

### Community 36 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 37 - "home_assistant/tables.go"
Cohesion: 0.09
Nodes (39): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+31 more)

### Community 38 - "Compare"
Cohesion: 0.07
Nodes (42): configPushLanded(), driftDescription(), Checker, highestDriftSeverity(), TestCompareIgnoresEntityAttributes(), TestCompareMapKeyOrderingDoesNotAffectResult(), TestCompareStillDetectsRuleContentChanges(), Compare() (+34 more)

### Community 39 - "dispatcher_test.go"
Cohesion: 0.08
Nodes (65): TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel() (+57 more)

### Community 40 - "ConnectorRecord"
Cohesion: 0.16
Nodes (13): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), scanSyncRun(), connectorWithRole (+5 more)

### Community 41 - "mountAPIRoutes"
Cohesion: 0.12
Nodes (23): chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes(), chi.Router, mountChatRoutes() (+15 more)

### Community 42 - "DecodeKey"
Cohesion: 0.10
Nodes (24): ProviderConfig, testHandler, Handler, Handler, Handler, primaryProviderConfig(), Handler, DecodeKey() (+16 more)

### Community 43 - "NotificationCenter.tsx"
Cohesion: 0.06
Nodes (40): 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 15. `system.resync`, 2. `sync.progress`, 3. `sync.complete` (+32 more)

### Community 44 - "go_pkg_os"
Cohesion: 0.10
Nodes (22): IsGeneratedName(), pruneStale(), TestIsGeneratedName(), writeExportState(), exportCursor, exportState, go_pkg_crypto_ed25519, go_pkg_encoding_pem (+14 more)

### Community 45 - "nilToStr"
Cohesion: 0.10
Nodes (11): ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus, Store (+3 more)

### Community 46 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 47 - "NewUser"
Cohesion: 0.15
Nodes (46): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestListFiltersGrantsBeforePagination(), Handler, snapshotFixture(), snapshotRequest(), snapshotResponse() (+38 more)

### Community 48 - "ConnectorForm.tsx"
Cohesion: 0.05
Nodes (38): Frontend shell & theme (decided 2026-06), react-error-boundary, sonner, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_connectors_connectors_usegetconnectorsschema, web_src_api_model_index_connector (+30 more)

### Community 49 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (37): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+29 more)

### Community 50 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 51 - "response.go"
Cohesion: 0.09
Nodes (24): decodeStoredSnapshot(), Handler, snapshotStoreError(), Handler, Cursor(), DecodeCursor(), EncodeCursor(), T (+16 more)

### Community 52 - "RulesPage.tsx"
Cohesion: 0.05
Nodes (45): web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid, web_src_api_generated_compliance_compliance_usegetcompliancerules, web_src_api_generated_compliance_compliance_usegetcomplianceschema (+37 more)

### Community 53 - "NewChecker"
Cohesion: 0.06
Nodes (72): contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity, Rule (+64 more)

### Community 54 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

### Community 55 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (31): statusInfo, upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo(), buildFiltering() (+23 more)

### Community 56 - "Sanitize"
Cohesion: 0.17
Nodes (11): Handler, isWritableField(), loggablePath(), loggableQuery(), elevationFailureReason(), ValidateElevationHeader(), ConfigPusher, Err() (+3 more)

### Community 57 - "src/theme.ts"
Cohesion: 0.14
Nodes (24): @fontsource/space-mono, @fontsource-variable/space-grotesk, AdvancedControls(), ColorMode, commit(), load(), Persisted, PRESETS_FONTS (+16 more)

### Community 58 - "traefik/tables.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+21 more)

### Community 59 - "NewStore"
Cohesion: 0.10
Nodes (38): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+30 more)

### Community 60 - "unifi/tables_test.go"
Cohesion: 0.15
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 61 - "main"
Cohesion: 0.12
Nodes (24): main(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), NewEmbedRegistry(), RegisterOpenAICompatible() (+16 more)

### Community 62 - "Store"
Cohesion: 0.04
Nodes (39): Provider, Config, routerDeps, Handler, Registry, NewHandler(), NewHandler(), NewHandler() (+31 more)

### Community 63 - "connector_permission.go"
Cohesion: 0.13
Nodes (15): APIKeyRestriction, auditConnectorGrantDiffJSON(), APIKeyRestrictionFromContext(), ClampConnectorRole(), Store, apiKeyConnectorFilter(), getConnectorGrant(), ConnectorGrantDiff (+7 more)

### Community 64 - "Connector"
Cohesion: 0.07
Nodes (11): init(), ConfigField, Connector, buildRouteTable(), buildGatewayTable(), primaryGatewayName(), wanInterfaceName(), PathSegment() (+3 more)

### Community 65 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 66 - "middleware_test.go"
Cohesion: 0.12
Nodes (20): ConnectorRoleChecker, fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireConnectorRole(), RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel() (+12 more)

### Community 67 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 68 - "config_test.go"
Cohesion: 0.08
Nodes (35): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), runHealthcheck(), Load() (+27 more)

### Community 69 - "Config"
Cohesion: 0.12
Nodes (20): newLogger(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote(), AISettings (+12 more)

### Community 70 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 71 - "SuggestRequest"
Cohesion: 0.13
Nodes (11): claudeProvider, openAICompatibleProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList(), TestStubProviderName() (+3 more)

### Community 72 - "rewritePlaceholders"
Cohesion: 0.11
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 73 - "handlers.ts"
Cohesion: 0.07
Nodes (26): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+18 more)

### Community 74 - "runbooks_test.go"
Cohesion: 0.14
Nodes (26): runbookResp, runbookStepResp, TestAttentionRunbookLinkForAlert(), TestAttentionRunbookLinkForFinding(), createRunbookWithStep(), testApp, seedProxmoxConnector(), seedRunbook() (+18 more)

### Community 75 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 76 - "NewMalformedResponseError"
Cohesion: 0.09
Nodes (10): NewMalformedResponseError(), WantsField(), Connector, putMetadata(), Connector, agentEnabled(), Connector, Connector (+2 more)

### Community 77 - "createTestConnector"
Cohesion: 0.11
Nodes (24): TestSnapshotKeysetSummaryAndConnectorOwnership(), TestDeleteOldHealthChecks(), TestGetConnectorUptimeDeterministicOutage(), TestGetConnectorUptimeNoData(), TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR(), TestRecordHealthCheckDefaults(), assertQueryPlanUsesIndex(), Store (+16 more)

### Community 78 - "newTestHandler"
Cohesion: 0.09
Nodes (49): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+41 more)

### Community 79 - "GetTypeSchema"
Cohesion: 0.12
Nodes (25): catalog(), TestRegisteredSchema(), TestSchemaConfigValidation(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+17 more)

### Community 80 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 81 - "timeline.ts"
Cohesion: 0.14
Nodes (17): installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env(), heartbeat() (+9 more)

### Community 82 - "newTestHandler"
Cohesion: 0.22
Nodes (13): Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound(), TestDismissSuccess() (+5 more)

### Community 83 - "router.go"
Cohesion: 0.08
Nodes (34): spaHandler(), migrationFiles(), TestMigrationVersionParity(), go_pkg_github_com_wiselabz_wiselabz_internal_api, go_pkg_github_com_wiselabz_wiselabz_internal_api_alerts, go_pkg_github_com_wiselabz_wiselabz_internal_api_apikeys, go_pkg_github_com_wiselabz_wiselabz_internal_api_attention, go_pkg_github_com_wiselabz_wiselabz_internal_api_auth (+26 more)

### Community 84 - "HashToken"
Cohesion: 0.15
Nodes (15): factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted(), TestGenerateTOTPSecretProducesScannableURL() (+7 more)

### Community 85 - "Store"
Cohesion: 0.07
Nodes (11): changeServiceIDs(), Store, placeholders(), scanBackupRun(), AlertRecord, ChangeRecord, Store, scanAlert() (+3 more)

### Community 86 - "Manager"
Cohesion: 0.15
Nodes (9): cron.EntryID, Manager, JobName(), LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 87 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 88 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 89 - "notifications/handlers_test.go"
Cohesion: 0.26
Nodes (16): TestCreate(), AuthedUser(), NewHandler(), decodePaginated(), jsonHasEmptyArrayItems(), newTestStore(), seedDelivery(), TestListDeliveriesEmpty() (+8 more)

### Community 90 - "backup/backup.go"
Cohesion: 0.17
Nodes (27): connectorIDs(), docIDs(), Export(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import() (+19 more)

### Community 91 - "store/mfa_test.go"
Cohesion: 0.21
Nodes (14): Store, mfaTestUser(), TestConfirmFactorRejectsSecondConfirmedTOTP(), TestConsumeTOTPStepReplayGuard(), TestDeleteFactorLastOneAlsoLeavesRecoveryCodesForCallerToWipe(), TestDeleteUserCascadesMFAFactorsAndRecoveryCodes(), TestDeleteUserFactorsForAdminReset(), TestGetRequire2FADefaultsToNone() (+6 more)

### Community 92 - "portainer_test.go"
Cohesion: 0.10
Nodes (35): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+27 more)

### Community 93 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 94 - "doc_test.go"
Cohesion: 0.12
Nodes (22): Store, mustCreateUser(), newPostgresTestStore(), skipOnPostgres(), TestDocLockAcquireAfterExpiry(), TestDocLockConflict(), TestDocLockReleaseOnlyByHolder(), TestDocLockRenewalByHolder() (+14 more)

### Community 95 - "NewHTTPClient"
Cohesion: 0.12
Nodes (16): newConnector(), Connector, newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback(), intConfig(), newConnector(), NewHTTPClient() (+8 more)

### Community 96 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 97 - "httpx/retry_test.go"
Cohesion: 0.34
Nodes (14): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+6 more)

### Community 98 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 99 - "lifecycleManager"
Cohesion: 0.12
Nodes (9): newLifecycleManager(), Election, context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, gatedStopScheduler, lifecycleDeps, lifecycleManager (+1 more)

### Community 100 - "NewService"
Cohesion: 0.10
Nodes (28): APIKeyChecker, UserStatusChecker, TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), treatAsSafeFromContext(), NewService() (+20 more)

### Community 101 - "SuggestWithFallback"
Cohesion: 0.23
Nodes (11): StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail(), TestSuggestWithFallbackFirstProviderSucceeds(), TestSuggestWithFallbackNoProviders() (+3 more)

### Community 102 - "net/http.Handler"
Cohesion: 0.14
Nodes (15): TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), SecurityHeaders(), TestSecurityHeaders(), chi.Router (+7 more)

### Community 103 - "validate.go"
Cohesion: 0.32
Nodes (5): Config, mask(), redactDSN(), redactKVPassword(), TestRedactDSN()

### Community 104 - "NewEngine"
Cohesion: 0.17
Nodes (25): RequestedFields(), TestBaseContext(), TestRunDueSyncsRespectsLimits(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector() (+17 more)

### Community 105 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 106 - "sshStdioConn"
Cohesion: 0.10
Nodes (12): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), sshStdioConn, bufio.ReadWriter, golang.org/x/crypto/ssh.Client, golang.org/x/crypto/ssh.ClientConfig, golang.org/x/crypto/ssh.Session (+4 more)

### Community 107 - "ws.ts"
Cohesion: 0.11
Nodes (18): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+10 more)

### Community 108 - "Handler"
Cohesion: 0.24
Nodes (7): NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), Handler, RuleEvaluator

### Community 109 - "ExportToFile"
Cohesion: 0.12
Nodes (38): ExportToFile(), newTestStore(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable(), TestExportToFilePermissions(), TestImportRejectsInvalidBundleBeforeWriting(), TestImportRollsBackOnLateFailure() (+30 more)

### Community 110 - "truenas_test.go"
Cohesion: 0.24
Nodes (17): Connector, newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchIsStableAcrossCalls() (+9 more)

### Community 111 - "diagnostics/diagnostics.go"
Cohesion: 0.14
Nodes (21): CheckHealth(), Collect(), collectVersions(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets() (+13 more)

### Community 112 - "ws/ws_test.go"
Cohesion: 0.16
Nodes (23): newHeartbeat(), NewHub(), normalizeOrigin(), assertEnvelope(), decodeEnvelope(), mustMarshalEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock() (+15 more)

### Community 113 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 114 - "Store"
Cohesion: 0.18
Nodes (7): testAPIKeyChecker, APIKeyClaims, validAPIKey(), decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 115 - "Register"
Cohesion: 0.13
Nodes (24): init(), init(), init(), init(), init(), init(), init(), init() (+16 more)

### Community 116 - "newTestHarness"
Cohesion: 0.16
Nodes (18): ContextWithAPIKeyRestriction(), TestClampConnectorRole(), seedAlert(), TestListAttentionItems(), seedChange(), TestListChanges(), TestListConnectors(), enableFakeEmbedding() (+10 more)

### Community 117 - "Get"
Cohesion: 0.10
Nodes (22): applyConnectorScalarUpdates(), Handler, validateConnectorConfig(), Get(), Connector, IsSecretFieldType(), MarshalConnectorConfig(), ParseConnectorConfig() (+14 more)

### Community 118 - "RunbookRecord"
Cohesion: 0.29
Nodes (5): RunbookRecord, RunbookStepRecord, Store, scanRunbook(), scanRunbookStep()

### Community 119 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 121 - "vectorCache"
Cohesion: 0.17
Nodes (11): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), TestVectorCachePutExistingKeyUpdatesInPlace(), vectorCache, vectorEntry, vectorKey (+3 more)

### Community 122 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 124 - "WiseLabz — Architecture & Technical Decisions"
Cohesion: 0.10
Nodes (20): ADR index, AI module, API design, Backend, Build pipeline, Changes / diff contract (decided 2026-06), Connector interface, Connector management via UI (decided 2026-06-27) (+12 more)

### Community 125 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 126 - "api/auth/webauthn_test.go"
Cohesion: 0.11
Nodes (33): secondFactorInput, virtualAuthenticator, flowCookieFrom(), beginWebAuthnLogin(), credentialResponse(), finishWebAuthnLogin(), testHandler, newVirtualAuthenticator() (+25 more)

### Community 127 - "New"
Cohesion: 0.09
Nodes (29): confirm(), formatCounts(), main(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle() (+21 more)

### Community 128 - "handlers_contract_test.go"
Cohesion: 0.28
Nodes (14): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+6 more)

### Community 129 - "Handler"
Cohesion: 0.24
Nodes (7): stepAuditDetail(), validTargetType(), validVerb(), SupportsLifecycleVerb(), Handler, runbookResponse, stepResponse

### Community 130 - "logging_test.go"
Cohesion: 0.17
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 131 - "Deps"
Cohesion: 0.34
Nodes (15): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+7 more)

### Community 132 - "time.Time"
Cohesion: 0.17
Nodes (20): digestDue(), TestDigestDue(), time.Time, ChangeEntry, ComplianceSection, ConnectorDrift, DefinitionSummary, DocChangeEntry (+12 more)

### Community 133 - "connector_permission_test.go"
Cohesion: 0.34
Nodes (14): Store, newTestConnector(), newTestUser(), TestDeleteConnectorGrant(), TestFilterConnectorIDsByGrant(), TestGetUserConnectorRoleHighestAcrossSources(), TestListConnectorGrants(), TestListConnectorIDs() (+6 more)

### Community 134 - "newTestLifecycle"
Cohesion: 0.26
Nodes (9): newTestLifecycle(), startTestLifecycle(), TestLeaderStartsSchedulerAndRunsJob(), TestLifecycleManagerOrderedShutdown(), TestLifecycleManagerShutdownCancelsWorkContext(), TestLifecycleManagerWaitsForSchedulerBeforeCancelAndDBClose(), TestStandbyIsUnreadyAndRunsNoScheduler(), waitForLifecycleSignal() (+1 more)

### Community 135 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 136 - "Handler"
Cohesion: 0.27
Nodes (6): definition(), record(), reportJSON(), valid(), Handler, input

### Community 137 - "render_test.go"
Cohesion: 0.31
Nodes (13): RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden(), TestRenderMarkdown_SectionUnavailable() (+5 more)

### Community 138 - "Contributing to WiseLabz"
Cohesion: 0.14
Nodes (14): Branch naming, Commit hooks, Commit messages, Contributing to WiseLabz, Getting help, Prerequisites, Pull request process, Releasing (+6 more)

### Community 139 - "0004 — PostgreSQL leader election for background workers"
Cohesion: 0.15
Nodes (8): 0004 — PostgreSQL leader election for background workers, Consequences, Context, Decision, Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 140 - "newTestHandler"
Cohesion: 0.24
Nodes (10): Handler, newTestHandler(), seedProxmoxConnector(), TestCreate(), TestCreateStepsValidation(), TestExecuteStepForbiddenWithoutOperatorGrant(), TestExecuteStepNotFound(), TestGetNotFound() (+2 more)

### Community 141 - "api/changes_test.go"
Cohesion: 0.26
Nodes (12): testApp, seedChange(), seedChangeWithSeverity(), TestChangesAcknowledgeRoleBoundary(), TestChangesAcknowledgeSuccess(), TestChangesBulkResolveEmptyIDs(), TestChangesBulkResolveInvalidStatus(), TestChangesBulkResolvePartialFailure() (+4 more)

### Community 142 - "chat/chat.go"
Cohesion: 0.17
Nodes (13): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), SplitSections(), TestCosineSimilarityRanksClosestVectorHighest(), TestPackUnpackVectorRoundTrips() (+5 more)

### Community 143 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 144 - "NewRegistry"
Cohesion: 0.50
Nodes (11): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+3 more)

### Community 145 - ".CreateShareLink"
Cohesion: 0.23
Nodes (6): contextWithShareLink(), Handler, newShareToken(), shareLinkFromContext(), shareLinkNode, shareLinkScope

### Community 146 - "newTestHandler"
Cohesion: 0.24
Nodes (12): templateRequest(), TestListPagination(), TestPreviewDoesNotPersist(), TestTemplateErrorPaths(), TestVersionLifecycle(), Handler, newTestHandler(), TestCreate() (+4 more)

### Community 147 - "net/http.Response"
Cohesion: 0.33
Nodes (6): retryable(), sleep(), net/http.Response, RetryPolicy, retryTransport, scripted

### Community 148 - "ListSchemas"
Cohesion: 0.28
Nodes (8): countLifecycle(), TestAllConnectorImplementationsRegister(), TestConnectorCapabilitiesMatchOptionalInterfaces(), TestConnectorFailureContract(), supportedLifecycleVerbs(), ListSchemas(), TestRegisterStubRoundTrips(), go_pkg_github_com_wiselabz_wiselabz_internal_connector_connectortest

### Community 149 - "connectors_maintenance_test.go"
Cohesion: 0.33
Nodes (9): testApp, seedMaintenanceConnector(), TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(), TestGetMaintenanceWindowAnyAuthenticatedUser(), TestListActiveMaintenanceWindowsEndpoint(), TestOpenMaintenanceWindowConnectorNotFound(), TestOpenMaintenanceWindowInvalidDuration(), TestOpenMaintenanceWindowNoElevationRequired() (+1 more)

### Community 150 - "maintenance_test.go"
Cohesion: 0.29
Nodes (10): Store, mustCreateMaintenanceConnector(), TestCloseMaintenanceWindow(), TestCloseMaintenanceWindowAlreadyClosed(), TestCloseMaintenanceWindowNotFound(), TestCreateAndGetActiveMaintenanceWindow(), TestGetActiveMaintenanceWindowExpired(), TestGetActiveMaintenanceWindowOverlapping() (+2 more)

### Community 151 - "Engine"
Cohesion: 0.25
Nodes (3): Engine, sync.Map, DocRegenerator

### Community 152 - "scripts"
Cohesion: 0.15
Nodes (13): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+5 more)

### Community 153 - "api/docs_test.go"
Cohesion: 0.36
Nodes (9): testApp, seedDoc(), TestDocLockConflict(), TestDocLockHappyPath(), TestDocLockRoleBoundary(), TestDocsListAndGetSuccess(), TestDocsSaveRoleBoundary(), TestDocsSaveSuccess() (+1 more)

### Community 154 - "templatefuncs.go"
Cohesion: 0.23
Nodes (10): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+2 more)

### Community 155 - "IsSecureRequest"
Cohesion: 0.27
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 156 - "store/backup_test.go"
Cohesion: 0.32
Nodes (11): newBackupTestStore(), TestCreateBackupRun(), TestGetBackupScheduleWhenNotExists(), TestListBackupRunsPaginated(), TestPruneBackupRunsByAge(), TestPruneBackupRunsByCount(), TestPruneBackupRunsCombinedLimits(), TestPruneBackupRunsNegativeMaxBackups() (+3 more)

### Community 157 - "Contributor Covenant Code of Conduct"
Cohesion: 0.17
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 158 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 159 - "Decision"
Cohesion: 0.17
Nodes (12): 0005 — Cross-replica WebSocket event relay for active/active, Authorization and secrets, Consequences, Context, Decision, Duplicate suppression and ordering, Event ownership: local first, then relay, Mechanism: PostgreSQL LISTEN/NOTIFY (+4 more)

### Community 160 - "Backend test performance"
Cohesion: 0.17
Nodes (12): Backend test performance, CI job times, Coverage strategy, Deterministic scheduled exports (#408), Follow-ups, Measurements and validation, Measuring, Rules for new tests (+4 more)

### Community 161 - "main.tsx"
Cohesion: 0.21
Nodes (8): react-dom, App(), USE_MOCKS, web_src_index, bootstrap(), worker, enableMocks(), handlers

### Community 162 - "Elector"
Cohesion: 0.24
Nodes (6): New(), postgresDB(), TestSecondElectorWaitsThenTakesOver(), TestWatchReportsTerminatedSession(), database/sql.Conn, Elector

### Community 163 - "retention/retention_test.go"
Cohesion: 0.61
Nodes (8): RunCleanupOnce(), newTestStore(), testLogger(), TestRunCleanupAllDBErrors(), TestRunCleanupIdempotent(), TestRunCleanupPartialFailure(), TestRunCleanupPrunesOldHealthChecks(), TestRunCleanupSkipsDisabledCategories()

### Community 164 - "keyset_test.go"
Cohesion: 0.33
Nodes (10): assertSameSet(), Store, T, queryPlan(), TestKeysetQueriesUseCoveringIndexes(), TestListAuditRecordsKeysetHonoursFilters(), TestListAuditRecordsKeysetTraversal(), TestListChangesKeysetTraversal() (+2 more)

### Community 165 - "time.Duration"
Cohesion: 0.16
Nodes (7): AuthSettings, Database, Server, WebAuthnSettings, time.Duration, OIDCProvider, PoolConfig

### Community 166 - "runbook_test.go"
Cohesion: 0.32
Nodes (7): Store, newCascadeTestStore(), TestGetRunbookByTarget(), TestRunbookRoundTrip(), TestRunbookStepsCascadeOnConnectorDelete(), TestRunbookStepsCascadeOnRunbookDelete(), TestRunbookStepsRoundTrip()

### Community 167 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 168 - "WiseLabz Connector Guide"
Cohesion: 0.18
Nodes (11): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Testing without a real instance, The Connector interface (+3 more)

### Community 169 - "Product"
Cohesion: 0.18
Nodes (10): Accessibility & Inclusion, Anti-references, Brand Personality, Design Principles, Locked frontend direction (planning session, 2026-06; revised 2026-09), Product, Product decisions (pre-planning, v1), Product Purpose (+2 more)

### Community 170 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 171 - "EmbedRegistry"
Cohesion: 0.38
Nodes (3): Embedder, EmbedRegistry, sync.RWMutex

### Community 172 - ".call"
Cohesion: 0.33
Nodes (7): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture, net/http.HandlerFunc

### Community 173 - "ReportData"
Cohesion: 0.40
Nodes (5): connectorFilter(), NewGenerator(), TestGeneratorPersistsPartialReportWhenASectionQueryFails(), Generator, ReportData

### Community 174 - "seedQualityConnector"
Cohesion: 0.22
Nodes (9): TestComplianceFindingRuleDedupAndResolve(), TestComplianceRuleCRUD(), Store, newConcurrentQualityTestStore(), seedQualityConnector(), TestListQualityFindingsFilters(), TestResolveThenReopenCreatesFreshRow(), TestUpsertQualityFindingConcurrentDedup() (+1 more)

### Community 175 - "RequirePermission"
Cohesion: 0.38
Nodes (5): PermissionChecker, chi.Router, mountDashboardRoutes(), mountWorkflowRoutes(), RequirePermission()

### Community 176 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 177 - "test-shards.sh"
Cohesion: 0.40
Nodes (8): cmd_check(), cmd_matrix(), cmd_profile(), cmd_run(), cmd_timings(), die(), test-shards.sh script, shard_json()

### Community 178 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 179 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 180 - "rowScanner"
Cohesion: 0.09
Nodes (16): changedFields(), ComplianceRuleRecord, Store, scanComplianceRule(), docSearchWhere(), escapeLike(), DocRecord, Store (+8 more)

### Community 181 - "openapi_contract_test.go"
Cohesion: 0.33
Nodes (8): normalizeParams(), routerOperations(), specOperations(), TestAPIV1AliasServesSameHandlers(), TestOpenAPIHealthProbeRoutes(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 182 - "BACKUP.md"
Cohesion: 0.25
Nodes (5): Bundle format, Diagnostics Bundle, Endpoint, What's excluded, and why, What's included

### Community 183 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 184 - "@vitejs/plugin-react"
Cohesion: 0.40
Nodes (3): @tailwindcss/vite, vite, @vitejs/plugin-react

### Community 185 - "TestComplianceRuleValidation"
Cohesion: 0.29
Nodes (7): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), complianceRule(), TestValidationErrorDetails()

### Community 187 - "dashboard/handlers_test.go"
Cohesion: 0.43
Nodes (7): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout()

### Community 188 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 189 - "ComputeWindow"
Cohesion: 0.39
Nodes (6): ComputeWindow(), TestComputeWindow_CappedAt31Days(), TestComputeWindow_ExactlyAtCap(), TestComputeWindow_FirstRun(), TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(), TestComputeWindow_Watermark()

### Community 190 - "notification_delivery_test.go"
Cohesion: 0.39
Nodes (7): createTestNotification(), Store, TestDeliveryCreateAndList(), TestListDeliveriesStatusFilterAndPagination(), TestListDueDeliveries(), TestUpdateDeliveryResultNotFound(), TestUpdateDeliveryResultTransitionsAndClearsNextAttempt()

### Community 192 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 193 - "Decision"
Cohesion: 0.15
Nodes (12): 0001 — Lab-mutating operation boundaries, Addendum (#282): runbook steps are an additional entry point, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision (+4 more)

### Community 194 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 195 - "WiseLabz"
Cohesion: 0.25
Nodes (8): Code of Conduct, Configuration, Contributing, Features, License, Quick start, Supported services, WiseLabz

### Community 196 - "newHandler"
Cohesion: 0.27
Nodes (13): TestList(), TestRevoke(), JWTService(), Token(), WithAuth(), Handler, newHandler(), serve() (+5 more)

### Community 198 - "change_pattern_test.go"
Cohesion: 0.48
Nodes (6): Store, seedConnectorForChanges(), TestChangeRelatedServiceIDsAndPatternIDRoundTrip(), TestChangeRelatedServiceIDsDefaultsToEmptyArray(), TestCountRecentChangePatterns(), TestCountRecentChangesByPattern()

### Community 199 - "log/slog.Logger"
Cohesion: 0.11
Nodes (14): formatDigest(), Dispatcher, Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep() (+6 more)

### Community 201 - "computeNextRun"
Cohesion: 0.43
Nodes (5): TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules(), TestComputeNextRun_SuccessSchedulesAtCadenceAndResetsRetries(), computeNextRun()

### Community 202 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 203 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.29
Nodes (7): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, What's excluded, and why, What's included

### Community 204 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 205 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 207 - "walkCursorPages"
Cohesion: 0.40
Nodes (6): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestChangesCursorPaginationTraversal(), walkCursorPages()

### Community 208 - "SnapshotEntity"
Cohesion: 0.09
Nodes (25): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), ServiceDependency, SnapshotEntity, TestBuildContainerTableAttributes(), buildContainerTable() (+17 more)

### Community 209 - ".GetConnectorUptime"
Cohesion: 0.33
Nodes (3): Store, HealthCheckRecord, UptimeStats

### Community 211 - "engine_maintenance_test.go"
Cohesion: 0.60
Nodes (5): driftingSnapshot(), setupMaintenanceTestConnector(), TestRunSyncExpiredMaintenanceWindowBehavesNormally(), TestRunSyncNoMaintenanceWindowBehavesNormally(), TestRunSyncSuppressesChangesDuringMaintenanceWindow()

### Community 212 - "Backup Recovery: What Comes Back, and What Doesn't"
Cohesion: 0.33
Nodes (6): 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real, 4. What a restore does *not* bring back, Backup Recovery: What Comes Back, and What Doesn't, Recovery runbook (suggested order)

### Community 213 - "WiseLabz — Deployment Guide"
Cohesion: 0.33
Nodes (6): Backups, PostgreSQL support, Scaling & high availability, systemd (bare binary), WebSocket behind a reverse proxy, WiseLabz — Deployment Guide

### Community 214 - "Scheduled Doc Export"
Cohesion: 0.33
Nodes (5): Behavior, Configuration, Failure notifications, Git mode, Scheduled Doc Export

### Community 216 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 218 - "connector/connector.go"
Cohesion: 0.09
Nodes (14): Capabilities(), CapabilityDescriptor, TimeoutError, GuardedDialer(), IsDangerousIP(), newWebhookClient(), AuthError, CredentialRefresher (+6 more)

### Community 220 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 221 - "golden_snapshot_test.go"
Cohesion: 0.60
Nodes (4): Store, mustCreateGoldenSnapshotConnector(), TestGetSnapshotByID(), TestPinGoldenSnapshotRoundTrip()

### Community 223 - "Authentication design"
Cohesion: 0.40
Nodes (5): Authentication design, Destructive-action pattern: confirm + blast radius (decided 2026-06-27), OIDC group→connector roles and IdP step-up (#279 part 3), OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only), Permissions & step-up for mutating actions (decided 2026-06-27)

### Community 224 - "Development workflow"
Cohesion: 0.40
Nodes (5): Branching, Code quality, Commit conventions, Commit hooks (`lefthook`), Development workflow

### Community 225 - "MfaEnrollDialog"
Cohesion: 0.50
Nodes (5): Sync flow, qrcode, MfaEnrollDialog(), close(), done()

### Community 226 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 230 - "MISSING — deferred & future frontend features"
Cohesion: 0.50
Nodes (3): Deferred from V1 (decided during planning), MISSING — deferred & future frontend features, Suggested-later (raised in build, not yet planned)

### Community 231 - "Fixture reuse and lifecycle tests (#405–#407)"
Cohesion: 0.50
Nodes (4): CI measurements, Fixture reuse and lifecycle tests (#405–#407), Local measurements, Validation and coverage

## Knowledge Gaps
- **595 isolated node(s):** `DocTreeProps`, `ButtonProps`, `IconButtonProps`, `Size`, `Variant` (+590 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1320 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `Handler`, `Deps`, `time.Time`, `net/http.Request`, `go_pkg_context`, `Handler`, `NewRegistry`, `ServiceSnapshot`, `git_test.go`, `DBTX`, `Engine`, `ErrorWithDetails`, `store/backup_test.go`, `NewEngine`, `RunMigrations`, `retention/retention_test.go`, `dispatcher_test.go`, `.call`, `DecodeKey`, `.call`, `ReportData`, `NewUser`, `NewChecker`, `NewStore`, `newHandler`, `rewritePlaceholders`, `engine_maintenance_test.go`, `Manager`, `notifications/handlers_test.go`, `backup/backup.go`, `lifecycleManager`, `NewEngine`, `Handler`, `ExportToFile`, `diagnostics/diagnostics.go`, `newTestHarness`, `Get`, `api/auth/webauthn_test.go`, `New`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Why does `UserIDFromContext()` connect `net/http.Request` to `go_pkg_strings`, `Handler`, `middleware_test.go`, `context.Context`, `Deps`, `Handler`, `NewStore`, `mountAPIRoutes`, `RequirePermission`, `.CreateShareLink`, `response.go`, `HashToken`, `DBTX`, `net/http.ResponseWriter`, `Sanitize`, `ErrorWithDetails`, `IsSecureRequest`, `Store`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Why does `gitFixture` connect `git_test.go` to `context.Context`, `testing.T`, `Store`, `log/slog.Logger`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **What connects `DocTreeProps`, `ButtonProps`, `IconButtonProps` to the rest of the system?**
  _595 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `react` be split into smaller, more focused modules?**
  _Cohesion score 0.03438857852265673 - nodes in this community are weakly interconnected._
- **Should `newTestApp` be split into smaller, more focused modules?**
  _Cohesion score 0.021849782743637493 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.025252525252525252 - nodes in this community are weakly interconnected._