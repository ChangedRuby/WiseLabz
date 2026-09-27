# Graph Report - WiseLabz  (2026-09-27)

## Corpus Check
- 845 files · ~514,611 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 18 file(s) not represented in the graph (top: (none) 9, .toml 2, .tmpl 2)

## Summary
- 6392 nodes · 19944 edges · 236 communities (203 shown, 33 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1743 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4b64b16a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Client
- context.Context
- Errorf
- ServiceSnapshot
- SuggestRequest
- Store
- Config
- sync.Mutex
- compliance/handlers.go
- middleware.go
- diagnostics/diagnostics.go
- NewEngine
- newTestApp
- compliance/engine.go
- time.Time
- newDocTestStore
- DecodeKey
- Dispatcher
- WiseLabz Project
- Registry
- testing.T
- CommandPalette
- initial database schema
- rewritePlaceholders
- Hub
- .OIDCCallback
- react
- Service
- go_pkg_testing
- icons.tsx
- Product
- web_src_api_model_index
- newTestHandler
- RunSync
- chat/chat.go
- Engine
- Application root
- go_pkg_context
- DashboardPage.tsx
- Bottom-dock shell
- RunMigrations
- ProfilePage.tsx
- WiseLabz
- react-i18next
- rowScanner
- UsersPage.tsx
- Runner
- go_pkg_os
- App.tsx
- package.json
- net/http.ResponseWriter
- ErrorWithDetails
- home_assistant/tables.go
- User
- ServicesPage.tsx
- fixtures.ts
- serveSSHDockerConn
- dispatcher_test.go
- data.go
- lefthook Commit Hooks
- TemplateEditorPage.tsx
- net/http.Request
- Get
- @tanstack/react-query
- routerDeps
- api/auth/oidc.go
- truenas/tables.go
- nilToStr
- AppShell — Bottom Dock Shell (single variant)
- share_links_test.go
- SnapshotEntity
- Connector
- ServiceDetailPage.tsx
- ConnectorEditPage.tsx
- newTestHandler
- dependencies
- NewMalformedResponseError
- ChatPage.tsx
- NewStore
- GetTypeSchema
- Register
- Diff viewer
- Destructive connector confirmation
- single-instance deployment model
- Graphify Knowledge Graph Rules
- Topbar Notification Center (deferred from V1)
- Database: SQLite + PostgreSQL
- React + Vite Frontend
- portainer/tables.go
- NewChecker
- home_assistant_test.go
- adguardhome/tables.go
- newRouterDeps
- traefik/tables.go
- Connector
- truenas_test.go
- unifi/tables.go
- Configuration & Documentation Backup (Export/Import)
- main
- NewEngine
- settings.mock.ts
- handlers.ts
- HashToken
- AlertRecord
- timeline.ts
- NewRegistry
- connector/connector.go
- git.go
- unifi_test.go
- SystemPage.tsx
- Manager
- src/theme.ts
- VerifyBundleFile
- WiseLabz — Design Contract
- devDependencies
- response.go
- ConnectorRecord
- portainer_test.go
- gitTarget
- NotificationRecord
- go_pkg_net_http
- middleware_test.go
- handlers_contract_test.go
- ExportToFile
- Store
- Template catalog
- Change detail synthesizer
- Settings mock data
- DocTree
- pgPlaceholderDB
- safe application defaults
- Branch Naming Convention
- viper Config Loader
- chi HTTP Router
- GHCR Container Registry
- GET /api/version (undocumented ops endpoint)
- adguardhome_test.go
- scheduler/health_test.go
- connector_permission.go
- NewService
- .Fetch
- config_test.go
- Connector
- export_test.go
- httpx/retry_test.go
- Deps
- logging_test.go
- Handler
- docker_test.go
- ws/ws_test.go
- ws.ts
- compilerOptions
- traefik_test.go
- time.Duration
- WiseLabz Connector Guide
- docdiffmodel.ts
- apikey_scope.go
- all.go
- Store
- templates.fixtures.ts
- compilerOptions
- testApp
- changes/handlers_test.go
- dashboard/handlers_test.go
- Handler
- vectorCache
- templatefuncs.go
- pagination_contract_test.go
- gitFixture
- New
- Config
- connectors_health_test.go
- Contributor Covenant Code of Conduct
- net/http.Handler
- IsSecureRequest
- Store
- Connector
- Decision
- Decision
- scripts
- runRestore
- lifecycleDeps
- newSSHDockerClient
- Decision
- .call
- cursor_pagination_test.go
- export.go
- handlers_bulk_test.go
- .call
- Changelog
- mockServiceWorker.js
- apikey_scopes_test.go
- connectors_hardening_test.go
- newDockerClient
- release-please-config.json
- RateLimit
- newTCPDockerClient
- Cache
- TimeAgo
- Panel
- store package
- OpenDB
- ErrNotFound
- ErrConflict
- slog (stdlib logging)
- Zustand State Management
- Tailwind CSS
- Docker Compose Deployment
- Step by step
- webAuthnUser
- openapi_contract_test.go
- .UpdateAuthConfig
- Collect
- scanMaintenanceWindow
- Audit Trail
- Bulk Review Actions
- PULL_REQUEST_TEMPLATE.md
- dialSSHStdio
- Mermaid.tsx
- main.tsx
- Security Policy
- seedScopeFixture
- Notification Channels
- compose-smoke.sh
- ClassifyHealth
- timeoutError
- RetentionSettings
- Saved Views
- WiseLabz — v2 Backlog
- tsconfig.json
- setup-env.sh
- CHANGE_PROVENANCE.md
- fakeNotifier
- snapshotChangingNotifier
- vite-env.d.ts
- github.com/WiseLabz/wiselabz

## God Nodes (most connected - your core abstractions)
1. `newTestApp()` - 219 edges
2. `Errorf()` - 179 edges
3. `Store` - 141 edges
4. `newDocTestStore()` - 140 edges
5. `UserIDFromContext()` - 80 edges
6. `SnapshotEntity` - 74 edges
7. `react` - 74 edges
8. `cn()` - 69 edges
9. `NewStore()` - 66 edges
10. `Dashboard widgets` - 62 edges

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

## Hyperedges (group relationships)
- **Step-Up Confirmation Flow for Destructive Actions** — docs_architecture_permissions_stepup, docs_architecture_destructive_confirm_pattern, docs_openapi_auth_elevate_endpoint, docs_openapi_removal_impact_endpoint [EXTRACTED 0.90]
- **Contract-First API Codegen Pipeline** — docs_architecture_orval, docs_openapi_spec_document, docs_architecture_react_query, docs_architecture_diff_contract [EXTRACTED 0.85]
- **Dual Local/OIDC Auth Mode System** — docs_architecture_auth_design, docs_architecture_oidc_provider_config, docs_openapi_oidc_provider_schema, docs_openapi_auth_config_endpoint [EXTRACTED 0.90]

## Communities (236 total, 33 thin omitted)

### Community 0 - "net/http.Client"
Cohesion: 0.04
Nodes (31): Connector, ollamaEmbedder, openAIEmbedder, TimeoutError, NewServiceUnavailableError(), NewTimeoutError(), setHeaders(), TestValidateCustomURL() (+23 more)

### Community 1 - "context.Context"
Cohesion: 0.03
Nodes (32): fakeStatusChecker, Handler, Connector, Connector, Store, existingIDs(), placeholders(), Store (+24 more)

### Community 2 - "Errorf"
Cohesion: 0.05
Nodes (31): Handler, newToken(), sanitize(), Handler, Handler, Handler, Handler, Handler (+23 more)

### Community 3 - "ServiceSnapshot"
Cohesion: 0.04
Nodes (21): healthFakeConnector, noopValidatedConnector, ServiceSnapshot, Connector, agentEnabled(), Connector, init(), RegisterTransformer() (+13 more)

### Community 4 - "SuggestRequest"
Cohesion: 0.13
Nodes (11): claudeProvider, openAICompatibleProvider, testProvider, SuggestChunk, SuggestRequest, TestRegistryGet(), TestRegistryList(), TestStubProviderName() (+3 more)

### Community 5 - "Store"
Cohesion: 0.19
Nodes (24): connectorIDs(), docIDs(), exportDocs(), exportTemplates(), exportWithin(), AIConfigSummary, Import(), importBundle() (+16 more)

### Community 6 - "Config"
Cohesion: 0.10
Nodes (23): newLogger(), NewHandler(), NewWebAuthnService(), TestWebAuthnRPConfig(), WebAuthnRPConfig(), Config, LogSettings, IsSSHRemote() (+15 more)

### Community 7 - "sync.Mutex"
Cohesion: 0.22
Nodes (4): sync.Mutex, fakeDocRegenerator, fakeNotifier, fakeQualityChecker

### Community 8 - "compliance/handlers.go"
Cohesion: 0.20
Nodes (12): catalog(), changedFields(), NewHandler(), response(), toRule(), validRecord(), writeRuleRejection(), ComplianceRuleRecord (+4 more)

### Community 9 - "middleware.go"
Cohesion: 0.15
Nodes (16): APIKeyChecker, AuditRecorder, contextKey, elevationError, UserStatusChecker, AuthMiddleware(), elevationFailureReason(), extractBearerToken() (+8 more)

### Community 10 - "diagnostics/diagnostics.go"
Cohesion: 0.35
Nodes (10): collectVersions(), AuthProviders, Bundle, Component, Health, OIDCProviderSummary, RecentFailures, SanitizedConfig (+2 more)

### Community 11 - "NewEngine"
Cohesion: 0.09
Nodes (43): entityNodeID(), renderLabMermaid(), renderMermaid(), shortHash(), TestRenderMermaid(), TestRenderMermaidNoLinks(), Engine, NewEngine() (+35 more)

### Community 12 - "newTestApp"
Cohesion: 0.02
Nodes (200): templateBody, testApp, seedAlert(), TestAlertsBulkSnoozePartialFailure(), TestAlertsBulkSnoozeRejectsTooManyIDs(), TestAlertsBulkSnoozeRoleBoundary(), TestAlertsBulkSnoozeValidation(), TestAlertsListDaysWindow() (+192 more)

### Community 13 - "compliance/engine.go"
Cohesion: 0.05
Nodes (52): configPushLanded(), contains(), equal(), Evaluate(), findAttribute(), Catalog, Condition, Entity (+44 more)

### Community 14 - "time.Time"
Cohesion: 0.05
Nodes (21): digestDue(), formatDigest(), Dispatcher, TestDigestDue(), Store, TestComputeNextRun_BackoffNeverExceedsScheduleCadence(), TestComputeNextRun_FailureUsesBackoffSchedule(), TestComputeNextRun_ManualOnlyNeverSchedules() (+13 more)

### Community 15 - "newDocTestStore"
Cohesion: 0.02
Nodes (158): TestAPIKeyLifecycle(), TestAPIKeyNotFound(), TestLookupAPIKeyReflectsLiveRole(), TestLookupAPIKeyRejectsDisabledUser(), TestRevokeAllAPIKeysForUser(), TestTouchAPIKeyLastUsed(), TestCreateAuditRecordAndListFiltering(), TestListAllAuditRecords() (+150 more)

### Community 16 - "DecodeKey"
Cohesion: 0.11
Nodes (19): Handler, Handler, Handler, DecodeKey(), Decrypt(), DeriveKey(), Encrypt(), TestDecodeKey() (+11 more)

### Community 17 - "Dispatcher"
Cohesion: 0.13
Nodes (15): discordPayload(), sendDiscordChannel(), sendGenericWebhookChannel(), sendSlackChannel(), slackPayload(), webhookPayload(), findChannel(), findRoute() (+7 more)

### Community 18 - "WiseLabz Project"
Cohesion: 0.08
Nodes (26): ADR 0001 — Monorepo, ADR Index (docs/adr/), AI Doc Generation Module (opt-in, provider-agnostic), API Design — REST + WebSocket split, Dual Auth Design (Local JWT + OIDC), Changes/Diff Contract (infra vs doc format), Change-Aware Diff Engine, Monorepo with Go Workspaces (+18 more)

### Community 19 - "Registry"
Cohesion: 0.12
Nodes (18): Provider, StatusError, StubProvider, SuggestResult, registerFailThenSucceed(), TestIsRetryable(), TestSuggestWithFallbackAdvancesOnRetryableError(), TestSuggestWithFallbackAllFail() (+10 more)

### Community 20 - "testing.T"
Cohesion: 0.02
Nodes (155): runConfigCommand(), setValidEnv(), TestConfigPrintRedacted(), TestConfigSchema(), TestConfigUnknown(), TestConfigValidate(), newTestLifecycle(), TestLifecycleManagerOrderedShutdown() (+147 more)

### Community 21 - "CommandPalette"
Cohesion: 0.12
Nodes (20): Axios API client, Button and IconButton, CommandPalette, theme cycling command, ConfirmDialog, Dialog, ElevationConfirm, English translation catalog (+12 more)

### Community 22 - "initial database schema"
Cohesion: 0.16
Nodes (20): alerts, changes, connector config JSON, connectors, dashboard layouts, doc versions, docs, HashToken (+12 more)

### Community 23 - "rewritePlaceholders"
Cohesion: 0.10
Nodes (14): TestAPIKeyLastUsedThrottle(), doRewritePlaceholders(), rewritePlaceholders(), TestRewritePlaceholders(), TestRewritePlaceholdersCached(), database/sql.Result, database/sql.Row, database/sql.Rows (+6 more)

### Community 24 - "Hub"
Cohesion: 0.09
Nodes (14): decodeBulkRequest(), Handler, loggablePath(), loggableQuery(), Sanitize(), TestSanitize(), Hub, bulkRequest (+6 more)

### Community 25 - ".OIDCCallback"
Cohesion: 0.17
Nodes (8): Handler, newOIDCUser(), validHostPort(), OIDCClaims, OIDCProvider, OIDCProvider, github.com/coreos/go-oidc/v3/oidc.Provider, golang.org/x/oauth2.Config

### Community 26 - "react"
Cohesion: 0.03
Nodes (131): motion, react, web_src_api_generated_compliance_compliance, web_src_api_generated_compliance_compliance_deletecompliancerulesid, web_src_api_generated_compliance_compliance_getgetcompliancerulesquerykey, web_src_api_generated_compliance_compliance_postcompliancerules, web_src_api_generated_compliance_compliance_postcompliancerulestest, web_src_api_generated_compliance_compliance_putcompliancerulesid (+123 more)

### Community 27 - "Service"
Cohesion: 0.16
Nodes (13): Claims, ElevationClaims, ElevationToken, IssuePairOptions, MFAClaims, MFATicket, Service, TokenPair (+5 more)

### Community 28 - "go_pkg_testing"
Cohesion: 0.04
Nodes (39): complianceRule(), TestValidationErrorDetails(), Schema(), schemaFor(), TestSchemaMatchesConfig(), TestBuildHostOverrideTableAttributes(), buildHostOverrideTable(), isIPv6() (+31 more)

### Community 29 - "icons.tsx"
Cohesion: 0.04
Nodes (89): Live dashboard state, i18next, zustand, web_src_api_generated_docs_docs, web_src_api_generated_docs_docs_getgetdocsdocidquerykey, web_src_api_generated_docs_docs_getgetdocsdocidversionsquerykey, web_src_api_generated_docs_docs_getgetdocstreequerykey, web_src_api_generated_docs_docs_postdocsdocidaisuggest (+81 more)

### Community 30 - "Product"
Cohesion: 0.09
Nodes (24): Connector Guide (docs/connectors/CONNECTOR_GUIDE.md), Connector Interface (Name/Fetch/Validate), Connector Management via UI (full CRUD), Destructive-Action Pattern: Confirm + Blast Radius, Manager Actions (v1 scope), Permissions & Step-Up for Mutating Actions, Role Model — viewer/operator, ServiceSnapshot Data Structure (+16 more)

### Community 31 - "web_src_api_model_index"
Cohesion: 0.03
Nodes (59): msw, react-router-dom, @testing-library/react, vitest, web_src_api_generated_auth_auth_postauthelevateoidcbegin, web_src_api_generated_auth_auth_postauthelevateoidccomplete, web_src_api_model_index, web_src_api_model_index_attentionpage (+51 more)

### Community 32 - "newTestHandler"
Cohesion: 0.05
Nodes (75): mockElevateOIDCServer, secondFactorInput, virtualAuthenticator, doJSON(), testHandler, req(), TestChangePassword(), TestChangePasswordRevokesAPIKeys() (+67 more)

### Community 33 - "RunSync"
Cohesion: 0.15
Nodes (14): Connector interface, connector schema registration, reverse proxy WebSocket support, OpenAPI REST contract, destructive-action step-up authentication, operational alerts, detected changes, Compare (+6 more)

### Community 34 - "chat/chat.go"
Cohesion: 0.17
Nodes (16): buildPrompt(), TestBuildPrompt(), cosineSimilarity(), Match, packVector(), Retrieve(), SplitSections(), SyncDocEmbeddings() (+8 more)

### Community 35 - "Engine"
Cohesion: 0.18
Nodes (5): Engine, sync.Map, AlertNotifier, DocRegenerator, QualityChecker

### Community 36 - "Application root"
Cohesion: 0.18
Nodes (12): OpenAPI client generation, Generated-code lint exclusions, Motion preference provider, Vite API and WebSocket proxy, Authentication and onboarding guards, Operator-only routes, Application root, Application router (+4 more)

### Community 37 - "go_pkg_context"
Cohesion: 0.08
Nodes (16): contains(), searchString(), go_pkg_context, go_pkg_database_sql, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_coreos_go_oidc_v3_oidc, go_pkg_github_com_go_webauthn_webauthn_protocol (+8 more)

### Community 38 - "DashboardPage.tsx"
Cohesion: 0.04
Nodes (79): Connector category icon map, Dashboard widget frame, 10. `doc.lock.acquired`, 11. `doc.lock.released`, 12. `doc.lock.expired`, 13. `system.health`, 14. `system.notice`, 2. `sync.progress` (+71 more)

### Community 39 - "Bottom-dock shell"
Cohesion: 0.40
Nodes (5): Authenticated app frame, Bottom-dock shell, Primary navigation, Non-React navigation bridge, Floating dock navigation

### Community 40 - "RunMigrations"
Cohesion: 0.05
Nodes (74): main(), TestExpireAlertsOnceNoExpiredAlertsIsNoop(), TestExpireAlertsOnceNotifiesViaDispatcher(), testLogger(), expireAlertsOnce(), RunCleanupOnce(), newTestStore(), testLogger() (+66 more)

### Community 41 - "ProfilePage.tsx"
Cohesion: 0.04
Nodes (60): @simplewebauthn/browser, setAccessToken(), web_src_api_generated_auth_auth, web_src_api_generated_auth_auth_deleteauthapikeysid, web_src_api_generated_auth_auth_getgetauthapikeysquerykey, web_src_api_generated_auth_auth_postauthapikeys, web_src_api_generated_auth_auth_postauthelevate, web_src_api_generated_auth_auth_postauthelevatewebauthnbegin (+52 more)

### Community 42 - "WiseLabz"
Cohesion: 0.20
Nodes (10): WCAG 2.2 AA accessibility, Docs-first information architecture, technical homelabbers, machine-honest interface, v1 narrow manager scope, trustworthy live documentation, WiseLabz, commit quality gates (+2 more)

### Community 43 - "react-i18next"
Cohesion: 0.05
Nodes (57): Frontend, 7. `quality.finding.created` and `quality.findings.changed`, react-i18next, web_src_api_generated_alerts_alerts, web_src_api_generated_alerts_alerts_getgetalertsquerykey, web_src_api_generated_alerts_alerts_postalertsalertiddismiss, web_src_api_generated_alerts_alerts_postalertsalertidresolve, web_src_api_generated_alerts_alerts_postalertsalertidsnooze (+49 more)

### Community 44 - "rowScanner"
Cohesion: 0.05
Nodes (31): actorRoleLabel(), auditFilterClause(), Store, scanAuditRecord(), scanAuditRecordRows(), Store, scanBackupRun(), changeFilterClause() (+23 more)

### Community 45 - "UsersPage.tsx"
Cohesion: 0.05
Nodes (52): axios, AXIOS_INSTANCE, BodyType, customInstance(), ErrorType, getAccessToken(), MfaEnrollmentRequiredFn, RefreshFn (+44 more)

### Community 46 - "Runner"
Cohesion: 0.16
Nodes (7): cron.EntryID, Runner, cron.Cron, HealthStore, jobEntry, JobInfo, Notifier

### Community 47 - "go_pkg_os"
Cohesion: 0.06
Nodes (46): main(), usage(), changePromptData(), stripPromptTags(), truncateUTF8(), bulkResolveItemResult, bulkResolveRequest, go_pkg_bufio (+38 more)

### Community 48 - "App.tsx"
Cohesion: 0.06
Nodes (43): react-error-boundary, sonner, setMfaEnrollmentRequiredHandler(), AiPage, AppearancePage, AppShell, AuditPage, AuthCallbackPage (+35 more)

### Community 49 - "package.json"
Cohesion: 0.04
Nodes (48): clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view, eslint, eslint-plugin-react-hooks (+40 more)

### Community 50 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (15): Handler, Handler, webAuthnFlow, Handler, Handler, Handler, cron.EntryID, Handler (+7 more)

### Community 51 - "ErrorWithDetails"
Cohesion: 0.08
Nodes (25): updateUserRequest, Handler, sanitizeUser(), writeUserWriteError(), Handler, mustHashDummyPassword(), Handler, randomOIDCToken() (+17 more)

### Community 52 - "home_assistant/tables.go"
Cohesion: 0.08
Nodes (40): jsonType(), TestAttributeCatalogCoversEmittedKeys(), unavailable(), attrIP(), attrNumber(), attrString(), buildEntities(), buildIntegrations() (+32 more)

### Community 53 - "User"
Cohesion: 0.06
Nodes (12): sanitizeSessions(), Store, Store, RunbookRecord, Store, scanRunbook(), boolToInt(), Session (+4 more)

### Community 54 - "ServicesPage.tsx"
Cohesion: 0.05
Nodes (37): match-sorter, @radix-ui/react-popover, web_src_api_generated_connectors_connectors_deleteconnectorsconnectorid, web_src_api_generated_connectors_connectors_deleteconnectorsconnectoridmaintenancewindow, web_src_api_generated_connectors_connectors_getgetconnectorsmaintenancewindowsquerykey, web_src_api_generated_connectors_connectors_postconnectorsbulkreauth, web_src_api_generated_connectors_connectors_postconnectorsbulkrestart, web_src_api_generated_connectors_connectors_postconnectorsbulksync (+29 more)

### Community 55 - "fixtures.ts"
Cohesion: 0.06
Nodes (41): web_src_api_model_index_alert, web_src_api_model_index_alertpage, web_src_api_model_index_changedetail, web_src_api_model_index_changepage, web_src_api_model_index_changesummary, web_src_api_model_index_connectortypeschema, web_src_api_model_index_dashboardoverview, web_src_api_model_index_doc (+33 more)

### Community 56 - "serveSSHDockerConn"
Cohesion: 0.29
Nodes (6): serveOneHTTPExchange(), serveSSHDockerConn(), bufio.ReadWriter, golang.org/x/crypto/ssh.Channel, golang.org/x/crypto/ssh.ServerConfig, net.Conn

### Community 57 - "dispatcher_test.go"
Cohesion: 0.20
Nodes (46): NewDispatcher(), deliveriesFor(), findDelivery(), Dispatcher, newTestStore(), setChannelAndRoutingConfig(), setChannelConfig(), setChannelConfigJSON() (+38 more)

### Community 58 - "data.go"
Cohesion: 0.09
Nodes (37): connectorFilter(), RenderHTML(), RenderMarkdown(), sampleData(), TestRenderHTML_EscapesDocTitles(), TestRenderHTML_SectionUnavailable(), TestRenderHTML_Truncated(), TestRenderMarkdown_Golden() (+29 more)

### Community 59 - "lefthook Commit Hooks"
Cohesion: 0.50
Nodes (5): commit-msg Hook, Conventional Commits Policy, lefthook Commit Hooks, pre-commit Hook, Commit Conventions & Hook Enforcement (dev workflow)

### Community 60 - "TemplateEditorPage.tsx"
Cohesion: 0.05
Nodes (36): web_src_api_generated_templates_templates, web_src_api_generated_templates_templates_getgettemplatestemplateidquerykey, web_src_api_generated_templates_templates_getgettemplatestemplateidversionsquerykey, web_src_api_generated_templates_templates_posttemplatestemplateidpreview, web_src_api_generated_templates_templates_posttemplatestemplateidversionsrevrestore, web_src_api_generated_templates_templates_puttemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateid, web_src_api_generated_templates_templates_usegettemplatestemplateidversions (+28 more)

### Community 61 - "net/http.Request"
Cohesion: 0.09
Nodes (21): oidcElevateFlow, clearFlowCookie(), clearOIDCFlowCookie(), clearOIDCElevateFlowCookie(), readOIDCElevateFlowCookie(), setOIDCElevateFlowCookie(), oidcFlowCookieName(), readOIDCFlowCookie() (+13 more)

### Community 62 - "Get"
Cohesion: 0.08
Nodes (28): Handler, isWritableField(), validateConfigPushRequest(), capitalize(), Handler, TestDiagnosticsRedactsSecrets(), WriteElevationError(), ConfigPusher (+20 more)

### Community 63 - "@tanstack/react-query"
Cohesion: 0.06
Nodes (39): Client dispatch model, Envelope, Mock emitter (frontend-first), Naming convention, Reconnect behavior, Transport, WiseLabz WebSocket Contract (`/ws`), @tanstack/react-query (+31 more)

### Community 64 - "routerDeps"
Cohesion: 0.09
Nodes (33): routerDeps, NewHandler(), chi.Router, mountAuthRoutes(), mountMeRoutes(), mountUserRoutes(), chi.Router, mountConnectorRoutes() (+25 more)

### Community 65 - "api/auth/oidc.go"
Cohesion: 0.06
Nodes (33): TestEmailDomainAllowed(), TestOIDCConnectorRolesForGroups(), TestOIDCRoleForGroups(), connectorRoleLess(), emailDomainAllowed(), oidcConnectorRolesForGroups(), oidcRoleForGroups(), Config (+25 more)

### Community 66 - "truenas/tables.go"
Cohesion: 0.13
Nodes (40): buildDatasets(), buildDisks(), buildInterfaces(), buildNFSShares(), buildPools(), buildReplicationTasks(), buildServices(), buildSMBShares() (+32 more)

### Community 67 - "nilToStr"
Cohesion: 0.10
Nodes (12): seedDelivery(), ChatConversationRecord, Store, nilToStr(), DocVersionRecord, Store, DeliveryRecord, DeliveryStatus (+4 more)

### Community 68 - "AppShell — Bottom Dock Shell (single variant)"
Cohesion: 0.50
Nodes (4): AppShell — Bottom Dock Shell (single variant), Theme Engine — Code Default, User-Overridable, Per-User Dashboard Layout with Admin Default (v2), DashboardLayout Schema (per-user widget layout)

### Community 69 - "share_links_test.go"
Cohesion: 0.17
Nodes (42): GrantConnectorRole(), instanceAdminRole(), NewUser(), TestCreateConversationDocVisibility(), TestListFiltersGrantsBeforePagination(), Handler, newTestHandler(), TestAISuggestInvalidJSON() (+34 more)

### Community 70 - "SnapshotEntity"
Cohesion: 0.10
Nodes (22): TestBuildDNSRecordTableAttributes(), TestBuildTunnelTableAttributes(), buildDNSRecordTable(), buildTunnelTable(), SnapshotEntity, SnapshotSection, TestBuildContainerTableAttributes(), buildContainerTable() (+14 more)

### Community 71 - "Connector"
Cohesion: 0.07
Nodes (12): ConfigField, Connector, buildRouteTable(), TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), primaryGatewayName(), wanInterfaceName() (+4 more)

