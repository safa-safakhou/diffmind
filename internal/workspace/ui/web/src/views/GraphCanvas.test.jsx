import test from 'node:test'
import assert from 'node:assert/strict'
import dagre from 'dagre'
import { layoutGraph } from './GraphCanvas.jsx'

function fixture(count, resourceCount = 0) {
  const top = new dagre.graphlib.Graph({ multigraph: true })
  top.setGraph({ rankdir: 'LR', nodesep: 100, ranksep: 190, marginx: 90, marginy: 90 })
  top.setDefaultEdgeLabel(() => ({}))
  const services = Array.from({ length: count }, (_, i) => ({ name: `service-${i}`, team: `team-${i % 8}` }))
  const models = new Map()
  const edges = []
  services.forEach((service, i) => {
    top.setNode(service.name, { type: 'service', width: 300, height: 116, expanded: false })
    models.set(service.name, { width: 1868, height: i === 0 ? 12000 : 1400 })
  })
  for (let i = 0; i < resourceCount; i++) {
    const id = `resource-${i}`
    top.setNode(id, { type: i % 2 ? 'db' : 'scheduler', width: 320, height: 100 })
    const edge = { from: services[i % count].name, to: id, type: i % 2 ? 'database' : 'scheduler' }
    edges.push(edge)
    top.setEdge(edge.from, edge.to, {})
  }
  return { top, services, models, edges }
}

function assertSeparated(top) {
  const nodes = top.nodes().map((id) => ({ id, ...top.node(id) }))
  for (let i = 0; i < nodes.length; i++) {
    const a = nodes[i]
    assert.ok(Number.isFinite(a.x) && Number.isFinite(a.y), a.id)
    for (const b of nodes.slice(i + 1)) {
      const overlap = Math.abs(a.x - b.x) < (a.width + b.width) / 2 && Math.abs(a.y - b.y) < (a.height + b.height) / 2
      assert.equal(overlap, false, `${a.id} overlaps ${b.id}`)
    }
  }
}

test('29-service overview clusters by team and separates resource attachments', () => {
  const { top, services, models, edges } = fixture(29, 116)
  layoutGraph(top, services, edges, models)
  assertSeparated(top)
  assert.equal(top.node('resource-0').teamFrame, 'team-0')
})

test('full detail arranges actual card dimensions, including unusually tall services', () => {
  for (const count of [2, 29]) {
    const { top, services, models, edges } = fixture(count, 30)
    layoutGraph(top, services, edges, models, { mode: 'detail' })
    assert.equal(top.node('service-0').height, 12000)
    assertSeparated(top)
  }
})

test('dense resource stacks remain separated beyond 80 attachments', () => {
  const { top, services, models, edges } = fixture(1, 200)
  layoutGraph(top, services, edges, models, { clustered: true })
  assertSeparated(top)
})

test('invalid or overlapping persisted positions are laid out again', () => {
  for (const x of [0, NaN, Infinity]) {
    const { top, services, models, edges } = fixture(2)
    const positions = new Map(services.map((s) => [s.name, { x, y: 0, width: 300, height: 116 }]))
    layoutGraph(top, services, edges, models, { positions })
    assertSeparated(top)
  }
})

test('valid saved positions are preserved', () => {
  const { top, services, models, edges } = fixture(2)
  const positions = new Map(services.map((s, i) => [s.name, { x: 300 + i * 800, y: 300, width: 300, height: 116 }]))
  layoutGraph(top, services, edges, models, { positions })
  services.forEach((s) => assert.equal(top.node(s.name).x, positions.get(s.name).x))
})
