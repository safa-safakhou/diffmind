import test from 'node:test'
import assert from 'node:assert/strict'
import { createProject, ApiError } from './api.js'

test('lost write responses retain safe recovery without replaying the request', async (t) => {
 const original=globalThis.fetch;t.after(()=>{globalThis.fetch=original});let calls=0
 globalThis.fetch=async()=>{calls++;throw new Error('private transport details')}
 await assert.rejects(createProject({name:'approved'}),error=>error instanceof ApiError && error.status===0 && error.recovery.retryable===false && !error.message.includes('private'))
 assert.equal(calls,1)
})

test('conflict feedback preserves recovery and validation without automatic write retries', async (t) => {
 const original=globalThis.fetch;t.after(()=>{globalThis.fetch=original});let calls=0
 globalThis.fetch=async()=>{calls++;return {ok:false,status:409,text:async()=>JSON.stringify({error:'Scope changed',validation:[{field:'root'}],recovery:{category:'conflict',retryable:false,next_action:'Preview current scope again.'}})}}
 await assert.rejects(createProject({name:'approved'}),error=>error.status===409 && error.validation[0].field==='root' && error.recovery.category==='conflict' && /Preview current scope/.test(error.message))
 assert.equal(calls,1)
})
