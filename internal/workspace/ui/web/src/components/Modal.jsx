import { useLayoutEffect, useId, useRef, useState } from 'preact/hooks'

// One keyboard and focus boundary for every dialog, including confirmations.
const dialogs = []
export function Modal({ title, onClose, children, wide, initialFocus }) {
  const ref = useRef(null)
  const close = useRef(onClose)
  close.current = onClose
  const titleID = useId()
  useLayoutEffect(() => {
    const node = ref.current
    const previous = document.activeElement
    dialogs.push(node)
    const top = () => dialogs.at(-1) === node
    const focusable = () => [...node.querySelectorAll('button, input, textarea, select, a[href], [tabindex]')]
      .filter((el) => !el.disabled && el.tabIndex !== -1 && !el.closest('[hidden], [inert]'))
    const focusFirst = () => (node.querySelector(initialFocus || '[data-dialog-focus]') || focusable().find((el) => el.tagName !== 'BUTTON') || focusable()[0] || node).focus()
    const inert = []
    for (let parent = node.parentElement; parent && parent !== document.body; parent = parent.parentElement) {
      for (const sibling of parent.parentElement?.children || []) {
        if (sibling !== parent && !sibling.contains(node)) {
          inert.push([sibling, sibling.hasAttribute('inert')])
          sibling.setAttribute('inert', '')
        }
      }
    }
    focusFirst()
    const keydown = (e) => {
      if (!top()) return
      if (e.key === 'Escape') { e.preventDefault(); e.stopImmediatePropagation(); close.current?.(); return }
      if (e.key !== 'Tab') return
      const elements = focusable(), first = elements[0], last = elements.at(-1)
      if (!first) { e.preventDefault(); node.focus(); return }
      if (e.shiftKey && (document.activeElement === first || !node.contains(document.activeElement))) { e.preventDefault(); last.focus() }
      else if (!e.shiftKey && (document.activeElement === last || !node.contains(document.activeElement))) { e.preventDefault(); first.focus() }
    }
    const focusin = (e) => { if (top() && !node.contains(e.target)) focusFirst() }
    document.addEventListener('keydown', keydown, true)
    document.addEventListener('focusin', focusin)
    return () => {
      dialogs.splice(dialogs.indexOf(node), 1)
      document.removeEventListener('keydown', keydown, true)
      document.removeEventListener('focusin', focusin)
      for (const [el, wasInert] of inert) if (!wasInert) el.removeAttribute('inert')
      if (previous?.isConnected && !previous.closest('[inert]')) previous.focus()
    }
  }, [])
  return <div class="modal-backdrop" onClick={() => close.current?.()}>
    <div ref={ref} class={'modal' + (wide ? ' wide' : '')} role="dialog" aria-modal="true" aria-labelledby={titleID} tabIndex="-1" onClick={(e) => e.stopPropagation()}>
      <div class="modal-head"><h2 id={titleID}>{title}</h2><button class="btn ghost tiny" aria-label="Close dialog" onClick={onClose}>✕</button></div>
      <div class="modal-body">{children}</div>
    </div>
  </div>
}

export function ConfirmDialog({ title, message, confirmLabel = 'Delete', onConfirm, onCancel }) {
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const submitting = useRef(false)
  const confirm = async (event) => {
    event.preventDefault(); event.stopPropagation()
    if (submitting.current) return
    submitting.current = true
    setBusy(true); setError('')
    try { await onConfirm() } catch (e) { setError(e.message) }
    finally { submitting.current = false; setBusy(false) }
  }
  return <Modal title={title} onClose={onCancel} initialFocus="[data-dialog-focus]">
    <p class="confirm-message">{message}</p>
    {error && <p class="banner error" role="alert">{error}</p>}
    {busy && <p role="status">Submitting… Check persisted work if the connection is lost.</p>}
    <div class="actions">
      <button class="btn danger" disabled={busy} onClick={confirm}>{confirmLabel}</button>
      <button class="btn ghost" data-dialog-focus onClick={onCancel}>Cancel</button>
    </div>
  </Modal>
}
