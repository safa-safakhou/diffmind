export function AgentConnectionHelp({ pid, role, endpoint }) {
 return <section class="banner" aria-label="Connect an agent">
   <h2>Connect an agent to this workspace</h2>
   <p>Use the shared server’s MCP endpoint: <code>{endpoint || '/mcp'}</code>. Select project <code>{pid}</code> explicitly when more than one project is accessible.</p>
   <p>Your browser login does not configure your agent’s identity. Ask your workspace administrator for a viewer project token or the supported identity-proxy connection. A viewer connection discovers the read tools, including readiness, and queries saved context.</p>
   <p>Your current project role is {role || 'being checked'}. Viewer agents cannot import or update context; ask an editor to refresh existing sources. A prepared company workspace needs no local ingestion setup.</p>
   <p>Store the one-time secret privately in your host’s secret configuration. Use HTTPS for shared connections. Your workspace administrator owns token expiry, renewal and support: obtain a replacement, verify the new connection, then revoke the old token. Keep admin recovery credentials and proxy secrets private.</p>
 </section>
}
