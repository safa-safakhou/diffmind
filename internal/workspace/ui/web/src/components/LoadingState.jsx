import './LoadingState.css'

// Unknown-duration requests use an indeterminate indicator, not a guessed percentage.
export function LoadingState({ label = 'Loading…', detail, compact = false, overlay = false }) {
  return <div class={`loading-state${compact ? ' compact' : ''}${overlay ? ' overlay' : ''}`} role="status" aria-live="polite">
    <span class="loading-state-spinner" aria-hidden="true" />
    <div><strong>{label}</strong>{detail && <p>{detail}</p>}</div>
  </div>
}