### Community 72 - "ServiceDetailPage.tsx"
Cohesion: 0.07
Nodes (37): ADR-0001, ADR-0003, 1. `service.status`, web_src_api_generated_connectors_connectors_postconnectorsconnectoridconfigpush, web_src_api_generated_connectors_connectors_postconnectorsconnectoridhealth, web_src_api_generated_connectors_connectors_postconnectorsconnectoridrestart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstart, web_src_api_generated_connectors_connectors_postconnectorsconnectoridstop (+29 more)

### Community 73 - "ConnectorEditPage.tsx"
Cohesion: 0.06
Nodes (31): RFC-3339, web_src_api_generated_connectors_connectors, web_src_api_generated_connectors_connectors_getgetconnectorsquerykey, web_src_api_generated_connectors_connectors_postconnectors, web_src_api_generated_connectors_connectors_postconnectorsconnectoridsync, web_src_api_generated_connectors_connectors_postconnectorsconnectoridtest, web_src_api_generated_connectors_connectors_postsync, web_src_api_generated_connectors_connectors_putconnectorsconnectorid (+23 more)

### Community 74 - "newTestHandler"
Cohesion: 0.11
Nodes (40): actionRequest(), actionResponse(), TestActionBulkGrantBoundaries(), TestActionInvalidConnectorConfig(), TestActionLifecyclePreviews(), TestActionMaintenanceLifecycle(), TestActionPermissions(), TestActionStoreFailures() (+32 more)

