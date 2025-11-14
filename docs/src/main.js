import './style.css'
import Prism from 'prismjs'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-yaml'
import 'prismjs/themes/prism.css'

// Prism.jsの自動ハイライトを初期化
document.addEventListener('DOMContentLoaded', () => {
  Prism.highlightAll()
})
