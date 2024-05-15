export const backend_url =
  window.location.origin === 'http://localhost:3000'
    ? `http://localhost:${import.meta.env.VITE_BACKEND_PORT}`
    : window.location.origin