### Community 75 - "dependencies"
Cohesion: 0.05
Nodes (42): dependencies, axios, clsx, codemirror, @codemirror/commands, @codemirror/lang-markdown, @codemirror/state, @codemirror/view (+34 more)

### Community 76 - "NewMalformedResponseError"
Cohesion: 0.11
Nodes (37): NewMalformedResponseError(), buildAdlistTable(), buildClientTable(), buildDomainTable(), buildGroupTable(), cell(), clientIP(), groupNames() (+29 more)

### Community 77 - "ChatPage.tsx"
Cohesion: 0.06
Nodes (36): web_src_api_generated_chat_chat, web_src_api_generated_chat_chat_getgetchatconversationsidquerykey, web_src_api_generated_chat_chat_getgetchatconversationsquerykey, web_src_api_generated_chat_chat_postchatconversations, web_src_api_generated_chat_chat_postchatconversationsidmessages, web_src_api_generated_chat_chat_usegetchatconversations, web_src_api_generated_chat_chat_usegetchatconversationsid, web_src_api_generated_reports_reports (+28 more)

### Community 78 - "NewStore"
Cohesion: 0.10
Nodes (37): NewHandler(), TestBulkSnooze(), TestDismissNotFound(), TestGetNotFound(), TestListEmpty(), TestResolveNotFound(), TestSnooze(), NewStore() (+29 more)

