-- Real data seeded from:
--   * Teguh ATS CV (experiences)
--   * github.com/wicahma public repositories (projects, excluding the private patok repo)
--   * public documentation sources (articles, see attribution inside each article)

-- 1. Experiences -------------------------------------------------------------
DELETE FROM experiences;

INSERT INTO experiences (id, company, location, total_duration, tags, is_current, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'Berijalan Technocenter', 'Yogyakarta, Indonesia', '2 yrs 5 mos',
 ARRAY['TypeScript', 'C#', '.NET', 'Go', 'Next.js', 'React', 'PostgreSQL', 'Microservices', 'Identity & Access Management', 'Microsoft Active Directory', 'API Gateway', 'SEO'],
 true, NOW()),
('22222222-2222-2222-2222-222222222222', 'DigistarClass', 'Remote, Indonesia', '3 mos',
 ARRAY['React', 'Software Analysis', 'Solution Design', 'Team Communication'],
 false, NOW() - INTERVAL '1 day'),
('33333333-3333-3333-3333-333333333333', 'Bank Mandiri x Rakamin Academy', 'Remote, Indonesia', '2 mos',
 ARRAY['Kotlin', 'Android', 'SQLite', 'Gradle', 'REST API'],
 false, NOW() - INTERVAL '2 days'),
('44444444-4444-4444-4444-444444444444', 'Bangkit Academy 2024 H1 Cohort', 'Remote, Indonesia', '3 mos',
 ARRAY['Kotlin', 'Android XML', 'Flutter', 'Express', 'TypeScript', 'Next.js', 'React', 'NextAuth', 'Tailwind CSS', 'Git'],
 false, NOW() - INTERVAL '3 days'),
('55555555-5555-5555-5555-555555555555', 'PT Sebangku Jaya Abadi', 'Indonesia', '6 mos',
 ARRAY['JavaScript', 'Express', 'React', 'Vite', 'Redux', 'React Router', 'JWT', 'OAuth 1'],
 false, NOW() - INTERVAL '4 days'),
('66666666-6666-6666-6666-666666666666', 'Arrizki Tour', 'Indonesia', '7 mos',
 ARRAY['Next.js', 'React', 'CMS', 'State Management', 'API Integration'],
 false, NOW() - INTERVAL '5 days');

