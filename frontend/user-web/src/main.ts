import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'
import '@fontsource-variable/lexend';

const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