### Community 79 - "GetTypeSchema"
Cohesion: 0.09
Nodes (34): TestRegisteredSchema(), TestSchemaConfigValidation(), TestAllConnectorImplementationsRegister(), TestRegisteredSchema(), TestAttributeCatalogCoversEmittedKeys(), TestAttributeCatalogCoversNewEntityKinds(), TestSchemaExposesAPIVersion(), TestAPIKeyIsStoredAsPassword() (+26 more)

### Community 80 - "Register"
Cohesion: 0.09
Nodes (33): init(), newConnector(), init(), Connector, init(), newGuardedClient(), TestGuardedClientRejectsLinkLocal(), TestGuardedClientRejectsLoopback() (+25 more)

### Community 81 - "Diff viewer"
Cohesion: 0.67
Nodes (3): Document diff model, Diff layout preference, Diff viewer

### Community 82 - "Destructive connector confirmation"
Cohesion: 0.67
Nodes (3): Destructive connector confirmation, Connector removal impact, Step-up reauthentication

### Community 83 - "single-instance deployment model"
Cohesion: 0.67
Nodes (3): PostgreSQL compose deployment, single-instance deployment model, SQLite compose deployment

### Community 84 - "Graphify Knowledge Graph Rules"
Cohesion: 0.67
Nodes (3): GRAPH_REPORT.md, Graphify Knowledge Graph Rules, Graphify Wiki Index

### Community 85 - "Topbar Notification Center (deferred from V1)"
Cohesion: 0.67
Nodes (3): Topbar Notification Center (deferred from V1), NotificationDelivery Schema (per-channel delivery/retry), Notification / NotificationPage Schemas

### Community 86 - "Database: SQLite + PostgreSQL"
Cohesion: 0.67
Nodes (3): Database: SQLite + PostgreSQL, golang-migrate, sqlc (type-safe SQL codegen)

### Community 87 - "React + Vite Frontend"
Cohesion: 0.67
Nodes (3): Frontend Testing Policy (deferred until rewrite), go:embed SPA Embedding, React + Vite Frontend

### Community 88 - "portainer/tables.go"
Cohesion: 0.12
Nodes (33): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEnvironmentTable(), buildStackTable(), cell(), containerRows(), environmentNames(), environmentTypeName() (+25 more)

### Community 89 - "NewChecker"
Cohesion: 0.19
Nodes (35): NewChecker(), RunStaleSweepOnce(), createComplianceRule(), createComplianceSnapshot(), createConnector(), findings(), newTestStore(), TestCheckEmptyDetectsAndAutoResolves() (+27 more)

### Community 90 - "home_assistant_test.go"
Cohesion: 0.10
Nodes (36): AllowLoopbackForTest(), Connector, homeAssistantAPI(), newTestConnector(), TestBearerHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchAppliesMaxEntities() (+28 more)

