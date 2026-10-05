import { useLayoutEffect, useId, useRef } from 'preact/hooks'
import { Button } from './Button.jsx'

// Modal is the shared dialog frame. size: normal|wide. Extracted from the
// inline Modal/EditorModal definitions so every dialog looks the same.
const dialogs = []
export function Modal({ title, onClose, size = 'normal', class: cls = '', children, footer, initialFocus }) {
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

  const classes = ['modal']
  if (size === 'wide') classes.push('wide')
  if (cls) classes.push(cls)
  return (
    <div class="modal-backdrop" onClick={onClose}>
      <div ref={ref} role="dialog" aria-modal="true" aria-labelledby={titleID} tabIndex="-1" class={classes.join(' ')} onClick={(e) => e.stopPropagation()}>
        <div class="modal-head">
          <h2 id={titleID}>{title}</h2>
          {onClose && <Button aria-label="Close dialog" variant="secondary" size="tiny" onClick={onClose}>✕</Button>}
        </div>
        <div class="modal-body">{children}</div>
        {footer && <div class="modal-foot">{footer}</div>}
      </div>
    </div>
  )
}

export function ConfirmDialog({ title, message, confirmLabel = 'Confirm', danger = true, onConfirm, onCancel }) {
  return (
    <Modal title={title} onClose={onCancel} class="confirm">
      <p class="confirm-message">{message}</p>
      <div class="actions">
        <Button variant={danger ? 'danger' : 'primary'} onClick={onConfirm}>{confirmLabel}</Button>
        <Button variant="secondary" onClick={onCancel}>Cancel</Button>
      </div>
    </Modal>
  )
}
