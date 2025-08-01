import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

// Debug log
console.log('Main script executing');

const rootElement = document.getElementById('root');

if (rootElement) {
  console.log('Root element found, rendering app');
  createRoot(rootElement).render(
    <StrictMode>
      <App />
    </StrictMode>,
  );
} else {
  console.error('Root element not found!');
}