### Community 91 - "adguardhome/tables.go"
Cohesion: 0.14
Nodes (32): statusInfo, unavailable(), upstreamDependencies(), jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildClientTable(), buildDHCP(), buildDNSInfo() (+24 more)

### Community 92 - "newRouterDeps"
Cohesion: 0.10
Nodes (34): NewHandler(), TestCreate(), TestList(), TestRevoke(), AuthedUser(), JWTService(), Token(), WithAuth() (+26 more)

### Community 93 - "traefik/tables.go"
Cohesion: 0.14
Nodes (30): jsonType(), TestAttributeCatalogCoversEmittedKeys(), buildEntryPointTable(), buildMiddlewareTable(), buildOverview(), buildRouterTable(), buildServiceTable(), cell() (+22 more)

### Community 94 - "Connector"
Cohesion: 0.10
Nodes (12): NewAuthError(), TestTypedErrorsAreDistinguishableByType(), Connector, apiMessage(), controllerName(), countByKind(), statusError(), unavailable() (+4 more)

### Community 95 - "truenas_test.go"
Cohesion: 0.11
Nodes (31): entityKinds(), sectionByTitle(), TestConfigPushV5(), TestFetchAuthFailureReturnsPlaceholderSnapshot(), TestFetchBothVersions(), TestFetchDegradesPerSection(), TestRestartUnsupportedOnV5(), TestStartStopV5() (+23 more)

### Community 96 - "unifi/tables.go"
Cohesion: 0.17
Nodes (29): jsonType(), TestAttributeCatalogCoversEmittedKeys(), boolOr(), buildClientSummary(), buildDeviceTable(), buildFirewallTable(), buildNetworkTable(), buildSiteTable() (+21 more)

### Community 97 - "Configuration & Documentation Backup (Export/Import)"
Cohesion: 0.06
Nodes (28): Bundle format, Configuration & Documentation Backup (Export/Import), Endpoints, Import behavior, Manifest, checksum, and verification, 1. Every export gets a manifest and a checksum, 2. Verifying a backup actually restores, 3. Restoring for real (+20 more)

### Community 98 - "main"
Cohesion: 0.10
Nodes (26): main(), runHealthcheck(), splitOrigins(), RegisterClaude(), TestRegisterClaudeDefaults(), RegisterOllamaEmbedder(), RegisterOpenAIEmbedder(), Embedder (+18 more)

### Community 99 - "NewEngine"
Cohesion: 0.15
Nodes (27): RequestedFields(), TestBaseContext(), TestSyncCancellationRecordsFailureAndReleasesGuard(), TestSyncExcludesConcurrentRuns(), TestRefreshCredentialsDirect(), TestRefreshCredentialsUnsupportedConnector(), TestRunSyncFieldsPassesHintToConnector(), TestRunSyncFieldsSurvivesCredentialRefresh() (+19 more)

### Community 100 - "settings.mock.ts"
Cohesion: 0.08
Nodes (27): web_src_api_model_index_aiconfig, web_src_api_model_index_aifallbackprovider, web_src_api_model_index_health, web_src_api_model_index_notificationchannel, web_src_api_model_index_notificationroute, web_src_api_model_index_profileupdate, web_src_api_model_index_role, web_src_api_model_index_session (+19 more)

### Community 101 - "handlers.ts"
Cohesion: 0.07
Nodes (28): web_src_api_generated_alerts_alerts_msw, web_src_api_generated_alerts_alerts_msw_getalertsmock, web_src_api_generated_auth_auth_msw, web_src_api_generated_auth_auth_msw_getauthmock, web_src_api_generated_changes_changes_msw, web_src_api_generated_changes_changes_msw_getchangesmock, web_src_api_generated_connectors_connectors_msw, web_src_api_generated_connectors_connectors_msw_getconnectorsmock (+20 more)

### Community 102 - "HashToken"
Cohesion: 0.13
Nodes (16): setRefreshCookie(), factorJSON(), Handler, GenerateRecoveryCodes(), GenerateTOTPSecret(), NormalizeRecoveryCode(), randomRecoveryChars(), TestGenerateRecoveryCodesAreUniqueAndFormatted() (+8 more)

### Community 103 - "AlertRecord"
Cohesion: 0.14
Nodes (9): changeServiceIDs(), AlertRecord, ChangeRecord, scanAlert(), scanChange(), changePatternID(), Engine, markError() (+1 more)

### Community 104 - "timeline.ts"
Cohesion: 0.13
Nodes (17): enableMocks(), installMockWebSocket(), Window, WsMockHandle, Listenerish, MockWebSocket, Emit, env() (+9 more)

### Community 105 - "NewRegistry"
Cohesion: 0.19
Nodes (25): NewRegistry(), NewHandler(), TestAIConfigRoundTrip(), testConfig(), TestGetAuthConfig(), TestGetDecryptedAPIKeyNoKeyStored(), TestNotificationsConfigRoundTrip(), TestNotificationsConfigSigningSecret() (+17 more)

### Community 106 - "connector/connector.go"
Cohesion: 0.09
Nodes (19): GuardedDialer(), IsDangerousIP(), buildEmailMessage(), sendSMTPChannel(), splitRecipients(), TestBuildEmailMessage_SanitizesSubjectNewlines(), TestSendSMTPChannel_MissingConfig(), TestSplitRecipients() (+11 more)

### Community 107 - "git.go"
Cohesion: 0.10
Nodes (18): TestCommitMessage(), keys(), dockerSSHAddr, go_pkg_crypto_ed25519, go_pkg_encoding_pem, go_pkg_github_com_go_git_go_git_v5, go_pkg_github_com_go_git_go_git_v5_config, go_pkg_github_com_go_git_go_git_v5_plumbing (+10 more)

### Community 108 - "unifi_test.go"
Cohesion: 0.19
Nodes (25): authorized(), decodeJSONBody(), Connector, newTestConnector(), passwordConfig(), TestAPIKeyIsNotSentInPasswordMode(), TestAutoDetectReportsUniFiOSError(), TestControllerErrorMessageIsSurfaced() (+17 more)

### Community 109 - "SystemPage.tsx"
Cohesion: 0.10
Nodes (19): web_src_api_generated_system_system_getgetsystembackuprunsquerykey, web_src_api_generated_system_system_getgetsystembackupschedulequerykey, web_src_api_generated_system_system_getsystembackupschedule, web_src_api_generated_system_system_postsystembackuprun, web_src_api_generated_system_system_putsystembackupschedule, web_src_api_generated_system_system_usegethealth, web_src_api_generated_system_system_usegetsystembackupruns, web_src_api_generated_system_system_usegetsysteminfo (+11 more)

### Community 110 - "Manager"
Cohesion: 0.15
Nodes (9): NewHandler(), cron.EntryID, Manager, LogPartial(), NewManager(), ReportDefinitionRecord, ReportRecord, Store (+1 more)

### Community 111 - "src/theme.ts"
Cohesion: 0.15
Nodes (23): @fontsource/space-mono, @fontsource-variable/space-grotesk, ColorMode, commit(), load(), Persisted, PRESETS_FONTS, ThemeState (+15 more)

### Community 112 - "VerifyBundleFile"
Cohesion: 0.15
Nodes (23): BuildManifest(), BundleCounts(), ChecksumBytes(), ReadManifest(), WriteManifest(), failVerification(), LatestBundle(), newScratchStore() (+15 more)

### Community 113 - "WiseLabz — Design Contract"
Cohesion: 0.08
Nodes (23): 10. Component conventions, 1. Identity, 2. Color tokens, 3. Status grammar, 4. Typography, 5. Radii & shadows, 6. Motion, 7. Z-index scale (+15 more)

### Community 114 - "devDependencies"
Cohesion: 0.08
Nodes (24): devDependencies, eslint, eslint-plugin-react-hooks, eslint-plugin-react-refresh, @faker-js/faker, jsdom, msw, orval (+16 more)

### Community 115 - "response.go"
Cohesion: 0.12
Nodes (14): TestWritePaginatedOmitsNextCursor(), Error(), HandleStoreError(), intQuery(), JSON(), Logger(), Paginate(), WritePaginated() (+6 more)

### Community 116 - "ConnectorRecord"
Cohesion: 0.16
Nodes (12): ConnectorRecord, Store, scanConnector(), scanConnectorRows(), nullInt64ToIntPtr(), nullStrToStr(), connectorWithRole, database/sql.NullInt64 (+4 more)

### Community 117 - "portainer_test.go"
Cohesion: 0.20
Nodes (21): dockerPath(), Connector, newTestConnector(), portainerAPI(), TestAPIKeyHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchContainersStillFetchesEnvironments() (+13 more)