-- 2. Company logos (inline SVG data URIs; small enough to live on the row)
UPDATE experiences SET logo = v.logo
FROM (VALUES
  ('Berijalan Technocenter', 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHJ4PSIxMCIgZmlsbD0iIzBCMkU1OSIvPgogIDxwYXRoIGQ9Ik0yMCA2LjggMzAuNCAxMXY4LjljMCA1LjktNC40IDEwLjUtMTAuNCAxMi4zLTYtMS44LTEwLjQtNi40LTEwLjQtMTIuM1YxMUwyMCA2LjhaIiBmaWxsPSIjMjBCRDVBIiBmaWxsLW9wYWNpdHk9IjAuMTQiIHN0cm9rZT0iIzIwQkQ1QSIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KICA8cGF0aCBkPSJNMjAgMTUuNCAxNi4xIDIzLjJNMjAgMTUuNCAyMy45IDIzLjJNMTYuMSAyMy4yaDcuOCIgc3Ryb2tlPSIjMjBCRDVBIiBzdHJva2Utb3BhY2l0eT0iMC41IiBzdHJva2Utd2lkdGg9IjEuMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIi8+CiAgPGNpcmNsZSBjeD0iMjAiIGN5PSIxNS4yIiByPSIyIiBmaWxsPSIjMjBCRDVBIi8+CiAgPGNpcmNsZSBjeD0iMTYuMSIgY3k9IjIzLjQiIHI9IjIiIGZpbGw9IiMyMEJENUEiLz4KICA8Y2lyY2xlIGN4PSIyMy45IiBjeT0iMjMuNCIgcj0iMiIgZmlsbD0iIzIwQkQ1QSIvPgo8L3N2Zz4K'),
  ('Bank Mandiri x Rakamin Academy', 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHJ4PSIxMCIgZmlsbD0iIzAwM0Q3OSIvPgogIDxwYXRoIGQ9Ik05LjUgMzBWMTIuNUgxNEwyMCAyMC44TDI2IDEyLjVIMzAuNVYzMEgyNi4zVjIwLjNMMjAgMjguNkwxMy43IDIwLjNWMzBIOS41WiIgZmlsbD0iI0Y5QjIzMyIvPgo8L3N2Zz4K'),
  ('Bangkit Academy 2024 H1 Cohort', 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8ZGVmcz4KICAgIDxsaW5lYXJHcmFkaWVudCBpZD0iYmFuZ2tpdC1yaXNlIiB4MT0iMTAiIHkxPSIzMSIgeDI9IjMwIiB5Mj0iOSIgZ3JhZGllbnRVbml0cz0idXNlclNwYWNlT25Vc2UiPgogICAgICA8c3RvcCBzdG9wLWNvbG9yPSIjMUE3M0U4Ii8+CiAgICAgIDxzdG9wIG9mZnNldD0iMSIgc3RvcC1jb2xvcj0iIzM0QTg1MyIvPgogICAgPC9saW5lYXJHcmFkaWVudD4KICA8L2RlZnM+CiAgPHJlY3Qgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiByeD0iMTAiIGZpbGw9IiMwRTFCMzMiLz4KICA8cGF0aCBkPSJNMjAgOC41IDMyIDIyaC02LjR2OS41SDE0LjRWMjJIOEwyMCA4LjVaIiBmaWxsPSJ1cmwoI2JhbmdraXQtcmlzZSkiLz4KPC9zdmc+Cg=='),
  ('DigistarClass', 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHJ4PSIxMCIgZmlsbD0iI0M4MTAyRSIvPgogIDxwYXRoIGQ9Ik0xMi41IDkuNWg2LjNjNi40IDAgMTAuNyA0LjIgMTAuNyAxMC41UzI1LjIgMzAuNSAxOC44IDMwLjVoLTYuM1Y5LjVabTQuNCAzLjl2MTMuMmgxLjZjNCAwIDYuMy0yLjUgNi4zLTYuNnMtMi4zLTYuNi02LjMtNi42aC0xLjZaIiBmaWxsPSIjRkZGRkZGIi8+CiAgPGNpcmNsZSBjeD0iMzAuOCIgY3k9IjEwLjQiIHI9IjIuMyIgZmlsbD0iI0ZGQzcwMCIvPgo8L3N2Zz4K'),
  ('PT Sebangku Jaya Abadi', 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHJ4PSIxMCIgZmlsbD0iI0I0NTMwOSIvPgogIDxyZWN0IHg9IjEwIiB5PSIxMy41IiB3aWR0aD0iMjAiIGhlaWdodD0iNCIgcng9IjIiIGZpbGw9IiNGRkZGRkYiLz4KICA8cmVjdCB4PSIxMiIgeT0iMTkuOCIgd2lkdGg9IjE2IiBoZWlnaHQ9IjQiIHJ4PSIyIiBmaWxsPSIjRkZGRkZGIiBmaWxsLW9wYWNpdHk9IjAuNzIiLz4KICA8cGF0aCBkPSJNMTQgMjYuNHY1LjFNMjYgMjYuNHY1LjEiIHN0cm9rZT0iI0ZGRkZGRiIgc3Ryb2tlLXdpZHRoPSIyLjQiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLW9wYWNpdHk9IjAuOSIvPgo8L3N2Zz4K'),
  ('Arrizki Tour', 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8ZGVmcz4KICAgIDxsaW5lYXJHcmFkaWVudCBpZD0iYXJyaXpraS1za3kiIHgxPSI2IiB5MT0iMzQiIHgyPSIzNCIgeTI9IjYiIGdyYWRpZW50VW5pdHM9InVzZXJTcGFjZU9uVXNlIj4KICAgICAgPHN0b3Agc3RvcC1jb2xvcj0iIzAyODRDNyIvPgogICAgICA8c3RvcCBvZmZzZXQ9IjEiIHN0b3AtY29sb3I9IiMyMkQzRUUiLz4KICAgIDwvbGluZWFyR3JhZGllbnQ+CiAgPC9kZWZzPgogIDxyZWN0IHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgcng9IjEwIiBmaWxsPSJ1cmwoI2Fycml6a2ktc2t5KSIvPgogIDxwYXRoIGQ9Ik0yMCA3LjhjLTQuNSAwLTguMSAzLjUtOC4xIDcuOSAwIDUuOSA4LjEgMTUuNSA4LjEgMTUuNXM4LjEtOS42IDguMS0xNS41YzAtNC40LTMuNi03LjktOC4xLTcuOVoiIGZpbGw9IiNGRkZGRkYiLz4KICA8Y2lyY2xlIGN4PSIyMCIgY3k9IjE1LjYiIHI9IjMuNSIgZmlsbD0iIzAyODRDNyIvPgo8L3N2Zz4K')
) AS v(company, logo)
WHERE experiences.company = v.company;

-- 2. Experience jobs ---------------------------------------------------------
INSERT INTO experience_jobs (id, experience_id, position, employment_type, description, competency, start_date, end_date, is_current) VALUES
('a1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'Product Continuity Analyst', 'Full-time',
 'Took over product continuity for ACC One once the microservice migration shipped: owning release stability, following up operational findings, and keeping requirements aligned across the ACC ecosystem.',
 ARRAY[
   'Own continuity of ACC One across the CMS, user side, IAM, and gateway services delivered during the migration.',
   'Monitor release stability and follow up operational findings with the teams that own each service.',
   'Keep business requirements aligned with the centralized IAM and gateway architecture serving roughly 40 internal applications.',
   'Turn user and business feedback into actionable product requirements for the next release cycle.'
 ],
 '2026-08-01', '2027-02-28', true),
('a1111111-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'Fullstack Developer', 'Full-time',
 'Fullstack developer on ACC One, Berijalan''s insurance platform. Led the CMS and user-side migration onto a microservice architecture, then hardened the gateway layer every service routes through.',
 ARRAY[
   'Migrated the ACC One CMS and user-facing application onto a microservice architecture, breaking one deployable into independently releasable services.',
   'Researched, upgraded, and fixed the Express Gateway and .NET Gateway that front the application, keeping service routing stable throughout the migration.',
   'Built a centralized Identity & Access Management system consumed by roughly 40 internal ACC applications, integrated directly with Microsoft Active Directory for user data management.',
   'Maintained, analyzed, and continuously improved Berijalan products for quality and runtime performance.',
   'Improved search visibility for the Berijalan and ACC One websites through applied SEO practice.',
   'Taught the Fullstack JavaScript bootcamp three times.',
   'Instructed the Berijalan x Himatif UPN Yogyakarta Fullstack Developer bootcamp.'
 ],
 '2024-09-01', '2026-07-31', false),
('a2222222-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Frontend Developer', 'Internship',
 'Studied and applied growth mindset practices, developed cross-team communication through a software development case study, contributed to requirements analysis and solution design, and learned web interface development with the React library.',
 ARRAY['React', 'Requirements Analysis', 'Solution Design', 'Team Communication'],
 '2024-08-01', '2024-10-31', false),
('a3333333-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333333', 'Mobile Developer', 'Project-based Internship',
 'Developed Android and iOS applications covering the Android lifecycle and environment, UI, design patterns, and platform security. Used SQLite as local storage, worked with the Gradle build compiler and Kotlin, and implemented API calls into the final bootcamp project.',
 ARRAY['Android Development', 'Kotlin', 'SQLite', 'Gradle', 'REST API'],
 '2023-11-01', '2023-12-31', false),
('a4444444-1111-1111-1111-111111111111', '44444444-4444-4444-4444-444444444444', 'Mobile Developer & Project Manager', 'Internship',
 'Built an Android application in Kotlin XML applying Android best practices such as DataStore, dependency injection, ViewModel providers, and database modelling. Kept development aligned with the business flow, applied Git version control for project management, and designed system requirements including ERD, flowcharts, prototypes, and documentation. Also built a Flutter inventory application for government data, an Express + TypeScript backend, and a Next.js/React/NextAuth/Tailwind admin site.',
 ARRAY['Android Development', 'Kotlin', 'Flutter', 'Project Management', 'System Design', 'Technical Documentation'],
 '2023-04-01', '2023-06-30', false),
('a5555555-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555', 'Backend Developer', 'Full-time',
 'Developed applications under Agile Software Development principles with the team. Built APIs and their documentation in JavaScript with the Express framework, and created an admin panel with React, Vite, Redux state management, and React Router, implementing OAuth 1 token validation with JWT.',
 ARRAY['REST API Design', 'API Documentation', 'Agile Development', 'Admin Panel Development', 'JWT Authentication'],
 '2023-02-01', '2023-07-31', false),
('a6666666-1111-1111-1111-111111111111', '66666666-6666-6666-6666-666666666666', 'Software Developer', 'Full-time',
 'Designed the travel website following company requirements and built the full system with a simple CMS so admins can manage product data independently. Implemented the site with Next.js, React, plus libraries for state management, navigation, error handling, and loading states. Maintained the stability of the website and API for the agreed period.',
 ARRAY['Web Design', 'CMS Development', 'Next.js', 'Error Handling', 'System Maintenance'],
 '2022-12-01', '2023-06-30', false);

-- 3. Projects ----------------------------------------------------------------
DELETE FROM projects;

INSERT INTO projects (id, title, description, problem, solution, image, tags, year, link) VALUES
('77777777-7777-7777-7777-777777777701', 'Onesist',
 'Full-stack desktop and web workspace for System Analysts: FSD documents, API specs, ERD schemas, test cases, task breakdowns, traceability matrices, and canvas sketches in one place, wired to local CLI AI agents.',
 'System analysts produce FSDs, API specs, ERDs, SIT cases, and traceability matrices across scattered files, while AI CLI agents generate artifacts that have no unified viewer or review workflow.',
 'Built a Tauri desktop app plus web dashboard with MDXEditor-based FSD analysis, an ERD studio (ReactFlow + Dagre over DBML), parsed OpenAPI spec viewer, sprint/phase task breakdown with Jira and Monday export, SIT case runner, BR-FR-DS-TC traceability with AI gap detection, Excalidraw canvas, and an embedded xterm.js agent terminal.',
 'https://picsum.photos/seed/onesist/800/600',
 ARRAY['TypeScript', 'Rust', 'Tauri', 'React', 'ReactFlow', 'DBML', 'Excalidraw', 'LLM Agents'],
 '2026', 'https://onesist.rogasper.com'),
('77777777-7777-7777-7777-777777777702', 'TypeScript Gateway',
 'Production-grade HTTP API gateway written in pure TypeScript on native node:http — zero runtime dependencies, one process, no frameworks.',
 'Most gateways pull in a framework, router, HTTP client, metrics library, and dozens of transitive packages before proxying a single byte, inflating the audit surface and supply-chain risk.',
 'Implemented every hot-path component — radix router, connection pool, circuit breaker, cache, compression — on the Node.js standard library, reaching sub-10ms P99 latency and 150k+ req/s throughput with a node_modules you can read in one sitting.',
 'https://picsum.photos/seed/tsgateway/800/600',
 ARRAY['TypeScript', 'Node.js', 'API Gateway', 'Radix Router', 'Circuit Breaker', 'Zero Dependencies'],
 '2026', 'https://tsgate.diama.dev'),
('77777777-7777-7777-7777-777777777703', 'Patok',
 'Indonesian spatial boundaries and address cleansing platform covering 38 provinces down to RT/RW, roads, and POIs by combining Kemendagri and OpenStreetMap data.',
 'Inconsistent regional boundary datasets and unstructured address data across Indonesia made geospatial queries, reverse geocoding, and jurisdiction lookups slow and inaccurate.',
 'Engineered a spatial data platform on PostGIS with automated GeoJSON harvesting pipelines, geometry compression, hierarchical boundary explorer APIs via Hono, and an interactive MapLibre GL web client.',
 'https://picsum.photos/seed/patok/800/600',
 ARRAY['TypeScript', 'PostGIS', 'PostgreSQL', 'Hono', 'MapLibre GL', 'Python', 'GIS'],
 '2026', 'https://patok.diama.dev'),
('77777777-7777-7777-7777-777777777704', 'pi-bridge',
 'HTTP bridge that supervises the Pi coding agent from any orchestrator: fire-and-forget dispatch, signed push callbacks, watchdog reconciliation, and heartbeats.',
 'A coding agent run takes minutes to hours, so blocking the caller wastes it, while a push-only design loses tasks whenever a callback is dropped.',
 'Made submit return a task id in milliseconds and deliver completion through a signed callback, added a watchdog that reconciles tasks whose process died or whose callback never arrived, and streamed heartbeat status lines. Standard library only, no third-party runtime dependencies.',
 'https://picsum.photos/seed/pibridge/800/600',
 ARRAY['Python', 'HTTP', 'Webhooks', 'systemd', 'Agent Orchestration'],
 '2026', 'https://github.com/wicahma/pi-bridge'),
('77777777-7777-7777-7777-777777777705', 'Blindfold',
 'Secret redaction boundary for LLM agents: registers actual .env values and their transforms so they are masked before the model ever sees them.',
 'Agent harnesses scan tool output for credential patterns such as sk- or xoxb-, so non-patterned secrets like SUDO_PASSWORD pass straight to the model, and base64/hex/reversed variants slip past the same scanner.',
 'Took the opposite approach: register real values and their transforms into the harness exact-substring scrub boundary, so no pattern guessing is involved and a byte match is masked. Python standard library only.',
 'https://picsum.photos/seed/blindfold/800/600',
 ARRAY['Python', 'Security', 'Secret Redaction', 'LLM Agents'],
 '2026', 'https://github.com/wicahma/blindfold'),
('77777777-7777-7777-7777-777777777706', 'SilverBullet MCP',
 'Tiny MCP server that gives any MCP-capable agent read/write access to a SilverBullet space, with no dependencies and nothing to patch inside SilverBullet.',
 'Agents had no standard way to read and write SilverBullet pages, and modifying SilverBullet itself would couple the knowledge base to the integration.',
 'Forwarded MCP tool calls to SilverBullet built-in /.fs HTTP API through a stdlib-only translator exposing list, read, write, append, and delete; speaks both stdio and streamable HTTP, and ships on npm for npx use.',
 'https://picsum.photos/seed/silverbullet/800/600',
 ARRAY['Python', 'MCP', 'SilverBullet', 'stdio', 'Streamable HTTP'],
 '2026', 'https://github.com/wicahma/silverbullet-mcp'),
('77777777-7777-7777-7777-777777777707', 'Corporate Underground',
 'Private anonymous social network for verified employees: real identity is used only for authentication and company verification, while all activity happens under company-scoped pseudonyms.',
 'Employee communities need honest, verified membership, but tying every post to a real name kills the conversation and makes moderation a deanonymization risk.',
 'Split identity from membership in a NestJS modular monolith (Prisma/PostgreSQL) and Next.js app, scoping content to verified company communities and designing moderation flows that enforce rules without revealing who wrote what.',
 'https://picsum.photos/seed/corpg/800/600',
 ARRAY['NestJS', 'Next.js', 'TypeScript', 'Prisma', 'PostgreSQL', 'Tailwind CSS'],
 '2026', 'https://github.com/wicahma/corporate-underground'),
('77777777-7777-7777-7777-777777777708', 'Monday MCP (unofficial)',
 'Open-source MCP bridge and agent skill suite for Monday.com, covering corporate instances that block enterprise developer admin rights via session token rotation.',
 'Official integrations assume developer admin access, so agents cannot reach corporate Monday.com subdomains that only expose personal access tokens or browser session tokens.',
 'Built a pure-Python stdio JSON-RPC MCP server with a pre-implemented GraphQL client, a token manager CLI that verifies validity and rotates expired session tokens without restarting agents, and 19 agent runbooks covering SDLC governance and board scaffolding.',
 'https://picsum.photos/seed/mondaymcp/800/600',
 ARRAY['Python', 'MCP', 'GraphQL', 'Monday.com', 'JSON-RPC'],
 '2026', 'https://github.com/wicahma/monday-mcp-unofficial'),
('77777777-7777-7777-7777-777777777709', 'diama.dev Portfolio Platform',
 'The portfolio site you are reading: Next.js frontend backed by a Go modular monolith API with PostgreSQL, Redis, and S3-backed media.',
 'Static portfolio sites cannot publish articles, projects, or experience data without a redeploy, and heavyweight frameworks add memory and dependency cost for a small content API.',
 'Split the product into a Next.js App Router frontend and a Go API built on net/http ServeMux, pgx, and sqlc with strict domain isolation, plus Redis caching, S3 image processing, and GitHub OAuth for the admin dashboard.',
 'https://picsum.photos/seed/portofolio/800/600',
 ARRAY['Next.js', 'TypeScript', 'Go', 'PostgreSQL', 'sqlc', 'Redis', 'AWS S3'],
 '2026', 'https://github.com/wicahma/porto');

-- 4. Articles ----------------------------------------------------------------
DELETE FROM articles;

INSERT INTO articles (id, slug, title, content, excerpt, category, image, tags, read_time, meta_title, meta_description, meta_keywords, og_image) VALUES
('88888888-8888-8888-8888-888888888801',
 'http-caching-private-shared-caches-and-freshness',
 'HTTP Caching: Private Caches, Shared Caches, and How Freshness Is Decided',
 '<p>The HTTP cache stores a response associated with a request and reuses the stored response for subsequent requests. The closer the client and the cache are, the faster the response, and when a response is reusable the origin server does not need to parse and route the request, restore the session from a cookie, query the database, or render a template.</p><h2>Private and shared caches</h2><p>The HTTP Caching specification defines two main types: <strong>private caches</strong>, tied to a specific client and typically living in the browser, and <strong>shared caches</strong>, which sit between client and server and store responses that can be shared among users. Shared caches split further into proxy caches and managed caches.</p><p>A private cache may store a personalized response for that user. If personalized content is stored anywhere else, other users may retrieve it, so responses containing personalized content must be marked with the <code>private</code> directive.</p><pre><code>Cache-Control: private</code></pre><p>Personalized content is usually controlled by cookies, but the presence of a cookie does not always mean the response is private.</p><h2>Managed caches</h2><p>Managed caches are deployed explicitly by service developers to offload the origin server: reverse proxies, CDNs, and service workers combined with the Cache API. Their behavior is usually controlled through <code>Cache-Control</code> plus the product own configuration or dashboard. The specification does not define a way to delete a stored response, but a managed cache can purge at any time through dashboard operations, API calls, or restarts.</p><p>Because managed caches can ignore the standard on purpose, avoid kitchen-sink headers such as <code>no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate</code>. Read the documentation of the cache you actually deployed and control it the way that mechanism provides.</p><h2>Heuristic caching</h2><p>HTTP is designed to cache as much as possible, so even without <code>Cache-Control</code> a response can be stored and reused when conditions are met. Content that has not been updated for a full year is assumed to stay unchanged for a while; the specification recommends reusing it for about 10% of the time since it was last modified. Heuristic caching predates wide <code>Cache-Control</code> support, and basically all responses should specify the header explicitly.</p><h2>Fresh, stale, and age</h2><p>Stored responses are either <strong>fresh</strong> or <strong>stale</strong>. The criterion is <em>age</em>: the time elapsed since the response was generated. With <code>Cache-Control: max-age=604800</code>, a response younger than one week is fresh and is used to fulfill requests; older than one week, it is stale.</p><p>A shared cache that has held the response for one day adds an <code>Age</code> header. The client receiving <code>Age: 86400</code> with <code>max-age=604800</code> treats the response as fresh for the remaining 518400 seconds.</p><h2>Expires or max-age</h2><p>HTTP/1.0 used <code>Expires</code> with an explicit time. The format is hard to parse, implementations had bugs, and shifting the system clock could induce problems, so HTTP/1.1 adopted <code>max-age</code> with an elapsed time. When both are present, <code>max-age</code> wins, so there is no need to send <code>Expires</code> today.</p><h2>Vary</h2><p>Responses are normally distinguished by URL, but content negotiation means the body can depend on <code>Accept</code>, <code>Accept-Language</code>, and <code>Accept-Encoding</code>. Adding <code>Vary: Accept-Language</code> keys the cache on URL plus language instead of URL alone. Avoid varying on <code>User-Agent</code>: it has so many variations that the cache will almost never be reused, so prefer feature detection. For cookie-protected personalized content, specify <code>Cache-Control: private</code> rather than listing cookies in <code>Vary</code>.</p><p><em>Adapted from MDN Web Docs, HTTP caching, licensed under CC BY-SA 2.5. Source: developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Caching</em></p>',
 'How private and shared caches differ, why heuristic caching exists, and how age, max-age, and Vary decide whether a stored response can be reused.',
 'Web Performance',
 'https://picsum.photos/seed/art-http-cache/800/600',
 ARRAY['HTTP', 'Caching', 'Performance', 'CDN', 'Web'],
 '8 min read',
 'HTTP Caching Explained: Private vs Shared Caches',
 'A practical walkthrough of HTTP caching semantics: private and shared caches, managed caches, heuristic caching, freshness, age, and Vary.',
 ARRAY['http', 'caching', 'cache-control', 'web performance', 'cdn'],
 'https://picsum.photos/seed/art-http-cache/800/600'),

('88888888-8888-8888-8888-888888888802',
 'go-concurrency-patterns-pipelines-and-cancellation',
 'Go Concurrency Patterns: Pipelines and Cancellation',
 '<p>Go concurrency primitives make it easy to construct streaming data pipelines that use I/O and multiple CPUs efficiently. A pipeline is a series of <em>stages</em> connected by channels, where each stage is a group of goroutines running the same function. In each stage, goroutines receive values from upstream via inbound channels, perform some function on that data, and send values downstream via outbound channels.</p><p>The first stage is the source or producer; the last is the sink or consumer. Each stage has any number of inbound and outbound channels, except the first and last, which have only outbound or inbound channels respectively.</p><h2>Squaring numbers</h2><p>The first stage converts a list of integers into a channel and closes the channel when all values have been sent:</p><pre><code>func gen(nums ...int) &lt;-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out &lt;- n
        }
        close(out)
    }()
    return out
}</code></pre><p>The second stage receives integers and emits their squares, closing its outbound channel once the inbound channel is closed and all values are sent. Because the inbound and outbound types match, stages compose any number of times.</p><h2>Fan-out, fan-in</h2><p>Multiple functions can read from the same channel until it is closed; this is <em>fan-out</em>, and it distributes work across a group of workers to parallelize CPU use and I/O. A function can read from multiple inputs and proceed until all are closed by multiplexing them onto a single channel: <em>fan-in</em>.</p><p>The merge function starts a goroutine per inbound channel that copies values to the single outbound channel. Because sends on a closed channel panic, all sends must finish before closing; <code>sync.WaitGroup</code> arranges that synchronization, and one more goroutine waits for the group and then closes the output.</p><h2>Stopping short</h2><p>Pipeline stages follow a pattern: stages close their outbound channels when all send operations are done, and keep receiving from inbound channels until those are closed. That lets every receiving stage be a <code>range</code> loop and guarantees goroutines exit once values are sent downstream.</p><p>Real pipelines do not always receive all inbound values. A receiver may only need a subset, or a stage exits early because an inbound value represents an error upstream. If a stage stops consuming, goroutines attempting to send block indefinitely. That is a resource leak: goroutines consume memory and runtime resources, and heap references in goroutine stacks keep data from being garbage collected. Goroutines are not garbage collected; they must exit on their own.</p><h2>Buffers are not the fix</h2><p>One way out is to buffer the outbound channel so sends complete immediately while there is room. When the number of values is known at channel creation time, a buffer simplifies code, and <code>gen</code> can copy the list into a buffered channel without a goroutine.</p><p>Adding a buffer to the channel returned by <code>merge</code> fixes the blocked goroutine in one program, but it is bad code: the size of 1 depends on knowing how many values merge will receive and how many downstream stages will consume. Pass one more value to <code>gen</code>, or have the downstream stage read one fewer value, and the goroutines block again.</p><h2>Explicit cancellation</h2><p>What is needed is a way for downstream stages to tell senders they will stop accepting input. When <code>main</code> decides to exit without receiving all values, it sends on a <code>done</code> channel to tell upstream stages to abandon the values they are trying to send. Each stage selects between sending its value downstream and receiving from <code>done</code>, so a closed done channel unblocks every sender and lets the pipeline unwind.</p><p><em>Adapted from Sameer Ajmani, Go Concurrency Patterns: Pipelines and cancellation, The Go Blog, 13 March 2014, licensed under CC BY 3.0. Source: go.dev/blog/pipelines</em></p>',
 'Stages, fan-out and fan-in, why buffers only hide the problem, and how an explicit done channel stops a pipeline cleanly.',
 'Backend',
 'https://picsum.photos/seed/art-go-pipeline/800/600',
 ARRAY['Go', 'Concurrency', 'Channels', 'Goroutines', 'Backend'],
 '11 min read',
 'Go Concurrency Patterns: Pipelines and Cancellation',
 'Build streaming pipelines in Go with channels, fan-out and fan-in, then cancel them safely without leaking goroutines.',
 ARRAY['golang', 'concurrency', 'goroutines', 'channels', 'pipelines'],
 'https://picsum.photos/seed/art-go-pipeline/800/600'),

('88888888-8888-8888-8888-888888888803',
 'postgresql-mvcc-why-reads-never-block-writes',
 'PostgreSQL MVCC: Why Reads Never Block Writes',
 '<p>PostgreSQL maintains data consistency with a multiversion model called Multiversion Concurrency Control (MVCC). Each SQL statement sees a snapshot of data as it was some time ago, regardless of the current state of the underlying rows. That prevents statements from viewing inconsistent data produced by concurrent transactions updating the same rows, and provides transaction isolation for each database session.</p><h2>What MVCC buys you</h2><p>By avoiding the locking methodologies of traditional database systems, MVCC minimizes lock contention so multiuser workloads keep reasonable performance. The main advantage over locking is that locks acquired for reading do not conflict with locks acquired for writing: reading never blocks writing, and writing never blocks reading.</p><p>PostgreSQL keeps that guarantee even at the strictest isolation level, Serializable, through Serializable Snapshot Isolation (SSI).</p><h2>When locks are still the right tool</h2><p>Table- and row-level locking facilities remain available for applications that do not need full transaction isolation and prefer to manage specific points of conflict explicitly. Proper use of MVCC generally performs better than locks. Application-defined advisory locks provide another mechanism for acquiring locks that are not tied to a single transaction.</p><h2>What this means for your schema</h2><p>Because an update writes a new row version instead of overwriting the old one, long-running transactions keep old versions alive. Dead tuples accumulate until vacuum reclaims them, and a transaction that holds a snapshot open for hours can hold back cleanup for the whole database. Monitoring long transactions is therefore part of keeping MVCC healthy, not an optional maintenance task.</p><p><em>Adapted from PostgreSQL 18 Documentation, Chapter 13.1 Introduction (Concurrency Control), PostgreSQL Global Development Group, PostgreSQL License. Source: postgresql.org/docs/current/mvcc-intro.html</em></p>',
 'How PostgreSQL snapshots isolate transactions, why readers never block writers, and what long-running transactions cost.',
 'Database',
 'https://picsum.photos/seed/art-pg-mvcc/800/600',
 ARRAY['PostgreSQL', 'MVCC', 'Transactions', 'Isolation', 'Database'],
 '5 min read',
 'PostgreSQL MVCC Explained: Snapshot Isolation and Locks',
 'Understand PostgreSQL multiversion concurrency control: snapshots, transaction isolation, Serializable Snapshot Isolation, and when explicit locks still matter.',
 ARRAY['postgresql', 'mvcc', 'transactions', 'isolation', 'database'],
 'https://picsum.photos/seed/art-pg-mvcc/800/600'),

('88888888-8888-8888-8888-888888888804',
 'react-server-components-rendering-before-the-bundle',
 'React Server Components: Rendering Before the Bundle',
 '<p>Server Components are a new type of Component that renders ahead of time, before bundling, in an environment separate from your client app or SSR server. That separate environment is the server in React Server Components. Server Components can run once at build time on a CI server, or for each request on a web server.</p><h2>Server Components without a server</h2><p>Because they can run at build time to read the filesystem or fetch static content, a web server is not required. Without Server Components, fetching static data on the client with an Effect is common, and it means users download and parse additional libraries and wait for a second request after the page loads, just to render content that will not change for the lifetime of the page.</p><p>With Server Components, the component can read the file during render at build time. The rendered output is then server-side rendered to HTML and uploaded to a CDN. The client never sees the original component or the expensive markdown libraries; it only sees the rendered output, so the content is visible on first page load and the bundle stays small.</p><p>Note that the Server Component is an async function. Async Components are a Server Components feature that allow <code>await</code> in render.</p><h2>Server Components with a server</h2><p>Server Components can also run on a web server during a request, letting you access your data layer without building an API. Without them, dynamic data is usually fetched on the client in an Effect, and nested components fetching their own data cause an expensive client-server waterfall.</p><p>With Server Components, the component reads the data and renders it directly, and the bundler combines the data, rendered Server Components, and dynamic Client Components into a bundle. Optionally that bundle is SSR-ed into initial HTML. When the page loads, the browser sees only the rendered output, not the original components. Re-fetching from the server makes them dynamic, combining the request/response model of server-centric multi-page apps with the interactivity of client-centric single-page apps.</p><h2>Context and interactivity</h2><p>Server Components cannot create context, but they can render a context provider imported from a Client Component module that exports the context under the <code>use client</code> directive. Client Components inside that provider read the value with <code>use</code> or <code>useContext</code>.</p><p>Server Components are not sent to the browser, so they cannot use interactive APIs like <code>useState</code>. To add interactivity, compose them with Client Components marked <code>use client</code>: the server renders the data, and the client component handles state such as an expand toggle, receiving server output as children.</p><p>A common misunderstanding is that Server Components are denoted by <code>use server</code>. There is no directive for Server Components; <code>use server</code> marks Server Functions.</p><h2>Async components and streaming</h2><p>Server Components introduce async/await in components. When you <code>await</code> in an async component, React suspends and waits for the promise to resolve before resuming rendering. This works across the server/client boundary with streaming support for Suspense, and a promise created on the server can be awaited on the client.</p><p><em>Adapted from React documentation, Server Components, Meta Platforms, licensed under CC BY 4.0. Source: react.dev/reference/rsc/server-components</em></p>',
 'How Server Components render before bundling, why that removes client waterfalls, and how they compose with Client Components for interactivity.',
 'Frontend',
 'https://picsum.photos/seed/art-rsc/800/600',
 ARRAY['React', 'Server Components', 'RSC', 'Next.js', 'Frontend'],
 '9 min read',
 'React Server Components Explained',
 'Server Components render before bundling: no client data waterfalls, smaller bundles, and clean composition with Client Components.',
 ARRAY['react', 'server components', 'rsc', 'nextjs', 'frontend'],
 'https://picsum.photos/seed/art-rsc/800/600'),

('88888888-8888-8888-8888-888888888805',
 'docker-multi-stage-builds-small-production-images',
 'Docker Multi-stage Builds: Small Production Images Without Extra Scripts',
 '<p>Multi-stage builds help anyone who has struggled to optimize Dockerfiles while keeping them easy to read and maintain. You use multiple <code>FROM</code> statements in one Dockerfile: each can use a different base and each begins a new stage of the build. You can selectively copy artifacts from one stage to another, leaving behind everything you do not want in the final image.</p><h2>A build stage and a runtime stage</h2><p>The following Dockerfile has two stages: one that builds a binary, and another that copies the binary out of the first stage.</p><pre><code># syntax=docker/dockerfile:1
FROM golang:1.26 AS build
WORKDIR /src
COPY main.go /src/main.go
RUN go build -o /bin/hello ./main.go

FROM scratch
COPY --from=build /bin/hello /bin/hello
CMD ["/bin/hello"]</code></pre><p>You only need the single Dockerfile and a plain <code>docker build</code>. The result is a tiny production image containing nothing but the binary. None of the build tools required to build the application are included: the second <code>FROM</code> starts a new stage on <code>scratch</code>, and <code>COPY --from=build</code> copies just the built artifact, leaving the Go SDK and intermediate artifacts behind.</p><h2>Name your stages</h2><p>By default stages are unnamed and you refer to them by integer, starting at 0. Adding <code>AS &lt;NAME&gt;</code> to <code>FROM</code> lets you use the name in <code>COPY</code>, so reordering instructions later does not break the copy.</p><h2>Stop at a specific stage</h2><p>You do not have to build every stage. <code>docker build --target build -t hello .</code> stops at the named stage. Useful scenarios: debugging one stage, a <code>debug</code> stage with symbols and tools plus a lean <code>production</code> stage, or a <code>testing</code> stage populated with test data while production builds with real data.</p><h2>External images and stage reuse</h2><p><code>COPY --from</code> can copy from a separate image, using a local image name, a tag from a registry, or an image ID; Docker pulls it if necessary. For example, <code>COPY --from=nginx:latest /etc/nginx/nginx.conf /nginx.conf</code>.</p><p>You can also continue from a previous stage by naming it in <code>FROM</code>, so a shared <code>builder</code> stage can feed several independent build stages.</p><h2>BuildKit only builds what it needs</h2><p>The legacy builder processes all stages leading up to the selected <code>--target</code>, even when the target does not depend on them. BuildKit builds only the stages the target depends on, so an unrelated sibling stage is skipped entirely, which makes targeted builds noticeably faster.</p><p><em>Adapted from Docker Docs, Multi-stage builds, Docker Inc., licensed under Apache 2.0. Source: docs.docker.com/build/building/multi-stage/</em></p>',
 'Why multi-stage Dockerfiles beat separate build scripts, how to name and target stages, and why BuildKit skips work the legacy builder does not.',
 'DevOps',
 'https://picsum.photos/seed/art-docker/800/600',
 ARRAY['Docker', 'BuildKit', 'Containers', 'DevOps', 'CI/CD'],
 '6 min read',
 'Docker Multi-stage Builds Explained',
 'Build small production images with multi-stage Dockerfiles: named stages, build targets, external images, and BuildKit dependency pruning.',
 ARRAY['docker', 'multi-stage build', 'buildkit', 'containers', 'devops'],
 'https://picsum.photos/seed/art-docker/800/600');