### Community 118 - "gitTarget"
Cohesion: 0.13
Nodes (14): commitMessage(), gitAuth(), Exporter, installHTTPS(), TestGitAuthHTTPSNoToken(), TestGitAuthHTTPSToken(), commitResult, GitOptions (+6 more)

### Community 119 - "NotificationRecord"
Cohesion: 0.14
Nodes (10): Dispatcher, Dispatcher, RunDeliveryRetries(), Store, RunDocLockSweep(), runDocLockSweep(), NotificationRecord, Store (+2 more)

### Community 120 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (31): bulkSnoozeItemResult, bulkSnoozeRequest, dashboardLayout, versionSections(), TemplateVersionSection, shareLinkContextKey, go_pkg_github_com_mark3labs_mcp_go_mcp, go_pkg_github_com_mark3labs_mcp_go_server (+23 more)

### Community 121 - "middleware_test.go"
Cohesion: 0.13
Nodes (18): fakeConnectorRoleChecker, testAuditCall, testAuditRecorder, RequireInstanceAdmin(), assertElevationAuditCalls(), boolLabel(), contextWithInstanceAdmin(), requestWithUser() (+10 more)

### Community 122 - "handlers_contract_test.go"
Cohesion: 0.18
Nodes (19): AssertMatchesSpec(), loadSpec(), specPath(), createForSpec(), decodeEnvelope(), fieldMsgs(), Handler, TestConnectorSuccessPayloadsMatchSpec() (+11 more)

### Community 123 - "ExportToFile"
Cohesion: 0.20
Nodes (20): Export(), ExportToFile(), newTestStore(), TestExportIncludesRecordsBeyondAPage(), TestExportRedactsConnectorSecrets(), TestExportToFile(), TestExportToFileCreatesDirectory(), TestExportToFileDirNotWritable() (+12 more)

### Community 136 - "adguardhome_test.go"
Cohesion: 0.20
Nodes (20): adguardAPI(), Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchDegradesWhenStatusFails() (+12 more)

### Community 137 - "scheduler/health_test.go"
Cohesion: 0.19
Nodes (10): newFakeHealthStore(), TestJobHealthOkToFailingNotifiesOnce(), TestJobHealthPanicCountsAsFailure(), TestJobHealthPersistsAcrossRestart(), JobHealthRecord, Store, scanJobHealth(), fakeHealthStore (+2 more)

### Community 138 - "connector_permission.go"
Cohesion: 0.19
Nodes (9): auditConnectorGrantDiffJSON(), getConnectorGrant(), ConnectorGrantDiff, Store, highestConnectorRole(), listOIDCConnectorGrants(), scanConnectorGrants(), upsertConnectorGrant() (+1 more)

### Community 139 - "NewService"
Cohesion: 0.18
Nodes (18): TestAuthMiddlewareAcceptsNonAdminAPIKey(), TestAuthMiddlewareAPIKeyLifecycle(), TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(), TestAuthMiddlewareThrottlesAPIKeyLastUsed(), NewService(), TestConcurrentIssuePairUniqueTokenIDs(), TestElevationExpired(), TestElevationRequiresOwner() (+10 more)

### Community 140 - ".Fetch"
Cohesion: 0.15
Nodes (10): ServiceDependency, WantsField(), TestRequestedFields(), TestWantsField(), environmentDependencies(), putMetadata(), unavailable(), networkDependencies() (+2 more)

### Community 141 - "config_test.go"
Cohesion: 0.16
Nodes (18): Load(), TestAccessTokenTTLDuration(), TestDocExportGitValidate(), TestLoadDefaults(), TestLoadEnvOverride(), TestLoadEnvOverrideAllFields(), TestLoadEnvOverrideDocExportGitSSH(), TestLoadFromYAML() (+10 more)

### Community 142 - "Connector"
Cohesion: 0.15
Nodes (7): TestBuildInterfaceTableAttributes(), buildGatewayTable(), buildInterfaceTable(), buildSystemContent(), primaryGatewayName(), wanInterfaceName(), Connector

### Community 143 - "export_test.go"
Cohesion: 0.20
Nodes (17): fetchAllDocs(), Exporter, IsGeneratedName(), NewExporter(), pruneStale(), RunExportOnce(), newTestStore(), readFile() (+9 more)

### Community 144 - "httpx/retry_test.go"
Cohesion: 0.29
Nodes (16): RetryTransport(), do(), fail(), status(), TestRetryTransportDisabled(), TestRetryTransportDoesNotRetry(), TestRetryTransportGivesUpAfterMaxRetries(), TestRetryTransportHonorsRetryAfterWithinCap() (+8 more)

### Community 145 - "Deps"
Cohesion: 0.23
Nodes (18): registerListAttentionItems(), registerListChanges(), jsonResult(), registerListConnectors(), registerSearchDocs(), connectorAllowSet(), registerListFindings(), NewHTTPHandler() (+10 more)

### Community 146 - "logging_test.go"
Cohesion: 0.18
Nodes (15): Logger(), captureLog(), TestLoggerCorrelatesErrorfWithRequestID(), TestLoggerRedactsShareToken(), TestLoggerRedactsWSTicket(), TestGetRequestIDMissing(), TestRecovererPassThrough(), TestRecovererReturns500OnPanic() (+7 more)

### Community 147 - "Handler"
Cohesion: 0.18
Nodes (4): Handler, stripLogControlChars(), Handler, BackupSchedule

### Community 148 - "docker_test.go"
Cohesion: 0.11
Nodes (17): TestConfigPush(), TestDockerWritableFields(), TestDoRequestContextTimeout(), TestDoRequestErrorCases(), TestFetchBuildsSectionsFromEndpoints(), TestFetchSurfacesMalformedSystemResponse(), TestFetchToleratesEndpointFailure(), TestFetchWithFieldsHintSkipsUnrequestedCalls() (+9 more)

### Community 149 - "ws/ws_test.go"
Cohesion: 0.20
Nodes (17): NewHub(), normalizeOrigin(), assertEnvelope(), setupWSConnection(), TestBroadcastFullQueueDoesNotBlock(), TestBroadcastToUserAfterUpgrade(), TestClientCloseDisconnect(), TestDocLockEventBroadcast() (+9 more)

### Community 150 - "ws.ts"
Cohesion: 0.11
Nodes (17): AlertCreatedPayload, AlertResolvedPayload, ChangeDetectedPayload, DocAiSuggestionPayload, DocGeneratedPayload, DocLockAcquiredPayload, DocLockExpiredPayload, DocLockReleasedPayload (+9 more)

### Community 151 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowImportingTsExtensions, isolatedModules, jsx, lib, module, moduleDetection, moduleResolution (+9 more)

### Community 152 - "traefik_test.go"
Cohesion: 0.25
Nodes (16): Connector, newTestConnector(), TestBasicAuthHeaderIsSent(), TestErrorMapping(), TestExpiredDeadlineMapsToTimeout(), TestFetchDegradesPerSection(), TestFetchHappyPath(), TestFetchSelectiveFields() (+8 more)

### Community 153 - "time.Duration"
Cohesion: 0.16
Nodes (8): retryable(), sleep(), Database, Server, time.Duration, RetryPolicy, retryTransport, PoolConfig

### Community 154 - "WiseLabz Connector Guide"
Cohesion: 0.12
Nodes (16): Conventions, Dependencies, Getting your connector merged, Health checks vs. sync, Keeping snapshots stable, Session-based and multi-flavour APIs, Sync flow, Testing without a real instance (+8 more)

### Community 155 - "docdiffmodel.ts"
Cohesion: 0.21
Nodes (14): diff, buildDocDiff(), DiffRowUnit, DocDiffModel, DocRow, fold(), toUnits(), DiffLine (+6 more)

### Community 156 - "apikey_scope.go"
Cohesion: 0.17
Nodes (12): APIKeyRestriction, testAPIKeyChecker, APIKeyRestrictionFromContext(), ClampConnectorRole(), ContextWithAPIKeyRestriction(), isSafeMethod(), TestClampConnectorRole(), treatAsSafeFromContext() (+4 more)

### Community 157 - "all.go"
Cohesion: 0.12
Nodes (15): go_pkg_github_com_wiselabz_wiselabz_internal_connector_adguardhome, go_pkg_github_com_wiselabz_wiselabz_internal_connector_cloudflare, go_pkg_github_com_wiselabz_wiselabz_internal_connector_custom, go_pkg_github_com_wiselabz_wiselabz_internal_connector_dnsresolver, go_pkg_github_com_wiselabz_wiselabz_internal_connector_docker, go_pkg_github_com_wiselabz_wiselabz_internal_connector_home_assistant, go_pkg_github_com_wiselabz_wiselabz_internal_connector_netbird, go_pkg_github_com_wiselabz_wiselabz_internal_connector_opnsense (+7 more)

### Community 158 - "Store"
Cohesion: 0.15
Nodes (4): SnapshotRecord, Store, Store, GoldenSnapshotRecord

### Community 159 - "templates.fixtures.ts"
Cohesion: 0.17
Nodes (13): web_src_api_model_index_docversion, web_src_api_model_index_template, web_src_api_model_index_templateinput, fillBody(), generatePreview(), PreviewConnector, previewConnectors, renderTemplate() (+5 more)

### Community 160 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 161 - "testApp"
Cohesion: 0.23
Nodes (9): testApp, TestBackupCreateManualRun(), TestBackupCreateManualRunFailsWhenDirNotCreatable(), TestBackupListRunsEmpty(), TestBackupRoutesRequireOperatorRole(), TestBackupScheduleGetDefaults(), TestBackupScheduleUpdate(), TestBackupScheduleUpdateDoesNotLeakSchedulerJobs() (+1 more)

### Community 162 - "changes/handlers_test.go"
Cohesion: 0.30
Nodes (14): NewHandler(), Handler, newTestHandler(), TestAcknowledgeNotFound(), TestAcknowledgeSuccess(), TestAIUpdate(), TestBulkResolve(), TestDismissNotFound() (+6 more)

### Community 163 - "dashboard/handlers_test.go"
Cohesion: 0.21
Nodes (13): Handler, newTestHandler(), TestGetAdminDefault(), TestGetLayoutFallsBackToAdminDefault(), TestOverview(), TestPutAdminDefault(), TestSaveAndResetLayout(), Handler (+5 more)

### Community 164 - "Handler"
Cohesion: 0.22
Nodes (7): definition(), record(), reportJSON(), valid(), JobName(), Handler, input

### Community 165 - "vectorCache"
Cohesion: 0.18
Nodes (10): newVectorCache(), TestVectorCacheBoundedLRU(), TestVectorCacheConcurrent(), TestVectorCacheInvalidateDocAndStalePut(), vectorCache, vectorEntry, vectorKey, go_pkg_container_list (+2 more)

### Community 166 - "templatefuncs.go"
Cohesion: 0.17
Nodes (12): dateFormat(), filterByTitle(), join(), TestDateFormat(), TestFilterByTitle(), TestJoin(), TestToJSON(), TestTruncate() (+4 more)

### Community 167 - "pagination_contract_test.go"
Cohesion: 0.21
Nodes (13): hasAllStringKeys(), httputilCalls(), receiverName(), TestBareArrayAllowlistIsCurrent(), TestListHandlersUseSharedPaginationWriter(), TestNoHandRolledPaginationEnvelopes(), writesEnvelope(), go_pkg_go_ast (+5 more)

### Community 168 - "gitFixture"
Cohesion: 0.32
Nodes (8): SetBeforePushForTest(), newGitFixture(), TestGitExportLifecycle(), TestGitExportPushRejectionReturnsError(), TestGitExportRefusesForeignDirectory(), gitFixture, Exporter, github.com/go-git/go-git/v5/plumbing/object.Commit

### Community 169 - "New"
Cohesion: 0.36
Nodes (13): TestJobHealthWithoutStoreDoesNothing(), New(), TestAddJobInvalidExpression(), TestAddJobRegistersAndFires(), TestContextGivenToJobFunction(), TestInvalidJobNameHandled(), TestJobContextDerivedFromStart(), TestJobSkipsOverlappingInvocations() (+5 more)

### Community 170 - "Config"
Cohesion: 0.19
Nodes (11): Config, TestEmbeddedSPAWithoutFrontendBuild(), CORS(), TestCORSMatchedOrigin(), TestCORSPreflightDisallowedOriginForbidden(), TestCORSUnlistedOriginGetsNoHeaders(), chi.Router, NewRouter() (+3 more)

### Community 171 - "connectors_health_test.go"
Cohesion: 0.32
Nodes (12): testApp, registerHealthFakeType(), seedHealthTestConnector(), TestConnectorsHealthDegraded(), TestConnectorsHealthDoesNotCreateSnapshot(), TestConnectorsHealthOffline(), TestConnectorsHealthOnline(), TestConnectorsHealthRecordsTimeSeriesRow() (+4 more)

### Community 172 - "Contributor Covenant Code of Conduct"
Cohesion: 0.15
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 173 - "net/http.Handler"
Cohesion: 0.20
Nodes (10): ConnectorRoleChecker, PermissionChecker, SecurityHeaders(), TestSecurityHeaders(), TreatAsSafeMethod(), RequireConnectorRole(), RequirePermission(), TestRequireConnectorRole() (+2 more)

### Community 174 - "IsSecureRequest"
Cohesion: 0.30
Nodes (10): ClientIP(), hostOnly(), IsSecureRequest(), isTrustedProxy(), TestClientIPRejectsNonIPForwardedFor(), TestClientIPTrustedPeerUsesForwardedFor(), TestClientIPUntrustedPeerIgnoresHeaders(), TestIsSecureRequestTLS() (+2 more)

### Community 175 - "Store"
Cohesion: 0.27
Nodes (4): decodeConnectorIDs(), APIKey, Store, scanAPIKey()

### Community 177 - "Decision"
Cohesion: 0.17
Nodes (11): 0001 — Lab-mutating operation boundaries, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 178 - "Decision"
Cohesion: 0.17
Nodes (11): 0002 — Start/stop lab-mutating operations, Audit, Authorization, Confirmation / step-up, Consequences, Context, Decision, Dry-run (+3 more)

### Community 179 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, dev, format, gen:api, gen:api:watch, lint, prebuild (+4 more)

### Community 180 - "runRestore"
Cohesion: 0.24
Nodes (11): confirm(), formatCounts(), runRestore(), runVerify(), newSeededStore(), TestRunRestoreImportsIntoConfiguredDatabase(), TestRunRestoreRejectsCorruptedBundle(), TestRunRestoreRequiresFileFlag() (+3 more)

### Community 181 - "lifecycleDeps"
Cohesion: 0.20
Nodes (8): newLifecycleManager(), context.CancelFunc, golang.org/x/sync/errgroup.Group, net/http.Server, sync/atomic.Bool, lifecycleDeps, lifecycleManager, ReadyState

### Community 182 - "newSSHDockerClient"
Cohesion: 0.25
Nodes (11): generateSSHHostKey(), startSSHDockerServer(), TestNewSSHDockerClientDialsAndExecutesDialStdio(), TestNewSSHDockerClientRejectsMissingHostKey(), TestNewSSHDockerClientRejectsWrongCredentials(), TestNewSSHDockerClientRejectsWrongHostKey(), TestNewSSHDockerClientSupportsSequentialRequests(), newSSHDockerClient() (+3 more)

### Community 183 - "Decision"
Cohesion: 0.18
Nodes (10): 0003 — Config-push lab-mutating operation, Authorization / confirmation / audit, Auto-revert-then-alert on mismatch, Consequences, Context, Decision, Field-level partial update via a per-connector whitelist, Out of scope (+2 more)

### Community 184 - ".call"
Cohesion: 0.47
Nodes (7): fixture, Handler, newFixture(), TestBulkSnoozeAuthzPerItem(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestMutationAuthz()

### Community 185 - "cursor_pagination_test.go"
Cohesion: 0.31
Nodes (9): cursorPage, decodeCursorPage(), testApp, TestAuditCursorPaginationTraversal(), TestAuditOffsetPaginationUnchanged(), TestAuditRejectsMalformedCursor(), TestChangesCursorPaginationTraversal(), TestSyncsCursorPaginationUsesHeader() (+1 more)

### Community 186 - "export.go"
Cohesion: 0.31
Nodes (8): badRegexMessage(), complianceCondition(), TestComplianceRulesCRUDAndAdminGate(), TestComplianceRuleValidation(), validComplianceRule(), fileName(), slugify(), go_pkg_regexp

### Community 187 - "handlers_bulk_test.go"
Cohesion: 0.47
Nodes (9): bulkReq(), bulkResults(), createBulkFakeConnector(), Handler, registerBulkFakeConnector(), TestBulkReauth(), TestBulkRestart(), TestBulkSync() (+1 more)

### Community 188 - ".call"
Cohesion: 0.44
Nodes (6): Handler, newFixture(), TestGetAuthz(), TestListFiltersByGrantAndPaginates(), TestResolveAuthz(), fixture

### Community 189 - "Changelog"
Cohesion: 0.20
Nodes (9): [0.2.0](https://github.com/WiseLabz/WiseLabz/compare/v0.1.0...v0.2.0) (2026-09-12), 0.3.0 (2026-09-14), ⚠ BREAKING CHANGES, Bug Fixes, Changelog, Changelog, Features, Unreleased (+1 more)

### Community 190 - "mockServiceWorker.js"
Cohesion: 0.36
Nodes (8): activeClientIds, getResponse(), handleRequest(), IS_MOCKED_RESPONSE, resolveMainClient(), respondWithMock(), sendToClient(), serializeRequest()

### Community 191 - "apikey_scopes_test.go"
Cohesion: 0.47
Nodes (8): createKey(), testApp, newConnector(), TestAPIKeyCreateValidation(), TestAPIKeyDefaultsToFullScope(), TestConnectorRestrictedAPIKey(), TestReadOnlyAPIKey(), TestReadOnlyAPIKeyCapsConnectorRoleAtViewer()

### Community 192 - "connectors_hardening_test.go"
Cohesion: 0.25
Nodes (8): testApp, init(), TestConnectorsCreateAcceptsValidConfig(), TestConnectorsCreateRejectsInvalidEnum(), TestConnectorsCreateRejectsMalformedConfig(), TestConnectorsSyncAcceptsFieldsHint(), TestConnectorsUpdateRejectsMalformedConfig(), waitForSyncRuns()

### Community 193 - "newDockerClient"
Cohesion: 0.22
Nodes (8): newDockerClient(), init(), TestNewDockerClientDialsUnixSocket(), TestNewDockerClientRejectsUnsupportedScheme(), TestValidateCompositeRef(), TestValidateRefSegment(), TestValidateUnixSocketPath(), ValidateUnixSocketPath()

### Community 194 - "release-please-config.json"
Cohesion: 0.22
Nodes (8): changelog-sections, changelog-type, extra-files, include-component-in-tag, last-release-sha, packages, release-type, $schema

### Community 195 - "RateLimit"
Cohesion: 0.29
Nodes (6): TestRateLimit(), RateLimit(), golang.org/x/time/rate.Limit, golang.org/x/time/rate.Limiter, limiterStore, visitor

### Community 196 - "newTCPDockerClient"
Cohesion: 0.25
Nodes (8): buildDockerTLSConfig(), newTCPDockerClient(), generateSelfSignedCert(), TestNewTCPDockerClientMutualTLS(), TestNewTCPDockerClientNoTLSWhenNoCert(), TestNewTCPDockerClientRejectsInvalidCertPair(), Unwrap(), crypto/tls.Config

### Community 197 - "Cache"
Cohesion: 0.43
Nodes (5): Cache, New(), Cache[V], entry, V

### Community 208 - "Step by step"
Cohesion: 0.25
Nodes (8): 1. Create the package, 2. Define your config schema, 3. Implement the interface, 4. Register the connector, 5. Add the barrel import, 6. Write tests, 7. Document config fields, Step by step

### Community 209 - "webAuthnUser"
Cohesion: 0.33
Nodes (3): webAuthnUser, github.com/go-webauthn/webauthn/webauthn.Credential, github.com/google/uuid.UUID

### Community 210 - "openapi_contract_test.go"
Cohesion: 0.48
Nodes (6): normalizeParams(), routerOperations(), specOperations(), TestOpenAPIMatchesRouter(), chi.Routes, go_pkg_go_yaml_in_yaml_v3

### Community 211 - ".UpdateAuthConfig"
Cohesion: 0.43
Nodes (3): Handler, oidcProviderJSON(), boolToInt()

### Community 212 - "Collect"
Cohesion: 0.43
Nodes (7): CheckHealth(), Collect(), newTestStore(), TestCheckHealthReportsDegradedOnClosedDB(), TestCollectIncludesHealthVersionsAndSchedule(), TestCollectListsRecentFailures(), TestCollectRedactsConnectorSecrets()

### Community 213 - "scanMaintenanceWindow"
Cohesion: 0.48
Nodes (3): Store, scanMaintenanceWindow(), MaintenanceWindowRecord

### Community 214 - "Audit Trail"
Cohesion: 0.29
Nodes (6): Audit Trail, Endpoint, Keyset (cursor) pagination, Retention, What's not recorded, What's recorded

### Community 215 - "Bulk Review Actions"
Cohesion: 0.29
Nodes (6): Auditability, Bulk Review Actions, Endpoint, Frontend, Partial failure is not batch failure, What counts as low-risk

### Community 216 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.29
Nodes (6): Breaking changes, Checklist, Description, For connector PRs only, Screenshots or logs, Type of change

### Community 217 - "dialSSHStdio"
Cohesion: 0.33
Nodes (5): TestDialSSHStdioHonorsContextCancel(), closeQuietly(), dialSSHStdio(), golang.org/x/crypto/ssh.ClientConfig, io.Closer

### Community 218 - "Mermaid.tsx"
Cohesion: 0.47
Nodes (4): mermaid, cssVar(), Mermaid(), resolveColor()

### Community 219 - "main.tsx"
Cohesion: 0.40
Nodes (4): react-dom, App(), USE_MOCKS, web_src_index

### Community 220 - "Security Policy"
Cohesion: 0.33
Nodes (5): Reporting a vulnerability, Security Policy, Supported versions, What counts as a security vulnerability, What we commit to

### Community 221 - "seedScopeFixture"
Cohesion: 0.60
Nodes (4): Store, seedScopeFixture(), TestListDocSectionEmbeddingsFiltersByGrant(), TestMergedAttentionItemsFiltersByGrant()

### Community 222 - "Notification Channels"
Cohesion: 0.40
Nodes (4): Adding a channel type, Channel reference, Notification Channels, Webhook signing (HMAC-SHA256)

### Community 223 - "compose-smoke.sh"
Cohesion: 0.40
Nodes (3): COMPOSE_SMOKE_ENV_FILE, COMPOSE_SMOKE_PORT, compose-smoke.sh script

### Community 224 - "ClassifyHealth"
Cohesion: 0.67
Nodes (3): ClassifyHealth(), TestClassifyHealth(), TestClassifyHealthPerTypeThreshold()

### Community 227 - "Saved Views"
Cohesion: 0.50
Nodes (3): Endpoints, Saved Views, Scope

## Ambiguous Edges - Review These
- `Topbar Notification Center (deferred from V1)` → `NotificationDelivery Schema (per-channel delivery/retry)`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to
- `Topbar Notification Center (deferred from V1)` → `Notification / NotificationPage Schemas`  [AMBIGUOUS]
  docs/MISSING.md · relation: conceptually_related_to

## Knowledge Gaps
- **573 isolated node(s):** `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult`, `bulkResolveRequest`, `bulkResolveItemResult` (+568 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1318 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `NotificationDelivery Schema (per-channel delivery/retry)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Topbar Notification Center (deferred from V1)` and `Notification / NotificationPage Schemas`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `Store` connect `Store` to `Errorf`, `Config`, `sync.Mutex`, `compliance/handlers.go`, `NewEngine`, `compliance/engine.go`, `time.Time`, `export_test.go`, `Deps`, `Dispatcher`, `rewritePlaceholders`, `newTestHandler`, `testApp`, `changes/handlers_test.go`, `chat/chat.go`, `Engine`, `go_pkg_context`, `Handler`, `RunMigrations`, `gitFixture`, `Config`, `rowScanner`, `net/http.ResponseWriter`, `ErrorWithDetails`, `runRestore`, `lifecycleDeps`, `.call`, `dispatcher_test.go`, `data.go`, `.call`, `net/http.Request`, `routerDeps`, `nilToStr`, `share_links_test.go`, `NewStore`, `Collect`, `NewChecker`, `newRouterDeps`, `main`, `NewEngine`, `AlertRecord`, `NewRegistry`, `Manager`, `VerifyBundleFile`, `response.go`, `ExportToFile`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `noopValidatedConnector` connect `ServiceSnapshot` to `connectors_hardening_test.go`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **What connects `github.com/WiseLabz/wiselabz`, `bulkSnoozeRequest`, `bulkSnoozeItemResult` to the rest of the system?**
  _573 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Client` be split into smaller, more focused modules?**
  _Cohesion score 0.03555686159271231 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.028192371475953566 - nodes in this community are weakly interconnected._